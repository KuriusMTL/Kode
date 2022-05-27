package kode

import (
	"strconv"
	"strings"
)

type ErrorStack struct {
	Message   string
	Line      int
	NextError *ErrorStack
	Index     int
}

var MAX_ERROR_DEPTH = 5

/**
 * Print the error stack.
 * @param e *ErrorStack - The error stack to print.
 * @return string - The error stack as a string.
 */
func (e *ErrorStack) Error() string {

	// Empty error stack.
	if e == nil {
		return ""
	}

	txt := e.Message + " on line " + strconv.Itoa(e.Line) + "."

	i := 1
	for e.NextError != nil {
		e = e.NextError
		txt += "\n" + strings.Repeat("  ", i) + "└" + e.Message + " on line " + strconv.Itoa(e.Line) + "."
		i++
	}
	return txt
}

/**
 * Create a new error stack.
 * @param message : string - The error message.
 * @param line : int - The line number where the error occurred.
 * @return *ErrorStack - The new error stack.
**/
func CreateError(message string, line int) *ErrorStack {
	return &ErrorStack{message, line, nil, 0}
}

/**
 * Add a new error to the error stack.
 * @param e *ErrorStack - The error stack.
 * @param err : *ErrorStack - The error to add.
**/
func (e *ErrorStack) AddError(err *ErrorStack) *ErrorStack {

	if (*e).Index == MAX_ERROR_DEPTH {
		// (*err).NextError = CreateError("(...)", 0)
		// return e
		new_err := CreateError("(...)", 0)
		new_err.NextError = e
		new_err.Index = MAX_ERROR_DEPTH + 1
		return new_err
	}

	if (*e).Index > MAX_ERROR_DEPTH {
		return e
	}

	(*err).NextError = e
	(*err).Index = (*e).Index + 1
	return err
}
