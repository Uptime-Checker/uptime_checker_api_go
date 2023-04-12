package gandalf

import (
	"fmt"

	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

func CanCreateMonitor(user *pkg.UserWithRoleAndSubscription, count, interval int32) error {
	if err := CanCreate(user); err != nil {
		return err
	}
	if err := HandleFeatureMax(user, FeatureAPICheckCount, count); err != nil {
		return fmt.Errorf("%w - %s", err, FeatureAPICheckCount)
	}
	if err := HandleFeatureMin(user, FeatureAPICheckInterval, interval); err != nil {
		return fmt.Errorf("%w - %s", err, FeatureAPICheckInterval)
	}
	return nil
}
