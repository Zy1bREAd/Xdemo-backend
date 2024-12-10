package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"time"

	resp "xdemo/internal/api/response"

	"github.com/gin-gonic/gin"
	"github.com/steambap/captcha"
)

type MyCaptcha struct {
	// 空结构体，方便扩展
}

var captchaInstance *MyCaptcha

// 单例模式
func NewCaptcha() *MyCaptcha {
	if captchaInstance == nil {
		captchaInstance = &MyCaptcha{}
	}
	return captchaInstance
}

func (mycc *MyCaptcha) CreateCaptchaImage(w io.Writer) (string, error) {
	ccData, err := captcha.New(150, 50)
	if err != nil {
		log.Printf("实例化captcha验证器发生错误 %s", err)
		return "", err
	}
	// 存储Redis
	// robj := NewMyRedis()
	// rdb := robj.NewRedisClient()
	randKey := GenerateRandKey()
	// 设置1分半过期时间的验证码key-value,key使用伪随机数
	err = RedisInstance.SetKey(randKey, ccData.Text, 90*time.Second)
	if err != nil {
		log.Println("存储验证码时出现错误", err)
		return "", err
	}
	fmt.Println(randKey)
	ccData.WriteImage(w)
	return randKey, nil
}

// 刷新验证码
func RefreshCaptcha() {

}

// API接口
func RequestLoginCaptcha(ctx *gin.Context) {
	bf := bytes.NewBuffer([]byte{})
	// loginCaptcha, err := captchaInstance.CreateCaptchaImage(ctx.Writer)
	rk, err := captchaInstance.CreateCaptchaImage(bf)
	// bf2Base64 := base64.StdEncoding.EncodeToString(bf.Bytes())
	if err != nil {
		respMsg := "获取验证码失败" + err.Error()
		log.Println(respMsg)
		ctx.JSON(200, resp.FailedRespJSON(5, "Failed", respMsg, nil))
	}
	ctx.JSON(200, resp.SuccessRespJSON("Success", "获取验证码成功", map[string]any{
		"rk":             rk,
		"captcha_base64": base64.StdEncoding.EncodeToString(bf.Bytes()),
	}))
}
