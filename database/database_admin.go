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

func GetNewUserTodayYesterday(db *sql.DB) (objData models.CountData, err error) {
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at > current_date - 1").Scan(&objData.Total)
	return
}

func GetNewUserThisWeek(db *sql.DB) (objData models.CountData, err error) {
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at > current_date - 7").Scan(&objData.Total)
	return
}

func GetNewUserThisMonth(db *sql.DB) (objData models.CountData, err error) {
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at > date_trunc('month', CURRENT_DATE)").Scan(&objData.Total)
	return
}

func GetTransactionAmount(db *sql.DB) (objData models.CountData, err error) {
	err = db.QueryRow("SELECT sum(tran_amount) FROM transaction_logs WHERE created_at > current_date - 1 ").Scan(&objData.Total)
	return
}

func GetLogTransactionToday(db *sql.DB) (objTransactionLog models.TransactionLogs, err error) {
	err = db.QueryRow("SELECT tl_id, account_num, from_account, dest_account, tran_amount, description, created_at FROM transaction_logs WHERE created_at > current_date order by 1 desc ").Scan(&objTransactionLog.TlId, &objTransactionLog.AccountNum, &objTransactionLog.DestAccount, &objTransactionLog.FromAccount, &objTransactionLog.TranAmount, &objTransactionLog.Description, &objTransactionLog.CreatedAt)
	return
}

func GetLogAdminToday(db *sql.DB) (objLogAdmin models.LogAdmin, err error) {
	err = db.QueryRow("SELECT id, username, account_num, action, action_time FROM log_admins WHERE action_time > current_date order by 1 desc ").Scan(&objLogAdmin.IDLogAdmin, &objLogAdmin.Username, &objLogAdmin.Action, &objLogAdmin.AccNum, &objLogAdmin.ActionTime)
	return
}
