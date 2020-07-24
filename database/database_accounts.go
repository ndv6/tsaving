package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

type AccountHandler struct {
	db *sql.DB
}

func NewAccountHandler(db *sql.DB) *AccountHandler {
	return &AccountHandler{
		db,
	}
}

func GetBalanceAcc(accNum string, db *sql.DB) (balance int, err error) {
	err = db.QueryRow("SELECT account_balance FROM accounts WHERE account_num = ($1) ", accNum).Scan(&balance)
	return
}

func (ah *AccountHandler) AddBalanceToMainAccount(balanceToAdd int, accountNumber string) (err error) {
	_, err = ah.db.Exec("UPDATE accounts SET account_balance = account_balance + ($1) WHERE account_num = ($2)", balanceToAdd, accountNumber)
	return
}

func (ah *AccountHandler) LogTransaction(log models.TransactionLogs) error {
	_, err := ah.db.Exec("INSERT INTO transaction_logs (account_num, dest_account, tran_amount, description, created_at) VALUES ($1, $2, $3, $4, $5);",
		log.AccountNum,
		log.DestAccount,
		log.TranAmount,
		log.Description,
		log.CreatedAt)
	return err
}
