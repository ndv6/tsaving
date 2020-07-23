package email

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/models"
)

func DeleteVerifiedEmailToken(id int, db *sql.DB) (err error) {
	_, err = db.Exec("DELETE FROM email_token WHERE et_id=$1;", id)
	return
}

func UpdateCustomerVerificationStatus(email string, db *sql.DB) (err error) {
	_, err = db.Exec("UPDATE CUSTOMERS SET is_verified = TRUE WHERE cust_email = $1;", email)
	return
}

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

		err = db.QueryRow("SELECT et_id, token, email FROM email_token WHERE token=$1 AND email=$2", et.Token, et.Email).Scan(&et.EtId, &et.Token, &et.Email)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to verify email token: "+err.Error())
			return
		}

		err = UpdateCustomerVerificationStatus(et.Email, db)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		}

		err = DeleteVerifiedEmailToken(et.EtId, db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to delete verified email: "+err.Error())
			return
		}
		fmt.Fprintf(w, "Email %v has been succesfully verified at %v\n", et.Email, time.Now())
	}
}