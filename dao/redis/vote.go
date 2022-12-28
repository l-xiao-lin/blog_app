package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"math"
	"strconv"
	"time"
)

const (
	OneWeekendTime = 7 * 24 * 3600
	ScorePerVote   = 432
)

var (
	ErrTimeExpire = errors.New("帖子时间已过")
	ErrRepeatVote = errors.New("重复投票")
)

func CreatePost(postID int64, communityID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//往KeyPostScoreZSet添加记录

	pipe := client.TxPipeline()
	pipe.ZAdd(ctx, GetRedisKey(KeyPostScoreZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//往KeyPostTimeZSet 添加数据

	pipe.ZAdd(ctx, GetRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//往KeyCommunitySet 添加数据
	cKey := GetRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))

	pipe.SAdd(ctx, cKey, postID)

	_, err := pipe.Exec(ctx)
	return err

}

/*
投票的几种情况,如果现在投的票值比之前来的大，op值为1，否则为-1
	direction = 1时
	之前没有投票，现在投赞成票  -->绝对差值  1 * 432
	之前投反对票，现在投造成票  -->绝对差值  2 * 432

	direction = 0 时
	之前投反对票，现在投取消票  -->绝对差值 1 * 432
	之前投赞成票，现在投取消票  -->绝对差值 1 * 432 * -1


	direction = -1时
	之前投赞成票，现在投反对票   --绝对差值 2 * 432 * -1
	之前没有投票，现在投反对票   --绝对差值 1 * 432 * -1



*/

func PostVote(userID, postID string, value float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//1、查询帖子是否过期
	postTime := client.ZScore(ctx, GetRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > OneWeekendTime {
		return ErrTimeExpire
	}
	//2、更新帖子分数
	//查询当前用户之前的投票记录
	ov := client.ZScore(ctx, GetRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	if ov == value {
		return ErrRepeatVote
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(value - ov)
	pipe := client.TxPipeline()
	pipe.ZIncrBy(ctx, GetRedisKey(KeyPostScoreZSet), op*diff*ScorePerVote, postID).Val()

	//3、更新当前用户为帖子投票的记录

	if value == 0 {
		pipe.ZRem(ctx, GetRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipe.ZAdd(ctx, GetRedisKey(KeyPostVotedZSetPF+postID), &redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipe.Exec(ctx)
	return err

}

func GetPostVote(postIds []string) (data []int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := client.Pipeline()
	for _, postId := range postIds {
		key := GetRedisKey(KeyPostVotedZSetPF + postId)
		pipe.ZCount(ctx, key, "1", "1")
	}
	data = make([]int64, 0, len(postIds))
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		return
	}
	for _, cmder := range cmders {
		value := cmder.(*redis.IntCmd).Val()
		data = append(data, value)
	}
	return

}
