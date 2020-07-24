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
		CreateDatabase(db)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return
}

func CreateDatabase(db *sql.DB) {
	fmt.Println("creating database")
	// var db_exist bool
	// _ = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1;", "db_tsaving").Scan(&db_exist)
	// fmt.Println("database exist")
	// fmt.Println(db_exist)
	// if !db_exist {
	db.Exec("CREATE DATABASE db_tsaving;")
	// fmt.Println("Success creating database db_tsaving")
	// }
}
