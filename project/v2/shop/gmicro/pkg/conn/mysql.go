package conn

import (
	"fmt"
	"shop/gmicro/pkg/log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type MySQLOptions struct {
	Host                  string        `mapstructure:"host"                   json:"host,omitempty"`
	Port                  int           `mapstructure:"port"                   json:"port,omitempty"`
	Username              string        `mapstructure:"username"               json:"username,omitempty"`
	Password              string        `mapstructure:"password"               json:"password,omitempty"`
	Database              string        `mapstructure:"database"               json:"database"`
	MaxIdleConnections    int           `mapstructure:"max_idle_connections"   json:"max_idle_connections,omitempty"`
	MaxOpenConnections    int           `mapstructure:"max_open_connections"   json:"max_open_connections,omitempty"`
	MaxConnectionLifetime time.Duration `mapstructure:"max_conection_lifetime" json:"max_conection_lifetime,omitempty"`
	EnableLog             bool          `mapstructure:"enable_log"             json:"enable_log"`
	LogLevel              int           `mapstructure:"log_level"              json:"log_level,omitempty"`
	SlowThreshold         time.Duration `mapstructure:"slow-threshold"         json:"slow-threshold"`
	EnableColorful        bool          `mapstructure:"enable-colorful"        json:"enable-color"`
}

func NewMySQLClient(opts *MySQLOptions) (*gorm.DB, error) {
	var newLogger logger.Interface

	//dsn := "root:56248123@tcp(192.168.16.110:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
	)

	if opts.EnableLog {
		newLogger = logger.New(
			log.StdInfoLogger(),
			logger.Config{
				SlowThreshold: opts.SlowThreshold,             // 慢 SQL 阈值
				LogLevel:      logger.LogLevel(opts.LogLevel), // Log level
				Colorful:      opts.EnableColorful,            // 禁用彩色打印
			},
		)
	}

	dbFactory, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁止表名复数形式, 例如User的表名默认是users
		},
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := dbFactory.DB()
	// 最大允许 连接数
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	// 最大空闲 连接数
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)
	// 设置连接重用的最大时间
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifetime)

	return dbFactory, nil
}
