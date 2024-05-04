package errs

import (
	"errors"
	"fmt"
)

var (
	CatErrNotFound = errors.New("Cat not found")
)

type CatError struct {
	Err error
}

func (e CatError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
