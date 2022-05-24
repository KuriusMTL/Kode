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

/**
 * Evaluate if a string is a boolean
 * @param str : string - The string to evaluate.
 * @return bool - True if the string is a boolean.
 */
func IsBoolean(str string) bool {
	return str == "true" || str == "false"
}

/**
 * While a queue is filled with a token, pop the next token.
 * @param queue : Queue - The queue to pop the next token from.
 * @return string - The next token.
 */
func SkipWhileTokenNext(queue *Queue, value string) {
	nextToken, _ := (*queue).Pop()
	for nextToken != nil && nextToken.(string) == value {
		nextToken, _ = (*queue).Pop()
	}
}

// ? Note by Eduard: I don't know if this is the best way to do this, but for now its works.
/**
 * Reassemble the tokens in a Queue into a string.
 * @param queue : Queue - The queue to reassemble.
 * @return string - The reassembled string.
 */
func InlineQueueToString(queue *Queue) string {
	value := ""
	isInString := false
	for !(*queue).IsEmpty() {
		tokenValue, _ := (*queue).Peek()

		// Enter or exist a string
		if tokenValue.(string) == "\"" {
			isInString = !isInString
		}

		// Exit if a comment is found.
		if tokenValue.(string) == "#" {
			break
		}

		// In order to keep the "#" in the queue if it exists, we only pop the next token now.
		(*queue).Pop()

		// If not in a string, add an additional space for proper parsing afterwards
		if !isInString {

			// ! No longer needed
			// peekedValue, _ := (*queue).Peek()

			// if tokenValue.(string) == "=" && peekedValue.(string) == "=" {
			// 	value += tokenValue.(string)
			// } else {
			// 	value += tokenValue.(string) + " "
			// }

			value += tokenValue.(string) + " "

		} else {
			value += tokenValue.(string)
		}
	}
	return value
}

/**
 * Copy a variable map.
 * @param map : map[string]*Variable - The map to copy.
 * @return map[string]*Variable - The copied map.
 */
func CopyVariableMap(originalMap map[string]*Variable) map[string]*Variable {
	// Create new map
	newMap := make(map[string]*Variable)
	// Copy values from the original map
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func CopyFunction(originalFunction *Function) *Function {
	// Create new function

	newVars := CopyVariableMap((*originalFunction).Variables)
	newFunction := &Function{
		Code:       (*originalFunction).Code,
		Return:     (*originalFunction).Return,
		Arguments:  (*originalFunction).Arguments,
		Variables:  newVars,
		IsInstance: (*originalFunction).IsInstance,
		Parent:     (*originalFunction).Parent,
	}
	return newFunction
}
