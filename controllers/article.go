package controllers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)


// Article 结构体
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ArticleController 控制器，内存存储文章
type ArticleController struct{}

var (
	articleStore = make(map[int]*Article)
	articleIDSeq = 1
	storeMutex   sync.Mutex
)

// 创建文章
// @Summary 创建文章
// @Description 创建一篇新文章
// @Accept json
// @Produce json
// @Param article body Article true "文章内容"
// @Success 200 {object} Article
// @Failure 400 {string} string "参数错误"
// @Router /article [post]
func (ac *ArticleController) CreateArticleHandler(c *gin.Context) {
var req Article
if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" || req.Content == "" {
	response(c, http.StatusBadRequest, "参数错误，title 和 content 必填", nil)
	return
}
storeMutex.Lock()
req.ID = articleIDSeq
articleStore[articleIDSeq] = &req
articleIDSeq++
storeMutex.Unlock()
response(c, http.StatusOK, "创建文章成功", req)
}

// 获取单篇文章
// @Summary 获取文章
// @Description 根据ID获取文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} Article
// @Failure 400 {string} string "ID格式错误"
// @Failure 404 {string} string "文章不存在"
// @Router /article/{id} [get]
func (ac *ArticleController) GetArticleHandler(c *gin.Context) {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
	response(c, http.StatusBadRequest, "ID格式错误", nil)
	return
}
storeMutex.Lock()
article, ok := articleStore[id]
storeMutex.Unlock()
if !ok {
	response(c, http.StatusNotFound, "文章不存在", nil)
	return
}
response(c, http.StatusOK, "获取文章成功", article)
}

// 更新文章
// @Summary 更新文章
// @Description 根据ID更新文章内容
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param article body Article true "文章内容"
// @Success 200 {object} Article
// @Failure 400 {string} string "参数错误"
// @Failure 404 {string} string "文章不存在"
// @Router /article/{id} [put]
func (ac *ArticleController) UpdateArticleHandler(c *gin.Context) {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
	response(c, http.StatusBadRequest, "ID格式错误", nil)
	return
}
var req Article
if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" || req.Content == "" {
	response(c, http.StatusBadRequest, "参数错误，title 和 content 必填", nil)
	return
}
storeMutex.Lock()
article, ok := articleStore[id]
if !ok {
	storeMutex.Unlock()
	response(c, http.StatusNotFound, "文章不存在", nil)
	return
}
article.Title = req.Title
article.Content = req.Content
storeMutex.Unlock()
response(c, http.StatusOK, "更新文章成功", article)
}

// 删除文章
// @Summary 删除文章
// @Description 根据ID删除文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {string} string "删除成功"
// @Failure 400 {string} string "ID格式错误"
// @Failure 404 {string} string "文章不存在"
// @Router /article/{id} [delete]
func (ac *ArticleController) DeleteArticleHandler(c *gin.Context) {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
	response(c, http.StatusBadRequest, "ID格式错误", nil)
	return
}
storeMutex.Lock()
_, ok := articleStore[id]
if ok {
delete(articleStore, id)
}
storeMutex.Unlock()
if !ok {
response(c, http.StatusNotFound, "文章不存在", nil)
return
}
response(c, http.StatusOK, "删除文章成功", nil)
}

// 文章列表
// @Summary 文章列表
// @Description 获取所有文章
// @Produce json
// @Success 200 {array} Article
// @Router /articles [get]
func (ac *ArticleController) ListArticlesHandler(c *gin.Context) {
storeMutex.Lock()
articles := make([]*Article, 0, len(articleStore))
for _, a := range articleStore {
	articles = append(articles, a)
}
storeMutex.Unlock()
response(c, http.StatusOK, "获取文章列表成功", articles)
}
