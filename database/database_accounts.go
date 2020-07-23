package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ndv6/tsaving/models"
)

func GetBalanceAcc(accNum string, db *sql.DB) (acc models.Accounts, err error) {
	err = db.QueryRow("SELECT account_balance FROM accounts WHERE account_num = ($1) ", accNum).Scan(&acc.AccountBalance)
	return
}

func TransferFromMainToVa(accNum, vaNum string, amount int, db *sql.DB) (err error) {
	//use tx to make sure all queries below success before tx.Commit()
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var sourceBalance int
	err = tx.QueryRow("SELECT account_balance FROM accounts WHERE account_num = $1 FOR UPDATE", accNum).Scan(&sourceBalance)
	if err != nil {
		return
	}

	status := CheckBalance("MAIN", accNum, amount, db)
	if !status {
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
	_, err = tx.Exec("UPDATE virtual_accounts SET va_balance = va_balance + $1 WHERE va_num = $2", amount, vaNum)
	if err != nil {
		return
	}

	logDesc := models.LogDescriptionMainToVaTemplate(amount, accNum, vaNum)
	logData := models.TransactionLogs{
		AccountNum:  accNum,
		DestAccount: vaNum,
		TranAmount:  amount,
		Description: logDesc,
		CreatedAt:   time.Now(),
	}

	err = models.TransactionLog(tx, logData)
	if err != nil {
		return
	}

	tx.Commit()

	return
}
