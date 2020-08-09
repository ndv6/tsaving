package admin

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

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

func (adm *AdminHandler) TransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	search := chi.URLParam(r, "search")
	accNum := chi.URLParam(r, "accNum")
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	day := chi.URLParam(r, "day")
	month := chi.URLParam(r, "month")
	year := chi.URLParam(r, "year")

	if accNum != "" && search == "" && day == "" && month == "" && year == "" {
		transactions, count, err := database.CustomerHistoryTransaction(adm.db, accNum, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && search != "" && day == "" && month == "" && year == "" {
		transactions, count, err := database.CustomerHistoryTransactionFiltered(adm.db, accNum, search, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
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
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && day != "" && month != "" && year != "" && search == "" {
		transactions, count, err := database.CustomerHistoryTransactionDateFiltered(adm.db, accNum, day, month, year, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
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
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && day != "" && month != "" && year != "" && search != "" {
		transactions, count, err := database.CustomerHistoryTransactionAllFiltered(adm.db, accNum, search, day, month, year, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
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
			return
		}

		fmt.Fprintln(w, string(res))
		return
	}

	transactions, err := database.AllHistoryTransaction(adm.db)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, transactions)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))
}

func (adm *AdminHandler) TransactionHistoryAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	date := chi.URLParam(r, "date")
	accNum := chi.URLParam(r, "accNum")
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	if accNum == "" && date == "" {
		transactions, count, err := database.AllHistoryTransactionPaged(adm.db, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && date == "" {
		transactions, count, err := database.AllHistoryTransactionFilteredAccNum(adm.db, accNum, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum == "" && date != "" {
		transactions, count, err := database.AllHistoryTransactionFilteredDate(adm.db, date, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && date != "" {
		transactions, count, err := database.AllHistoryTransactionFilteredAccNumDate(adm.db, accNum, date, page)

		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := GetTransactionResponse{
			Total:           count,
			TransactionList: transactions,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	}

}

func (ah *AdminHandler) GetDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		logTransactionToday, err := database.GetLogTransactionToday(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		logAdminToday, err := database.GetLogAdminToday(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		transactionMonthbyWeek, err := database.GetTransactionByWeek(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		dashboardUser := models.DashboardUserResponse{
			ActUser:          act,
			InactUser:        inact,
			TotalTransaction: total,
			NewUserToday:     newUserToday,
			NewUserYesterday: newUserYesterday,
			NewUserThisWeek:  newUserThisWeek,
			NewUserThisMonth: newUserThisMonth,
		}

		dashboardTransaction := models.DashboardTransactionResponse{
			TotalTransactionMonth:     totalTransactionAmountMonth,
			TotalTransactionToday:     totalTransactionAmountToday,
			TotalTransactionYesterday: totalTransactionAmountYesterday,
			TransactionMonth:          transactionMonthbyWeek,
		}

		dashboardAdm := models.DashboardAdmin{
			DashboardUser:        dashboardUser,
			DashboardTransaction: dashboardTransaction,
			LogTransactionToday:  logTransactionToday,
			LogAdminToday:        logAdminToday,
		}

		w, res, err := helpers.NewResponseBuilder(w, true, "Success fetching dashboard data", dashboardAdm)
		fmt.Fprint(w, res)
		return
	}
}
