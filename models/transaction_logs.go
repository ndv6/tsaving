package models

import (
	"database/sql"
	"fmt"
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

type Execer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

func TransactionLog(db Execer, log TransactionLogs) error {
	_, err := db.Exec("INSERT INTO transaction_logs (account_num, dest_account, tran_amount, description, created_at) VALUES ($1, $2, $3, $4, $5);",
		log.AccountNum,
		log.DestAccount,
		log.TranAmount,
		log.Description,
		log.CreatedAt)
	return err
}
func LogDescriptionVaToMainTemplate(amount int, vaNum, accountNum string) string {
	return fmt.Sprintf("Transfer %d from Virtual Account %s to Main Account %s", amount, vaNum, accountNum)
}

func LogDescriptionPartnerToMainTemplate(amount, partnerId int, accountNum string) string {
	return fmt.Sprintf("Transfer %d from Partner %d to Main Account %s", amount, partnerId, accountNum)
}

func LogDescriptionMainToVaTemplate(amount int, accountNum, vaNum string) string {
	return fmt.Sprintf("Transfer %d from Main Account %s to Virtual Account %s", amount, accountNum, vaNum)
}
