package points

import (
	"fmt"

	"github.com/kkbase/rtdb_api"
)

// Manager 测点管理器
type Manager struct {
	conn *rtdb_api.RtdbConnect
}

// NewManager 创建测点管理器
func NewManager(conn *rtdb_api.RtdbConnect) *Manager {
	return &Manager{conn: conn}
}

// CreateTable 创建表
func (m *Manager) CreateTable(name, desc string) (*rtdb_api.RtdbTable, error) {
	// 先检查表是否已存在
	tables, err := m.conn.GetTables()
	if err != nil {
		return nil, fmt.Errorf("get tables failed: %w", err)
	}

	for _, t := range tables {
		if t.Name == name {
			// 表已存在，返回已有表
			return &t, nil
		}
	}

	// 创建新表
	table, err := m.conn.CreateTable(name, desc)
	if err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}

	return table, nil
}

// CreatePoints 批量创建测点
func (m *Manager) CreatePoints(tableName string, defs []PointDef) ([]*rtdb_api.PointInfo, []error, error) {
	// 获取表信息
	tables, err := m.conn.GetTables()
	if err != nil {
		return nil, nil, fmt.Errorf("get tables failed: %w", err)
	}

	var tableID rtdb_api.TableID
	found := false
	for _, t := range tables {
		if t.Name == tableName {
			tableID = t.ID
			found = true
			break
		}
	}

	if !found {
		return nil, nil, fmt.Errorf("table '%s' not found", tableName)
	}

	// 批量创建测点
	var infos []*rtdb_api.PointInfo
	var errs []error

	for _, def := range defs {
		info := m.createPointInfo(tableID, def)
		newInfo, err := m.conn.AddPoint(info)
		if err != nil {
			errs = append(errs, fmt.Errorf("create point '%s' failed: %w", def.Name, err))
			infos = append(infos, nil)
		} else {
			errs = append(errs, nil)
			infos = append(infos, newInfo)
		}
	}

	return infos, errs, nil
}

// DeletePoints 批量删除测点
func (m *Manager) DeletePoints(defs []PointDef) error {
	// 获取所有表
	tables, err := m.conn.GetTables()
	if err != nil {
		return fmt.Errorf("get tables failed: %w", err)
	}

	// 构建表名到ID的映射
	tableMap := make(map[string]rtdb_api.TableID)
	for _, t := range tables {
		tableMap[t.Name] = t.ID
	}

	// 逐个删除测点
	for _, def := range defs {
		_, ok := tableMap[def.Table]
		if !ok {
			continue // 表不存在，跳过
		}

		// 构造表.点名查找点
		tableDotTag := def.Table + "." + def.Name
		pointInfos, _, err := m.conn.FindPoints([]string{tableDotTag})
		if err != nil || len(pointInfos) == 0 || pointInfos[0] == nil {
			continue // 点不存在，跳过
		}

		// 删除测点
		if err := m.conn.DeletePoint(pointInfos[0].ID); err != nil {
			fmt.Printf("Warning: delete point '%s' failed: %v\n", def.Name, err)
		}
	}

	return nil
}

// DeleteTable 删除表
func (m *Manager) DeleteTable(name string) error {
	// 获取表ID
	tables, err := m.conn.GetTables()
	if err != nil {
		return fmt.Errorf("get tables failed: %w", err)
	}

	for _, t := range tables {
		if t.Name == name {
			if err := m.conn.DeleteTable(t.ID); err != nil {
				return fmt.Errorf("delete table failed: %w", err)
			}
			return nil
		}
	}

	return nil // 表不存在，视为成功
}

// createPointInfo 根据定义创建PointInfo
func (m *Manager) createPointInfo(tableID rtdb_api.TableID, def PointDef) *rtdb_api.PointInfo {
	// 解析类型
	valueType := parseValueType(def.Type)

	// 解析类别
	class := rtdb_api.PointBase
	switch def.Class {
	case "scan":
		class = rtdb_api.PointScan
	case "calc":
		class = rtdb_api.PointCalc
	}

	// 解析精度
	precision := rtdb_api.RtdbPrecisionSecond
	switch def.Precision {
	case "milli":
		precision = rtdb_api.RtdbPrecisionMilli
	case "micro":
		precision = rtdb_api.RtdbPrecisionMicro
	case "nano":
		precision = rtdb_api.RtdbPrecisionNano
	}

	// 创建基础点信息
	info := rtdb_api.NewPointInfo(
		def.Name,
		tableID,
		valueType,
		class,
		precision,
		def.Unit,
		def.Desc,
	)

	// 设置存档相关属性（关闭压缩）
	info.Archive = rtdb_api.ON
	info.Compress = rtdb_api.OFF
	info.Step = rtdb_api.OFF
	info.Summary = rtdb_api.OFF

	// 如果是采集点，设置采集属性
	if class == rtdb_api.PointScan && def.Source != "" {
		info.Source = def.Source
		info.Scan = rtdb_api.ON
		info.Instrument = def.Instrument
	}

	// 如果是计算点，设置计算属性
	if class == rtdb_api.PointCalc && def.Equation != "" {
		info.SetCalc(def.Equation, rtdb_api.RtdbEventTrigger, rtdb_api.RtdbCalcTime, 0)
	}

	return info
}

// parseValueType 解析类型字符串
func parseValueType(typ string) rtdb_api.ValueType {
	switch typ {
	case "i8":
		return rtdb_api.ValueTypeInt8
	case "i16":
		return rtdb_api.ValueTypeInt16
	case "i32":
		return rtdb_api.ValueTypeInt32
	case "i64":
		return rtdb_api.ValueTypeInt64
	case "u8":
		return rtdb_api.ValueTypeUint8
	case "u16":
		return rtdb_api.ValueTypeUint16
	case "u32":
		return rtdb_api.ValueTypeUint32
	case "f32":
		return rtdb_api.ValueTypeFloat32
	case "f64":
		return rtdb_api.ValueTypeFloat64
	case "string":
		return rtdb_api.ValueTypeString
	case "blob":
		return rtdb_api.ValueTypeBlob
	case "datetime":
		return rtdb_api.ValueTypeDatetime
	case "coordinates":
		return rtdb_api.ValueTypeCoor
	default:
		return rtdb_api.ValueTypeFloat64
	}
}
