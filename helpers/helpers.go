package helpers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

//untuk ngehandle error"
func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

//untuk ngecek input rekening apakah benar atau tidak.
func CheckAccountVA(db *sql.DB, VaNum string) (err error) {
	VaNumber := 0
	err = db.QueryRow("SELECT va_num FROM virtual_accounts WHERE va_num = $1", VaNum).Scan(&VaNumber)
	return
}

func CheckAccount(db *sql.DB, AccountNum string) (err error) {
	AccountNumber := 0
	err = db.QueryRow("SELECT account_num FROM customers WHERE account_num = $1", AccountNum).Scan(&AccountNumber)
	return
}
