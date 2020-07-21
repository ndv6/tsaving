package main

import (
	"fmt"
)

type Config struct {
	Addr      string `json:"addr"`
	JWTSecret string `json:"jwt_secret"`
	Database  string `json:"database"`
}

func main() {
	fmt.Println("branch develop")

}
