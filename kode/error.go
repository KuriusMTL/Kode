package kode

import (
	"strconv"
	"strings"
)

type ErrorStack struct {
	Message   string
	Line      int
	NextError *ErrorStack
}

/**
 * Print the error stack.
 * @param e *ErrorStack - The error stack to print.
 * @return string - The error stack as a string.
 */
func (e *ErrorStack) Error() string {
	txt := e.Message
	if e.NextError != nil {
		txt = "\n└" + txt + " on line " + strconv.Itoa(e.Line) + "."
	}
	for e.NextError != nil {
		e = e.NextError
		txt = strings.ReplaceAll(txt, "\n", "\n  ") // Add indentation
		if e.NextError != nil {
			txt = "\n└" + e.Message + " on line " + strconv.Itoa(e.Line) + "." + txt
		}
	}
	return txt
}

/**
 * Create a new error stack and add an error to it.
 * @param e : *ErrorStack - The error stack to add the error to.
 * @param message : string - The error message.
 * @param line : int - The line number where the error occurred.
**/
func (e *ErrorStack) AddError(message string, line int) {
	e = &ErrorStack{message, line, e}
}
