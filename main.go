package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ndv6/tsaving/api"
	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/tokens"
)

var jwt *tokens.JWT

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("$JWT_SECRET must be set")
	}
	jwt := tokens.New([]byte(jwtSecret))

	dbUri := os.Getenv("DATABASE_URL")
	if dbUri == "" {
		log.Fatal("$DATABASE_URL must be set")
	}
	db, err := database.DatabaseConnect(dbUri)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is now accepting request from port " + port)
	err = http.ListenAndServe(":"+port, api.Router(jwt, db))
	if err != nil {
		log.Fatal(err)
	}
}
