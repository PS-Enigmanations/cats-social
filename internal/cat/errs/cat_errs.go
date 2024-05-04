package errs

import (
	"errors"
	"fmt"
)

var (
	CatErrNotFound       = errors.New("Cat not found")
	CatErrSexNotEditable = errors.New("Cat sex is not editable, already requested to match")
)

type CatError struct {
	Err error
}

func (e CatError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
