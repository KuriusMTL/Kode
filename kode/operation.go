package kode

import (
	"errors"
	"math"
	"strconv"
)

/**
 * Evaluate the precedence of an operator.
 * @param op : string - The operator to evaluate.
 * @return int - The precedence of the operator.
 */
func OperatorPrecedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^", "%":
		return 3
	case "¬":
		return 4
	default:
		return 0
	}
}

/**
* Determine if a string is an operator.
* @param str : string - The string to evaluate.
* @return bool - True if the string is an operator. Otherwise, return false.
 */
func isOperator(op string) bool {
	switch op {
	case "+", "-", "*", "/", "¬", "^", "%":
		return true
	default:
		return false
	}
}

func ApplyOperator(op string, val1 string, val2 string) (string, error) {
	v1, _ := strconv.Atoi(val1)
	v2, _ := strconv.Atoi(val2)

	switch op {
	case "+":
		return strconv.Itoa(v1 + v2), nil
	case "-":
		return strconv.Itoa(v1 - v2), nil
	case "*":
		return strconv.Itoa(v1 * v2), nil
	case "/":
		if v2 == 0 {
			return "", errors.New("Error: Division by zero.")
		}
		return strconv.Itoa(v1 / v2), nil
	case "^":
		return strconv.Itoa(int(math.Pow(float64(v1), float64(v2)))), nil
	case "%":
		return strconv.Itoa(v1 % v2), nil
	case "¬":
		return strconv.Itoa(-v2), nil
	default:
		return "", errors.New("Error: Invalid operation or operation not implemented.")
	}
}

func CheckForNegation(tokens []string) []string {
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "-" {
			if i == 0 {
				tokens[i] = "¬"
			} else if isOperator(tokens[i-1]) {
				tokens[i] = "¬"
			}
		}
	}
	return tokens
}
