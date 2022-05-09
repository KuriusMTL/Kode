package kode

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		return NullVariable(), errors.New("Too many arguments")
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
		return NullVariable(), errors.New("Too many arguments")
	}

	switch args[0].Type {
	case "string":
		i, err := strconv.ParseInt(args[0].Value.(string), 10, 64)
		if err != nil {
			return NullVariable(), errors.New("String is not a number")
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
		return NullVariable(), errors.New("Argument must be a string or a number")
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
		return NullVariable(), errors.New("Too many arguments")
	}

	switch args[0].Type {
	case "string":
		f, err := strconv.ParseFloat(args[0].Value.(string), 64)
		if err != nil {
			return NullVariable(), errors.New("String is not a number")
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
		return NullVariable(), errors.New("Argument must be a string or a number")
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
		return NullVariable(), errors.New("Too many arguments")
	}

	if args[0].Type != "string" {
		return NullVariable(), errors.New("Argument must be a string")
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
		return NullVariable(), errors.New("Too many arguments")
	}

	if args[0].Type != "string" {
		return NullVariable(), errors.New("Argument must be a string")
	}

	// Uppercase the string
	s := strings.ToUpper(args[0].Value.(string))
	variable := CreateVariable(s)
	return &variable, nil
}
