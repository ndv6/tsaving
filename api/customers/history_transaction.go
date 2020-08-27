package customers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/models"
)

func (ch *CustomerHandler) HistoryTransactionHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ch.jwt.GetToken(r)
		err := token.Valid()
		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		page, err := strconv.Atoi(chi.URLParam(r, "page"))
		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseURLParams)
			return
		}

		listHistoryTransaction, err := models.ListTransactionLog(db, token.CustId, page)
		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, "Cannot get history transaction")
			helpers.HTTPError(w, http.StatusBadRequest, "Cannot get history transaction")
			return
		}
		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetListSuccess, listHistoryTransaction)
		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		w.Header().Set(constants.ContentType, constants.Json)
		fmt.Fprint(w, string(res))
	})

}
