package customers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/models"
)

func HistoryTransactionHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		custID := 1
		listHistoryTransaction, err := models.ListTransactionLog(db, custID)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Cannot get history transaction")
			return
		}
		err = json.NewEncoder(w).Encode(listHistoryTransaction)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Can not parse response")
			return
		}
	})
}
