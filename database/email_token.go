package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

func DeleteVerifiedEmailTokenById(id int, db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM email_token WHERE et_id=$1;", id)
	return
}

func UpdateCustomerVerificationStatusByEmail(email string, db *sql.DB) (err error) {
	_, err = db.Exec("UPDATE CUSTOMERS SET is_verified = TRUE WHERE cust_email = $1;", email)
	return
}

func GetEmailTokenByTokenAndEmail(db *sql.DB, token, email string) (et models.EmailToken, err error) {
	err = db.QueryRow("SELECT et_id, token, email FROM email_token WHERE token=$1 AND email=$2", token, email).Scan(&et.EtId, &et.Token, &et.Email)
	return
}
