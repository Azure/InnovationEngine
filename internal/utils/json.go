package utils

import (
	"encoding/json"
)

func OrderJsonFields(jsonStr string) (string, error) {
	expectedMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &expectedMap)
	if err != nil {
		return "", err
	}

	orderedJson, err := json.Marshal(expectedMap)
	if err != nil {
		return "", err
	}
	return string(orderedJson), nil
}
