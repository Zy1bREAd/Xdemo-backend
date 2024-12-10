package global

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ！定义全局变量

var GDB *gorm.DB // Gorm的DB操作对象
var VA *validator.Validate
