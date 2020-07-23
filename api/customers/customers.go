package customers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

type CustomerHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

func NewCustomerHandler(jwt *tokens.JWT, db *sql.DB) *CustomerHandler {
	return &CustomerHandler{jwt, db}
}

type RegisterResponse struct {
	Status string `json:"status"`
}

type EmailResponse struct {
	Email string `json:"email"`
}

func (ch *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable Request Body")
		return
	}
	var cus models.Customers
	err = json.Unmarshal(b, &cus)

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to parse JSON Request")
		return
	}

	if len(cus.CustPassword) < 6 {
		helpers.HTTPError(w, http.StatusBadRequest, "Password Min 6 Character")
		return
	}

	date := time.Now()
	now := date.Format("060102")
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(9999)

	// Account Number YYMMDDNNNN
	AccNum := fmt.Sprint(now, randomNumber)

	// Password Hash
	Pass := helpers.HashString(cus.CustPassword)

	if err := models.RegisterCustomer(ch.db, cus, AccNum, Pass); err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to Register, Your Phone Number Or Email Has Been Used")
		return
	}

	_, tokenRegister, _ := ch.jwt.Encode(&tokens.Token{
		AccountNum: AccNum,
	})

	data := RegisterResponse{
		Status: "Register Succedeed",
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to Encode response")
		return
	}

	if err := models.AddEmailTokens(ch.db, tokenRegister, cus.CustEmail); err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Email Token Failed")
		return
	}

	if err := models.AddAccountsWhileRegister(ch.db, AccNum); err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Account Failed")
		return
	}

	requestBody, err := json.Marshal(map[string]string{
		"email": cus.CustEmail,
		"token": tokenRegister,
	})

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to Parse JSON Email")
		return
	}

	resp, err := http.Post("http://localhost:8082/sendMail", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Can't Send Email")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "False Body")
		return
	}

	dataemail := EmailResponse{
		Email: cus.CustEmail,
	}

	err = json.NewEncoder(w).Encode(dataemail)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to Encode response")
		return
	}

	log.Println(string(body))

}

type GetProfileResult struct {
	Customers models.Customers `json:"customers"`
	Accounts  models.Accounts  `json:"accounts"`
}

func (ch *CustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	token := ch.jwt.GetToken(r)
	err := token.Valid()
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	cus, err := models.GetProfile(ch.db, token.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	acc, err := models.GetMainAccount(ch.db, cus.AccountNum)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := GetProfileResult{
		Customers: cus,
		Accounts:  acc,
	}

	res, err := json.Marshal(result)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintln(w, string(res))
}
