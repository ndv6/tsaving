package tokens

import (
	"errors"
	"time"
)

type Token struct {
	CustId     int       `json:"cust_id"`
	AccountNum string    `json:"account_num"`
	Expired    time.Time `json:"expired"`
}

func (t *Token) Valid() error {
	// Cek Expired Token
	if t.Expired.Before(time.Now()) {
		return errors.New("Token Expired")
	}
	return nil
}
