package limitmgr

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"modernc.org/mathutil"
	"time"
)

func (l *LimitManager) GetLimit(ctx context.Context) int {
	limit, err := l.rdb.Get(ctx, l.limitKey).Int()
	if err != nil {
		if err == redis.Nil {
			limit = 0
		} else {
			log.Panicf("Failed to fetch data: %e\n", err)
		}
	}

	return limit
}

func (l *LimitManager) SetLimit(ctx context.Context, newLimit int64) {
	log.Printf("Updating limit to %v.\n", newLimit)
	err := l.rdb.Set(ctx, l.limitKey, newLimit, 0).Err()
	if err != nil {
		log.Panicf("Failed to set limit: %e\n", err)
	}
}

func (l *LimitManager) GetCurrentConnectionCount(ctx context.Context) int64 {
	t := time.Now().Unix()
	count, err := l.rdb.ZCount(ctx, l.connectionsKey, fmt.Sprintf("%v", t), "inf").Result()
	if err != nil {
		log.Panicf("Unable to count connections: %s\n", err)
	}
	return count
}

func (l *LimitManager) GetAutoscaleHardLimit(ctx context.Context) int {
	limit, err := l.rdb.Get(ctx, l.autoscaleLimitKey).Int()
	if err != nil {
		if err == redis.Nil {
			limit = 0
		} else {
			log.Panicf("Failed to fetch data: %e\n", err)
		}
	}

	return limit
}

func (l *LimitManager) SetAutoscaleHardLimit(ctx context.Context, newLimit int64) {
	log.Printf("Updating autoscale hard limit to %v.\n", newLimit)
	err := l.rdb.Set(ctx, l.autoscaleLimitKey, newLimit, 0).Err()
	if err != nil {
		log.Panicf("Failed to set autoscale hard limit: %e\n", err)
	}
}

func (l *LimitManager) GetAutoscaleMaxLoad(ctx context.Context) int {
	limit, err := l.rdb.Get(ctx, l.autoscaleMaxLoadKey).Int()
	if err != nil {
		if err == redis.Nil {
			limit = 0
		} else {
			log.Panicf("Failed to fetch data: %e\n", err)
		}
	}

	return mathutil.Clamp(limit, 0, 100)
}

func (l *LimitManager) SetAutoscaleMaxLoad(ctx context.Context, newLimit int) {
	log.Printf("Updating autoscale max load to %v.\n", newLimit)
	err := l.rdb.Set(ctx, l.autoscaleMaxLoadKey, newLimit, 0).Err()
	if err != nil {
		log.Panicf("Failed to set autoscale max load: %e\n", err)
	}
}

func (l *LimitManager) GetAutoscaleMinLimit(ctx context.Context) int {
	limit, err := l.rdb.Get(ctx, l.autoscaleLowLimitKey).Int()
	if err != nil {
		if err == redis.Nil {
			limit = 0
		} else {
			log.Panicf("Failed to fetch data: %e\n", err)
		}
	}

	return limit
}

func (l *LimitManager) SetAutoscaleMinLimit(ctx context.Context, newLimit int64) {
	log.Printf("Updating autoscale min limit to %v.\n", newLimit)
	err := l.rdb.Set(ctx, l.autoscaleLowLimitKey, newLimit, 0).Err()
	if err != nil {
		log.Panicf("Failed to set autoscale min limit: %e\n", err)
	}
}
