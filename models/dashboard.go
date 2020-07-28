package models

type Dashboard struct {
	CustName       string            `json:"cust_name"`
	CustPhone      string            `json:"cust_phone"`
	AccountNum     string            `json:"account_num"`
	AccountBalance int               `json:"account_balance"`
	ListVA         []VirtualAccounts `json:"virtual_accounts"`
}
