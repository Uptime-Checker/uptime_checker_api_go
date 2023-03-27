package cache

import (
	"encoding/json"
	"fmt"

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
