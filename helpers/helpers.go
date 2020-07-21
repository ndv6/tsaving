package helper

import (
	"encoding/json"
	"net/http"
)

//untuk ngehandle error"
func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

//untuk ngecek input rekening apakah benar atau tidak.
func check_rekening() {

	// token := ch.jwt.GetToken(r)

}
