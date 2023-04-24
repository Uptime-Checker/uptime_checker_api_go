package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stripe/stripe-go/v74"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

// BearerClaims contains claims
type BearerClaims struct {
	UserID int64
	Email  string
	*jwt.RegisteredClaims
}

type ShadowAssertion struct {
	Source     int32
	Property   *string
	Comparison int32
	Value      string
}

type SubscriptionFeature struct {
	*model.ProductFeature
	*model.Feature
}

type UserWithRoleAndSubscription struct {
	*model.User
	Organization *model.Organization

	Role struct {
		*model.Role
		Claims []*model.RoleClaim
	}

	Subscription struct {
		*model.Subscription
		Plan     *model.Plan
		Product  *model.Product
		Features []*SubscriptionFeature
	}
}

type MonitorRegionWithAssertions struct {
	*model.MonitorRegion
	Region *model.Region

	Monitor struct {
		*model.Monitor
		Assertions []model.Assertion
	}
}

type OrganizationUserRole struct {
	*model.OrganizationUser
	Role         *model.Role
	Organization *model.Organization
}

type PlanWithProduct struct {
	*model.Plan
	*model.Product
}

type ProductWithPlansAndFeatures struct {
	*model.Product
	Plans    []*model.Plan
	Features []*SubscriptionFeature
}

type BillingProduct struct {
	*stripe.Product
	Prices []*stripe.Price
}

type WebhookData struct {
	EventType, EventID string
	EventAt            time.Time
	Data               map[string]any
}
