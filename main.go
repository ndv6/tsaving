package main

import (
	"fmt"

	_ "github.com/vicimilenia/tsaving/database"
)

func main() {
	//cara connect ke db
	//db, _ := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db sslmode=disable")
	fmt.Println("ini branch develop")
}
