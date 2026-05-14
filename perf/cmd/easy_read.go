package cmd

import (
	"fmt"
	"time"

	"github.com/golden-datas/rtdb-go-sdk"
	"github.com/golden-datas/rtdb-go-sdk/perf/config"
	"github.com/golden-datas/rtdb-go-sdk/perf/points"
	"github.com/golden-datas/rtdb-go-sdk/perf/tester"
	"github.com/golden-datas/rtdb-go-sdk/perf/tester/easy"
	"github.com/spf13/cobra"
)

var (
	easyReadConfig     string
	easyReadPointsFile string
	easyReadMode       string
	easyReadWorkers    int
	easyReadPointCount int
	easyReadRange      time.Duration
	easyReadInterpo    int
	easyReadPlot       int
	easyReadSamples    int
)

var easyReadCmd = &cobra.Command{
	Use:   "test-read",
	Short: "Easy API 数据读取测试",
	Long: `使用 Easy API 进行数据读取性能测试。

测试模式：
  rt       - 实时数据查询（最新值）
  raw      - 历史原始值查询
  interpo  - 历史插值查询
  plot     - 趋势曲线查询
  summary  - 历史统计值查询
  section  - 断面数据查询

示例：
  # 实时数据查询
  perf easy test-read -c db.yaml -f points.csv --mode=rt --workers=10

  # 历史原始值查询（查询最近1天）
  perf easy test-read -c db.yaml -f points.csv --mode=raw --range=24h --workers=10`,
	RunE: runEasyRead,
}

func init() {
	easyCmd.AddCommand(easyReadCmd)
	easyReadCmd.Flags().StringVarP(&easyReadConfig, "config", "c", "db.yaml", "数据库配置文件路径")
	easyReadCmd.Flags().StringVarP(&easyReadPointsFile, "points-file", "f", "points.csv", "测点CSV文件路径")
	easyReadCmd.Flags().StringVarP(&easyReadMode, "mode", "m", "rt", "测试模式: rt|raw|interpo|plot|summary|section")
	easyReadCmd.Flags().IntVar(&easyReadWorkers, "workers", 1, "并发工作线程数")
	easyReadCmd.Flags().IntVar(&easyReadPointCount, "point-count", 0, "查询测点数(0表示全部)")
	easyReadCmd.Flags().DurationVar(&easyReadRange, "range", 24*time.Hour, "查询时间范围")
	easyReadCmd.Flags().IntVar(&easyReadInterpo, "interpo-count", 600, "插值数量")
	easyReadCmd.Flags().IntVar(&easyReadPlot, "plot-interval", 1000, "趋势曲线区间数")
	easyReadCmd.Flags().IntVar(&easyReadSamples, "samples", 1, "采样次数(rt模式)")
}

func runEasyRead(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(easyReadConfig)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	defs, err := points.ReadCSV(easyReadPointsFile)
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

	reader := easy.NewReader(conn, pointInfos)
	metrics := tester.NewMetrics()

	end := time.Now()
	start := end.Add(-easyReadRange)

	switch easyReadMode {
	case "rt":
		fmt.Printf("\n开始实时数据查询测试：\n")
		fmt.Printf("  并发数: %d\n", easyReadWorkers)
		fmt.Printf("  采样次数: %d\n", easyReadSamples)
		if err := reader.ReadLast(easyReadWorkers, easyReadSamples, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "raw":
		fmt.Printf("\n开始历史原始值查询测试：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  并发数: %d\n", easyReadWorkers)
		if err := reader.ReadRaw(start, end, easyReadWorkers, easyReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "interpo":
		fmt.Printf("\n开始历史插值查询测试：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  插值数量: %d\n", easyReadInterpo)
		fmt.Printf("  并发数: %d\n", easyReadWorkers)
		if err := reader.ReadInterpo(start, end, easyReadInterpo, easyReadWorkers, easyReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "plot":
		fmt.Printf("\n开始趋势曲线查询测试：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  区间数: %d\n", easyReadPlot)
		fmt.Printf("  并发数: %d\n", easyReadWorkers)
		if err := reader.ReadPlot(start, end, easyReadPlot, easyReadWorkers, easyReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "summary":
		fmt.Printf("\n开始历史统计值查询测试：\n")
		fmt.Printf("  时间范围: %v ~ %v\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  并发数: %d\n", easyReadWorkers)
		if err := reader.ReadSummary(start, end, easyReadWorkers, easyReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	case "section":
		fmt.Printf("\n开始断面数据查询测试：\n")
		fmt.Printf("  查询时刻: %v\n", end.Format("2006-01-02 15:04:05"))
		fmt.Printf("  并发数: %d\n", easyReadWorkers)
		if err := reader.ReadSection(end, easyReadWorkers, easyReadPointCount, metrics); err != nil {
			return fmt.Errorf("测试失败: %w", err)
		}
	default:
		return fmt.Errorf("未知的测试模式: %s", easyReadMode)
	}

	metrics.Print()
	return nil
}
