package admin

import (
	"database/sql"
	"fmt"
	"net/http"

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

	if accNum != "" && search == "" {
		transactions, err := database.CustomerHistoryTransaction(adm.db, accNum)

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
	} else if accNum != "" && search != "" {
		transactions, err := database.CustomerHistoryTransactionFiltered(adm.db, accNum, search)

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
