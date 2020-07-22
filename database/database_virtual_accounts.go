package database

import (
	"database/sql"

	"github.com/ndv6/tsaving/models"
)

func GetBalanceVA(vaNum string, db *sql.DB) (va models.VirtualAccounts, err error) {
	var balanceFloat float32
	err = db.QueryRow("SELECT va_balance FROM virtual_accounts WHERE va_num = ($1) ", vaNum).Scan(&balanceFloat)
	if err != nil {
		return
	}
	va.VaBalance = int(balanceFloat)
	return
}
