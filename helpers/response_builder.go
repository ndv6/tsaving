package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/ndv6/tsaving/constants"
)

type ResponseBuilder struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//This function will immediately return a templated json response
func NewResponseBuilder(w http.ResponseWriter, isAPICallSuccess bool, message string, obj interface{}) (rw http.ResponseWriter, responseJson string, err error) {
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
	return
}
