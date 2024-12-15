package api

type QueueProvider interface {
	InitQueue()
}

type QueueForRedis struct {
	RedisClient MyRedis
}

func InitQueue() {

}
