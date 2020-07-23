package helpers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/models"
)

func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
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

func CheckBalance(target string, accNumber string, amount int, db *sql.DB) (status bool) {
	if target == "MAIN" {
		sourceBalance, err := database.GetBalanceAcc(accNumber, db)
		if err != nil {
			return
		}
		if sourceBalance.AccountBalance < amount || amount <= 0 {
			return
		}
		status = true
	}
	if target == "VA" {
		sourceBalance, err := database.GetBalanceVA(accNumber, db)
		if err != nil {
			return
		}
		if sourceBalance.VaBalance < amount || amount <= 0 {
			return
		}
		status = true
	}
	return
}

func HashString(toHash string) string {
	hashed := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(hashed[:])
}
