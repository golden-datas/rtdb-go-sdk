package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/kkbase/rtdb_api"
	"github.com/kkbase/rtdb_api/perf/config"
	"github.com/kkbase/rtdb_api/perf/points"
	"github.com/kkbase/rtdb_api/perf/tester"
	"github.com/kkbase/rtdb_api/perf/tester/raw"
)

var (
	rawReadConfig     string
	rawReadPointsFile string
	rawReadMode       string
	rawReadWorkers    int
	rawReadPointCount int
	rawReadRange      time.Duration
	rawReadInterpo    int
	rawReadPlot       int
	rawReadSamples    int
)

var rawReadCmd = &cobra.Command{
	Use:   "test-read",
	Short: "Raw API 数据读取测试",
	Long: `使用 Raw API 进行数据读取性能测试。

测试模式：
  rt       - 实时数据查询（最新值）
  raw      - 历史原始值查询
  interpo  - 历史插值查询
  plot     - 趋势曲线查询
  summary  - 历史统计值查询
  section  - 断面数据查询

示例：
  # 实时数据查询
  perf raw test-read -c db.yaml -f points.csv --mode=rt --workers=10

  # 历史原始值查询（查询最近1天）
  perf raw test-read -c db.yaml -f points.csv --mode=raw --range=24h --workers=10`,
	RunE: runRawRead,
}

func init() {
	rawCmd.AddCommand(rawReadCmd)
	rawReadCmd.Flags().StringVarP(&rawReadConfig, "config", "c", "db.yaml", "数据库配置文件路径")
	rawReadCmd.Flags().StringVarP(&rawReadPointsFile, "points-file", "f", "points.csv", "测点CSV文件路径")
	rawReadCmd.Flags().StringVarP(&rawReadMode, "mode", "m", "rt", "测试模式: rt|raw|interpo|plot|summary|section")
	rawReadCmd.Flags().IntVar(&rawReadWorkers, "workers", 1, "并发工作线程数")
	rawReadCmd.Flags().IntVar(&rawReadPointCount, "point-count", 0, "查询测点数(0表示全部)")
	rawReadCmd.Flags().DurationVar(&rawReadRange, "range", 24*time.Hour, "查询时间范围")
	rawReadCmd.Flags().IntVar(&rawReadInterpo, "interpo-count", 600, "插值数量")
	rawReadCmd.Flags().IntVar(&rawReadPlot, "plot-interval", 1000, "趋势曲线区间数")
	rawReadCmd.Flags().IntVar(&rawReadSamples, "samples", 1, "采样次数(rt模式)")
}

func runRawRead(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(rawReadConfig)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	defs, err := points.ReadCSV(rawReadPointsFile)
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

	reader := raw.NewReader(conn, pointInfos)
	metrics := tester.NewMetrics()

	end := time.Now()
	start := end.Add(-rawReadRange)

	switch rawReadMode {
	case "rt":
		fmt.Printf("\n开始实时数据查询测试（Raw API）：\n")
		fmt.Printf("  并发数: %d\n", rawReadWorkers)
		fmt.Printf("  采样次数: %d\n", rawReadSamples)
		if err := reader.ReadLast(rawReadWorkers, rawReadSamples, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "raw":
		fmt.Printf("\n开始历史原始值查询测试（Raw API）：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  并发数: %d\n", rawReadWorkers)
		if err := reader.ReadRaw(start, end, rawReadWorkers, rawReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "interpo":
		fmt.Printf("\n开始历史插值查询测试（Raw API）：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  插值数量: %d\n", rawReadInterpo)
		fmt.Printf("  并发数: %d\n", rawReadWorkers)
		if err := reader.ReadInterpo(start, end, rawReadInterpo, rawReadWorkers, rawReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "plot":
		fmt.Printf("\n开始趋势曲线查询测试（Raw API）：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  区间数: %d\n", rawReadPlot)
		fmt.Printf("  并发数: %d\n", rawReadWorkers)
		if err := reader.ReadPlot(start, end, rawReadPlot, rawReadWorkers, rawReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "summary":
		fmt.Printf("\n开始历史统计值查询测试（Raw API）：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  并发数: %d\n", rawReadWorkers)
		if err := reader.ReadSummary(start, end, rawReadWorkers, rawReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "section":
		fmt.Printf("\n开始断面数据查询测试（Raw API）：\n")
		fmt.Printf("  查询时刻: %v\n", end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  并发数: %d\n", rawReadWorkers)
		if err := reader.ReadSection(end, rawReadWorkers, rawReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	default:
		return fmt.Errorf("未知的测试模式: %s", rawReadMode)
	}

	metrics.Print()
	return nil
}
