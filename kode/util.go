package kode

import (
	"strconv"
)

/**
 * Evaluate if a string is a number.
 * @param str : string - The string to evaluate.
 * @return bool - True if the string is a number.
 */
func IsNumber(str string) bool {
	// Check if parse float causes an error
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}
