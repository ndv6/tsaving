// Request, interface, and API for deposit to main account, made by Vici
package customers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ndv6/tsaving/constants"

	"github.com/ndv6/tsaving/models"

	"github.com/ndv6/tsaving/helpers"
)

type DepositRequest struct {
	AccountNumber string `json:"account_number"`
	BalanceAdded  int    `json:"balance_added"`
	AuthCode      string `json:"auth_code"`
	ClientId      int    `json:"client_id"`
}

type PartnerInterface interface {
	GetSecret(id int) (string, error)
}

type Transactor interface {
	DepositToMainAccountDatabaseAccessor(balanceToAdd int, accountNumber string, log models.TransactionLogs) error
}

//This is an API that partner bank(s) will call when one of our customers make a deposit through our bank, made by Vici
func DepositToMainAccount(partner PartnerInterface, trx Transactor) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		var request DepositRequest

		w.Header().Set(constants.ContentType, constants.Json)

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
			return
		}

		err = json.Unmarshal(b, &request)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
			return
		}

		if !helpers.IsRequestValid(request.AccountNumber, request.AuthCode) || !helpers.IsValidInt(request.BalanceAdded, request.ClientId) {
			helpers.HTTPError(w, http.StatusBadRequest, constants.RequestHasInvalidFields)
			return
		}

		isValidAuthCode, err := isValidPartnerAuthCode(request, partner)
		if !isValidAuthCode || err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, constants.UnauthorizedRequest)
			return
		}

		log := models.TransactionLogs{
			AccountNum:  request.AccountNumber,
			DestAccount: request.AccountNumber,
			FromAccount: strconv.Itoa(request.ClientId),
			TranAmount:  request.BalanceAdded,
			Description: constants.Deposit,
			CreatedAt:   time.Now(),
		}
		err = trx.DepositToMainAccountDatabaseAccessor(request.BalanceAdded, request.AccountNumber, log)
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.InsertFailed)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		w, responseJson, err := helpers.NewResponseBuilder(w, true, constants.DepositSuccess, nil)
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintf(w, responseJson)
	})

}

func isValidPartnerAuthCode(request DepositRequest, partner PartnerInterface) (isValidCode bool, err error) {
	isValidCode = false
	clientSecret, err := partner.GetSecret(request.ClientId)
	if err != nil {
		return
	}

	//hashing is a one-way process, so we can only reconstruct the hashed string
	hashedData := helpers.HashString(request.AccountNumber + strconv.Itoa(request.BalanceAdded) + clientSecret)

	if hashedData == request.AuthCode {
		isValidCode = true
	}

	return
}
