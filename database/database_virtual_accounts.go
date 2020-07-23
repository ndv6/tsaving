package database

import (
	"database/sql"
	"fmt"

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

func RevertVacBalanceToMainAccount(db *sql.DB, va models.VirtualAccounts) (err error) {
	acc, err := GetAccountByAccountNum(db, va.AccountNum)
	fmt.Println(err)
	if err == nil {
		_, err = db.Exec("UPDATE accounts SET account_balance=$1 WHERE account_id=$2;", acc.AccountBalance+va.VaBalance, acc.AccountId)
		fmt.Println(err)
	}
	return
}

func DeleteVacById(db *sql.DB, vId int) (err error) {
	_, err = db.Exec("DELETE FROM virtual_accounts WHERE va_id=$1;", vId)
	return
}

func GetAccountByAccountNum(db *sql.DB, accountNum string) (acc models.Accounts, err error) {
	err = db.QueryRow("SELECT account_id, account_num, account_balance FROM accounts WHERE account_num=$1", accountNum).Scan(&acc.AccountId, &acc.AccountNum, &acc.AccountBalance)
	return
}

func GetCustomerById(db *sql.DB, id int) (cust models.Customers, err error) {
	err = db.QueryRow("SELECT cust_id, account_num, cust_email FROM customers WHERE cust_id=3;").Scan(&cust.CustId, &cust.AccountNum, &cust.CustEmail)
	return
}

func GetVacByAccountNum(db *sql.DB, accountNum string) (va models.VirtualAccounts, err error) {
	err = db.QueryRow("SELECT va_id, va_num, account_num, va_balance FROM virtual_accounts WHERE account_num=$1", accountNum).Scan(&va.VaId, &va.VaNum, &va.AccountNum, &va.VaBalance)
	return
}
