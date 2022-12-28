package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var CxtUserIDKey = "user_id"

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CxtUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, _ = uid.(int64)
	return

}

func GetPageInfo(c *gin.Context) (page, size int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return

}
