package models

import (
	"database/sql"
	"time"
)

type VirtualAccounts struct {
	VaId       int       `json:"var_id"`
	VaNum      string    `json:"va_num"`
	AccountNum string    `json:"account_num"`
	VaBalance  int       `json:"va_balance"`
	VaColor    string    `json:"va_color"`
	VaLabel    string    `json:"va_label"`
	CreatedAt  time.Time `json:"created_at"`
	UpdateAt   time.Time `json:"updated_at"`
}

func GetBalanceVA(vaNum string, db *sql.DB) (va VirtualAccounts, err error) {
	var balancefloat float32
	err = db.QueryRow("SELECT va_balance from virtual_accounts where va_num = ($1) ", vaNum).Scan(&balancefloat)
	va.VaBalance = int(balancefloat)
	return
}
