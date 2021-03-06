package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ndv6/tsaving/constants"
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
	err = tx.QueryRow("SELECT va_balance FROM virtual_accounts WHERE account_num = $1 AND va_num = $2 FOR UPDATE", accountNum, vacNum).Scan(&balanceVA)
	if err != nil {
		fmt.Print(err)
		err = errors.New(constants.CannotTransferVaToMain)
		tx.Rollback()
		return
	}

	if balanceVA < balanceInput {
		err = errors.New(constants.InvalidBalance)
		fmt.Print(err)
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE accounts SET account_balance = account_balance + $1 WHERE account_num = $2", balanceInput, accountNum)
	if err != nil {
		fmt.Print(err)
		err = errors.New(constants.CannotTransferVaToMain)
		tx.Rollback()
		return
	}
	_, err = tx.Exec("UPDATE virtual_accounts SET va_balance = va_balance - $1 WHERE va_num = $2", balanceInput, vacNum)
	if err != nil {
		fmt.Print(err)
		err = errors.New(constants.CannotTransferVaToMain)
		tx.Rollback()
		return
	}

	logDesc := constants.TransferToMainAccount

	//inpu transaction log
	tLogs := models.TransactionLogs{
		AccountNum:  accountNum,
		FromAccount: vacNum,
		DestAccount: accountNum,
		TranAmount:  balanceInput,
		Description: logDesc,
		CreatedAt:   time.Now(),
	}

	err = models.TransactionLog(tx, tLogs)
	if err != nil {
		fmt.Print(err)
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
	rows, err := db.Query("SELECT va_id, va_num, virtual_accounts.account_num, COALESCE(va_label,'') as va_label, COALESCE(va_color,'') as va_color, va_balance, COALESCE(virtual_accounts.created_at,now()) as created_at, COALESCE(virtual_accounts.updated_at,now()) as updated_at FROM virtual_accounts INNER JOIN customers ON virtual_accounts.account_num = customers.account_num WHERE cust_id = $1 ORDER BY virtual_accounts.va_balance DESC", id)
	if err != nil {
		return VirAcc, err
	}

	defer rows.Close()
	//defer -> yang harus dipanggil di akhir (biar ga lugitpa

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

func GetListVAAdmin(db *sql.DB, id int, page int) (virAcc []models.VirtualAccounts, count int, err error) {

	offset := (page - 1) * 20
	rows, err := db.Query("SELECT va_id, va_num, virtual_accounts.account_num, COALESCE(va_label,'') as va_label, COALESCE(va_color,'') as va_color, va_balance, COALESCE(virtual_accounts.created_at,now()) as created_at, COALESCE(virtual_accounts.updated_at,now()) as updated_at FROM virtual_accounts INNER JOIN customers ON virtual_accounts.account_num = customers.account_num WHERE cust_id = $1 ORDER BY virtual_accounts.created_at OFFSET $2 LIMIT 20", id, offset)
	if err != nil {
		return
	}

	var accNum string

	err = db.QueryRow("SELECT account_num FROM customers WHERE cust_id = $1", id).Scan(&accNum)
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM virtual_accounts WHERE account_num = $1", accNum).Scan(&count)
	if err != nil {
		return
	}

	defer rows.Close()
	//defer -> yang harus dipanggil di akhir (biar ga lugitpa

	// res := make([]models.VirtualAccounts, 0)

	for rows.Next() {
		var va models.VirtualAccounts
		err = rows.Scan(&va.VaId, &va.VaNum, &va.AccountNum, &va.VaLabel, &va.VaColor, &va.VaBalance, &va.CreatedAt, &va.UpdatedAt)
		if err != nil {
			return
		}
		virAcc = append(virAcc, va)
	}

	return
}

func GetListVAAdminFilter(db *sql.DB, id int, color string, page int) (virAcc []models.VirtualAccounts, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT va_id, va_num, virtual_accounts.account_num, va_label, va_color, va_balance, virtual_accounts.created_at, virtual_accounts.updated_at FROM virtual_accounts INNER JOIN customers ON virtual_accounts.account_num = customers.account_num WHERE cust_id = $1 AND va_color = $2 ORDER BY virtual_accounts.created_at OFFSET $3 LIMIT 20", id, color, offset)
	if err != nil {
		return
	}

	var accNum string

	err = db.QueryRow("SELECT account_num FROM customers WHERE cust_id = $1", id).Scan(&accNum)
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM virtual_accounts WHERE account_num = $1 AND va_color = $2", accNum, color).Scan(&count)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var va models.VirtualAccounts
		err = rows.Scan(&va.VaId, &va.VaNum, &va.AccountNum, &va.VaLabel, &va.VaColor, &va.VaBalance, &va.CreatedAt, &va.UpdatedAt)
		if err != nil {
			return
		}
		virAcc = append(virAcc, va)
	}

	return
}

//untuk ngecek input rekening apakah benar atau tidak.
func CheckAccountVA(db *sql.DB, VaNum string, id int) (err error) {
	var exist bool
	err = db.QueryRow("SELECT EXISTS(SELECT va_num FROM virtual_accounts INNER JOIN customers ON virtual_accounts.account_num = customers.account_num WHERE va_num = $1 AND cust_id = $2)", VaNum, id).Scan(&exist)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New(constants.InvalidVA)
		return
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
	rows, err := db.Query("SELECT va_num FROM virtual_accounts WHERE account_num = $1 ORDER BY va_num DESC LIMIT 1", accNum)
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

func RevertVacBalanceToMainAccount(trx *sql.Tx, va models.VirtualAccounts) (err error) {
	if err == nil {
		_, err = trx.Exec(fmt.Sprintf("UPDATE accounts SET account_balance = account_balance + subquery.va_balance FROM (SELECT va_balance FROM virtual_accounts WHERE va_num = '%s') as subquery WHERE account_num = '%s'; DELETE FROM virtual_accounts WHERE va_num = '%s';", va.VaNum, va.AccountNum, va.VaNum))
	}
	return
}

func GetVacByAccountNum(trx *sql.Tx, accountNum string, vaNum string) (va models.VirtualAccounts, err error) {
	err = trx.QueryRow("SELECT va_id, va_num, account_num, va_balance FROM virtual_accounts WHERE account_num=$1 AND va_num = $2 FOR UPDATE;", accountNum, vaNum).Scan(&va.VaId, &va.VaNum, &va.AccountNum, &va.VaBalance)
	return
}

func GetVaNumber(db *sql.DB, vaNum string) (va models.VirtualAccounts, err error) {
	err = db.QueryRow("SELECT va_num FROM virtual_accounts WHERE va_num=$1", vaNum).Scan(&va.VaNum)
	return
}
