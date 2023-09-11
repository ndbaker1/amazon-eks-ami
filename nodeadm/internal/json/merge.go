package json

import (
	"strings"

	"github.com/RaveNoX/go-jsonmerge"
)

var (
	jsonPrefix = ""
	jsonIndent = strings.Repeat(" ", 4)
)

// Merge merges two JSON documents.
func Merge(a string, b string) (*string, error) {
	mergedBytes, _, err := jsonmerge.MergeBytesIndent([]byte(a), []byte(b), jsonPrefix, jsonIndent)
	if err != nil {
		return nil, err
	}
	merged := string(mergedBytes)
	return &merged, nil
}
