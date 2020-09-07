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
	CustPassword string    `json:"cust_password"`
	CustPict     string    `json:"cust_pict"`
	IsVerified   bool      `json:"is_verified"`
	Channel      string    `json:"channel"`
	CardNum      string    `json:"card_num"`
	Cvv          string    `json:"cvv"`
	Expired      time.Time `json:"expired"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsDeleted    time.Time `json:"is_deleted"`
}

type Card struct {
	Number  string
	Cvv     string
	Expired time.Time
}

func RegisterCustomer(db *sql.DB, objCustomer Customers, AccNum string, Pass string, cardNum string, cvv string, expired time.Time) error {
	_, err := db.Exec(`INSERT into customers (account_num, cust_name, cust_address, cust_phone, cust_email, cust_password, cust_pict, is_verified, channel, card_num, cvv, expired, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		AccNum,
		objCustomer.CustName,
		objCustomer.CustAddress,
		objCustomer.CustPhone,
		objCustomer.CustEmail,
		Pass,
		"",
		false,
		objCustomer.Channel,
		cardNum,
		cvv,
		expired,
		time.Now(),
		time.Now(),
	)
	return err
}

func LoginCustomer(db *sql.DB, email string, password string) (objCustomer Customers, err error) {
	err = db.QueryRow("SELECT cust_id, account_num, cust_name, cust_address, cust_phone, cust_email, cust_password, is_verified, expired, channel, created_at, updated_at from customers where is_verified = true and cust_email = ($1) and cust_password = ($2)", email, password).Scan(&objCustomer.CustId, &objCustomer.AccountNum, &objCustomer.CustName, &objCustomer.CustAddress, &objCustomer.CustPhone, &objCustomer.CustEmail, &objCustomer.CustPassword, &objCustomer.IsVerified, &objCustomer.Expired, &objCustomer.Channel, &objCustomer.CreatedAt, &objCustomer.UpdatedAt)
	return
}

func CheckLoginVerified(db *sql.DB, email string) (isVerified bool, err error) {
	err = db.QueryRow("SELECT is_verified FROM customers WHERE cust_email = ($1)", email).Scan(&isVerified)
	return
}

func GetDetailsCard(db *sql.DB, accountnum string) (objCustomer Customers, err error) {
	err = db.QueryRow("SELECT card_num, cvv, expired FROM customers WHERE account_num = ($1)", accountnum).Scan(&objCustomer.CardNum, &objCustomer.Cvv, &objCustomer.Expired)
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

func UpdateProfile(db *sql.DB, cus Customers) error {
	_, err := db.Exec("UPDATE customers SET cust_name = $1, cust_address = $2, cust_phone = $3, cust_email = $4, is_verified = $5, channel = $6, updated_at = NOW() WHERE cust_id = $7", cus.CustName, cus.CustAddress, cus.CustPhone, cus.CustEmail, cus.IsVerified, cus.Channel, cus.CustId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCustomerPicture(db *sql.DB, path string, id int) error {
	_, err := db.Exec("UPDATE customers SET cust_pict = $1, updated_at = NOW() WHERE cust_id = $2", path, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCustomerPassword(db *sql.DB, pass string, id int) error {
	_, err := db.Exec("UPDATE customers SET cust_password = $1, updated_at = NOW() WHERE cust_id = $2", pass, id)
	if err != nil {
		return err
	}
	return nil
}

func IsOldPasswordCorrect(db *sql.DB, pass string, id int) (bool, error) {
	var res bool
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM customers WHERE cust_password = $1 AND cust_id = $2)", pass, id)
	err := row.Scan(&res)
	if err != nil {
		return true, err
	}
	return res, nil
}

func IsEmailExist(db *sql.DB, email string, id int) (bool, error) {
	var res bool
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM customers WHERE cust_email = $1 AND cust_id <> $2)", email, id)
	err := row.Scan(&res)
	if err != nil {
		return true, err
	}
	return res, nil
}

func IsEmailChanged(db *sql.DB, email string, id int) (bool, error) {
	var res bool
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM customers WHERE cust_email = $1 AND cust_id = $2)", email, id)
	err := row.Scan(&res)
	if err != nil {
		return true, err
	}
	return !res, nil
}

func IsPhoneExist(db *sql.DB, phone string, id int) (bool, error) {
	var res bool
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM customers WHERE cust_phone = $1 AND cust_id <> $2)", phone, id)
	err := row.Scan(&res)
	if err != nil {
		return true, err
	}
	return res, nil
}
