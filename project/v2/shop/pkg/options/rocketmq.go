package options

import (
	"fmt"
	"shop/gmicro/pkg/common/util/net"
	"shop/gmicro/pkg/host"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type RocketmqOptions struct {
	Addr      []string `mapstructure:"addr"       json:"addr,omitempty"`
	GroupName string   `mapstructure:"group-name" json:"group-name,omitempty"`
	Retry     int      `mapstructure:"retry"      json:"retry,omitempty"`
}

func NewRocketmqOptions() *RocketmqOptions {
	return &RocketmqOptions{
		Addr:      []string{"127.0.0.1:9876"},
		GroupName: "default",
	}
}

func (ro *RocketmqOptions) Validate() []error {
	errs := []error{}

	for _, addr := range ro.Addr {
		adr := strings.Split(addr, ":")
		port, err := strconv.Atoi(adr[1])
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid addr: %s", addr))
			break
		}
		if !net.IsValidPort(port) {
			errs = append(errs, fmt.Errorf("not a valid http port: %d", port))
		}
		if !host.IsValidIP(adr[0]) {
			errs = append(errs, fmt.Errorf("not a valid ip: %s", adr[0]))
		}
	}
	return errs
}

func (ro *RocketmqOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringArrayVar(&ro.Addr, "rocketmq.addr", ro.Addr, "RocketMQ server address")
	fs.StringVar(&ro.GroupName, "rocketmq.group-name", ro.GroupName, "RocketMQ group name")
	fs.IntVar(&ro.Retry, "rocketmq.retry", ro.Retry, "RocketMQ retry")
}
