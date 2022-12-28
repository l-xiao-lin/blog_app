package mysql

import (
	"blog_app/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CreatePost(p *model.Post) error {
	sqlStr := "insert into post(post_id,title,content,author_id,community_id)values (?,?,?,?,?)"
	db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorId, p.CommunityId)

	return nil
}

func PostList(page, size int64) (postList []*model.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post order by create_time DESC limit ?,? `

	//不用make分配内存也是可以的
	postList = make([]*model.Post, 0, 2)
	err = db.Select(&postList, sqlStr, (page-1)*size, size)
	return

}

func PostDetailByID(pid int64) (post *model.Post, err error) {
	post = new(model.Post)
	sqlStr := `select post_id,title,content,author_id,community_id from post where post_id=?`
	err = db.Get(post, sqlStr, pid)
	return
}

func PostListInOrder(ids []string) (postList []*model.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time  from post where post_id in (?) order by find_in_set(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return

}
