package api

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandKey() string {
	// 使用本地时间戳作为随机种子进行生成伪随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.FormatInt(r.Int63(), 10)
}
