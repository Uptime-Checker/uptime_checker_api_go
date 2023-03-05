package constant

import "errors"

var (
	ErrGuestUserCodeExpired   = errors.New("magic link expired")
	ErrInvalidProvider        = errors.New("invalid provider value")
	ErrInvalidUserContactMode = errors.New("invalid user contact mode")
	ErrInvalidUserToken       = errors.New("invalid user token")
)
