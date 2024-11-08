package config

import (
	cliflag "go-learn/project/v2/shop/pkg/common/cli/flag"
	"go-learn/project/v2/shop/pkg/log"
)

type Config struct {
	Log *log.Options `json:"log" mapstructure:"log"`
}

func (c *Config) Validate() []error {
	var errors []error
	errors = append(errors, c.Log.Validate()...)
	return errors
}

func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	return fss
}

func New() *Config {
	return &Config{
		Log: log.NewOptions(),
	}
}
