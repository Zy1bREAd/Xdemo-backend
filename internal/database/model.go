package database

import (
	"time"

	"gorm.io/gorm"
)

type UserAccount struct {
	Account string `gorm:"account;type:varchar(64) not null;uniqueIndex" json:"account"  validate:"required"`
}

// 用户表
type User struct {
	gorm.Model
	Name string `gorm:"name;type:varchar(64); not null" json:"username" validate:"required"`
	// Account     UserAccount
	Account     string      `gorm:"column:account;type:varchar(64) not null;uniqueIndex" json:"account"  validate:"required"`
	Password    string      `gorm:"column:password;type:varchar(64) not null" json:"password"  validate:"required"`
	LastLoginAt time.Time   `gorm:"column:datetime;type:datetime(3);default:'1970-01-01 00:00:00'"`
	Status      string      `gorm:"column:status;type:varchar(64) not null"`
	Team        string      `gorm:"column:team;type:varchar(64)" json:"team"`
	Email       string      `gorm:"column:email;type:varchar(64)" json:"email"`
	Phone       string      `gorm:"column:phone;type:varchar(20);uniqueIndex" json:"phone"`
	HandleLogs  []HandleLog `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:Handler"`
	// isDeleted bool 用于记录是否被删除

	RK      string `gorm:"-" json:"rk"`
	Captcha string `gorm:"-" json:"validate_code"`
}

type ExecCmdBody struct {
	Command []string `json:"command"`
}

// 操作日志表
type HandleLog struct {
	gorm.Model
	Detail   string `gorm:"column:handle_detail;type:varchar(64);"`
	Handler  uint   `gorm:"column:handle_user;comment:User_ID;index;"`
	Type     int    `gorm:"column:handle_type;type:int(8);not null"`
	Resource string `gorm:"column:handle_resource;type:varchar(64);not null"`
}

func (hl HandleLog) TableName() string {
	return TableNamePrefix + "handle_log"
}

// 校验器
type RegisterUser struct {
	Name            string `json:"username" validate:"required"`
	Email           string `json:"email"  validate:"required,email,isRepeat"`
	Password        string `json:"password"  validate:"required,min=8,max=24,pwdLevel"`
	ConfirmPassword string `json:"confirm_password"  validate:"required,eqfield=Password,pwdLevel"`
}
