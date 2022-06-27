package lib

import "encoding/json"

func ErrorResponseHelper(requestId string, message string) []byte {
	result := map[string]string{
		"requestId": requestId,
		"message":   message,
	}
	b, _ := json.Marshal(result)
	return b
}
