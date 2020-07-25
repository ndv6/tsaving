// Unit test for deposit API, made by Vici
package customers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ndv6/tsaving/constants"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/models"
)

const (
	testValidClientSecret    = "ASDF123456"
	testInvalidClientSecret  = "2bb34e46cf2d0c23bf2eca8564"
	testValidPartnerId       = 1
	testInvalidPartnerId     = 2
	testValidDepositAmount   = 100000
	testInvalidDepositAmount = -100000
	testIncompleteRequest    = `{}`
	testUnauthorizedRequest  = `{ 
									"balance_added": 1000000,
    								"account_number": "202007221",
    								"auth_code": "asdfghj",
									"client_id": 2
								}`
	testValidAuthorizedRequest = `{ 
									"balance_added": 1000000,
									"account_number": "202007221",
									"auth_code": "ba5cac8f9f665451b75b357b4bcc419720ce3eda25833a944e928d1870308569",
									"client_id": 1
								}`
)

type testTransactor struct {
}

func (trx testTransactor) DepositToMainAccountDatabaseAccessor(balanceToAdd int, accountNumber string, log models.TransactionLogs) error {
	return nil
}

func (trx testTransactor) LogTransaction(log models.TransactionLogs) error {
	return nil
}

type testPartnerInterface struct {
}

func (pi testPartnerInterface) GetSecret(id int) (string, error) {
	if id == testValidPartnerId {
		return testValidClientSecret, nil
	}
	return testInvalidClientSecret, nil
}

func trimResponseMessage(message string) string {
	splitStrings := strings.Split(message, "\"")
	return splitStrings[3]
}

func TestShouldReturnBadRequestError(t *testing.T) {
	partnerInterface := testPartnerInterface{}
	transactor := testTransactor{}

	mockRequest, err := http.NewRequest("POST", "/deposit", bytes.NewBuffer([]byte(testIncompleteRequest)))
	if err != nil {
		log.Fatal(err)
	}

	responseReader := httptest.NewRecorder()
	handler := customers.DepositToMainAccount(partnerInterface, transactor)
	handler.ServeHTTP(responseReader, mockRequest)

	if status := responseReader.Code; status != http.StatusBadRequest {
		t.Fatalf("Expected: %v, Received: %v", http.StatusBadRequest, status)
	}

	if errorMsg := trimResponseMessage(responseReader.Body.String()); errorMsg != constants.RequestHasInvalidFields {
		t.Fatalf("Expected: %v, Received: %v", constants.RequestHasInvalidFields, errorMsg)
	}
}

func TestShouldReturnUnauthorizedError(t *testing.T) {
	partnerInterface := testPartnerInterface{}
	transactor := testTransactor{}

	mockRequest, err := http.NewRequest("POST", "/deposit", bytes.NewBuffer([]byte(testUnauthorizedRequest)))
	if err != nil {
		log.Fatal(err)
	}

	responseReader := httptest.NewRecorder()
	handler := customers.DepositToMainAccount(partnerInterface, transactor)
	handler.ServeHTTP(responseReader, mockRequest)

	if status := responseReader.Code; status != http.StatusUnauthorized {
		t.Fatalf("Expected: %v, Received: %v", http.StatusUnauthorized, status)
	}

	if errorMsg := trimResponseMessage(responseReader.Body.String()); errorMsg != constants.UnauthorizedRequest {
		t.Fatalf("Expected: %v, Received: %v", constants.UnauthorizedRequest, errorMsg)
	}
}

func TestShouldDepositSuccess(t *testing.T) {
	partnerInterface := testPartnerInterface{}
	transactor := testTransactor{}

	mockRequest, err := http.NewRequest("POST", "/deposit", bytes.NewBuffer([]byte(testValidAuthorizedRequest)))
	if err != nil {
		log.Fatal(err)
	}

	responseReader := httptest.NewRecorder()
	handler := customers.DepositToMainAccount(partnerInterface, transactor)
	handler.ServeHTTP(responseReader, mockRequest)

	if status := responseReader.Code; status != http.StatusOK {
		t.Fatalf("Expected: %v, Received: %v", http.StatusOK, status)
	}

	if errorMsg := trimResponseMessage(responseReader.Body.String()); errorMsg != constants.Success {
		t.Fatalf("Expected: %v, Received: %v", constants.Success, errorMsg)
	}
}
