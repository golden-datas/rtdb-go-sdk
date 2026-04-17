package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/kkbase/rtdb_api/perf/config"
	"github.com/kkbase/rtdb_api/perf/points"
)

var (
	createConfig    string
	createPointsFile string
	createTableDesc string
)

var createCmd = &cobra.Command{
	Use:   "create-points",
	Short: "根据CSV创建表和测点",
	Long: `根据CSV文件中的测点定义，自动创建表和测点。

如果表不存在，会自动创建；如果表已存在，则直接在现有表中创建测点。

示例：
  # 使用默认配置创建测点
  perf create-points -c db.yaml -f points.csv

  # 指定表描述
  perf create-points -c db.yaml -f points.csv --table-desc="性能测试表"`,
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&createConfig, "config", "c", "db.yaml", "数据库配置文件路径")
	createCmd.Flags().StringVarP(&createPointsFile, "points-file", "f", "points.csv", "测点CSV文件路径")
	createCmd.Flags().StringVar(&createTableDesc, "table-desc", "性能测试表", "表描述")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg, err := config.Load(createConfig)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 读取CSV
	defs, err := points.ReadCSV(createPointsFile)
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

	// 创建表
	fmt.Printf("正在创建/获取表 '%s'...\n", tableName)
	table, err := mgr.CreateTable(tableName, createTableDesc)
	if err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}
	fmt.Printf("表 '%s' (ID=%d) 准备就绪\n", table.Name, table.ID)

	// 创建测点
	fmt.Printf("正在创建 %d 个测点...\n", len(defs))
	infos, errs, err := mgr.CreatePoints(tableName, defs)
	if err != nil {
		return fmt.Errorf("创建测点失败: %w", err)
	}

	// 统计结果
	successCount := 0
	failCount := 0
	for i, e := range errs {
		if e != nil {
			failCount++
			fmt.Printf("  失败: %s - %v\n", defs[i].Name, e)
		} else {
			successCount++
		}
	}

	fmt.Printf("\n创建结果: 成功 %d 个, 失败 %d 个\n", successCount, failCount)
	_ = infos

	return nil
}
