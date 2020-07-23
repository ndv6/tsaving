package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ndv6/tsaving/database"
	helper "github.com/ndv6/tsaving/helpers"
)

type InputVa struct {
	BalanceChange int    `json:"balance_change"`
	VaNum         string `json:"va_num"`
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

func (va *VAHandler) VacToMain(w http.ResponseWriter, r *http.Request) {

	//ambil input dari jsonnya (no rek VAC dan saldo input)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to read request body")
		return
	}

	// di parse dan dimasukkan kedalam struct InputVac
	var VirAcc InputVa
	err = json.Unmarshal(b, &VirAcc)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// cek rekening
	err = database.CheckAccountVA(va.db, VirAcc.VaNum, 4)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	//cek input apakah melebihi saldo
	var BalanceChange int = VirAcc.BalanceChange
	returnValue := database.CheckBalance("VA", VirAcc.VaNum, BalanceChange, va.db)
	if returnValue == false {
		helper.HTTPError(w, http.StatusBadRequest, "your input is bigger than virtual account balance.")
		return
	}

	//get no rekening by rekening vac
	AccountNumber, _ := database.GetAccountByVA(va.db, VirAcc.VaNum)

	//update balance at both accounts
	err = database.UpdateVacToMain(va.db, VirAcc.BalanceChange, VirAcc.VaNum, AccountNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "transfer error")
		return
	}

	response := VAResponse{
		Status:       1,
		Notification: fmt.Sprintf("successfully move balance to your main account : %v", VirAcc.BalanceChange),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to encode response")
		return
	}
	return

}

func (va *VAHandler) VacList(w http.ResponseWriter, r *http.Request) {

	res, err := database.GetListVA(va.db, 4)

	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "id must be integer")
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

}
