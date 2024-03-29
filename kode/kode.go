package kode

import (
	"errors"
	"strings"
)

func Run(code string) error {

	// Return if the code is empty.
	if code == "" {
		return nil
	}

	// Create a new main scope.
	_debug := CreateVariable(false)               // DEBUG variable prints debug info to console
	_max_recursion := CreateVariable(int64(5000)) // Max recursion depth for functions
	scope := CreateFunction("main", 1, []Argument{}, map[string]*Variable{"_DEBUG": &_debug, "_MAX_RECURSION": &_max_recursion}, "null", nil, strings.ReplaceAll(code, "\r", " "))

	// Enter the main scope.
	_, _, err := scope.Run([]*Variable{}, map[string]*Variable{}, 0, 0)

	// Convert the Kode stack error and return it.

	if err != nil {
		return errors.New((*err).Error())
	}

	return nil

}
