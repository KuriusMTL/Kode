package kode

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/**
 * Check if an embedded function exists for Kode.
 * @param name : string - The name of the function.
 * @return boolean - True if the function exists, false otherwise.
**/
func ExistsIncluded(name string) bool {
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
func RunIncluded(name string, args []*Variable) (*Variable, error) {
	switch name {
	case "print":
		return Print(args)
	case "toString":
		return ToString(args)
	case "toInt":
		return ToInt(args)
	case "toFloat":
		return ToFloat(args)
	case "yell":
		return Yell(args)
	case "whisper":
		return Whisper(args)
	case "typeOf":
		return TypeOf(args)
	case "len":
		return Len(args)
	case "random":
		return Random(args)
	case "append":
		return Append(args)
	case "truncate":
		return Truncate(args)
	case "round":
		return Round(args)
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
func Print(args []*Variable) (*Variable, error) {
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
			msg += "unknown"
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
func ToString(args []*Variable) (*Variable, error) {
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
		return NullVariable(), errors.New("Error: Expected 1 argument for \"toString\"")
	}
}

/**
 * Convert a variable to an int.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func ToInt(args []*Variable) (*Variable, error) {
	if len(args) != 1 {
		return NullVariable(), errors.New("Error: Expected 1 argument for \"toInt\"")
	}

	switch args[0].Type {
	case "string":
		i, err := strconv.ParseInt(args[0].Value.(string), 10, 64)
		if err != nil {
			return NullVariable(), errors.New("Error: String is not a number or is too large to be converted to an int for \"toInt\"")
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
		return NullVariable(), errors.New("Error: Argument must be a string or a number for \"toInt\"")
	}

}

/**
 * Convert a variable to a float.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func ToFloat(args []*Variable) (*Variable, error) {
	if len(args) != 1 {
		return NullVariable(), errors.New("Error: Expected 1 argument for \"toFloat\"")
	}

	switch args[0].Type {
	case "string":
		f, err := strconv.ParseFloat(args[0].Value.(string), 64)
		if err != nil {
			return NullVariable(), errors.New("Error: String is not a number or is too large to be converted to a float for \"toFloat\"")
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
		return NullVariable(), errors.New("Error: Argument must be a string or a number for \"toFloat\"")
	}
}

/**
 * Lowercase a string.
 * @param args :[]*Variable - The arguments to the function.
 * @return *Variable - The result of the function.
 * @return error - The error if one occurs.
**/
func Whisper(args []*Variable) (*Variable, error) {

	if len(args) != 1 {
		return NullVariable(), errors.New("Error: Expected 1 argument for \"whisper\"")
	}

	if args[0].Type != "string" {
		return NullVariable(), errors.New("Error: Argument must be a string for \"whisper\"")
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
func Yell(args []*Variable) (*Variable, error) {

	if len(args) != 1 {
		return NullVariable(), errors.New("Error: Expected 1 argument ofr \"yell\"")
	}

	if args[0].Type != "string" {
		return NullVariable(), errors.New("Error: Argument must be a string for \"yell\"")
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
func TypeOf(args []*Variable) (*Variable, error) {
	if len(args) != 1 {
		return NullVariable(), errors.New("Error: Expected 1 argument for \"typeof\"")
	}

	variable := CreateVariable(args[0].Type)
	return &variable, nil
}

func Len(args []*Variable) (*Variable, error) {
	if len(args) != 1 {
		return NullVariable(), errors.New("Error: Expected 1 argument for \"len\"")
	}

	if !isArrayType(args[0].Type) && args[0].Type != "string" {
		return NullVariable(), errors.New("Error: Expected an array or string as the argument for \"len\"")
	}

	if args[0].Type == "string" {
		variable := CreateVariable(int64(len(args[0].Value.(string))))
		return &variable, nil
	} else {
		variable := CreateVariable(int64(len(args[0].Value.([]Variable))))
		return &variable, nil
	}
}

func Random(args []*Variable) (*Variable, error) {
	if len(args) > 0 {
		return NullVariable(), errors.New("Error: Expected 0 argument for \"random\"")
	}

	rand.Seed(time.Now().UnixNano())
	variable := CreateVariable(rand.Float64())
	return &variable, nil
}

func Append(args []*Variable) (*Variable, error) {

	if len(args) != 2 {
		return NullVariable(), errors.New("Error: Expected 2 arguments (array, val) for \"append\"")
	}

	if !isArrayType(args[0].Type) {
		return NullVariable(), errors.New("Error: Expected an array as the first argument for \"append\"")
	}

	// Get allowed types
	allowedType := strings.ReplaceAll(args[0].Type, "[]", "")

	if allowedType != "val" {
		if args[1].Type != allowedType {
			return NullVariable(), errors.New("Error: Expected a " + allowedType + " as the second argument for \"append\"")
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

func Truncate(args []*Variable) (*Variable, error) {
	if len(args) != 2 {
		return NullVariable(), errors.New("Error: Expected 2 arguments (array, int) for \"truncate\"")
	}

	if !isArrayType(args[0].Type) {
		return NullVariable(), errors.New("Error: Expected an array as the first argument for \"truncate\"")
	}

	if args[1].Type != "int" {
		return NullVariable(), errors.New("Error: Expected second argument must be an int for \"truncate\"")
	}

	// Get size and index
	array := args[0]
	size, err := GetArraySize(array)
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

func Round(args []*Variable) (*Variable, error) {
	if len(args) != 1 {
		return NullVariable(), errors.New("Error: Expected 1 argument (float) for \"round\"")
	}

	if args[0].Type != "float" {
		return NullVariable(), errors.New("Error: Argument must be a float for \"round\"")
	}

	// Round the float
	f := args[0].Value.(float64)
	variable := CreateVariable(math.Round(f))
	return &variable, nil
}
