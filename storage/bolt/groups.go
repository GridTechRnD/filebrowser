package bolt

import (
	"reflect"
	"errors"

	"github.com/asdine/storm/v3"

	fbErrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/groups"
)


type groupsBackend struct {
	db *storm.DB
}


func (st groupsBackend) Create( group *groups.Group ) error {
	
	err := st.db.Save( group )

	if errors.Is(err, storm.ErrAlreadyExists) {
		return fbErrors.ErrExist
	}

	return err

}

func (st groupsBackend) Update( group *groups.Group ) error {

	v := reflect.ValueOf(group).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {

		fieldName := t.Field(i).Name

		if fieldName != "ID" {

			fieldValue := v.Field(i).Interface()
			if err := st.db.UpdateField(group, fieldName, fieldValue) ; err != nil {
				return err
			}
		} 
	}

	return nil
}

func (st groupsBackend) GetAllGroups() ([]*groups.Group, error) {

	var allGroups []*groups.Group

	err := st.db.All(&allGroups)

	if errors.Is(err, storm.ErrNotFound){

		return nil, fbErrors.ErrNotExist
	}

	if err != nil {

		return allGroups, err
	}

	return allGroups, err
}

func (st groupsBackend) DeleteByID( groupId int ) error {
	err := st.db.DeleteStruct( &groups.Group{ ID: uint(groupId) } )

	if err != nil {
		return err
	}

	return nil
}