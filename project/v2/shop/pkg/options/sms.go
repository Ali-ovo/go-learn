package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

type SmsOptions struct {
	APIKey      string `mapstructure:"key"       json:"key"`
	APISecret   string `mapstructure:"secret"    json:"secret"`
	APISignName string `mapstructure:"sign-name" json:"sign-name"`
	APICode     string `mapstructure:"code"      json:"code"`
}

func NewSmsOptions() *SmsOptions {
	return &SmsOptions{}
}

func (so *SmsOptions) Validate() []error {
	var errs []error
	if so.APIKey == "" {
		errs = append(errs, fmt.Errorf("APIKey is required"))
	}
	if so.APISecret == "" {
		errs = append(errs, fmt.Errorf("APISecret is required"))
	}
	if so.APICode == "" {
		errs = append(errs, fmt.Errorf("APICode is required"))
	}
	if so.APISignName == "" {
		errs = append(errs, fmt.Errorf("APISignName is required"))
	}
	return errs
}

func (so *SmsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&so.APIKey, "sms.apikey", so.APIKey, "sms API Key")
	fs.StringVar(&so.APISecret, "sms.secret", so.APISecret, "sms API Secret")
	fs.StringVar(&so.APICode, "sms.code", so.APICode, "sms API Code")
	fs.StringVar(&so.APISignName, "sms.sign-name", so.APISignName, "sms API Sign Name")
}
