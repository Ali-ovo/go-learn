package flag

import (
	goflag "flag"
	"go-learn/project/v2/shop/pkg/log"
	"strings"

	"github.com/spf13/pflag"
)

// WordSepNormalizeFunc changes all flags that contain "_" separators.
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

// WarnWordSepNormalizeFunc changes and warns for flags that contain "_" separators.
func WarnWordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		nname := strings.ReplaceAll(name, "_", "-")
		log.Warnf("%s is DEPRECATED and will be removed in a future version. Use %s instead.", name, nname)

		return pflag.NormalizedName(nname)
	}
	return pflag.NormalizedName(name)
}

// InitFlags normalizes, parses, then logs the command line flags.
func InitFlags(flags *pflag.FlagSet) {
	// SetNormalizeFunc 方法为命令的所有标志设置一个标准化函数。
	// 标准化函数可以用来规范化标志的名称和值。
	// WordSepNormalizeFunc 是一个由 Cobra 库提供的函数，它将标志名称中的短横线转换为下划线
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	// AddGoFlagSet 方法将 Go 标准库 flag 包的标志集添加到 Cobra 命令对象中。
	// goflag.CommandLine 是一个 *flag.FlagSet 类型的对象，它表示命令行参数和标志集。
	// 通过将 goflag.CommandLine 添加到 Cobra 命令对象中，可以让 Cobra 库支持 Go 标准库 flag 包的所有功能，包括命令行参数解析、帮助信息生成等
	flags.AddGoFlagSet(goflag.CommandLine)
}

// PrintFlags logs the flags in the flagset.
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.Debugf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
