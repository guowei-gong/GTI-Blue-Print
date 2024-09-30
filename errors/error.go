package errors

import "errors"

var (
	ErrInvalidDecoder        = New("invalid decoder")
	ErrInvalidScanner        = New("invalid scanner")
	ErrInvalidFormat         = New("invalid format")
	ErrNoOperationPermission = New("no operation permission")
)

// New Wrapping for errors.New standard library
func New(text string) error { return errors.New(text) }
