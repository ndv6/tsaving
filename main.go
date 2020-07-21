package main

import (
	"fmt"
	"log"

	_ "github.com/vicimilenia/tsaving/database"
)

func main() {
	//cara connect ke db
	db, _ := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ini branch develop")
}
