package models

type Dashboard struct {
	CustName       string            `json:"cust_name"`
	CustEmail      string            `json:"cust_email"`
	AccountNum     string            `json:"account_num"`
	AccountBalance int               `json:"account_balance"`
	ListVA         []VirtualAccounts `json:"virtual_accounts"`
}
