package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	api "xdemo/internal/api"
	resp "xdemo/internal/api/response"
	model "xdemo/internal/database"
	global "xdemo/internal/global"
	middleware "xdemo/internal/middleware"
)

func UserLogin(ctx *gin.Context) {
	// TODO：登录逻辑(获取登录信息 -> 验证鉴权 -> 颁发token)
	// 错误处理逻辑,遇到panic统一处理返回错误json给前端
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(200, resp.FailedRespJSON(global.InternalServerError, "InternalServerError", fmt.Sprintln(err), nil))
			log.Fatalf("[Login Failed] - Internal Server Error is %s", err)
		}
	}()

	var user model.User
	ctx.ShouldBind(&user)
	inputPwd := user.Password
	// gorm获取DB记录
	userData := global.GDB.Where("account = ?", user.Account).First(&user)
	if userData.RowsAffected == 0 {
		respMsg := fmt.Sprintf("[Login Failed] - %s Account or Password is incorrect.", user.Account)
		ctx.JSON(200, resp.FailedRespJSON(global.ParamsError, "VerifyFailed", respMsg, nil))
		log.Println(respMsg)
		return
	}
	// 验证账号密码
	dbPwd := user.Password
	if !api.ComparePasswordWithHash(inputPwd, dbPwd) {
		respMsg := fmt.Sprintf("[Login Failed] - %s Account or Password is incorrect.", user.Account)
		ctx.JSON(200, resp.FailedRespJSON(global.ParamsError, "VerifyFailed", respMsg, nil))
		log.Println(respMsg)
		return
	}
	// 校验验证码
	val, err := api.RedisInstance.GetKey(ctx, user.RK)
	switch {
	case err == redis.Nil:
		log.Println("验证码Key不存在")
		respMsg := fmt.Sprintln("验证码已过期,请重新输入")
		ctx.JSON(200, resp.FailedRespJSON(global.DeafultFailed, "ValidateFailed", respMsg, nil))
		return
	case err != nil:
		log.Println("操作Redis发生内部错误")
		panic(err)
	}

	// 无论大小写都通过！
	if !strings.EqualFold(user.Captcha, val) {
		// 验证码不通过
		log.Println("登录验证码校验不一致")
		respMsg := fmt.Sprintln("输入的验证码不一致")
		ctx.JSON(200, resp.FailedRespJSON(global.DeafultFailed, "ValidateFailed", respMsg, nil))
		return
	}

	// 鉴权并授予Token，并存储redis中
	tokenStr, err := middleware.GenerateJWT(user.ID, user.Account, user.Status)
	if err != nil {
		log.Println("创建Token时发生内部错误")
		panic(err)
	}
	err = api.RedisInstance.SetKey(ctx, middleware.USER_LOGIN_TOKEN_KEY_PREFIX+strconv.Itoa(int(user.ID)), tokenStr)
	if err != nil {
		log.Println("存储User Token时发生内部错误")
		panic(err)
	}

	// 校验用户输入验证码
	// cc := NewCaptcha()
	// loginCaptcha, err := cc.CreateCaptchaImage(ctx.Writer)
	// if err != nil {
	// 	panic(err)
	// }
	// userToken := ctx.Request.Header.Get("token")
	// fmt.Println(userToken)
	// isVaild := auth.IsUserTokenVaild(userToken)
	// if !isVaild {
	// 	log.Println("User Token is not Vaild!")
	// 	ctx.JSON(200, gin.H{
	// 		"code":   403,
	// 		"status": "Success",
	// 		"msg":    fmt.Sprintf("Your Account [%s] Token is not Vaild.", user.Name),
	// 	})
	// }

	// 生成jwt时需要将user信息传入进来
	respData := map[string]any{
		"user_info": map[string]any{
			"account":       user.Account,
			"username":      user.Name,
			"status":        user.Status,
			"team":          user.Team,
			"last_login_at": user.LastLoginAt,
		},
		"token": tokenStr,
		// "captcha": loginCaptcha,
	}
	ctx.JSON(200, resp.SuccessRespJSON("Success", "Login Success", respData))
}

/*
request json body

	{
	    "username": "wangwenqiang",
	    "account": "1186405248",
	    "password": "whypassword"
	}
*/
func UserRegister(ctx *gin.Context) {
	// TODO: 注册逻辑
	var registerForm *model.RegisterUser
	ctx.ShouldBind(&registerForm)
	fmt.Println(registerForm)

	// 表单校验
	errList := middleware.RegisterValidator(registerForm)
	if len(errList) > 0 {
		log.Println(errList)
		ctx.JSON(200, resp.FailedRespJSON(global.ValidateError, "ValidateFailed", "User Register Form Validate Failed", map[string]any{
			"errList": errList,
		}))
		return
	}

	// 为注册新用户添加数据(并且保护密码加密存储)
	userPwd, err := api.EncryptWithHash(registerForm.Password)
	if err != nil {
		log.Fatal(err)
		return
	}
	now := time.Now()
	user := &model.User{
		Name:     registerForm.Name,
		Account:  registerForm.Email,
		Email:    registerForm.Email,
		Password: string(userPwd),
		Status:   "Pending",
	}
	user.CreatedAt = now
	// Gorm操作DB存储记录
	global.GDB.Create(&user)
	// 后期需要封装响应体
	respMsg := fmt.Sprintf("Account %s Register Success", user.Account)
	ctx.JSON(200, resp.SuccessRespJSON("Success", respMsg, map[string]any{
		"user":   user,
		"timeAt": time.Now(),
	}))
}

func ValidateUserInfo(ctx *gin.Context) {
	// 暂时只有校验用户账号是否重复
	var uc model.UserAccount
	ctx.ShouldBindJSON(&uc)
	validateAccount := uc.Account
	result := global.GDB.Table("xdemo_user").First(&uc, "account = ?", validateAccount)
	// log.Println(result.Statement, result.Error.Error(), result.RowsAffected)
	if result.Error != nil && result.Error.Error() != "record not found" {
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.ParamsError, "SQLQueryError", "检查是否重名出现错误", nil))
		return
	}
	if result.RowsAffected > 0 {
		respMsg := fmt.Sprintln("该邮箱账号已存在")
		ctx.JSON(http.StatusOK, resp.FailedRespJSON(resp.DeafultFailed, "ValidateFailed", respMsg, nil))
		return
	}
	ctx.JSON(http.StatusOK, resp.SuccessRespJSON("ValidateSuccess", "该账号可以使用", nil))
}
