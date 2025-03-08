package main

import (
	// "go_demo/controllers"
	// "go_demo/registry"
	"go_demo/routes"
	// "github.com/gin-gonic/gin"
)

func main() {
	// 初始化路由
	r := routes.SetupRouter()
	// 启动服务
	r.Run(":8080")
}
