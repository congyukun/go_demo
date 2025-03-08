package controllers

import "github.com/gin-gonic/gin"

type ArticleController struct{}

// 以下处理函数暂时留空，后续实现
func (articl *ArticleController) CreateArticleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建文章成功!",
	})
}
func (articl *ArticleController) GetArticleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建文章成功!",
	})
}
func (articl *ArticleController) UpdateArticleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建文章成功!",
	})
}
func (articl *ArticleController) DeleteArticleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建文章成功!",
	})
}
func (articl *ArticleController) ListArticlesHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "创建文章成功!",
	})
}
