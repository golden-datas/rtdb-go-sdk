package cmd

import (
	"fmt"
	"time"

	"github.com/golden-datas/rtdb-go-sdk"
	"github.com/golden-datas/rtdb-go-sdk/perf/config"
	"github.com/golden-datas/rtdb-go-sdk/perf/points"
	"github.com/golden-datas/rtdb-go-sdk/perf/tester"
	"github.com/golden-datas/rtdb-go-sdk/perf/tester/raw"
	"github.com/spf13/cobra"
)

var (
	rawWriteHisConfig     string
	rawWriteHisPointsFile string
	rawWriteHisMode       string
	rawWriteHisBatchSize  int
	rawWriteHisInterval   time.Duration
	rawWriteHisDuration   time.Duration
	rawWriteHisBatches    int
	rawWriteHisWorkers    int
)

var rawWriteHisCmd = &cobra.Command{
	Use:   "test-write-his",
	Short: "Raw API 历史数据写入测试",
	Long:  `使用 Raw API 进行历史数据写入性能测试。`,
	RunE:  runRawWriteHis,
}

func init() {
	rawCmd.AddCommand(rawWriteHisCmd)
	rawWriteHisCmd.Flags().StringVarP(&rawWriteHisConfig, "config", "c", "db.yaml", "数据库配置文件路径")
	rawWriteHisCmd.Flags().StringVarP(&rawWriteHisPointsFile, "points-file", "f", "points.csv", "测点CSV文件路径")
	rawWriteHisCmd.Flags().StringVarP(&rawWriteHisMode, "mode", "m", "periodic", "测试模式: periodic|burst")
	rawWriteHisCmd.Flags().IntVar(&rawWriteHisBatchSize, "batch-size", 10000, "每批写入的测点数")
	rawWriteHisCmd.Flags().DurationVar(&rawWriteHisInterval, "interval", time.Second, "写入间隔(periodic模式)")
	rawWriteHisCmd.Flags().DurationVar(&rawWriteHisDuration, "duration", 30*time.Minute, "测试持续时间(periodic模式)")
	rawWriteHisCmd.Flags().IntVar(&rawWriteHisBatches, "batches", 3600, "每点写入批次数(burst模式)")
	rawWriteHisCmd.Flags().IntVar(&rawWriteHisWorkers, "workers", 1, "并发工作线程数")
}

func runRawWriteHis(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(rawWriteHisConfig)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	defs, err := points.ReadCSV(rawWriteHisPointsFile)
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

	writer := raw.NewWriter(conn, pointInfos)
	metrics := tester.NewMetrics()

	switch rawWriteHisMode {
	case "periodic":
		fmt.Printf("\n开始周期性历史写入测试（Raw API）：\n")
		fmt.Printf("  每批写入: %d 点\n", rawWriteHisBatchSize)
		fmt.Printf("  写入间隔: %v\n", rawWriteHisInterval)
		fmt.Printf("  持续时间: %v\n", rawWriteHisDuration)
		if err := writer.WriteHisPeriodic(rawWriteHisBatchSize, rawWriteHisInterval, rawWriteHisDuration, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "burst":
		fmt.Printf("\n开始急速历史写入测试（Raw API）：\n")
		fmt.Printf("  每点写入: %d 批\n", rawWriteHisBatches)
		fmt.Printf("  并发数: %d\n", rawWriteHisWorkers)
		if err := writer.WriteHisBurst(rawWriteHisBatchSize, rawWriteHisBatches, rawWriteHisWorkers, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	default:
		return fmt.Errorf("未知的测试模式: %s", rawWriteHisMode)
	}

	metrics.Print()
	return nil
}
