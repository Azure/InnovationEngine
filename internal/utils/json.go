package utils

import (
	"encoding/json"

	"github.com/xrash/smetrics"
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

// Compute the Jaro-Winkler score for two JSON strings. The score is computed
// by ordering the fields alphabetically and then comparing the strings using
// the Jaro-Winkler algorithm.
func ComputeJaroWinklerScore(actualJson string, expectedJson string) (float64, error) {
	actualOutput, err := OrderJsonFields(actualJson)
	if err != nil {
		return 0, err
	}

	expectedOutput, err := OrderJsonFields(expectedJson)
	if err != nil {
		return 0, err
	}

	return smetrics.JaroWinkler(expectedOutput, actualOutput, 0.7, 4), nil
}

// Compare two JSON strings by ordering the fields alphabetically and then
// comparing the strings using the Jaro-Winkler algorithm to compute a score.
// If the score is greater than the threshold, return true.
func CompareJsonStrings(actualJson string, expectedJson string, threshold float64) (bool, error) {
	score, err := ComputeJaroWinklerScore(actualJson, expectedJson)
	if err != nil {
		return false, err
	}
	return score > threshold, nil
}
