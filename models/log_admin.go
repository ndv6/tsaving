package models

import "time"

type LogAdmin struct {
	IDLogAdmin    string    `json:"admin_username"`
	AdminUsername string    `json:"admin_password"`
	CreatedAt     time.Time `json:"created_at"`
}
