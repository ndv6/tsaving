package database

import (
	"database/sql"
	"errors"

	"github.com/ndv6/tsaving/constants"
	"github.com/ndv6/tsaving/models"
)

func GetListCustomers(db *sql.DB, page int) (list []models.Customers, err error) {
	offset := (page - 1) * 20
	rows, err := db.Query("SELECT COALESCE(account_num,'') as account_num, cust_name, cust_address, cust_phone, cust_email, cust_password, is_verified, COALESCE(channel,'') as channel, COALESCE(card_num,'') as card_num, COALESCE(cvv,'') as cvv, COALESCE(expired,now()) as expired, COALESCE(created_at,now()) as created_at, COALESCE(updated_at,now()) as updated_at, is_deleted FROM customers ORDER BY created_at DESC OFFSET $1 LIMIT 20", offset)
	if err != nil {
		err = errors.New(constants.CustomersNotFound)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var cus models.Customers
		err = rows.Scan(&cus.AccountNum, &cus.CustName, &cus.CustAddress, &cus.CustPhone, &cus.CustEmail, &cus.CustPassword, &cus.IsVerified, &cus.Channel, &cus.CardNum, &cus.Cvv, &cus.Expired, &cus.CreatedAt, &cus.UpdatedAt, &cus.IsDeleted)
		if err != nil {
			return
		}
		list = append(list, cus)
	}
	return
}
