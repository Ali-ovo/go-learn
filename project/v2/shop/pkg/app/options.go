package app

import (
	cliflag "go-learn/project/v2/shop/pkg/common/cli/flag"
)

// CliOptions abstracts configuration options for reading parameters from the
// command line.
type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets) // 生成一个被 cliflag.NamedFlagSets 封装后的日志对象
	Validate() []error                  // 添加你想要的验证函数
	// AddFlags 将标志添加到指定的 FlagSet 对象中。
	// AddFlags(fs *pflag.FlagSet)
}

// ConfigurableOptions abstracts configuration options for reading parameters
// from a configuration file.
type ConfigurableOptions interface {
	// ApplyFlags parsing parameters from the command line or configuration file
	// to the options instance.
	ApplyFlags() []error
}

// CompleteableOptions abstracts options which can be completed.
type CompleteableOptions interface {
	Complete() error
}

// PrintableOptions abstracts options which can be printed.
type PrintableOptions interface {
	String() string
}
