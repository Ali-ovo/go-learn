package log

import (
	"fmt"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	FORMAT_CONSOLE = "console"
	OUTPUT_STD     = "stdout"
	OUTPUT_STD_ERR = "stderr"

	flagLevel = "log.level"
)

type Options struct {
	OutputPaths      []string `json:"output_path" mapstructure:"output_path"`
	ErrorOutputPaths []string `json:"error_output_path" mapstructure:"error_output_path"`
	Level            string   `json:"level" mapstructure:"level"`
	Format           string   `json:"format" mapstructure:"format"`
	Name             string   `json:"name" mapstructure:"name"`
}

type Option func(o *Options)

func NewOptions(opts ...Option) *Options {
	options := &Options{
		Level:            zapcore.InfoLevel.String(),
		Format:           FORMAT_CONSOLE,
		OutputPaths:      []string{OUTPUT_STD},
		ErrorOutputPaths: []string{OUTPUT_STD_ERR},
	}

	for _, o := range opts {
		o(options)
	}

	return options
}

func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func (p *Options) Validate() []error {
	var errs []error
	if p.Format != FORMAT_CONSOLE {
		errs = append(errs, fmt.Errorf("invalid format: %s", p.Format))
	}
	if p.Level == "" {
		errs = append(errs, fmt.Errorf("empty level"))
	}
	if p.Name == "" {
		errs = append(errs, fmt.Errorf("empty name"))
	}
	return errs
}

func (o *Options) AddFlags(fs pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Log level")
}
