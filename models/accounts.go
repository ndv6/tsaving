package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Accounts struct {
	AccountId      int       `json:"account_id"`
	AccountNum     string    `json:"account_num"`
	AccountBalance int       `json:"account_balance"`
	CreatedAt      time.Time `json:"created_at"`
}

func GetMainAccount(db *sql.DB, accNum string) (Accounts, error) {
	var acc Accounts
	row := db.QueryRow("SELECT account_num FROM accounts WHERE account_num = $1", accNum)
	err := row.Scan(&acc.AccountNum)
	if err != nil {
		fmt.Println(err)
		return Accounts{}, err
	}
	return acc, nil
}
