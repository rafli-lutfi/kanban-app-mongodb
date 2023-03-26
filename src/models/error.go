package models

import "errors"

var (
	// method is request not allowed
	ErrMethodNotAllowed = errors.New("method is not allowed")

	// Unauthorized user
	ErrUnauthorized = errors.New("Unauthorized user")

	// failed to decode body
	ErrFailedDecodeBody = errors.New("failed to decode body")

	// there is empty data on body
	ErrEmptyDataBody = errors.New("there is empty data on body")

	ErrQueryParamEmpty = errors.New("query param is empty")

	// Email or password is not match
	ErrEmailPasswordNotMatched = errors.New("email or password is not matched")

	// Record not found
	ErrRecordNotFound = errors.New("record not found")
)
