package db

import (
	"fmt"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	"shop/pkg/code"
	"shop/pkg/options"

	"gorm.io/gorm"
)

func NewOrderSQLClient(mysqlOpts *options.MySQLOptions) (*gorm.DB, error) {
	if mysqlOpts == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}

	msqDB, err := conn.NewMySQLClient((*conn.MySQLOptions)(mysqlOpts))
	if err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}

	//// 定义一个表结构, 将表结构直接生成对应的表 - migrations
	//// 迁移 schema
	//_ = msqDB.AutoMigrate(
	//	&do.OrderInfoDO{},
	//	&do.OrderGoods{},
	//	&do.ShoppingCartDO{},
	//)

	log.Info("[order-srv] 初始化 Mysql 完成")
	return msqDB, nil
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

	// 排序 如: age desc, name 等  age 为降序, name 为升序
	for _, v := range orderBy {
		db = db.Order(v)
	}
	db.Count(&count)
	return db.Offset(offset).Limit(limit), count
}
