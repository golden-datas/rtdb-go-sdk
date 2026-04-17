package points

import (
	"encoding/csv"
	"fmt"
	"os"
)

// PointDef 定义测点结构
type PointDef struct {
	Name      string
	Table     string
	Type      string // i8,i16,i32,i64,u8,u16,u32,f32,f64,string,blob,datetime,coordinates
	Desc      string
	Unit      string
	Class     string // base,scan,calc
	Precision string // second,milli,micro,nano

	// scan专用
	Source     string
	Instrument string

	// calc专用
	Equation string
}

// CSVHeaders CSV文件表头
var CSVHeaders = []string{"name", "table", "type", "desc", "unit", "class", "precision", "source", "instrument", "equation"}

// ReadCSV 从CSV文件读取测点定义
func ReadCSV(path string) ([]PointDef, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open csv file failed: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv file failed: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("csv file is empty or has no data rows")
	}

	// 解析表头
	headerMap := make(map[string]int)
	for i, h := range records[0] {
		headerMap[h] = i
	}

	// 检查必需字段
	requiredFields := []string{"name", "table", "type"}
	for _, field := range requiredFields {
		if _, ok := headerMap[field]; !ok {
			return nil, fmt.Errorf("required field '%s' not found in csv header", field)
		}
	}

	// 解析数据行
	var points []PointDef
	for i, record := range records[1:] {
		point := PointDef{
			Name:      getField(record, headerMap, "name"),
			Table:     getField(record, headerMap, "table"),
			Type:      getField(record, headerMap, "type"),
			Desc:      getField(record, headerMap, "desc"),
			Unit:      getField(record, headerMap, "unit"),
			Class:     getFieldDefault(record, headerMap, "class", "base"),
			Precision: getFieldDefault(record, headerMap, "precision", "second"),
			Source:    getField(record, headerMap, "source"),
			Instrument: getField(record, headerMap, "instrument"),
			Equation:  getField(record, headerMap, "equation"),
		}

		if point.Name == "" {
			return nil, fmt.Errorf("row %d: name is empty", i+2)
		}
		if point.Table == "" {
			return nil, fmt.Errorf("row %d: table is empty", i+2)
		}
		if point.Type == "" {
			return nil, fmt.Errorf("row %d: type is empty", i+2)
		}

		points = append(points, point)
	}

	return points, nil
}

// WriteCSV 将测点定义写入CSV文件
func WriteCSV(path string, points []PointDef) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create csv file failed: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write(CSVHeaders); err != nil {
		return fmt.Errorf("write csv header failed: %w", err)
	}

	// 写入数据
	for _, p := range points {
		record := []string{
			p.Name,
			p.Table,
			p.Type,
			p.Desc,
			p.Unit,
			p.Class,
			p.Precision,
			p.Source,
			p.Instrument,
			p.Equation,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write csv record failed: %w", err)
		}
	}

	return nil
}

// GeneratePoints 生成测点定义列表
func GeneratePoints(table string, typeCounts map[string]int) []PointDef {
	var points []PointDef

	// 类型到名称前缀的映射
	typePrefixes := map[string]string{
		"i8":          "i8",
		"i16":         "i16",
		"i32":         "i32",
		"i64":         "i64",
		"u8":          "u8",
		"u16":         "u16",
		"u32":         "u32",
		"f32":         "f32",
		"f64":         "f64",
		"string":      "str",
		"blob":        "blob",
		"datetime":    "dt",
		"coordinates": "coor",
	}

	for typ, count := range typeCounts {
		if count <= 0 {
			continue
		}

		prefix, ok := typePrefixes[typ]
		if !ok {
			prefix = typ
		}

		for i := 0; i < count; i++ {
			point := PointDef{
				Name:      fmt.Sprintf("tag_%s_%03d", prefix, i+1),
				Table:     table,
				Type:      typ,
				Desc:      fmt.Sprintf("测点 %s_%03d", prefix, i+1),
				Class:     "base",
				Precision: "second",
			}
			points = append(points, point)
		}
	}

	return points
}

// getField 获取字段值，如果不存在返回空字符串
func getField(record []string, headerMap map[string]int, field string) string {
	if idx, ok := headerMap[field]; ok && idx < len(record) {
		return record[idx]
	}
	return ""
}

// getFieldDefault 获取字段值，如果不存在返回默认值
func getFieldDefault(record []string, headerMap map[string]int, field, defaultValue string) string {
	value := getField(record, headerMap, field)
	if value == "" {
		return defaultValue
	}
	return value
}

// ParseTypeCounts 解析命令行传入的类型数量参数
func ParseTypeCounts(
	i8, i16, i32, i64 int,
	u8, u16, u32 int,
	f32, f64 int,
	str, blob, dt, coor int,
) map[string]int {
	return map[string]int{
		"i8":          i8,
		"i16":         i16,
		"i32":         i32,
		"i64":         i64,
		"u8":          u8,
		"u16":         u16,
		"u32":         u32,
		"f32":         f32,
		"f64":         f64,
		"string":      str,
		"blob":        blob,
		"datetime":    dt,
		"coordinates": coor,
	}
}

// GetPointCount 获取测点总数
func GetPointCount(counts map[string]int) int {
	total := 0
	for _, c := range counts {
		total += c
	}
	return total
}

// FormatPointCount 格式化输出测点数量统计
func FormatPointCount(counts map[string]int) string {
	result := ""
	for typ, count := range counts {
		if count > 0 {
			if result != "" {
				result += ", "
			}
			result += fmt.Sprintf("%s=%d", typ, count)
		}
	}
	return result
}
