package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DataBase *gorm.DB

// 连接 MySQL 数据库并自动创建视频表
func init() {
	dsn := "root:123456@tcp(192.168.1.163:3306)/video_process?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DataBase, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("myaql连接成功")
	sqlDB, err := DataBase.DB()
	if err != nil {
		panic(err)
	}

	// 最大连接数为 10
	sqlDB.SetMaxOpenConns(10)

	// 最大空闲连接数为 5
	sqlDB.SetMaxIdleConns(5)

	// 设置连接可复用的最长时间为 5 分钟
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
