package utils

import (
	"strconv"
)

// Atoi64 is shorthand for ParseInt(s, 10, 0).
func Atoi64(s string) (int64, error) {
	i64, err := strconv.ParseInt(s, 10, 0)
	return i64, err
}
