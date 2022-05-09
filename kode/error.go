package kode

type ErrorStack struct {
	Message   string
	Line      int
	NextError *ErrorStack
}

func (e *ErrorStack) Error() string {
	return ""
}
