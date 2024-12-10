package router

import (
	app "xdemo/internal/app"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		rgPublic.POST("/login", app.UserLogin)
		rgPublic.POST("/register", app.UserRegister)
		rgPublic.POST("/validate", app.ValidateUserInfo)
		rgAuth.GET("/user", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "get all user data Success!",
				"data": map[string]any{
					"id":   123,
					"name": "boniu",
				},
			})
		})
		rgAuth.GET("/user/:id", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "get user data Success!",
				"data": map[string]any{
					"id":   123,
					"name": "boniu",
				},
			})
		})
	})
}
