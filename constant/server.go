package constant

const UTCTimeZone = "UTC"

const (
	APIKeyHeader        = "X_API_KEY"
	OriginalIPHeader    = "X-Forwarded-For"
	AuthorizationHeader = "Authorization"
)

const AuthScheme = "Bearer"

const (
	MaxRequestPerMinute            = 10
	ServerShutdownTimeoutInSeconds = 5
	SentryTraceSampleRate          = 0.2
	CronCheckIntervalInSeconds     = 30
	WatchDogCheckIntervalInSeconds = 5
	MonitorStartDelayInSeconds     = 10
)

// Environment type
type Environment string

// List of environments
const (
	EnvironmentDev  Environment = "DEV"
	EnvironmentProd Environment = "PROD"
)
