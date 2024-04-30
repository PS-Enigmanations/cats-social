package errs

import (
	"errors"
	"fmt"
)

var (
	Unauthenticated = errors.New("Unauthenticated")
	WrongPassword   = errors.New("Wrong password")
)

type AuthError struct {
	Err error
}

func (e AuthError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
