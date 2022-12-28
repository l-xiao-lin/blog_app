package logic

import (
	"blog_app/dao/mysql"
	"blog_app/dao/redis"
	"blog_app/model"
	"blog_app/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *model.Post) (err error) {
	//1、生成post_id
	p.ID = snowflake.GenID()

	//2、写到表中
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}

	//3、往redis中的post:time post:score ZSet添加数据
	err = redis.CreatePost(p.ID, p.CommunityId)
	return
}

func PostList(page, size int64) (data []*model.ApiPostDetail, err error) {
	postList, err := mysql.PostList(page, size)

	data = make([]*model.ApiPostDetail, 0, len(postList))
	for _, post := range postList {
		//1、通过author_id查询用户
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Int64("author_id", post.AuthorId), zap.Error(err))

			continue
		}
		//2、通过community_id查询社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityId)", zap.Int64("community_id", post.CommunityId), zap.Error(err))
			continue
		}

		postDetail := &model.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return

}

func PostDetailByID(pid int64) (data *model.ApiPostDetail, err error) {
	//1、通过pid查询数据库
	post, err := mysql.PostDetailByID(pid)

	//2、通过author_id查询作者信息
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		return
	}
	//3、通过community_id查询社区信息
	community, err := mysql.GetCommunityDetailById(post.CommunityId)
	if err != nil {
		return
	}
	data = &model.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return

}

func PostListInOrder(p *model.ParamPostList) (data []*model.ApiPostDetail, err error) {
	//1、从redis中查找指定的post_id列表
	postIDs, err := redis.GetPostListIDs(p)
	if err != nil {
		return
	}
	if len(postIDs) == 0 {
		return
	}

	//2、将上面返回的post_id作为参数，往数据库中查询
	postList, err := mysql.PostListInOrder(postIDs)
	if err != nil {
		return
	}

	data = make([]*model.ApiPostDetail, 0, len(postList))

	//3、查询每篇帖子的投票数
	voteData, err := redis.GetPostVote(postIDs)
	if err != nil {
		return
	}
	for index, post := range postList {

		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			continue
		}

		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {

			continue
		}

		postDetail := &model.ApiPostDetail{
			VoteNum:         voteData[index],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)

	}
	return
}

func GetPostListByCommunityIdsInOrder(p *model.ParamPostList) (data []*model.ApiPostDetail, err error) {
	//1、通过community参数获取postIDS
	postIds, err := redis.GetPostListIDsByCommunity(p)
	if err != nil {
		return
	}
	if len(postIds) == 0 {
		return
	}

	//2、将上面获取到的postIds查询数据库
	postList, err := mysql.PostListInOrder(postIds)
	if err != nil {
		return
	}

	//3、查询每个帖子的投票数
	voteData, err := redis.GetPostVote(postIds)
	if err != nil {
		return
	}

	//4、补充帖子的其他信息
	for index, post := range postList {
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			continue
		}
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			continue
		}
		postDetail := &model.ApiPostDetail{
			VoteNum:         voteData[index],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)

	}
	return
}
