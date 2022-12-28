package logic

import (
	"blog_app/dao/redis"
	"blog_app/model"
	"strconv"
)

func PostVote(userID int64, p *model.ParamVoteData) error {
	userIDstr := strconv.Itoa(int(userID))

	return redis.PostVote(userIDstr, p.PostID, float64(p.Direction))

}
