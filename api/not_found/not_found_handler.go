package not_found

import (
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/helpers"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	helpers.HTTPError(w, 404, fmt.Sprintf("Seems like %v is not available or does not exist", r.URL))
}
