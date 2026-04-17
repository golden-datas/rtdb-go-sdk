package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/kkbase/rtdb_api"
	"github.com/kkbase/rtdb_api/perf/config"
	"github.com/kkbase/rtdb_api/perf/points"
	"github.com/kkbase/rtdb_api/perf/tester"
	"github.com/kkbase/rtdb_api/perf/tester/easy"
)

var (
	easyWriteHisConfig     string
	easyWriteHisPointsFile string
	easyWriteHisMode       string
	easyWriteHisBatchSize  int
	easyWriteHisInterval   time.Duration
	easyWriteHisDuration   time.Duration
	easyWriteHisBatches    int
	easyWriteHisWorkers    int
)

var easyWriteHisCmd = &cobra.Command{
	Use:   "test-write-his",
	Short: "Easy API 历史数据写入测试",
	Long:  `使用 Easy API 进行历史数据写入性能测试。`,
	RunE:  runEasyWriteHis,
}

func init() {
	easyCmd.AddCommand(easyWriteHisCmd)
	easyWriteHisCmd.Flags().StringVarP(&easyWriteHisConfig, "config", "c", "db.yaml", "数据库配置文件路径")
	easyWriteHisCmd.Flags().StringVarP(&easyWriteHisPointsFile, "points-file", "f", "points.csv", "测点CSV文件路径")
	easyWriteHisCmd.Flags().StringVarP(&easyWriteHisMode, "mode", "m", "periodic", "测试模式: periodic|burst")
	easyWriteHisCmd.Flags().IntVar(&easyWriteHisBatchSize, "batch-size", 10000, "每批写入的测点数")
	easyWriteHisCmd.Flags().DurationVar(&easyWriteHisInterval, "interval", time.Second, "写入间隔(periodic模式)")
	easyWriteHisCmd.Flags().DurationVar(&easyWriteHisDuration, "duration", 30*time.Minute, "测试持续时间(periodic模式)")
	easyWriteHisCmd.Flags().IntVar(&easyWriteHisBatches, "batches", 3600, "每点写入批次数(burst模式)")
	easyWriteHisCmd.Flags().IntVar(&easyWriteHisWorkers, "workers", 1, "并发工作线程数")
}

func runEasyWriteHis(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(easyWriteHisConfig)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	defs, err := points.ReadCSV(easyWriteHisPointsFile)
	if err != nil {
		return fmt.Errorf("读取CSV文件失败: %w", err)
	}

	if len(defs) == 0 {
		return fmt.Errorf("CSV文件中没有测点定义")
	}

	conn, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer conn.Logout()

	fmt.Printf("成功连接到数据库 %s:%d\n", cfg.Host, cfg.Port)

	pointInfos := make([]*rtdb_api.PointInfo, 0, len(defs))
	for _, def := range defs {
		tableDotTag := def.Table + "." + def.Name
		infos, _, err := conn.FindPoints([]string{tableDotTag})
		if err != nil || len(infos) == 0 || infos[0] == nil {
			fmt.Printf("警告: 无法找到测点 %s\n", tableDotTag)
			continue
		}
		pointInfos = append(pointInfos, infos[0])
	}

	if len(pointInfos) == 0 {
		return fmt.Errorf("没有有效的测点可供测试")
	}

	fmt.Printf("成功获取 %d 个测点信息\n", len(pointInfos))

	writer := easy.NewWriter(conn, pointInfos)
	metrics := tester.NewMetrics()

	switch easyWriteHisMode {
	case "periodic":
		fmt.Printf("\n开始周期性历史写入测试：\n")
		fmt.Printf("  每批写入: %d 点\n", easyWriteHisBatchSize)
		fmt.Printf("  写入间隔: %v\n", easyWriteHisInterval)
		fmt.Printf("  持续时间: %v\n", easyWriteHisDuration)
		if err := writer.WriteHisPeriodic(easyWriteHisBatchSize, easyWriteHisInterval, easyWriteHisDuration, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "burst":
		fmt.Printf("\n开始急速历史写入测试：\n")
		fmt.Printf("  每点写入: %d 批\n", easyWriteHisBatches)
		fmt.Printf("  并发数: %d\n", easyWriteHisWorkers)
		if err := writer.WriteHisBurst(easyWriteHisBatchSize, easyWriteHisBatches, easyWriteHisWorkers, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	default:
		return fmt.Errorf("未知的测试模式: %s", easyWriteHisMode)
	}

	metrics.Print()
	return nil
}
