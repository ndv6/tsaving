package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
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

	//cek balance 2008210001 ini perlu diupdate ambilnya dari token
	// status := helpers.CheckBalance("MAIN", "2008210001", vac.VaBalance, va.db)
	// if !status {
	// 	helpers.HTTPError(w, http.StatusBadRequest, "insufficient balance")
	// 	return
	// }
	//perlu diupdate ambil dari token
	updateBalanceVA := database.TransferFromMainToVa("2008210001", vac.VaNum, vac.VaBalance, va.db)
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
