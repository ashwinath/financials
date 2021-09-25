package mediator

import "fmt"

var (
	// ErrorWrongPassword is when the user enters a wrong password
	ErrorWrongPassword = fmt.Errorf("wrong password")
)
