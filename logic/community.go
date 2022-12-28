package logic

import (
	"blog_app/dao/mysql"
	"blog_app/model"
)

func CommunityList() ([]*model.Community, error) {
	//查询数据库 查找所有的community数据，并返回
	return mysql.CommunityList()
}

func CommunityDetail(id int64) (*model.CommunityDetail, error) {

	return mysql.GetCommunityDetailById(id)
}
