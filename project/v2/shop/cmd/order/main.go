package main

import (
	"fmt"
	"go-learn/project/v2/shop/pkg/app"
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

var _ app.CliOptions = &Config{}

func main() {
	cfg := &Config{
		Log: log.NewOptions(),
	}
	appl := app.NewApp("order", "shop",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		app.WithNoConfig(),
	)

	appl.Run()

}

func run(cfg *Config) app.RunFunc {
	return func(basename string) error {
		fmt.Println(cfg.Log.Level)
		return nil
	}
}
