package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/models"
)

func GetListCustomers(db *sql.DB, page int, date string, keyword string) (list []models.Customers, total int, err error) {
	var where = ""
	if date != "" {
		where = where + " AND created_at::date = date '" + date + "' "
	}
	if keyword != "" {
		where = where + " AND (LOWER(account_num) LIKE LOWER('%" + keyword + "%') OR LOWER(cust_name) LIKE LOWER('%" + keyword + "%') OR LOWER(cust_address) LIKE LOWER('%" + keyword + "%') OR LOWER(cust_phone) LIKE LOWER('%" + keyword + "%') OR LOWER(cust_email) LIKE LOWER('%" + keyword + "%') OR LOWER(channel) LIKE LOWER('%" + keyword + "%') OR LOWER(card_num) LIKE LOWER('%" + keyword + "%') ) "
	}

	offset := (page - 1) * 20

	query := fmt.Sprintf("SELECT cust_id, COALESCE(account_num,'') as account_num, cust_name, cust_address, cust_phone, cust_email, cust_password, COALESCE(is_verified, false) as is_verified, COALESCE(channel,'') as channel, COALESCE(card_num,'') as card_num, COALESCE(cvv,'') as cvv, COALESCE(expired,now()) as expired, COALESCE(created_at,now()) as created_at, COALESCE(updated_at,now()) as updated_at, COALESCE(is_deleted,'1970-01-01 08:00:00') as is_deleted FROM customers WHERE 1 = 1 %v ORDER BY created_at DESC OFFSET %v LIMIT 20", where, strconv.Itoa(offset))

	rows, err := db.Query(query)
	if err != nil {
		err = errors.New(constants.CustomersNotFound)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var cus models.Customers
		err = rows.Scan(&cus.CustId, &cus.AccountNum, &cus.CustName, &cus.CustAddress, &cus.CustPhone, &cus.CustEmail, &cus.CustPassword, &cus.IsVerified, &cus.Channel, &cus.CardNum, &cus.Cvv, &cus.Expired, &cus.CreatedAt, &cus.UpdatedAt, &cus.IsDeleted)
		if err != nil {
			return
		}
		list = append(list, cus)
	}

	querytotal := "SELECT COUNT(cust_id) as total FROM customers WHERE 1 = 1 " + where
	err = db.QueryRow(querytotal).Scan(&total)
	if err != nil {
		err = errors.New(constants.CustomersNotFound)
		return
	}

	return
}

func CheckAccount(db *sql.DB, AccountNum string) (err error) {
	// AccountNumber := 0
	var exist bool
	err = db.QueryRow("SELECT EXISTS(SELECT account_num FROM customers WHERE account_num = $1)", AccountNum).Scan(&exist)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("invalid account number")
		return
	}
	return
}
