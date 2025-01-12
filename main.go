package main

import (
	"log"
	"xdemo/internal/api"
	"xdemo/internal/config"
	"xdemo/internal/database"
	"xdemo/internal/middleware"
	"xdemo/internal/router"
	// api "xdemo/internal/api"
	// "xdemo/internal/config"
	// db "xdemo/internal/database"
	// middleware "xdemo/internal/middleware"
	// router "xdemo/internal/router"
)

func main() {
	// 异常处理
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			return
		}
	}()
	// 准备日志记录的配置
	middleware.PrepLogger()
	// 初始化校验器
	middleware.InitGlobalValidator()
	// 加载配置文件
	// yamlConfig := config.LoadConfig()
	config.InitConfigEnv()
	database.LoadDB()
	api.InitRedis()
	api.InitQueue()
	api.InitDocker()
	api.InitK8sClient()

	// 尝试读取env 变量
	router.InitRouter()
}
