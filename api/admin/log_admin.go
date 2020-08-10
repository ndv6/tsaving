package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/database"
	helper "github.com/ndv6/tsaving/helpers"
)

type GetLogAdminResponse struct {
	Total        int               `json:"count"`
	LogAdminList []models.LogAdmin `json:"list"`
}

type LogAdminHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

func NewLogAdminHandler(jwt *tokens.JWT, db *sql.DB) *LogAdminHandler {
	return &LogAdminHandler{jwt, db}
}

func (la *LogAdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	LogAdmin, count, err := database.GetLogAdmin(la.db, page)
	if err != nil {
		fmt.Fprint(w, err)
		helper.HTTPError(w, http.StatusBadRequest, constants.LogAdminFailed)
		return
	}

	responseBody := GetLogAdminResponse{
		Total:        count,
		LogAdminList: LogAdmin,
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.GetLogAdminSuccess, responseBody)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprint(w, res)

}

func (la *LogAdminHandler) Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	var username = "admin" //get from token (later)

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}

	var lar models.LogAdmin
	err = json.Unmarshal(req, &lar)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	err = database.InsertLogAdmin(la.db, lar, username)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.InsertAdminLogFailed)
		return
	}

	_, res, err := helper.NewResponseBuilder(w, true, constants.AddLogAdminSuccess, nil)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}
	fmt.Fprint(w, string(res))

	return

}

func (la *LogAdminHandler) GetFilteredLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	date := chi.URLParam(r, "date")
	username := chi.URLParam(r, "username")
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	if date != "" && username == "" {
		LogAdmin, count, err := database.GetLogAdminFilteredDate(la.db, date, page)

		if err != nil {
			helper.HTTPError(w, http.StatusBadRequest, constants.LogAdminFailed)
			return
		}

		responseBody := GetLogAdminResponse{
			Total:        count,
			LogAdminList: LogAdmin,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetLogAdminSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if date == "" && username != "" {
		LogAdmin, count, err := database.GetLogAdminFilteredUsername(la.db, username, page)

		if err != nil {
			helper.HTTPError(w, http.StatusBadRequest, constants.LogAdminFailed)
			return
		}

		responseBody := GetLogAdminResponse{
			Total:        count,
			LogAdminList: LogAdmin,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetLogAdminSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if date != "" && username != "" {
		LogAdmin, count, err := database.GetLogAdminFilteredUsernameDate(la.db, username, date, page)

		if err != nil {
			helper.HTTPError(w, http.StatusBadRequest, constants.LogAdminFailed)
			return
		}

		responseBody := GetLogAdminResponse{
			Total:        count,
			LogAdminList: LogAdmin,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.GetLogAdminSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	}
}
