package api

import (
	"context"
	"log"
)

type QueueProvider interface {
}

type QueueForRedis struct {
	RedisClient *MyRedis // 结构体直接存储该指针，到时声明就不需要&解引用
}

var QueueInstance *QueueForRedis

// 初始化消息队列
func InitQueue() {
	QueueInstance := NewMyQueue()
	ctx := context.Background()
	// 测试Redis连通性
	err := QueueInstance.RedisClient.Ping(ctx)
	if err != nil {
		log.Println("Init MessageQueue Failed,", err)
		panic(err)
	}
	log.Println("Init MessageQueue Success.")
}

// 直接调用该函数创建instance，然后使用这个实例去操作消息队列
func NewMyQueue() *QueueForRedis {
	if QueueInstance == nil {
		return &QueueForRedis{
			RedisClient: RedisInstance,
		}
	}
	return QueueInstance
}

func (q *QueueForRedis) SetJob() {
	//
}
