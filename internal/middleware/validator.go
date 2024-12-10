package middleware

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	model "xdemo/internal/database"
	global "xdemo/internal/global"

	"github.com/go-playground/validator/v10"
)

func InitGlobalValidator() {
	global.VA = validator.New()
	// 注册自定义校验函数
	global.VA.RegisterValidation("pwdLevel", ValidateUserPasswordLevel)
	global.VA.RegisterValidation("isRepeat", ValidateIsAccountRepeat)
}

// 自定义校验规则
func ValidateUserPasswordLevel(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	hasLower := regexp.MustCompile(`[a-z].{1,}`).MatchString(pwd)
	hasUpper := regexp.MustCompile(`[A-Z].{1,}`).MatchString(pwd)
	hasNumber := regexp.MustCompile(`[0-9].{1,}`).MatchString(pwd)
	return hasLower && hasUpper && hasNumber
}

// 查询用户账号是否有重复的记录
func ValidateIsAccountRepeat(fl validator.FieldLevel) bool {
	queryData := fl.Field().String()
	result := global.GDB.First(&model.User{}, "account = ?", queryData)
	log.Println(result.RowsAffected, result.Error, reflect.TypeOf(result.Error))
	if result.Error != nil && result.Error.Error() != "record not found" {
		log.Println("查询数据时发生错误", result.Error)
		return false
	} else if result.RowsAffected > 0 {
		return false
	}
	return true
}

// 判断字段校验Error来返回消息
func ValidatorFailedMessage(fe validator.FieldError) string {
	logMsg := ""
	switch fe.Tag() {
	case "email":
		logMsg = fmt.Sprintf("%s=%s 邮箱格式不正确", fe.Field(), fe.Value())
	case "eqfield":
		logMsg = fmt.Sprintf("%s=%s 两次输入的密码不一致", fe.Field(), fe.Value())
	case "min":
		logMsg = fmt.Sprintf("%s 密码长度少于%s位", fe.Field(), fe.Param())
	case "max":
		logMsg = fmt.Sprintf("%s 密码长度多于%s位", fe.Field(), fe.Param())
	case "required":
		logMsg = fmt.Sprintf("%s 不能为空", fe.Field())
	case "pwdLevel":
		logMsg = fmt.Sprintf("%s 密码至少有1个大写字母、小写字母和数字", fe.Field())
	case "isRepeat":
		logMsg = fmt.Sprintf("该账号 %s 已存在", fe.Field())
	default:
		logMsg = fe.Error()
	}
	return logMsg
}

// 用户注册验证
func RegisterValidator(ru *model.RegisterUser) []string {
	err := global.VA.Struct(ru)
	if err != nil {
		errList := []string{}
		// 自定义验证错误消息
		validErr := err.(validator.ValidationErrors)
		for _, v := range validErr {
			errMsg := ValidatorFailedMessage(v)
			log.Println(errMsg)
			errList = append(errList, errMsg)
		}
		return errList
	}
	return nil
}
