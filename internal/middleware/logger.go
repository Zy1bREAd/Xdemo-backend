package middleware

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	db "xdemo/internal/database"

	global "xdemo/internal/global"

	"github.com/gin-gonic/gin"
)

type MyLogger struct {
	CTX       context.Context
	CancelAt  time.Time
	HandleLog db.HandleLog
}

// 用于日志记录的准备工作
func PrepLogger() {
	//  暂时不写入到文件中
	// logFile, err := os.OpenFile("xdemo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logFile)
	log.SetPrefix("[XDemo] ")
	log.Println("Demo with Golang -- Ocean Wang")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in main", r)
		}
	}()

}

// 将Log记录下来
func DoLog(ctx *gin.Context, operationType int, operationResource string, operationDetails string) error {
	userInfo := ctx.GetStringMapString("userInfo")
	handlerId, err := strconv.ParseUint(userInfo["id"], 10, 32)

	if err != nil {
		return err
	}
	// 写入DB记录
	record := &db.HandleLog{
		Handler:  uint(handlerId), // 暂时使用user id作为用户标识
		Type:     operationType,
		Resource: operationResource,
		Detail:   operationDetails,
	}
	resultObj := global.GDB.Create(record)
	if resultObj.Error != nil {
		log.Printf("Insert Account:[%s] Handle Log Data Failed,The Error is %s\n", userInfo["account"], resultObj.Error)
		return resultObj.Error
	}
	if resultObj.RowsAffected != 1 {
		log.Println("Insert Handle Log Data Failed!!!")
		return fmt.Errorf("insert row is zero")
	}
	return nil
}
