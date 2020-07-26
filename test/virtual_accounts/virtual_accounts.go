package virtualaccounts

import (
	"testing"
	"time"

	"github.com/ndv6/tsaving/api/virtual_accounts"
	"github.com/ndv6/tsaving/models"
)

var cust_id = 1
var AccountVA = "2009110001001"
var InputBalance = 1000000

var InputVac = []virtual_accounts.InputVac{
	{
		10000,
		"2009110001",
	},
	{
		100000,
		"2009110002",
	},
	{
		100000,
		"2009110001",
	},
}

var Acc = []models.Accounts{
	{
		1,
		"2009110001",
		100000000,
		time.Now(),
	},
	{
		2,
		"2009110002",
		200000000,
		time.Now(),
	},
}

var Cust = []models.Customers{
	{
		1,
		"2009110001",
		"Lyra",
		"Jalan Dr. Satrio 88, Jakarta",
		"081293829092",
		"lyra@gmail.com",
		"lyra.jpg",
		"lyra",
		true,
		"web",
		time.Now(),
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	},
	{
		2,
		"2009110002",
		"Taylor",
		"Jalan Kuningan, Jakarta",
		"08192839209",
		"taylor@gmail.com",
		"taylor.jpg",
		"taylor",
		"android",
		true,
		time.Now(),
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	},
}

var VirAcc = []models.VirtualAccounts{
	{
		1,
		"2009110001001",
		"2009110001",
		1000000,
		"RED",
		"Tabungan Darurat",
		time.Now(),
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	},
	{
		2,
		"2009110001002",
		"2009110001",
		200000,
		"BLUE",
		"Tabungan Liburan",
		time.Now(),
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	},
	{
		3,
		"2009110002001",
		"2009110002",
		2500000,
		"PURPLE",
		"Tabungan Hadiah",
		time.Now(),
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	},
	{
		4,
		"2009110002002",
		"2009110002",
		250000,
		"YELLOW",
		"Tabungan Cafe",
		time.Now(),
		time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	},
}

func TestVacToMain(t *testing.T) {

	//
	var status = CheckAccountVA(AccountVA)
	if status != true {
		t.Fatal("Invalid Virtual Account Number")
	}

	err := CheckBalance()
	if err != nil {
		t.Fatal("Input is bigger than VA Balance")
	}

	AccountNumber := GetAccountById()

	status = UpdateVacBalance(InputBalance, AccountVA)
	if status != true {
		t.Fatal("Update Failed")
	}

	status = UpdateMainBalance(InputBalance, AccountNumber)

	if len(res) != 2 {
		t.Fatalf("Expect 2 Events, got: %v", len(res))
	}
	if res[0].Name != "Training REST" {
		t.Fatalf("Expected event : Training REST, got : %v", res[0].Name)
	}
	if res[1].Name != "Training Android" {
		t.Fatalf("Expected event : Training Android, got : %v", res[1].Name)
	}

}

func TestVacList(t *testing.T) {

}

// function test support

func VacToMain() {

}

func CheckAccountVA(AccVa string) (status bool) {

	var AccountNumber = GetAccountById()
	for _, v := range VirAcc {
		if v.VaNum == AccVa && v.AccountNum == AccountNumber {
			return true
		}
	}
	return

}

func VacList() {

	res := []models.VirtualAccounts{}
	AccNum := GetAccountById()
	for _, v := range VirAcc {
		if v.AccountNum == AccNum {
			res = append(res, v)
		}
	}
}

func GetAccountById() (AccountNumber string) {
	res := []models.Customers{}
	for _, v := range Cust {
		if v.CustId == cust_id {
			AccountNumber = v.AccountNum
			return
		}
	}
	return
}

func CheckBalance() {
	return
}

const (
	DbAccNum  = "12345678"
	DbVaNum   = "12345678001"
	DbBalance = 50000
)

type PayloadAddBalanceVAC struct {
	AccNum string
	VaNum  string
	Amount int
}

var listPayload = []PayloadAddBalanceVAC{
	{"12345678", "12345678001", 20000},
	{"12345678", "12345678001", 40000},
	{"12345678", "12345678001", 50000},
}

func TestActiveEvents(t *testing.T) {
	var caseNum int = 1
	for _, v := range listPayload {

		if v.AccNum != DbAccNum {
			t.Fatalf("Found Error at case %v : account number '%v' not exist", caseNum, v.AccNum)
		}
		if v.VaNum != DbVaNum {
			t.Fatalf("Found Error at case %v : virtual account number '%v' not exist", caseNum, v.VaNum)
		}
		if v.Amount > DbBalance {
			t.Fatalf("Found Error at case %v : insufficient balance", caseNum)
		}
		caseNum++
	}

}
