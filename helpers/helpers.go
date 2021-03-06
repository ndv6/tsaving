package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/models"
)

func SendMessageToTelegram(r *http.Request, status int, errorMessage string) error {

	current_time := time.Now()
	chat_id := os.Getenv("CHATID")
	text := "There has been an exception.\n" +
		"<b>HTTP Status</b>:" + strconv.Itoa(status) + "\n" +
		"<b>Message</b> : " + errorMessage + "\n" +
		"<b>Timestamp</b> :" + current_time.Format(time.RFC1123) + "\n" +
		"<b>Endpoint</b> :" + html.EscapeString(r.URL.Path) + "\n" +
		"<b>Method</b> :" + r.Method
	data, err := json.Marshal(map[string]string{
		"chat_id":    chat_id,
		"text":       text,
		"parse_mode": "HTML",
	})
	if err != nil {
		return err
	}

	url_bot_telegram := os.Getenv("TELEGRAM")
	resp, err := http.Post(url_bot_telegram, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}

// made by Joseph
func HTTPError(w http.ResponseWriter, r *http.Request, status int, errorMessage string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	w, resp, err := NewResponseBuilder(w, r, false, errorMessage, make(map[string]string))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": constants.CannotEncodeResponse})
		w.WriteHeader(http.StatusInternalServerError)
	}

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
