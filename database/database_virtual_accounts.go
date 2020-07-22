package database

import (
	"database/sql"

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

func UpdateVacBalance(db *sql.DB, balanceInput float64, vacNum string) (err error) {

	_, err = db.Exec("UPDATE virtual_accounts SET va_balance = va_balance - $1 WHERE va_num = $2", balanceInput, vacNum)
	return
}

func UpdateMainBalance(db *sql.DB, balanceInput float64, accountNum string) (err error) {
	_, err = db.Exec("UPDATE accounts SET account_balance = account_balance + $1 WHERE account_num = $2", balanceInput, accountNum)
	return
}

func GetAccountByVA(db *sql.DB, vacNum string) (AccountNum string, err error) {
	err = db.QueryRow("SELECT account_num from virtual_accounts WHERE va_num = $1", vacNum).Scan(&AccountNum)
	return
}

func GetVaStatus(db *sql.DB, vacNum string) (va models.VirtualAccounts, err error) {
	var balance float64
	err = db.QueryRow("SELECT va_label, va_balance FROM virtual_accounts WHERE va_num = $1", vacNum).Scan(&va.VaLabel, &balance)
	va.VaBalance = int(balance)
	return
}
