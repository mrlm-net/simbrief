package types

import "errors"

// Common validation errors
var (
	ErrMissingOrigin      = errors.New("origin airport (orig) is required")
	ErrMissingDestination = errors.New("destination airport (dest) is required")
	ErrMissingAircraft    = errors.New("aircraft type (type) is required")
	ErrMissingUserID      = errors.New("user ID or username is required")
	ErrInvalidUserID      = errors.New("invalid user ID format")
	ErrInvalidAPIKey      = errors.New("invalid or missing API key")
)
