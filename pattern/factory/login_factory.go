package factory

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ndv6/tsaving/constants"

	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

type LoginHandler interface {
	ManageLogin(r *http.Request) (obj interface{}, err error)
}

type AdminLoginHandler struct {
	db  *sql.DB
	jwt *tokens.JWT
}

type CustomerLoginHandler struct {
	db  *sql.DB
	jwt *tokens.JWT
}

type LoginRequest struct {
	CustEmail    string `json:"cust_email"`
	CustPassword string `json:"cust_password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	CustEmail string `json:"cust_email"`
	CustName  string `json:"cust_name"`
}

type LoginAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginAdminResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func LoginHandlerFactory(database *sql.DB, token *tokens.JWT, role string) (lh LoginHandler) {
	switch role {
	case "admin":
		lh = AdminLoginHandler{
			db:  database,
			jwt: token,
		}
		break
	case "customer":
		lh = CustomerLoginHandler{
			db:  database,
			jwt: token,
		}
		break
	default:
		return nil
	}
	return
}

func (ah AdminLoginHandler) ManageLogin(r *http.Request) (obj interface{}, err error) {
	var l LoginAdminRequest
	err = json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		err = errors.New(constants.CannotParseRequest)
		return
	}

	Pass := helpers.HashString(l.Password)

	objAdmin, err := models.LoginAdmin(ah.db, l.Username, Pass)
	if err != nil {
		err = errors.New("Invalid username or password")
		return
	}
	_, token, _ := ah.jwt.JWTAuth.Encode(&tokens.TokenAdmin{
		AdminId:  objAdmin.AdminId,
		Username: objAdmin.Username,
		Expired:  time.Now().Add(120 * time.Minute),
	})

	obj = LoginAdminResponse{
		Token:    token,
		Username: objAdmin.Username,
	}
	return
}

func (ch CustomerLoginHandler) ManageLogin(r *http.Request) (obj interface{}, err error) {
	var l LoginRequest
	err = json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		err = errors.New(constants.CannotParseRequest)
		return
	}

	Pass := helpers.HashString(l.CustPassword)

	isVerified, err := models.CheckLoginVerified(ch.db, l.CustEmail, Pass)
	if err != nil {
		err = errors.New("Failed to check verified status")
		return
	}

	if isVerified == false {
		err = errors.New("This account is not verified")
		return
	}

	objCustomer, err := models.LoginCustomer(ch.db, l.CustEmail, Pass)
	if err != nil {
		err = errors.New("Wrong Email or Password")
		return
	}
	_, token, _ := ch.jwt.JWTAuth.Encode(&tokens.Token{
		CustId:            objCustomer.CustId,
		AccountNum:        objCustomer.AccountNum,
		AccountExpiration: objCustomer.Expired,
		Expired:           time.Now().Add(120 * time.Minute),
	})

	obj = LoginResponse{
		Token:     token,
		CustEmail: objCustomer.CustEmail,
		CustName:  objCustomer.CustName,
	}
	return
}
