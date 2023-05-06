package limitmgr

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func (l *LimitManager) GetLimit(ctx context.Context) int {
	limit, err := l.rdb.Get(ctx, l.limitKey).Int()
	if err != nil {
		if err == redis.Nil {
			limit = 0
		} else {
			log.Fatalf("Failed to fetch data: %e\n", err)
		}
	}

	return limit
}

func (l *LimitManager) SetLimit(ctx context.Context, newLimit int64) {
	log.Printf("Updating limit to %v.\n", newLimit)
	err := l.rdb.Set(ctx, l.limitKey, newLimit, 0).Err()
	if err != nil {
		log.Fatalf("Failed to set limit: %e\n", err)
	}
}

func (l *LimitManager) GetCurrentConnectionCount(ctx context.Context) int64 {
	t := time.Now().Unix()
	count, err := l.rdb.ZCount(ctx, l.connectionsKey, fmt.Sprintf("%v", t), "inf").Result()
	if err != nil {
		log.Fatalf("Unable to count connections: %s\n", err)
	}
	return count
}
