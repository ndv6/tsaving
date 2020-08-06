package models

type DashboardAdmin struct {
	ActUser   int `json:"active_user"`
	InactUser int `json:"inactive_user"`
}
