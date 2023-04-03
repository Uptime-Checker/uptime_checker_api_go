package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/coocood/freecache"
	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

var cache *freecache.Cache

const cacheSize = 10 * 1024 * 1024 // 10 MB

const (
	KeyUser = "user"
)

func SetupCache() {
	cache = freecache.NewCache(cacheSize)
}

func set(key string, value []byte, expireSeconds int) {
	if err := cache.Set([]byte(key), value, expireSeconds); err != nil {
		sentry.CaptureException(err)
	}
}

// SetUserWithRoleAndSubscription caches for 7 days
func SetUserWithRoleAndSubscription(user *pkg.UserWithRoleAndSubscription) {
	serializedUser, err := json.Marshal(user)
	if err != nil {
		sentry.CaptureException(err)
	}
	set(getUserCacheKey(user.ID), serializedUser, 7*24*60*60)
}

func GetUserWithRoleAndSubscription(userID int64) *pkg.UserWithRoleAndSubscription {
	serializedUser, err := cache.Get([]byte(getUserCacheKey(userID)))
	if err != nil {
		return nil
	}
	var user pkg.UserWithRoleAndSubscription
	if err := json.Unmarshal(serializedUser, &user); err != nil {
		sentry.CaptureException(err)
	}
	return &user
}

func DeleteUserWithRoleAndSubscription(userID int64) {
	cache.Del([]byte(getUserCacheKey(userID)))
}

func getUserCacheKey(userID int64) string {
	return fmt.Sprintf("%s_%d", KeyUser, userID)
}

// SetMonitorToRun caches for 8 seconds
func SetMonitorToRun(monitorID int64, nextCheckAt time.Time) {
	serializedNextCheckAt, err := json.Marshal(nextCheckAt)
	if err != nil {
		sentry.CaptureException(err)
	}
	set(getUserCacheKey(monitorID), serializedNextCheckAt, 8)
}

func GetMonitorToRun(monitorID int64) *time.Time {
	serializedNextCheckAt, err := cache.Get([]byte(getMonitorCacheKey(monitorID)))
	if err != nil {
		return nil
	}
	var nextCheckAt time.Time
	if err := json.Unmarshal(serializedNextCheckAt, &nextCheckAt); err != nil {
		sentry.CaptureException(err)
	}
	return &nextCheckAt
}

func getMonitorCacheKey(monitorID int64) string {
	return fmt.Sprintf("monitor_to_run_%d", monitorID)
}
