package hybrid

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/pattern/factory"

	"github.com/ndv6/tsaving/constants"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/tokens"
)

func V2Loginhandler(jwt *tokens.JWT, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)

		lh := factory.LoginHandlerFactory(db, jwt, chi.URLParam(r, "role"))
		if lh == nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.InvalidUrlParams)
			return
		}

		obj, err := lh.ManageLogin(r)
		if err != nil {
			// returning err.Error() for more detailed error message
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.LoginSucceed, obj)
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprint(w, res)
	}
}
