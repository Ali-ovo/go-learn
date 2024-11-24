package db

import (
	"shop/pkg/options"
	"testing"
	"time"
)

// TestGetDBfactoryOr 测试 mysql 连接  和 创建表
//
//	@Description:
//	@param t
func TestGetDBfactoryOr(t *testing.T) {
	mysqlOpts := &options.MySQLOptions{
		Host:                  "192.168.101.49",
		Port:                  3306,
		Username:              "root",
		Password:              "56248123",
		Database:              "demo",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Second * time.Duration(10),
		LogLevel:              4,
	}
	_, err := GetDBfactoryOr(mysqlOpts)
	if err != nil {
		panic(err)
	}
}
