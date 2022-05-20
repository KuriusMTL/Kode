package kode

import (
	"errors"
	"strconv"
	"strings"
)

// ! Argument : A function argument.
// -------------------------
// !Name : string - The name of the argument.
// -------------------------
// !Variable : *Variable - The variable of the argument.
type Argument struct {
	Name     string
	Variable *Variable
}

// ! Function : A function with scope.
// -------------------------
// ! Arguments : The arguments template of the function as a map of variables.
// The "value" of the variables are the default values.
// The "type" of the variables are the types of the variables.
// -------------------------
// ! Variables : The local variables of the function as a map of variables.
// -------------------------
// ! Return : Expected return of the function.
// -------------------------
// ! Code : The code of the function.
// ? TODO (Eduard): Precompile the functions instead of parsing them every time.
type Function struct {
	Arguments  []Argument
	Variables  map[string](*Variable)
	Return     string
	Code       string
	Parent     *Function
	IsInstance bool
}

/**
 * Create a new function with scope.
 * @param arguments : map[string]Variable - The arguments of the function.
 * @param variables : map[string](*Variable) - The variables of the function.
 * @param code : string - The code of the function.
 * @return Function - The new function.
 */
func CreateFunction(argumentsTemplate []Argument, variables map[string](*Variable), returnType string, parent *Function, code string) Function {

	// Make sur to initialize the variables
	if variables == nil {
		variables = map[string](*Variable){}
	}

	// Create a copy of the variables map
	vars := CopyVariableMap(variables)

	// Update the variables with the arguments
	for _, arg := range argumentsTemplate {

		vars[arg.Name] = &Variable{
			Value: GetDefaultValue((*arg.Variable).Type),
			Type:  (*arg.Variable).Type,
		}
	}

	function := Function{
		Arguments:  argumentsTemplate,
		Variables:  vars,
		Return:     returnType,
		Code:       code,
		Parent:     parent,
		IsInstance: false,
	}

	if parent == nil {
		function.Parent = &function
	}

	return function
}

/**
 * Add argument variables to the current scope.
 * @param args : []*Variable - The arguments to add.
 * @param scope : *Function - The current scope.
 * @return error - The error, if any.
 */
func (scope *Function) ArgumentsToVariables(args []*Variable) error {

	// Set the variables
	// Loop through each argument and set the appropriate variable
	for i, arg := range args {

		// Check if there are more arguments than variables for the function
		if i >= len((*scope).Arguments) {
			return errors.New("Too many arguments")
		}

		// Determine if the type of the variable is compatible with the function argument
		// ! Exception: If the function argument is "val" then it is compatible with any type.
		if (*scope).Arguments[i].Variable.Type == (*arg).Type || (*scope).Arguments[i].Variable.Type == "val" {

			if (*arg).Type == "string" || (*arg).Type == "int" || (*arg).Type == "float" || (*arg).Type == "bool" {

				// Create a copy of the variable
				// And set the variable value inside the function scope
				varCopy := *arg
				(*scope).Variables[(*scope).Arguments[i].Name] = &varCopy

			} else {

				// Other types only need to pass their reference
				(*scope).Variables[(*scope).Arguments[i].Name] = arg
			}

		} else {
			return errors.New("Argument type mismatch for the argument \"" + (*scope).Arguments[i].Name + "\"")
		}
	}

	return nil // No error
}

func (scope *Function) ExtractFunctionArgs(queue *Queue) ([]*Variable, error) {

	// Evaluate the function
	paranthesis, _ := (*queue).Pop()
	if paranthesis == nil || paranthesis.(string) != "(" {
		return nil, errors.New("Error: Missing paranthesis for the function")
	}

	parameters := []*Variable{}
	closedFunction := false
	tempVariable := ""
	isInString := false
	nestedParenthesis := 0
	// Get the function parameters
	for !queue.IsEmpty() {

		// Get the next token
		nextToken, _ := queue.Pop()

		//println(nextToken.(string))

		// Check if the token is the end of the function parameters
		if (nextToken.(string) == "(" || nextToken.(string) == "[") && !isInString {
			nestedParenthesis += 1
		} else if (nextToken.(string) == ")" || nextToken.(string) == "]") && !isInString {
			if nestedParenthesis == 0 && nextToken.(string) == ")" {
				closedFunction = true
				break
			}
			nestedParenthesis -= 1
		} else if nextToken.(string) == "," && nestedParenthesis == 0 && !isInString {

			//println(tempVariable)

			// Evaluate the parameter
			parameter, err := EvaluateExpression(scope, tempVariable)
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, &parameter)

			tempVariable = ""
			continue
		}

		// Check if the token is a string
		if nextToken.(string) == "\"" {
			isInString = !isInString
		}

		// If the token is not a string, and the string is not closed, add the token to the temp variable
		if isInString {
			tempVariable += nextToken.(string)
		} else {
			tempVariable += nextToken.(string) + " "
		}
	}

	if !closedFunction {
		return nil, errors.New("Error: Missing closing paranthesis for the function call")
	}

	// Evaluate the last parameter
	if tempVariable != "" {
		//println(tempVariable)
		parameter, err := EvaluateExpression(scope, tempVariable)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, &parameter)
	}

	return parameters, nil
}

/**
 * Run a function. The function will be executed in the scope of the function.
 * @param scope : Function - The scope of the referenced function.
 * @param args : interface{} - The arguments of the function.
 * @return *Variable - Return of the function.
 */
func (scope *Function) Run(args []*Variable, vars map[string](*Variable)) (*Variable, bool, error) {

	// Set the variables
	// Loop through each argument and set the appropriate variable
	for key, value := range vars {
		(*scope).Variables[key] = value
	}

	// Set the argument variables
	err := scope.ArgumentsToVariables(args)
	if err != nil {
		return nil, false, err
	}

	// Split the code into lines.
	lines := LineParse((*scope).Code)

	// Loop through the lines.
	for currentLine := 0; currentLine < len(lines); currentLine++ {
		line := lines[currentLine]

		// Skip empty lines.
		if line == "" || line == "\r" {
			continue
		}

		// Get the command tokens of the line.
		// Add the tokens to the queue.
		tokens := Queue{}
		tokensArr := InlineParse(line, []string{" ", "\t", "\r", "[", "]", "!=", "==", ">=", "<=", "=", ":=", "#", "\"", "\\\"", ",", ".", "(", ")", "\""}, true)
		for _, token := range tokensArr {
			tokens.Push(token)
		}

		// Loop through the tokens.
		for !tokens.IsEmpty() {
			// First token is the command.
			// The command determines the action.
			command, _ := tokens.Pop()

			switch command {

			// ? Comment
			// The comment is ignored by the interpreter.
			// ? TODO (Eduard): Add support for multiline comments.
			case "#":

				// Ignore the rest of the line.
				tokens.Clear()
				continue

			// ? Variable creation
			// The variable is created in the current scope of the function.
			// "val" <name> = <value> where the type is inferred from the value.
			case "val", "int", "float", "string", "bool":

				// Get the dimensions of the variable.
				// The dimensions are optional.
				dimension, err := ExtractArrayDimensionFromDeclaration(&tokens)
				if err != nil {
					return nil, false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1))
				}

				command = command.(string) + strings.Repeat("[]", dimension)

				// Get the provided variable name
				name, nameProvided := tokens.Pop()

				// Check if the name for the variable was provided
				if !nameProvided {
					return NullVariable(), false, errors.New("Error: Missing variable name on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Check if the variable name is valid
				if !HasValidVariableName(name.(string)) {
					return nil, false, errors.New("Error: Invalid variable name on line " + strconv.Itoa(currentLine+1) + ". The name must be alphanumeric and start with a letter.")
				}

				// Check if the variable name is already in use in the current scope
				if (*scope).VariableExists(name.(string)) {
					return NullVariable(), false, errors.New("Error: Variable was already declared on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Get the expected variable declaration format
				// e.g. "val <name> = <value>"
				equal, assign := tokens.Pop()

				// Check if the variable has an assignment
				if !assign || equal.(string) != "=" {
					return NullVariable(), false, errors.New("Error: Missing variable assignment \"=\" on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Get the rest of the line tokens and join them to feed the variable value
				// Convert queue tokens back into a string
				// Also, account that quotes define a string
				value := InlineQueueToString(&tokens)

				// Make sure the variable value is not empty
				if value == "" {
					return NullVariable(), false, errors.New("Error: Missing variable value on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Create the variable and evaluate the value
				evaluatedValue, err := EvaluateExpression(scope, value)
				if err != nil {
					return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
				}

				if command.(string) != "val" && evaluatedValue.Type != command.(string) {
					return NullVariable(), false, errors.New("Error: Invalid variable type on line " + strconv.Itoa(currentLine+1) + ". The type of the variable must be " + command.(string) + ".")
				}

				// Create the variable in the current scope.
				(*scope).Variables[name.(string)] = &evaluatedValue
				if (*scope).GetVariable("_DEBUG").Type == "bool" && (*scope).GetVariable("_DEBUG").Value.(bool) {
					println("Created variable " + name.(string) + "(" + evaluatedValue.Type + ").")
				}
				break

			// ? If condition
			// The condition is evaluated and if it is true, the code is executed.
			case "if":

				// Parse the condition blocks
				conditionBlocks, nextLine, err := ParseConditionBlocks(&tokens, currentLine, lines)
				if err != nil {
					return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
				}

				currentLine = nextLine // Update the current line to skip the code block

				// Visit the condition blocks
				// goToNextLine ensures that the loop is properly exited
				goToNextLine := false
				for _, conditionBlock := range conditionBlocks {

					// Evaluate the condition

					// Is an else statement
					if conditionBlock.Condition == "else" {

						// Directly execute the code
						ifCondition := CreateFunction([]Argument{}, (*scope).Variables, "val", scope, conditionBlock.Code)
						returnValue, toReturn, err := ifCondition.Run([]*Variable{}, map[string]*Variable{})

						// If the code returns a value, return it
						if err != nil {
							return NullVariable(), false, err
						}

						if toReturn {
							return returnValue, true, nil
						}

						goToNextLine = true
						break

					} else {

						// Is an if statement or if else statement
						// Evaluate the condition
						evaluatedCondition, err := EvaluateExpression(scope, conditionBlock.Condition)
						if err != nil {
							return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(conditionBlock.ConditionIndex+1) + ".")
						}

						// Check the type of the evaluated condition
						// If it is not a boolean, return an error
						if evaluatedCondition.Type != "bool" {
							return NullVariable(), false, errors.New("Error: Invalid condition on line " + strconv.Itoa(conditionBlock.ConditionIndex+1) + ". The condition must be a boolean.")
						}

						// If the condition is true, execute the block of code
						if evaluatedCondition.Value.(bool) {

							// Create a new scope for the block

							ifCondition := CreateFunction([]Argument{}, (*scope).Variables, "val", scope, conditionBlock.Code)
							returnValue, toReturn, err := ifCondition.Run([]*Variable{}, map[string]*Variable{})
							if err != nil {
								return NullVariable(), false, err
							}

							if toReturn {
								return returnValue, true, nil
							}

							goToNextLine = true
							break
						}

						// Else, visit the next condition

					}

				}

				// This statement makes sure to skip to the next line immediately
				if goToNextLine {
					break
				}

			// ? Function creation
			case "func":

				// ! Functions act like variables
				// Get the provided function name
				name, nameProvided := tokens.Pop()

				// Check if the name for the function was provided
				if !nameProvided {
					return NullVariable(), false, errors.New("Error: Missing function name at decleration on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Check if the function name is valid
				// Again, they act like variables
				if !HasValidVariableName(name.(string)) {
					return NullVariable(), false, errors.New("Error: Invalid function name on line " + strconv.Itoa(currentLine+1) + ". The name must be alphanumeric and start with a letter.")
				}

				// Check if the function name is already in use in the current scope
				// A function and primitive variable cannot have the same name
				if (*scope).VariableExists(name.(string)) {
					return NullVariable(), false, errors.New("Error: Function was already declared by another variable on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Get the parameters for the function
				// Check if the function parameters start with a paranthesis
				char, charProvided := tokens.Pop()
				if !charProvided || char.(string) != "(" {
					return NullVariable(), false, errors.New("Error: Missing opening parenthesis for the parameters on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Parameters list
				parameters := []Argument{}

				for !tokens.IsEmpty() {

					// Check if the next token is a closing parenthesis
					// If it is, break the loop
					token, tokenProvided := tokens.Pop()

					if !tokenProvided {
						return NullVariable(), false, errors.New("Error: Missing closing parenthesis for the parameters on line " + strconv.Itoa(currentLine+1) + ".")
					}

					if token.(string) == ")" {
						break
					} else {
						// Check if the parameter type is valid
						// If it is not, return an error
						if token.(string) != "val" && token.(string) != "int" && token.(string) != "float" && token.(string) != "bool" && token.(string) != "string" {
							return NullVariable(), false, errors.New("Error: Invalid parameter type \"" + token.(string) + "\" on line " + strconv.Itoa(currentLine+1) + ".")
						}

						// Get the dimensions of the variable.
						// The dimensions are optional.
						dimension, err := ExtractArrayDimensionFromDeclaration(&tokens)
						if err != nil {
							return nil, false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1))
						}

						token = token.(string) + strings.Repeat("[]", dimension)

						// Get the parameter name
						parameterName, parameterNameProvided := tokens.Pop()
						if !parameterNameProvided {
							return NullVariable(), false, errors.New("Error: Missing parameter name on line " + strconv.Itoa(currentLine+1) + ".")
						}

						// Check if the parameter name is valid
						// If it is not, return an error
						if !HasValidVariableName(parameterName.(string)) {
							return NullVariable(), false, errors.New("Error: Invalid parameter name \"" + parameterName.(string) + "\" on line " + strconv.Itoa(currentLine+1) + ". The name must be alphanumeric and start with a letter.")
						}

						// Create the parameter
						parameters = append(parameters, Argument{
							Name: parameterName.(string),
							Variable: &Variable{
								Type:  token.(string),
								Value: nil,
							},
						})

						// Check if the next token is a comma or a closing parenthesis
						// If it is, continue the loop
						token, tokenProvided = tokens.Pop()
						if !tokenProvided {
							return NullVariable(), false, errors.New("Error: Missing closing parenthesis for the parameters on line " + strconv.Itoa(currentLine+1) + ".")
						}

						if token.(string) == "," {
							continue
						} else if token.(string) == ")" {
							break
						} else {
							return NullVariable(), false, errors.New("Error: Invalid function syntax \"" + token.(string) + "\" on line " + strconv.Itoa(currentLine+1) + ".")
						}

					}

				}

				// Get return type
				returnType, returnTypeProvided := tokens.Pop()
				if !returnTypeProvided {
					returnType = "null"
				} else if returnType.(string) != "val" && returnType.(string) != "int" && returnType.(string) != "float" && returnType.(string) != "bool" && returnType.(string) != "string" && returnType.(string) != "func" && isArrayType(returnType.(string)) == false {
					return NullVariable(), false, errors.New("Error: Invalid return type \"" + returnType.(string) + "\" on line " + strconv.Itoa(currentLine+1) + ".")
				}

				// Get the block of code for the function
				funcEnded := false
				functionStart := currentLine
				functionCode := ""
				currentLine++
				for currentLine < len(lines) {
					parsed := InlineParse(lines[currentLine], []string{" "}, true)

					if len(parsed) > 1 && parsed[0] == "end" && parsed[1] == name.(string) {
						funcEnded = true
						break
					} else {
						functionCode += lines[currentLine] + "\n"
					}

					currentLine++
				}

				if !funcEnded {
					return NullVariable(), false, errors.New("Error: Missing end for function \"" + name.(string) + "\" on line " + strconv.Itoa(functionStart+1) + ".")
				}

				// Create the function and add it to the scope
				//function := CreateVariable(CreateFunction(parameters, (*scope).Variables, returnType.(string), functionCode))
				function := CreateVariable(CreateFunction(parameters, make(map[string]*Variable), returnType.(string), scope, functionCode))
				(*scope).Variables[name.(string)] = &function

			case "return":
				// Retrive the return value
				// Get the return value expression
				expression := InlineQueueToString(&tokens)

				// For loop inside expression array
				// println(expression)

				returnValue := Variable{}
				// Evaluate the expression
				if expression != "" {
					// Evaluate the expression
					value, err := EvaluateExpression(scope, expression)
					if err != nil {
						return NullVariable(), false, err
					}
					// Set the return value
					returnValue = value
				}

				// Check if the return value is valid type
				// If it is not, return an error
				if returnValue == (Variable{}) {
					return NullVariable(), true, nil
				} else if returnValue.Type == (*scope).Return || (*scope).Return == "val" {
					return &returnValue, true, nil
				} else {
					return NullVariable(), true, errors.New("Error: Invalid return type \"" + returnValue.Type + "\" on line " + strconv.Itoa(currentLine+1) + ".")
				}

			case "break":
				return NullVariable(), false, nil

			case "for":

				loopBlock, nextLine, err := ParseLoopBlock(&tokens, currentLine, lines)

				if err != nil {
					return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
				}

				currentLine = nextLine

				evaluatedCondition, err := EvaluateExpression(scope, loopBlock.Condition)
				if err != nil {
					return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
				}

				if evaluatedCondition.Type != "bool" {
					return NullVariable(), false, errors.New("Error: Invalid condition type \"" + evaluatedCondition.Type + "\" on line " + strconv.Itoa(currentLine+1) + ".")
				}

				for evaluatedCondition.Value.(bool) {

					forLoop := CreateFunction([]Argument{}, (*scope).Variables, "val", scope, loopBlock.Code)
					returnValue, toReturn, err := forLoop.Run([]*Variable{}, map[string]*Variable{})

					// If the code returns a value, return it
					if err != nil {
						return NullVariable(), false, err
					}

					if toReturn {
						return returnValue, true, nil
					}

					// Evaluate the condition
					evaluatedCondition, err = EvaluateExpression(scope, loopBlock.Condition)
					if err != nil {
						return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
					}

					if evaluatedCondition.Type != "bool" {
						return NullVariable(), false, errors.New("Error: Invalid condition type \"" + evaluatedCondition.Type + "\" on line " + strconv.Itoa(currentLine+1) + ".")
					}

					// Exit the loop if the condition is false
					if !evaluatedCondition.Value.(bool) {
						break
					}

				}

			default:

				// Is the command a variable (check inside the scope)?
				// If so, update the variable value.
				if (*scope).VariableExists(command.(string)) {

					variable := (*scope).GetVariable(command.(string))

					// Check for array index
					peeked, validPeek := tokens.Peek()
					for validPeek && peeked.(string) == "[" {

						// Check if array
						if !isArrayType((*variable).Type) {
							return NullVariable(), false, errors.New("Error: '" + peeked.(string) + "' cannot access that index because it might not be an array on line " + strconv.Itoa(currentLine+1) + ".")
						}

						tokens.Pop()
						indexArray, err := (*scope).ExtractArrayValues(&tokens)
						if err != nil {
							return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
						}
						if len(indexArray) != 1 || indexArray[0].Type != "int" {
							return NullVariable(), false, errors.New("Error: Invalid array index on line " + strconv.Itoa(currentLine+1) + ".")
						}
						size := int64(len((*variable).Value.([]Variable)))
						if size == 0 {
							return NullVariable(), false, errors.New("Error: Array is empty on line " + strconv.Itoa(currentLine+1) + ".")
						}
						index := indexArray[0].Value.(int64) % size
						if index < 0 { // Handle negative indexes
							index += size
						}
						variable = &(*variable).Value.([]Variable)[index]
						peeked, validPeek = tokens.Peek()

					}

					// Get the provided definition token
					// i.e. "=" or ":="
					equal, assign := tokens.Peek()

					// Check if the variable has a valid assignment
					if !assign || (equal.(string) != "=" && equal.(string) != ":=") {

						// Simply execute the command and return NO value

						output := command.(string) + InlineQueueToString(&tokens)
						_, err := EvaluateExpression(scope, output)

						if err != nil {
							return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
						}

					} else {

						// Get the rest of the line tokens and join them to feed the variable value.

						// Remove the equal sign
						tokens.Pop()

						value := InlineQueueToString(&tokens)

						// Make sure the variable value is valid (not empty)
						if value == "" {
							return NullVariable(), false, errors.New("Error: Missing variable value on line " + strconv.Itoa(currentLine+1) + ".")
						}

						// Create the new variable.
						// Evaluate the value and update the variable.
						evaluatedValue, err := EvaluateExpression(scope, value)
						if err != nil {
							return NullVariable(), false, errors.New(err.Error() + " on line " + strconv.Itoa(currentLine+1) + ".")
						}

						// Check safe assignment
						if ((*variable).Type != evaluatedValue.Type) && equal.(string) != ":=" {
							return NullVariable(), false, errors.New("Error: Variable type mismatch on line " + strconv.Itoa(currentLine+1) + ". Expected type " + (*scope).Variables[command.(string)].Type + " but got type " + evaluatedValue.Type + ".")
						}

						// Update the variable in the current scope.
						*variable = evaluatedValue
						// *((*scope).Variables[command.(string)]) = evaluatedValue
						if (*scope).GetVariable("_DEBUG").Type == "bool" && (*scope).GetVariable("_DEBUG").Value.(bool) {
							println("Updated variable " + command.(string) + "(" + evaluatedValue.Type + ").")
						}

					}

				} else if ExistsIncluded(command.(string)) {

					// Extract the function's arguments
					args, err := (*scope).ExtractFunctionArgs(&tokens)

					if err != nil {
						return NullVariable(), false, err
					}

					// Call the function

					_, err = RunIncluded(command.(string), args)
					if err != nil {
						return NullVariable(), false, err
					}

				} else {
					// Command is unknown
					return NullVariable(), false, errors.New("Error: Unknown command \"" + command.(string) + "\" on line " + strconv.Itoa(currentLine+1) + ".")
				}

			}

		}

	}

	return NullVariable(), false, nil
}
