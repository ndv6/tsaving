package helpers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/david1312/tsaving/models"
)

func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

func CheckBalance(target string, accNumber string, amount int, db *sql.DB) (status bool) {
	if target == "MAIN" {
		saldoasal, err := models.GetBalanceAcc(accNumber, db) // println(err.Error())
		//cek dapet engga balancenya
		if err != nil {
			status = false
			return
		}
		//kalo balancenya dapet lanjut cek saldo yg sesuai engga sama gaboleh negatif yg diinput
		if saldoasal.AccountBalance < amount || amount <= 0 {
			status = false
			return
		}
		status = true
	}
	if target == "VA" {
		saldoasal, err := models.GetBalanceVA(accNumber, db)
		//cek dapet engga balancenya
		println(saldoasal.VaBalance)
		if err != nil {
			status = false
			return
		}
		//kalo balancenya dapet lanjut cek saldo yg sesuai engga sama gaboleh negatif yg diinput
		if saldoasal.VaBalance < amount || amount <= 0 {
			status = false
			return
		}
		status = true
	}
	return
}
