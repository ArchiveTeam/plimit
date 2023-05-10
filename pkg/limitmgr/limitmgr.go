package limitmgr

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

type LimitManager struct {
	rdb                  *redis.Client
	connectionsKey       string
	limitKey             string
	autoscaleLimitKey    string
	autoscaleMaxLoadKey  string
	autoscaleLowLimitKey string
}

var defaultconnectionsKey = "limiter:connections"
var defaultlimitKey = "limiter:limit"
var defaultAutoscaleLimitKey = "limiter:autoscale:hardlimit"
var defaultAutoscaleMaxLoadKey = "limiter:autoscale:maxload"
var defaultAutoscaleLowLimitKey = "limiter:autoscale:minlimit"

func NewLimitManagerFromViper() *LimitManager {
	redisConnString := viper.GetString("redis-url")
	opt, err := redis.ParseURL(redisConnString)
	if err != nil {
		log.Panicf("Failed to parse REDIS_URL: %e\n", err)
	}

	rdb := redis.NewClient(opt)

	return &LimitManager{
		rdb:                  rdb,
		connectionsKey:       defaultconnectionsKey,
		limitKey:             defaultlimitKey,
		autoscaleLimitKey:    defaultAutoscaleLimitKey,
		autoscaleMaxLoadKey:  defaultAutoscaleMaxLoadKey,
		autoscaleLowLimitKey: defaultAutoscaleLowLimitKey,
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

func (l *LimitManager) GetAutoscaleLimitKey() string {
	return l.autoscaleLimitKey
}

func (l *LimitManager) SetAutoscaleLimitKey(newKey string) {
	l.autoscaleLimitKey = newKey
}

func (l *LimitManager) GetAutoscaleMaxLoadKey() string {
	return l.autoscaleMaxLoadKey
}

func (l *LimitManager) SetAutoscaleMaxLoadKey(newKey string) {
	l.autoscaleMaxLoadKey = newKey
}

func (l *LimitManager) GetAutoscaleLowLimitKey() string {
	return l.autoscaleLowLimitKey
}

func (l *LimitManager) SetAutoscaleLowLimitKey(newKey string) {
	l.autoscaleLowLimitKey = newKey
}
