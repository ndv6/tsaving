package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ndv6/tsaving/models"

	"github.com/ndv6/tsaving/database"
	helper "github.com/ndv6/tsaving/helpers"
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
	var BalanceChange int = int(VirAcc.BalanceChange)
	fmt.Println(VirAcc.BalanceChange)
	fmt.Println(VirAcc.VacNumber)
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

	// get virtual account info, status
	var VaObj models.VirtualAccounts
	var MaObj models.Accounts
	VaObj, err = database.GetVaStatus(va.db, VirAcc.VacNumber)
	MaObj, err = database.GetBalanceAcc(AccountNumber, va.db)

	fmt.Fprintf(w, "%v VA Balance: %v, ", VaObj.VaLabel, VaObj.VaBalance)
	fmt.Fprintf(w, "Main Account Balance: %v", MaObj.AccountBalance)

	return

}

func (va *VAHandler) VacList(w http.ResponseWriter, r *http.Request) {

	res, err := database.GetListVA(va.db, 4)

	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "id must be integer")
		return
	}

	// err = json.NewEncoder(w).Encode(res) // kalau pake ini, tidak usah ada *

	for i, va := range res {

		fmt.Fprintf(w, "Row %v: Number = %v,  Color = %v, Label = %v, Balance = %v,\n", i, va.VaNum, va.VaColor, va.VaLabel, va.VaBalance)
		t := va.CreatedAt
		fmt.Fprintf(w, "Created : %02d-%02d-%d %02d:%02d:%02d,",
			t.Day(), t.Month(), t.Year(),
			t.Hour(), t.Minute(), t.Second())
		t = va.UpdateAt
		fmt.Fprintf(w, " Updated : %02d-%02d-%d %02d:%02d:%02d\n",
			t.Day(), t.Month(), t.Year(),
			t.Hour(), t.Minute(), t.Second())
	}

	// b, err := json.Marshal(cus) // *
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to encode data to json")
	}

}
