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

// ProductTier type
type ProductTier int

// list of product tiers
const (
	ProductTierFree ProductTier = iota + 1
	ProductTierDeveloper
	ProductTierStartup
	ProductTierEnterprise
)

// Valid checks if the ProductTier is valid
func (p ProductTier) Valid() bool {
	productTiers := []ProductTier{
		ProductTierFree,
		ProductTierDeveloper,
		ProductTierStartup,
		ProductTierEnterprise,
	}
	for _, c := range productTiers {
		if c == p {
			return true
		}
	}
	return false
}

func (p ProductTier) String() string {
	return [...]string{"Free", "Developer", "Startup", "Enterprise"}[p-1]
}
