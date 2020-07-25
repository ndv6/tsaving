package models

type Partners struct {
	PartnerId int    `json:"partner_id"`
	ClientId  int    `json:"client_id"`
	Secret    string `json:"secret"`
}
