package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

func GetBalanceAcc(accNum string, db *sql.DB) (acc models.Accounts, err error) {
	err = db.QueryRow("SELECT account_balance FROM accounts WHERE account_num = ($1) ", accNum).Scan(&acc.AccountBalance)
	return
}
