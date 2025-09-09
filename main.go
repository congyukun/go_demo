package main

import (
	"go_demo/config"
	"go_demo/db"
	"go_demo/routes"
	"strconv"
)

func main() {

	// 初始化配置
	if err := config.InitConfig(""); err != nil {
		panic("config init error:" + err.Error())
	}
	cfg := config.GlobalConfig
	serverConfig := cfg.Server
	if serverConfig.Port == 0 {
		panic("port is invalid")
	}
	// 初始化数据库
	if err := db.InitMySQLConn(); err != nil {
		panic("mysql connect error:" + err.Error())
	}
	defer db.Close()

	// 初始化路由
	r := routes.SetupRouter()

	// 启动服务
	addr := ":" + strconv.Itoa(serverConfig.Port)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
