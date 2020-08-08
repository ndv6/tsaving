package database

import (
	"database/sql"
	"time"

	"github.com/ndv6/tsaving/models"
)

func InsertLogAdmin(db *sql.DB, la models.LogAdmin, username string) (err error) {
	_, err = db.Exec("INSERT INTO log_admins (username,account_num,action,action_time) VALUES ($1, $2, $3, $4);",
		username,
		la.AccNum,
		la.Action,
		time.Now())
	return err
}

func GetLogAdmin(db *sql.DB, page int) (LogAdmin []models.LogAdmin, err error) {

	offset := (page - 1) * 20
	rows, err := db.Query("SELECT id,username,action,account_num,action_time FROM log_admins ORDER BY action_time OFFSET $1 LIMIT 20", offset)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var la models.LogAdmin
		err = rows.Scan(&la.IDLogAdmin, &la.Username, &la.Action, &la.AccNum, &la.ActionTime)
		if err != nil {
			return
		}
		LogAdmin = append(LogAdmin, la)
	}
	return
}

func GetFilteredDateLogAdmin(db *sql.DB, date string, offset int) (res []models.LogAdmin, err error) {
	rows, err := db.Query("SELECT id,username,action,account_num,action_time FROM log_admins WHERE action_time LIKE $1% ORDER BY action_time OFFSET $2 LIMIT 20", date, offset)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var mla models.LogAdmin
		err = rows.Scan(&mla.IDLogAdmin, &mla.Username, &mla.Action, &mla.AccNum, &mla.ActionTime)
		if err != nil {
			return
		}
		res = append(res, mla)
	}
	return res, nil
}

func GetFilteredUsernameLogAdmin(db *sql.DB, username string, offset int) (res []models.LogAdmin, err error) {
	rows, err := db.Query("SELECT id,username,action,account_num,action_time FROM log_admins WHERE username = $1 ORDER BY action_time OFFSET $2 LIMIT 20", username, offset)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var mla models.LogAdmin
		err = rows.Scan(&mla.IDLogAdmin, &mla.Username, &mla.Action, &mla.AccNum, &mla.ActionTime)
		if err != nil {
			return
		}
		res = append(res, mla)
	}
	return res, nil
}

func GetFilteredDateAndUsernameLogAdmin(db *sql.DB, date string, username string, offset int) (res []models.LogAdmin, err error) {
	rows, err := db.Query("SELECT id,username,action,account_num,action_time FROM log_admins WHERE DATE(action_time) = $1 AND username = $2 ORDER BY action_time LIMIT 20 OFFSET $3", date, username, offset)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var mla models.LogAdmin
		err = rows.Scan(&mla.IDLogAdmin, &mla.Username, &mla.Action, &mla.AccNum, &mla.ActionTime)
		if err != nil {
			return
		}
		res = append(res, mla)
	}
	return res, nil
}
