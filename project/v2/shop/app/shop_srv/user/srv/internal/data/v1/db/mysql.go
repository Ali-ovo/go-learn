package db

import (
	"fmt"
	"shop/gmicro/pkg/errors"
	"shop/pkg/code"
	"shop/pkg/options"
	"sync"
	"time"

	"gorm.io/gorm/logger"

	"shop/gmicro/pkg/log"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	dbFactory *gorm.DB
	once      sync.Once
)

func GetDBfactoryOr(mysqlOpts *options.MySQLOptions) (*gorm.DB, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}

	var err error
	var newLogger logger.Interface
	once.Do(func() {
		//dsn := "root:56248123@tcp(192.168.16.110:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database,
		)

		if mysqlOpts.EnableLog {
			newLogger = logger.New(
				log.StdInfoLogger(),
				logger.Config{
					SlowThreshold: time.Second, // 慢 SQL 阈值
					LogLevel:      logger.Info, // Log level
					Colorful:      true,        // 禁用彩色打印
				},
			)
		}

		dbFactory, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 禁止表名复数形式, 例如User的表名默认是users
			},
			Logger: newLogger,
		})
		if err != nil {
			return
		}

		zap.S().Debug("[user-srv] 初始化 Mysql 完成")
		// 定义一个表结构, 将表结构直接生成对应的表 - migrations
		// 迁移 schema
		// _ = dbFactory.AutoMigrate(&data.UserDO{})

		sqlDB, _ := dbFactory.DB()
		// 最大允许 连接数
		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		// 最大空闲 连接数
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		// 设置连接重用的最大时间
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)
	})

	if dbFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil
}

// paginate 处理分页逻辑 先排序后返回分页后的查询结果
//
//	@Description:
//	@param db
//	@param page
//	@param pageSize
//	@param orderBy
//	@param conditions: 添加额外自己的逻辑
//	@return *gorm.DB
//	@return int
//	@return error
func paginate(db *gorm.DB, page int, pageSize int, orderBy []string, conditions ...func(*gorm.DB) *gorm.DB) (*gorm.DB, int64) {
	var (
		count  int64
		limit  int
		offset int
	)
	db = db.Scopes(conditions...)

	// 分页逻辑
	if pageSize == 0 {
		limit = 10
	} else {
		limit = pageSize
	}
	if page > 0 {
		offset = (page - 1) * limit
	}

	// 排序
	for _, v := range orderBy {
		db = db.Order(v)
	}
	db.Count(&count)
	return db.Offset(offset).Limit(limit), count
}
