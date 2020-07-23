package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
	helper "github.com/ndv6/tsaving/helpers"
)

type InputVac struct {
	BalanceChange float64 `json:"balance_change"`
	VacNumber     string  `json:"va_num"`
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
	var VirAcc InputVac
	err = json.Unmarshal(b, &VirAcc)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// cek rekening
	err = helper.CheckAccountVA(va.db, VirAcc.VacNumber, 4)
	// fmt.Fprint(w, err)
	if err != nil {
		// fmt.Fprint(w, err)
		helper.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	//cek input apakah melebihi saldo
	var BalanceChange int = int(VirAcc.BalanceChange)
	returnValue := helper.CheckBalance("VA", VirAcc.VacNumber, BalanceChange, va.db)
	if returnValue == false {
		helper.HTTPError(w, http.StatusBadRequest, "your input is bigger than virtual account balance.")
		return
	}

	//get no rekening by rekening vac
	AccountNumber, _ := database.GetAccountByVA(va.db, VirAcc.VacNumber)

	//update balance at both accounts
	err = database.UpdateVacBalance(va.db, VirAcc.BalanceChange, VirAcc.VacNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "error updating virtual account balance")
		return
	}

	err = database.UpdateMainBalance(va.db, VirAcc.BalanceChange, AccountNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "error updating account balance")
		return
	}

	response := VAResponse{
		Status:       1,
		Notification: fmt.Sprintf("successfully move balance to your main account : %v", VirAcc.BalanceChange),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "unable to encode response")
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
