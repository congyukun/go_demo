package routes

import (
	"go_demo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterArticleRoutes(g *gin.Engine) {
	articleController := &controllers.ArticleController{}
	g.POST("/article", articleController.CreateArticleHandler)
	g.GET("/article/:id", articleController.GetArticleHandler)
	g.PUT("/article/:id", articleController.UpdateArticleHandler)
	g.DELETE("/article/:id", articleController.DeleteArticleHandler)
	g.GET("/articles", articleController.ListArticlesHandler)
}

// 以下处理函数已由 ArticleController 实现，详见 controllers/article.go
