package models

type DashboardAdmin struct {
	ActUser                   int               `json:"active_user"`
	InactUser                 int               `json:"inact_user"`
	TotalTransaction          int               `json:"total_transaction"`
	NewUserToday              int               `json:"new_user_today"`
	NewUserYesterday          int               `json:"new_user_yesterday"`
	NewUserThisWeek           int               `json:"new_user_this_week"`
	NewUserThisMonth          int               `json:"new_user_this_month"`
	TotalTransactionToday     int               `json:"total_transaction_today"`
	TotalTransactionYesterday int               `json:"total_transaction_yesterday"`
	LogTransactionToday       []TransactionLogs `json:"log_transaction_totay"`
	LogAdminToday             []LogAdmin        `json:"log_admin_today"`
}
