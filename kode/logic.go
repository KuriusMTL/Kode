package kode

import "errors"

/**
 * Evaluate an expression.
 * @param str : string - The expression to evaluate.
 * @return string - The result of the expression.
 * @return error - The error if any.
 */
func EvaluateExpression(str string) (string, error) {

	// Tokenize the line
	tokens := inlineParse(str)
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
			values.Push(token.(string)) // Push the number to the values stack

			// Check if the token is a left parenthesis
		} else if token.(string) == "(" {
			operators.Push(token.(string)) // Push the token to the operators stack

			// Check if the token is a right parenthesis
		} else if token.(string) == ")" {
			peeked, valid := operators.Peek()

			// Check if the operators stack is empty
			if !valid {
				return "", errors.New("Error: Invalid expression. Missing a \"(\".")
			}

			for peeked.(string) != "(" {

				// Pop the operator from the stack
				operator, valid := operators.Pop()

				if !valid {
					return "", errors.New("Error: Invalid expression. Missing a \"(\".")
				}

				// Check for negation
				if operator.(string) == "¬" {
					val2, exists2 := values.Pop()
					if !exists2 {
						return "", errors.New("Error: Invalid expression.")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), "0", val2.(string))

					// Handle any operation errors
					if opError != nil {
						return "", opError
					}

					values.Push(result) // Push the result to the values stack
				} else {

					// Normal operation
					val2, exists2 := values.Pop()
					val1, exists1 := values.Pop()
					if !exists1 || !exists2 {
						return "", errors.New("Error: Invalid expression.")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), val1.(string), val2.(string))

					// Handle any operation errors
					if opError != nil {
						return "", opError
					}

					values.Push(result) // Push the result to the values stack
				}

				peeked, valid = operators.Peek()

				if !valid {
					return "", errors.New("Error: Invalid expression.")
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
						return "", errors.New("Error: Invalid expression.")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), "0", val2.(string))

					// Handle any operation errors
					if opError != nil {
						return "", opError
					}

					values.Push(result) // Push the result to the values stack
				} else {

					// Normal operation
					val2, exists2 := values.Pop()
					val1, exists1 := values.Pop()
					if !exists1 || !exists2 {
						return "", errors.New("Error: Invalid expression.")
					}

					// Compute the result
					result, opError := ApplyOperator(operator.(string), val1.(string), val2.(string))

					// Handle any operation errors
					if opError != nil {
						return "", opError
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
				return "", errors.New("Error: Invalid expression.")
			}

			// Compute the result
			result, opError := ApplyOperator(operator.(string), "0", val2.(string))

			// Handle any operation errors
			if opError != nil {
				return "", opError
			}

			values.Push(result) // Push the result to the values stack
		} else {

			// Normal operation
			val2, exists2 := values.Pop()
			val1, exists1 := values.Pop()
			if !exists1 || !exists2 {
				return "", errors.New("Error: Invalid expression.")
			}

			// Compute the result
			result, opError := ApplyOperator(operator.(string), val1.(string), val2.(string))

			// Handle any operation errors
			if opError != nil {
				return "", opError
			}

			values.Push(result) // Push the result to the values stack
		}

	}

	value, exists := values.Pop()

	if !exists {
		return "", errors.New("Error: Empty expression.")
	}

	return value.(string), nil

}
