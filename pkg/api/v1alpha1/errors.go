package v1alpha1

import (
	"errors"
)

var (
	ErrInvalidKey   = errors.New("key is invalid")
	ErrInvalidValue = errors.New("value is invalid")
	ErrDuplicateKey = errors.New("key is duplicated")
	ErrCacheEmpty   = errors.New("cache is empty")
)
