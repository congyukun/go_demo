package routes

import (
	"go_demo/controllers"

	"github.com/gin-gonic/gin"
)

// 注册用户路由

func RegisterUserRoutes(g *gin.Engine) {
	// 用户注册路由
	loginController := &controllers.LoginController{}
	g.POST("/register", loginController.Register)
	g.POST("/login", loginController.Login)
}
