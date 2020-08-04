package home

import (
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/helpers"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w, responseJson, _ := helpers.NewResponseBuilder(w, true, "Welcome to tsaving", nil)
	fmt.Println(w.Header())
	fmt.Fprintln(w, responseJson)
}
