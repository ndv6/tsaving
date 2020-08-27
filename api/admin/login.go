package admin

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

type LoginAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginAdminResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func LoginAdminHandler(jwt *tokens.JWT, db *sql.DB) http.HandlerFunc { // Handle by Caesar Gusti
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.ContentType, constants.Json)
		var l LoginAdminRequest // Ngambil dari body API
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest) //Format JSON Tidak Sesuai
			return
		}

		fmt.Println(l.Username, l.Password)
		//Membuat Hash Password
		Pass := helpers.HashString(l.Password)

		objAdmin, err := models.LoginAdmin(db, l.Username, Pass)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Wrong Username or Password")
			return
		}
		_, tokenLoginAdmin, _ := jwt.JWTAuth.Encode(&tokens.TokenAdmin{
			AdminId:  objAdmin.AdminId,
			Username: objAdmin.Username,
			Expired:  time.Now().Add(120 * time.Minute),
		})

		data := LoginAdminResponse{
			Token:    tokenLoginAdmin,
			Username: objAdmin.Username,
		}

		_, res, err := helpers.NewResponseBuilder(w, true, constants.LoginSucceed, data)

		if err != nil {
			helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
			helpers.SendMessageToTelegram(r, http.StatusInternalServerError, constants.CannotEncodeResponse)
			return
		}

		fmt.Fprint(w, string(res))
	}
}
