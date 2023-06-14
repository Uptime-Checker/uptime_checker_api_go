package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

var (
	remoteCache *cache.Cache
	rdb         *redis.Client
)

func SetupRemoteCache() {
	opt, err := redis.ParseURL(config.App.RedisCache)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opt)
	remoteCache = cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, 10*time.Minute),
	})
}

func Shutdown() error {
	return rdb.Close()
}

// SetUserWithRoleAndSubscription caches for 7 days
func SetUserWithRoleAndSubscription(ctx context.Context, user *pkg.UserWithRoleAndSubscription) {
	serializedUser, err := json.Marshal(user)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	if err := remoteCache.Set(&cache.Item{
		Ctx: ctx, Key: getUserCacheKey(user.ID), Value: serializedUser, TTL: 7 * 24 * time.Hour,
	}); err != nil {
		sentry.CaptureException(err)
	}
}

func GetUserWithRoleAndSubscription(ctx context.Context, userID int64) *pkg.UserWithRoleAndSubscription {
	var serializedUser string
	if err := remoteCache.Get(ctx, getUserCacheKey(userID), &serializedUser); err != nil {
		return nil
	}
	var user pkg.UserWithRoleAndSubscription
	if err := json.Unmarshal([]byte(serializedUser), &user); err != nil {
		sentry.CaptureException(err)
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

// ========================================================================

// SetReceiptEventForCustomer caches for 1 hour
func SetPaymentEventForCustomer(ctx context.Context, key string, eventAt time.Time) {
	serializedTime, err := json.Marshal(eventAt)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	if err := remoteCache.Set(&cache.Item{
		Ctx: ctx, Key: key, Value: serializedTime, TTL: 1 * time.Hour,
	}); err != nil {
		sentry.CaptureException(err)
	}
}

func GetPaymentEventForCustomer(ctx context.Context, key string) *time.Time {
	var serializedTime string
	if err := remoteCache.Get(ctx, key, &serializedTime); err != nil {
		return nil
	}
	var eventAt time.Time
	if err := json.Unmarshal([]byte(serializedTime), &eventAt); err != nil {
		sentry.CaptureException(err)
	}
	return &eventAt
}

func GetReceiptEventKey(customerID string) string {
	return fmt.Sprintf("receipt_event_%s", customerID)
}

func GetSubscriptionEventKey(customerID string) string {
	return fmt.Sprintf("subscription_event_%s", customerID)
}

// ========================================================================
// SetInternalProducts caches for 7 days
func SetInternalProducts(ctx context.Context, products []pkg.ProductWithPlansAndFeatures) {
	serializedProducts, err := json.Marshal(products)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	if err := remoteCache.Set(&cache.Item{
		Ctx: ctx, Key: getProductsCacheKey(), Value: serializedProducts, TTL: 7 * 24 * time.Hour,
	}); err != nil {
		sentry.CaptureException(err)
	}
}

func GetInternalProducts(ctx context.Context) *[]pkg.ProductWithPlansAndFeatures {
	var serializedProducts string
	if err := remoteCache.Get(ctx, getProductsCacheKey(), &serializedProducts); err != nil {
		return nil
	}
	var products []pkg.ProductWithPlansAndFeatures
	if err := json.Unmarshal([]byte(serializedProducts), &products); err != nil {
		sentry.CaptureException(err)
	}
	return &products
}

func DeleteInternalProducts(ctx context.Context) {
	if err := remoteCache.Delete(ctx, getProductsCacheKey()); err != nil {
		sentry.CaptureException(err)
	}
}

func getProductsCacheKey() string {
	return "products_internal"
}
