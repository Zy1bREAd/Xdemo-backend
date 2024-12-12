package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	config "xdemo/internal/config"
	db "xdemo/internal/database"
	middleware "xdemo/internal/middleware"

	"github.com/gin-gonic/gin"
)

// rgPublic用于公共API，而rgAuth则需要鉴权（需要token验证）的API
type FnRegisterRoute func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

// 这是一个收集各个模块路由注册函数的切片
var fnRoutes []FnRegisterRoute

// 初始化路由
func InitRouter() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	r := gin.Default()
	// 添加中间件到handleFunc链
	cc := &middleware.CorsConfig{
		AllowOrigin:  []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Accept", "Authorization", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Bearer"},
	}
	r.Use(cc.Cors())
	// 定义路由组
	rgPublic := r.Group("/api/v1/public")
	rgAuth := r.Group("/api/v1")
	// 添加鉴权中间件
	rgAuth.Use(middleware.AuthMiddleware())

	InitBaseRoutes()
	for _, fn := range fnRoutes {
		fn(rgPublic, rgAuth)
	}
	configProvider := config.NewConfigEnvProvider()
	srv := &http.Server{
		// Addr: config.System().ListenHost().Env() + ":" + config.System().ListenPort().Env(),
		Addr:    configProvider.System.Host + ":" + configProvider.System.Port,
		Handler: r,
	}
	go func() {
		// 建立http服务端通信连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Start Gin Server Error is %s", err)
		}
	}()
	// 等待信号量的出现
	<-ctx.Done()
	// 接收到信号量后进行优雅关闭服务器
	ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown Gin Server Failed,Error :%s", err)
		return
	}
	log.Println("Shutdown Gin Success!!!")
	// 顺便关闭db连接
	db.CloseDB()

}

// 对路由函数进行注册
func RegisterRoute(fn FnRegisterRoute) {
	if fn == nil {
		fmt.Println("不需要注册")
		return
	}
	fnRoutes = append(fnRoutes, fn)
}

// 初始化基础路由信息
func InitBaseRoutes() {
	InitPublicRoutes()
	// 里面就是各个模块的注册路由函数，由他们去实际控制对应路由的handleFunc
	InitUserRoutes()
	InitContainerRoute()
	InitK8sClusterRoutes()
}
