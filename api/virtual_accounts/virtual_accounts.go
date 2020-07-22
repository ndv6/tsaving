package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	helper "github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
)

type InputVac struct {
	SaldoInput float64 `json:"perubahan_saldo"`
	RekVac     string  `json:"rekening_vac"`
}

// {
// 	"perubahan_saldo" : 1000.00,
// 	"rekening_vac" : "2009110001001"
// }

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
		helper.HTTPError(w, http.StatusBadRequest, "unable to request body")
		return
	}

	// di parse dan dimasukkan kedalam struct InputVac
	var VirAc InputVac
	err = json.Unmarshal(b, &VirAc)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// cek rekening
	err = helper.CheckRekeningVA(va.db, VirAc.RekVac)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "invalid virtual account number")
		return
	}

	//cek input apakah melebihi saldo

	var ViA models.VirtualAccounts

	//get no rekening by rekening vac
	NoRek, _ := ViA.GetRekeningByVA(va.db, VirAc.RekVac)

	//update balance at both accounts
	err = ViA.UpdateVacBalance(va.db, VirAc.SaldoInput, VirAc.RekVac)
	if err != nil {
		fmt.Fprint(w, err)
		// helper.HTTPError(w, http.StatusBadRequest, "error updating virtual account balance")
		return
	}

	err = ViA.UpdateMainBalance(va.db, VirAc.SaldoInput, NoRek)
	if err != nil {
		fmt.Fprint(w, err)
		// helper.HTTPError(w, http.StatusBadRequest, "error updating account balance")
		return
	}

	return

}

//function di model
