package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
	fmt.Fprintln(w, resp)
}

func Base64Decoder(b64 string) (dec io.Reader, err error) {
	idx := strings.Index(b64, ",")
	if idx < 0 {
		err = errors.New(constants.Base64DecodeFailed)
		return
	}
	dec = base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64[idx+1:]))
	return
}

func GenerateStaticImagePath(id int) string {
	return fmt.Sprintf("%v/customer_%v/", constants.StaticImagePath, id)
}

func RemoveAllFilesInDir(dirPath string) (err error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, file := range files {
		err = os.Remove(dirPath + file.Name())
		if err != nil {
			return
		}
	}
	return
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
