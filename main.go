package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ndv6/tsaving/api"
	"github.com/ndv6/tsaving/database"
)

type Config struct {
	Addr      string `json:"addr"`
	JWTSecret string `json:"jwt_secret"`
	Database  string `json:"database"`
}

func main() {
	//cara connect ke db
	db, err := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db_tsaving sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is now accepting request from port 8001")
	err = http.ListenAndServe("127.0.0.1:8001", api.Router(db))
	if err != nil {
		log.Fatal("Can not listen to port 8001", err)
	}
}
