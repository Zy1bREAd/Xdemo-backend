package router

import (
	app "xdemo/internal/app"

	"github.com/gin-gonic/gin"
)

func InitContainerRoute() {
	RegisterRoute(func(rgPublic, rgAuth *gin.RouterGroup) {
		rgAuth.GET("/containers/list", app.ListContainers)
		rgAuth.GET("/containers/inspect/:cid", app.InspectContainer)
		rgAuth.POST("/container/create", app.CreateContainer)
		rgAuth.POST("/container/run", app.CreateAndRunContainer)
		rgAuth.POST("/container/start", app.StartContainer)
		rgAuth.POST("/container/stop", app.StopContainer)
		rgAuth.POST("/container/restart", app.RestartContainer)
		rgAuth.DELETE("/container/delete", app.DeleteContainer)

		rgAuth.POST("/container/exec", app.ExecCmdContainer)
		rgAuth.GET("/container/enter", app.EnterContainer)
		rgAuth.POST("/container/keepconnection", app.KeepConnection)
		rgAuth.POST("/image/pull", app.PullImages)
		rgAuth.POST("/image/tag", app.TagImage)
		rgAuth.POST("/image/push", app.PushImage)

		rgAuth.POST("/docker/login", app.LoginDockerRegistry)
	})
}
