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
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db}
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
		transactions, err := database.CustomerHistoryTransaction(adm.db, accNum, page)

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
		return
	} else if accNum != "" && search != "" && day == "" && month == "" && year == "" {
		transactions, err := database.CustomerHistoryTransactionFiltered(adm.db, accNum, search, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, transactions)
		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && day != "" && month != "" && year != "" && search == "" {
		transactions, err := database.CustomerHistoryTransactionDateFiltered(adm.db, accNum, day, month, year, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, transactions)
		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if accNum != "" && day != "" && month != "" && year != "" && search != "" {
		transactions, err := database.CustomerHistoryTransactionAllFiltered(adm.db, accNum, search, day, month, year, page)

		if err != nil {
			fmt.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetAllTransactionSuccess, transactions)
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

		dashboardAdm := models.DashboardAdmin{
			ActUser:          act,
			InactUser:        inact,
			TotalTransaction: total,
		}

		w, res, err := helpers.NewResponseBuilder(w, true, "Success fetching dashboard data", dashboardAdm)
		fmt.Fprint(w, res)
		return
	}
}
