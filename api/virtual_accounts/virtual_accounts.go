package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/tokens"
)

type AddBalanceVARequest struct {
	VaNum     string `json:"va_num"`
	VaBalance int    `json:"va_balance"`
}

type VAResponse struct {
	Status       int    `json:"status"`
	Notification string `json:"notification"`
}

type VAHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

func NewVAHandler(jwt *tokens.JWT, db *sql.DB) *VAHandler {
	return &VAHandler{jwt, db}
}

func (va *VAHandler) AddBalanceVA(w http.ResponseWriter, r *http.Request) {
	var vac AddBalanceVARequest
	token := va.jwt.GetToken(r)
	err := json.NewDecoder(r.Body).Decode(&vac)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}
	//check if va number is exist and valid to its owner
	err = database.CheckAccountVA(va.db, vac.VaNum, token.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateBalanceVA := database.TransferFromMainToVa(token.AccountNum, vac.VaNum, vac.VaBalance, va.db)
	if updateBalanceVA != nil {
		helpers.HTTPError(w, http.StatusBadRequest, updateBalanceVA.Error())
		return
	}

	response := VAResponse{
		Status:       1,
		Notification: fmt.Sprintf("successfully add balance to your virtual account : %v", vac.VaBalance),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "unable to encode response")
		return
	}

}
