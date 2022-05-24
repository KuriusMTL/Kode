package kode

import "strings"

func Run(code string) error {

	// Return if the code is empty.
	if code == "" {
		return nil
	}

	// Create a new main scope.
	_debug := CreateVariable(false)
	_max_recursion := CreateVariable(5000)
	scope := CreateFunction([]Argument{}, map[string]*Variable{"_DEBUG": &_debug, "_MAX_RECURSION": &_max_recursion}, "null", nil, strings.ReplaceAll(code, "\r", " "))

	// Enter the main scope.
	_, _, err := scope.Run([]*Variable{}, map[string]*Variable{})

	return err

}
