package routers

import (
	"blog_app/controller"
	"blog_app/logger"
	"blog_app/middleware"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	//加载logger日志中间件
	r.Use(ginzap.Ginzap(logger.Lg, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.Lg, true))

	v1 := r.Group("/api/v1")

	//注册
	v1.POST("/signup", controller.SingUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)

	//获取社区列表
	v1.GET("/community", controller.CommunityHandler)
	//获取社区详情
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	//获取帖子列表
	v1.GET("/post", controller.PostListHandler)
	//按时间或分数获取帖子列表
	v1.GET("/post_list_in_order", controller.PostListInOrder)
	//选中指定社区，同时按时间或分数获取帖子列表
	v1.GET("/post_community_list_in_order", controller.PostCommunityListInOrder)

	//获取帖子详情
	v1.GET("/post/:id", controller.PostDetailHandler)

	v1.Use(middleware.JwtAuthMiddleware())
	{
		v1.POST("/post", controller.CreatePostHandler)
		v1.POST("/vote", controller.PostVoteHandler)

	}

	return r
}
