package test

import (
	"testing"
	"time"

	"github.com/ndv6/tsaving/models"
)

var cust_id = 1

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

	// // ini dibikin skenario, kalau inputnya begini, dan hasilnya oke, nanti akan keluar oke.
	// res := VacToMain()
	// if len(res) != 2 {
	// 	t.Fatalf("Expect 2 Events, got: %v", len(res))
	// }
	// if res[0].Name != "Training REST" {
	// 	t.Fatalf("Expected event : Training REST, got : %v", res[0].Name)
	// }
	// if res[1].Name != "Training Android" {
	// 	t.Fatalf("Expected event : Training Android, got : %v", res[1].Name)
	// }

}

func TestVacList(t *testing.T) {

}

func TestCheckAccountVA(t *testing.T) {

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

func CheckVaColor(VaColor string, VaAccount string) (status bool) {
	for _, v := range VirAcc {
		if v.VaNum == AccVa && v.AccountNum == AccountNumber {
			return true
		}
	}
	return

}

func CheckVaLabel(VaLabel string, VaAccount string) (status bool) {
	var AccountNumber = GetAccountById()
	for _, v := range VirAcc {
		if v.VaNum == AccVa && v.AccountNum == AccountNumber {
			return true
		}
	}
	return

}

func UpdateVac(t *testing.T) {
	var VaColor = "red"
	var VaLabel = "holiday"
	var VaAccount = "2009110001001"

	var status = CheckAccountVA(VaAccount)
	if status != true {
		t.Fatalf("Invalid VA Number")
	}

	var status = CheckVaColor(VaColor, VaAccount)
	if status != true {
		t.Fatalf("Invalid VA Color")
	}

	status = CheckVaLabel(VaLabel)
	if status != true {
		t.Fatalf("Invalid VA Label")
	}

	res := []models.VirtualAccounts{}
	for _, v := range VirAcc {
		if v.VaNum == VaAccount {
			v.VaColor = VaColor
			v.VaLabel = VaLabel
		}
	}

}
