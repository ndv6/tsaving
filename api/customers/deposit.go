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

type DepositResponse struct {
	DepositStatus string `json:"status"`
}

type PartnerInterface interface {
	GetSecret(id int) (string, error)
}

type Transactor interface {
	AddBalanceToMainAccount(amtToDeposit int, accNum string) error
	LogTransaction(log models.TransactionLogs) error
}

//This is an API that partner bank(s) will call when one of our customers make a deposit through our bank
func DepositToMainAccount(partner PartnerInterface, trx Transactor) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		var request DepositRequest
		var response DepositResponse

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

		if !isValidRequest(request) {
			helpers.HTTPError(w, http.StatusBadRequest, constants.RequestHasInvalidFields)
			return
		}

		isValidAuthCode, err := isValidPartnerAuthCode(request, partner)
		if !isValidAuthCode || err != nil {
			helpers.HTTPError(w, http.StatusUnauthorized, constants.UnauthorizedRequest)
			return
		}

		err = trx.AddBalanceToMainAccount(request.BalanceAdded, request.AccountNumber)
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.InsertFailed)
			return
		}

		response = DepositResponse{
			DepositStatus: constants.DepositSuccess,
		}

		log := models.TransactionLogs{
			AccountNum:  request.AccountNumber,
			DestAccount: request.AccountNumber,
			TranAmount:  request.BalanceAdded,
			Description: constants.Deposit,
			CreatedAt:   time.Now(),
		}

		err = trx.LogTransaction(log)
		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusInternalServerError, constants.InitLogFailed)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}

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

func isValidRequest(request DepositRequest) (isValidRequest bool) {
	isValidRequest = false

	if !(request.AccountNumber == "") && !(request.AuthCode == "") && !(request.BalanceAdded <= 0) && !(request.ClientId == 0) {
		isValidRequest = true
	}
	return
}