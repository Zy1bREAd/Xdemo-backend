package app

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	api "xdemo/internal/api"
	resp "xdemo/internal/api/response"
	reqData "xdemo/internal/app/tmp"
	model "xdemo/internal/database"
	"xdemo/internal/middleware"

	"github.com/gin-gonic/gin"
)

// 业务逻辑接口
func ListContainers(ctx *gin.Context) {
	result, err := api.DockerInstance.ContainersList()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "获取Docker容器列表失败", nil))
		return
	}
	for _, ctr := range result {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "List Success", result))
	middleware.DoLog(ctx, 4, "Container", fmt.Sprintln("获取容器列表"))
}

func InspectContainer(ctx *gin.Context) {
	cid := ctx.Param("cid")
	result, err := api.DockerInstance.ContainerInspect(cid)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "查询容器详细信息失败", nil))
		return
	}
	fmt.Println(result)
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Docker Inspect Success", result))
	middleware.DoLog(ctx, 4, "Container", fmt.Sprintf("检查容器详情 %s", cid))
}

// 拉取镜像
func PullImages(ctx *gin.Context) {
	defer func() {
		// 异常最终处理
		if err := recover(); err != nil {
			respMsg := fmt.Sprintf("拉取镜像发生错误: %s", err)
			log.Println(respMsg)
			ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", respMsg, nil))
			return
		}
	}()

	// 最终应该包含一个或多个的拉取（字符串切片）
	// 使用HTTP协议长连接 + 分块传输实现 =》 需要设置响应头，启用分块传输编码
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	var i reqData.ImageBody
	ctx.ShouldBindJSON(&i)
	// pullDataSlice, err := api.DockerInstance.ImagePull(i.Name, i.Tag, i.IsCover)
	reader, err := api.DockerInstance.ImagePull(i.Name, i.Tag, i.IsCover)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	//将reader中的数据传输到ctx.writer的response中
	err = resp.WriteChunkStringToClient(ctx, reader, 1024)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Image Pull Success", nil))
	middleware.DoLog(ctx, 1, "Image", fmt.Sprintln("拉取镜像"))
}

// 运行镜像成容器
func CreateContainer(ctx *gin.Context) {
	var createCfg reqData.ContainerCreateConfig
	ctx.ShouldBindJSON(&createCfg)
	fmt.Println(createCfg)
	// 判断是否传入空json body
	if createCfg.Image == "" {
		log.Println("无法创建，因为没有传入Image Name")
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.ParamsError, "Params Error or Null", "没有传入Image名字", nil))
		return
	}
	// 解析请求body，将数据传到结构体中
	err := api.DockerInstance.ParseContainerCreateConfig(&createCfg)
	if err != nil {
		log.Println("解析创建容器配置发生错误", err)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "拉取镜像失败", nil))
		return
	}
	// 使用异步机制处理
	queueClient, _ := api.GetMyQueueForRedis()
	jobInfo := &api.Job{
		UUID:       api.GenerateRandKey(),
		Name:       "ContainerCreate",
		QueueName:  "xdemo_default_task",
		Type:       "container",
		Parameters: []string{createCfg.Name, createCfg.Image},
		Service:    "create_container",
		CreateAt:   time.Now(),
	}
	jobId, err := queueClient.JobProducer(ctx, jobInfo)
	if err != nil {
		log.Println("(异步)创建容器发生错误", err)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "容器创建失败", nil))
		return
	}
	// （同步）
	// cid, err := api.DockerInstance.ContainerCreate(createCfg.Name, createCfg.Image)
	// cInfo, err := api.DockerInstance.ContainerInspectBaseInfo(cid)
	// if err != nil {
	// 	log.Println("获取容器基础信息发生错误", err)
	// 	ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "获取容器信息失败", nil))
	// 	return
	// }
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Create Container Success", jobId))
	middleware.DoLog(ctx, 1, "Container", fmt.Sprintf("创建容器,Job:%s", jobId))

}

func CreateAndRunContainer(ctx *gin.Context) {

}

func StartContainer(ctx *gin.Context) {
	cid := ctx.DefaultQuery("cid", "")
	err := api.DockerInstance.ContainerStart(cid)
	if err != nil {
		log.Println("启动容器发生错误", err)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "启动容器失败", nil))
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Contianer Start Success", nil))
	middleware.DoLog(ctx, 0, "Container", fmt.Sprintf("启动容器 %s", cid))
}

func DeleteContainer(ctx *gin.Context) {
	delOpts := map[string]bool{
		"clean": false,
		"force": false,
	}
	cid := ctx.DefaultQuery("cid", "")
	delOpts["force"] = ctx.DefaultQuery("force", "") == "true"
	delOpts["clean"] = ctx.DefaultQuery("clean", "") == "true"
	err := api.DockerInstance.ContainerDelete(cid, delOpts)
	if err != nil {
		log.Println("删除容器发生错误", err)
		errType := reflect.TypeOf(err)
		fmt.Println(errType, errType.Name())
		if errType.Name() == "errNotFound" {
			ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "删除容器失败,容器不存在", nil))
		} else if errType.Name() == "errConflict" {
			ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "删除容器失败,容器正在运行中...", nil))
		}
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "delete Success", nil))
	middleware.DoLog(ctx, 2, "Container", fmt.Sprintf("删除容器 %s", cid))
}

func RestartContainer(ctx *gin.Context) {
	var forceOpt bool
	cid := ctx.DefaultQuery("cid", "")
	force := ctx.DefaultQuery("force", "")
	if force != "" {
		forceOpt = true
	}
	err := api.DockerInstance.ContainerRestart(cid, forceOpt)
	if err != nil {
		log.Println("重启容器发生错误", err)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "重启容器失败", nil))
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Container Restart Success", nil))
	middleware.DoLog(ctx, 0, "Container", fmt.Sprintf("重启容器 %s", cid))
}

func StopContainer(ctx *gin.Context) {
	var forceOpt bool
	cid := ctx.DefaultQuery("cid", "")
	force := ctx.DefaultQuery("force", "")
	if force != "" {
		forceOpt = true
	}
	err := api.DockerInstance.ContainerStop(cid, forceOpt)
	if err != nil {
		log.Println("停止容器发生错误", err)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", "停止容器失败", nil))
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Container Stop Success", nil))
	middleware.DoLog(ctx, 0, "Container", fmt.Sprintf("停止容器 %s", cid))
}

// Exec到容器
func ExecCmdContainer(ctx *gin.Context) {
	cid := ctx.DefaultQuery("cid", "")
	var execBody model.ExecCmdBody
	err := ctx.ShouldBindBodyWithJSON(&execBody)
	if err != nil {
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.ParamsError, "Params Error", "执行命令语法有误", nil))
		return
	}
	resultSlice, err := api.DockerInstance.ContainerExecCmd(cid, execBody.Command)
	if err != nil {
		respMsg := fmt.Sprintln("进入容器Exec发生错误", err.Error())
		log.Println(respMsg)
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", respMsg, nil))
		return
	}
	// ctx.Writer.Header().Set("Content-Type", "application/vnd.docker.raw-stream")
	ctx.Writer.Write(resultSlice)
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Enter Container Exec Cmd Success", nil))
	middleware.DoLog(ctx, 0, "Container", fmt.Sprintf("Exec容器 %s - %s", cid, execBody.Command))
}

// 测试长连接
func KeepConnection(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	var counter int
	for {
		data := "data: " + fmt.Sprintf("%d\n\n", counter)
		_, err := ctx.Writer.WriteString(data)
		if err != nil {
			log.Println("发送SSE数据失败:", err)
			// ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.DeafultFailed, "Success", "Test Keep Connection", nil))
			break
		}
		ctx.Writer.Flush()
		counter++
		time.Sleep(1 * time.Second)
		if counter == 10 {
			break

		}
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Test Keep Connection", nil))
}

func EnterContainer(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			respMsg := fmt.Sprintln("Enter in Container Error,", err)
			log.Println(respMsg)
			ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "Internal Server Error", respMsg, nil))
		}
	}()

	cid := ctx.DefaultQuery("cid", "")
	if cid == "" {
		panic(fmt.Errorf("contianer Name Not Empty"))
	}
	// 升级Websocket协议
	wsDefaultConfig := api.DefaultWSConfig()
	wsConn, err := wsDefaultConfig.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	// 初始化WebSocket通信
	wsInstance := api.NewMyWebSocket(wsConn)
	wsInstance.Start()
	// 正式开始双方收发消息
	for {
		clientCmd, err := wsInstance.ReadMessage()
		if err != nil {
			panic(err)
		}
		if string(clientCmd) == "exit" {
			// 退出逻辑
			wsInstance.Close()
			return
		}
		// 将输入字节切片数据转成命令（字符串切片）
		cmdStringSlice := api.ConvertToCmdSlice(clientCmd)

		execResult, err := api.DockerInstance.ContainerExecCmd(cid, cmdStringSlice)
		if err != nil {
			panic(err)
		}
		// 返回的字节切片数据，总是前面会有莫名的前缀
		err = wsInstance.WriteMessage(execResult[8:])
		if err != nil {
			panic(err)
		}
	}

}

func TagImage(ctx *gin.Context) {
	newImage := ctx.DefaultPostForm("target", "")
	if newImage == "" {
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.ParamsNull, "ParamsIsNull", "Target Image is Null", nil))
		return
	}
	orginalImageName := ctx.DefaultQuery("iname", "")
	orginalImageTag := ctx.DefaultQuery("itag", "latest")
	source := orginalImageName + ":" + orginalImageTag
	// 是否需要校验新的Image Tag合法性
	err := api.DockerInstance.UpdateImageTag(source, newImage)
	if err != nil {
		respMsg := fmt.Sprintf("UpdateImageTag is Failed,Error:%s", err.Error())
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "UpdateFailed", respMsg, nil))
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("Success", "Update Success", map[string]string{
		"new_image": newImage,
	}))
	middleware.DoLog(ctx, 4, "Image", "Tag Image")
}

func LoginDockerRegistry(ctx *gin.Context) {
	var registryAuth struct {
		Username string `json:"login_usr"`
		Password string `json:"login_pwd"`
		Addr     string `json:"registry_addr"`
	}
	ctx.ShouldBindJSON(&registryAuth)
	fmt.Println(registryAuth)
	// ctx.Writer.Header().Set(registry.AuthHeader)
	err := api.DockerInstance.RegistryLogin(ctx, registryAuth.Username, registryAuth.Password, registryAuth.Addr)
	if err != nil {
		respMsg := fmt.Sprintf("Login Docker Registry is Failed,%s", err.Error())
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "RegistryLoginFailed", respMsg, nil))
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("RegistryLoginSuccess", "Registry Login Success", map[string]string{
		"registry_auth": api.DockerInstance.RegistryAuthToken,
	}))
	middleware.DoLog(ctx, 4, "Docker Registry", "Login Docker Registry")
}

// 针对前端表单数据中传递的image进行推送到仓库中
func PushImage(ctx *gin.Context) {
	target := ctx.DefaultPostForm("image", "")
	if target == "" {
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.ParamsNull, "ParamsIsNull", "Target Image is Null", nil))
		return
	}
	pushResultReader, err := api.DockerInstance.ImagePush(ctx, target)
	if err != nil {
		respMsg := fmt.Sprintf("Push Image is Failed,%s", err.Error())
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "ImagePushnFailed", respMsg, nil))
		return
	}
	defer pushResultReader.Close()
	// 设置长连接的头
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	err = resp.WriteChunkStringToClient(ctx, pushResultReader, 2048)
	if err != nil {
		respMsg := fmt.Sprintf("Write Data To Gin Response Failed,Error:%s", err.Error())
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "ImagePushnFailed", respMsg, nil))
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("ImagePushSuccess", "Image Push Success", nil))
	middleware.DoLog(ctx, 4, "Image", "Push Image")
}
