package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ndv6/tsaving/database"
	helper "github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
	"github.com/ndv6/tsaving/tokens"
)

type VirtualAcc struct {
	VaNumber string `json:"va_num"` //ini berarti di request jsonnya "va_num" disimpen di variable VaNum.
	VaColor  string `json:"va_color"`
	VaLabel  string `json:"va_label"`
}

type VirtualAccHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

func NewVirtualAccHandler(jwt *tokens.JWT, db *sql.DB) *VirtualAccHandler {
	// return &VirtualAccHandler{jwt, db}
	return &VirtualAccHandler{jwt, db}
}

func (vah *VirtualAccHandler) Create(w http.ResponseWriter, r *http.Request) {
	// read request body

	token := vah.jwt.GetToken(r)
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to read request body")
		return
	}

	// parse json request
	var vac VirtualAcc
	err = json.Unmarshal(req, &vac)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// initialize model
	var vam models.VirtualAccounts
	// var am models.Accounts

	// validasi
	am, err := models.GetMainAccount(vah.db, token.AccountNum)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to validate account. please try again and make sure account is correct")
		return
	}

	// generate random va number
	// vaNum := valAccNum + strconv.Itoa(rand.Intn(999)) // combine account number with random number (0-999)

	// generate va number
	res, err := database.GetListVANum(token.AccountNum, vah.db)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to get virtual account list")
		return
	}

	log.Println(res)
	suffixVaNum := "000"
	// get the last of VaNum
	if len(res) > 0 {
		suffixVaNumLast := []rune(res[len(res)-1])
		suffixVaNum = string(suffixVaNumLast[10:])
	}

	lastVaNum, err := strconv.Atoi(suffixVaNum)
	if err != nil {
		return
	}

	newSuffix := ""
	if lastVaNum+1 < 10 {
		newSuffix = "00" + strconv.Itoa(lastVaNum+1)
	} else if (lastVaNum + 1) < 100 {
		newSuffix = "0" + strconv.Itoa(lastVaNum+1)
	} else {
		newSuffix = strconv.Itoa(lastVaNum + 1)
	}
	newVaNum := am.AccountNum + newSuffix
	log.Println(newSuffix)
	log.Println(am.AccountNum)
	log.Println(newVaNum)

	// insert to db
	vam, err = database.CreateVA(newVaNum, token.AccountNum, vac.VaColor, vac.VaLabel, vah.db)

	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "failed insert data to db")
		return
	}

	fmt.Fprintf(w, "VA Number: %v Created!\n", vam.VaNum)
}

// to edit VA
func (vah *VirtualAccHandler) Edit(w http.ResponseWriter, r *http.Request) {

	// read request body
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to read request body")
		return
	}

	// parse json request
	var vac VirtualAcc
	err = json.Unmarshal(req, &vac)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "unable to parse json request")
		return
	}

	// update to db
	fmt.Printf(vac.VaNumber + " " + " " + vac.VaColor + " " + vac.VaLabel)
	var vam models.VirtualAccounts
	vam, err = database.UpdateVA(vac.VaNumber, vac.VaColor, vac.VaLabel, vah.db)

	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, "failed insert data to db")
		return
	}

	fmt.Fprintf(w, "Virtual Account: %v Updated!\n", vam.VaNum)
}
