package customers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
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
	Token string `json:"token"`
}

func LoginHandler(jwt *tokens.JWT, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var l LoginRequest // Ngambil dari body API
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable parse Request") //Format JSON Tidak Sesuai
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

		_, tokenLogin, _ := jwt.Encode(&tokens.Token{
			CustId:     objCustomer.CustId,
			AccountNum: objCustomer.AccountNum,
			Expired:    time.Now().Add(120 * time.Minute),
		})

		data := LoginResponse{
			Token: tokenLogin,
		}

		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Unable to Encode response")
			return
		}
	}
}
