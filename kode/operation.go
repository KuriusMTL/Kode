package kode

import (
	"errors"
	"math"
	"strings"
)

/**
 * Replace appropriate substraction operators in an expression with negation operators.
 * @param tokens : []string - The tokens to evaluate.
 * @return []string - The tokens with appropriate substraction operators replaced with negation operators.
 */
func CheckForNegation(tokens []string) []string {

	betweenQuotes := false

	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "-" && !betweenQuotes {
			if i == 0 {
				tokens[i] = "¬"
			} else if isOperator(tokens[i-1]) {
				tokens[i] = "¬"
			}
		} else if tokens[i] == "\"" {
			betweenQuotes = !betweenQuotes
		}
	}
	return tokens
}

/**
 * Evaluate the precedence of an operator.
 * @param op : string - The operator to evaluate.
 * @return int - The precedence of the operator.
 */
func OperatorPrecedence(op string) int {
	switch op {
	case "or", "and":
		return 1
	case "is", "==", "!=", ">", "<", ">=", "<=", "not":
		return 2
	case "+", "-":
		return 3
	case "*", "/":
		return 4
	case "^", "%":
		return 5
	case "¬":
		return 6
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
	case "+", "-", "*", "/", "¬", "^", "%", "==", "!=", ">", "<", ">=", "<=", "is", "not", "or", "and":
		return true
	default:
		return false
	}
}

func ApplyOperator(op string, val1 Variable, val2 Variable) (Variable, error) {

	switch op {
	case "+", "or":
		return val1.Add(&val2)
	case "-":
		return val1.Sub(&val2)
	case "*", "and":
		return val1.Mult(&val2)
	case "/":
		return val1.Div(&val2)
	case "^":
		return val1.Pow(&val2)
	case "%":
		return val1.Mod(&val2)
	case "¬", "not":
		return val2.Neg()
	case "==", "is":
		return val1.Equal(&val2)
	case "!=":
		return val1.NotEqual(&val2)
	case ">":
		return val1.Greater(&val2)
	case "<":
		return val1.Less(&val2)
	case ">=":
		return val1.GreaterEqual(&val2)
	case "<=":
		return val1.LessEqual(&val2)
	default:
		return Variable{}, errors.New("Error: Invalid operator (" + op + ")")
	}
}

func (val1 *Variable) Add(val2 *Variable) (Variable, error) {
	switch (*val1).Type {
	case "int":

		if (*val2).Type == "int" {
			return Variable{Type: "int", Value: (*val1).Value.(int64) + (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "float", Value: float64((*val1).Value.(int64)) + (*val2).Value.(float64)}, nil
		} else {
			break
		}

	case "float":

		if (*val2).Type == "int" {
			return Variable{Type: "float", Value: (*val1).Value.(float64) + float64((*val2).Value.(int64))}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "float", Value: (*val1).Value.(float64) + (*val2).Value.(float64)}, nil
		} else {
			break
		}

	case "string":

		if (*val2).Type == "string" {
			return Variable{Type: "string", Value: (*val1).Value.(string) + (*val2).Value.(string)}, nil
		} else if isArrayType((*val2).Type) {
			return CreateVariable(append([]Variable{(*val1)}, val2.Value.([]Variable)...)), nil
		} else {
			break
		}

	case "bool":

		if (*val2).Type == "bool" {
			return Variable{Type: "bool", Value: (*val1).Value.(bool) || (*val2).Value.(bool)}, nil
		} else {
			break
		}

	default:

		// Array type
		if isArrayType((*val1).Type) && isArrayType((*val2).Type) && (*val1).Type == (*val2).Type {
			return CreateVariable(append(val1.Value.([]Variable), val2.Value.([]Variable)...)), nil
		}

		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Invalid type (" + (*val1).Type + " + " + (*val2).Type + ") operation with addition")
}

func (val1 *Variable) Sub(val2 *Variable) (Variable, error) {
	switch (*val1).Type {
	case "int":

		if (*val2).Type == "int" {

			return Variable{Type: "int", Value: (*val1).Value.(int64) - (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {

			return Variable{Type: "float", Value: float64((*val1).Value.(int64)) - (*val2).Value.(float64)}, nil
		} else {
			break
		}

	case "float":

		if (*val2).Type == "int" {

			return Variable{Type: "float", Value: (*val1).Value.(float64) - float64((*val2).Value.(int64))}, nil
		} else if (*val2).Type == "float" {

			return Variable{Type: "float", Value: (*val1).Value.(float64) - (*val2).Value.(float64)}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Invalid type (" + (*val1).Type + " - " + (*val2).Type + ") operation with substraction")
}

func (val1 *Variable) Mult(val2 *Variable) (Variable, error) {
	switch (*val1).Type {
	case "int":

		if (*val2).Type == "int" {
			return Variable{Type: "int", Value: (*val1).Value.(int64) * (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "float", Value: float64((*val1).Value.(int64)) * (*val2).Value.(float64)}, nil
		} else {
			break
		}

	case "float":

		if (*val2).Type == "int" {
			return Variable{Type: "float", Value: (*val1).Value.(float64) * float64((*val2).Value.(int64))}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "float", Value: (*val1).Value.(float64) * (*val2).Value.(float64)}, nil
		} else {
			break
		}

	case "string":

		if (*val2).Type == "int" {
			return Variable{Type: "string", Value: strings.Repeat((*val1).Value.(string), int((*val2).Value.(int64)))}, nil
		} else {
			break
		}

	case "bool":

		if (*val2).Type == "bool" {
			return Variable{Type: "bool", Value: (*val1).Value.(bool) && (*val2).Value.(bool)}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Invalid type (" + (*val1).Type + " * " + (*val2).Type + ") operation with multiplication")
}

func (val1 *Variable) Div(val2 *Variable) (Variable, error) {
	switch (*val1).Type {
	case "int":

		if (*val2).Type == "int" {

			if (*val2).Value.(int64) == 0 {
				return Variable{}, errors.New("Error: Division by zero")
			}

			return Variable{Type: "int", Value: (*val1).Value.(int64) / (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {

			if (*val2).Value.(float64) == 0 {
				return Variable{}, errors.New("Error: Division by zero")
			}

			return Variable{Type: "float", Value: float64((*val1).Value.(int64)) / (*val2).Value.(float64)}, nil
		} else {
			break
		}

	case "float":

		if (*val2).Type == "int" {

			if (*val2).Value.(int64) == 0 {
				return Variable{}, errors.New("Error: Division by zero")
			}

			return Variable{Type: "float", Value: (*val1).Value.(float64) / float64((*val2).Value.(int))}, nil
		} else if (*val2).Type == "float" {

			if (*val2).Value.(float64) == 0 {
				return Variable{}, errors.New("Error: Division by zero")
			}

			return Variable{Type: "float", Value: (*val1).Value.(float64) / (*val2).Value.(float64)}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Invalid type (" + (*val1).Type + " / " + (*val2).Type + ") operation with division")
}

func (val1 *Variable) Pow(val2 *Variable) (Variable, error) {
	switch (*val1).Type {
	case "int":

		if (*val2).Type == "int" {
			return Variable{Type: "int", Value: int64(math.Pow(float64((*val1).Value.(int64)), float64((*val2).Value.(int64))))}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "float", Value: math.Pow(float64((*val1).Value.(int64)), (*val2).Value.(float64))}, nil
		} else {
			break
		}

	case "float":

		if (*val2).Type == "int" {
			return Variable{Type: "float", Value: math.Pow((*val1).Value.(float64), float64((*val2).Value.(int64)))}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "float", Value: math.Pow((*val1).Value.(float64), (*val2).Value.(float64))}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Invalid type (" + (*val1).Type + " ^ " + (*val2).Type + ") operation with exponent")
}

func (val1 *Variable) Mod(val2 *Variable) (Variable, error) {
	switch (*val1).Type {
	case "int":

		if (*val2).Type == "int" {

			if (*val2).Value.(int64) == 0 {
				return Variable{}, errors.New("Error: Modulo by zero")
			}

			return Variable{Type: "int", Value: (*val1).Value.(int64) % (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {

			if (*val2).Value.(float64) == 0 {
				return Variable{}, errors.New("Error: Modulo by zero")
			}

			k := math.Floor(float64((*val1).Value.(int64)) / (*val2).Value.(float64))

			return Variable{Type: "float", Value: float64((*val1).Value.(int64)) - (*val2).Value.(float64)*k}, nil
		} else {
			break
		}

	case "float":

		if (*val2).Type == "int" {

			if (*val2).Value.(int64) == 0 {
				return Variable{}, errors.New("Error: Modulo by zero")
			}

			k := math.Floor((*val1).Value.(float64) / float64((*val2).Value.(int64)))

			return Variable{Type: "float", Value: (*val1).Value.(float64) - float64((*val2).Value.(int64))*k}, nil
		} else if (*val2).Type == "float" {

			if (*val2).Value.(float64) == 0 {
				return Variable{}, errors.New("Error: Modulo by zero")
			}

			k := math.Floor((*val1).Value.(float64) / val2.Value.(float64))

			return Variable{Type: "float", Value: (*val1).Value.(float64) - (*val2).Value.(float64)*k}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Invalid type (" + (*val1).Type + " % " + (*val2).Type + ") operation with modulo")
}

func (val1 *Variable) Neg() (Variable, error) {

	switch (*val1).Type {
	case "int":

		return Variable{Type: "int", Value: -(*val1).Value.(int64)}, nil

	case "float":

		return Variable{Type: "float", Value: -(*val1).Value.(float64)}, nil

	case "bool":

		return Variable{Type: "bool", Value: !(*val1).Value.(bool)}, nil

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Invalid type (" + (*val1).Type + ") operation with negation")
}

func (val1 *Variable) Equal(val2 *Variable) (Variable, error) {

	switch (*val1).Type {
	case "int":
		if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(int64) == (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: float64((*val1).Value.(int64)) == (*val2).Value.(float64)}, nil
		} else {
			break
		}
	case "float":
		if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) == (*val2).Value.(float64)}, nil
		} else if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) == float64((*val2).Value.(int64))}, nil
		} else {
			break
		}

	case "string":
		if (*val2).Type == "string" {
			return Variable{Type: "bool", Value: (*val1).Value.(string) == (*val2).Value.(string)}, nil
		} else {
			break
		}

	case "bool":
		if (*val2).Type == "bool" {
			return Variable{Type: "bool", Value: (*val1).Value.(bool) == (*val2).Value.(bool)}, nil
		} else {
			break
		}
	case "null":
		if (*val2).Type == "null" {
			return Variable{Type: "bool", Value: true}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Cannot compare " + (*val1).Type + " with " + (*val2).Type + " type")
}

func (val1 *Variable) NotEqual(val2 *Variable) (Variable, error) {
	result, err := val1.Equal(val2)
	if err != nil {
		return Variable{}, err
	} else {
		return Variable{Type: "bool", Value: !(result.Value.(bool))}, nil
	}
}

func (val1 *Variable) Greater(val2 *Variable) (Variable, error) {

	switch (*val1).Type {
	case "int":
		if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(int64) > (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: float64((*val1).Value.(int64)) > (*val2).Value.(float64)}, nil
		} else {
			break
		}
	case "float":
		if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) > (*val2).Value.(float64)}, nil
		} else if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) > float64((*val2).Value.(int64))}, nil
		} else {
			break
		}

	case "string":
		if (*val2).Type == "string" {
			return Variable{Type: "bool", Value: len((*val1).Value.(string)) > len((*val2).Value.(string))}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Cannot compare " + (*val1).Type + " with " + (*val2).Type + " type")
}

func (val1 *Variable) Less(val2 *Variable) (Variable, error) {

	switch (*val1).Type {
	case "int":
		if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(int64) < (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: float64((*val1).Value.(int64)) < (*val2).Value.(float64)}, nil
		} else {
			break
		}
	case "float":
		if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) < (*val2).Value.(float64)}, nil
		} else if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) < float64((*val2).Value.(int64))}, nil
		} else {
			break
		}

	case "string":
		if (*val2).Type == "string" {
			return Variable{Type: "bool", Value: len((*val1).Value.(string)) < len((*val2).Value.(string))}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Cannot compare " + (*val1).Type + " with " + (*val2).Type + " type")
}

func (val1 *Variable) GreaterEqual(val2 *Variable) (Variable, error) {

	switch (*val1).Type {
	case "int":
		if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(int64) >= (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: float64((*val1).Value.(int64)) >= (*val2).Value.(float64)}, nil
		} else {
			break
		}
	case "float":
		if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) >= (*val2).Value.(float64)}, nil
		} else if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) >= float64((*val2).Value.(int64))}, nil
		} else {
			break
		}

	case "string":
		if (*val2).Type == "string" {
			return Variable{Type: "bool", Value: len((*val1).Value.(string)) >= len((*val2).Value.(string))}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Cannot compare " + (*val1).Type + " with " + (*val2).Type + " type")
}

func (val1 *Variable) LessEqual(val2 *Variable) (Variable, error) {

	switch (*val1).Type {
	case "int":
		if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(int64) <= (*val2).Value.(int64)}, nil
		} else if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: float64((*val1).Value.(int64)) <= (*val2).Value.(float64)}, nil
		} else {
			break
		}
	case "float":
		if (*val2).Type == "float" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) <= (*val2).Value.(float64)}, nil
		} else if (*val2).Type == "int" {
			return Variable{Type: "bool", Value: (*val1).Value.(float64) <= float64((*val2).Value.(int64))}, nil
		} else {
			break
		}

	case "string":
		if (*val2).Type == "string" {
			return Variable{Type: "bool", Value: len((*val1).Value.(string)) <= len((*val2).Value.(string))}, nil
		} else {
			break
		}

	default:
		break
	}

	// If incompatible types, return error
	return Variable{}, errors.New("Error: Cannot compare " + (*val1).Type + " with " + (*val2).Type + " type")
}
