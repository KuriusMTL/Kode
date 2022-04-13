package kode

import (
	"errors"
	"strconv"
	"strings"
)

/**
 * Evaluate an expression.
 * @param str : string - The expression to evaluate.
 * @return string - The result of the expression.
 * @return error - The error if any.
 */
func EvaluateExpression(scope *Function, str string) (Variable, error) {

	// Tokenize the line
	tokens := InlineParse(str, " \t\n\r*/+-()¬^%")
	// Replace proper substractions with negation
	tokens = CheckForNegation(tokens)

	// Add the tokens to the queue
	queue := Queue{}
	for _, token := range tokens {
		queue.Push(token)
	}

	values := Stack{}    // Store values to evaluate
	operators := Stack{} // Store operators to evaluate values

	// Loop through the tokens
	for !queue.IsEmpty() {

		// Get the next token
		token, _ := queue.Pop()

		// Check if the token is a number
		if IsNumber(token.(string)) {

			// Is float?
			if strings.Contains(token.(string), ".") {
				val, _ := strconv.ParseFloat(token.(string), 64)
				values.Push(CreateVariable(val)) // Push the number to the values stack
			} else {
				val, _ := strconv.ParseInt(token.(string), 10, 64)
				values.Push(CreateVariable(val)) // Push the number to the values stack
			}

			// Check if the token is a boolean
		} else if IsBoolean(token.(string)) {

			if token.(string) == "true" {
				values.Push(CreateVariable(true))
			} else {
				values.Push(CreateVariable(false))
			}

			// Check if the token is a variable
		} else if (*scope).VariableExists(token.(string)) {

			values.Push((*scope).GetVariable(token.(string)))

			// Check if the token is a left parenthesis
		} else if token.(string) == "(" {
			operators.Push(token.(string)) // Push the token to the operators stack

			// Check if the token is a right parenthesis
		} else if token.(string) == ")" {
			peeked, valid := operators.Peek()

			// Check if the operators stack is empty
			if !valid {
				return Variable{}, errors.New("Error: Invalid expression. Missing a \"(\"")
			}

			for peeked.(string) != "(" {

				// Pop the operator from the stack
				operator, valid := operators.Pop()

				if !valid {
					return Variable{}, errors.New("Error: Invalid expression. Missing a \"(\"")
				}

				// Check for negation
				if operator.(string) == "¬" {
					val2, exists2 := values.Pop()

					if !exists2 {
						return Variable{}, errors.New("Error: Invalid expression 1")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), Variable{}, val2.(Variable))

					// Handle any operation errors
					if opError != nil {
						return Variable{}, opError
					}

					values.Push(result) // Push the result to the values stack
				} else {

					// Normal operation
					val2, exists2 := values.Pop()
					val1, exists1 := values.Pop()
					if !exists1 || !exists2 {
						return Variable{}, errors.New("Error: Invalid expression 2")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), val1.(Variable), val2.(Variable))

					// Handle any operation errors
					if opError != nil {
						return Variable{}, opError
					}

					values.Push(result) // Push the result to the values stack
				}

				peeked, valid = operators.Peek()

				if !valid {
					return Variable{}, errors.New("Error: Invalid expression 3")
				}
			}

			// Pop the left parenthesis from the stack
			operators.Pop()

		} else if isOperator(token.(string)) {

			currOp := token.(string)
			peeked, _ := operators.Peek()

			for !operators.IsEmpty() && OperatorPrecedence(peeked.(string)) >= OperatorPrecedence(currOp) {

				// Pop the operator from the stack
				operator, _ := operators.Pop()

				// Check for negation
				if operator.(string) == "¬" {
					val2, exists2 := values.Pop()
					if !exists2 {
						return Variable{}, errors.New("Error: Invalid expression 4")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), Variable{}, val2.(Variable))

					// Handle any operation errors
					if opError != nil {
						return Variable{}, opError
					}

					values.Push(result) // Push the result to the values stack
				} else {

					// Normal operation
					val2, exists2 := values.Pop()
					val1, exists1 := values.Pop()
					if !exists1 || !exists2 {
						return Variable{}, errors.New("Error: Invalid expression 5")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), val1.(Variable), val2.(Variable))

					// Handle any operation errors
					if opError != nil {
						return Variable{}, opError
					}

					values.Push(result) // Push the result to the values stack
				}

				peeked, _ = operators.Peek()

			}

			operators.Push(currOp) // Push the operator to the operators stack
		}

	}

	for !operators.IsEmpty() {
		operator, _ := operators.Pop()

		// Check for negation
		if operator.(string) == "¬" {
			val2, exists2 := values.Pop()
			if !exists2 {
				return Variable{}, errors.New("Error: Invalid expression 7")
			}

			// Compute the result
			result, opError := ApplyOperator(operator.(string), Variable{}, val2.(Variable))

			// Handle any operation errors
			if opError != nil {
				return Variable{}, opError
			}

			values.Push(result) // Push the result to the values stack
		} else {

			// Normal operation
			val2, exists2 := values.Pop()
			val1, exists1 := values.Pop()
			if !exists1 || !exists2 {
				return Variable{}, errors.New("Error: Invalid expression 8")
			}

			// Compute the result
			result, opError := ApplyOperator(operator.(string), val1.(Variable), val2.(Variable))

			// Handle any operation errors
			if opError != nil {
				return Variable{}, opError
			}

			values.Push(result) // Push the result to the values stack
		}

	}

	value, exists := values.Pop()

	if !exists {
		return Variable{}, errors.New("Error: Empty expression")
	}

	return value.(Variable), nil

}
