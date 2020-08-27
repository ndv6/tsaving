package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/xlzd/gotp"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

type TransactionHistoryResponse struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

type AdminHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

type GetTransactionResponse struct {
	Total           int                      `json:"count"`
	TransactionList []models.TransactionLogs `json:"list"`
}

func NewAdminHandler(jwt *tokens.JWT, db *sql.DB) *AdminHandler {
	return &AdminHandler{jwt, db}
}

type AdminInterface interface {
	EditCustomerData(customerData models.Customers, adminUsername string) error
	SendMail(w http.ResponseWriter, OTPEmail string, cusEmail string) error
}

type TokenInterface interface {
	UpsertEmailToken(token string, email string) error
}

type EditCustomerDataRequest struct {
	AdminUsername  string `json:"username"`
	AccountNum     string `json:"account_num"`
	CustEmail      string `json:"cust_email"`
	CustPhone      string `json:"cust_phone"`
	IsVerified     bool   `json:"is_verified"`
	IsEmailChanged bool   `json:"is_email_changed"`
}

func (adm *AdminHandler) EditCustomerData(admInterface AdminInterface, tokenInterface TokenInterface) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		var request EditCustomerDataRequest
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

		if !helpers.IsRequestValid(request.AccountNum, request.AdminUsername, request.CustEmail, request.CustPhone) {
			helpers.HTTPError(w, http.StatusBadRequest, constants.RequestHasInvalidFields)
			return
		}

		customerData := models.Customers{
			AccountNum: request.AccountNum,
			CustEmail:  request.CustEmail,
			CustPhone:  request.CustPhone,
			IsVerified: request.IsVerified,
		}

		err = admInterface.EditCustomerData(customerData, request.AdminUsername)
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.InsertFailed)
			return
		}

		if request.IsEmailChanged {
			OTPEmail := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO").Now()

			if err := tokenInterface.UpsertEmailToken(OTPEmail, request.CustEmail); err != nil {
				helpers.HTTPError(w, http.StatusInternalServerError, constants.GenerateEmailTokenFailed)
				return
			}

			if err := admInterface.SendMail(w, OTPEmail, request.CustEmail); err != nil {
				w.Header().Set(constants.ContentType, constants.Json)
				helpers.HTTPError(w, http.StatusInternalServerError, constants.EditSuccessMailNotSent)
				return
			}
		}

		w, responseJson, err := helpers.NewResponseBuilder(w, true, constants.EditCustomerDataSuccess, nil)
		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintf(w, responseJson)
	})
}

func (adm *AdminHandler) TransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	search := chi.URLParam(r, "search")
	accNum := chi.URLParam(r, "accNum")
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		return
	}

	day := chi.URLParam(r, "day")
	month := chi.URLParam(r, "month")
	year := chi.URLParam(r, "year")

	if accNum != "" && search == "" && day == "" && month == "" && year == "" {
		transactions, count, err := database.CustomerHistoryTransaction(adm.db, accNum, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && search != "" && day == "" && month == "" && year == "" {
		transactions, count, err := database.CustomerHistoryTransactionFiltered(adm.db, accNum, search, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && day != "" && month != "" && year != "" && search == "" {
		transactions, count, err := database.CustomerHistoryTransactionDateFiltered(adm.db, accNum, day, month, year, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())

			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && day != "" && month != "" && year != "" && search != "" {
		transactions, count, err := database.CustomerHistoryTransactionAllFiltered(adm.db, accNum, search, day, month, year, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	}

	transactions, err := database.AllHistoryTransaction(adm.db)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, transactions)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))
}

func (adm *AdminHandler) TransactionHistoryAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	date := chi.URLParam(r, "date")
	search := chi.URLParam(r, "search")
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		return
	}

	if search == "" && date == "" {
		transactions, count, err := database.AllHistoryTransactionPaged(adm.db, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if search != "" && date == "" {
		transactions, count, err := database.AllHistoryTransactionFilteredAccNum(adm.db, search, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if search == "" && date != "" {
		transactions, count, err := database.AllHistoryTransactionFilteredDate(adm.db, date, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if search != "" && date != "" {
		transactions, count, err := database.AllHistoryTransactionFilteredAccNumDate(adm.db, search, date, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	}

}

func (ah *AdminHandler) GetDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		act, inact, err := database.GetActInActUserCount(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		total, err := database.GetTotalTransactionCount(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		newUserToday, err := database.GetNewUserToday(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		newUserYesterday, err := database.GetNewUserYesterday(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		newUserThisWeek, err := database.GetNewUserThisWeek(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		newUserThisMonth, err := database.GetNewUserThisMonth(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		totalTransactionAmountMonth, err := database.GetTransactionAmountMonth(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		totalTransactionAmountToday, err := database.GetTransactionAmountToday(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		totalTransactionAmountYesterday, err := database.GetTransactionAmountYesterday(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		transactionMonthbyWeek, err := database.GetTransactionByWeek(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		transactionAmount, err := database.GetTransactionAmount(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		totalVa, amountVa, err := database.GetTotalFromLog(ah.db, constants.TransferToMainAccount)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}
		totalMain, amountMain, err := database.GetTotalFromLog(ah.db, constants.TransferToVirtualAccount)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}
		totalDeposit, amountDeposit, err := database.GetTotalFromLog(ah.db, constants.Deposit)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		transactionAmountWeek, err := database.GetTransactionAmountWeek(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		dashboardUser := models.DashboardUserResponse{
			ActUser:          act,
			InactUser:        inact,
			NewUserToday:     newUserToday,
			NewUserYesterday: newUserYesterday,
			NewUserThisWeek:  newUserThisWeek,
			NewUserThisMonth: newUserThisMonth,
		}

		dashboardTransaction := models.DashboardTransactionResponse{
			TotalTransactionMonth:         totalTransactionAmountMonth,
			TotalTransactionToday:         totalTransactionAmountToday,
			TotalTransactionAmount:        transactionAmount,
			TotalTransactionAmountVa:      amountVa,
			TotalTransactionAmountMain:    amountMain,
			TotalTransactionAmountDeposit: amountDeposit,
			TotalTransactionYesterday:     totalTransactionAmountYesterday,
			TotalTransactionWeek:          transactionAmountWeek,
			TransactionMonth:              transactionMonthbyWeek,
		}

		dashboardTotalTransaction := models.DashboardTotalTransactionResponse{
			TotalTransaction:            total,
			TotalTransactionVa:          totalVa,
			TotalTransactionMainAccount: totalMain,
			TotalTransactionDeposit:     totalDeposit,
		}

		dashboardAdm := models.DashboardAdmin{
			DashboardUser:             dashboardUser,
			DashboardTransaction:      dashboardTransaction,
			DashboardTotalTransaction: dashboardTotalTransaction,
		}

		w, res, err := helpers.NewResponseBuilder(w, true, "Success fetching dashboard data", dashboardAdm)
		fmt.Fprint(w, res)
		return
	}
}
