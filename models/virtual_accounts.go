package models

import (
	"database/sql"
	"time"
)

type VirtualAccounts struct {
	VaId       int    `json:"var_id"`
	VaNum      string `json:"va_num"`
	AccountNum string `json:"account_num"`
	VaBalance  int    `json:"va_balance"`
	VaColor    string `json:"va_color"`
	VaLabel    string `json:"va_label"`
	CreatedAt  time.Time
	UpdateAt   time.Time
}

//
func (va *VirtualAccounts) UpdateVacBalance(db *sql.DB, saldoInput float64, RekVac string) error {

	_, err := db.Exec("UPDATE virtual_accounts SET va_balance = va_balance - $1 WHERE va_num = '$2'", saldoInput, RekVac)

	return err
}

func (va *VirtualAccounts) UpdateMainBalance(db *sql.DB, saldoInput float64, accountNum string) error {
	_, err := db.Exec("UPDATE accounts WHERE account_balance = account_balance + $1 WHERE account_num = $2", saldoInput, accountNum)

	return err
}

func (va *VirtualAccounts) GetRekeningByVA(db *sql.DB, rekVA string) (NoRek string, err error) {
	err = db.QueryRow("SELECT account_num from virtual_accounts WHERE va_num = $1", rekVA).Scan(&NoRek)

	if err != nil {
		return "", err
	}

	return NoRek, nil
}
