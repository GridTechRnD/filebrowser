package http

import (
	"log"
	"net/http"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/filebrowser/filebrowser/v2/runner"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage"
	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/tomasen/realip"
)

type handleFunc func(w http.ResponseWriter, r *http.Request, d *data) (int, error)

type data struct {
	*runner.Runner
	settings *settings.Settings
	server   *settings.Server
	store    *storage.Storage
	user     *users.User
	raw      interface{}
}

type PathMeta struct {
	FullPath	[]string
	Parent		[]string
	Object		string
}

func GetPathMeta( path string ) PathMeta {

	splitted_path := strings.Split(path, "/")[ 1 : ]

	return PathMeta{
		FullPath: splitted_path,
		Parent: splitted_path[ : len(splitted_path) - 1],
		Object: splitted_path[ len(splitted_path) - 1 ],
	}

}

func AppendRules( a []rules.Rule, b []rules.Rule ) []rules.Rule {

	if len(a) == 0 {

		return b
	}

	for _, rule := range b {

		contains := false

		for _, a_rule := range a {
			
			if rule.Path == a_rule.Path {

				contains = true
			}
		}

		if !contains {
			a = append(a, rule)
		}
	}

	return a
}

// Check implements rules.Checker.
func (d *data) Check(path string) bool {
	if d.user.HideDotfiles && rules.MatchHidden(path) {
		return false
	}

	if d.user.Perm.Admin {
		return true
	}

	if path[ len(path) - 1 ] == '/' {
		return true
	}

	unique_rules := d.settings.Rules

	// Append user rules
	unique_rules = AppendRules( unique_rules, d.user.Rules )

	// Append group rules
	groups, err := d.store.Groups.GetAll()
	if err != nil {
		
		return false
	}

	for _, group := range groups {
		
		if slices.Contains(group.UsersIds, d.user.ID) {
			
			unique_rules = AppendRules( unique_rules, group.Rules )
		}
	}

	sort.Slice(unique_rules, func(i, j int) bool {
		return len(strings.Split(unique_rules[i].Path, "/")) > len(strings.Split(unique_rules[j].Path, "/"))
	})
	
	log.Println(unique_rules)

	path_meta := GetPathMeta(path)

	allow_rules := []string{}
	deny_rules := []string{}
	
	for _, rule := range(unique_rules) {
		rule_meta := GetPathMeta(rule.Path)
		
		if len(path_meta.Parent) < len(rule_meta.Parent) {
			
			if path_meta.Object == rule_meta.FullPath[ len(path_meta.Parent) ] && !slices.Contains(allow_rules, path_meta.Object) {

				allow_rules = append(allow_rules, path_meta.Object)
			}
		} else if len(path_meta.Parent) == len(rule_meta.Parent) && reflect.DeepEqual(path_meta.Parent, rule_meta.Parent) {
			
			if rule.Allow && !slices.Contains(allow_rules, rule_meta.Object) {
	
				allow_rules = append(allow_rules, rule_meta.Object)
			}

			if !rule.Allow && !slices.Contains(deny_rules, rule_meta.Object) {
					
				deny_rules = append(deny_rules, rule_meta.Object)
			}
			
		} else if len(path_meta.Parent) > len(rule_meta.Parent) {
			
			if rule.Allow {

				if reflect.DeepEqual(path_meta.Parent[ : len(rule_meta.FullPath) ], rule_meta.FullPath) {
					
					return true
				}
			}

			if !rule.Allow {

				if reflect.DeepEqual(path_meta.Parent[ : len(rule_meta.FullPath) - 1 ], rule_meta.Parent) && path_meta.Parent[ len(rule_meta.FullPath) - 1 ] != rule_meta.Object {
					
					return true
				}
			}
		}
	}
	
	allow := false
	
	if len(allow_rules) >= 1 {
		
		allow = slices.Contains(allow_rules, path_meta.Object)
	}

	if len(deny_rules) >= 1 {
		
		allow =  !slices.Contains(deny_rules, path_meta.Object)
	}
	
	return allow
}

func handle(fn handleFunc, prefix string, store *storage.Storage, server *settings.Server) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range globalHeaders {
			w.Header().Set(k, v)
		}

		settings, err := store.Settings.Get()
		if err != nil {
			log.Fatalf("ERROR: couldn't get settings: %v\n", err)
			return
		}
		
		status, err := fn(w, r, &data{
			Runner:   &runner.Runner{Enabled: server.EnableExec, Settings: settings},
			store:    store,
			settings: settings,
			server:   server,
		})

		if status >= 400 || err != nil {
			clientIP := realip.FromRequest(r)
			log.Printf("%s: %v %s %v", r.URL.Path, status, clientIP, err)
		}

		if status != 0 {
			txt := http.StatusText(status)
			if status == http.StatusBadRequest && err != nil {
				txt += " (" + err.Error() + ")"
			}
			http.Error(w, strconv.Itoa(status)+" "+txt, status)
			return
		}
	})

	return stripPrefix(prefix, handler)
}
