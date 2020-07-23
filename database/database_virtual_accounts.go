package database

import (
	"database/sql"
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
		return
	}
	va.VaNum = vaNum
	return va, err
}

func UpdateVA(vaNum string, vaColor string, vaLabel string, db *sql.DB) (va models.VirtualAccounts, err error) {
	_, err = db.Exec("UPDATE virtual_accounts SET va_color = $1, va_label = $2"+
		" WHERE va_num = $3 ", vaColor, vaLabel, vaNum)
	if err != nil {
		return
	}
	va.VaNum = vaNum
	va.VaColor = vaColor
	va.VaLabel = vaLabel
	return va, err
}

func GetListVANum(accNum string, db *sql.DB) (res []string, err error) {
	rows, err := db.Query("SELECT va_num FROM virtual_accounts WHERE account_num = $1", accNum)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var va_num string
		err = rows.Scan(&va_num)
		if err != nil {
			return
		}
		res = append(res, va_num)
	}
	return res, nil
}

func GetMaxVANum(accNum string, db *sql.DB) (maxId int, err error) {
	row := db.QueryRow("SELECT max(va_id) FROM virtual_accounts WHERE account_num = $1", accNum)
	err = row.Scan(&maxId)
	if err != nil {
		return
	}
	return maxId, nil
}
