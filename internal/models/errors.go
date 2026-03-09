package models

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrNoUpdateFields       = errors.New("no fields to update")
	ErrIncorrectDateRange   = errors.New("incorrect date range")
)
