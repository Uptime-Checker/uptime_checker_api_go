package constant

const (
	GuestUserRateLimitInMinutes     = 5
	GuestUserCodeExpiryInMinutes    = 10
	BearerTokenExpirationInDays     = 180
	FreeSubscriptionDurationInDays  = 120
	TrialSubscriptionDurationInDays = 14
)

// DateCompare type
type DateCompare int

// list of DateCompares
const (
	Date1AfterDate2  DateCompare = 1
	Date1BeforeDate2 DateCompare = -1
	Date1EqualDate2  DateCompare = 0
)
