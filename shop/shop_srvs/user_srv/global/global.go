package global

import (
	"go-learn/shop/shop_srvs/user_srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

// func init() {
// 	dsn := "root:123456@tcp(192.168.189.128:3306)/user_srv?charset=utf8mb4&parseTime=True&loc=Local"

// 	newLogger := logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
// 		logger.Config{
// 			SlowThreshold:             time.Second, // Slow SQL threshold
// 			LogLevel:                  logger.Info, // Log level
// 			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
// 			ParameterizedQueries:      true,        // Don't include params in the SQL log
// 			Colorful:                  true,        // Disable color
// 		},
// 	)

// 	var err error
// 	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: newLogger,
// 		NamingStrategy: schema.NamingStrategy{
// 			SingularTable: true,
// 		},
// 	})

// 	if err != nil {
// 		panic(err)
// 	}
// }
