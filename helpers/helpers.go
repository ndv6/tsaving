package helpers

import (
	"encoding/json"
	"net/http"
	"os"

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
