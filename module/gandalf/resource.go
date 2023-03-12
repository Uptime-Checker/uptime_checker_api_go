package gandalf

import (
	. "github.com/samber/lo"

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

func CanCreate(user *pkg.UserWithRoleAndSubscription) bool {
	return handle(user, ClaimCreateResource)
}

func CanUpdate(user *pkg.UserWithRoleAndSubscription) bool {
	return handle(user, ClaimUpdateResource)
}

func CanDelete(user *pkg.UserWithRoleAndSubscription) bool {
	return handle(user, ClaimDeleteResource)
}

func CanClaimBilling(user *pkg.UserWithRoleAndSubscription) bool {
	return handle(user, ClaimBilling)
}

func CanClaimInviteUser(user *pkg.UserWithRoleAndSubscription) bool {
	return handle(user, ClaimInviteUser)
}

func CanClaimDestroyOrganization(user *pkg.UserWithRoleAndSubscription) bool {
	return handle(user, ClaimDestroyOrganization)
}

func handle(user *pkg.UserWithRoleAndSubscription, requiredClaim ClaimType) bool {
	_, found := Find(user.Role.Claims, func(item *model.RoleClaim) bool {
		return item.Claim == string(requiredClaim)
	})
	return found
}
