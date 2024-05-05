package errs

import (
	"errors"
	"fmt"
)

var (
	Unauthenticated = errors.New("Unauthenticated")
	WrongPassword   = errors.New("Wrong password")
)

type SessionError struct {
	Err error
}

func (e SessionError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
