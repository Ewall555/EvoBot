package utils

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func AcquireRedisLockWithRetry(ctx context.Context, rds redis.UniversalClient, key string, ttl time.Duration, retryInterval time.Duration, maxWait time.Duration) (lockValue string, ok bool) {
	lockValue = uuid.NewString()
	start := time.Now()

	for {
		ok = rds.SetNX(ctx, key, lockValue, ttl).Val()
		if ok {
			return lockValue, true
		}
		if time.Since(start) > maxWait {
			return "", false
		}
		time.Sleep(retryInterval)
	}
}

func ReleaseRedisLock(ctx context.Context, rds redis.UniversalClient, key string, lockValue string) {
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	rds.Eval(ctx, luaScript, []string{key}, lockValue)
}
