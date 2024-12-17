package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

const TableNamePrefix = "xdemo_"

type InitTable struct {
	UpdateAt time.Time `gorm:"datetime;type:datetime(6)"`
	Version  string    `gorm:"primarykey;type:varchar(64)"`
}

// 自定义表名（设置后将不会使用表前缀）
func (it InitTable) TableName() string {
	return TableNamePrefix + "init_version"
}

// 初始化数据库操作（建立表等）
func InitDBTable(db *gorm.DB) error {
	// 先判断是否需要初始化，其次判断是否连接正常，最后进行初始化
	err := db.AutoMigrate(&InitTable{}, &User{}, &HandleLog{}, &Job{}) // 后期只能第一次进行automigrate，否则影响线上生产环境
	if err != nil {
		return fmt.Errorf("init DB Error: %s", err)
	}
	// layout := "2009-09-09 14:10:30"
	//通过判断是否有记录来决定是否需要添加记录
	version := "v0.0.1"
	db.FirstOrCreate(&InitTable{
		Version:  version,
		UpdateAt: time.Now(),
	})
	var lastVersion InitTable
	resultObj := db.Last(&lastVersion)
	if resultObj.RowsAffected > 0 && resultObj.Error != nil {
		log.Println("XDemo Version is ", lastVersion.Version)
	}
	return nil
}
