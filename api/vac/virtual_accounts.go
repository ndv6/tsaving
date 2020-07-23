package vac

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/models"
)

type VacHandler struct {
	Db *sql.DB
}

func (vh VacHandler) DeleteVac(w http.ResponseWriter, r *http.Request) {
	// decode jwt to get id

	cust, err := database.GetCustomerById(vh.Db, 1)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "User not found")
		return
	}

	vac, err := database.GetVacByAccountNum(vh.Db, cust.AccountNum)

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Virtual account not found")
		return
	}

	if vac.VaBalance > 0 {
		err = database.RevertBalanceToAccount(vh.Db, vac)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Fail to revert balance to main account")
			return
		}

		err = models.CreateTransactionLog(vh.Db, models.TransactionLogs{
			AccountNum:  vac.AccountNum,
			DestAccount: vac.VaNum,
			TranAmount:  vac.VaBalance,
			Description: models.LogDescriptionVaToMainTemplate(vac.VaBalance, vac.VaNum, vac.AccountNum),
		})
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Fail to create log transaction")
			return
		}
	} else {
		err = database.DeleteVacById(vh.Db, vac.VaId)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Fail to revert balance to main account")
			return
		}
	}

	fmt.Fprintf(w, "Success deleting VAC and reverting %d amount of balance to main account", vac.VaBalance)
}
