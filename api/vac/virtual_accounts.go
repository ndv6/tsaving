package vac

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ndv6/tsaving/tokens"

	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"

	"github.com/ndv6/tsaving/models"
)

type VaHandler struct {
	Db  *sql.DB
	Jwt *tokens.JWT
}

func NewVaHandler(db *sql.DB, jwt *tokens.JWT) *VaHandler {
	return &VaHandler{
		Db:  db,
		Jwt: jwt,
	}
}

type DeleteVacRequest struct {
	VaNum string `json:"va_num"`
}

func (vh VaHandler) DeleteVac(w http.ResponseWriter, r *http.Request) {
	token := vh.Jwt.GetToken(r)
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

	cust, err := database.GetCustomerById(vh.Db, token.CustId)
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
		err = database.RevertVacBalanceToMainAccount(vh.Db, vac)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Fail to revert balance to main account")
			return
		}

		err = models.CreateTransactionLog(vh.Db, models.TransactionLogs{
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
	err = database.DeleteVacById(vh.Db, vac.VaId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Fail to delete virtual account")
		return
	}

	fmt.Fprintf(w, "Success deleting VAC and reverting %d amount of balance to main account", vac.VaBalance)
}
