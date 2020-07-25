package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ndv6/tsaving/models"
)

//untuk ngehandle error"
func HTTPError(w http.ResponseWriter, status int, errorMessage string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

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
