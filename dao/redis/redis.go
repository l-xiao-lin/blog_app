package redis

import (
	"blog_app/settings"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var client *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "",
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.Ping(ctx).Result()
	return err

}

func Close() {
	_ = client.Close()
}
