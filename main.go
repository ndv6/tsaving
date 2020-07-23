package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ndv6/tsaving/api"
	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/tokens"
)

var jwt *tokens.JWT

func main() {
	config, err := helpers.LoadConfig("configs/configs.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.GetDatabaseConnection(config.DbCfg)
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := tokens.New([]byte(config.SecretKey))
	fmt.Println("Server is now accepting request from port " + config.Port)
	err = http.ListenAndServe("127.0.0.1:"+config.Port, api.Router(jwtSecret, db))
	if err != nil {
		log.Fatal("Can not listen to port 8000: ", err)
	}
}
