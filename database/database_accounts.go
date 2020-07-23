package database

import (
	"database/sql"
	"errors"

	"github.com/ndv6/tsaving/helpers"
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
	sourceBalanceInt := int(sourceBalance)

	status := helpers.CheckBalance("MAIN", sourceBalanceInt, vac.VaBalance, va.db)
	//check balance here
	if sourceBalanceInt < amount || amount <= 0 {
		err = errors.New("insufficient balance")
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
