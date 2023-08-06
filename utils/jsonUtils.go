package utils

import "encoding/json"

func ConvertToJsonString(object interface{}) (string, error) {
	responseBytes, err := json.Marshal(object)
	return string(responseBytes), err
}
