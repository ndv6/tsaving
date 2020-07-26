package customers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

type EmailResponse struct {
	Email string `json:"email"`
}

type StatusResult struct {
	Status string `json:"status"`
}

type GetProfileResult struct {
	Customers models.Customers `json:"customers"`
	Accounts  models.Accounts  `json:"accounts"`
}
type CustomerHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

func NewCustomerHandler(jwt *tokens.JWT, db *sql.DB) *CustomerHandler {
	return &CustomerHandler{jwt, db}
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

func (ch *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) { // Handle by Caesar Gusti
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}
	var cus models.Customers
	err = json.Unmarshal(b, &cus)

	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	if len(cus.CustPassword) < 6 {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.PasswordRequirement)
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
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.DupeEmailorPhone)
		return
	}

	tokenRegister := ch.jwt.Encode(tokens.Token{
		AccountNum: AccNum,
	})

	if err := models.AddEmailTokens(ch.db, tokenRegister, cus.CustEmail); err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.EmailToken)
		return
	}

	if err := models.AddAccountsWhileRegister(ch.db, AccNum); err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.AccountFailed)
		return
	}

	if err := ch.sendMail(w, tokenRegister, cus.CustEmail); err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, constants.MailFailed)
		return
	}

	data := EmailResponse{
		Email: cus.CustEmail,
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.RegisterSucceed, data)

	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprint(w, string(res))

}

func (ch *CustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userToken := ch.jwt.GetToken(r)
	err := userToken.Valid()
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	requestedBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to read the requested body")
		return
	}

	var cus models.Customers
	err = json.Unmarshal(requestedBody, &cus)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Invalid json type")
		return
	}

	//check if email address is valid
	isValid := isEmailValid(cus.CustEmail)
	isEmailChanged, err := models.IsEmailChanged(ch.db, cus.CustEmail, cus.CustId)
	if isValid {
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}

		if isEmailChanged {
			isExist, err := models.IsEmailExist(ch.db, cus.CustEmail, cus.CustId)
			if err != nil {
				helpers.HTTPError(w, http.StatusBadRequest, err.Error())
				return
			}
			if isExist {
				helpers.HTTPError(w, http.StatusBadRequest, "Email already taken")
				return
			}
			cus.IsVerified = false
		}

	} else {
		helpers.HTTPError(w, http.StatusBadRequest, "Invalid email")
		return
	}

	if len(cus.CustPassword) < 6 {
		helpers.HTTPError(w, http.StatusBadRequest, "Password Min 6 Character")
		return
	}

	cus.CustPassword = helpers.HashString(cus.CustPassword)

	err = models.UpdateProfile(ch.db, cus)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Error updating customer data"+err.Error())
	}

	if isEmailChanged {
		tokenRegister := ch.jwt.Encode(tokens.Token{
			AccountNum: cus.AccountNum,
		})
		if err := models.AddEmailTokens(ch.db, tokenRegister, cus.CustEmail); err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Email Token Failed")
			return
		}
		ch.sendMail(w, tokenRegister, cus.CustEmail)
	}

	result := StatusResult{
		Status: "success",
	}

	res, err := json.Marshal(result)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintln(w, string(res))
}

func (ch *CustomerHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	tokens := ch.jwt.GetToken(r)
	err := tokens.Valid()
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	file, _, err := r.FormFile("myPhoto")
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Error Retrieving the File")
		return
	}
	defer file.Close()

	// Create a temporary file within our temp-images directory with particular naming pattern
	folderLocation := "temp-images"
	newFileName := tokens.AccountNum + ".png"
	tempFile, err := ioutil.TempFile(folderLocation, newFileName)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	pictPath := folderLocation + newFileName
	err = models.UpdateCustomerPicture(ch.db, pictPath, tokens.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, "Error updating customer picure"+err.Error())
	}

	result := StatusResult{
		Status: "success",
	}

	res, err := json.Marshal(result)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintln(w, string(res))
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 4 || len(e) > 64 {
		return false
	}
	return emailRegex.MatchString(e)
}

func (ch *CustomerHandler) sendMail(w http.ResponseWriter, tokenRegister string, cusEmail string) (err error) {

	requestBody, err := json.Marshal(map[string]string{
		"email": cusEmail,
		"token": tokenRegister,
	})

	if err != nil {
		return
	}

	_, err = http.Post("http://localhost:8082/sendMail", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return
	}

	return
}
