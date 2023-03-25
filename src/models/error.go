package models

import "errors"

var (
	// method is request not allowed
	ErrMethodNotAllowed = errors.New("method is not allowed")

	// failed to decode body
	ErrFailedDecodeBody = errors.New("failed to decode body")

	// there is empty data on body
	ErrEmptyDataBody = errors.New("there is empty data on body")
)
