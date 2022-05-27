package kode

import (
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
// -------------------------
// ! Parent : The reference to the parent function.
// -------------------------
// ! Name : The name of the function.
// ? TODO (Eduard): Precompile the functions instead of parsing them every time.
type Function struct {
	Arguments []Argument
	Variables map[string](*Variable)
	Return    string
	Code      string
	Parent    *Function
	Name      string
	Index     int
}

/**
 * Create a new function with scope.
 * @param arguments : map[string]Variable - The arguments of the function.
 * @param variables : map[string](*Variable) - The variables of the function.
 * @param returnType : string - The return type of the function.
 * @param parent : *Function - The parent function.
 * @param code : string - The code of the function.
 * @return Function - The new function.
 */
func CreateFunction(name string, index int, argumentsTemplate []Argument, variables map[string](*Variable), returnType string, parent *Function, code string) Function {

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
		Arguments: argumentsTemplate,
		Variables: vars,
		Return:    returnType,
		Code:      code,
		Parent:    parent,
		Name:      name,
		Index:     index,
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
 * @return *ErrorStack - The error, if any.
 */
func (scope *Function) ArgumentsToVariables(args []*Variable, startLine int) *ErrorStack {

	// Set the variables
	// Loop through each argument and set the appropriate variable
	for i, arg := range args {

		// Check if there are more arguments than variables for the function
		if i >= len((*scope).Arguments) {
			return CreateError("Error: Argument count (expected at most "+strconv.FormatInt(int64(len((*scope).Arguments)), 10)+") mismatch for \""+(*scope).Name+"\"", startLine)
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
			return CreateError("Error: Argument type mismatch for the argument \""+(*scope).Arguments[i].Name+"\"", startLine)
		}
	}

	return nil // No error
}

func (scope *Function) ExtractFunctionArgs(queue *Queue, depth int64, startLine int) ([]*Variable, *ErrorStack) {

	// Evaluate the function
	parentheses, _ := (*queue).Pop()
	if parentheses == nil || parentheses.(string) != "(" {
		return nil, CreateError("Error: Missing parentheses for the function \""+(*scope).Name+"\"", startLine)
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
			parameter, err := EvaluateExpression(scope, tempVariable, depth, startLine)
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
		return nil, CreateError("Error: Missing closing parentheses for the function call \""+(*scope).Name+"\"", startLine)
	}

	// Evaluate the last parameter
	if tempVariable != "" {
		//println(tempVariable)
		parameter, err := EvaluateExpression(scope, tempVariable, depth, startLine)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, &parameter)
	}

	return parameters, nil
}

/**
 * Run a function. The function will be executed within its scope.
 * @param scope : Function - The scope of the referenced function.
 * @param args : interface{} - The arguments of the function.
 * @param vars : map[string](*Variable) - The variables of the function.
 * @param depth : int - The current depth of the function.
 * @return *Variable - Return of the function.
 * @return int - The return type of the function.
		** 0: Immediate return
		** 1: Return until root function
		** 2: Return until for loop
 * @return *ErrorStack - The error, if any.
*/
func (scope *Function) Run(args []*Variable, vars map[string](*Variable), depth int64, startLine int) (*Variable, int, *ErrorStack) {

	// Limit the depth of the function recursion
	if (*scope).VariableExists("_MAX_RECURSION") && EvaluateType((*(*scope).GetVariable("_MAX_RECURSION")).Value) == "int" {
		if depth > (*(*scope).GetVariable("_MAX_RECURSION")).Value.(int64) {
			return nil, 0, CreateError("Error: Recursion limit reached (_MAX_RECURSION)", startLine)
		}
	} else {
		// Could not find the variable _MAX_RECURSION, default max depth to 5000
		if depth > 5000 {
			return nil, 0, CreateError("Error: Recursion limit reached (5000)", startLine+(*scope).Index)
		}
	}

	// Set the variables
	// Loop through each argument and set the appropriate variable
	for key, value := range vars {
		(*scope).Variables[key] = value
	}

	// Set the argument variables
	err := scope.ArgumentsToVariables(args, startLine)
	if err != nil {
		return nil, 0, err
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
				dimension, err := ExtractArrayDimensionFromDeclaration(&tokens, currentLine+1)
				if err != nil {
					return nil, 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
				}

				command = command.(string) + strings.Repeat("[]", dimension)

				// Get the provided variable name
				name, nameProvided := tokens.Pop()

				// Check if the name for the variable was provided
				if !nameProvided {
					return NullVariable(), 0, CreateError("Error: Missing variable name", startLine+currentLine+(*scope).Index)
				}

				// Check if the variable name is valid
				if !HasValidVariableName(name.(string)) {
					return nil, 0, CreateError("Error: Variable names must be alphanumeric and start with a letter. Invalid variable name \""+name.(string)+"\"", startLine+currentLine+(*scope).Index)
				}

				// Check if the variable name is already in use in the current scope
				if (*scope).VariableExists(name.(string)) {
					return NullVariable(), 0, CreateError("Error: Variable \""+name.(string)+"\" already exists in the current scope", startLine+currentLine+(*scope).Index)
				}

				// Get the expected variable declaration format
				// e.g. "val <name> = <value>"
				equal, assign := tokens.Pop()

				// Check if the variable has an assignment
				if !assign || equal.(string) != "=" {
					return NullVariable(), 0, CreateError("Error: Missing assignment for variable \""+name.(string)+"\"", startLine+currentLine+(*scope).Index)
				}

				// Get the rest of the line tokens and join them to feed the variable value
				// Convert queue tokens back into a string
				// Also, account that quotes define a string
				value := InlineQueueToString(&tokens)

				// Make sure the variable value is not empty
				if value == "" {
					return NullVariable(), 0, CreateError("Error: Missing value for variable \""+name.(string)+"\"", startLine+currentLine+(*scope).Index)
				}

				// Create the variable and evaluate the value
				evaluatedValue, err := EvaluateExpression(scope, value, depth, startLine+currentLine+(*scope).Index)
				if err != nil {
					return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
				}

				if command.(string) != "val" && evaluatedValue.Type != command.(string) {

					// If the evaluated value is an empty array (e.g. []), then its evaluated type would be "val[]" and its length would be 0.
					// Empty arrays are allowed to be assigned to any array type.
					if !isArrayType(command.(string)) || evaluatedValue.Type != "val[]" || len(evaluatedValue.Value.([]Variable)) != 0 {
						return NullVariable(), 0, CreateError("Error: Variable \""+name.(string)+"\" cannot be assigned to type \""+command.(string)+"\"", startLine+currentLine+(*scope).Index)
					} else {
						// Properly assign the variable type for the empty array
						evaluatedValue.Type = command.(string)
					}

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
				conditionBlocks, nextLine, err := ParseConditionBlocks(&tokens, currentLine, lines, startLine+(*scope).Index)
				if err != nil {
					return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
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
						ifCondition := CreateFunction("if", startLine+currentLine, []Argument{}, (*scope).Variables, "val", scope, conditionBlock.Code)
						returnValue, toReturn, err := ifCondition.Run([]*Variable{}, map[string]*Variable{}, depth, startLine+currentLine+(*scope).Index)

						// If the code returns a value, return it
						if err != nil {
							return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
						}

						// Continue returning the value
						if toReturn > 0 {
							return returnValue, toReturn, nil
						}

						goToNextLine = true
						break

					} else {

						// Is an if statement or if else statement
						// Evaluate the condition
						evaluatedCondition, err := EvaluateExpression(scope, conditionBlock.Condition, depth, startLine+currentLine+(*scope).Index)
						if err != nil {
							return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
						}

						// Check the type of the evaluated condition
						// If it is not a boolean, return an error
						if evaluatedCondition.Type != "bool" {
							return NullVariable(), 0, CreateError("Error: Condition must be a boolean", startLine+currentLine+(*scope).Index)
						}

						// If the condition is true, execute the block of code
						if evaluatedCondition.Value.(bool) {

							// Create a new scope for the block

							ifCondition := CreateFunction("if", startLine+currentLine, []Argument{}, (*scope).Variables, "val", scope, conditionBlock.Code)
							returnValue, toReturn, err := ifCondition.Run([]*Variable{}, map[string]*Variable{}, depth, startLine+currentLine+(*scope).Index)

							if err != nil {
								return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
							}

							if toReturn > 0 {
								return returnValue, toReturn, nil
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
					return NullVariable(), 0, CreateError("Error: Function name not provided", startLine+currentLine+(*scope).Index)
				}

				// Check if the function name is valid
				// Again, they act like variables
				if !HasValidVariableName(name.(string)) {
					return NullVariable(), 0, CreateError("Error: The function name must be alphanumeric. Invalid function name \""+name.(string)+"\"", startLine+currentLine+(*scope).Index)
				}

				// Check if the function name is already in use in the current scope
				// A function and primitive variable cannot have the same name
				if (*scope).VariableExists(name.(string)) {
					return NullVariable(), 0, CreateError("Error: The function/variable name \""+name.(string)+"\" is already in use", startLine+currentLine+(*scope).Index)
				}

				// Get the parameters for the function
				// Check if the function parameters start with a parentheses
				char, charProvided := tokens.Pop()
				if !charProvided || char.(string) != "(" {
					return NullVariable(), 0, CreateError("Error: Function parameters must start with a parentheses", startLine+currentLine+(*scope).Index)
				}

				// Parameters list
				parameters := []Argument{}

				for !tokens.IsEmpty() {

					// Check if the next token is a closing parenthesis
					// If it is, break the loop
					token, tokenProvided := tokens.Pop()

					if !tokenProvided {
						return NullVariable(), 0, CreateError("Error: Expected a closing parenthesis", startLine+currentLine+(*scope).Index)
					}

					if token.(string) == ")" {
						break
					} else {
						// Check if the parameter type is valid
						// If it is not, return an error
						if token.(string) != "val" && token.(string) != "int" && token.(string) != "float" && token.(string) != "bool" && token.(string) != "string" {
							return NullVariable(), 0, CreateError("Error: Invalid parameter type \""+token.(string)+"\"", startLine+currentLine+(*scope).Index)
						}

						// Get the dimensions of the variable.
						// The dimensions are optional.
						dimension, err := ExtractArrayDimensionFromDeclaration(&tokens, startLine)
						if err != nil {
							return nil, 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
						}

						token = token.(string) + strings.Repeat("[]", dimension)

						// Get the parameter name
						parameterName, parameterNameProvided := tokens.Pop()
						if !parameterNameProvided {
							return NullVariable(), 0, CreateError("Error: Expected a parameter name", startLine+currentLine+(*scope).Index)
						}

						// Check if the parameter name is valid
						// If it is not, return an error
						if !HasValidVariableName(parameterName.(string)) {
							return NullVariable(), 0, CreateError("Error: The parameter name must be alphanumeric. Invalid parameter name \""+parameterName.(string)+"\"", startLine+currentLine+(*scope).Index)
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
							return NullVariable(), 0, CreateError("Error: Expected a comma or a closing parenthesis", startLine+currentLine+(*scope).Index)
						}

						if token.(string) == "," {
							continue
						} else if token.(string) == ")" {
							break
						} else {
							return NullVariable(), 0, CreateError("Error: Expected a comma or a closing parenthesis", startLine+currentLine+(*scope).Index)
						}

					}

				}

				// Get return type
				returnType, returnTypeProvided := tokens.Pop()

				if !returnTypeProvided {
					returnType = "null"
				} else if returnType.(string) != "val" && returnType.(string) != "int" && returnType.(string) != "float" && returnType.(string) != "bool" && returnType.(string) != "string" && returnType.(string) != "func" {
					return NullVariable(), 0, CreateError("Error: Invalid return type \""+returnType.(string)+"\"", startLine+currentLine+(*scope).Index)
				}

				// If its an array, get the dimensions
				dimensions, err := ExtractArrayDimensionFromDeclaration(&tokens, startLine+currentLine+(*scope).Index)
				if err != nil {
					return nil, 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
				}
				// Apply the dimensions to the return type
				returnType = returnType.(string) + strings.Repeat("[]", dimensions)

				// Get the block of code for the function
				funcEnded := false
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
					return NullVariable(), 0, CreateError("Error: Expected \"end "+name.(string)+"\"", startLine+(*scope).Index)
				}

				// Create the function and add it to the scope
				function := CreateVariable(CreateFunction(name.(string), startLine+currentLine, parameters, make(map[string]*Variable), returnType.(string), scope, functionCode))
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
					value, err := EvaluateExpression(scope, expression, depth, startLine+currentLine+(*scope).Index)
					if err != nil {
						return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
					}
					// Set the return value
					returnValue = value
				}

				// Check if the return value is valid type
				// If it is not, return an error
				if returnValue == (Variable{}) {
					return NullVariable(), 1, nil
				} else if returnValue.Type == (*scope).Return || (*scope).Return == "val" {
					return &returnValue, 1, nil
				} else {
					return NullVariable(), 1, CreateError("Error: Invalid return type \""+returnValue.Type+"\"", startLine+currentLine+(*scope).Index)
				}

			case "break":

				// Return until exit the for loop
				return NullVariable(), 2, nil

			case "for":

				loopBlock, nextLine, err := ParseLoopBlock(&tokens, currentLine, lines, startLine)

				if err != nil {
					return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
				}

				evaluatedCondition, err := EvaluateExpression(scope, loopBlock.Condition, depth, startLine+currentLine+(*scope).Index)
				if err != nil {
					return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
				}

				if evaluatedCondition.Type != "bool" {
					return NullVariable(), 0, CreateError("Error: Invalid condition type \""+evaluatedCondition.Type+"\"", startLine+currentLine+(*scope).Index)
				}

				for evaluatedCondition.Value.(bool) {

					forLoop := CreateFunction("for", startLine+currentLine, []Argument{}, (*scope).Variables, "val", scope, loopBlock.Code)
					returnValue, toReturn, err := forLoop.Run([]*Variable{}, map[string]*Variable{}, depth, startLine+currentLine+(*scope).Index)

					// If the code returns a value, return it
					if err != nil {
						return NullVariable(), 0, err
					}

					// If the code returns a value, return it to the caller
					if toReturn == 1 {
						return returnValue, toReturn, nil
					}

					// Exit the for loop if the break statement was called
					if toReturn == 2 {
						return returnValue, 0, nil
					}

					// Evaluate the condition
					evaluatedCondition, err = EvaluateExpression(scope, loopBlock.Condition, depth, startLine+currentLine+(*scope).Index)
					if err != nil {

						return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
					}

					if evaluatedCondition.Type != "bool" {
						return NullVariable(), 0, CreateError("Error: Invalid condition type \""+evaluatedCondition.Type+"\"", startLine+currentLine+(*scope).Index)
					}

					// Exit the loop if the condition is false
					if !evaluatedCondition.Value.(bool) {
						break
					}

				}

				currentLine = nextLine

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
							return NullVariable(), 0, CreateError("Error: \""+peeked.(string)+"\" cannot access that index because it might not be an array, startLine+currentLine+(*scope).Index", startLine+currentLine+(*scope).Index)
						}

						tokens.Pop()
						indexArray, err := (*scope).ExtractArrayValues(&tokens, depth, startLine+currentLine+(*scope).Index)
						if err != nil {
							return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
						}
						if len(indexArray) != 1 || indexArray[0].Type != "int" {
							return NullVariable(), 0, CreateError("Error: Invalid array index type \""+indexArray[0].Type+"\"", startLine+currentLine+(*scope).Index)
						}
						size := int64(len((*variable).Value.([]Variable)))
						if size == 0 {
							return NullVariable(), 0, CreateError("Error: Array is empty", startLine+currentLine+(*scope).Index)
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
						_, err := EvaluateExpression(scope, output, depth, startLine+currentLine+(*scope).Index)

						if err != nil {
							return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
						}

					} else {

						// Get the rest of the line tokens and join them to feed the variable value.

						// Remove the equal sign
						tokens.Pop()

						value := InlineQueueToString(&tokens)

						// Make sure the variable value is valid (not empty)
						if value == "" {
							return NullVariable(), 0, CreateError("Error: Variable value cannot be empty", startLine+currentLine+(*scope).Index)
						}

						// Create the new variable.
						// Evaluate the value and update the variable.
						evaluatedValue, err := EvaluateExpression(scope, value, depth, startLine+currentLine+(*scope).Index)
						if err != nil {
							return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
						}

						// Check safe assignment
						if ((*variable).Type != evaluatedValue.Type) && equal.(string) != ":=" {

							// Accept to store type[] inside val[]
							// Although, do not change the type of the variable
							if !isArrayType((*variable).Type) && !isArrayType(evaluatedValue.Type) && strings.ReplaceAll((*variable).Type, "[]", "") != "val" {
								return NullVariable(), 0, CreateError("Error:"+"Expected type "+(*scope).Variables[command.(string)].Type+" but got type "+evaluatedValue.Type+"."+" Invalid assignment type \""+evaluatedValue.Type+"\"", startLine+currentLine+(*scope).Index)
							} else {
								evaluatedValue.Type = (*variable).Type
							}

						}

						// Update the variable in the current scope.
						*variable = evaluatedValue
						// *((*scope).Variables[command.(string)]) = evaluatedValue
						if (*scope).GetVariable("_DEBUG").Type == "bool" && (*scope).GetVariable("_DEBUG").Value.(bool) {
							println("Updated variable " + command.(string) + "(" + evaluatedValue.Type + ").")
						}

					}

				} else if ExistsBuiltIn(command.(string)) {

					// Extract the function's arguments
					args, err := (*scope).ExtractFunctionArgs(&tokens, depth, startLine+currentLine+(*scope).Index)

					if err != nil {
						return NullVariable(), 0, err.AddError(CreateError("In function \""+(*scope).Name+"\"", startLine+(*scope).Index))
					}

					// Call the function

					_, err = RunBuiltIn(command.(string), args, startLine+currentLine+(*scope).Index)
					if err != nil {
						return NullVariable(), 0, err
					}

				} else {
					// Command is unknown
					return NullVariable(), 0, CreateError("Error: Unknown command \""+command.(string)+"\"", startLine+currentLine+(*scope).Index)
				}

			}

		}

	}

	return NullVariable(), 0, nil
}
