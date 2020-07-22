package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	helpers "github.com/ndv6/tsaving/helpers"
)

type AddBalanceVARequest struct {
	VaNum      string `json:"va_num"`
	VaBalance  int    `json:"va_balance"`
	AccountNum string `json:"account_num"`
}

type AddBalanceVAResponse struct {
	Status       int    `json:"status"`
	Notification string `json:"notification"`
}

type VAHandler struct {
	db *sql.DB
}

func NewVAHandler(db *sql.DB) *VAHandler {
	return &VAHandler{db}
}

func (va *VAHandler) AddBalanceVA(w http.ResponseWriter, r *http.Request) {
	//CEK inputan body dari apinya dulu sesuai format json apa gak
	var vac AddBalanceVARequest
	err := json.NewDecoder(r.Body).Decode(&vac)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	//cek balance
	status := helpers.CheckBalance("MAIN", vac.AccountNum, vac.VaBalance, va.db)
	if !status {
		helpers.HTTPError(w, http.StatusBadRequest, "insufficient balance")
		return
	}
	token := AddBalanceVAResponse{
		Status:       1,
		Notification: fmt.Sprintf("successfully add balance to your virtual account : %v", vac.VaBalance),
	}
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "unable to encode response")
		return
	}

	// json.NewEncoder(w).Encode({"status" : 1})

}
