package utils

import (
	"strconv"
)

// Atoi64 is shorthand for ParseInt(s, 10, 0).
func Atoi64(s string) (int64, error) {
	i64, err := strconv.ParseInt(s, 10, 0)
	return i64, err
}

// StringsToJSON Format String to JSON
func StringsToJSON(str string) string {
	rs := []rune(str)
	jsons := ""

	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}
