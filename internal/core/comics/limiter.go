package comics

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func TokenBucketRateLimit(ctx context.Context, redisClient *redis.Client, userId string, refillWindow int64, maximumTokens int64) bool {
	tokenKey := "token"
	lastRefillTimeKey := "last_refill_time"
	tokenBucket := "rate_limiting:" + userId
	tokenCountStr := redisClient.HGet(ctx, tokenBucket, tokenKey)
	lastRefillTimeStr := redisClient.HGet(ctx, tokenBucket, lastRefillTimeKey)

	tokenCount, _ := strconv.ParseInt(tokenCountStr.Val(), 10, 64)
	lastRefillTime, _ := strconv.ParseInt(lastRefillTimeStr.Val(), 10, 64)

	currentTime := time.Now().Unix()
	timeElapsed := currentTime - lastRefillTime

	if timeElapsed >= refillWindow {
		tokenCount = maximumTokens
		lastRefillTime = currentTime
	}

	if tokenCount <= 0 {
		return false
	}

	tokenCount--
	redisClient.HSet(ctx, tokenBucket, map[string]interface{}{
		tokenKey:          tokenCount,
		lastRefillTimeKey: currentTime,
	}).Val()
	return true
}
