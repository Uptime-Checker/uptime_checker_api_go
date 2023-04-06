package cache

import (
	"fmt"
	"time"

	"github.com/mdaliyan/icache/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

var userPot icache.Pot[*pkg.UserWithRoleAndSubscription]
var monitorCheckPot icache.Pot[int64]

func SetupCache() {
	userPot = icache.NewPot[*pkg.UserWithRoleAndSubscription](icache.WithTTL(7 * 24 * time.Hour))
	monitorCheckPot = icache.NewPot[int64](icache.WithTTL(8 * time.Second))
}

// SetUserWithRoleAndSubscription caches for 7 days
func SetUserWithRoleAndSubscription(user *pkg.UserWithRoleAndSubscription) {
	userPot.Set(getUserCacheKey(user.ID), user)
}

func GetUserWithRoleAndSubscription(userID int64) *pkg.UserWithRoleAndSubscription {
	user, err := userPot.Get(getUserCacheKey(userID))
	if err != nil {
		return nil
	}
	return user
}

func DeleteUserWithRoleAndSubscription(userID int64) {
	userPot.Drop(getUserCacheKey(userID))
}

func getUserCacheKey(userID int64) string {
	return fmt.Sprintf("%s_%d", "user", userID)
}

// SetMonitorToRun caches for 8 seconds
func SetMonitorToRun(monitorID, regionID int64) {
	monitorCheckPot.Set(getMonitorCacheKey(monitorID), regionID)
}

func GetMonitorToRun(monitorID int64) *int64 {
	value, err := monitorCheckPot.Get(getMonitorCacheKey(monitorID))
	if err != nil {
		return nil
	}
	return &value
}

func getMonitorCacheKey(monitorID int64) string {
	return fmt.Sprintf("monitor_to_run_%d", monitorID)
}
