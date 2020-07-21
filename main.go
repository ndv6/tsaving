package main

import (
	"fmt"
	"log"

	"github.com/david1312/tsaving/helpers"
)

func main() {
	fmt.Println("branch develop")

	var check = helpers.CheckBalance("VA", "2828271", 40000)
	if !check {
		log.Fatal("uang tidak valid")
	}
	println("ok")
}
