package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

func TransactionLog(db *sql.DB, log models.TransactionLogs) error {
	_, err := db.Exec("INSERT INTO transaction_logs (account_num, dest_account, tran_amount, description, created_at) VALUES ($1, $2, $3, $4, $5);",
		log.AccountNum,
		log.DestAccount,
		log.TranAmount,
		log.Description,
		log.CreatedAt)
	return err
}
