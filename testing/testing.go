package testing

import (
	"database/sql"
	"log"
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
	row := db.QueryRow("SELECT account_id, account_num, account_balance, created_at FROM accounts WHERE account_num = $1", accNum)
	var balance float32
	err := row.Scan(&acc.AccountId, &acc.AccountNum, &balance, &acc.CreatedAt)
	acc.AccountBalance = int(balance)
	if err != nil {
		log.Println(err)
		return Accounts{}, err
	}
	return acc, nil
}

func AddAccountsWhileRegister(db *sql.DB, AccNum string) error {
	Create := time.Now()
	Ammount := 0
	_, err := db.Exec("INSERT into accounts(account_num, account_balance, created_at) values ($1, $2, $3)", AccNum,
		Ammount,
		Create,
	)
	return err
}
