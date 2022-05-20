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
 * @param delimiters : []string - Set of delimiters to use. Order defines precedence.
 * @param includeDelimiter : bool - True if the delimiter should be included in the token.
 * @return []string - Resulting tokens.
 */
func InlineParse(txt string, delimiters []string, includeDelimiter bool) []string {
	tokens := []string{}

	tempToken := ""

	for i := 0; i < len(txt); i++ {
		char := string(txt[i])
		isDelimiter := false

		// Check if the current character is forming a delimiter
		for _, delimiter := range delimiters {
			delimiterSize := len(delimiter)

			// Check if the following characters are the delimiter
			if i+delimiterSize <= len(txt) && txt[i:i+delimiterSize] == delimiter {

				// Dump the current token if it is not empty
				if tempToken != "" {
					tokens = append(tokens, tempToken)
					// Reset the current token
					tempToken = ""
				}

				// Check if the delimiter should be included
				if includeDelimiter {
					tokens = append(tokens, delimiter)
				}

				// Skip the delimiter characters
				i += delimiterSize - 1
				isDelimiter = true
				break
			}
		}

		// Did not find a delimiter, add the character to the token
		if !isDelimiter {
			tempToken += char
		}

	}

	// Dump the last token if it is not empty
	if tempToken != "" {
		tokens = append(tokens, tempToken)
	}

	// Remove certain words for strings
	result := []string{}
	removeEmpty := true
	for _, word := range tokens {

		if word == "\"" {
			removeEmpty = !removeEmpty
		}

		if !removeEmpty || (word != " " && word != "\t") {
			result = append(result, word)
		}
	}

	return result
}

/**
 * Replace the special characters in a string. This is used to replace escaped characters like \n and \t.
 * @param txt : string - The text to parse.
 * @return string - The formatted text.
 */
func HandleEscapeCharacters(txt string) string {
	return strings.Replace(txt, "\\\"", "\"", -1)
}
