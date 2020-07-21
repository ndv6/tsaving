package main

import (
	"fmt"
	"log"

	"github.com/jocelyntjahyadi/tsaving/database"
	_ "github.com/jocelyntjahyadi/tsaving/database"
)

type Config struct {
	Addr      string `json:"addr"`
	JWTSecret string `json:"jwt_secret"`
	Database  string `json:"database"`
}

func main() {
	//cara connect ke db
	jwt := token.New([]byte(cfg.JWTSecret))

	db, err := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("branch develop")

}
