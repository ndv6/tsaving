package customers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/database"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/helpers"
)

func (ch *CustomerHandler) GetDashboardData(db *sql.DB) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		token := ch.jwt.GetToken(r)
		result, err := database.GetDashboardData(token.CustId, token.AccountNum, db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.Success, result)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, res)
	})
}
