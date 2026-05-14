package cmd

import (
	"fmt"

	"github.com/golden-datas/rtdb-go-sdk/perf/config"
	"github.com/golden-datas/rtdb-go-sdk/perf/points"
	"github.com/spf13/cobra"
)

var (
	deleteConfig     string
	deletePointsFile string
	deleteDropTable  bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete-points",
	Short: "根据CSV删除测点和表",
	Long: `根据CSV文件中的测点定义，删除测点。可选择是否同时删除表。

示例：
  # 仅删除测点，保留表
  perf delete-points -c db.yaml -f points.csv

  # 删除测点并删除表
  perf delete-points -c db.yaml -f points.csv --drop-table`,
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&deleteConfig, "config", "c", "db.yaml", "数据库配置文件路径")
	deleteCmd.Flags().StringVarP(&deletePointsFile, "points-file", "f", "points.csv", "测点CSV文件路径")
	deleteCmd.Flags().BoolVar(&deleteDropTable, "drop-table", false, "是否同时删除表")
}

func runDelete(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg, err := config.Load(deleteConfig)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 读取CSV
	defs, err := points.ReadCSV(deletePointsFile)
	if err != nil {
		return fmt.Errorf("读取CSV文件失败: %w", err)
	}

	if len(defs) == 0 {
		return fmt.Errorf("CSV文件中没有测点定义")
	}

	// 连接数据库
	conn, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer conn.Logout()

	fmt.Printf("成功连接到数据库 %s:%d\n", cfg.Host, cfg.Port)

	// 获取表名（假设所有测点都在同一个表）
	tableName := defs[0].Table

	// 创建测点管理器
	mgr := points.NewManager(conn)

	// 删除测点
	fmt.Printf("正在删除 %d 个测点...\n", len(defs))
	if err := mgr.DeletePoints(defs); err != nil {
		return fmt.Errorf("删除测点失败: %w", err)
	}
	fmt.Println("测点删除完成")

	// 如果需要，删除表
	if deleteDropTable {
		fmt.Printf("正在删除表 '%s'...\n", tableName)
		if err := mgr.DeleteTable(tableName); err != nil {
			return fmt.Errorf("删除表失败: %w", err)
		}
		fmt.Println("表删除完成")
	}

	return nil
}
