package resource

// AssertionComparison type
type AssertionComparison int

// list of assertion comparisons
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

// Valid checks if the UserLoginProvider is valid
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
