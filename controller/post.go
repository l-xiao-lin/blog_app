package controller

import (
	"blog_app/logic"
	"blog_app/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
)

func CreatePostHandler(c *gin.Context) {
	//1、参数校验
	p := new(model.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}

	//2、从context中获取当前用户的ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorId = userID

	//3、创建post
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//4、返回状态码
	ResponseSuccess(c, nil)
}

func PostListHandler(c *gin.Context) {
	//1、获取分页参数
	page, size := GetPageInfo(c)

	//2、查询数据库
	data, err := logic.PostList(page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	//3、返回状态码
	ResponseSuccess(c, data)
}

func PostDetailHandler(c *gin.Context) {
	//1、获取url中的post_id参数
	strPid := c.Param("id")
	pid, err := strconv.ParseInt(strPid, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2、查询数据库
	data, err := logic.PostDetailByID(pid)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	//3、返回
	ResponseSuccess(c, data)
}

func PostListInOrder(c *gin.Context) {
	//1、url参数校验
	p := &model.ParamPostList{
		Page:  1,
		Size:  5,
		Order: model.OrderTime,
	}
	if err := c.ShouldBind(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2、调用业务逻辑
	data, err := logic.PostListInOrder(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	//3、返回状态码
	ResponseSuccess(c, data)

}

func PostCommunityListInOrder(c *gin.Context) {
	//1、参数校验
	p := &model.ParamPostList{
		Page:  1,
		Size:  5,
		Order: model.OrderTime,
	}
	if err := c.ShouldBind(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2、业务逻辑
	data, err := logic.GetPostListByCommunityIdsInOrder(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	//3、返回状态码
	ResponseSuccess(c, data)
}
