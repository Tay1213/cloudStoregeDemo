package models

import (
	"cloudStoregeDemo/pkg/setting"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

var db *gorm.DB

type Model struct {
	IsDelete int       `json:"is_delete"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
	Atime    time.Time `json:"atime"`
}

func Setup() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))
	if err != nil {
		fmt.Printf("数据库连接出错了！: %#v", err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}
