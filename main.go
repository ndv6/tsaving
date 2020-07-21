package main

import (
	"log"

	"github.com/david1312/tsaving/database"
	_ "github.com/david1312/tsaving/database"
)

func main() {
	//cara connect ke db
	_, err := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db sslmode=disable")
	if err != nil {
		log.Fatal("gagal connect ke db")
	}
	//kalo konek sukses
	println("sukses")
}
