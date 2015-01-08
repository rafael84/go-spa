package api

import "errors"

var (
	EInternalServerError = errors.New("Internal server error")
	ENotImplemented      = errors.New("Not Implemented")
)
