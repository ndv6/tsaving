package database

import (
	"database/sql"
)

type PartnerHandler struct {
	db *sql.DB
}

func NewPartnerHandler(db *sql.DB) *PartnerHandler {
	return &PartnerHandler{
		db,
	}
}

func (ph *PartnerHandler) GetSecret(clientId int) (clientSecret string, err error) {
	err = ph.db.QueryRow("SELECT secret FROM partners WHERE client_id = ($1)", clientId).Scan(&clientSecret)
	return
}
