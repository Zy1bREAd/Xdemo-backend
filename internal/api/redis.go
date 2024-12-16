package api

import (
	"context"
	"fmt"
	"log"
	"time"
	config "xdemo/internal/config"

	"github.com/redis/go-redis/v9"
)

const DeafultExpireWithRedis int = 60 * 24

type MyRedis struct {
	// Addr     string
	// Password string
	// DB       int
	// TLS      bool
}

var RedisInstance *MyRedis
var RDBClient *redis.Client

// 初始化Redis操作
func InitRedis() {
	RedisInstance = NewMyRedis()
	RDBClient = RedisInstance.NewRedisClient()
	if RDBClient == nil {
		panic(fmt.Errorf("初始化Redis出现错误"))
	}
}

func NewMyRedis() *MyRedis {
	if RedisInstance == nil {
		RedisInstance = &MyRedis{}
	}
	return RedisInstance
}

func (r *MyRedis) NewRedisClient() *redis.Client {
	// 获取redis配置信息
	configProvider := config.NewConfigEnvProvider()
	client := redis.NewClient(&redis.Options{
		Addr:     configProvider.Redis.Addr + ":" + configProvider.Redis.Port,
		DB:       configProvider.Redis.DB,
		Password: configProvider.Redis.Password,
	})
	// 健康检查
	ctx := context.Background()
	re := client.Ping(ctx)
	if _, err := re.Result(); err != nil {
		log.Println("连接Redis出现错误 ", err)
		return nil
	}
	log.Println("Redis Connect Status is", re.Val())
	return client
}

func (r *MyRedis) Ping(ctx context.Context) error {
	_, err := RDBClient.Ping(ctx).Result()
	return err
}

func (r *MyRedis) SetKey(ctx context.Context, k string, v any, expireMin ...any) error {
	expiration := time.Duration(DeafultExpireWithRedis * int(time.Minute))
	if len(expireMin) > 0 {
		// 判断传入的过期时间是int还是time.duration类型
		if v, ok := expireMin[0].(int); ok {
			expiration = time.Duration(v * int(time.Minute))
		} else if v, ok := expireMin[0].(time.Duration); ok {
			expiration = v
		}
	}

	return RDBClient.Set(ctx, k, v, expiration).Err()
}

func (r *MyRedis) GetKey(ctx context.Context, k string) (string, error) {
	return RDBClient.Get(ctx, k).Result()
}

func (r *MyRedis) DelKey(ctx context.Context, k string) error {
	return RDBClient.Del(ctx, k).Err()
}

// 以秒为单位
func (r *MyRedis) CheckExpiration(ctx context.Context, k string) (time.Duration, error) {
	return RDBClient.TTL(ctx, k).Result()
}
