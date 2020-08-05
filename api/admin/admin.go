package admin

import (
"database/sql"
"fmt"
"net/http"

"github.com/ndv6/tsaving/constants"
"github.com/ndv6/tsaving/helpers"
"github.com/ndv6/tsaving/database"
"github.com/ndv6/tsaving/tokens"

)

type AdminHandler struct {
	jwt *tokens.JWT
	db *sql.DB
}

func NewAdminHandler(jwt *tokens.JWT, db *sql.DB) *AdminHandler {
	return &AdminHandler{jwt, db}
}

func (adm *AdminHandler) TransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
	token := adm.jwt.GetToken(r)
	err := token.Valid()
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

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