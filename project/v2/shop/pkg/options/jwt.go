package options

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/pflag"
)

type JwtOptions struct {
	Realm      string        `json:"realm" mapstructure:"realm"`             // 区分项目用
	Method     string        `json:"method" mapstructure:"method"`           // 加密方式
	Key        string        `json:"key" mapstructure:"key"`                 // jwt密钥
	Timeout    time.Duration `json:"timeout" mapstructure:"timeout"`         // jwt 超时时间
	MaxRefresh time.Duration `json:"max_refresh" mapstructure:"max_refresh"` // jwt 刷新时间
}

func NewJwtOptions() *JwtOptions {
	return &JwtOptions{
		Realm:      "ali",
		Method:     "HS256",
		Key:        "ali",
		Timeout:    time.Duration(24) * time.Hour,
		MaxRefresh: time.Duration(24) * time.Hour,
	}
}

func (j *JwtOptions) Validate() []error {
	var errs []error

	if !govalidator.StringLength(j.Key, "6", "32") {
		errs = append(errs, fmt.Errorf("key length should be 6 or 32"))
	}
	return errs
}

func (j *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.StringVar(&j.Realm, "jwt.realm", j.Realm, "Realm name to display to the user.")
	fs.StringVar(&j.Method, "jwt.method", j.Method, "Method used to sign jwt token.")
	fs.StringVar(&j.Key, "jwt.key", j.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&j.Timeout, "jwt.timeout", j.Timeout, "JWT token timeout.")
	fs.DurationVar(&j.MaxRefresh, "jwt.max-refresh", j.MaxRefresh, "This field allows clients to refresh their token until MaxRefresh has passed.")
}
