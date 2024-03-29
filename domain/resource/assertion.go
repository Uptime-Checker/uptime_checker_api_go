package resource

// AssertionComparison type
type AssertionComparison int

// List of assertion comparisons
const (
	AssertionComparisonEqual AssertionComparison = iota + 1
	AssertionComparisonNotEqual
	AssertionComparisonGreater
	AssertionComparisonLesser
	AssertionComparisonEmpty
	AssertionComparisonNotEmpty
	AssertionComparisonContain
	AssertionComparisonNotContain
	AssertionComparisonNull
	AssertionComparisonNotNull
	AssertionComparisonHasKey
	AssertionComparisonNotHasKey
	AssertionComparisonHasValue
	AssertionComparisonNotHasValue
)

// Valid checks if the AssertionComparison is valid
func (a AssertionComparison) Valid() bool {
	assertionComparisons := []AssertionComparison{
		AssertionComparisonEqual,
		AssertionComparisonNotEqual,
		AssertionComparisonGreater,
		AssertionComparisonLesser,
		AssertionComparisonEmpty,
		AssertionComparisonNotEmpty,
		AssertionComparisonContain,
		AssertionComparisonNotContain,
		AssertionComparisonNull,
		AssertionComparisonNotNull,
		AssertionComparisonHasKey,
		AssertionComparisonNotHasKey,
		AssertionComparisonHasValue,
		AssertionComparisonNotHasValue,
	}
	for _, p := range assertionComparisons {
		if p == a {
			return true
		}
	}
	return false
}

func (a AssertionComparison) String() string {
	return [...]string{
		"equal", "not-equal", "greater-than", "lesser-than", "empty", "not-empty", "contains", "not-contains",
		"null", "not-null", "has-key", "not-has-key", "has-value", "not-has-value",
	}[a-1]
}

// AssertionSource type
type AssertionSource int

// List of assertion sources
const (
	AssertionSourceStatusCode AssertionSource = iota + 1
	AssertionSourceTextBody
	AssertionSourceHeaders
	AssertionSourceResponseTime
	AssertionSourceJSONBody
)

// Valid checks if the AssertionSource is valid
func (a AssertionSource) Valid() bool {
	assertionSources := []AssertionSource{
		AssertionSourceStatusCode,
		AssertionSourceTextBody,
		AssertionSourceHeaders,
		AssertionSourceResponseTime,
		AssertionSourceJSONBody,
	}
	for _, p := range assertionSources {
		if p == a {
			return true
		}
	}
	return false
}

func (a AssertionSource) String() string {
	return [...]string{"status-code", "text-body", "headers", "response-time", "json-body"}[a-1]
}
