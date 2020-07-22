package models

type TransactionLogs struct {
	TlId        int    `json:"tl_id"`
	AccountNum  string `json:"account_num"`
	DestAccount string `json:"dest_account"`
	TranAmount  int    `json:"tran_amount"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
