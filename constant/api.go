package constant

const (
	GuestUserRateLimitInMinutes = 5
)

// DateCompare type
type DateCompare int

// list of DateCompares
const (
	Date1AfterDate2  DateCompare = 1
	Date1BeforeDate2 DateCompare = -1
	Date1EqualDate2  DateCompare = 0
)
