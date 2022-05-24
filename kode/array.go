package kode

import (
	"errors"
	"math"
	"strings"
)

/**
 * Evaluate the type of an array variable.
 * @param value : interface{} - The array
 * @return string - The type of the array.
 */
func EvaluateArrayType(array []Variable) string {

	// If empty, return make generic array
	if len(array) == 0 {
		return "val[]"
	}

	currType := EvaluateType(array[0].Value)
	for i := 1; i < len(array); i++ {

		// If two different types, return generic array
		if EvaluateType(array[i].Value) != currType {
			return "val[]"
		}
	}

	// Every element is the same type, return array of that type
	return currType + "[]"
}

/**
 * Extract the values of a declared array.
 * @param queue : *Queue - The queue to extract the values from.
 * @return []Variable - Resulting array.
 * @return error - The error if one occurs.
 */
func (scope *Function) ExtractArrayValues(queue *Queue) ([]Variable, error) {

	// Extract the array's value
	nestedArray := 0
	closedArray := false
	isInString := false
	tmpVal := ""
	array := []Variable{}

	for !queue.IsEmpty() {
		nextToken, _ := queue.Pop()

		if nextToken.(string) == "[" && !isInString {
			nestedArray++
		}
		if nextToken.(string) == "]" && !isInString {
			if nestedArray == 0 {
				closedArray = true
				break
			}
			nestedArray--
		}
		if nextToken.(string) == "," && nestedArray == 0 && !isInString {
			// Evaluate the parameter
			//println(tmpVal)
			parameter, err := EvaluateExpression(scope, tmpVal)
			if err != nil {
				return []Variable{}, err
			}

			array = append(array, parameter)
			tmpVal = ""
			continue
		}

		if nextToken.(string) == "\"" {
			isInString = !isInString
		}

		if isInString {
			tmpVal += nextToken.(string)
		} else {
			tmpVal += nextToken.(string) + " "
		}
	}

	// Add last value
	if tmpVal != "" {
		//println(tmpVal)
		parameter, err := EvaluateExpression(scope, tmpVal)
		if err != nil {
			return []Variable{}, err
		}
		array = append(array, parameter)
	}

	if !closedArray {
		return []Variable{}, errors.New("Error: Missing closing bracket for array")
	}

	return array, nil
}

/**
 * Extract the specified array dimension from a variable declaration.
 * @param tokens : Queue - The queue to extract the dimension from.
 * @return int - The dimension.
 * @return error - The error if one occurs.
 */
func ExtractArrayDimensionFromDeclaration(tokens *Queue) (int, error) {
	// Check for array dimension
	// Check for square brackets
	dimension := 0.0
	squareBracket, _ := (*tokens).Peek()
	if squareBracket.(string) == "[" {
		dimension = 0.5
		(*tokens).Pop()
		for !(*tokens).IsEmpty() {
			nextBracket, _ := (*tokens).Peek()
			if nextBracket.(string) == "]" && squareBracket.(string) == "[" {
				dimension += 0.5
				squareBracket = nextBracket
				(*tokens).Pop()
			} else if nextBracket.(string) == "[" && squareBracket.(string) == "]" {
				squareBracket = nextBracket
				dimension += 0.5
				(*tokens).Pop()
			} else {
				break
			}
		}
	}

	// Update the type according to dimmension
	_, fraction := math.Modf(dimension)
	if fraction != 0 {
		return 0, errors.New("Error: Invalid array dimension decleration")
	}
	return int(dimension), nil
}

/**
 * Check if the type is an array type.
 * @param strType : string - The type to check.
 * @return bool - True if the type is an array type.
 */
func isArrayType(strType string) bool {
	return strings.Contains(strType, "[]")
}

/**
 * Check the size of an array or a string.
 * @param variable : *Variable - The array or string variable to check.
 * @return int64 - The size of the array or string.
 * @return error - The error if one occurs.
 */
func GetArraySize(variable *Variable) (int64, error) {
	if isArrayType(variable.Type) {

		size := int64(len((*variable).Value.([]Variable)))
		if size == 0 {
			return 0, errors.New("Error: Array is empty")
		}
		return size, nil

	} else if variable.Type == "string" {

		size := int64(len((*variable).Value.(string)))
		if size == 0 {
			return 0, errors.New("Error: String is empty")
		}
		return size, nil

	} else {
		return 0, errors.New("Error: Cannot get array size of non-array type")
	}
}
