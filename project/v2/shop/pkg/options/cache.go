package options

type RedisOptions struct {
	Host                  string   `json:"host" mapstructure:"host"`
	Port                  int      `json:"port" mapstructure:"port"`
	Addrs                 []string `json:"addrs" mapstructure:"addrs"`
	Username              string   `json:"username" mapstructure:"username"`
	Password              string   `json:"password" mapstructure:"password"`
	Database              int      `json:"database" mapstructure:"database"`
	MaxIdle               int      `json:"max_idle" mapstructure:"max_idle"`                                 // 最大空闲连接数
	MaxActive             int      `json:"max_active" mapstructure:"max_active"`                             // 最大连接次数
	TimeOut               int      `json:"time_out" mapstructure:"time_out"`                                 // 连接超时时间
	MasterName            string   `json:"master_name" mapstructure:"master_name"`                           // 设置 Sentinel 分布式 redis 连接 主从架构
	EnableCluster         bool     `json:"enable_cluster" mapstructure:"enable_cluster"`                     // 是否是 Cluster 分布式 redis 连接 分片架构
	UseSSL                bool     `json:"use_ssl" mapstructure:"use_ssl"`                                   // 是否支持 SSL连接 (redis 需要做对应的配置)
	SSLInsecureSkipVerify bool     `json:"ssl_insecure_skip_verify" mapstructure:"ssl_insecure_skip_verify"` // 是否跳过 SSL验证
	EnableTracing         bool     `json:"enable_tracing" mapstructure:"enable_tracing"`                     // 是否开启 tracing
}
