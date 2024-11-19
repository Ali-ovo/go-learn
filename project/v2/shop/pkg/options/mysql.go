package options

import (
	"fmt"
	"shop/gmicro/pkg/common/util/net"
	"shop/gmicro/pkg/host"
	"time"

	"github.com/spf13/pflag"
)

type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"`
	Port                  int           `json:"port,omitempty"`
	Username              string        `json:"username,omitempty"`
	Password              string        `json:"password,omitempty"`
	Database              string        `json:"database"`
	MaxIdleConnections    int           `json:"max_idle_connections,omitempty"`
	MaxOpenConnections    int           `json:"max_open_connections,omitempty"`
	MaxConnectionLifetime time.Duration `json:"max_conection_lifetime,omitempty"`
	LogLevel              int           `json:"log_level,omitempty"`
}

// NewMySQLOptions
//
//	@Description: create a `zero` value instance.
//	@return *MySQLOptions
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1",
		Port:                  3306,
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Duration(10) * time.Second,
		LogLevel:              1,
	}
}

// Validate
//
//	@Description: 校验输入是否正确
//	@receiver o
//	@return []error
func (mo *MySQLOptions) Validate() []error {
	errs := []error{}
	if net.IsValidPort(mo.Port) {
		errs = append(errs, fmt.Errorf("not a valid http port: %d", mo.Port))
	}
	if host.IsValidIP(mo.Host) {
		errs = append(errs, fmt.Errorf("not a valid ip: %s", mo.Host))
	}
	return errs
}

func (mo *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&mo.Host, "mysql.host", mo.Host, "MySQL service host address. If left blank, the following related mysql options will be ignored.")
	fs.IntVar(&mo.Port, "mysql.port", mo.Port, "MySQL service port")
	fs.StringVar(&mo.Username, "mysql.username", mo.Username, "Username for access to mysql service.")
	fs.StringVar(&mo.Password, "mysql.password", mo.Password, "Password for access to mysql, should be used pair with password.")
	fs.StringVar(&mo.Database, "mysql.database", mo.Database, "Database name for the server to use.")
	fs.IntVar(&mo.MaxIdleConnections, "mysql.max-idle-connections", mo.MaxOpenConnections, "Maximum idle connections allowed to connect to mysql.")
	fs.IntVar(&mo.MaxOpenConnections, "mysql.max-open-connections", mo.MaxOpenConnections, "Maximum open connections allowed to connect to mysql.")
	fs.DurationVar(&mo.MaxConnectionLifetime, "mysql.max-connection-life-time", mo.MaxConnectionLifetime, "Maximum connection life time allowed to connecto to mysql.")
	fs.IntVar(&mo.LogLevel, "mysql.log-mode", mo.LogLevel, "Specify gorm log level.")
}
