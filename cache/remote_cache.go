package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

var remoteCache *cache.Cache

func SetupRemoteCache() {
	opt, err := redis.ParseURL(config.App.RedisCache)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	remoteCache = cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, 10*time.Minute),
	})
}

// SetUserWithRoleAndSubscription caches for 7 days
func SetUserWithRoleAndSubscription(ctx context.Context, user *pkg.UserWithRoleAndSubscription) {
	if err := remoteCache.Set(&cache.Item{
		Ctx: ctx, Key: getUserCacheKey(user.ID), Value: user, TTL: 7 * 24 * time.Hour,
	}); err != nil {
		sentry.CaptureException(err)
	}
}

func GetUserWithRoleAndSubscription(ctx context.Context, userID int64) *pkg.UserWithRoleAndSubscription {
	var user pkg.UserWithRoleAndSubscription
	if err := remoteCache.Get(ctx, getUserCacheKey(userID), &user); err != nil {
		return nil
	}
	return &user
}

func DeleteUserWithRoleAndSubscription(ctx context.Context, userID int64) {
	if err := remoteCache.Delete(ctx, getUserCacheKey(userID)); err != nil {
		sentry.CaptureException(err)
	}
}

func getUserCacheKey(userID int64) string {
	return fmt.Sprintf("%s_%d", "user", userID)
}
