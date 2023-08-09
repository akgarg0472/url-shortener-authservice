package utils

import "encoding/json"

func ConvertToJsonString(object interface{}) (string, error) {
	responseBytes, err := json.Marshal(object)
	return string(responseBytes), err
}

func ConvertToJsonBytes(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}
