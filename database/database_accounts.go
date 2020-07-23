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

func TransferFromMainToVa(accNum, vaNum string, amount int, db *sql.DB) (err error) {
	//use tx to make sure all queries below success before tx.Commit()
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var sourceBalance float32
	err = tx.QueryRow("SELECT account_balance FROM accounts WHERE account_num = $1 FOR UPDATE", accNum).Scan(&sourceBalance)
	if err != nil {
		return
	}
	//check balance here
	var sourceBalanceVA float32
	err = tx.QueryRow("SELECT va_balance FROM virtual_accounts WHERE va_num = $1 FOR UPDATE", vaNum).Scan(&sourceBalanceVA)
	if err != nil {
		return
	}

	_, err = tx.Exec("UPDATE accounts SET account_balance = account_balance - $1 WHERE account_num = $2", amount, accNum)
	if err != nil {
		return
	}
	_, err = tx.Exec("UPDATE virtual_accounts SET va_balance = va_balance + $1 WHERE va_num = $2", amount, vaNum)
	if err != nil {
		return
	}
	tx.Commit()
	return

}
