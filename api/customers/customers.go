package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
	"net/http"
)

type CustomerHandler struct {
	jwt *token.JWT
	db  *sql.DB
}

func NewCustomerHandler(db *sql.DB, jwt *token.JWT) *CustomerHandler {
	return &CustomerHandler{db: db, jwt: jwt}
}

type GetProfileResult struct {
	Customers models.Customers `json:"customers"`
	Accounts  models.Accounts  `json:"accounts"`
}

func (ch *CustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	token := ch.jwt.GetToken(r)
	err := token.Valid()
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	cus, err := models.GetProfile(ch.db, token.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	acc, err := models.GetMainAccount(ch.db, cus.AccountNum)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := GetProfileResult{
		Customers: cus,
		Accounts:  acc,
	}

	res, err := json.Marshal(result)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintln(w, string(res))
}
