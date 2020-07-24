package email

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ndv6/tsaving/constants"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/models"
)

func VerifyEmailToken(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		var et models.EmailToken

		err = json.Unmarshal(requestBody, &et)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		et, err = database.GetEmailTokenByTokenAndEmail(db, et.Token, et.Email)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.VerifyEmailTokenFailed+err.Error())
			return
		}

		err = database.UpdateCustomerVerificationStatusByEmail(et.Email, db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = database.DeleteVerifiedEmailTokenById(et.EtId, db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.DeleteEmailTokenFailed+err.Error())
			return
		}

		b, err := json.Marshal(models.VerifiedEmailResponse{
			Email:  et.Email,
			Status: constants.Verified,
		})

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}
		fmt.Fprintf(w, string(b))
		return
	}
}
