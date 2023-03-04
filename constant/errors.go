package constant

import "errors"

var (
	ErrGuestUserCodeExpired = errors.New("magic link expired")
	ErrInvalidProvider      = errors.New("invalid provider value")
)
