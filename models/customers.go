package models

import (
	"database/sql"
	"time"
)

type Customers struct {
	CustId       int       `json:"cust_id"`
	AccountNum   string    `json:"account_num"`
	CustName     string    `json:"cust_name"`
	CustAddress  string    `json:"cust_address"`
	CustPhone    string    `json:"cust_phone"`
	CustEmail    string    `json:"cust_email"`
	CustPict     string    `json:"cust_pict"`
	CustPassword string    `json:"cust_password"`
	IsVerified   bool      `json:"is_verified"`
	Channel      string    `json:"channel"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func RegisterCustomer(db *sql.DB, objCustomer Customers, AccNum string) error {
	Create := time.Now()
	Update := time.Now()
	Verified := false
	_, err := db.Exec("INSERT into customers (account_num, cust_name, cust_address, cust_phone, cust_email, cust_password, is_verified, channel, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", AccNum,
		objCustomer.CustName,
		objCustomer.CustAddress,
		objCustomer.CustPhone,
		objCustomer.CustEmail,
		objCustomer.CustPassword,
		Verified,
		objCustomer.Channel,
		Create,
		Update,
	)
	return err
}

func LoginCustomer(db *sql.DB, email string, password string) (objCustomer Customers, err error) {
	err = db.QueryRow("SELECT cust_id, account_num, cust_name, cust_address, cust_phone, cust_email, cust_password, is_verified, channel, created_at, updated_at from customers where is_verified = true and cust_email = ($1) and cust_password = ($2)", email, password).Scan(&objCustomer.CustId, &objCustomer.AccountNum, &objCustomer.CustName, &objCustomer.CustAddress, &objCustomer.CustPhone, &objCustomer.CustEmail, &objCustomer.CustPassword, &objCustomer.IsVerified, &objCustomer.Channel, &objCustomer.CreatedAt, &objCustomer.UpdatedAt)
	return
}

func GetAccNumber(db *sql.DB, id int) (acc string, err error) {
	err = db.QueryRow("SELECT account_num FROM customers WHERE cust_id = $1", id).Scan(&acc)
	return
}

func GetProfile(db *sql.DB, id int) (Customers, error) {
	var cus Customers
	row := db.QueryRow("SELECT cust_id, account_num, cust_name, cust_address, cust_phone, cust_email, cust_pict, is_verified, channel, created_at, updated_at FROM customers WHERE cust_id = $1", id)
	err := row.Scan(&cus.CustId, &cus.AccountNum, &cus.CustName, &cus.CustAddress, &cus.CustPhone, &cus.CustEmail, &cus.CustPict, &cus.IsVerified, &cus.Channel, &cus.CreatedAt, &cus.UpdatedAt)
	if err != nil {
		return Customers{}, err
	}
	return cus, nil
}
