package helpers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ndv6/tsaving/models"
)

func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

func CheckBalance(target string, accNumber string, amount int, db *sql.DB) (status bool) {
	if target == "MAIN" {
		sourceBalance, err := models.GetBalanceAcc(accNumber, db) // println(err.Error())
		//cek dapet engga balancenya
		if err != nil {
			return
		}
		//kalo balancenya dapet lanjut cek saldo yg sesuai engga sama gaboleh negatif yg diinput
		if sourceBalance.AccountBalance < amount || amount <= 0 {
			return
		}
		status = true
	}
	if target == "VA" {
		sourceBalance, err := models.GetBalanceVA(accNumber, db)
		//cek dapet engga balancenya
		if err != nil {
			return
		}
		//kalo balancenya dapet lanjut cek saldo yg sesuai engga sama gaboleh negatif yg diinput
		if sourceBalance.VaBalance < amount || amount <= 0 {
			return
		}
		status = true
	}
	return
}
