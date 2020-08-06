package tokens

import (
	"errors"
	"time"
)

type TokenAdmin struct {
	AdminId  int       `json:"id"`
	Username string    `json:"username"`
	Expired  time.Time `json:"expired"`
}

func (ta *TokenAdmin) Valid() error {
	// Cek Expired Token
	if ta.Expired.Before(time.Now()) {
		return errors.New("Token Expired")
	}
	return nil
}
