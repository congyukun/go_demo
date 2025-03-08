package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterArticleRoutes(g *gin.Engine) {
	g.POST("/article", CreateArticleHandler)
	g.GET("/article/:id", GetArticleHandler)
	g.PUT("/article/:id", UpdateArticleHandler)
	g.DELETE("/article/:id", DeleteArticleHandler)
	g.GET("/articles", ListArticlesHandler)
}

// 以下处理函数暂时留空，后续实现
func CreateArticleHandler(c *gin.Context) {}
func GetArticleHandler(c *gin.Context)    {}
func UpdateArticleHandler(c *gin.Context) {}
func DeleteArticleHandler(c *gin.Context) {}
func ListArticlesHandler(c *gin.Context)  {}
