package cache

import (
	"encoding/json"

	"github.com/coocood/freecache"
	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

var cache *freecache.Cache

const cacheSize = 10 * 1024 * 1024 // 10 MB

const (
	KeyUser = "user_"
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
	set(KeyUser, serializedUser, 7*24*60*60)
}

func GetUserWithRoleAndSubscription() *pkg.UserWithRoleAndSubscription {
	serializedUser, err := cache.Get([]byte(KeyUser))
	if err != nil {
		return nil
	}
	var user pkg.UserWithRoleAndSubscription
	if err := json.Unmarshal(serializedUser, &user); err != nil {
		sentry.CaptureException(err)
	}
	return &user
}

func DeleteUserWithRoleAndSubscription() {
	cache.Del([]byte(KeyUser))
}
