package models

type DashboardAdmin struct {
	ActUser          int `json:"active_user"`
	InactUser        int `json:"inact_user"`
	TotalTransaction int `json:"total_transaction"`
}
