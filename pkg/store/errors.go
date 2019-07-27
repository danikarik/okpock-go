package store

import "errors"

var (
	// ErrNilStruct raises when db struct is nil.
	ErrNilStruct = errors.New("store: nil struct given")
	// ErrZeroID raises when db struct has zero id.
	ErrZeroID = errors.New("store: id not set")
	// ErrEmptyQueryParam raises when query param is empty.
	ErrEmptyQueryParam = errors.New("store: empty query param")
	// ErrZeroRowsAffected raises when zero rows affected.
	ErrZeroRowsAffected = errors.New("store: zero rows affected")
	// ErrNotFound raises when record not found.
	ErrNotFound = errors.New("store: record not found")
	// ErrWrongPassword raises when input password doesn't match.
	ErrWrongPassword = errors.New("store: wrong password")
)
