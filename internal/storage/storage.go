package storage

import "github.com/nsmak/bannersRotation/internal/app"

var (
	ErrBannerNotFound           = NewError("banner not found", nil)
	ErrSlotNotFound             = NewError("slot not found", nil)
	ErrBannerInSlotNotFound     = NewError("banner in slot not found", nil)
	ErrBannerInSlotAlreadyExist = NewError("banner in slot already exist", nil)
)

type Error struct {
	app.BaseError
}

func NewError(msg string, err error) *Error {
	return &Error{BaseError: app.BaseError{Message: msg, Err: err}}
}
