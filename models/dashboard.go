package models

import "time"

type Dashboard struct {
	CustName       string            `json:"cust_name"`
	CustEmail      string            `json:"cust_email"`
	AccountNum     string            `json:"account_num"`
	AccountBalance int               `json:"account_balance"`
	CardNum        string            `json:"card_num"`
	CVV            string            `json:"cvv"`
	Expired        time.Time         `json:"expired"`
	ListVA         []VirtualAccounts `json:"virtual_accounts"`
}
