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

func CustomerHistoryTransaction(db *sql.DB, accNum string, page int) (res []models.TransactionLogs, err error) {
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
	return res, nil
}

func CustomerHistoryTransactionFiltered(db *sql.DB, accNum, search string, page int) (res []models.TransactionLogs, err error) {
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
	return res, nil
}

func CustomerHistoryTransactionDateFiltered(db *sql.DB, accNum, day, month, year string, page int) (res []models.TransactionLogs, err error) {
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
	return res, nil
}

func CustomerHistoryTransactionAllFiltered(db *sql.DB, accNum, search, day, month, year string, page int) (res []models.TransactionLogs, err error) {
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
