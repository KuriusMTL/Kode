package kode

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/**
 * Check if an embedded function exists for Kode.
 * @param name : string - The name of the function.
 * @return boolean - True if the function exists, false otherwise.
**/
func ExistsBuiltIn(name string) bool {
	switch name {
	case "print":
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
	case "len":
		return true
	case "random":
		return true
	case "append":
		return true
	case "truncate":
		return true
	case "round":
		return true
	case "sqrt":
		return true
	case "isNumeric":
		return true
	case "isAlphaNumeric":
		return true
	case "toUnicode":
		return true
	case "fromUnicode":
		return true
	default:
		return false
	}
}

/**
 * Run a Kode embedded function.
 * @param name : string - The name of the function.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func RunBuiltIn(name string, args []*Variable, startLine int) (*Variable, *ErrorStack) {
	switch name {
	case "print":
		return Print(args, startLine)
	case "toString":
		return ToString(args, startLine)
	case "toInt":
		return ToInt(args, startLine)
	case "toFloat":
		return ToFloat(args, startLine)
	case "yell":
		return Yell(args, startLine)
	case "whisper":
		return Whisper(args, startLine)
	case "typeOf":
		return TypeOf(args, startLine)
	case "len":
		return Len(args, startLine)
	case "random":
		return Random(args, startLine)
	case "append":
		return Append(args, startLine)
	case "truncate":
		return Truncate(args, startLine)
	case "round":
		return Round(args, startLine)
	case "sqrt":
		return Sqrt(args, startLine)
	case "isNumeric":
		return IsNumeric(args, startLine)
	case "isAlphaNumeric":
		return IsAlphaNumeric(args, startLine)
	case "toUnicode":
		return ToUnicode(args, startLine)
	case "fromUnicode":
		return FromUnicode(args, startLine)
	default:
		return NullVariable(), nil
	}
}

/**
 * Print variable(s) to the console.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func Print(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	msg := ""
	for i, arg := range args {

		switch arg.Type {
		case "string":
			msg += arg.Value.(string)
			break
		case "int":
			msg += strconv.FormatInt(arg.Value.(int64), 10)
			break
		case "float":
			msg += fmt.Sprintf("%g", arg.Value.(float64))
			break
		case "bool":
			msg += strconv.FormatBool(arg.Value.(bool))
			break
		case "null":
			msg += "null"
			break
		case "function":
			msg += "function"
			break
		default:
			if isArrayType(arg.Type) {
				msg += "array"
			} else {
				msg += "unknown"
			}
			break
		}
		if i != 0 || i != len(args)-1 {
			msg += " "
		}
	}
	fmt.Println(msg)
	variable := CreateVariable(msg)
	return &variable, nil
}

/**
 * Convert a variable to a string.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func ToString(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) == 1 {
		switch args[0].Type {
		case "string":
			variable := CreateVariable(args[0].Value.(string))
			return &variable, nil
		case "int":
			variable := CreateVariable(strconv.FormatInt(args[0].Value.(int64), 10))
			return &variable, nil
		case "float":
			variable := CreateVariable(fmt.Sprintf("%g", args[0].Value.(float64)))
			return &variable, nil
		case "bool":
			variable := CreateVariable(strconv.FormatBool(args[0].Value.(bool)))
			return &variable, nil
		case "null":
			variable := CreateVariable("null")
			return &variable, nil
		case "function":
			variable := CreateVariable("function")
			return &variable, nil
		default:
			return NullVariable(), nil
		}
	} else {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"toString\"", startLine)
	}
}

/**
 * Convert a variable to an int.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func ToInt(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"toInt\"", startLine)
	}

	switch args[0].Type {
	case "string":
		i, err := strconv.ParseInt(args[0].Value.(string), 10, 64)
		if err != nil {
			return NullVariable(), CreateError("Error: String is not a number or is too large to be converted to an int for \"toInt\"", startLine)
		}
		variable := CreateVariable(i)
		return &variable, nil
	case "int":
		variable := CreateVariable(args[0].Value.(int64))
		return &variable, nil
	case "float":
		i := int64(args[0].Value.(float64))
		variable := CreateVariable(i)
		return &variable, nil
	default:
		return NullVariable(), CreateError("Error: Argument must be a string or a number for \"toInt\"", startLine)
	}

}

/**
 * Convert a variable to a float.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func ToFloat(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"toFloat\"", startLine)
	}

	switch args[0].Type {
	case "string":
		f, err := strconv.ParseFloat(args[0].Value.(string), 64)
		if err != nil {
			return NullVariable(), CreateError("Error: String is not a number or is too large to be converted to a float for \"toFloat\"", startLine)
		}
		variable := CreateVariable(f)
		return &variable, nil
	case "int":
		f := float64(args[0].Value.(int64))
		variable := CreateVariable(f)
		return &variable, nil
	case "float":
		variable := CreateVariable(args[0].Value.(float64))
		return &variable, nil
	default:
		return NullVariable(), CreateError("Error: Argument must be a string or a number for \"toFloat\"", startLine)
	}
}

/**
 * Lowercase a string.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func Whisper(args []*Variable, startLine int) (*Variable, *ErrorStack) {

	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 string argument for \"whisper\"", startLine)
	}

	if args[0].Type != "string" {
		return NullVariable(), CreateError("Error: Argument must be a string for \"whisper\"", startLine)
	}

	// Lowercase the string
	s := strings.ToLower(args[0].Value.(string))
	variable := CreateVariable(s)
	return &variable, nil
}

/**
 * Uppercase a string.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func Yell(args []*Variable, startLine int) (*Variable, *ErrorStack) {

	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 string argument for \"yell\"", startLine)
	}

	if args[0].Type != "string" {
		return NullVariable(), CreateError("Error: Argument must be a string for \"yell\"", startLine)
	}

	// Uppercase the string
	s := strings.ToUpper(args[0].Value.(string))
	variable := CreateVariable(s)
	return &variable, nil
}

/**
 * Get the type of a variable.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
 */
func TypeOf(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"typeOf\"", startLine)
	}

	variable := CreateVariable(args[0].Type)
	return &variable, nil
}

func Len(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"len\"", startLine)
	}

	if !isArrayType(args[0].Type) && args[0].Type != "string" {
		return NullVariable(), CreateError("Error: Argument must be an array or a string for \"len\"", startLine)
	}

	if args[0].Type == "string" {
		variable := CreateVariable(int64(len(args[0].Value.(string))))
		return &variable, nil
	} else {
		variable := CreateVariable(int64(len(args[0].Value.([]Variable))))
		return &variable, nil
	}
}

func Random(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) > 0 {
		return NullVariable(), CreateError("Error: Expected 0 arguments for \"random\"", startLine)
	}

	rand.Seed(time.Now().UnixNano())
	variable := CreateVariable(rand.Float64())
	return &variable, nil
}

func Append(args []*Variable, startLine int) (*Variable, *ErrorStack) {

	if len(args) != 2 {
		return NullVariable(), CreateError("Error: Expected 2 arguments for \"append\"", startLine)
	}

	if !isArrayType(args[0].Type) {
		return NullVariable(), CreateError("Error: Argument 1 must be an array for \"append\"", startLine)
	}

	// Get allowed types
	allowedType := strings.ReplaceAll(args[0].Type, "[]", "")

	if allowedType != "val" {
		if args[1].Type != allowedType {
			return NullVariable(), CreateError("Error: Argument 2 must be a "+allowedType+" for \"append\"", startLine)
		}

		// Append the value
		array := append(args[0].Value.([]Variable), *(args[1]))

		// Return the array
		variable := CreateVariable(array)
		return &variable, nil
	} else {
		// Append the value
		array := append(args[0].Value.([]Variable), *(args[1]))

		// Return the array
		variable := CreateVariable(array)
		return &variable, nil
	}

}

func Truncate(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 2 {
		return NullVariable(), CreateError("Error: Expected 2 arguments for \"truncate\"", startLine)
	}

	if !isArrayType(args[0].Type) {
		return NullVariable(), CreateError("Error: Argument 1 must be an array for \"truncate\"", startLine)
	}

	if args[1].Type != "int" {
		return NullVariable(), CreateError("Error: Argument 2 must be an int for \"truncate\"", startLine)
	}

	// Get size and index
	array := args[0]
	size, err := GetArraySize(array, startLine)
	if err != nil {
		return NullVariable(), err
	}

	index := args[1].Value.(int64) % size
	if index < 0 { // Handle negative indexes
		index += size
	}

	// Truncate the array
	// Remove the element

	// Copy the array
	newArray := make([]Variable, size-1)
	i := int64(0)
	removed := false
	for i < int64(len(newArray)) {

		if i != index || removed {
			if removed {
				newArray[i] = array.Value.([]Variable)[i+1]
			} else {
				newArray[i] = array.Value.([]Variable)[i]
			}
			i++
		} else {
			removed = true
		}

	}

	variable := CreateVariable(newArray)
	return &variable, nil
}

func Round(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"round\"", startLine)
	}

	if args[0].Type != "float" {
		return NullVariable(), CreateError("Error: Argument must be a float for \"round\"", startLine)
	}

	// Round the float
	f := args[0].Value.(float64)
	variable := CreateVariable(math.Round(f))
	return &variable, nil
}

func Sqrt(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"sqrt\"", startLine)
	}

	if args[0].Type != "float" && args[0].Type != "int" {
		return NullVariable(), CreateError("Error: Argument must be a float or an int for \"sqrt\"", startLine)
	}

	// Square root the float
	if args[0].Type == "float" {
		f := args[0].Value.(float64)

		if f < 0 {
			return NullVariable(), CreateError("Error: Argument must be a positive float for \"sqrt\"", startLine)
		}

		variable := CreateVariable(math.Sqrt(f))
		return &variable, nil
	} else {
		i := args[0].Value.(int64)

		if i < 0 {
			return NullVariable(), CreateError("Error: Argument must be a positive int for \"sqrt\"", startLine)
		}

		variable := CreateVariable(math.Sqrt(float64(i)))
		return &variable, nil
	}
}

func IsNumeric(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"isNumeric\"", startLine)
	}

	if args[0].Type != "string" {
		return NullVariable(), CreateError("Error: Argument must be a string for \"isNumeric\"", startLine)
	}

	// Check if the string is numeric
	str := args[0].Value.(string)
	_, err := strconv.ParseFloat(str, 64)
	if err == nil {
		variable := CreateVariable(true)
		return &variable, nil
	} else {
		variable := CreateVariable(false)
		return &variable, nil
	}
}

func IsAlphaNumeric(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"isAlphaNumeric\"", startLine)
	}

	if args[0].Type != "string" {
		return NullVariable(), CreateError("Error: Argument must be a string for \"isAlphaNumeric\"", startLine)
	}

	for _, c := range args[0].Value.(string) {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			variable := CreateVariable(false)
			return &variable, nil
		}
	}

	variable := CreateVariable(true)
	return &variable, nil

}

func ToUnicode(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"toUnicode\"", startLine)
	}

	if args[0].Type != "string" {
		return NullVariable(), CreateError("Error: Argument must be a string for \"toUnicode\"", startLine)
	}

	if len(args[0].Value.(string)) != 1 {
		return NullVariable(), CreateError("Error: String argument must be of size 1 for \"toUnicode\"", startLine)
	}

	// Convert the string to unicode
	variable := CreateVariable(int64(args[0].Value.(string)[0]))
	return &variable, nil
}

func FromUnicode(args []*Variable, startLine int) (*Variable, *ErrorStack) {
	if len(args) != 1 {
		return NullVariable(), CreateError("Error: Expected 1 argument for \"fromUnicode\"", startLine)
	}

	if args[0].Type != "int" {
		return NullVariable(), CreateError("Error: Argument must be an integer for \"fromUnicode\"", startLine)
	}

	// Convert the unicode to string
	variable := CreateVariable(string(rune(args[0].Value.(int64))))
	return &variable, nil
}
