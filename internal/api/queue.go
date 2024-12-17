package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	config "xdemo/internal/config"
	"xdemo/internal/database"
	"xdemo/internal/global"

	"github.com/redis/go-redis/v9"
)

// ! 关于 Job 的信息
type Job struct {
	UUID       string    `json:"job_id"`        // Job的唯一标识
	Name       string    `json:"job_name"`      // Job的名字
	QueueName  string    `json:"job_queue_ame"` // Job所属队列名字
	Type       string    `json:"job_type"`      // 操作资源的类型
	Service    string    `json:"job_service"`   // 相对应要操作函数的名字
	Parameters []string  `json:"job_params"`    // 传递到消费者要接入去执行函数的参数
	CreateAt   time.Time `json:"job_create_at"` // Job创建时间
}

type QueueProvider interface {
	Ping() error
}

type QueueForRedis struct {
	RedisClient *MyRedis // 结构体直接存储该指针，到时声明就不需要&解引用
	Processer   int
}

var queueInstance QueueProvider

// 初始化消息队列
func InitQueue() {
	queueInstance = NewMyQueue()
	// 测试Redis连通性
	err := queueInstance.Ping()
	if err != nil {
		log.Println("JobQueue Init Failed,", err)
		panic(err)
	}
	ProcesserReady()
	log.Println("JobQueue Init Success ~")
}

// 直接调用该函数创建instance，然后使用这个实例去操作消息队列
func NewMyQueue() QueueProvider {
	config := config.NewConfigEnvProvider()
	if queueInstance == nil {
		switch strings.ToLower(config.TaskQueue.Provider) {
		case "redis":
			return &QueueForRedis{
				RedisClient: RedisInstance,
				Processer:   config.TaskQueue.Processer,
			}
		case "rabbitmq":
			fmt.Println("rabbitmq Code Slot")
		default:
			log.Println("消息队列暂不支持 ", config.TaskQueue.Provider)
		}
	}
	return queueInstance
}

// 判断接口是否属于某个结构体
func isRedisQueueProvider(q QueueProvider) bool {
	_, ok := q.(*QueueForRedis)
	return ok
}

func GetMyQueueForRedis() (*QueueForRedis, error) {
	if queueInstance == nil {
		return nil, fmt.Errorf("QueueInstance未初始化")
	}
	ok := isRedisQueueProvider(queueInstance)
	if ok {
		return queueInstance.(*QueueForRedis), nil
	}
	return nil, fmt.Errorf("QueueInstance Type不正确")
}

// 测试Task Queue连通性
func (q *QueueForRedis) Ping() error {
	ctx := context.Background()
	return q.RedisClient.Ping(ctx)
}

func (q *QueueForRedis) GetQueueClient() *MyRedis {
	return q.RedisClient
}

func (q *QueueForRedis) JobEntryList(ctx context.Context, queueName string, taskId string) (string, error) {
	// 为task创建一个uuid
	client := q.GetQueueClient()
	err := client.LPush(ctx, queueName, taskId)
	if err != nil {
		log.Println("Create a Task in Queue Failed,", err)
		return "", err
	}
	// 存入db记录
	jobInfo := &database.Job{
		UUID:   taskId,
		Name:   queueName,
		Status: "Pending",
	}
	result := global.GDB.Create(&jobInfo)
	if result.Error != nil {
		return "", result.Error
	}
	return taskId, nil
}

func (q *QueueForRedis) JobDoHandle(ctx context.Context, taskID string, timeout time.Duration, queueName ...string) map[string]string {
	// 消费者取出进行处理
	resultCh := make(chan map[string]string)
	// 构建一个超时ctx
	ctx, cancel := context.WithTimeout(ctx, time.Second*300)
	defer cancel()

	go func() {
		resultSlice, err := q.RedisClient.BRPop(ctx, timeout, queueName...)
		if err != nil {
			log.Println("获取Task失败,", err)
			resultCh <- map[string]string{
				"Error": err.Error(),
			}
			return
		}
		// 更新Job的状态
		dbUpdate := global.GDB.Model(&database.Job{}).Where("job_uuid = ?", taskID).Update("job_status", "Running")
		if dbUpdate.Error != nil {
			resultCh <- map[string]string{
				"Error": dbUpdate.Error.Error(),
			}
			return
		}
		resultMap := map[string]string{
			resultSlice[0]: resultSlice[1],
		}
		resultCh <- resultMap
	}()

	select {
	case <-ctx.Done():
		log.Println("context主动timeout，结束Task处理")
		return nil
	case re := <-resultCh:
		return re
	}
}

func (q *QueueForRedis) JobComplete(ctx context.Context, taskID string) error {
	// 更新Job完成状态
	dbUpdate := global.GDB.Model(&database.Job{}).Where("job_uuid = ?", taskID).Where("job_status = ?", "Running").Update("job_status", "Completed")
	if dbUpdate.Error != nil {
		log.Println("Job Update Status \"Completed\" is Failed", dbUpdate.Error.Error())
		return dbUpdate.Error
	}
	return nil
}

// job生产者(返回task ID和error)
func (q *QueueForRedis) JobProducer(ctx context.Context, jobInfo *Job) (string, error) {
	// 序列化成JSON传递JobInfo
	jobJson, err := json.Marshal(jobInfo)
	if err != nil {
		log.Println("JSON序列化Job信息Failed", err)
		return "", err
	}
	redisClient := q.GetQueueClient()
	err = redisClient.LPush(ctx, jobInfo.QueueName, jobJson)
	if err != nil {
		log.Println("生产者Push Job入队Failed", err)
		return "", err
	}
	return jobInfo.UUID, nil
}

// 初始化并准备好Job消费者
func ProcesserReady() {
	// 启动多少个processer由配置决定
	config := config.NewConfigEnvProvider()
	// config.TaskQueue.Processer
	q, err := GetMyQueueForRedis()
	if err != nil {
		panic(err)
	}
	for i := 0; i < config.TaskQueue.Processer; i++ {
		// goroutine启动Processer函数
		go q.startJobConsumer(i + 1)
	}
}

// Job消费者（监听并等待执行）
func (q *QueueForRedis) startJobConsumer(number int) {
	// 获取redis客户端
	redisClient := q.GetQueueClient()
	ctx := context.Background()
	QueueName := "xdemo_default_task"
	log.Printf("启动Job消费者 %d 号\n", number)
	for {
		qResult, err := redisClient.BRPop(ctx, time.Second*3600, QueueName)
		if err != nil {
			// Redis队列为Null的Error不打印
			if err != redis.Nil {
				log.Println("Redis队列取出Job消息失败,", err)
			}
			continue
		}
		fmt.Printf("Job消费者 %d 号取到Job %s\n", number, qResult[1])
		// 序列化Job的JSON信息
		strToByte := []byte(qResult[1])
		var job Job
		err = json.Unmarshal(strToByte, &job)
		if err != nil {
			log.Println("消费者解析Job信息失败,", err)
			continue
		}
		fmt.Println(job)
		// 根据Job类型寻找合适的函数进行处理Job
		fmt.Printf("Job消费者 %d 号进行处理\n", number)
		switch strings.ToLower(job.Type) {
		case "container":
			task := &DockerTaskExecutor{}
			task.Execute(&job)
		default:
			log.Println("无法处理")
		}
	}
}

// 实现一个任务执行者接口
type TaskExecutor interface {
	Execute(task *Job) error
}

// 实现一个任务执行者的Docker接口
type DockerTaskExecutor struct {
}

func (d *DockerTaskExecutor) Execute(job *Job) error {
	// 根据Job的Service来处理
	switch job.Service {
	case "create_container":
		// 创建容器的实际函数
		cid, err := DockerInstance.ContainerCreate(job.Parameters[0], job.Parameters[1])
		if err != nil {
			log.Println("创建容器发生错误", err)
			return err
		}
		fmt.Println("成功（异步）,cid:", cid)
	case "test_test":
		fmt.Println("测试测试测试异步场景！")
	}
	return nil
}
