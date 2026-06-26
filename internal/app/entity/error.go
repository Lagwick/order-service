package entity

import (
	"errors"
	"net/http"
)

type AppError struct {
	err        error
	httpStatus int
}

func (e AppError) Error() string {
	return e.err.Error()
}

func (e AppError) HTTPStatus() int {
	return e.httpStatus
}

func NewAppError(err error, httpStatus int) AppError {
	return AppError{
		err:        err,
		httpStatus: httpStatus,
	}
}

var (
	ErrNotFound            = NewAppError(errors.New("not found"), http.StatusNotFound)
	ErrAlreadyExists       = NewAppError(errors.New("already exists"), http.StatusConflict)
	ErrIncorrectParameters = NewAppError(errors.New("incorrect parameters"), http.StatusBadRequest)
)
