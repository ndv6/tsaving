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
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	LogAdmin, count, err := database.GetLogAdmin(la.db, page)
	if err != nil {
		fmt.Fprint(w, err)
		helper.HTTPError(w, r, http.StatusBadRequest, constants.LogAdminFailed)
		return
	}

	responseBody := GetLogAdminResponse{
		Total:        count,
		LogAdminList: LogAdmin,
	}

	_, res, err := helpers.NewResponseBuilder(w, r, true, constants.GetLogAdminSuccess, responseBody)
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprint(w, res)

}

func (la *LogAdminHandler) Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
	tokens := la.jwt.GetTokenAdmin(r)
	err := tokens.Valid()
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var username = tokens.Username

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}

	var lar models.LogAdmin
	err = json.Unmarshal(req, &lar)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	err = database.InsertLogAdmin(la.db, lar, username)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.InsertAdminLogFailed)
		return
	}

	_, res, err := helper.NewResponseBuilder(w, r, true, constants.AddLogAdminSuccess, nil)
	if err != nil {
		helper.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}
	fmt.Fprint(w, string(res))

	return

}

func (la *LogAdminHandler) GetFilteredLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	date := chi.URLParam(r, "date")
	search := chi.URLParam(r, "search")
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	if date != "" && search == "" {
		LogAdmin, count, err := database.GetLogAdminFilteredDate(la.db, date, page)

		if err != nil {
			helper.HTTPError(w, r, http.StatusBadRequest, constants.LogAdminFailed)
			return
		}

		responseBody := GetLogAdminResponse{
			Total:        count,
			LogAdminList: LogAdmin,
		}

		_, res, err := helpers.NewResponseBuilder(w, r, true, constants.GetLogAdminSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if date == "" && search != "" {
		LogAdmin, count, err := database.GetLogAdminFilteredSearch(la.db, search, page)

		if err != nil {
			helper.HTTPError(w, r, http.StatusBadRequest, constants.LogAdminFailed)
			return
		}

		responseBody := GetLogAdminResponse{
			Total:        count,
			LogAdminList: LogAdmin,
		}

		_, res, err := helpers.NewResponseBuilder(w, r, true, constants.GetLogAdminSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	} else if date != "" && search != "" {
		LogAdmin, count, err := database.GetLogAdminFilteredSearchDate(la.db, search, date, page)

		if err != nil {
			helper.HTTPError(w, r, http.StatusBadRequest, constants.LogAdminFailed)
			return
		}

		responseBody := GetLogAdminResponse{
			Total:        count,
			LogAdminList: LogAdmin,
		}

		_, res, err := helpers.NewResponseBuilder(w, r, true, constants.GetLogAdminSuccess, responseBody)
		if err != nil {
			helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprintln(w, string(res))
		return
	}
}
