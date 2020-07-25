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

var InputVac = []virtual_accounts.InputVa{
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
		true,
		"android",
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

	AccountNumber := GetAccountById()

	status = UpdateVacToMainBalance(InputBalance, AccountVA)
	if status != true {
		t.Fatal("Update Failed")
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

func UpdateVacToMainBalance(InputBal int, AccountVA string) (status bool) {
	return
}
