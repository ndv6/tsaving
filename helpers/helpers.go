package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ndv6/tsaving/models"
)

type ResponseBuilder struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

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

func NewResponseBuilder(w http.ResponseWriter, status bool, message string, obj interface{}) (rw http.ResponseWriter, jsoned string, err error) {
	stat := "failed"
	if status {
		stat = "success"
	}
	if obj == nil {
		obj = make(map[string]string)
	}

	b, err := json.Marshal(ResponseBuilder{
		Status:  stat,
		Message: message,
		Data:    obj,
	})

	if err != nil {
		b = []byte(`{}`)
	}
	jsoned = string(b)

	rw = w
	rw.Header().Set("Content-Type", "application/json")
	return
}
