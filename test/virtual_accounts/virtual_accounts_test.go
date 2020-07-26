package virtual_accounts

import (
	"testing"

	"github.com/ndv6/tsaving/api/virtual_accounts"
)

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

// unit testing created by Joseph
var (
	// length 12
	invalidVa = "200725347100"
	validVa   = "2007253471001"
)

func TestVaNumValid(t *testing.T) {
	res := virtual_accounts.CheckVaNumValid(validVa)
	if !res {
		t.Fatalf("Va Num Valid Testing: Expected %v Got %v", true, res)
	}
}

func TestVaNumInvalid(t *testing.T) {
	res := virtual_accounts.CheckVaNumValid(invalidVa)
	if res {
		t.Fatalf("Va Num Invalid Testing: Expected %v Got %v", false, res)
	}
}
