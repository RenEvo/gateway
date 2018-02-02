package env

import (
	"os"
	"strconv"
	"strings"
)

// Bool reads the supplied environmental variable and attempts to parse it.
//
// When the value is 1, true, or yes, will return true
// Any unknown values or empty, will return false
func Bool(name string) bool {
	value := os.Getenv(name)

	if value == "" {
		return false
	}

	if value == "1" {
		return true
	}

	if strings.EqualFold(value, "true") {
		return true
	}

	if strings.EqualFold(value, "yes") {
		return true
	}

	return false
}

// Int64 reads the supplied environmental variable and attempts to parse it.
//
// When the value is not a valid int64, this method will return a zero
func Int64(name string) int64 {
	value := os.Getenv(name)
	if value == "" {
		return 0
	}

	ivalue, _ := strconv.ParseInt(value, 10, 64)

	return ivalue
}
