package kode

import "regexp"

type Variable struct {
	Value interface{}
	Type  string
}

var varFormat, _ = regexp.Compile("^[a-zA-Z_][a-zA-Z0-9_]*$")

/**
 * Create a new variable.
 * @param value : interface{} - The value of the variable.
 * @return Variable - The new variable.
 */
func CreateVariable(value interface{}) Variable {
	return Variable{
		Value: value,
		Type:  EvaluateType(value),
	}
}

/**
 * Check if a variable exists.
 * @param name : string - The name of the variable.
 * @return bool - True if the variable exists, false otherwise.
 */
func (function *Function) VariableExists(name string) bool {
	_, ok := (*function).Variables[name]
	return ok
}

/**
 * Get the value of a variable. If the variable does not exist, return a nil.
 * @param name : string - The name of the variable.
 * @return Variable - The value of the variable.
 */
func (function *Function) GetVariable(name string) Variable {
	return (*function).Variables[name]
}

/**
 * Verify if a variable has a valid name
 * @param name : string - The name of the variable.
 * @return bool - True if the variable has a valid name, false otherwise.
 */
func HasValidVariableName(name string) bool {
	return varFormat.MatchString(name)
}

/**
 * Get the value of a variable. If the variable does not exist, "unknown" is returned.
 * @param value : interface{} - The value of the variable.
 * @return string - The type of the variable.
 */
func EvaluateType(value interface{}) string {
	switch value.(type) {
	case int64:
		return "int"
	case float64:
		return "float"
	case string:
		return "string"
	case bool:
		return "bool"
	case int:
		return "illegal_int"
	default:
		return "unknown"
	}
}
