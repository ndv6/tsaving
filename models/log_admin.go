package models

import "time"

type LogAdmin struct {
	IDLogAdmin  string    `json:"id_logadmin"`
	Username 	string    `json:"username"`
	AccNum		string	  `json:"acc_num"`
	Action 		string 	  `json:"action"`
	ActionTime   time.Time `json:"action_time"`
}
