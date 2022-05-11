package router

import (
	"gin_demo/controllers"
	"gin_demo/logger"
	"gin_demo/middlewares"
	"net/http"

	_ "gin_demo/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 引入日志中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 每两秒钟添加一个令牌  全局限流
	//r.Use(logger.GinLogger(), logger.GinRecovery(true),middlewares.RateLimitMiddleware(2*time.Second , 1))

	// 接口文档：http://localhost:8083/swagger/index.html#/
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	v1.POST("/signup", controllers.SignUpHandler)
	v1.POST("/login", controllers.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/post/list", controllers.GetPostListHandler)
		v1.GET("/post/lists", controllers.GetPostListByCreateTimeOrScoreHandler)

		v1.POST("/vote", controllers.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
