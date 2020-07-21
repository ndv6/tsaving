package not_found_page

import (
	"fmt"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Seems like %v is not available or does not exist", r.URL)
}
