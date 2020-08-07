package models

type DashboardAdmin struct {
	DashboardUser        DashboardUserResponse        `json:"dashboard_user"`
	DashboardTransaction DashboardTransactionResponse `json:"dashboard_transaction"`
	LogTransactionToday  []TransactionLogs            `json:"log_transaction_today"`
	LogAdminToday        []LogAdmin                   `json:"log_admin_today"`
}

type DashboardUserResponse struct {
	ActUser          int `json:"active_user"`
	InactUser        int `json:"inact_user"`
	TotalTransaction int `json:"total_transaction"`
	NewUserToday     int `json:"new_user_today"`
	NewUserYesterday int `json:"new_user_yesterday"`
	NewUserThisWeek  int `json:"new_user_this_week"`
	NewUserThisMonth int `json:"new_user_this_month"`
}

type DashboardTransactionResponse struct {
	TotalTransactionMonth     int                `json:"total_transaction_month"`
	TotalTransactionToday     int                `json:"total_transaction_today"`
	TotalTransactionYesterday int                `json:"total_transaction_yesterday"`
	TransactionMonth          []TransactionMonth `json:"transaction_month"`
}

type TransactionMonth struct {
	Week     int `json:"week"`
	RealWeek int `json:"real_week"`
	Amount   int `json:"amount"`
}
