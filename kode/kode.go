package kode

import "strings"

func Run(code string) error {

	// Return if the code is empty.
	if code == "" {
		return nil
	}

	// Create a new main scope.
	scope := CreateFunction([]Argument{}, nil, "null", nil, strings.ReplaceAll(code, "\r", " "))

	// Enter the main scope.
	_, _, err := scope.Run([]*Variable{}, map[string]*Variable{})

	return err

}
