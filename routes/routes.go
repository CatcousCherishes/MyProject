package routes

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"web_app/controller"
	_ "web_app/docs" // 千万不要忘了导入把你上一步生成的docs
	"web_app/logger"
	"web_app/middlewares"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(mode string) *gin.Engine {
	//gin 设置为发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	//  zap日志库的两个中间件的使用
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	//user 注册业务 路由
	v1.POST("/SingUp", controller.SignUpHander)
	//user 登录
	v1.POST("/login", controller.LoginHandler)
	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件

	//功能需求，先从router开始，按照请求处理流程开始
	{ //查询社区信息(社区是由后端数据库设定的，不用创建，只要在数据库中添加即可)
		v1.GET("/community", controller.CommunityHandler)
		//根据id查询到社区详情
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		//创建帖子 post请求
		v1.POST("/post", controller.CreatePostHandler)
		//查询单个帖子详情 请求接口
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		//查询帖子列表[]详情 请求接口
		v1.GET("/posts/", controller.GetPostListHandler)
		//post2
		//根据时间和分数来获取到 post列表
		v1.GET("/post2", controller.GetPostListHandler2)

		//投票功能
		v1.POST("/vote", controller.PostVoteController)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
