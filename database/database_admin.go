package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/models"
)

type AdminDatabaseHandler struct {
	db *sql.DB
}

func NewAdminDatabaseHandler(db *sql.DB) *AdminDatabaseHandler {
	return &AdminDatabaseHandler{
		db,
	}
}

func (adm *AdminDatabaseHandler) SendMail(w http.ResponseWriter, OTPEmail string, cusEmail string) (err error) {
	requestBody, err := json.Marshal(map[string]string{
		"email": cusEmail,
		"token": OTPEmail,
	})

	if err != nil {
		return
	}

	_, err = http.Post(constants.TnotifLocal, constants.Json, bytes.NewBuffer(requestBody))
	if err != nil {
		return
	}

	return
}

func (adh *AdminDatabaseHandler) EditCustomerData(customerData models.Customers, adminUsername string) (err error) {
	tx, err := adh.db.Begin()
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE customers SET cust_phone = ($1), cust_email = ($2), is_verified=($3) WHERE account_num=($4) AND cust_phone IS DISTINCT FROM ($1) OR cust_email IS DISTINCT FROM ($2) OR is_verified IS DISTINCT FROM ($3)", customerData.CustPhone, customerData.CustEmail, customerData.IsVerified, customerData.AccountNum)
	if err != nil {
		tx.Rollback()
		return
	}

	editCustomerDataLog := models.LogAdmin{
		Username:   adminUsername,
		AccNum:     customerData.AccountNum,
		Action:     constants.EditCustomerData,
		ActionTime: time.Now(),
	}
	err = InsertLogAdminWithDbTransaction(tx, editCustomerDataLog, adminUsername)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}

func AllHistoryTransaction(db *sql.DB) (res []models.TransactionLogs, err error) {
	rows, err := db.Query("SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs")

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

func AllHistoryTransactionPaged(db *sql.DB, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs OFFSET $1 LIMIT 20", offset)

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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs").Scan(&count)
	if err != nil {
		return
	}

	return res, count, nil
}

func AllHistoryTransactionFilteredAccNum(db *sql.DB, search string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE (account_num like '%'||$1||'%' OR from_account like '%'||$1||'%' OR dest_account like '%'||$1||'%' OR description like '%'||$1||'%') OFFSET $2 LIMIT 20", search, offset)

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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs WHERE (account_num like '%'||$1||'%' OR from_account like '%'||$1||'%' OR dest_account like '%'||$1||'%' OR description like '%'||$1||'%')", search).Scan(&count)
	if err != nil {
		return
	}

	return res, count, nil
}

func AllHistoryTransactionFilteredDate(db *sql.DB, date string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE CAST(created_at as VARCHAR) like '%'||$1||'%' OFFSET $2 LIMIT 20", date, offset)

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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs WHERE CAST(created_at as VARCHAR) like '%'||$1||'%'", date).Scan(&count)
	if err != nil {
		return
	}

	return res, count, nil
}

func AllHistoryTransactionFilteredAccNumDate(db *sql.DB, search string, date string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE (account_num like '%'||$1||'%' OR from_account like '%'||$1||'%' OR dest_account like '%'||$1||'%' OR description like '%'||$1||'%') AND CAST(created_at as VARCHAR) like '%'||$2||'%' OFFSET $3 LIMIT 20", search, date, offset)

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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs WHERE (account_num like '%'||$1||'%' OR from_account like '%'||$1||'%' OR dest_account like '%'||$1||'%' OR description like '%'||$1||'%') AND CAST(created_at as VARCHAR) like '%'||$2||'%'", search, date).Scan(&count)
	if err != nil {
		return
	}

	return res, count, nil
}

func CustomerHistoryTransaction(db *sql.DB, accNum string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE account_num = $1 OFFSET $2 LIMIT 20", accNum, offset)
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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs WHERE account_num = $1", accNum).Scan(&count)
	if err != nil {
		return
	}

	return
}

func CustomerHistoryTransactionFiltered(db *sql.DB, accNum, search string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query(`SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE account_num = $1 AND (from_account like '%'||$2||'%' OR dest_account like '%'||$2||'%' OR description like '%'||$2||'%') OFFSET $3 LIMIT 20`, accNum, search, offset)
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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs WHERE account_num = $1 AND (from_account like '%'||$2||'%' OR dest_account like '%'||$2||'%' OR description like '%'||$2||'%')", accNum, search).Scan(&count)
	if err != nil {
		return
	}

	return
}

func CustomerHistoryTransactionDateFiltered(db *sql.DB, accNum, day, month, year string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	date := year + "-" + month + "-" + day
	rows, err := db.Query(`SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE account_num = $1 AND DATE(created_at) = $2 OFFSET $3 LIMIT 20`, accNum, date, offset)
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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs WHERE account_num = $1 AND DATE(created_at) = $2", accNum, date).Scan(&count)
	if err != nil {
		return
	}
	return
}

func CustomerHistoryTransactionAllFiltered(db *sql.DB, accNum, search, day, month, year string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	date := year + "-" + month + "-" + day

	rows, err := db.Query(`SELECT tl_id, account_num, dest_account, from_account, tran_amount, description, created_at FROM transaction_logs WHERE account_num = $1 AND (from_account like '%'||$2||'%' OR dest_account like '%'||$2||'%' OR description like '%'||$2||'%') AND DATE(created_at) = $3 OFFSET $4 LIMIT 20`, accNum, search, date, offset)
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

	err = db.QueryRow("SELECT COUNT(*) FROM transaction_logs WHERE account_num = $1 AND (from_account like '%'||$2||'%' OR dest_account like '%'||$2||'%' OR description like '%'||$2||'%') AND DATE(created_at) = $3", accNum, search, date).Scan(&count)
	if err != nil {
		return
	}

	return
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
	err = db.QueryRow("SELECT count(cust_id) FROM customers WHERE created_at::date = current_date").Scan(&total)
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
	err = db.QueryRow("SELECT COALESCE(sum(tran_amount),0) FROM transaction_logs WHERE created_at::date = current_date").Scan(&total)
	return
}

func GetTransactionAmountYesterday(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT COALESCE(sum(tran_amount),0) FROM transaction_logs WHERE created_at::date = current_date - 1").Scan(&total)
	return
}

func GetTransactionAmountMonth(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT COALESCE(sum(tran_amount),0) FROM transaction_logs WHERE created_at > date_trunc('month', CURRENT_DATE)").Scan(&total)
	return
}

func GetTransactionAmount(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT COALESCE(sum(tran_amount), 0) FROM transaction_logs").Scan(&total)
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
	rows, err := db.Query("SELECT tl_id, account_num, from_account, dest_account, tran_amount, description, created_at FROM transaction_logs WHERE created_at::date = current_date order by 1 desc ")
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
	rows, err := db.Query("SELECT id, username, account_num, action, action_time FROM log_admins WHERE action_time::date = current_date order by 1 desc ")

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

func GetTransactionAmountWeek(db *sql.DB) (total int, err error) {
	err = db.QueryRow("SELECT COALESCE(SUM(tran_amount), 0) FROM transaction_logs WHERE created_at > current_date - 7").Scan(&total)
	return
}

func GetTotalFromLog(db *sql.DB, cons string) (totalVa int, totalAmount int, err error) {
	err = db.QueryRow("SELECT COUNT(tl_id), COALESCE(SUM(tran_amount), 0) FROM transaction_logs WHERE description = $1", cons).Scan(&totalVa, &totalAmount)
	return
}

//SOFT DELETE
func SoftDeleteCustomer(db *sql.DB, AccNum string) (err error) {
	_, err = db.Exec("UPDATE customers SET is_deleted = $1 WHERE account_num = $2", time.Now(), AccNum)
	return
}
