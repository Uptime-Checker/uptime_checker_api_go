package resource

// SubscriptionStatus type
type SubscriptionStatus int

// List of subscription statuses
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

// List of product tiers
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

func GetProductTier(tier string) ProductTier {
	switch tier {
	case "Free":
		return ProductTierFree
	case "Developer":
		return ProductTierDeveloper
	case "Startup":
		return ProductTierStartup
	case "Enterprise":
		return ProductTierEnterprise
	}
	return ProductTierFree
}

// PlanType type
type PlanType int

// List of plan types
const (
	PlanTypeMonthly PlanType = iota + 1
	PlanTypeYearly
)

// Valid checks if the PlanType is valid
func (p PlanType) Valid() bool {
	planTypes := []PlanType{
		PlanTypeMonthly,
		PlanTypeYearly,
	}
	for _, c := range planTypes {
		if c == p {
			return true
		}
	}
	return false
}

func (p PlanType) String() string {
	return [...]string{"Monthly", "Yearly"}[p-1]
}
