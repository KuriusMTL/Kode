package kode

import (
	"regexp"
)

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
 * Check if a variable exists within a scope (funcion).
 * @param name : string - The name of the variable.
 * @return bool - True if the variable exists, false otherwise.
 */
func (function *Function) VariableExists(name string) bool {
	val, ok := (*function).Variables[name]
	if ok && val != nil {
		return true
	} else {
		return false
	}
}

/**
 * Get the value of a variable. If the variable does not exist, return a nil.
 * @param name : string - The name of the variable.
 * @return Variable - The value of the variable.
 */
func (function *Function) GetVariable(name string) *Variable {
	return (*function).Variables[name]
}

/**
 * Verify if a variable has a valid name
 * @param name : string - The name of the variable.
 * @return bool - True if the variable has a valid name, false otherwise.
 */
func HasValidVariableName(name string) bool {
	return varFormat.MatchString(name) && !IsReservedWord(name)
}

func IsReservedWord(name string) bool {
	switch name {
	case "null":
		return true
	case "true":
		return true
	case "false":
		return true
	case "if":
		return true
	case "else":
		return true
	case "val":
		return true
	case "string":
		return true
	case "int":
		return true
	case "float":
		return true
	case "bool":
		return true
	case "func":
		return true
	case "return":
		return true
	case "for":
		return true
	case "break":
		return true
	case "self":
		return true
	case "super":
		return true
	case "new":
		return true
	case "print":
		return true
	case "len":
		return true
	case "toString":
		return true
	case "toInt":
		return true
	case "toFloat":
		return true
	case "yell":
		return true
	case "whisper":
		return true
	case "typeOf":
		return true
	case "random":
		return true
	case "truncate":
		return true
	case "append":
		return true
	case "round":
		return true
	default:
		return false
	}
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
	case Function:
		return "func"
	case int:
		return "illegal_int"
	case []Variable:
		return EvaluateArrayType(value.([]Variable)) // i.e. val[], int[], float[], string[], bool[], func[]
	default:
		return "null"
	}
}

/**
 * Get the default value of a variable type.
 * @param typeName : string - The name of the variable type.
 * @return interface{} - The default value of the variable type.
 */
func GetDefaultValue(typeName string) interface{} {
	switch typeName {
	case "int":
		return int64(0)
	case "float":
		return float64(0)
	case "string":
		return ""
	case "bool":
		return false
	case "func":
		return nil
	default:
		return nil
	}
}

/**
 * Get a null variable.
 * @return *Variable - The null variable.
 */
func NullVariable() *Variable {
	return &Variable{
		Value: nil,
		Type:  "null",
	}
}
