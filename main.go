package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/polipopoliko/tsaving/tsaving/api"
	"github.com/polipopoliko/tsaving/tsaving/database"
)

func main() {
	db, err := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)

	fmt.Println("Server is now accepting request from port 8000")
	err = http.ListenAndServe("127.0.0.1:8000", api.Router(db))
	if err != nil {
		log.Fatal("Can not listen to port 8000", err)
	}
}
