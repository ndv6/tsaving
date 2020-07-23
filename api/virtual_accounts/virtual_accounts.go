package virtual_accounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/database"
	helper "github.com/ndv6/tsaving/helpers"
	"github.com/ndv6/tsaving/models"
)

type VirtualAcc struct {
	VaNumber   string `json:"va_num"` //ini berarti di request jsonnya "va_num" disimpen di variable VaNum.
	AccountNum string `json:"acc_num"`
	VaColor    string `json:"va_color"`
	VaLabel    string `json:"va_label"`
}

type VirtualAccHandler struct {
	//jwt *token.JWT
	db *sql.DB
}

func NewVirtualAccHandler(db *sql.DB) *VirtualAccHandler {
	// return &VirtualAccHandler{jwt, db}
	return &VirtualAccHandler{db}
}

func (vah *VirtualAccHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	//rows, err := db.Query("SELECT va_num FROM virtual_accounts WHERE account_num = $1;", vac.AccountNum)
	//if err != nil {
	// return
	//}

	// generate random va number
	vaNum := vac.AccountNum + strconv.Itoa(rand.Intn(999)) // combine account number with random number (0-999)

	// insert to db
	var vam models.VirtualAccounts
	// id := select max(id) from db
	// check if VaID < 10 -> vaNum := string(norek + "-00%s", string(vaID))
	// check if VaID < 100 -> vaNum := string(norek + "-0%s", string(vaID))
	// check else -> vaNum := string(norek + "-%s", string(vaID))
	// vaNum := string(norek + "-%s", vaID)
	vam, err = database.CreateVA(vaNum, vac.AccountNum, vac.VaColor, vac.VaLabel, vah.db)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "failed insert data to db")
	}

	fmt.Fprintf(w, "VA Number: %v Created!\n", vam.VaNum)
}

func (vah *VirtualAccHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vaNum := chi.URLParam(r, "vaNumber")

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
	fmt.Printf(vaNum + " " + vac.AccountNum + " " + vac.VaColor + " " + vac.VaLabel)
	var vam models.VirtualAccounts
	vam, err = database.UpdateVA(vaNum, vac.VaColor, vac.VaLabel, vah.db)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "failed insert data to db")
	}

	fmt.Fprintf(w, "VA Number: %v Updated!\n", vam.VaNum)
}
