package limitmgr

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var acquireScript = redis.NewScript(`
local key = KEYS[1]
local current_time = tonumber(ARGV[1])
local limit_key = ARGV[3]
local connections_key = ARGV[4]
local connection_duration = tonumber(ARGV[2])

local limit = tonumber(redis.pcall("GET", limit_key))
if limit == nil then
	limit = 0
end

local num_locks = tonumber(redis.pcall('ZCOUNT', connections_key, current_time, 'inf'))
if num_locks == nil then
	num_locks = 0
end

if num_locks < limit then
	local expire_time = current_time + connection_duration
	redis.call('ZADD', connections_key, expire_time, key)
	return true
else
	return false
end
`)

func (l *LimitManager) TryAcquireLock(ctx context.Context, id uuid.UUID, duration time.Duration) (bool, error) {
	log.Printf("Attempting to acquire lock %s...\n", id.String())
	t := time.Now().Unix()
	acquired, err := acquireScript.Run(ctx, l.rdb, []string{id.String()}, t, int(duration.Seconds()), l.limitKey, l.connectionsKey).Bool()
	if err == redis.Nil {
		return false, nil
	} else {
		return acquired, err
	}
}

func (l *LimitManager) RefreshLock(ctx context.Context, id uuid.UUID, duration time.Duration) {
	log.Printf("Refreshing lock %s...\n", id.String())
	t := time.Now().Add(duration).Unix()
	err := l.rdb.ZAddXX(ctx, l.connectionsKey, redis.Z{Score: float64(t), Member: id.String()}).Err()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to renew lock: %s\n", err)
	}
}

func (l *LimitManager) ReleaseLock(ctx context.Context, id uuid.UUID) {
	log.Printf("Attempting to release lock %s...\n", id.String())
	err := l.rdb.ZRem(ctx, l.connectionsKey, id.String()).Err()
	if err != nil {
		log.Panicf("Failed to release lock %s: %s\n", id.String(), err)
	}
	log.Printf("Released lock %s.\n", id.String())
}
