package main

import (
	"fmt"
	"log"
	"net/http"

<<<<<<< HEAD
	"github.com/jocelyntjahyadi/tsaving/database"
	_ "github.com/jocelyntjahyadi/tsaving/database"
=======
	"github.com/ndv6/tsaving/api"
	"github.com/ndv6/tsaving/database"
>>>>>>> 68e54b1bf40f1d82ad6d51f0c2092e1fda8b51c9
)

type Config struct {
	Addr      string `json:"addr"`
	JWTSecret string `json:"jwt_secret"`
	Database  string `json:"database"`
}

func main() {
<<<<<<< HEAD
	//cara connect ke db
	jwt := token.New([]byte(cfg.JWTSecret))

	db, err := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db sslmode=disable")
=======
	db, err := database.GetDatabaseConnection("host=127.0.0.1 port=5432 user=postgres password=password dbname=db_tsaving sslmode=disable")
>>>>>>> 68e54b1bf40f1d82ad6d51f0c2092e1fda8b51c9
	if err != nil {
		log.Fatal(err)
	}

<<<<<<< HEAD
	fmt.Println("branch develop")

=======
	fmt.Println("Server is now accepting request from port 8000")
	err = http.ListenAndServe("127.0.0.1:8000", api.Router(db))
	if err != nil {
		log.Fatal("Can not listen to port 8000", err)
	}
>>>>>>> 68e54b1bf40f1d82ad6d51f0c2092e1fda8b51c9
}
