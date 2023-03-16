package constant

const UTCTimeZone = "UTC"

const (
	APIKeyHeader        = "X_API_KEY"
	OriginalIPHeader    = "X-Forwarded-For"
	AuthorizationHeader = "Authorization"
)

const AuthScheme = "Bearer"

const MaxRequestPerMinute = 10
const ServerShutdownTimeout = 5

// Environment type
type Environment string

// list of environments
const (
	EnvironmentDev  Environment = "DEV"
	EnvironmentProd Environment = "PROD"
)
