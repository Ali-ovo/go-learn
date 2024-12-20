package options

import (
	"shop/gmicro/pkg/errors"

	"github.com/spf13/pflag"
)

type RegistryOptions struct {
	Address string `mapstructure:"address" json:"address,omitempty"`
	Scheme  string `mapstructure:"scheme"  json:"scheme,omitempty"`
	//Version string `mapstructure:"version" json:"version"` // 设置版本等级
}

func NewRegistryOptions() *RegistryOptions {
	return &RegistryOptions{
		Address: "127.0.0.1:8500",
		Scheme:  "http",
	}
}

func (o *RegistryOptions) Validate() []error {
	errs := []error{}
	if o.Address == "" || o.Scheme == "" {
		errs = append(errs, errors.New("address and scheme is empty"))
	}
	return errs
}

func (o *RegistryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "consul.address", o.Address, "consul address, if left, default is 127.0.0.1:8500")
	fs.StringVar(&o.Scheme, "consul.schema", o.Scheme, "registry schema, if left, default is http")
	//fs.StringVar(&o.Version, "consul.version", o.Version, "server version default is")
}
