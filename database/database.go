package database

import (
	"database/sql"
	"fmt"

	"github.com/ndv6/tsaving/models"

	_ "github.com/lib/pq"
)

func GetDatabaseConnection(dbCfg models.DatabaseConfig) (db *sql.DB, err error) {
	db, err = sql.Open(dbCfg.DbType, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.DbName))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return
}
