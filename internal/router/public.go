package router

import (
	"fmt"
	"log"
	"net/http"
	api "xdemo/internal/api"
	resp "xdemo/internal/api/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func InitPublicRoutes() {
	// 其中定义验证码请求接口
	RegisterRoute(func(rgPublic, rgAuth *gin.RouterGroup) {
		rgPublic.GET("/getCaptcha", api.RequestLoginCaptcha)

		// 测试API
		rgPublic.POST("/test", func(ctx *gin.Context) {
			tokenStr := ctx.Request.Header.Get("Authorization")
			log.Println(tokenStr)

			// 测试解析
		})

		// 测试Websocket
		rgAuth.GET("/wstest", func(ctx *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					respMsg := fmt.Sprintln("Recovered in Websocket test", err)
					log.Println(respMsg)
					ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.InternalServerError, "InternalServerError", respMsg, nil))
					return
				}
			}()
			var wsConfig = websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			ws, err := wsConfig.Upgrade(ctx.Writer, ctx.Request, nil)
			if err != nil {
				fmt.Println("发生ws升级error", err)
				return
			}
			hcCancel := make(chan struct{})
			defer func() {
				// 关闭websocket连接以及health-check检测
				ws.Close()
				hcCancel <- struct{}{}
			}()

			// WS健康检测

			ws.SetPingHandler(func(appData string) error {
				fmt.Println("收到ping，确认存活")
				return ws.WriteMessage(websocket.PongMessage, []byte(appData))
			})
			// ws.SetPongHandler(func(appData string) error {
			// 	fmt.Println("收到pong！")
			// 	return nil
			// })

			conn := api.NewMyWebSocket(ws)
			conn.Start()
			for {
				fmt.Println("ws start....")
				clientData, err := conn.ReadMessage()
				fmt.Println("read completed")
				if err != nil {
					break
				}
				fmt.Println("读取到消息2 CLient Data:", string(clientData))

				err = conn.WriteMessage([]byte("???????????????a"))
				fmt.Println("write completed")
				if err != nil {
					break
				}
			}
			fmt.Println("ws end...")

			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "Success Use WebSocket!",
			})
		})
	})
}

// func wsHealthCheck(cancel chan struct{}) {
// 	// select信号控制
// 	for {
// 		fmt.Println("正在进行心跳检测！！！！")

// 		select {
// 		case <-cancel:
// 			fmt.Println("收到channel消息，取消当前goroutine")
// 			return
// 		default:
// 			fmt.Println("default nothing")
// 		}
// 		time.Sleep(2 * time.Second)
// 	}
// }
