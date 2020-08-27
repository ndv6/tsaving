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
func VerifyEmailToken(eh database.EmailHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
			return
		}

		var et models.EmailToken

		err = json.Unmarshal(requestBody, &et)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
			return
		}

		dbEt, err := eh.GetEmailTokenByEmail(et.Email)
		if err != nil {
			helpers.HTTPError(w, http.StatusNotFound, constants.EmailTokenNotFound)
			helpers.SendMessageToTelegram(r, http.StatusNotFound, constants.EmailTokenNotFound)
			return
		}

		if et.Token != dbEt.Token {
			helpers.HTTPError(w, http.StatusBadRequest, constants.VerifyEmailFailed)
			return
		}

		err = eh.UpdateCustomerVerificationStatusByEmail(dbEt.Email)
		if err != nil {
			helpers.HTTPError(w, http.StatusNotFound, constants.UpdateEmailStatusFailed)
			helpers.SendMessageToTelegram(r, http.StatusNotFound, constants.UpdateEmailStatusFailed)
			return
		}

		err = eh.DeleteVerifiedEmailTokenById(dbEt.EtId)
		if err != nil {
			helpers.HTTPError(w, http.StatusNotFound, constants.DeleteEmailTokenFailed)
			helpers.SendMessageToTelegram(r, http.StatusNotFound, constants.DeleteEmailTokenFailed)
			return
		}

		w, resp, err := helpers.NewResponseBuilder(w, true, constants.SuccessVerifyEmail, models.VerifiedEmailResponse{Email: et.Email})
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}
		fmt.Fprintf(w, resp)
		return
	}
}

func GetEmailToken(eh database.EmailHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		var req models.GetTokenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotReadRequest)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
			return
		}

		dbEt, err := eh.GetEmailTokenByEmail(req.Email)
		if err != nil {
			helpers.SendMessageToTelegram(r, http.StatusNotFound, constants.EmailTokenNotFound)
			helpers.HTTPError(w, http.StatusNotFound, constants.EmailTokenNotFound)
			return
		}

		w, resp, err := helpers.NewResponseBuilder(w, true, constants.SuccessGetToken, models.EmailToken{
			Email: dbEt.Email,
			Token: dbEt.Token,
		})
		if err != nil {
			helpers.SendMessageToTelegram(r, http.StatusInternalServerError, constants.CannotEncodeResponse)
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}
		fmt.Fprintf(w, resp)
		return
	}
}
