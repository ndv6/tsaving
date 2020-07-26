package customers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

var JWT = jwtauth.New("HS256", []byte("secret"), nil)

const (
	LoginSucceed = "Login Succeed"
)

type LoginRequest struct {
	CustEmail    string `json:"cust_email"`
	CustPassword string `json:"cust_password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(jwt *tokens.JWT, db *sql.DB) http.HandlerFunc { // Handle by Caesar Gusti
	return func(w http.ResponseWriter, r *http.Request) {
		var l LoginRequest // Ngambil dari body API
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest) //Format JSON Tidak Sesuai
			return
		}

		//Membuat Hash Password
		Pass := helpers.HashString(l.CustPassword)
		objCustomer, err := models.LoginCustomer(db, l.CustEmail, Pass)
		if err != nil {
			log.Println(err)
			helpers.HTTPError(w, http.StatusBadRequest, "Wrong Email or Password")
			return
		}
		_, tokenLogin, _ := jwt.JWTAuth.Encode(&tokens.Token{
			CustId:     objCustomer.CustId,
			AccountNum: objCustomer.AccountNum,
			Expired:    time.Now().Add(120 * time.Minute),
		})

		data := LoginResponse{
			Token: tokenLogin,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, LoginSucceed, data)

		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprint(w, string(res))
	}
}
