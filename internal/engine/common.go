package engine

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
func compareCommandOutputs(
	actualOutput string,
	expectedOutput string,
	expectedSimilarity float64,
	expectedRegex *regexp.Regexp,
	expectedOutputLanguage string,
) error {
	if expectedRegex != nil {
		if !expectedRegex.MatchString(actualOutput) {
			return fmt.Errorf(
				ui.ErrorMessageStyle.Render(
					fmt.Sprintf("Expected output does not match: %q.", expectedRegex),
				),
			)
		}

		return nil
	}

	if strings.ToLower(expectedOutputLanguage) == "json" {
		logging.GlobalLogger.Debugf(
			"Comparing JSON strings:\nExpected: %s\nActual%s",
			expectedOutput,
			actualOutput,
		)
		results, err := lib.CompareJsonStrings(actualOutput, expectedOutput, expectedSimilarity)
		if err != nil {
			return err
		}

		if !results.AboveThreshold {
			return fmt.Errorf(
				ui.ErrorMessageStyle.Render("Expected output does not match actual output."),
			)
		}

		logging.GlobalLogger.Debugf(
			"Expected Similarity: %f, Actual Similarity: %f",
			expectedSimilarity,
			results.Score,
		)

		return nil
	}

	// Default case, using similarity on non JSON block.
	score := smetrics.JaroWinkler(expectedOutput, actualOutput, 0.7, 4)

	if expectedSimilarity > score {
		return fmt.Errorf(
			ui.ErrorMessageStyle.Render("Expected output does not match actual output."),
		)
	}

	return nil
}
