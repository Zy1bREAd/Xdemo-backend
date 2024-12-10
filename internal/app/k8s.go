package app

import (
	"net/http"
	api "xdemo/internal/api"
	resp "xdemo/internal/api/response"

	"github.com/gin-gonic/gin"
)

func GetAllPods(ctx *gin.Context) {
	err := api.K8sInstance.GetPodsForDefault()
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Get Success", "Get Pods Success", nil))
}
