package main

import (
	"fmt"
	"log"
	"os"
	api "xdemo/internal/api"
	"xdemo/internal/config"
	db "xdemo/internal/database"
	middleware "xdemo/internal/middleware"
	router "xdemo/internal/router"
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
	config.InitConfigEnv("local")
	db.LoadDB()
	api.InitRedis()
	api.InitDocker()
	api.InitK8sClient()

	// 尝试读取env 变量
	result, exits := os.LookupEnv("XDEMO_VERSION")
	fmt.Printf("变量结果是:%s,是否存在该env:%v", result, exits)
	router.InitRouter()
}
