package storage

import "github.com/nsmak/bannersRotation/internal/app"

var (
	ErrObjectNotFound = NewError("object not found", nil)
)

type Error struct {
	app.BaseError
}

func NewError(msg string, err error) *Error {
	return &Error{BaseError: app.BaseError{Message: msg, Err: err}}
}
