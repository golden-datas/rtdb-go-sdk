package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "perf",
	Short: "RTDB 性能测试工具",
	Long: `RTDB 性能测试工具 - 用于测试实时数据库的读写性能

主要功能：
  - 生成测点CSV文件
  - 创建/删除测点和表
  - 执行各类性能测试（实时/历史 写入/查询）

使用示例：
  perf init-config -o db.yaml
  perf gen-csv -o points.csv --i64=1000 --f64=2000
  perf create-points -c db.yaml -f points.csv
  perf test-write-rt -c db.yaml -f points.csv --mode=burst --workers=10`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 这里可以添加全局标志
}
