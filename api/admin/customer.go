package admin

import (
	"database/sql"
)

type CustomerHandler struct {
	db *sql.DB
}

func NewCustomerHandler(db *sql.DB) *LogAdminHandler {
	return &LogAdminHandler{db}
}
