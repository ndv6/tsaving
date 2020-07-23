package models

type EmailToken struct {
	EtId  int    `json:"et_id"`
	Token string `json:"token"`
	Email string `json:"email"`
}
