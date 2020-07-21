package models

import (
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

// func (va *Virtual_Account) TesA() {

// }
