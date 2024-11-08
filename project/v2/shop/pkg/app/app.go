package app

import (
	"fmt"
	"go-learn/project/v2/shop/pkg/common/cli/flag"
	"go-learn/project/v2/shop/pkg/common/cli/globalflag"
	"go-learn/project/v2/shop/pkg/common/term"
	"go-learn/project/v2/shop/pkg/common/version"
	"go-learn/project/v2/shop/pkg/common/version/verflag"
	"go-learn/project/v2/shop/pkg/errors"
	"go-learn/project/v2/shop/pkg/log"

	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var progressMessage = color.GreenString("==>")

//nolint: deadcode,unused,varcheck

// App 是一个命令行应用程序的主要结构。
// 建议使用 app.NewApp() 函数创建一个 app
type App struct {
	basename    string
	name        string
	description string
	options     CliOptions
	runFunc     RunFunc
	silence     bool
	noVersion   bool
	noConfig    bool
	commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

// Option 定义了初始化应用程序结构的可选参数
type Option func(*App)

// WithOptions 用于开启应用程序从命令行读取参数或从配置文件中读取参数的功能
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// RunFunc 定义应用程序的启动回调函数
type RunFunc func(basename string) error

// WithRunFunc 用于设置应用程序启动回调函数的选项
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithDescription 用于设置应用程序的描述
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithSilence 用于将应用程序设置为静默模式，在该模式下，程序启动信息、配置信息和版本信息不会在控制台中打印出来
func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoVersion 设置应用程序不提供版本标志
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithNoConfig 设置应用程序 不使用 config 文件
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithValidArgs set the validation function to valid non-flag arguments.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// NewApp 根据给定的应用程序名称、项目文件名称以及其他选项创建一个新的应用程序实例
func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,     // 启动 app 名字
		basename: basename, // 启动 项目 名字
	}

	for _, o := range opts { // 执行附加参数
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:   FormatBaseName(a.basename),
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true, // 表示当命令出现错误时，不会打印使用情况
		SilenceErrors: true, // 表示当命令出现错误时，不会输出任何内容
		Args:          a.args,
	}
	//cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOut(os.Stdout)        // 输入到标准输出流(cobra 默认 不填也没事 )
	cmd.SetErr(os.Stderr)        // 输出到 标准错误输出流(cobra 默认 不填也没事)
	cmd.Flags().SortFlags = true // 参数配置排序 按照字母顺序显示在帮助文档中
	flag.InitFlags(cmd.Flags())  // 规范化标志的名称和值 并且让 Cobra 库支持 Go 标准库 flag 包的所有功能

	// 设置 cobra 子命令
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		// 设置帮助命令
		cmd.SetHelpCommand(helpCommand(a.name))
	}
	// 设置 corbra.RunE  解析成功后 运行此方法
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets flag.NamedFlagSets
	if a.options != nil {
		// 获取分装了 *pflag.FlagSet 的 cliflag.NamedFlagSets 结构体 这个结构体可以存在多个 *pflag.FlagSet
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets { // 遍历 map[string]*pflag.FlagSet
			fs.AddFlagSet(f) // 集成到 cobra 中  AddFlagSet将一个 FlagSet 添加到另一个 FlagSet
			//fs.StringVar()	// 类似此操作 将 pflag 设置的表示全部添加到 fs 中
		}

		usageFmt := "Usage:\n  %s\n"
		cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine()) // 输出到指定输出流中
			flag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)                 // cols：一个整数，表示终端的列数，用于控制输出的宽度。
		})
		cmd.SetUsageFunc(func(cmd *cobra.Command) error {
			fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
			flag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols) // cols：一个整数，表示终端的列数，用于控制输出的宽度。

			return nil
		})
	}

	if !a.noVersion {
		// 添加版本号标志 到 global flagSet 下
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	if !a.noConfig {
		// 读取配置 获取参数
		// namedFlagSets 是是需要的所有标识 目前还没被赋值 为空
		// namedFlagSets 为封装后的 pflag 的结构体 FlagSet 是这个结构体的一个自定义方法
		// 可以想象成为创建一个新的命名空间
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}

	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	a.cmd = &cmd
}

// Run 用于启动应用程序
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command 方法返回应用程序中的 cobra.Command 实例
func (a *App) Command() *cobra.Command {
	return a.cmd
}

// cobra 最后执行的函数
func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	// 打印工作的目录
	printWorkingDir()
	// 打印所有的 Flags
	flag.PrintFlags(cmd.Flags())
	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	// 如果设置了 文件读取
	if !a.noConfig {
		// 与 cmd.Flags 参数 绑定
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		// 将 文件中的值 赋值给 a.options
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	// 这里可以添加 自定义规则 只需要 添加 Complete 方法
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	// 验证 字段是否合规
	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	// 这里可以添加 String 方法 就会打印输出以下的内容
	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}
