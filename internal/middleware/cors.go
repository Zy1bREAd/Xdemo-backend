package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AllowOriginTitle  = "Access-Control-Allow-Origin"
	AllowMethodsTitle = "Access-Control-Allow-Methods"
	AllowHeadersTitle = "Access-Control-Allow-Headers"
)

type CorsConfig struct {
	AllowOrigin  []string
	AllowMethods []string
	AllowHeaders []string
}

// 解析并拼接结构体成员成string返回
func ParseConfig(member []string) string {
	var builder strings.Builder
	for k, v := range member {
		builder.WriteString(v)
		if k != len(member)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}

func (cc *CorsConfig) Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set(AllowOriginTitle, ParseConfig(cc.AllowOrigin))
		ctx.Writer.Header().Set(AllowMethodsTitle, ParseConfig(cc.AllowMethods))
		ctx.Writer.Header().Set(AllowHeadersTitle, ParseConfig(cc.AllowHeaders))
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 执行下一个handleFunc链的函数
		ctx.Next()
	}
}
