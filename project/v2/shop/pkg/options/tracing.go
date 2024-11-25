package options

import (
	"shop/gmicro/pkg/errors"

	"github.com/spf13/pflag"
)

type TelemetryOptions struct {
	Name     string  `json:"name"     mapstructure:"name"`
	Endpoint string  `json:"endpoint" mapstructure:"endpoint"`
	Sampler  float64 `json:"sampler"  mapstructure:"sampler"`
	Batcher  string  `json:"batcher"  mapstructure:"batcher"`
}

func NewTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Name:     "telemetry",
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler:  1.0,
		Batcher:  "jaeger",
	}
}

func (to *TelemetryOptions) Validate() []error {
	var errs []error
	if to.Batcher != "jaeger" && to.Batcher != "zipkin" {
		errs = append(errs, errors.New("opentelemetry batcher must be 'jaeger' or 'zipkin'"))
	}
	return errs
}

// AddFlags adds flags related to open telemetry for a specific tracing to the specified FlagSet.
func (to *TelemetryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&to.Name, "telemetry.name", to.Name, "opentelemetry name")
	fs.StringVar(&to.Endpoint, "telemetry.endpoint", to.Endpoint, "opentelemetry endpoint")
	fs.Float64Var(&to.Sampler, "telemetry.sampler", to.Sampler, "telemetry sampler")
	fs.StringVar(&to.Batcher, "telemetry.batcher", to.Batcher, "telemetry batcher, only support 'jaeger' or 'zipkin'")
}
