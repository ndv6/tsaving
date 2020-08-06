package admin

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"

	"github.com/ndv6/tsaving/database"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db}
}

func (ah *AdminHandler) GetDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		act, inact, err := database.GetActInActUserCount(ah.db)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		}

		dashboardAdm := models.DashboardAdmin{
			ActUser:   act,
			InactUser: inact,
		}

		w, res, err := helpers.NewResponseBuilder(w, true, "Success fetching dashboard data", dashboardAdm)
		fmt.Fprint(w, res)
		return
	}
}
