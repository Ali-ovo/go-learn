package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"strings"
)

type MapValue map[string]string

func (m *MapValue) Type() string {
	return "map[string]string"
}

func (m *MapValue) String() string {
	// 实现 String 方法，返回自定义类型的字符串表示
	// 根据您的实际需求进行实现
	return fmt.Sprintf("%v", *m)
}

func (m *MapValue) Set(value string) error {
	// 实现 Set 方法，用于解析和设置命令行标志的值到自定义类型
	// 根据您的实际需求进行实现
	// 示例中，我们将以逗号分隔的键值对字符串解析为 map[string]string
	pairs := strings.Split(value, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) == 2 {
			(*m)[kv[0]] = kv[1]
		}
	}
	return nil
}

type DtmOptions struct {
	GrpcServer string   `mapstructure:"grpc" json:"grpc"`
	HttpServer string   `mapstructure:"http" json:"http"`
	AccessPath MapValue `mapstructure:"access" json:"access"`
}

func NewDtmOptions() *DtmOptions {
	return &DtmOptions{
		GrpcServer: "127.0.0.1:36790",
		HttpServer: "http://127.0.0.1:36789/api/dtmsvr",
	}
}

func (do *DtmOptions) Validate() []error {
	errs := []error{}
	return errs
}

func (do *DtmOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&do.GrpcServer, "dtm.grpc", do.GrpcServer, "dtm grpc server. If left blank, the following related Dtm options will be ignored.")
	fs.StringVar(&do.HttpServer, "dtm.http", do.HttpServer, "dtm http server. If left blank, the following related Dtm options will be ignored.")
	fs.VarP(&do.AccessPath, "dtm.access", "", "dtm access path.")
}
