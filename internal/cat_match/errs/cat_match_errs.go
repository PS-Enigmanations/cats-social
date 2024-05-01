package errs

import (
	"errors"
	"fmt"
)

var (
	CatMatchErrGender = errors.New("Cat genders should not equal")
	CatMatchErrOwner  = errors.New("Cat is from the same owner")
)

type CatMatchError struct {
	Err error
}

func (e CatMatchError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
