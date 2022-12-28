package redis

import (
	"blog_app/model"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func GetPostListIDs(p *model.ParamPostList) ([]string, error) {
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	key := GetRedisKey(KeyPostTimeZSet)
	if p.Order == "score" {
		key = GetRedisKey(KeyPostScoreZSet)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.ZRevRange(ctx, key, start, end).Result()

}

func GetPostListIDsByCommunity(p *model.ParamPostList) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	//判断是否已经有key
	orderKey := GetRedisKey(KeyPostTimeZSet)
	if p.Order == model.OrderScore {
		orderKey = GetRedisKey(KeyPostScoreZSet)
	}

	cKey := GetRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	key := orderKey + strconv.Itoa(int(p.CommunityID))

	if client.Exists(ctx, key).Val() < 1 {
		pipe := client.Pipeline()
		pipe.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{orderKey, cKey},
			Aggregate: "MAX",
		})
		pipe.Expire(ctx, key, 3600*time.Second)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	return client.ZRevRange(ctx, key, start, end).Result()

}
