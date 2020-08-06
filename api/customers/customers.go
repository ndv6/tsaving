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
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/database"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
	"github.com/theplant/luhn"
	"github.com/xlzd/gotp"
)

type EmailResponse struct {
	Email string `json:"email"`
}

type CardResponse struct {
	CardNum string    `json:"card_num"`
	CVV     string    `json:"cvv"`
	Expired time.Time `json:"expired"`
}

type CustomerHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

type GetPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func NewCustomerHandler(jwt *tokens.JWT, db *sql.DB) *CustomerHandler {
	return &CustomerHandler{jwt, db}
}

func (ch *CustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
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

	_, res, err := helpers.NewResponseBuilder(w, true, constants.GetProfilSuccess, cus)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))
}

func (ch *CustomerHandler) GetCardCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	AccountNum := chi.URLParam(r, "account_num")
	cardDetails, err := models.GetDetailsCard(ch.db, AccountNum)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := CardResponse{
		CardNum: cardDetails.CardNum,
		CVV:     cardDetails.Cvv,
		Expired: cardDetails.Expired,
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.GetCardSuccess, data)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))
}

func (ch *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) { // Handle by Caesar Gusti
	b, err := ioutil.ReadAll(r.Body)
	w.Header().Set(constants.ContentType, constants.Json)

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
		helpers.HTTPError(w, http.StatusBadRequest, constants.PasswordRequirement)
		return
	}

	date := time.Now()
	now := date.Format("060102")
	number := string(gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO").Now()[0:4])
	// Account Number YYMMDDNNNN
	AccNum := fmt.Sprint(now, number)

	// Password Hash
	Pass := helpers.HashString(cus.CustPassword)

	// Generate card number
	cardNum, err := GenerateCardNumber(AccNum, date)

	if err != nil {
		w.Header().Set(constants.ContentType, constants.Json)
		helpers.HTTPError(w, http.StatusBadRequest, "Cannot generate card number")
		return
	}

	if err := models.RegisterCustomer(ch.db, cus, AccNum, Pass, cardNum.Number, cardNum.Cvv, cardNum.Expired); err != nil {
		fmt.Println(err)
		helpers.HTTPError(w, http.StatusBadRequest, constants.DupeEmailorPhone)
		return
	}

	OTPEmail := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO").Now()

	if err := models.AddEmailTokens(ch.db, OTPEmail, cus.CustEmail); err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.EmailToken)
		return
	}

	if err := models.AddAccountsWhileRegister(ch.db, AccNum); err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.AccountFailed)
		return
	}

	if err := ch.sendMail(w, OTPEmail, cus.CustEmail); err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.MailFailed)
		return
	}

	data := EmailResponse{
		Email: cus.CustEmail,
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.RegisterSucceed, data)

	if err != nil {
		helpers.HTTPError(w, http.StatusInternalServerError, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprint(w, string(res))

}

func (ch *CustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
	userToken := ch.jwt.GetToken(r)
	err := userToken.Valid()
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
	cus.CustId = userToken.CustId

	//check if email address is valid
	isValid := isEmailValid(cus.CustEmail)
	isEmailChanged, err := models.IsEmailChanged(ch.db, cus.CustEmail, userToken.CustId)
	if isValid {
		if err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		}
		if isEmailChanged {
			isExist, err := models.IsEmailExist(ch.db, cus.CustEmail, userToken.CustId)
			if err != nil {
				helpers.HTTPError(w, http.StatusBadRequest, err.Error())
				return
			}
			if isExist {
				helpers.HTTPError(w, http.StatusBadRequest, constants.EmailTaken)
				return
			}
			cus.IsVerified = false
		} else {
			cus.IsVerified = true
		}
	} else {
		helpers.HTTPError(w, http.StatusBadRequest, constants.InvalidEmail)
		return
	}

	isPhoneExist, err := models.IsPhoneExist(ch.db, cus.CustPhone, userToken.CustId)
	if isPhoneExist {
		helpers.HTTPError(w, http.StatusBadRequest, constants.PhoneTaken)
		return
	}

	err = models.UpdateProfile(ch.db, cus)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.UpdateFailed+err.Error())
	}

	if isEmailChanged {
		OTPEmail := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO").Now()

		if err := models.AddEmailTokens(ch.db, OTPEmail, cus.CustEmail); err != nil {
			helpers.HTTPError(w, http.StatusBadRequest, "Email Token Failed")
			return
		}
		if err := ch.sendMail(w, OTPEmail, cus.CustEmail); err != nil {
			w.Header().Set(constants.ContentType, constants.Json)
			helpers.HTTPError(w, http.StatusBadRequest, constants.MailFailed)
			return
		}
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.UpdateProfileSuccess, nil)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))
}

func (ch *CustomerHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
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
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
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
		helpers.HTTPError(w, http.StatusBadRequest, constants.UpdateFailed+err.Error())
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.UpdatePhotoSuccess, nil)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprintln(w, string(res))
}

func (ch *CustomerHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)
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

	var reqPass GetPasswordRequest
	err = json.Unmarshal(requestedBody, &reqPass)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	if len(reqPass.OldPassword) < 6 || len(reqPass.NewPassword) < 6 {
		helpers.HTTPError(w, http.StatusBadRequest, constants.MinimumPassword)
		return
	}

	hashedOldPass := helpers.HashString(reqPass.OldPassword)
	hashedNewPass := helpers.HashString(reqPass.NewPassword)

	isOldPasswordCorrect, err := models.IsOldPasswordCorrect(ch.db, hashedOldPass, tokens.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !isOldPasswordCorrect {
		helpers.HTTPError(w, http.StatusBadRequest, "Incorrect password")
		return
	}

	err = models.UpdateCustomerPassword(ch.db, hashedNewPass, tokens.CustId)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.UpdatePasswordSuccess, nil)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
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

func (ch *CustomerHandler) sendMail(w http.ResponseWriter, OTPEmail string, cusEmail string) (err error) {

	requestBody, err := json.Marshal(map[string]string{
		"email": cusEmail,
		"token": OTPEmail,
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

func GenerateCardNumber(accNum string, date time.Time) (card models.Card, err error) {
	intAccNum, err := strconv.Atoi(constants.Mastercard + accNum)
	if err != nil {
		return
	}

	cvv := GenerateRandomNumber(100, 999)
	validDigit := luhn.CalculateLuhn(intAccNum)
	cardNumber := strconv.Itoa(intAccNum) + strconv.Itoa(validDigit)
	fmt.Println(cardNumber)

	expired := date.AddDate(5, 0, 0)

	card = models.Card{
		Number:  cardNumber,
		Cvv:     strconv.Itoa(cvv),
		Expired: expired,
	}

	return card, nil
}

func GenerateRandomNumber(min, max int) int {
	return min + rand.Intn(max-min)
}

func (ch *CustomerHandler) GetListCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.ContentType, constants.Json)

	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotParseURLParams)
		return
	}

	listCustomers, err := database.GetListCustomers(ch.db, page)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, true, constants.Success, listCustomers)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}

	fmt.Fprint(w, res)
}
