package database

import (
	"database/sql"
	"fmt"
)

func GetActInActUserCount(db *sql.DB) (act, inact int, err error) {
	err = db.QueryRow("SELECT COUNT(cust_id) FROM customers WHERE is_verified=true").Scan(&act)
	if err != nil {
		return
	}
	fmt.Println(act)

	err = db.QueryRow("SELECT COUNT(cust_id) FROM customers WHERE is_verified=false").Scan(&inact)
	if err != nil {
		return
	}
	fmt.Println(inact)
	return
}

// func GetTotalTransactionCount(db *sql.DB) (total int, err error) {
// 	err = db.QueryRow("SELECT COUNT()")
// }
