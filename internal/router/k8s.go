package router

import (
	"xdemo/internal/app"

	"github.com/gin-gonic/gin"
)

func InitK8sClusterRoutes() {
	RegisterRoute(func(rgPublic, rgAuth *gin.RouterGroup) {
		rgAuth.GET("/k8s/pods/get", app.GetAllPods)
	})
}
