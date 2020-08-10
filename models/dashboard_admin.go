package models

type DashboardAdmin struct {
	DashboardUser             DashboardUserResponse             `json:"dashboard_user"`
	DashboardTransaction      DashboardTransactionResponse      `json:"dashboard_amount"`
	DashboardTotalTransaction DashboardTotalTransactionResponse `json:"dashboard_transaction"`
}

type DashboardUserResponse struct {
	ActUser          int `json:"active_user"`
	InactUser        int `json:"inact_user"`
	NewUserToday     int `json:"new_user_today"`
	NewUserYesterday int `json:"new_user_yesterday"`
	NewUserThisWeek  int `json:"new_user_this_week"`
	NewUserThisMonth int `json:"new_user_this_month"`
}

type DashboardTransactionResponse struct {
	TotalTransactionAmount        int                `json:"total_transaction_amount"`
	TotalTransactionMonth         int                `json:"total_transaction_month"`
	TotalTransactionWeek          int                `json:"total_transaction_week"`
	TotalTransactionToday         int                `json:"total_transaction_today"`
	TotalTransactionYesterday     int                `json:"total_transaction_yesterday"`
	TotalTransactionAmountVa      int                `json:"total_transaction_amount_va"`
	TotalTransactionAmountMain    int                `json:"total_transaction_amount_main"`
	TotalTransactionAmountDeposit int                `json:"total_transaction_amount_deposit"`
	TransactionMonth              []TransactionMonth `json:"transaction_month"`
}

type DashboardTotalTransactionResponse struct {
	TotalTransaction            int `json:"total_transaction"`
	TotalTransactionVa          int `json:"total_transaction_va"`
	TotalTransactionMainAccount int `json:"total_transaction_main"`
	TotalTransactionDeposit     int `json:"total_transaction_deposit"`
}

type TransactionMonth struct {
	Week     int `json:"week"`
	RealWeek int `json:"real_week"`
	Amount   int `json:"amount"`
}
