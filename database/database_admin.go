package database

import (
	"database/sql"
	"fmt"
	"time"

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

func CustomerHistoryTransaction(db *sql.DB, accNum string, page int) (res []models.TransactionLogs, count int, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT * FROM transaction_logs WHERE account_num = $1 OFFSET $2 LIMIT 20", accNum, offset)
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
	rows, err := db.Query(`SELECT * FROM transaction_logs WHERE account_num = $1 AND (from_account like '%'||$2||'%' OR dest_account like '%'||$2||'%' OR description like '%'||$2||'%') OFFSET $3 LIMIT 20`, accNum, search, offset)
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
	rows, err := db.Query(`SELECT * FROM transaction_logs WHERE account_num = $1 AND DATE(created_at) = $2 OFFSET $3 LIMIT 20`, accNum, date, offset)
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

	rows, err := db.Query(`SELECT * FROM transaction_logs WHERE account_num = $1 AND (from_account like '%'||$2||'%' OR dest_account like '%'||$2||'%' OR description like '%'||$2||'%') AND DATE(created_at) = $3 OFFSET $4 LIMIT 20`, accNum, search, date, offset)
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

//SOFT DELETE
func SoftDeleteCustomer(db *sql.DB, AccNum string) (err error) {
	_, err = db.Exec("UPDATE customers SET is_deleted = $1 WHERE account_num = $2", time.Now(), AccNum)
	return
}
