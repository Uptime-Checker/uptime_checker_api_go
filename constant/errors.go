package constant

import "errors"

var (
	ErrGuestUserCodeExpired       = errors.New("magic link expired")
	ErrInvalidProvider            = errors.New("invalid provider value")
	ErrInvalidUserContactMode     = errors.New("invalid user contact mode")
	ErrInvalidUserToken           = errors.New("invalid user token")
	ErrInvalidSubscriptionStatus  = errors.New("invalid subscription status")
	ErrInvalidJobStatus           = errors.New("invalid job status")
	ErrExpiresAtInThePast         = errors.New("expires at in the past")
	ErrInvalidMonitorType         = errors.New("invalid monitor type")
	ErrInvalidMonitorStatus       = errors.New("invalid monitor status")
	ErrInvalidMonitorBodyFormat   = errors.New("invalid monitor body format")
	ErrInvalidAlarmPolicy         = errors.New("invalid alarm policy")
	ErrInvalidAssertionSource     = errors.New("invalid assertion source")
	ErrInvalidAssertionComparison = errors.New("invalid assertion comparison")
	ErrInvalidErrorLogType        = errors.New("invalid error log type")
)

var (
	ErrUpgradeSubscription = errors.New("upgrade subscription")
	ErrUpgradePermission   = errors.New("upgrade permission")
	ErrSubscriptionExpired = errors.New("subscription expired")
)

const (
	ErrRemoteServerMaxTimeout = "remote server exceeded max timeout"
	ErrResponseMalformed      = "response malformed"
)
