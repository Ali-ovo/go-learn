package options

import (
	"fmt"
	"shop/gmicro/pkg/common/util/net"
	"shop/gmicro/pkg/host"

	"github.com/spf13/pflag"
)

type EsOptions struct {
	Host     string `json:"host"     mapstructure:"host"`
	Port     int    `json:"port"     mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

func NewEsOptions() *EsOptions {
	return &EsOptions{
		Host: "127.0.0.1",
		Port: 9200,
	}
}

func (eo *EsOptions) Validate() []error {
	errs := []error{}
	if !net.IsValidPort(eo.Port) {
		errs = append(errs, fmt.Errorf("not a valid http port: %d", eo.Port))
	}
	if !host.IsValidIP(eo.Host) {
		errs = append(errs, fmt.Errorf("not a valid ip: %s", eo.Host))
	}
	return errs
}

func (eo *EsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&eo.Host, "es.host", eo.Host, "es service host address. If left blank, the following related es options will be ignored.")
	fs.IntVar(&eo.Port, "es.port", eo.Port, "es service port. If left blank, the following related es options will be ignored.")
	fs.StringVar(&eo.Username, "es.username", eo.Username, "es service username. If left blank, the following related es options will be ignored.")
	fs.StringVar(&eo.Password, "es.password", eo.Password, "es service password. If left blank, the following related es options will be ignored.")
}
