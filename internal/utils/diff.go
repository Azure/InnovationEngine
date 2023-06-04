package utils

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func GetDifferenceBetweenStrings(a, b string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(a, b, false)
	return fmt.Sprintf("%s", dmp.DiffPrettyText(diffs))
}
