package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/kkbase/rtdb_api/perf/points"
)

var (
	genOutput   string
	genTable    string
	genI8       int
	genI16      int
	genI32      int
	genI64      int
	genU8       int
	genU16      int
	genU32      int
	genF32      int
	genF64      int
	genString   int
	genBlob     int
	genDatetime int
	genCoor     int
)

var genCmd = &cobra.Command{
	Use:   "gen-csv",
	Short: "生成测点CSV文件",
	Long: `根据指定的类型和数量生成测点CSV文件。

支持的数据类型：
  i8, i16, i32, i64    - 有符号整数
  u8, u16, u32         - 无符号整数  
  f32, f64             - 浮点数
  string               - 字符串
  blob                 - 二进制数据
  datetime             - 日期时间
  coordinates          - 坐标

示例：
  # 生成1000个i64类型和2000个f64类型的测点
  perf gen-csv -o points.csv --i64=1000 --f64=2000

  # 生成多种类型的测点
  perf gen-csv -o points.csv --i64=500 --f64=1000 --string=100 --blob=50`,
	RunE: runGen,
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&genOutput, "output", "o", "points.csv", "输出CSV文件路径")
	genCmd.Flags().StringVarP(&genTable, "table", "t", "perf_test", "测点所属表名")

	genCmd.Flags().IntVar(&genI8, "i8", 0, "i8类型测点数量")
	genCmd.Flags().IntVar(&genI16, "i16", 0, "i16类型测点数量")
	genCmd.Flags().IntVar(&genI32, "i32", 0, "i32类型测点数量")
	genCmd.Flags().IntVar(&genI64, "i64", 0, "i64类型测点数量")
	genCmd.Flags().IntVar(&genU8, "u8", 0, "u8类型测点数量")
	genCmd.Flags().IntVar(&genU16, "u16", 0, "u16类型测点数量")
	genCmd.Flags().IntVar(&genU32, "u32", 0, "u32类型测点数量")
	genCmd.Flags().IntVar(&genF32, "f32", 0, "f32类型测点数量")
	genCmd.Flags().IntVar(&genF64, "f64", 0, "f64类型测点数量")
	genCmd.Flags().IntVar(&genString, "string", 0, "string类型测点数量")
	genCmd.Flags().IntVar(&genBlob, "blob", 0, "blob类型测点数量")
	genCmd.Flags().IntVar(&genDatetime, "datetime", 0, "datetime类型测点数量")
	genCmd.Flags().IntVar(&genCoor, "coordinates", 0, "coordinates类型测点数量")
}

func runGen(cmd *cobra.Command, args []string) error {
	// 解析类型数量
	counts := points.ParseTypeCounts(
		genI8, genI16, genI32, genI64,
		genU8, genU16, genU32,
		genF32, genF64,
		genString, genBlob, genDatetime, genCoor,
	)

	// 检查是否指定了任何类型
	total := points.GetPointCount(counts)
	if total == 0 {
		return fmt.Errorf("请至少指定一种类型的测点数量，使用 --help 查看可用选项")
	}

	// 生成测点定义
	pointDefs := points.GeneratePoints(genTable, counts)

	// 写入CSV文件
	if err := points.WriteCSV(genOutput, pointDefs); err != nil {
		return fmt.Errorf("写入CSV文件失败: %w", err)
	}

	// 输出统计信息
	fmt.Printf("成功生成 %d 个测点定义\n", total)
	fmt.Printf("表名: %s\n", genTable)
	fmt.Printf("类型分布: %s\n", points.FormatPointCount(counts))
	fmt.Printf("输出文件: %s\n", genOutput)

	return nil
}
