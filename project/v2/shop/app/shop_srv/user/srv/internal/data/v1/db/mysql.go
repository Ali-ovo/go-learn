package db

import (
	"fmt"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/errors"
	"shop/pkg/code"
	"shop/pkg/options"
	"sync"

	"gorm.io/gorm"
)

var (
	dbFactory *gorm.DB
	once      sync.Once
)

// GetDBfactoryOr
//
//	@Description: 返回 gorm 连接 并且返回的是全局的 gorm 连接, 只初始一次, 后续直接拿到这个变量
//	@param myqslOpts
//	@return *gorm.DB
//	@return error
func GetDBfactoryOr(mysqlOpts *options.MySQLOptions) (*gorm.DB, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}

	var err error

	once.Do(func() {
		dbFactory, err = conn.NewMySQLClient((*conn.MySQLOptions)(mysqlOpts))
		if err != nil {
			return
		}
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
