package admin

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db}
}

func (adm *AdminHandler) TransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	transactions, err := database.AllHistoryTransaction(adm.db)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, transactions)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))
}
