package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HealthCheckInterval = 30
)

// var WSConn MyWebSocket

type MyWebSocket struct {
	WS       *websocket.Conn
	InChan   chan []byte
	OutChan  chan []byte
	ExitChan chan struct{}
	IsClose  bool
}

func DefaultWSConfig() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func NewMyWebSocket(ws *websocket.Conn) MyWebSocket {
	return MyWebSocket{
		WS:       ws,
		InChan:   make(chan []byte),
		OutChan:  make(chan []byte),
		ExitChan: make(chan struct{}),
	}
}

func (x *MyWebSocket) Close() {
	if x.IsClose {
		log.Println("WebSocket is Closed")
		return
	}
	x.ExitChan <- struct{}{}
	close(x.InChan)
	close(x.OutChan)
	close(x.ExitChan)
	x.WS.Close()
	x.IsClose = true
}

func (x *MyWebSocket) Start() {
	// 接收和发送消息通过goroutine进入循环，因此还有一个类似守护进程的goroutine去完成心跳健康检查、重连等操作
	// go x.HealthCheck()
	go x.ReadLoop()
	go x.WriteLoop()
}

func (x *MyWebSocket) ReadMessage() ([]byte, error) {
	// 从Channel中返回读取到的数据
	if x.IsClose {
		return nil, fmt.Errorf("webSocket Connection is Closed")
	}
	return <-x.InChan, nil
}

func (x *MyWebSocket) WriteMessage(data []byte) error {
	if x.IsClose {
		return fmt.Errorf("webSocket Connection is Closed")
	}
	x.OutChan <- data
	return nil
}

func (x *MyWebSocket) ReadLoop() {
	for {
		// 超时检查（若长时间没有接收消息则关闭连接，节省资源）
		x.WS.SetReadDeadline(time.Now().Add(HealthCheckInterval * 60 * time.Second))
		_, msg, err := x.WS.ReadMessage()
		if err != nil {
			// 有错误就关闭连接
			log.Println("ReadMessage occur a Error: ", err)
			x.Close()
			return
		}
		if string(msg) == "quit" {
			x.Close()
			break
		}
		x.InChan <- msg
	}
}

func (x *MyWebSocket) WriteLoop() {
	for {
		select {
		case msg := <-x.OutChan:
			err := x.WS.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("WriteMessage occur Error: ", err)
				x.Close()
				return
			}
		case <-x.ExitChan:
			// 收到退出信号，退出Loop
			log.Println("收到退出信号Exit")
			return
		}
	}
}

func (x *MyWebSocket) HealthCheck() {
	for {
		go func() {

		}()
		// 健康检查的逻辑(需要前端回复pong消息)
		select {
		case <-x.ExitChan:
			log.Println("Quit HC")
			return
		case <-time.After(time.Second * 10):
			x.Close()
			return
		default:
			err := x.WS.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("Write Ping Message Error", err)
			}

			log.Println("心跳存活！！")
		}
		// HC间隔(秒)
		time.Sleep(HealthCheckInterval * time.Second)
	}
}
