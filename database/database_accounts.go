package database

import (
	"database/sql"
)

func GetBalanceAcc(accNum string, db *sql.DB) (balance int, err error) {
	err = db.QueryRow("SELECT account_balance FROM accounts WHERE account_num = ($1) ", accNum).Scan(&balance)
	return
}
