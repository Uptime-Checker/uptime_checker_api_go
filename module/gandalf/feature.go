package gandalf

import (
	. "github.com/samber/lo"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
)

// FeatureType type
type FeatureType string

const (
	FeatureAPICheckCount    FeatureType = "API_CHECK_COUNT"
	FeatureAPICheckInterval FeatureType = "API_CHECK_INTERVAL"
	FeatureUserCount        FeatureType = "USER_COUNT"
)

func HandleFeatureMax(user *pkg.UserWithRoleAndSubscription, feature FeatureType, count int32) error {
	subscriptionFeature, err := handleFeature(user, feature)
	if err != nil {
		return err
	}
	if count > *subscriptionFeature.Count {
		return constant.ErrUpgradeSubscription
	}
	return nil
}

func HandleFeatureMin(user *pkg.UserWithRoleAndSubscription, feature FeatureType, count int32) error {
	subscriptionFeature, err := handleFeature(user, feature)
	if err != nil {
		return err
	}
	if count < *subscriptionFeature.Count {
		return constant.ErrUpgradeSubscription
	}
	return nil
}

func handleFeature(user *pkg.UserWithRoleAndSubscription, feature FeatureType) (*pkg.SubscriptionFeature, error) {
	now := times.Now()
	if times.CompareDate(now, *user.Subscription.ExpiresAt) == constant.Date1AfterDate2 { // Expired
		return nil, constant.ErrSubscriptionExpired
	}

	subscriptionFeature, found := Find(user.Subscription.Features, func(item *pkg.SubscriptionFeature) bool {
		return item.Feature.Name == string(feature)
	})
	if !found {
		return nil, constant.ErrUpgradeSubscription
	}
	return subscriptionFeature, nil
}
