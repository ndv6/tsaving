package models

import (
	"database/sql"
	"time"
)

type TransactionLogs struct {
	TlId        int       `json:"tl_id"`
	AccountNum  string    `json:"account_num"`
	DestAccount string    `json:"dest_account"`
	TranAmount  int       `json:"tran_amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type HistoryTransaction struct {
	AccountNum  string    `json:"account_num"`
	DestAccount string    `json:"dest_account"`
	TranAmount  int       `json:"tran_amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func TransactionLog(db *sql.DB, log TransactionLogs) error {
	_, err := db.Exec("INSERT INTO transaction_logs (account_num, dest_account, tran_amount, description, created_at) VALUES ($1, $2, $3, $4, $5);",
		log.AccountNum,
		log.DestAccount,
		log.TranAmount,
		log.Description,
		log.CreatedAt)
	return err
}

func ListTransactionLog(db *sql.DB, id int) (list []HistoryTransaction, err error) {
	accNumber, err := GetAccNumber(db, id)
	if err != nil {
		return list, err
	}

	rows, err := db.Query("SELECT account_num, dest_account, tran_amount, description, created_at FROM transaction_logs WHERE account_num = $1", accNumber)
	if err != nil {
		return list, err
	}

	defer rows.Close()
	for rows.Next() {
		var ht HistoryTransaction
		err = rows.Scan(&ht.AccountNum, &ht.DestAccount, &ht.TranAmount, &ht.Description, &ht.CreatedAt)
		if err != nil {
			return list, err
		}
		list = append(list, ht)
	}
	return
}
