package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	helper "github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
)

type InputVac struct {
	BalanceChange float64 `json:"balance_change"`
	VacNumber     string  `json:"vac_number"`
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
	err = helper.CheckAccountVA(va.db, VirAcc.VacNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "invalid virtual account number")
		return
	}

	//cek input apakah melebihi saldo

	var ViA models.VirtualAccounts

	//get no rekening by rekening vac
	AccountNumber, _ := ViA.GetAccountByVA(va.db, VirAcc.VacNumber)

	//update balance at both accounts
	err = ViA.UpdateVacBalance(va.db, VirAcc.BalanceChange, VirAcc.VacNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "error updating virtual account balance")
		return
	}

	err = ViA.UpdateMainBalance(va.db, VirAcc.BalanceChange, AccountNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "error updating account balance")
		return
	}

	return

}

//function di model
