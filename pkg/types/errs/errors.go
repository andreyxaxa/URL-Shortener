package errs

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrInvalidInterval   = errors.New("invalid interval")
	ErrAliasAlreadyTaken = errors.New("alias already taken")
)
