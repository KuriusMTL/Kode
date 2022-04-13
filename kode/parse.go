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
 * Parse the individual tokens of a single line of code.
 * @param txt : string - The code to parse.
 * @return []string - The parsed code.
 */
func InlineParse(txt string, delimiters string) []string {

	word := ""
	words := []string{}

	for _, char := range txt {
		// Add char if it isn't one of the delimiters

		if !strings.ContainsAny(string(char), delimiters) {
			word += string(char)
			continue
		}

		if word != "" {
			words = append(words, word)
			word = ""
		}

		// Add the delimiter to the words array
		words = append(words, string(char))

	}

	// Add the last word to the words array
	if word != "" {
		words = append(words, word)
	}

	// Remove empty words
	for i, word := range words {
		if word == " " {
			words = append(words[:i], words[i+1:]...)
		}
	}

	return words
}
