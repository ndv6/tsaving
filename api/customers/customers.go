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

type RegisterResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

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

func (ch *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}
	var cus models.Customers
	err = json.Unmarshal(b, &cus)

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	if len(cus.CustPassword) < 6 {
		helpers.HTTPError(w, http.StatusBadRequest, constants.InvalidPassword)
		return
	}

	date := time.Now()
	now := date.Format(constants.DateFormat)
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(9999)

	// Account Number YYMMDDNNNN
	AccNum := fmt.Sprint(now, randomNumber)

	// Password Hash
	Pass := helpers.HashString(cus.CustPassword)

	if err := models.RegisterCustomer(ch.db, cus, AccNum, Pass); err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.NotUniquePhoneNumberAndEmail)
		return
	}

	tokenRegister := ch.jwt.Encode(tokens.Token{
		AccountNum: AccNum,
	})

	data := RegisterResponse{
		Token: tokenRegister,
		Email: cus.CustEmail,
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	if err := models.AddEmailTokens(ch.db, tokenRegister, cus.CustEmail); err != nil {
		helpers.HTTPError(w, http.StatusInternalServerError, constants.InsertTokenFailed)
		return
	}

	if err := models.AddAccountsWhileRegister(ch.db, AccNum); err != nil {
		helpers.HTTPError(w, http.StatusInternalServerError, constants.CreateAccountFailed)
		return
	}

	requestBody, err := json.Marshal(map[string]string{
		"email": cus.CustEmail,
		"token": tokenRegister,
	})

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseTnotifRequest)
		return
	}

	_, err = http.Post(constants.TnotifLocalhost+constants.TnotifEndpoint, constants.ApplicationJson, bytes.NewBuffer(requestBody))

	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.ErrorWhenCallingTnotif)
		return
	}

	dataemail := EmailResponse{
		Email: cus.CustEmail,
	}

	err = json.NewEncoder(w).Encode(dataemail)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeTnotifResponse)
		return
	}
}

func (ch *CustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	tokens := ch.jwt.GetToken(r)
	err := tokens.Valid()
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	requestedBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}

	var cus models.Customers
	err = json.Unmarshal(requestedBody, &cus)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	//check if email address is valid
	isValid := isEmailValid(cus.CustEmail)
	if isValid {
		isExist, err := models.IsEmailExist(ch.db, cus.CustEmail, cus.CustId)
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}
		if isExist {
			helpers.HTTPError(w, http.StatusBadRequest, constants.NotUniquePhoneNumberAndEmail)
			return
		}
	} else {
		helpers.HTTPError(w, http.StatusBadRequest, constants.InvalidUserEmail)
		return
	}

	err = models.UpdateProfile(ch.db, cus)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.UpdateUserDataFailed+err.Error())
	}

	result := StatusResult{
		Status: constants.Success,
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
	file, _, err := r.FormFile(constants.UserPhotoFileName)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.RetrieveUserFileFailed)
		return
	}
	defer file.Close()

	// Create a temporary file within our temp-images directory with particular naming pattern
	folderLocation := constants.UserPhotoFolderName
	newFileName := tokens.AccountNum + constants.UserPhotoFileExtension
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
		helpers.HTTPError(w, http.StatusBadRequest, constants.UpdateUserDataFailed+err.Error())
	}

	result := StatusResult{
		Status: constants.Success,
	}

	res, err := json.Marshal(result)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintln(w, string(res))
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(constants.EmailRegex)

	if len(e) < 4 || len(e) > 64 {
		return false
	}
	return emailRegex.MatchString(e)
}
