package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

type SliceCmd = redis.SliceCmd
type StringStringMapCmd = redis.StringStringMapCmd

// Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	if err != nil {
		zap.L().Error("connet redis DB failed", zap.Error(err))
		return err
	}
	return nil
}

func Close() {
	_ = client.Close()
}
