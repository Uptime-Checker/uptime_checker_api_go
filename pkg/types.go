package pkg

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

// BearerClaims contains claims
type BearerClaims struct {
	UserID int64
	Email  string
	*jwt.RegisteredClaims
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
		Features []*model.Feature
	}
}
