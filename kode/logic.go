package kode

import (
	"errors"
	"strconv"
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
	// println("Evaluating expression: " + str)
	tokens := InlineParse(str, []string{" ", "\t", "\r", "\n", ",", ".", "!=", "*", "/", "+", "-", "(", ")", "[", "]", "¬", "^", "%", "\"", "\\\"", "<=", ">=", ">", "<", "=="}, true)

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
		// println("Token: " + token.(string))

		// ! SELF
		if token.(string) == "self" {
			(*scope).IsInstance = true
			values.Push(CreateVariable(*scope))
			// ! PARENT
		} else if token.(string) == "super" {
			(*scope).IsInstance = true
			values.Push(CreateVariable(*(*scope).Parent))
			// Check if the token is a number
			// ! NULL
		} else if token.(string) == "null" {
			values.Push(CreateVariable(nil))

			// ! SUB VARIABLE
		} else if token.(string) == "." {
			// Get the latest value from the stack
			value, hasValue := values.Pop()

			if !hasValue {
				return CreateVariable(nil), errors.New("Improper use of '.'")
			}

			if value.(Variable).Type != "func" {
				return CreateVariable(nil), errors.New("Could not access a non-function (" + value.(Variable).Type + ") using '.'")
			}

			// Get the function
			function := value.(Variable).Value.(Function)

			// Get the next token being the function name
			varName, hasVar := queue.Pop()

			if !hasVar {
				return CreateVariable(nil), errors.New("Improper use of '.'")
			}

			// Get the function
			variable := function.GetVariable(varName.(string))

			// Check if the variable exists in the function
			if variable == nil {
				return CreateVariable(nil), errors.New("Could not find variable '" + varName.(string) + "' in the function")
			}

			// Evaluate the sub variable
			if (*variable).Type == "func" {

				// Check if the function is called
				// Check if the next token is a parenthesis
				nextToken, hasNextToken := queue.Peek()
				if !hasNextToken || nextToken.(string) != "(" {
					// Add the variable to the values stack
					values.Push(*(*scope).GetVariable(token.(string)))
				} else {

					// Extract the function's arguments
					args, err := (*scope).ExtractFunctionArgs(&queue)

					if err != nil {
						return Variable{}, err
					}

					// Call the function
					function := (*variable).Value.(Function)
					copyFunc := CopyFunction(&function)
					newVars := CopyVariableMap((*copyFunc).Parent.Variables)
					(*copyFunc).Variables = newVars
					instance, _, err := (*copyFunc).Run(args, map[string]*Variable{})
					if err != nil {
						return Variable{}, err
					}

					values.Push(*instance)
				}

			} else {
				// Else push the variable to the values stack
				values.Push(*variable)
			}

			// ! NUMBER
		} else if IsNumber(token.(string)) {

			// Extract full number from the following tokens
			// Since the token "." is split, it is necessary to check if the next token is a number
			strNumber := token.(string)
			isFloat := false

			nextToken, hasNextToken := queue.Peek()

			for hasNextToken && (nextToken.(string) == "." || IsNumber(nextToken.(string))) {

				queue.Pop()

				// Check if the number is a float
				if nextToken.(string) == "." {

					// If it meets another ".", it is not a float and has an invalid format (error)
					if isFloat {
						return CreateVariable(nil), errors.New("Invalid floating point number")
					} else {
						isFloat = true
					}
				}
				strNumber += nextToken.(string)
				nextToken, hasNextToken = queue.Peek()

				// if !hasNextToken {
				// 	break
				// } else {
				// 	token = nextToken
				// }
			}

			// Is float?
			// ? TODO (Eduard): Select multiple decimal points automatically
			// println("num: " + strNumber)
			if isFloat {
				val, _ := strconv.ParseFloat(strNumber, 64)
				// println(val)
				values.Push(CreateVariable(val)) // Push the number to the values stack
			} else {
				val, _ := strconv.ParseInt(strNumber, 10, 64)
				// println(val)
				values.Push(CreateVariable(val)) // Push the number to the values stack
			}

			// Check if the token is a boolean
			// ! BOOLEAN
		} else if IsBoolean(token.(string)) {

			// Is boolean true?
			// Through IsBoolean(), the token is already lowercase and either "true" or "false"
			if token.(string) == "true" {
				values.Push(CreateVariable(true))
			} else {
				values.Push(CreateVariable(false))
			}

			// Check if the token is the beginning of a string with quotes
			// ! STRING
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
			// ! NEW VARIABLE
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

			// Check if the function is called
			// Check if the next token is a parenthesis
			nextToken2, hasNextToken := queue.Peek()
			if !hasNextToken || nextToken2.(string) != "(" {
				// Add the variable to the values stack
				values.Push(*(*scope).GetVariable(nextToken.(string)))
			} else {

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
			}

			// ! VARIABLE
		} else if (*scope).VariableExists(token.(string)) {

			// If its a function, evaluate it and push the result to the values stack
			if (*scope).GetVariable(token.(string)).Type == "func" {

				// Check if the function is called
				// Check if the next token is a parenthesis
				nextToken, hasNextToken := queue.Peek()
				if !hasNextToken || nextToken.(string) != "(" {
					// Add the variable to the values stack
					values.Push(*(*scope).GetVariable(token.(string)))
				} else {
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
				}

				// Variable is an array
				// Or a string
			} else if isArrayType((*scope).GetVariable(token.(string)).Type) || (*scope).GetVariable(token.(string)).Type == "string" {
				variable := *(*scope).GetVariable(token.(string))

				// Check if the next token is a square bracket
				peeked, validPeek := queue.Peek()
				for validPeek && peeked.(string) == "[" {

					// Check if array
					if !isArrayType(variable.Type) && variable.Type != "string" {
						return Variable{}, errors.New("Error: '" + token.(string) + "' cannot access that index because it might not be an array or a string")
					}

					queue.Pop()
					indexArray, err := (*scope).ExtractArrayValues(&queue)
					if err != nil {
						return Variable{}, err
					}
					if len(indexArray) != 1 || indexArray[0].Type != "int" {
						return Variable{}, errors.New("Error: Invalid array index")
					}

					// Extract the max index
					size, err := GetArraySize(&variable)
					if err != nil {
						return Variable{}, err
					}

					index := indexArray[0].Value.(int64) % size
					if index < 0 { // Handle negative indexes
						index += size
					}

					if variable.Type == "string" {
						variable = CreateVariable(string(variable.Value.(string)[index]))
						peeked, validPeek = queue.Peek()
					} else {
						variable = variable.Value.([]Variable)[index]
						peeked, validPeek = queue.Peek()
					}

				}

				values.Push(variable)

			} else {
				// Else push the variable to the values stack
				values.Push(*(*scope).GetVariable(token.(string)))
			}

			// ! ARRAY
		} else if token.(string) == "[" {

			// Extract the array's value
			array, err := (*scope).ExtractArrayValues(&queue)
			if err != nil {
				return Variable{}, err
			}

			// Create the array variable
			arrayVar := CreateVariable(array)
			values.Push(arrayVar)

			// Check if the token is a left parenthesis
			// ! LEFT PARENTHESIS
		} else if token.(string) == "(" {
			operators.Push(token.(string)) // Push the token to the operators stack

			// Check if the token is a right parenthesis
			// ! RIGHT PARENTHESIS
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

			// ! OPERATOR
		} else if isOperator(token.(string)) {

			// println("Operator: " + token.(string))

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

			// ! PREBUILT FUNCTION
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

			// ! UNKNOWN
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
