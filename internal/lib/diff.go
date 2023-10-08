package lib

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

func GetDifferenceBetweenStrings(a, b string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(a, b, false)
	return dmp.DiffPrettyText(diffs)
}
