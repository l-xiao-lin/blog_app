package mysql

import (
	"blog_app/model"
	"database/sql"
)

func CommunityList() (communityList []*model.Community, err error) {
	//返回值如果申明的是切片类型，在申明时，已经重新分配内存
	sqlStr := `select community_id,community_name from community`
	err = db.Select(&communityList, sqlStr)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func GetCommunityDetailById(id int64) (community *model.CommunityDetail, err error) {
	community = new(model.CommunityDetail)
	sqlStr := "select community_id,community_name from community where community_id=?"
	err = db.Get(community, sqlStr, id)
	return
}
