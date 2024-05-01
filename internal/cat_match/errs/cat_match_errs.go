package errs

import (
	"errors"
	"fmt"
)

var (
	CatMatchErrGender         = errors.New("Cat genders should not equal")
	CatMatchErrOwner          = errors.New("Cat is from the same owner")
	CatMatchErrNotFound       = errors.New("Cat not found")
	CatMatchErrOwnerNotFound  = errors.New("Cat owner not found")
	CatMatchErrInvalidAuthor  = errors.New("Cat owner is equal with matches")
	CatMatchErrAlreadyMatched = errors.New("Cat already matched")
)

type CatMatchError struct {
	Err error
}

func (e CatMatchError) Error() error {
	return fmt.Errorf(e.Err.Error())
}
