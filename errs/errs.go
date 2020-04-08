package errs

import "errors"

var (
	ErrIsNotInt   = errors.New("the value is not an integer")
	ErrIsNotFloat = errors.New("the value is not a valid float")
)
