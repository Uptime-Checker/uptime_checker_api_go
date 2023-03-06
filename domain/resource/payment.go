package resource

// SubscriptionStatus type
type SubscriptionStatus int

// list of subscription statuses
const (
	SubscriptionStatusIncomplete SubscriptionStatus = iota + 1
	SubscriptionStatusIncompleteExpired
	SubscriptionStatusTrialing
	SubscriptionStatusActive
	SubscriptionStatusPastDue
	SubscriptionStatusCanceled
	SubscriptionStatusUnpaid
)

// Valid checks if the SubscriptionStatus is valid
func (s SubscriptionStatus) Valid() bool {
	subscriptionStatuses := []SubscriptionStatus{
		SubscriptionStatusIncomplete,
		SubscriptionStatusIncompleteExpired,
		SubscriptionStatusTrialing,
		SubscriptionStatusActive,
		SubscriptionStatusPastDue,
		SubscriptionStatusCanceled,
		SubscriptionStatusUnpaid,
	}
	for _, c := range subscriptionStatuses {
		if c == s {
			return true
		}
	}
	return false
}

func (s SubscriptionStatus) String() string {
	return [...]string{"incomplete", "incompleteExpired", "trialing", "active", "pastDue", "canceled", "unpaid"}[s-1]
}
