package cmd

import (
	"github.com/spf13/cobra"
)

// easyCmd Easy API 测试命令组
var easyCmd = &cobra.Command{
	Use:   "easy",
	Short: "Easy API 性能测试",
	Long: `使用 Easy API 进行性能测试。

Easy API 提供了面向对象的 Go 风格接口，使用更便捷。

子命令：
  test-write-rt   实时数据写入测试
  test-write-his  历史数据写入测试
  test-read-rt    实时数据查询测试
  test-read-his   历史数据查询测试`,
}

func init() {
	rootCmd.AddCommand(easyCmd)
}
