package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "a brief description of your application",
	Long:  "a longer description of your application",
}

var mockMsgCmd = &cobra.Command{
	Use:   "mockMsg",
	Short: "批量发送",
	Long:  "mockMsg",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("call mockMsg")
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出数据",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("call export")
		return nil
	},
}

func main() {
	rootCmd.AddCommand(mockMsgCmd)
	rootCmd.AddCommand(exportCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
