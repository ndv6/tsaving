package models

import (
	"database/sql"
	"time"
)

type Accounts struct {
	AccountId      int       `json:"account_id"`
	AccountNum     string    `json:"account_num"`
	AccountBalance int       `json:"account_balance"`
	CreatedAt      time.Time `json:"created_at"`
}

func GetBalanceAcc(accNum string, db *sql.DB) (acc Accounts, err error) {
	var balanceFloat float32
	err = db.QueryRow("SELECT account_balance FROM accounts WHERE account_num = ($1) ", accNum).Scan(&balanceFloat)
	if err != nil {
		return
	}
	acc.AccountBalance = int(balanceFloat)
	return
}
