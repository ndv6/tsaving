package database

import (
	"database/sql"
	"fmt"

	"github.com/ndv6/tsaving/models"
)

func AllHistoryTransaction(db *sql.DB) (res []models.TransactionLogs, err error) {
	rows, err := db.Query("SELECT * FROM transaction_logs")

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var mtl models.TransactionLogs
		err = rows.Scan(&mtl.TlId, &mtl.AccountNum, &mtl.DestAccount, &mtl.FromAccount, &mtl.TranAmount, &mtl.Description, &mtl.CreatedAt)
		if err != nil {
			return
		}
		res = append(res, mtl)
	}
	return res, nil
}

func CustomerHistoryTransaction(db *sql.DB, accNum string) (res []models.TransactionLogs, err error) {
	rows, err := db.Query("SELECT * FROM transaction_logs WHERE account_num = $1", accNum)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var mtl models.TransactionLogs
		err = rows.Scan(&mtl.TlId, &mtl.AccountNum, &mtl.DestAccount, &mtl.FromAccount, &mtl.TranAmount, &mtl.Description, &mtl.CreatedAt)
		if err != nil {
			return
		}
		res = append(res, mtl)
	}
	return res, nil
}

func CustomerHistoryTransactionFiltered(db *sql.DB, accNum, search string) (res []models.TransactionLogs, err error) {
	rows, err := db.Query(`SELECT * FROM transaction_logs WHERE account_num = $1 AND (from_account like '%'||$2||'%' OR dest_account like '%'||$2||'%' OR description like '%'||$2||'%')`, accNum, search)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var mtl models.TransactionLogs
		err = rows.Scan(&mtl.TlId, &mtl.AccountNum, &mtl.DestAccount, &mtl.FromAccount, &mtl.TranAmount, &mtl.Description, &mtl.CreatedAt)
		if err != nil {
			return
		}
		res = append(res, mtl)
	}
	return res, nil
}

func GetActInActUserCount(db *sql.DB) (act, inact int, err error) {
	rows, err := db.Query("SELECT is_verified FROM customers")
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var is_verified bool
		err = rows.Scan(&is_verified)
		if err != nil {
			act = 0
			inact = 0
			return
		}
		if is_verified {
			act += 1
			continue
		}
		inact += 1
	}
	return
}

func GetTotalTransactionCount(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT COUNT(tl_id) FROM transaction_logs").Scan(&total)
	return
}

func GetNewUserToday(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at > current_date").Scan(&total)
	return
}

func GetNewUserYesterday(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at::date = current_date - 1").Scan(&total)
	return
}

func GetNewUserThisWeek(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at > current_date - 7").Scan(&total)
	return
}

func GetNewUserThisMonth(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at > date_trunc('month', CURRENT_DATE)").Scan(&total)
	return
}

func GetTransactionAmountToday(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT sum(tran_amount) FROM transaction_logs WHERE created_at > current_date").Scan(&total)
	return
}

func GetTransactionAmountYesterday(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT sum(tran_amount) FROM transaction_logs WHERE created_at::date = current_date - 1").Scan(&total)
	return
}

func GetTransactionAmountMonth(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT sum(tran_amount) FROM transaction_logs WHERE created_at > date_trunc('month', CURRENT_DATE)").Scan(&total)
	return
}

func GetTransactionByWeek(db *sql.DB) (res []models.TransactionMonth, err error) {
	rows, err := db.Query("SELECT ROW_NUMBER () OVER (ORDER BY extract(week from created_at)) as week, extract(week from created_at) as realweek, sum(tran_amount) as amount FROM transaction_logs where created_at > date_trunc('month', CURRENT_DATE) group by 2 order by 1 asc")

	defer rows.Close()

	for rows.Next() {
		var objTransactionByWeek models.TransactionMonth
		err = rows.Scan(&objTransactionByWeek.Week, &objTransactionByWeek.RealWeek, &objTransactionByWeek.Amount)
		if err != nil {
			return
		}
		res = append(res, objTransactionByWeek)
	}
	return res, nil
}

func GetLogTransactionToday(db *sql.DB) (res []models.TransactionLogs, err error) {
	rows, err := db.Query("SELECT tl_id, account_num, from_account, dest_account, tran_amount, description, created_at FROM transaction_logs WHERE created_at > current_date order by 1 desc ")
	defer rows.Close()

	for rows.Next() {
		var objTransactionLog models.TransactionLogs
		err = rows.Scan(&objTransactionLog.TlId, &objTransactionLog.AccountNum, &objTransactionLog.DestAccount, &objTransactionLog.FromAccount, &objTransactionLog.TranAmount, &objTransactionLog.Description, &objTransactionLog.CreatedAt)
		if err != nil {
			return
		}
		res = append(res, objTransactionLog)
	}
	return res, nil
}

func GetLogAdminToday(db *sql.DB) (res []models.LogAdmin, err error) {
	rows, err := db.Query("SELECT id, username, account_num, action, action_time FROM log_admins WHERE action_time > current_date order by 1 desc ")

	defer rows.Close()

	for rows.Next() {
		var objLogAdmin models.LogAdmin
		err = rows.Scan(&objLogAdmin.IDLogAdmin, &objLogAdmin.Username, &objLogAdmin.Action, &objLogAdmin.AccNum, &objLogAdmin.ActionTime)
		if err != nil {
			return
		}
		res = append(res, objLogAdmin)
	}
	return res, nil
}
