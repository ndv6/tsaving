package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ndv6/tsaving/constants"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/models"
)

// this verify_email_token.go is made by Joseph

// constantnya sebelu di merge nanti bakal gw adjust lagi jadi nvm it for now okay? :)
const (
	EmailTokenNotFound      = "Can not find requested email"
	VerifyEmailFailed       = "Email fail to be verified with given token"
	UpdateEmailStatusFailed = "Fail to change email status to verified"
	VerifyEmailTokenFailed  = "Unable to verify email token: "
	DeleteEmailTokenFailed  = "Unable to delete verified email"
	CannotReadRequest       = "Cannot read request body"
	CannotParseRequest      = "Unable to parse request"
	CannotEncodeResponse    = "Failed to encode response."
)

func VerifyEmailToken(eh database.EmailHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, CannotReadRequest)
			return
		}

		var et models.EmailToken

		err = json.Unmarshal(requestBody, &et)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, CannotParseRequest)
			return
		}

		dbEt, err := eh.GetEmailTokenByEmail(et.Email)
		if err != nil {
			helpers.HTTPError(w, http.StatusNotFound, EmailTokenNotFound)
			return
		}

		if et.Token != dbEt.Token {
			helpers.HTTPError(w, http.StatusBadRequest, VerifyEmailFailed)
			return
		}

		err = eh.UpdateCustomerVerificationStatusByEmail(et.Email)
		if err != nil {
			helpers.HTTPError(w, http.StatusNotFound, UpdateEmailStatusFailed)
			return
		}

		err = eh.DeleteVerifiedEmailTokenById(et.EtId)
		if err != nil {
			helpers.HTTPError(w, http.StatusNotFound, DeleteEmailTokenFailed)
			return
		}

		w, resp, err := helpers.NewResponseBuilder(w, true, "Email has been successfully verified", models.VerifiedEmailResponse{Email: et.Email})
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
		}
		fmt.Fprintf(w, resp)
		return
	}
}
