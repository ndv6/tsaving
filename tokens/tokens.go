package token

import (
	"errors"
	"time"
)

type Token struct {
	User     string    `json :"user"`   //define usernya
	Expired  time.Time `json:"expired"` //define timenya.
	Id_admin int       `json:"id_admin"`
}

func (t *Token) Valid() error {
	if t.Expired.Before(time.Now()) { //kalau sebelum saat ini, token expired
		return errors.New("token expired")
	}
	return nil
}
