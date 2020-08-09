package models

import (
	"database/sql"
	"time"
)

type Admins struct {
	AdminId   int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func LoginAdmin(db *sql.DB, username string, password string) (objAdmin Admins, err error) {
	err = db.QueryRow("SELECT id, username, created_at from admins where username = ($1) and password = ($2)", username, password).Scan(&objAdmin.AdminId, &objAdmin.Username, &objAdmin.CreatedAt)
	return
}
