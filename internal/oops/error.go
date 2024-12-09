package oops

import "errors"

// Database errors
// ErrNoData ...
var ErrNoData = errors.New("no data in db")

// Business erros

// Custom error
type DBError struct {
	Err error
	ID  string
}

func NewError(err error, id string) *DBError {
	return &DBError{
		Err: err,
		ID:  id,
	}
}
