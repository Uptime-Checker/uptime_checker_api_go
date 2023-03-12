package constant

import "errors"

var (
	ErrGuestUserCodeExpired      = errors.New("magic link expired")
	ErrInvalidProvider           = errors.New("invalid provider value")
	ErrInvalidUserContactMode    = errors.New("invalid user contact mode")
	ErrInvalidUserToken          = errors.New("invalid user token")
	ErrInvalidSubscriptionStatus = errors.New("invalid subscription status")
	ErrExpiresAtInThePast        = errors.New("expires at in the past")
)

var (
	ErrUpgradeSubscription = errors.New("upgrade subscription")
	ErrUpgradePermission   = errors.New("upgrade permission")
	ErrSubscriptionExpired = errors.New("subscription expired")
)
