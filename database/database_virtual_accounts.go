package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ndv6/tsaving/models"
)

func GetBalanceVA(vaNum string, db *sql.DB) (va models.VirtualAccounts, err error) {
	var balanceFloat float32
	err = db.QueryRow("SELECT va_balance FROM virtual_accounts WHERE va_num = ($1) ", vaNum).Scan(&balanceFloat)
	if err != nil {
		return
	}
	va.VaBalance = int(balanceFloat)
	return
}

func CreateVA(vaNum string, accNum string, vaColor string, vaLabel string, db *sql.DB) (va models.VirtualAccounts, err error) {
	_, err = db.Exec("INSERT INTO virtual_accounts (va_num, account_num, va_balance, va_color, va_label, created_at, updated_at)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7) ", vaNum, accNum, 0, vaColor, vaLabel, time.Now(), time.Now())
	if err != nil {
		return va, err
	}
	va.VaNum = vaNum
	return va, err
}

func UpdateVA(vaNum string, vaColor string, vaLabel string, db *sql.DB) (va models.VirtualAccounts, err error) {
	_, err = db.Exec("UPDATE virtual_accounts SET va_color = $1, va_label = $2"+
		" WHERE va_num = $3 ", vaColor, vaLabel, vaNum)
	if err != nil {
		fmt.Println(err)
		return va, err
	}
	va.VaNum = vaNum
	va.VaColor = vaColor
	va.VaLabel = vaLabel
	return va, err
}
