package virtual_accounts

import (
	"testing"
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
	{"12345678", "12345678001", 0},
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
