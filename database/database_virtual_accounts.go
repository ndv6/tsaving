package database

import (
	"database/sql"
	"errors"
)

func GetBalanceVA(vaNum string, db *sql.DB) (balance int, err error) {
	err = db.QueryRow("SELECT va_balance FROM virtual_accounts WHERE va_num = ($1) ", vaNum).Scan(&balance)
	return
}

func CheckBalance(target string, accNumber string, amount int, db *sql.DB) (status bool) {
	if target == "MAIN" {
		sourceBalance, err := GetBalanceAcc(accNumber, db)
		if err != nil {
			return
		}
		if sourceBalance < amount || amount <= 0 {
			return
		}
		status = true
	}
	if target == "VA" {
		sourceBalance, err := GetBalanceVA(accNumber, db)
		if err != nil {
			return
		}
		if sourceBalance < amount || amount <= 0 {
			return
		}
		status = true
	}
	return
}

func CheckAccountVA(db *sql.DB, VaNum string, id int) (err error) {
	var exist bool
	err = db.QueryRow("SELECT EXISTS(SELECT va_num FROM virtual_accounts INNER JOIN customers ON virtual_accounts.account_num = customers.account_num WHERE va_num = $1 AND cust_id = $2)", VaNum, id).Scan(&exist)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("invalid virtual account number")
		return
	}
	return
}

func CheckAccount(db *sql.DB, AccountNum string, id int) (err error) {
	var exist bool
	err = db.QueryRow("SELECT EXISTS(SELECT account_num FROM customers WHERE account_num = $1 AND cust_id = $2)", AccountNum, id).Scan(&exist)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("invalid account number")
		return
	}
	return
}
