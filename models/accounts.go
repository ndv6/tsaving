package models

import (
	"time"
)

type Accounts struct {
	AccountId      int       `json:"account_id"`
	AccountNum     string    `json:"account_num"`
	AccountBalance int       `json:"account_balance"`
	CreatedAt      time.Time `json:"created_at"`
}
