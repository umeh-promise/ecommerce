package utils

import "errors"

var (
	ErrorNotFound          = errors.New("resource not found")
	ErrorInvalidID         = errors.New("invalid post id")
	ErrorDuplicateEmail    = errors.New("a user with that email already exists")
	ErrorDuplicateUsername = errors.New("a user with that username already exists")
)
