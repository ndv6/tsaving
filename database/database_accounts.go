package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

func GetBalanceAcc(accNum string, db *sql.DB) (acc models.Accounts, err error) {
	var balanceFloat float32
	err = db.QueryRow("SELECT account_balance FROM accounts WHERE account_num = ($1) ", accNum).Scan(&balanceFloat)
	if err != nil {
		return
	}
	acc.AccountBalance = int(balanceFloat)
	return
}
