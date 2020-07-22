package models

import (
	"database/sql"
	"time"
)

type VirtualAccounts struct {
	VaId       int       `json:"va_id"`
	VaNum      string    `json:"va_num"`
	AccountNum string    `json:"account_num"`
	VaBalance  int       `json:"va_balance"`
	VaColor    string    `json:"va_color"`
	VaLabel    string    `json:"va_label"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func GetBalanceVA(vaNum string, db *sql.DB) (va VirtualAccounts, err error) {
	var balanceFloat float32
	err = db.QueryRow("SELECT va_balance FROM virtual_accounts WHERE va_num = ($1) ", vaNum).Scan(&balanceFloat)
	if err != nil {
		return
	}
	va.VaBalance = int(balanceFloat)
	return
}
