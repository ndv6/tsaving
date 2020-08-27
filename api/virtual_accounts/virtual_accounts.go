package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ndv6/tsaving/constants"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
	helper "github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

type VirtualAcc struct {
	VaNumber string `json:"va_num"` //ini berarti di request jsonnya "va_num" disimpen di variable VaNum.
	VaColor  string `json:"va_color"`
	VaLabel  string `json:"va_label"`
}

type InputVa struct {
	BalanceChange int `json:"balance_change"`
}

type AddBalanceVARequest struct {
	VaNum     string `json:"va_num"`
	VaBalance int    `json:"va_balance"`
}

type VAResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type VAListAdminResponse struct {
	Total  int                      `json:"total"`
	VAList []models.VirtualAccounts `json:"data"`
}

type VAHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

type DeleteVacRequest struct {
	VaNum string `json:"va_num"`
}

func NewVAHandler(jwt *tokens.JWT, db *sql.DB) *VAHandler {
	return &VAHandler{jwt, db}
}

// Delete VAC and checkVaNumValid made by Joseph
// Delete VAC relies heavily on query logic so unit testing is not really the perfect choice
// unless we create factory and interface that enables memory storage
func CheckVaNumValid(vaNum string) bool {
	if len(vaNum) > 12 {
		return true
	}
	return false
}

func (vh VAHandler) DeleteVac(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
	token := vh.jwt.GetToken(r)
	err := token.Valid()
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.TokenExpires)
		return
	}

	trx, err := vh.db.Begin()
	if err != nil {
		helpers.HTTPError(w, r, http.StatusInternalServerError, constants.FailSqlTransaction)
	}

	vaNum := chi.URLParam(r, "va_num")
	if !CheckVaNumValid(vaNum) {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.InvalidVaNumber)
		trx.Rollback()
		return
	}

	vac, err := database.GetVacByAccountNum(trx, token.AccountNum, vaNum)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusNotFound, constants.VANotFound)
		trx.Rollback()
		return
	}

	err = database.RevertVacBalanceToMainAccount(trx, vac)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.FailToRevertBalance)
		trx.Rollback()
		return
	}

	if vac.VaBalance > 0 {
		err = models.TransactionLog(trx, models.TransactionLogs{
			AccountNum:  vac.AccountNum,
			FromAccount: vac.VaNum,
			DestAccount: vac.AccountNum,
			TranAmount:  vac.VaBalance,
			Description: constants.TransferToMainAccount,
			CreatedAt:   time.Now(),
		})
		if err != nil {
			helpers.HTTPError(w, r, http.StatusInternalServerError, constants.InitLogFailed)
			trx.Rollback()
			return
		}
	}

	err = trx.Commit()
	if err != nil {
		helpers.HTTPError(w, r, http.StatusInternalServerError, constants.InsertFailed)
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (va *VAHandler) VacToMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
	token := va.jwt.GetToken(r)
	//ambil input dari jsonnya (no rek VAC dan saldo input)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}
	var VirAcc InputVa
	err = json.Unmarshal(b, &VirAcc)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	vaNum := chi.URLParam(r, "va_num")
	// cek rekening
	err = database.CheckAccountVA(va.db, vaNum, token.CustId)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.InvalidVA)
		return
	}

	//get no rekening by rekening vac
	AccountNumber, _ := database.GetAccountByVA(va.db, vaNum)

	//update balance at both accounts
	err = database.UpdateVacToMain(va.db, VirAcc.BalanceChange, vaNum, AccountNumber)
	if err != nil {
		helper.HTTPError(w, r, http.StatusOK, err.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, r, true, fmt.Sprintf("successfully move balance to your main account : %v", VirAcc.BalanceChange), nil)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}
	fmt.Fprint(w, string(res))

	return

}

func (va *VAHandler) AddBalanceVA(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
	token := va.jwt.GetToken(r)

	var vac AddBalanceVARequest
	err := json.NewDecoder(r.Body).Decode(&vac)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}
	//check if va number is exist and valid to its owner
	err = database.CheckAccountVA(va.db, vac.VaNum, token.CustId)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.InvalidVA)
		return
	}

	updateBalanceVA := database.TransferFromMainToVa(token.AccountNum, vac.VaNum, vac.VaBalance, va.db)
	if updateBalanceVA != nil {
		helpers.HTTPError(w, r, http.StatusOK, updateBalanceVA.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, r, true, fmt.Sprintf("successfully add balance to your virtual account : %v", vac.VaBalance), nil)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprint(w, res)
}

func (va *VAHandler) Create(w http.ResponseWriter, r *http.Request) {

	// set response header
	w.Header().Set("Content-Type", "application/json")

	// read request body
	token := va.jwt.GetToken(r)
	req, err := ioutil.ReadAll(r.Body)

	// parse json request
	var vac VirtualAcc
	err = json.Unmarshal(req, &vac)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// initialize model
	var vam models.VirtualAccounts

	// validasi
	am, err := models.GetMainAccount(va.db, token.AccountNum)
	fmt.Println(token.AccountNum)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// generate va number
	res, err := database.GetListVANum(token.AccountNum, va.db)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "unable to get virtual account list")
		return
	}

	suffixVaNum := "000"
	// get the last of VaNum
	if len(res) > 0 {
		suffixVaNumLast := []rune(res[len(res)-1])
		suffixVaNum = string(suffixVaNumLast[10:])
	}

	lastVaNum, err := strconv.Atoi(suffixVaNum)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "error generating va number")
		return
	}

	newSuffix := ""
	if lastVaNum+1 < 10 {
		newSuffix = "00" + strconv.Itoa(lastVaNum+1)
	} else if (lastVaNum + 1) < 100 {
		newSuffix = "0" + strconv.Itoa(lastVaNum+1)
	} else {
		newSuffix = strconv.Itoa(lastVaNum + 1)
	}
	newVaNum := am.AccountNum + newSuffix

	// insert to db
	vam, err = database.CreateVA(newVaNum, token.AccountNum, vac.VaColor, vac.VaLabel, va.db)

	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "failed insert data to db")
		return
	}

	response := VAResponse{
		Status:  "SUCCESS",
		Message: fmt.Sprintf("successfully create virtual account! virtual account number : %v", vam.VaNum),
	}

	// set response http status
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, "unable to encode response")
		return
	}
}

// to edit VA
func (va *VAHandler) Update(w http.ResponseWriter, r *http.Request) {

	// set response header
	w.Header().Set("Content-Type", "application/json")

	// get va number in url
	vaNumber := chi.URLParam(r, "va_num")

	// read request body
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "unable to read request body")
		helper.SendMessageToTelegram(r, http.StatusBadRequest, "unable to read request body")
		return
	}

	// parse json request
	var vac VirtualAcc
	err = json.Unmarshal(req, &vac)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "unable to parse json request")
		helper.SendMessageToTelegram(r, http.StatusBadRequest, "unable to read json body")
		return
	}

	// validasi
	var vam models.VirtualAccounts
	vam, err = database.GetVaNumber(va.db, vaNumber)
	if err != nil {
		helper.SendMessageToTelegram(r, http.StatusBadRequest, "validate va number failed, make sure va number is correct")
		helper.HTTPError(w, r, http.StatusBadRequest, "validate va number failed, make sure va number is correct")
		return
	}

	// update to db
	vam, err = database.UpdateVA(vam.VaNum, vac.VaColor, vac.VaLabel, va.db)

	if err != nil {
		helper.SendMessageToTelegram(r, http.StatusBadRequest, "failed insert data to db")
		helper.HTTPError(w, r, http.StatusBadRequest, "failed insert data to db")
		return
	}

	response := VAResponse{
		Status:  "SUCCESS",
		Message: fmt.Sprintf("successfully edit your virtual account! virtual account number : %v", vam.VaNum),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, "unable to encode response")
		helper.SendMessageToTelegram(r, http.StatusBadRequest, "unable to encode response")
		return
	}
}

func (va *VAHandler) VacList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
	token := va.jwt.GetToken(r)
	res, err := database.GetListVA(va.db, token.CustId)

	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "id must be integer")
		helper.SendMessageToTelegram(r, http.StatusBadRequest, "id must be integer")
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, "unable to parse json request")
		return
	}

}

func (va *VAHandler) VacListAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	custId, err := strconv.Atoi(chi.URLParam(r, "cust_id"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	data, count, err := database.GetListVAAdmin(va.db, custId, page)
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, r, http.StatusBadRequest, "Cannot get va list")
		return
	}

	responseBody := VAListAdminResponse{
		Total:  count,
		VAList: data,
	}

	_, res, err := helpers.NewResponseBuilder(w, r, true, constants.GetListSuccess, responseBody)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))

}

func (va *VAHandler) VacListAdminFilter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	custId, err := strconv.Atoi(chi.URLParam(r, "cust_id"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	color := chi.URLParam(r, "color")
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	data, count, err := database.GetListVAAdminFilter(va.db, custId, color, page)
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		fmt.Println(err)
		helpers.HTTPError(w, r, http.StatusBadRequest, "Cannot get va list")
		return
	}

	responseBody := VAListAdminResponse{
		Total:  count,
		VAList: data,
	}

	_, res, err := helpers.NewResponseBuilder(w, r, true, constants.GetListSuccess, responseBody)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))

}
