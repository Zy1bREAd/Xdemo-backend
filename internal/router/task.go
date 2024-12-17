package router

import (
	"xdemo/internal/app"

	"github.com/gin-gonic/gin"
)

func InitTaskPlatformRoute() {
	RegisterRoute(func(rgPublic, rgAuth *gin.RouterGroup) {
		rgAuth.POST("/task/create", app.CreateDefaultTask)
	})
}
