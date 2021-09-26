package mediator

import "fmt"

var (
	// ErrorWrongPassword is when the user enters a wrong password
	ErrorWrongPassword = fmt.Errorf("wrong password")
	// ErrorDuplicateUser is when the username has been taken
	ErrorDuplicateUser = fmt.Errorf("duplicate user")
	// ErrorNoSuchUser is when the username has been taken
	ErrorNoSuchUser = fmt.Errorf("unknown user")
	// ErrorExpiredSession is an expired session
	ErrorExpiredSession = fmt.Errorf("expired session")
)
