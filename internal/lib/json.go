package lib

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

type ComparisonResult struct {
	AboveThreshold bool
	Score          float64
}

// Compare two JSON strings by ordering the fields alphabetically and then
// comparing the strings using the Jaro-Winkler algorithm to compute a score.
// If the score is greater than the threshold, return true.
func CompareJsonStrings(actualJson string, expectedJson string, threshold float64) (ComparisonResult, error) {
	actualOutput, err := OrderJsonFields(actualJson)
	if err != nil {
		return ComparisonResult{}, err
	}

	expectedOutput, err := OrderJsonFields(expectedJson)
	if err != nil {
		return ComparisonResult{}, err
	}

	score := smetrics.Jaro(actualOutput, expectedOutput)

	return ComparisonResult{
		AboveThreshold: score >= threshold,
		Score:          score,
	}, nil
}
