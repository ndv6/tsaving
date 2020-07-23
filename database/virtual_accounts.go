package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

func RevertBalanceToAccount(db *sql.DB, va models.VirtualAccounts) (err error) {
	acc, err := GetAccountByAccountNum(db, va.AccountNum)
	if err == nil {
		_, err = db.Exec("UPDATE ACCOUNTS SET account_balance=$1 WHERE account_id=$2;", acc.AccountBalance+va.VaBalance, acc.AccountId)
	}
	return
}

func DeleteVacById(db *sql.DB, vId int) (err error) {
	_, err = db.Exec("DELETE FROM VIRTUAL_ACCOUNTS WHERE va_id=$1;", vId)
	return
}

func GetAccountByAccountNum(db *sql.DB, accountNum string) (acc models.Accounts, err error) {
	err = db.QueryRow("SELECT account_id, account_num, account_balance, created_at FROM ACCOUNTS WHERE account_num=$1", accountNum).Scan(&acc.AccountId, &acc.AccountNum, &acc.AccountBalance, &acc.CreatedAt)
	return
}

func GetCustomerById(db *sql.DB, id int) (cust models.Customers, err error) {
	err = db.QueryRow("SELECT cust_id, account_num, email FROM CUSTOMERS WHERE cust_id=$1;", id).Scan(&cust.CustId, &cust.AccountNum, &cust.CustEmail)
	return
}

func GetVacByAccountNum(db *sql.DB, accountNum string) (va models.VirtualAccounts, err error) {
	err = db.QueryRow("SELECT va_id, va_num, account_num, va_balance FROM VIRTUAL_ACCOUNTS WHERE account_num=$1", accountNum).Scan(&va.VaId, &va.VaNum, &va.AccountNum, &va.VaBalance)
	return
}
