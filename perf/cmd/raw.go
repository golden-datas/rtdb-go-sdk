package cmd

import (
	"github.com/spf13/cobra"
)

// rawCmd Raw API 测试命令组
var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Raw API 性能测试",
	Long: `使用 Raw API 进行性能测试。

Raw API 提供了直接调用 C 库的接口，性能更高但使用较复杂。

子命令：
  test-write-rt   实时数据写入测试
  test-write-his  历史数据写入测试
  test-read-rt    实时数据查询测试
  test-read-his   历史数据查询测试`,
}

func init() {
	rootCmd.AddCommand(rawCmd)
}
