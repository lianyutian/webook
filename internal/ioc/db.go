package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	dsn := "cloud:li.ming9518@tcp(116.198.217.158:3306)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			LogLevel:      logger.Info,            // 日志级别
			SlowThreshold: 200 * time.Millisecond, // 慢查询阈值
			// IgnoreRecordNotFoundError: true, // 忽略未找到记录的错误
			Colorful: true, // 颜色输出
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
