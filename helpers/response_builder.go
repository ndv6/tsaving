package helpers

import (
	"encoding/json"
	"net/http"

	mid "github.com/ndv6/tsaving/api/middleware"
	"github.com/ndv6/tsaving/constants"
)

// Made by Joseph and refactoring is done by Vici which consist of constants message

type ResponseBuilder struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//This function will immediately return a templated json response
func NewResponseBuilder(w http.ResponseWriter, r *http.Request, isAPICallSuccess bool, message string, obj interface{}) (rw http.ResponseWriter, responseJson string, err error) {
	status := constants.Failed
	if isAPICallSuccess {
		status = constants.Success
	}

	if obj == nil {
		obj = make(map[string]string)
	}

	b, err := json.Marshal(ResponseBuilder{
		Status:  status,
		Message: message,
		Data:    obj,
	})

	if err != nil {
		b = []byte(`{}`)
	}

	responseJson = string(b)
	rw = w
	rw.Header().Set(constants.ContentType, constants.Json)
	mid.InsertApplicationLog(r, status, message)
	return
}
