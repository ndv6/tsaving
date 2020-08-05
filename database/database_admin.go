package database

import (
	"database/sql"
	"github.com/ndv6/tsaving/models"
	//"errors"
	//"time"

)

func AllHistoryTransaction(db *sql.DB)(res []models.TransactionLogs, err error) {
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