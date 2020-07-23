package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

func HashString(toHash string) string {
	hashed := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(hashed[:])
}
