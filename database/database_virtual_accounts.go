package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ndv6/tsaving/models"
)

func GetBalanceVA(vaNum string, db *sql.DB) (balance int, err error) {
	err = db.QueryRow("SELECT va_balance FROM virtual_accounts WHERE va_num = ($1) ", vaNum).Scan(&balance)
	return
}

func UpdateVacBalance(db *sql.DB, balanceInput int, vacNum string) (err error) {

	_, err = db.Exec("UPDATE virtual_accounts SET va_balance = va_balance - $1 WHERE va_num = $2", balanceInput, vacNum)
	return
}

func UpdateMainBalance(db *sql.DB, balanceInput int, accountNum string) (err error) {
	_, err = db.Exec("UPDATE accounts SET account_balance = account_balance + $1 WHERE account_num = $2", balanceInput, accountNum)
	return
}

func UpdateVacToMain(db *sql.DB, balanceInput int, vacNum string, accountNum string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	var balanceVA int
	err = tx.QueryRow("SELECT va_balance FROM virtual_accounts WHERE account_num = $1 FOR UPDATE", accountNum).Scan(&balanceVA)
	if err != nil {
		tx.Rollback()
		return
	}

	if balanceVA < balanceInput {
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE accounts SET account_balance = account_balance + $1 WHERE account_num = $2", balanceInput, accountNum)
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = tx.Exec("UPDATE virtual_accounts SET va_balance = va_balance - $1 WHERE va_num = $2", balanceInput, vacNum)
	if err != nil {
		tx.Rollback()
		return
	}

	logDesc := models.LogDescriptionVaToMainTemplate(balanceInput, vacNum, accountNum)

	//inpu transaction log
	tLogs := models.TransactionLogs{
		AccountNum:  accountNum,
		DestAccount: vacNum,
		TranAmount:  balanceInput,
		Description: logDesc,
		CreatedAt:   time.Now(),
	}

	err = models.TransactionLog(db, tLogs)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func GetAccountByVA(db *sql.DB, vacNum string) (AccountNum string, err error) {
	err = db.QueryRow("SELECT account_num from virtual_accounts WHERE va_num = $1", vacNum).Scan(&AccountNum)
	return
}

func GetListVA(db *sql.DB, id int) (VirAcc []models.VirtualAccounts, err error) {
	rows, err := db.Query("SELECT va_id, va_num, virtual_accounts.account_num, va_label, va_color, va_balance, virtual_accounts.created_at, virtual_accounts.created_at FROM virtual_accounts INNER JOIN customers ON virtual_accounts.account_num = customers.account_num WHERE cust_id = $1", id)
	if err != nil {
		return VirAcc, err
	}

	defer rows.Close()
	//defer -> yang harus dipanggil di akhir (biar ga lupa

	res := make([]models.VirtualAccounts, 0)

	for rows.Next() {
		var va models.VirtualAccounts
		err := rows.Scan(&va.VaId, &va.VaNum, &va.AccountNum, &va.VaLabel, &va.VaColor, &va.VaBalance, &va.CreatedAt, &va.UpdatedAt)

		if err != nil {
			return VirAcc, err
		}
		res = append(res, va)
	}

	return res, nil

}

//untuk ngecek input rekening apakah benar atau tidak.
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
	// AccountNumber := 0
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

func RevertVacBalanceToMainAccount(db *sql.DB, va models.VirtualAccounts) (err error) {
	acc, err := GetAccountByAccountNum(db, va.AccountNum)

	if err == nil {
		_, err = db.Exec("UPDATE accounts SET account_balance=$1 WHERE account_id=$2;", acc.AccountBalance+va.VaBalance, acc.AccountId)
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
	err = db.QueryRow("SELECT cust_id, account_num, cust_email FROM customers WHERE cust_id=$1;", id).Scan(&cust.CustId, &cust.AccountNum, &cust.CustEmail)
	return
}

func GetVacByAccountNum(db *sql.DB, accountNum string) (va models.VirtualAccounts, err error) {
	err = db.QueryRow("SELECT va_id, va_num, account_num, va_balance FROM virtual_accounts WHERE account_num=$1", accountNum).Scan(&va.VaId, &va.VaNum, &va.AccountNum, &va.VaBalance)
	return
}
