package limitmgr

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

type LimitManager struct {
	rdb            *redis.Client
	connectionsKey string
	limitKey       string
}

var defaultconnectionsKey = "limiter:connections"
var defaultlimitKey = "limiter:limit"

func NewLimitManagerFromViper() *LimitManager {
	redisConnString := viper.GetString("redis_url")
	opt, err := redis.ParseURL(redisConnString)
	if err != nil {
		log.Fatalf("Failed to parse REDIS_URL: %e\n", err)
	}

	rdb := redis.NewClient(opt)

	return &LimitManager{
		rdb:            rdb,
		connectionsKey: defaultconnectionsKey,
		limitKey:       defaultlimitKey,
	}
}

func (l *LimitManager) GetLimitKey() string {
	return l.limitKey
}

func (l *LimitManager) SetLimitKey(newKey string) {
	l.limitKey = newKey
}

func (l *LimitManager) GetConnectionsKey() string {
	return l.connectionsKey
}

func (l *LimitManager) SetConnectionsKey(newKey string) {
	l.connectionsKey = newKey
}
