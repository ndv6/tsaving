package email

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to verify email token: "+err.Error())
			return
		}

		et, err = database.GetEmailTokenByTokenAndEmail(db, et.Token, et.Email)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to verify email token: "+err.Error())
			return
		}

		err = database.UpdateCustomerVerificationStatus(et.Email, db)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = database.DeleteVerifiedEmailToken(et.Et_id, db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to delete verified email: "+err.Error())
			return
		}

		b, err := json.Marshal(models.VerifiedEmailResponse{
			Email:  et.Email,
			Status: "verified",
		})

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to parse to json")
			return
		}
		fmt.Fprintf(w, string(b))
		return
	}
}
