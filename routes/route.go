package routes

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	// 注册路由
	RegisterUserRoutes(r)    // 用户相关路由
	RegisterArticleRoutes(r) // 文章相关路由

	return r
}
