package app

import (
	"fmt"
	"shop/gmicro/pkg/common/util/homedir"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configFlagName = "config"

var cfgFile string

// nolint: gochecknoinits
func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}

// addConfigFlag adds flags for a specific server to the specified FlagSet
// object.
func addConfigFlag(basename string, fs *pflag.FlagSet) {
	// 添加 config 标志到 FlagSet
	// pflag.Lookup 方法查找现有的标志 查找 pflag 是否设置了 config 这个标志
	fs.AddFlag(pflag.Lookup(configFlagName))

	// viper.AutomaticEnv() 是 Go 语言中 Viper 配置库的一个方法，它可以自动读取环境变量并将其添加到配置中
	viper.AutomaticEnv()
	// SetEnvPrefix 方法将这个字符串设置为环境变量的前缀，
	// 这意味着 Viper 将在读取环境变量时只考虑以该前缀开头的变量。
	// 这使得在不同的应用程序中使用相同的环境变量名称变得更加容易，因为它们可以使用不同的前缀来区分它们。
	// 例如，如果 basename 是 "my-app"，则 SetEnvPrefix 方法将设置环境变量前缀为 MY_APP_。
	// 然后，当使用 viper.GetString("port") 读取 port 配置选项时，Viper 将尝试读取环境变量 MY_APP_PORT。
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	// 该方法接受一个字符串替换器 strings.NewReplacer(".", "_", "-", "_")，
	// 它将所有的点 . 和连字符 - 替换为下划线 _。这样做是因为在环境变量中不能使用点和连字符作为键名，所以必须使用下划线来替换它们。
	// 例如，如果使用 viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) 方法后，
	// 当使用 viper.GetString("database.username") 读取 database.username 配置选项时，
	// Viper 将尝试读取环境变量 DATABASE_USERNAME，并将点和连字符替换为下划线。
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			// 设置 config 读取文件地址
			viper.SetConfigFile(cfgFile)
		} else {
			// 在此路径下 检索配置文件
			viper.AddConfigPath(".")

			if names := strings.Split(basename, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homedir.HomeDir(), "."+names[0]))
			}
			// 设置 config 文件名称
			viper.SetConfigName(basename)
		}
		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}

	})
}

// nolint: deadcode,unused
func printConfig() {
	keys := viper.AllKeys()
	if len(keys) > 0 {
		fmt.Printf("%v Configuration items:\n", progressMessage)
		table := uitable.New()
		table.Separator = " "
		table.MaxColWidth = 80
		table.RightAlign(0)
		for _, k := range keys {
			table.AddRow(fmt.Sprintf("%s:", k), viper.Get(k))
		}
		fmt.Printf("%v", table)
	}
}
