package times

import "time"

// CompareDate compares date1 and date2. If date1 is before date2 it returns -1,
// if date1 is after date2 it returns 1 otherwise 0
func CompareDate(date1, date2 time.Time) int {
	dt1 := date1.UTC()
	dt2 := date2.UTC()
	y1, m1, d1 := dt1.Date()
	y2, m2, d2 := dt2.Date()
	h1, n1, s1 := dt1.Clock()
	h2, n2, s2 := dt2.Clock()
	t1 := time.Date(y1, m1, d1, h1, n1, s1, dt1.Nanosecond(), time.UTC)
	t2 := time.Date(y2, m2, d2, h2, n2, s2, dt2.Nanosecond(), time.UTC)
	if t1.Before(t2) {
		return -1
	}
	if t1.After(t2) {
		return 1
	}
	return 0
}

// Now returns in utc
func Now() time.Time {
	return time.Now().UTC()
}

// IsSameDay returns if two dates are same day
func IsSameDay(date1, date2 time.Time) bool {
	if date1.Year() == date2.Year() && date1.Month() == date2.Month() && date1.Day() == date2.Day() {
		return true
	}
	return false
}

// BeginningOfTheDay returns beginning from current time
func BeginningOfTheDay(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

// RouteDate returns time
func RouteDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func DaysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	return days
}

func ParseAdaptyDate(unparsedTime string) (*time.Time, error) {
	layout := "2006-01-02T15:04:05-0700"
	parsedTime, err := time.Parse(layout, unparsedTime)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}
