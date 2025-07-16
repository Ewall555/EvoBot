package redis

import (
	"EvoBot/backend/global"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func Init() {
	RDS := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        []string{fmt.Sprintf("%s:%d", global.CONF.RedisConfig.Host, global.CONF.RedisConfig.Port)},
		Password:     global.CONF.RedisConfig.Password,
		DB:           global.CONF.RedisConfig.Database,
		PoolSize:     global.CONF.RedisConfig.PoolSize,
		MinIdleConns: global.CONF.RedisConfig.MinIdleConns,
		IdleTimeout:  time.Second * time.Duration(global.CONF.RedisConfig.IdleTimeout),
	})
	if _, err := RDS.Ping(context.Background()).Result(); err != nil {
		panic(fmt.Errorf("connect redis, err: %s ", err))
	}
	global.RDS = RDS
	global.ZAPLOG.Info("redis init success")
}
