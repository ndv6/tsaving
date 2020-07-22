package helper

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
func CheckRekeningVA(db *sql.DB, RekVA string) error {
	VirtualRek := 0
	err := db.QueryRow("SELECT va_num FROM virtual_accounts WHERE va_num = $1", RekVA).Scan(&VirtualRek)

	if VirtualRek == 0 {
		return err
	}

	return nil

	// token := ch.jwt.GetToken(r)

}

func CheckRekening(db *sql.DB, Rek string) error {
	NoRek := 0
	err := db.QueryRow("SELECT account_num FROM customers WHERE account_num = $1", Rek).Scan(&NoRek)

	if NoRek == 0 {
		return err
	}

	return nil

}
