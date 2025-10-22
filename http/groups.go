package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	fbErrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/groups"
	"github.com/gorilla/mux"
)


type modifyGroupRequest struct {
	modifyRequest
	Data *groups.Group `json:"data"`
}


func getGroup(_ http.ResponseWriter, r *http.Request ) (*modifyGroupRequest, error) {
	if r.Body == nil {
		return nil, fbErrors.ErrEmptyRequest
	}

	req := &modifyGroupRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		return nil, err
	}

	if req.What != "groups" {
		return nil, fbErrors.ErrInvalidDataType
	}

	return req, nil
}


var groupPostHandler = withAdmin( func( w http.ResponseWriter, r *http.Request, d *data ) (int, error) {

	req, err := getGroup(w, r)

	if err != nil {
		return http.StatusBadRequest, err
	}

	err = rulesValidate(req.Data.Rules)
	if err != nil {

		return http.StatusBadRequest, err
	}

	err = d.store.Groups.SaveGroup( req.Data )

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
} )

var groupsGetHandler = withAdmin( func( w http.ResponseWriter, r *http.Request, d *data ) (int, error) {

	groups, err := d.store.Groups.GetAll()
	if err != nil {

		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, groups)
})

var groupsPutHandler = withAdmin( func (w http.ResponseWriter, r *http.Request, d *data) (int, error) {

	req, err := getGroup(w, r)
	
	err = rulesValidate(req.Data.Rules) 
	if err != nil {

		return http.StatusBadRequest, err
	}

	if err != nil {

		return http.StatusBadRequest, err
	}

	err = d.store.Groups.UpdateGroup( req.Data )

	return http.StatusAccepted, nil
})

var groupsDeleteHandler = withAdmin( func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	
	req_parameters := mux.Vars(r)

	num, err := strconv.Atoi(req_parameters["id"])
	
	if err != nil {

		return http.StatusBadRequest, err
	}


	err = d.store.Groups.Delete( num )

	if err != nil {

		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
} )