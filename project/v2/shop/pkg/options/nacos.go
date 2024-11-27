package options

import (
	"fmt"
	"shop/gmicro/pkg/common/util/net"
	"shop/gmicro/pkg/host"

	"github.com/spf13/pflag"
)

type NacosOptions struct {
	Host      string `mapstructure:"host"     json:"host,omitempty"`
	Port      int    `mapstructure:"port"     json:"port,omitempty"`
	NameSpace string `mapstructure:"name"     json:"name,omitempty"`
	User      string `mapstructure:"user"     json:"user,omitempty"`
	Password  string `mapstructure:"password" json:"password,omitempty"`
	DataID    string `mapstructure:"data-id"  json:"data-id,omitempty"`
	Group     string `mapstructure:"group"    json:"group,omitempty"`
}

func NewNacosOptions() *NacosOptions {
	return &NacosOptions{
		Host:      "127.0.0.1",
		Port:      8848,
		NameSpace: "public",
		User:      "nacos",
		Password:  "nacos",
		DataID:    "flow",
		Group:     "sentinel-go",
	}
}

func (no *NacosOptions) Validate() []error {
	var errs []error
	if !net.IsValidPort(no.Port) {
		errs = append(errs, fmt.Errorf("not a valid http port: %d", no.Port))
	}
	if !host.IsValidIP(no.Host) {
		errs = append(errs, fmt.Errorf("not a valid ip: %s", no.Host))
	}
	return errs
}

func (no *NacosOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&no.Host, "nacos.host", no.Host, "nacos host")
	fs.IntVar(&no.Port, "nacos.port", no.Port, "nacos port")
	fs.StringVar(&no.NameSpace, "nacos.namespace", no.NameSpace, "nacos namespace")
	fs.StringVar(&no.User, "nacos.user", no.User, "nacos user")
	fs.StringVar(&no.Password, "nacos.password", no.Password, "nacos password")
	fs.StringVar(&no.DataID, "nacos.data-id", no.DataID, "nacos data-id")
	fs.StringVar(&no.Group, "nacos.group", no.Group, "nacos group")
}
