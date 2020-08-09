package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/models"
)

// made by Joseph
func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	w, resp, err := NewResponseBuilder(w, false, errorMessage, make(map[string]string))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": constants.CannotEncodeResponse})
		w.WriteHeader(http.StatusInternalServerError)
	}
	// fmt.Println(w.Header())
	fmt.Fprintln(w, resp)
}

// Function to hash string, made by Vici
func HashString(toHash string) string {
	hashed := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(hashed[:])
}

func LoadConfig(file string) (models.Config, error) {
	var cfg models.Config
	fm, err := os.Open(file)
	if err != nil {
		return models.Config{}, err
	}
	err = json.NewDecoder(fm).Decode(&cfg)
	if err != nil {
		return models.Config{}, err
	}
	return cfg, err
}

// Function to check if one of request field is empty
func IsRequestValid(strings ...string) (result bool) {
	result = true
	for _, s := range strings {
		if s == "" {
			result = false
			return
		}
	}
	return
}

func IsValidInt(ints ...int) (result bool) {
	result = true
	for _, i := range ints {
		if i <= 0 {
			result = false
			return
		}
	}
	return
}
