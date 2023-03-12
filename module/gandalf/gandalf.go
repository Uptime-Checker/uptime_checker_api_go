package gandalf

import "github.com/Uptime-Checker/uptime_checker_api_go/pkg"

func CanCreateMonitor(user *pkg.UserWithRoleAndSubscription, count, interval int32) error {
	if err := CanCreate(user); err != nil {
		return err
	}
	if err := HandleFeatureMax(user, FeatureAPICheckCount, count); err != nil {
		return err
	}
	if err := HandleFeatureMin(user, FeatureAPICheckInterval, interval); err != nil {
		return err
	}
	return nil
}
