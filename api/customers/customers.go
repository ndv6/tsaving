package customers

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
	"database/sql"
	"math/rand"
	"time"
)

type CustomerHandler struct{
	jwt *tokens.JWT
	db *sql.DB
}

func NewCustomerHandler(jwt *tokens.JWT, db *sql.DB) *CustomerHandler{
	return &CustomerHandler{jwt,db}
}

type RegisterResponse struct{
	Token string `json:"token"`
	Email string `json:"email"`
}

func (ch *CustomerHandler) Create(w http.ResponseWriter, r *http.Request){
	b, err := ioutil.ReadAll(r.Body)
	if err != nil{
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

	if err := models.RegisterCustomer(ch.db, cus, AccNum, Pass); err != nil{
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to Register, Your Phone Number Or Email Has Been Used")
		return 
	}

	_, tokenRegister, _  := ch.jwt.Encode(&tokens.Token{
		AccountNum: AccNum,
	})

	data := RegisterResponse{
		Token: tokenRegister,
		Email: cus.CustEmail,
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil{
		helpers.HTTPError(w, http.StatusBadRequest, "Unable to Encode response")
		return
	}

	if err := models.AddEmailTokens(ch.db, tokenRegister, cus.CustEmail); err != nil{
		helpers.HTTPError(w, http.StatusBadRequest, "Email Token Failed")
		return 
	}
}