package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/xrash/smetrics"
)

// Compares the actual output of a command to the expected output of a command.
func CompareCommandOutputs(
	actualOutput string,
	expectedOutput string,
	expectedSimilarity float64,
	expectedRegex *regexp.Regexp,
	expectedOutputLanguage string,
) (float64, error) {
	if expectedRegex != nil {
		if !expectedRegex.MatchString(actualOutput) {
			return 0.0, fmt.Errorf(
				ui.ErrorMessageStyle.Render(
					fmt.Sprintf("Expected output does not match: %q.", expectedRegex),
				),
			)
		}

		return 0.0, nil
	}

	if strings.ToLower(expectedOutputLanguage) == "json" {
		logging.GlobalLogger.Debugf(
			"Comparing JSON strings:\nExpected: %s\nActual%s",
			expectedOutput,
			actualOutput,
		)
		results, err := lib.CompareJsonStrings(actualOutput, expectedOutput, expectedSimilarity)
		if err != nil {
			return results.Score, err
		}

		logging.GlobalLogger.Debugf(
			"Expected Similarity: %f, Actual Similarity: %f",
			expectedSimilarity,
			results.Score,
		)

		if !results.AboveThreshold {
			return results.Score, fmt.Errorf(
				ui.ErrorMessageStyle.Render(
					"Expected output does not match actual output. Got: %s\n Expected: %s"),
				actualOutput,
				expectedOutput,
			)
		}

		return results.Score, nil
	}

	// Default case, using similarity on non JSON block.
	score := smetrics.JaroWinkler(expectedOutput, actualOutput, 0.7, 4)

	if expectedSimilarity > score {
		return score, fmt.Errorf(
			ui.ErrorMessageStyle.Render(
				"Expected output does not match actual output.\nGot:\n%s\nExpected:\n%s\nExpected Score:%s\nActualScore:%s",
			),
			ui.VerboseStyle.Render(actualOutput),
			ui.VerboseStyle.Render(expectedOutput),
			ui.VerboseStyle.Render(fmt.Sprintf("%f", expectedSimilarity)),
			ui.VerboseStyle.Render(fmt.Sprintf("%f", score)),
		)
	}

	return score, nil
}
