package customers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

var JWT = jwtauth.New("HS256", []byte("secret"), nil)

type LoginRequest struct {
	CustEmail    string `json:"cust_email"`
	CustPassword string `json:"cust_password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	CustEmail string `json:"cust_email"`
	CustName  string `json:"cust_name"`
}

func LoginHandler(jwt *tokens.JWT, db *sql.DB) http.HandlerFunc { // Handle by Caesar Gusti
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		var l LoginRequest // Ngambil dari body API
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, r, http.StatusBadRequest, constants.CannotReadRequest) //Format JSON Tidak Sesuai
			return
		}

		//Membuat Hash Password
		Pass := helpers.HashString(l.CustPassword)

		isVerified, err := models.CheckLoginVerified(db, l.CustEmail)
		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, r, http.StatusBadRequest, "Failed to check verified status")
			helpers.SendMessageToTelegram(r, http.StatusBadRequest, err.Error())
			return
		}

		if isVerified == false {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, r, http.StatusUnauthorized, "This account is not verified")
			return
		}

		objCustomer, err := models.LoginCustomer(db, l.CustEmail, Pass)
		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, r, http.StatusBadRequest, "Wrong Email or Password")
			return
		}
		_, tokenLogin, _ := jwt.JWTAuth.Encode(&tokens.Token{
			CustId:            objCustomer.CustId,
			AccountNum:        objCustomer.AccountNum,
			AccountExpiration: objCustomer.Expired,
			Expired:           time.Now().Add(120 * time.Minute),
		})

		data := LoginResponse{
			Token:     tokenLogin,
			CustEmail: objCustomer.CustEmail,
			CustName:  objCustomer.CustName,
		}

		_, res, err := helpers.NewResponseBuilder(w, r, true, constants.LoginSucceed, data)

		if err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, r, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprint(w, string(res))
	}
}
