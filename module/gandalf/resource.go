package gandalf

import (
	"github.com/samber/lo"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

// ClaimType type
type ClaimType string

const (
	ClaimCreateResource      ClaimType = "CREATE_RESOURCE"
	ClaimUpdateResource      ClaimType = "UPDATE_RESOURCE"
	ClaimDeleteResource      ClaimType = "DELETE_RESOURCE"
	ClaimBilling             ClaimType = "BILLING"
	ClaimInviteUser          ClaimType = "INVITE_USER"
	ClaimDestroyOrganization ClaimType = "DESTROY_ORGANIZATION"
)

func CanCreate(user *pkg.UserWithRoleAndSubscription) error {
	return handleClaim(user, ClaimCreateResource)
}

func CanUpdate(user *pkg.UserWithRoleAndSubscription) error {
	return handleClaim(user, ClaimUpdateResource)
}

func CanDelete(user *pkg.UserWithRoleAndSubscription) error {
	return handleClaim(user, ClaimDeleteResource)
}

func CanClaimBilling(user *pkg.UserWithRoleAndSubscription) error {
	return handleClaim(user, ClaimBilling)
}

func CanClaimInviteUser(user *pkg.UserWithRoleAndSubscription) error {
	return handleClaim(user, ClaimInviteUser)
}

func CanClaimDestroyOrganization(user *pkg.UserWithRoleAndSubscription) error {
	return handleClaim(user, ClaimDestroyOrganization)
}

func handleClaim(user *pkg.UserWithRoleAndSubscription, requiredClaim ClaimType) error {
	_, found := lo.Find(user.Role.Claims, func(item *model.RoleClaim) bool {
		return item.Claim == string(requiredClaim)
	})
	if !found {
		return constant.ErrUpgradePermission
	}
	return nil
}
