package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/kkbase/rtdb_api/perf/config"
)

var initOutput string

var initCmd = &cobra.Command{
	Use:   "init-config",
	Short: "生成数据库配置文件模板",
	Long: `生成数据库配置文件模板，用于存储数据库连接信息。

生成的配置文件包含以下字段：
  host     - 数据库服务器地址
  port     - 数据库服务器端口
  username - 用户名
  password - 密码

示例：
  perf init-config -o db.yaml`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "db.yaml", "输出配置文件路径")
}

func runInit(cmd *cobra.Command, args []string) error {
	if err := config.GenerateTemplate(initOutput); err != nil {
		return fmt.Errorf("生成配置文件失败: %w", err)
	}

	fmt.Printf("成功生成配置文件: %s\n", initOutput)
	fmt.Println("请修改配置文件中的连接信息后再使用其他命令")
	return nil
}
