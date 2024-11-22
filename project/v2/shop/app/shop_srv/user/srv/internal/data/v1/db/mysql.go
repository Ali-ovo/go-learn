package db

import (
	"fmt"
	"shop/app/shop_srv/user/srv/internal/data/v1"
	"shop/gmicro/pkg/errors"
	"shop/pkg/code"
	"shop/pkg/options"
	"sync"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	dbFactory *gorm.DB
	once      sync.Once
)

func GetDBfactoryOr(myqslOpts *options.MySQLOptions) (*gorm.DB, error) {
	if myqslOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}

	var err error
	once.Do(func() {
		//dsn := "root:56248123@tcp(192.168.16.110:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			myqslOpts.Username,
			myqslOpts.Password,
			myqslOpts.Host,
			myqslOpts.Port,
			myqslOpts.Database,
		)

		dbFactory, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 禁止表名复数形式, 例如User的表名默认是users
			},
		})
		if err != nil {
			return
		}

		zap.S().Debug("[user-srv] 初始化 Mysql 完成")
		// 定义一个表结构, 将表结构直接生成对应的表 - migrations
		// 迁移 schema
		_ = dbFactory.AutoMigrate(&data.UserDO{})

		sqlDB, _ := dbFactory.DB()
		// 最大允许 连接数
		sqlDB.SetMaxOpenConns(myqslOpts.MaxOpenConnections)
		// 最大空闲 连接数
		sqlDB.SetMaxIdleConns(myqslOpts.MaxIdleConnections)
		// 设置连接重用的最大时间
		sqlDB.SetConnMaxLifetime(myqslOpts.MaxConnectionLifetime)
	})

	if dbFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil
}
