package kode

import (
	"strings"
)

/**
 * Parse the individual lines of code.
 * @param txt : string - The code to parse.
 * @return []string - The parsed code.
 */
func LineParse(txt string) []string {
	return strings.Split(txt, "\n")
}

/**
 * ! DEPRECATED
 * Parse the individual tokens of a single line of code.
 * @param txt : string - The code to parse.
 * @return []string - The parsed code.
 */
// func InlineParse(txt string, delimiters string) []string {

// 	word := ""
// 	words := []string{}

// 	for _, char := range txt {
// 		// Add char if it isn't one of the delimiters

// 		if !strings.ContainsAny(string(char), delimiters) {
// 			word += string(char)
// 			continue
// 		}

// 		if word != "" {
// 			words = append(words, word)
// 			word = ""
// 		}

// 		// Add the delimiter to the words array
// 		words = append(words, string(char))

// 	}

// 	// Add the last word to the words array
// 	if word != "" {
// 		words = append(words, word)
// 	}

// 	// Remove empty words
// 	for i, word := range words {
// 		if word == " " {
// 			words = append(words[:i], words[i+1:]...)
// 		}
// 	}

// 	return words
// }

/**
 * Parse the individual tokens of a single line of code.
 * @param txt : string - The code to parse.
 * @return []string - The parsed code.
 */
func InlineParse(txt string, delimiters []string, includeDelimiter bool) []string {

	word := ""
	tempDelimiter := ""
	words := []string{}

	for _, char := range txt {

		count := IsComplete(delimiters, tempDelimiter+string(char))

		if count > 1 {
			tempDelimiter += string(char)
		} else if count == 1 {
			if word != "" {
				words = append(words, word)
			}
			word = ""
			if includeDelimiter {
				words = append(words, tempDelimiter+string(char))
			}
			tempDelimiter = ""
		} else {
			if IsInArray(delimiters, tempDelimiter) {
				if word != "" {
					words = append(words, word)
				}
				word = ""

				if includeDelimiter {
					words = append(words, tempDelimiter)
				}
				tempDelimiter = ""

				if IsComplete(delimiters, string(char)) >= 1 {
					tempDelimiter = string(char)
				} else {
					word = string(char)
				}

			} else {
				word += string(char)
			}
		}

	}

	// Add the last word to the words array
	if word != "" {
		words = append(words, word)
	}

	// Remove certain words
	result := []string{}
	removeEmpty := true
	for _, word := range words {

		if word == "\"" {
			removeEmpty = !removeEmpty
		}

		if !removeEmpty || word != " " {
			result = append(result, word)
		}
	}

	return result
}

func IsInArray(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func IsComplete(arr []string, str string) int {
	count := 0
	for _, a := range arr {
		if strings.HasPrefix(a, str) {
			count++
		}
	}
	return count
}
