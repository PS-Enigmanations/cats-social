package errs

import (
	"errors"
	"fmt"
)

var (
	UserErrEmailExist         = errors.New("Email already exists")
	UserErrEmailInvalidFormat = errors.New("Email format invalid")
)

type UserError struct {
	Err error
}

func (e UserError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
