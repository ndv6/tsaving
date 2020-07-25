package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ndv6/tsaving/constants"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/helpers"
	helper "github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"

	"github.com/ndv6/tsaving/database"
)

type VirtualAcc struct {
	VaNumber string `json:"va_num"` //ini berarti di request jsonnya "va_num" disimpen di variable VaNum.
	VaColor  string `json:"va_color"`
	VaLabel  string `json:"va_label"`
}

type InputVa struct {
	BalanceChange int `json:"balance_change"`
	// VaNum         string `json:"va_num"`
}

type AddBalanceVARequest struct {
	VaNum     string `json:"va_num"`
	VaBalance int    `json:"va_balance"`
}

type VAResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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

func (vh VAHandler) DeleteVac(w http.ResponseWriter, r *http.Request) {
	token := vh.jwt.GetToken(r)
	err := token.Valid()
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
	}

	var reqBody DeleteVacRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to decode request body")
		return
	}

	cust, err := database.GetCustomerById(vh.db, token.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "User not found")
		return
	}

	vac, err := database.GetVacByAccountNum(vh.db, cust.AccountNum)

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Virtual account not found")
		return
	}

	if vac.VaBalance > 0 {
		err = database.RevertVacBalanceToMainAccount(vh.db, vac)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Fail to revert balance to main account")
			return
		}

		err = models.CreateTransactionLog(vh.db, models.TransactionLogs{
			AccountNum:  vac.AccountNum,
			DestAccount: vac.VaNum,
			TranAmount:  vac.VaBalance,
			Description: models.LogDescriptionVaToMainTemplate(vac.VaBalance, vac.VaNum, vac.AccountNum),
			CreatedAt:   time.Now(),
		})
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Fail to create log transaction")
			return
		}
	}
	err = database.DeleteVacById(vh.db, vac.VaId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Fail to delete virtual account")
		return
	}

	fmt.Fprintf(w, "Success deleting VAC and reverting %d amount of balance to main account", vac.VaBalance)
}

func (va *VAHandler) VacToMain(w http.ResponseWriter, r *http.Request) {
	token := va.jwt.GetToken(r)
	//ambil input dari jsonnya (no rek VAC dan saldo input)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}
	var VirAcc InputVa
	err = json.Unmarshal(b, &VirAcc)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	vaNum := chi.URLParam(r, "va_num")
	// cek rekening
	err = database.CheckAccountVA(va.db, vaNum, token.CustId)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	//get no rekening by rekening vac
	AccountNumber, _ := database.GetAccountByVA(va.db, vaNum)

	//update balance at both accounts
	err = database.UpdateVacToMain(va.db, VirAcc.BalanceChange, vaNum, AccountNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "transfer error")
		return
	}

	response := VAResponse{
		Status:  "sucess",
		Message: fmt.Sprintf("successfully move balance to your main account : %v", VirAcc.BalanceChange),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	logDesc := models.LogDescriptionVaToMainTemplate(VirAcc.BalanceChange, vaNum, token.AccountNum)

	//inpu transaction log
	tLogs := models.TransactionLogs{
		AccountNum:  token.AccountNum,
		DestAccount: vaNum,
		TranAmount:  VirAcc.BalanceChange,
		Description: logDesc,
		CreatedAt:   time.Now(),
	}

	err = models.CreateTransactionLog(va.db, tLogs)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.InitLogFailed)
		return
	}

	return

}

func (va *VAHandler) AddBalanceVA(w http.ResponseWriter, r *http.Request) {
	var vac AddBalanceVARequest
	token := va.jwt.GetToken(r)
	err := json.NewDecoder(r.Body).Decode(&vac)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}
	//check if va number is exist and valid to its owner
	err = database.CheckAccountVA(va.db, vac.VaNum, token.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateBalanceVA := database.TransferFromMainToVa(token.AccountNum, vac.VaNum, vac.VaBalance, va.db)
	if updateBalanceVA != nil {
		helpers.HTTPError(w, http.StatusBadRequest, updateBalanceVA.Error())
		return
	}

	response := VAResponse{
		Status:  1,
		Message: fmt.Sprintf("successfully add balance to your virtual account : %v", vac.VaBalance),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "unable to encode response")
		return
	}

}

func (va *VAHandler) Create(w http.ResponseWriter, r *http.Request) {
	// read request body

	token := va.jwt.GetToken(r)
	req, err := ioutil.ReadAll(r.Body)

	// parse json request
	var vac VirtualAcc
	err = json.Unmarshal(req, &vac)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// initialize model
	var vam models.VirtualAccounts

	// validasi
	am, err := models.GetMainAccount(va.db, token.AccountNum)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "validate account failed, make sure account number is correct")
		return
	}

	// generate va number
	res, err := database.GetListVANum(token.AccountNum, va.db)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to get virtual account list")
		return
	}

	log.Println(res)
	suffixVaNum := "000"
	// get the last of VaNum
	if len(res) > 0 {
		suffixVaNumLast := []rune(res[len(res)-1])
		suffixVaNum = string(suffixVaNumLast[10:])
	}

	lastVaNum, err := strconv.Atoi(suffixVaNum)
	if err != nil {
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
	log.Println(newSuffix)
	log.Println(am.AccountNum)
	log.Println(newVaNum)

	// insert to db
	vam, err = database.CreateVA(newVaNum, token.AccountNum, vac.VaColor, vac.VaLabel, va.db)

	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "failed insert data to db")
		return
	}

	fmt.Fprintf(w, "VA Number: %v Created!\n", vam.VaNum)
}

// to edit VA
func (va *VAHandler) Edit(w http.ResponseWriter, r *http.Request) {

	// read request body
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to read request body")
		return
	}

	// parse json request
	var vac VirtualAcc
	err = json.Unmarshal(req, &vac)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// update to db
	fmt.Printf(vac.VaNumber + " " + " " + vac.VaColor + " " + vac.VaLabel)
	var vam models.VirtualAccounts
	vam, err = database.UpdateVA(vac.VaNumber, vac.VaColor, vac.VaLabel, va.db)

	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "failed insert data to db")
		return
	}

	fmt.Fprintf(w, "Virtual Account: %v Updated!\n", vam.VaNum)
}

// print virtual account list
func (va *VAHandler) VacList(w http.ResponseWriter, r *http.Request) {

	token := va.jwt.GetToken(r)
	res, err := database.GetListVA(va.db, token.CustId)

	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "id must be integer")
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

}
