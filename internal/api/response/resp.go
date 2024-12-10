package api

import (
	"fmt"
	"io"
	"log"

	// constant "xdemo/internal/global"

	"github.com/gin-gonic/gin"
)

// // 定义响应代码
const (
	Success             = 0
	ParamsNull          = 1
	ParamsError         = 2
	SignatureError      = 3
	RequestTimeOut      = 4
	InternalServerError = 5
	DeafultFailed       = 8 // 该错误通常要重新发起请求

	TokenExpired    = 1001
	TokenNotFound   = 1002
	TokenError      = 1003
	TokenRenewError = 1004

	Malformed = 1101

	Unauthorized = 41
	Unknown      = 666
)

// 封装API响应体
type APIResp struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"msg"`
	Data    any
}

// !封装gin.H来返回响应数据
func SuccessRespJSON(status string, msg string, data any) gin.H {
	return gin.H{
		"code":   Success,
		"status": status,
		"msg":    msg,
		"data":   data,
	}
}

func FailedRespJSON(errCode int, status string, msg string, data any) gin.H {
	log.Printf("[%s] - 操作发生错误,%s", status, msg)
	return gin.H{
		"code":   errCode,
		"status": status,
		"msg":    msg,
		"data":   data,
	}
}

// 未知的也算错误的一种
func UnknownRespJSON(status string, msg string, data any) gin.H {
	return gin.H{
		"code":   Unknown,
		"status": status,
		"msg":    msg,
		"data":   data,
	}
}

// 将字节流数据返回给Client
func ConvertByteToSSE(msg []byte) string {
	return fmt.Sprintf("data:%s\n\n", string(msg))
}

// 将输入流的数据转换成字符串，写入到writer中
func WriteChunkStringToClient(ctx *gin.Context, reader interface{}, chunkSize int) error {
	if val, ok := reader.(io.ReadCloser); ok {
		tempBf := make([]byte, chunkSize)
		// var tempBf []byte
		for {
			n, err := val.Read(tempBf)
			fmt.Println(string(tempBf[:n]))
			if err == io.EOF {
				respData := ConvertByteToSSE(tempBf[:n])
				_, err = ctx.Writer.WriteString(respData)
				return err
			} else if err != nil {
				return err
			}
			respData := ConvertByteToSSE(tempBf[:n])
			_, err = ctx.Writer.WriteString(respData)
			if err != nil {
				return err
			}
			// ctx.Writer.Flush()
		}
	}
	return fmt.Errorf("仅支持io.ReadCloser类型")
}
