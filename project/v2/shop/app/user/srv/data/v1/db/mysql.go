package db

import (
	"fmt"
	"log"
	"os"
	code2 "shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
	"shop/pkg/options"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GetDBfactoryOr(myqslOpts *options.MySQLOptions) (*gorm.DB, error) {
	//dsn := "root:56248123@tcp(192.168.16.110:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		myqslOpts.Username,
		myqslOpts.Password,
		myqslOpts.Host,
		myqslOpts.Port,
		myqslOpts.Database,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁止表名复数形式, 例如User的表名默认是users
		},
		Logger: newLogger,
	})
	if err != nil {
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	//zap.S().Debug("[user-srv] 初始化 Mysql 完成")
	//// 定义一个表结构, 将表结构直接生成对应的表 - migrations
	//// 迁移 schema
	//_ = db.AutoMigrate(&model.User{})

	sqlDB, err := db.DB()
	// 最大允许 连接数
	sqlDB.SetMaxOpenConns(myqslOpts.MaxOpenConnections)
	// 最大空闲 连接数
	sqlDB.SetMaxIdleConns(myqslOpts.MaxIdleConnections)
	// 设置连接重用的最大时间
	sqlDB.SetConnMaxLifetime(myqslOpts.MaxConnectionLifetime)
	if err != nil {
		return nil, err
	}
	return db, nil
}
