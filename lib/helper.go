package lib

import "encoding/json"

type ErrorResponse struct {
	RequestId string `json:"requestId" example:"b8fce75f-a603-4f30-8896-29bf54e6ce25"`
	Message   string `json:"message" example:"some error message"`
}

func ErrorResponseHelper(requestId string, message string) []byte {
	result := ErrorResponse{
		RequestId: requestId,
		Message:   message,
	}
	b, _ := json.Marshal(result)
	return b
}
