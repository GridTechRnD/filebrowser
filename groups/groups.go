package groups

import "github.com/filebrowser/filebrowser/v2/rules"

type Group struct {
	ID 			uint	`storm:"id,increment" json:"id"`
	GroupName	string	`storm:"unique" json:"groupName"`
	UsersIds 	[]uint	`json:"usersIds"`
	Rules		[]rules.Rule `json:"groupRules"`
}

