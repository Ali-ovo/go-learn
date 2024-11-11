package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Command 是 cli 应用程序的子命令结构。建议使用 app.NewCommand() 函数创建命令
type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command
	runFunc  RunCommandFunc
}

// CommandOption 定义了初始化命令结构的可选参数
type CommandOption func(*Command)

// WithCommandOptions 方法开启应用程序的命令行读取功能
func WithCommandOptions(opt CliOptions) CommandOption {
	return func(c *Command) {
		c.options = opt
	}
}

// RunCommandFunc 定义了应用程序的子命令启动回调函数
type RunCommandFunc func(args []string) error

// WithCommandRunFunc 用于设置应用程序命令的启动回调函数选项
func WithCommandRunFunc(run RunCommandFunc) CommandOption {
	return func(c *Command) {
		c.runFunc = run
	}
}

// NewCommand 根据给定的命令名称和其他选项创建一个新的子命令实例
func NewCommand(usage string, desc string, opts ...CommandOption) *Command {
	c := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// AddCommand 将子命令添加到当前命令中
func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

// AddCommands 用于将多个子命令添加到当前命令
func (c *Command) AddCommands(cmds ...*Command) {
	c.commands = append(c.commands, cmds...)
}

func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stdout)
	cmd.Flags().SortFlags = false
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}
	if c.runFunc != nil {
		// cobra Command 结构体 添加 runCommand 函数 之后执行
		cmd.Run = c.runCommand
	}
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
		// c.options.AddFlags(cmd.Flags())
	}
	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

// 此函数分装着需要执行的真正函数  c.runFunc(args)
// 如果有错误 打印错误 并退出
func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

// AddCommand 函数用于向应用中添加子命令
func (a *App) AddCommand(cmd *Command) {
	a.commands = append(a.commands, cmd)
}

// AddCommands 向应用程序添加多个子命令
func (a *App) AddCommands(cmds ...*Command) {
	a.commands = append(a.commands, cmds...)
}

// FormatBaseName 根据给定的名称，根据不同操作系统格式化为可执行文件名
func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}
