package models

type EmailToken struct {
	Et_id int    `json:"et_id"`
	Token string `json:"token"`
	Email string `json:"email"`
}
