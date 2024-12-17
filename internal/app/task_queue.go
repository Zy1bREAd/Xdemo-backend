package app

import (
	"fmt"
	api "xdemo/internal/api"
	resp "xdemo/internal/api/response"

	"github.com/gin-gonic/gin"
)

func CreateDefaultTask(ctx *gin.Context) {
	queueClient, _ := api.GetMyQueueForRedis()
	taskUUID := api.GenerateRandKey()
	task, err := queueClient.JobEntryList(ctx, "xdemo_default_1", taskUUID)
	if err != nil {
		respMsg := fmt.Sprintf("创建默认Task失败,%s", err)
		resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", respMsg, nil)
		return
	}
	resp.SuccessRespJSON("CreateTaskSuccess", "创建默认任务成功", task)
}

func HandleTask(ctx *gin.Context) {
	queueClient, _ := api.GetMyQueueForRedis()
	taskUUID := api.GenerateRandKey()
	task, err := queueClient.JobEntryList(ctx, "xdemo_default_1", taskUUID)
	if err != nil {
		respMsg := fmt.Sprintf("创建默认Task失败,%s", err)
		resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", respMsg, nil)
		return
	}
	resp.SuccessRespJSON("CreateTaskSuccess", "创建默认任务成功", task)
}
