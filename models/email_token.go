package models

import (
	"database/sql"
)

type EmailToken struct {
	EtId 			int 		`json:"et_id"`
	Token 			string		`json:"token"`
 	Email 			string 		`json:"email"`
}

func AddEmailTokens(db *sql.DB, Token string, Email string) error {
	_, err := db.Exec("INSERT into email_token(token, email) values ($1, $2)", Token, Email,)
		return err
}

type VerifiedEmailResponse struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}
