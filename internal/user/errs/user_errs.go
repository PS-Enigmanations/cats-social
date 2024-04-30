package errs

import (
	"errors"
	"fmt"
)

var (
	UserErrEmailExist         = errors.New("Email already exists")
	UserErrEmailInvalidFormat = errors.New("Email format invalid")
	UserErrNotFound           = errors.New("User not found")
)

type UserError struct {
	Err error
}

func (e UserError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
