package main

import (
	"crypto/md5"
	"encoding/hex"
	"go-learn/shop/shop_srvs/inventory_srv/model"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:123456@tcp(192.168.189.128:3306)/inventory_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&model.Inventory{},
		&model.StockSellDetail{},
	)

	// orderDetail := model.StockSellDetail{
	// 	OrderSn: "ali-123",
	// 	Status:  1,
	// 	Detail:  []model.GoodsDetail{{1, 2}, {2, 4}},
	// }

	// db.Create(&orderDetail)

	// var sellDetail model.StockSellDetail
	// db.Where(model.StockSellDetail{OrderSn: "ali-123"}).First(&sellDetail)
	// fmt.Println(sellDetail)
}
