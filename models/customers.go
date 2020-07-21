package models

import (
	"database/sql"
	"time"

	token "github.com/david1312/tsaving/tokens"
)

type Customers struct {
	CustId       int       `json:"cust_id"`
	AccountNum   string    `json:"account_num"`
	CustName     string    `json:"cust_name"`
	CustAddress  int       `json:"cust_address"`
	CustPhone    string    `json:"cust_phone"`
	CustEmail    string    `json:"cust_email"`
	CustPassword string    `json:"cust_password"`
	CustPict     string    `json:"cust_pict"`
	IsVerified   string    `json:"is_verified"`
	Channel      string    `json:"channel"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

type CustomersHandler struct {
	jwt *token.JWT
	db  *sql.DB
}
