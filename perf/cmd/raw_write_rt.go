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
	rawWriteRtConfig     string
	rawWriteRtPointsFile string
	rawWriteRtMode       string
	rawWriteRtBatchSize  int
	rawWriteRtInterval   time.Duration
	rawWriteRtDuration   time.Duration
	rawWriteRtBatches    int
	rawWriteRtWorkers    int
)

var rawWriteRtCmd = &cobra.Command{
	Use:   "test-write-rt",
	Short: "Raw API 实时数据写入测试",
	Long: `使用 Raw API 进行实时数据写入性能测试。

测试模式：
  periodic - 周期性写入：按指定频率持续写入
  burst    - 急速写入：尽快完成指定数量的写入

示例：
  # 周期性写入：每秒写入10000点，持续30分钟
  perf raw test-write-rt -c db.yaml -f points.csv --mode=periodic --batch-size=10000 --interval=1s --duration=30m

  # 急速写入：100并发，每个点写入3600批
  perf raw test-write-rt -c db.yaml -f points.csv --mode=burst --batches=3600 --workers=100`,
	RunE: runRawWriteRt,
}

func init() {
	rawCmd.AddCommand(rawWriteRtCmd)
	rawWriteRtCmd.Flags().StringVarP(&rawWriteRtConfig, "config", "c", "db.yaml", "数据库配置文件路径")
	rawWriteRtCmd.Flags().StringVarP(&rawWriteRtPointsFile, "points-file", "f", "points.csv", "测点CSV文件路径")
	rawWriteRtCmd.Flags().StringVarP(&rawWriteRtMode, "mode", "m", "periodic", "测试模式: periodic|burst")
	rawWriteRtCmd.Flags().IntVar(&rawWriteRtBatchSize, "batch-size", 10000, "每批写入的测点数")
	rawWriteRtCmd.Flags().DurationVar(&rawWriteRtInterval, "interval", time.Second, "写入间隔(periodic模式)")
	rawWriteRtCmd.Flags().DurationVar(&rawWriteRtDuration, "duration", 30*time.Minute, "测试持续时间(periodic模式)")
	rawWriteRtCmd.Flags().IntVar(&rawWriteRtBatches, "batches", 3600, "每点写入批次数(burst模式)")
	rawWriteRtCmd.Flags().IntVar(&rawWriteRtWorkers, "workers", 1, "并发工作线程数")
}

func runRawWriteRt(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg, err := config.Load(rawWriteRtConfig)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 读取CSV
	defs, err := points.ReadCSV(rawWriteRtPointsFile)
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

	// 获取测点信息
	fmt.Printf("正在获取 %d 个测点信息...\n", len(defs))
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

	// 创建写入测试器
	writer := raw.NewWriter(conn, pointInfos)
	metrics := tester.NewMetrics()

	// 执行测试
	switch rawWriteRtMode {
	case "periodic":
		fmt.Printf("\n开始周期性写入测试（Raw API）：\n")
		fmt.Printf("  每批写入: %d 点\n", rawWriteRtBatchSize)
		fmt.Printf("  写入间隔: %v\n", rawWriteRtInterval)
		fmt.Printf("  持续时间: %v\n", rawWriteRtDuration)
		if err := writer.WritePeriodic(rawWriteRtBatchSize, rawWriteRtInterval, rawWriteRtDuration, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "burst":
		fmt.Printf("\n开始急速写入测试（Raw API）：\n")
		fmt.Printf("  每点写入: %d 批\n", rawWriteRtBatches)
		fmt.Printf("  并发数: %d\n", rawWriteRtWorkers)
		if err := writer.WriteBurst(rawWriteRtBatchSize, rawWriteRtBatches, rawWriteRtWorkers, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	default:
		return fmt.Errorf("未知的测试模式: %s", rawWriteRtMode)
	}

	// 打印结果
	metrics.Print()
	return nil
}
