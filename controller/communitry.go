package controller

import (
	"blog_app/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CommunityHandler(c *gin.Context) {
	//获取社区列表
	data, err := logic.CommunityList()
	if err != nil {
		zap.L().Error("logic.CommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//返回状态码
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	//1、获取请求参数
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2、业务处理
	data, err := logic.CommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.CommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3、返回状态码及数据
	ResponseSuccess(c, data)
}
