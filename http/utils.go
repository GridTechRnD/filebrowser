package http

import (
	"encoding/json"
	"errors"
	"reflect"

	"net/http"
	"net/url"
	"os"
	"strings"

	libErrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/rules"
)

func renderJSON(w http.ResponseWriter, _ *http.Request, data interface{}) (int, error) {
	marsh, err := json.Marshal(data)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(marsh); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func errToStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case os.IsPermission(err):
		return http.StatusForbidden
	case os.IsNotExist(err), errors.Is(err, libErrors.ErrNotExist):
		return http.StatusNotFound
	case os.IsExist(err), errors.Is(err, libErrors.ErrExist):
		return http.StatusConflict
	case errors.Is(err, libErrors.ErrPermissionDenied):
		return http.StatusForbidden
	case errors.Is(err, libErrors.ErrInvalidRequestParams):
		return http.StatusBadRequest
	case errors.Is(err, libErrors.ErrRootUserDeletion):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// This is an adaptation if http.StripPrefix in which we don't
// return 404 if the page doesn't have the needed prefix.
func stripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" || prefix == "/" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		h.ServeHTTP(w, r2)
	})
}

func rulesValidate(rulesList []rules.Rule) error {
	
	//Compare each rules with others
	for r_i, r := range rulesList {
		
		if !r.Regex {
			
			if r.Path[ 0 ] != '/' {

				return errors.New("[" + r.Path + "]: rules must start with /")
			}

			if len(r.Path) == 0 || len( strings.Split(r.Path, "/") ) == 0 {

				return errors.New("empty rule")
			}

			r_meta := GetPathMeta(r.Path)
			
			if len(r_meta.FullPath) == 0 {
				return errors.New(r.Path + " invalid rule")
			}
			
			for sr_i, sr := range rulesList {

				if !sr.Regex {

					if r_i != sr_i {
	
						if len(sr.Path) == 0 || len( strings.Split(sr.Path, "/") ) == 0 {
	
							return errors.New("empty rule")
						}
	
						sr_meta := GetPathMeta(sr.Path)
		
						if reflect.DeepEqual(r_meta.FullPath, sr_meta.FullPath) && r_i != sr_i {
							return errors.New(r.Path + "and" + sr.Path + " is duplicated rule")
						}
		
						if len(r_meta.FullPath) <= len(sr_meta.FullPath) {
							
							if reflect.DeepEqual(r_meta.FullPath, sr_meta.FullPath[ : len(r_meta.FullPath) ]) && !r.Allow {
								return errors.New(r.Path + "and" + sr.Path + " have conflicted rules")
							}
		
						}
	
						if reflect.DeepEqual(r_meta.Parent, sr_meta.Parent) && r.Allow != sr.Allow {
							return errors.New(r.Path + "and" + sr.Path + " have conflicted rules")
						}
					}
				}
			}
		}
	}

	return nil

}