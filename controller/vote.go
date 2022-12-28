package controller

import (
	"blog_app/logic"
	"blog_app/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PostVoteHandler(c *gin.Context) {
	//1、参数校验
	p := new(model.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}
	//2、获取当前用户userID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	//3、调用投票功能
	if err = logic.PostVote(userID, p); err != nil {

		ResponseError(c, CodeServerBusy)
		return
	}

	//4、返回状态
	ResponseSuccess(c, nil)
}
