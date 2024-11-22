package options

import (
	"fmt"
	"shop/gmicro/pkg/common/util/net"
	"shop/gmicro/pkg/host"

	"github.com/spf13/pflag"
)

type RedisOptions struct {
	Host                  string   `json:"host"                     mapstructure:"host"`
	Port                  int      `json:"port"                     mapstructure:"port"`
	Addrs                 []string `json:"addrs"                    mapstructure:"addrs"`
	Username              string   `json:"username"                 mapstructure:"username"`
	Password              string   `json:"password"                 mapstructure:"password"`
	Database              int      `json:"database"                 mapstructure:"database"`
	MaxIdle               int      `json:"max_idle"                 mapstructure:"max_idle"`                 // 最大空闲连接数
	MaxActive             int      `json:"max_active"               mapstructure:"max_active"`               // 最大连接次数
	TimeOut               int      `json:"time_out"                 mapstructure:"time_out"`                 // 连接超时时间
	MasterName            string   `json:"master_name"              mapstructure:"master_name"`              // 设置 Sentinel 分布式 redis 连接 主从架构
	EnableCluster         bool     `json:"enable_cluster"           mapstructure:"enable_cluster"`           // 是否是 Cluster 分布式 redis 连接 分片架构
	UseSSL                bool     `json:"use_ssl"                  mapstructure:"use_ssl"`                  // 是否支持 SSL连接 (redis 需要做对应的配置)
	SSLInsecureSkipVerify bool     `json:"ssl_insecure_skip_verify" mapstructure:"ssl_insecure_skip_verify"` // 是否跳过 SSL验证
	EnableTracing         bool     `json:"enable_tracing"           mapstructure:"enable_tracing"`           // 是否开启 tracing
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Host:                  "127.0.0.1",
		Port:                  6379,
		Addrs:                 []string{},
		Username:              "",
		Password:              "",
		Database:              0,
		MaxIdle:               2000,
		MaxActive:             4000,
		TimeOut:               0,
		MasterName:            "",
		EnableCluster:         false,
		UseSSL:                false,
		SSLInsecureSkipVerify: false,
		EnableTracing:         false,
	}
}

func (ro *RedisOptions) Validate() []error {
	errs := []error{}
	if !net.IsValidPort(ro.Port) {
		errs = append(errs, fmt.Errorf("not a valid http port: %d", ro.Port))
	}
	if !host.IsValidIP(ro.Host) {
		errs = append(errs, fmt.Errorf("not a valid ip: %s", ro.Host))
	}
	return errs
}

// AddFlags adds flags related to redis storage for a specific APIServer to the specified FlagSet.
func (ro *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&ro.Host, "redis.host", ro.Host, "Hostname of your Redis server.")
	fs.IntVar(&ro.Port, "redis.port", ro.Port, "The port the Redis server is listening on.")
	fs.StringSliceVar(&ro.Addrs, "redis.addrs", ro.Addrs, "A set of redis address(format: 127.0.0.1:6379).")
	fs.StringVar(&ro.Username, "redis.username", ro.Username, "Username for access to redis service.")
	fs.StringVar(&ro.Password, "redis.password", ro.Password, "Optional auth password for Redis db.")

	fs.IntVar(&ro.Database, "redis.database", ro.Database, ""+
		"By default, the database is 0. Setting the database is not supported with redis cluster. "+
		"As such, if you have --redis.enable-cluster=true, then this value should be omitted or explicitly set to 0.")

	fs.StringVar(&ro.MasterName, "redis.master-name", ro.MasterName, "The name of master redis instance.")

	fs.IntVar(&ro.MaxIdle, "redis.optimisation-max-idle", ro.MaxIdle, ""+
		"This setting will configure how many connections are maintained in the pool when idle (no traffic). "+
		"Set the --redis.optimisation-max-active to something large, we usually leave it at around 2000 for "+
		"HA deployments.")

	fs.IntVar(&ro.MaxActive, "redis.optimisation-max-active", ro.MaxActive, ""+
		"In order to not over commit connections to the Redis server, we may limit the total "+
		"number of active connections to Redis. We recommend for production use to set this to around 4000.")

	fs.IntVar(&ro.TimeOut, "redis.timeout", ro.TimeOut, "Timeout (in seconds) when connecting to redis service.")

	fs.BoolVar(&ro.EnableCluster, "redis.enable-cluster", ro.EnableCluster, ""+
		"If you are using Redis cluster, enable it here to enable the slots mode.")

	fs.BoolVar(&ro.UseSSL, "redis.use-ssl", ro.UseSSL, ""+
		"If set, IAM will assume the connection to Redis is encrypted. "+
		"(use with Redis providers that support in-transit encryption).")

	fs.BoolVar(&ro.SSLInsecureSkipVerify, "redis.ssl-insecure-skip-verify", ro.SSLInsecureSkipVerify, ""+
		"Allows usage of self-signed certificates when connecting to an encrypted Redis database.")
}
