package helpers

func checkBalance(target string, acc_number string, amount int) bool {

	if amount > 0 && amount <= 50000 {
		return true
	}
	return false
}
