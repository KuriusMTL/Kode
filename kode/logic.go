package kode

import (
	"errors"
	"strconv"
	"strings"
)

/**
 * Evaluate an expression.
 * @param scope : *Function - The scope of the expression.
 * @param str : string - The expression to evaluate.
 * @return Variable - The result of the expression.
 * @return error - The error if any.
 */
func EvaluateExpression(scope *Function, str string) (Variable, error) {

	// Tokenize the line
	// "is" and "not" are implicitly parsed as well
	tokens := InlineParse(str, []string{" ", "\t", "\r", "\n", ",", "!=", "*", "/", "+", "-", "(", ")", "¬", "^", "%", "\"", "\\\"", "<=", ">=", ">", "<", "=="}, true)

	// Replace proper substractions with negation
	// If there is a minus sign, and the next token is a number or a variable,
	tokens = CheckForNegation(tokens)

	// Add the tokens to the queue
	queue := Queue{}
	for _, token := range tokens {
		queue.Push(token)
	}

	values := Stack{}    // Store values to evaluate
	operators := Stack{} // Store operators to evaluate values

	// Loop through the all the tokens
	for !queue.IsEmpty() {

		// Get the next token
		token, _ := queue.Pop()

		if token.(string) == "self" {
			(*scope).IsInstance = true
			values.Push(CreateVariable(*scope))

			// Check if the token is a number
		} else if token.(string) == "null" {
			values.Push(CreateVariable(nil))
		} else if IsNumber(token.(string)) {

			// Is float?
			// ? TODO (Eduard): Select multiple decimal points automatically
			if strings.Contains(token.(string), ".") {
				val, _ := strconv.ParseFloat(token.(string), 64)
				values.Push(CreateVariable(val)) // Push the number to the values stack
			} else {
				val, _ := strconv.ParseInt(token.(string), 10, 64)
				values.Push(CreateVariable(val)) // Push the number to the values stack
			}

			// Check if the token is a boolean
		} else if IsBoolean(token.(string)) {

			// Is boolean true?
			// Through IsBoolean(), the token is already lowercase and either "true" or "false"
			if token.(string) == "true" {
				values.Push(CreateVariable(true))
			} else {
				values.Push(CreateVariable(false))
			}

			// Check if the token is the beginning of a string with quotes
		} else if token.(string) == "\"" {

			token = "" // Remove the quote

			// Loop through the tokens until the next quote to extract the string
			nextToken, hasNextToken := queue.Pop()

			for hasNextToken && nextToken.(string) != "\"" {

				// Replace escape characters
				token = token.(string) + HandleEscapeCharacters(nextToken.(string))
				nextToken, hasNextToken = queue.Pop()
			}

			// Check for incomplete string errors
			if nextToken == nil || nextToken.(string) != "\"" {
				return Variable{}, errors.New("Error: Missing closing quote for string")
			} else {

				// No errors found, push the string to the values stack
				values.Push(CreateVariable(token.(string)))
			}

			// Check if the token is a variable
		} else if token.(string) == "new" {

			// Get the next token
			nextToken, valid := queue.Pop()

			if !valid {
				return Variable{}, errors.New("Error: Missing function name after 'new'")
			}

			if !(*scope).VariableExists(nextToken.(string)) {
				return Variable{}, errors.New("Error: Unknown variable '" + nextToken.(string) + "'")
			}

			// Check if the variable is a function
			if (*scope).GetVariable(nextToken.(string)).Type != "func" {
				return Variable{}, errors.New("Error: '" + nextToken.(string) + "' is not a function")
			}

			// Extract the function's arguments
			args, err := (*scope).ExtractFunctionArgs(&queue)

			if err != nil {
				return Variable{}, err
			}

			// Get the function
			function := (*scope).GetVariable(nextToken.(string)).Value.(Function)
			copyFunc := CopyFunction(&function)
			(*copyFunc).Variables = map[string]*Variable{}
			(*copyFunc).Parent = copyFunc
			instance, _, err := (*copyFunc).Run(args, map[string]*Variable{})
			if err != nil {
				return Variable{}, err
			}

			values.Push(*instance)

		} else if (*scope).VariableExists(token.(string)) {

			// If its a function, evaluate it and push the result to the values stack
			if (*scope).GetVariable(token.(string)).Type == "func" {

				// Extract the function's arguments
				args, err := (*scope).ExtractFunctionArgs(&queue)

				if err != nil {
					return Variable{}, err
				}

				// Call the function
				function := (*scope).GetVariable(token.(string)).Value.(Function)
				copyFunc := CopyFunction(&function)
				newVars := CopyVariableMap((*copyFunc).Parent.Variables)
				(*copyFunc).Variables = newVars
				instance, _, err := (*copyFunc).Run(args, map[string]*Variable{})
				if err != nil {
					return Variable{}, err
				}

				values.Push(*instance)

			} else {
				// Else push the variable to the values stack
				values.Push(*(*scope).GetVariable(token.(string)))
			}

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
				if operator.(string) == "¬" || operator.(string) == "not" {
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
				if operator.(string) == "¬" || operator.(string) == "not" {
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
		} else if ExistsIncluded(token.(string)) {

			// Extract the function's arguments
			args, err := (*scope).ExtractFunctionArgs(&queue)

			if err != nil {
				return Variable{}, err
			}

			// Call the function

			result, err := RunIncluded(token.(string), args)
			if err != nil {
				return Variable{}, err
			}

			values.Push(*result)

		} else {
			return Variable{}, errors.New("Error: Invalid expression \"" + token.(string) + "\"")
		}

	}

	for !operators.IsEmpty() {
		operator, _ := operators.Pop()

		// Check for negation
		if operator.(string) == "¬" || operator.(string) == "not" {
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

	// Check if the values stack is empty. If it is, then the expression is invalid
	if !exists {
		return Variable{}, errors.New("Error: Empty expression")
	}

	return value.(Variable), nil

}
