package home

import (
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/constants"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, constants.HomePage)
}
