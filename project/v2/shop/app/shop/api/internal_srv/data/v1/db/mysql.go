package db

import (
	"fmt"
	"shop/app/shop/api/internal_srv/data/v1"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	"shop/pkg/code"
	"shop/pkg/options"
	"sync"

	"gorm.io/gorm"
)

var (
	dbFactory data.DataFactory
	once      sync.Once
)

type mysqlFactory struct {
	db *gorm.DB
}

func (mf *mysqlFactory) Goods() data.GoodsStore {
	return newGoods(mf)
}

func (mf *mysqlFactory) Category() data.CategoryStore {
	return newCategory(mf)
}

func (mf *mysqlFactory) Brands() data.BrandsStore {
	return newBrand(mf)
}

func (mf *mysqlFactory) Banner() data.BannerStore {
	return newBanner(mf)
}

func (mf *mysqlFactory) CategoryBrands() data.CategoryBrandStore {
	return newCategoryBrands(mf)
}

func (mf *mysqlFactory) Begin() *gorm.DB {
	return mf.db.Begin()
}

var _ data.DataFactory = (*mysqlFactory)(nil)

// GetDBfactoryOr
//
//	@Description: 返回 gorm 连接 并且返回的是全局的 gorm 连接, 只初始一次, 后续直接拿到这个变量
//	@param myqslOpts
//	@return *gorm.DB
//	@return error
func GetDBfactoryOr(mysqlOpts *options.MySQLOptions) (data.DataFactory, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}

	var err error

	once.Do(func() {
		msqDB, err := conn.NewMySQLClient((*conn.MySQLOptions)(mysqlOpts))
		if err != nil {
			return
		}

		//// 定义一个表结构, 将表结构直接生成对应的表 - migrations
		//// 迁移 schema
		//_ = msqDB.AutoMigrate(
		//	&do.CategoryDO{},
		//	&do.BrandsDO{},
		//	&do.CategoryBrandDO{},
		//	&do.BannerDO{},
		//	&do.GoodsDO{},
		//)

		dbFactory = &mysqlFactory{
			db: msqDB,
		}
		log.Info("[goods-srv] 初始化 Mysql 完成")
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

	// 排序 如: age desc, name 等  age 为降序, name 为升序
	for _, v := range orderBy {
		db = db.Order(v)
	}
	db.Count(&count)
	return db.Offset(offset).Limit(limit), count
}
