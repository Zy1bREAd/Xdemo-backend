package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	api "xdemo/internal/api"
	resp "xdemo/internal/api/response"

	"github.com/gin-gonic/gin"
)

const (
	USER_LOGIN_TOKEN_KEY_PREFIX string = "USER_LOGIN_TOKEN_"
	TokenHeader                 string = "Authorization"
	TokenPrefix                 string = "Bearer: "

	ExpirationSecond int = 60
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从Header获取Token，验证token是否存在，解析token是否正确
		headerValue := ctx.GetHeader(TokenHeader)
		if headerValue == "" || !strings.HasPrefix(headerValue, TokenPrefix) {
			respMsg := "Token is Malformed"
			ctx.JSON(http.StatusUnauthorized, resp.FailedRespJSON(resp.Malformed, "TokenMalformed", respMsg, nil))
			ctx.Abort()
			return
		}
		reqToken := headerValue[len(TokenPrefix):]
		tokenClaim, err := ParseJWT(reqToken, SECRETKEY, &UserClaim{})
		if err != nil {
			respMsg := fmt.Sprintln("鉴权失败，验证用户Token出错", err)
			ctx.JSON(http.StatusUnauthorized, resp.FailedRespJSON(resp.TokenError, "Unauthorized", respMsg, nil))
			ctx.Abort()
			return
		}
		ri := api.NewMyRedis()
		userID := strconv.FormatUint(uint64(tokenClaim.UserID), 10)
		tokenKeyName := USER_LOGIN_TOKEN_KEY_PREFIX + userID
		userToken, err := ri.GetKey(tokenKeyName)
		log.Println(USER_LOGIN_TOKEN_KEY_PREFIX + userID)
		if err != nil {
			respMsg := "鉴权失败，不存在该用户的Token"
			ctx.JSON(http.StatusUnauthorized, resp.FailedRespJSON(resp.TokenNotFound, "Unauthorized", respMsg, nil))
			ctx.Abort()
			return
		}
		// fmt.Println(reqToken, userToken)
		if reqToken != userToken {
			respMsg := "鉴权失败，Token不一致"
			ctx.JSON(http.StatusUnauthorized, resp.FailedRespJSON(resp.TokenError, "Handle Unauthorized", respMsg, nil))
			ctx.Abort()
			return
		}
		// 检查过期时间以及给予续签
		tokenExpiration, err := ri.CheckExpiration(tokenKeyName)
		// fmt.Println(tokenExpiration.Nanoseconds())
		if err != nil || (int(tokenExpiration.Seconds()) <= 0 && tokenExpiration.Nanoseconds() != -1) {
			respMsg := "鉴权失败，Token已过期或不存在"
			ctx.JSON(http.StatusUnauthorized, resp.FailedRespJSON(resp.TokenExpired, "Handle Unauthorized", respMsg, nil))
			ctx.Abort()
			return
		} else if int(tokenExpiration.Seconds()) < ExpirationSecond && int(tokenExpiration.Seconds()) > 0 {
			// 小于一分钟，将自动重新续签
			newToken, err := GenerateJWT(tokenClaim.UserID, tokenClaim.Account, tokenClaim.Status)
			if err != nil {
				respMsg := "续签新Token失败"
				ctx.JSON(http.StatusUnauthorized, resp.FailedRespJSON(resp.TokenRenewError, "Handle Unauthorized", respMsg, nil))
				ctx.Abort()
				return
			}
			ri.SetKey(tokenKeyName, newToken)
			ctx.Header(TokenHeader, TokenPrefix+newToken)
			log.Println("续签Token成功")
		}
		// 将认证后的信息存储context
		ctx.Set("userInfo", map[string]string{
			"id":      userID,
			"account": tokenClaim.Account,
		})
		ctx.Next()

	}
}
