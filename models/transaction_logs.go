package models

import (
	"database/sql"
	"time"
)

type TransactionLogs struct {
	TlId        int       `json:"tl_id"`
	AccountNum  string    `json:"account_num"`
	FromAccount string    `json:"from_account"`
	DestAccount string    `json:"dest_account"`
	TranAmount  int       `json:"tran_amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type HistoryTransaction struct {
	AccountNum  string    `json:"account_num"`
	FromAccount string    `json:"from_account"`
	DestAccount string    `json:"dest_account"`
	TranAmount  int       `json:"tran_amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Execer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

func CreateTransactionLog(db *sql.DB, log TransactionLogs) error {
	_, err := db.Exec("INSERT INTO transaction_logs (account_num, from_account, dest_account, tran_amount, description, created_at) VALUES ($1, $2, $3, $4, $5);",
		log.AccountNum,
		log.FromAccount,
		log.DestAccount,
		log.TranAmount,
		log.Description,
		log.CreatedAt)
	return err
}

func TransactionLog(db Execer, log TransactionLogs) error {
	_, err := db.Exec("INSERT INTO transaction_logs (account_num, from_account, dest_account, tran_amount, description, created_at) VALUES ($1, $2, $3, $4, $5, $6);",
		log.AccountNum,
		log.FromAccount,
		log.DestAccount,
		log.TranAmount,
		log.Description,
		log.CreatedAt)
	return err
}

func ListTransactionLog(db *sql.DB, id int, page int) (list []HistoryTransaction, err error) {
	accNumber, err := GetAccNumber(db, id)
	if err != nil {
		return
	}

	offset := (page - 1) * 20
	rows, err := db.Query("SELECT account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE account_num = $1 ORDER BY created_at DESC OFFSET $2 LIMIT 20", accNumber, offset)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var ht HistoryTransaction
		err = rows.Scan(&ht.AccountNum, &ht.DestAccount, &ht.FromAccount, &ht.TranAmount, &ht.Description, &ht.CreatedAt)
		if err != nil {
			return
		}
		list = append(list, ht)
	}
	return
}
