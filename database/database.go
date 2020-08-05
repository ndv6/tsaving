package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func DatabaseConnect(url string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return
}
