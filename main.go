package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/ndv6/tsaving/api"
	"github.com/ndv6/tsaving/database"
	token "github.com/ndv6/tsaving/tokens"
)

type Config struct {
	Addr      string `json:"addr"`
	JWTSecret string `json:"jwt_secret"`
	Database  string `json:"database"`
}

func main() {
	//cara connect ke db
	jwt := token.New([]byte(cfg.JWTSecret))
	// cfg, err := LoadConfig("config/configs.json")
	// jwt := token.New([]byte(cfg.JWTSecret))
	db, err := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db_tsaving sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is now accepting request from port 8000")
	err = http.ListenAndServe("127.0.0.1:8000", api.Router(db))
	if err != nil {
		log.Fatal("Can not listen to port 8000", err)
	}
}

func LoadConfig(file string) (Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	return cfg, err
}
