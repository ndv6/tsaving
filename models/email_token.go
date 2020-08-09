package models

import (
	"database/sql"
)

type EmailToken struct {
	EtId  int    `json:"et_id"`
	Token string `json:"token"`
	Email string `json:"email"`
}

type TokenHandler struct {
	db *sql.DB
}

func NewTokenHandler(db *sql.DB) *TokenHandler {
	return &TokenHandler{
		db,
	}
}

type GetTokenRequest struct {
	Email string `json:"email"`
}

func AddEmailTokens(db *sql.DB, Token string, Email string) error {
	_, err := db.Exec("INSERT into email_token(token, email) values ($1, $2)", Token, Email)
	return err
}

func (t *TokenHandler) UpsertEmailToken(token string, email string) error {
	_, err := t.db.Exec("INSERT INTO email_token(token, email) VALUES ($1, $2) ON CONFLICT (email) DO UPDATE SET token = ($1)", token, email)
	return err
}

type VerifiedEmailResponse struct {
	Email string `json:"email"`
}
