package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

// This email_token.go file is made by Joseph

type EmailHandler struct {
	Db *sql.DB
}

func NewEmailHandler(db *sql.DB) EmailHandler {
	return EmailHandler{Db: db}
}

func (eh EmailHandler) DeleteVerifiedEmailTokenById(id int) (err error) {
	_, err = eh.Db.Exec("DELETE FROM email_token WHERE et_id=$1;", id)
	return
}

func (eh EmailHandler) UpdateCustomerVerificationStatusByEmail(email string) (err error) {
	_, err = eh.Db.Exec("UPDATE CUSTOMERS SET is_verified = TRUE WHERE cust_email = $1;", email)
	return
}

func (eh EmailHandler) GetEmailTokenByEmail(email string) (et models.EmailToken, err error) {
	err = eh.Db.QueryRow("SELECT et_id, token, email FROM email_token WHERE email=$1", email).Scan(&et.EtId, &et.Token, &et.Email)
	return
}
