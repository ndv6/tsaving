package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func GetDatabaseConnection(connectionConfig string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connectionConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return
}
