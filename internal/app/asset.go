package app

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Asset struct {
	ID       int32
	Name     string
	Arena    string
	CreateAt time.Time
	CreateBy string
	UpdateAt time.Time
	UpdateBy string
	Status   string
	IsPublic bool
	Type     string
}

// 操作资产
// map[string]any中存放选项，比如需要操作的一些数据体
/*
e.g. => {
	data: {
		xxx:"xxx",
		xx: 1234,
	},
	status: "xxx",
	xxx:xxxx
}
*/
func HandleAsset(db *gorm.DB, handle string, opts ...map[string]any) {
	db.AutoMigrate(&Asset{})
	switch strings.ToLower(handle) {
	case "add":
		AddAsset(db, opts...)
	case "update":
		fmt.Println("222")
	}
}

// 增删改查
func AddAsset(db *gorm.DB, opts ...map[string]any) {
	// db.Create()
}

// 增删改查
func ListAsset(db *gorm.DB, opts ...map[string]any) {
	fmt.Println("22")
}
