package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	helpers "github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
)

type AddBalanceVARequest struct {
	VaNum     string `json:"va_num"`
	VaBalance int    `json:"va_balance"`
}

type AddBalanceVAResponse struct {
	Token string `json:"token"`
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
	//disini karena inputannya gak ada rekeningnya maka cari norek asal dulu ,diisi angka 1 karena blm ada token simpen id customer jadi manual input
	noRek, err := models.GetAccountNumById(1, va.db)
	if err != nil {
		helpers.HTTPError(w, http.StatusUnauthorized, "acc number not found")
		return
	}
	//cek balance
	// fmt.Fprintf(w, "Api Called Sukses : %v", noRek.AccountNum)
	status := helpers.CheckBalance("MAIN", noRek.AccountNum, vac.VaBalance, va.db)
	if !status {
		helpers.HTTPError(w, http.StatusBadRequest, "saldo tidak mencukupi")
		return
	}
	fmt.Fprintln(w, "lanjut coding nambahin balance ke va")
	// json.NewEncoder(w).Encode({"status" : 1})
}
