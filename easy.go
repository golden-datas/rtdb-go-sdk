package rtdb_api

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"sort"
	"strconv"
	"time"
	"unsafe"
)

const (
	// MaxBlockSize 单次读取文件最大大小
	MaxBlockSize = 5 * 1024 * 1024

	// MaxFileSize 允许读取文件最大大小
	MaxFileSize = 512 * 1024 * 1024
)

// ServerOption 服务端配置
type ServerOption struct {
	IsString     bool
	StringOption ParamString
	IntOption    ParamInt
}

// NewServerOption 新建服务端类型（通过字面值新建服务端配置, 会自动推断配置类型是String或Int）
func NewServerOption(option string) ServerOption {
	intOption, err := strconv.Atoi(option)
	if err != nil {
		return NewStringServerOption(ParamString(option))
	} else {
		return NewIntServerOption(ParamInt(intOption))
	}
}

// NewStringServerOption 新建String类型服务端配置
func NewStringServerOption(option ParamString) ServerOption {
	return ServerOption{StringOption: option, IsString: true}
}

// NewIntServerOption 新建Int类型服务端配置
func NewIntServerOption(option ParamInt) ServerOption {
	return ServerOption{IntOption: option, IsString: false}
}

// GetString 获取String类型配置，如果配置为Int类型则会报错
func (o *ServerOption) GetString() (ParamString, error) {
	if o.IsString {
		return o.StringOption, nil
	} else {
		return "", errors.New("配置为Int类型")
	}
}

// GetInt 获取Int类型配置，如果配置为String类型则会报错
func (o *ServerOption) GetInt() (ParamInt, error) {
	if o.IsString {
		return 0, errors.New("配置为String类型")
	} else {
		return o.IntOption, nil
	}
}

// GetLiteralValue 获取字面值，无论是String还是Int都会转换为字符串，方便前端显示
func (o *ServerOption) GetLiteralValue() string {
	if o.IsString {
		return string(o.StringOption)
	} else {
		return strconv.Itoa(int(o.IntOption))
	}
}

// SocketInfo Socket基本信息
type SocketInfo struct {
	SocketHandle SocketHandle // Socket句柄
	IpAddr       string       // IP地址
	Port         int32        // 端口号
	JobId        int32        // 连接最近处理的任务编号
	JobTime      DateTimeType // 最近处理任务的时间
	ConnectTime  DateTimeType // 客户端连接时间
	Timeout      DateTimeType // 连接超时时间
	Client       string       // 连接的客户端主机名称
	Process      string       // 连接的客户端程序名
	User         string       // 登录的用户
}

func getSocketInfo(handle ConnectHandle, nodeNumber int32, socket SocketHandle) (*SocketInfo, error) {
	connInfo, rte := RawRtdbGetConnectionInfoIpv6Warp(handle, nodeNumber, socket)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	ipAddr := connInfo.IpAddr6
	if ipAddr == "" {
		ipAddr = fmt.Sprintf("%d.%d.%d.%d", byte(connInfo.IpAddr>>24), byte(connInfo.IpAddr>>16), byte(connInfo.IpAddr>>8), byte(connInfo.IpAddr))
	}
	timeout, rte := RawRtdbGetTimeoutWarp(handle, socket)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	info := SocketInfo{
		SocketHandle: socket,
		IpAddr:       ipAddr,
		Port:         int32(connInfo.Port),
		JobId:        connInfo.Job,
		JobTime:      connInfo.JobTime,
		ConnectTime:  connInfo.ConnectTime,
		Timeout:      timeout,
		Client:       connInfo.Client,
		Process:      connInfo.Process,
		User:         connInfo.User,
	}
	return &info, nil
}

// NamedType 自定义类型
type NamedType struct {
	Name   string              // 自定义类型名称
	Fields []RtdbDataTypeField // 字段列表
	Desc   string              // 自定义类型描述
	Length int32               // 自定义类型长度(所有字段长度的累加和)
}

// ValueType 数值类型
type ValueType string

// 基本数值类型
const (
	// ValueTypeBool 布尔类型
	ValueTypeBool = ValueType("bool")

	// ValueTypeUint8 无符号8位整数
	ValueTypeUint8 = ValueType("uint8")

	// ValueTypeInt8 有符号8位整数
	ValueTypeInt8 = ValueType("int8")

	// ValueTypeChar 单字节字符
	ValueTypeChar = ValueType("char")

	// ValueTypeUint16 无符号16位整数
	ValueTypeUint16 = ValueType("uint16")

	// ValueTypeInt16 有符号16位整数
	ValueTypeInt16 = ValueType("int16")

	// ValueTypeUint32 无符号32位整数
	ValueTypeUint32 = ValueType("uint32")

	// ValueTypeInt32 有符号32位整数
	ValueTypeInt32 = ValueType("int32")

	// ValueTypeInt64 有符号64位整数
	ValueTypeInt64 = ValueType("int64")

	// ValueTypeFloat16 16位浮点数
	ValueTypeFloat16 = ValueType("float16")

	// ValueTypeFloat32 32位浮点数
	ValueTypeFloat32 = ValueType("float32")

	// ValueTypeFloat64 64位浮点数
	ValueTypeFloat64 = ValueType("float64")

	// ValueTypeCoor 二维坐标
	ValueTypeCoor = ValueType("coor")

	// ValueTypeString 字符串
	ValueTypeString = ValueType("string")

	// ValueTypeBlob 数据块
	ValueTypeBlob = ValueType("blob")

	// ValueTypeDatetime 时间
	ValueTypeDatetime = ValueType("datetime")

	// ValueTypeFp16 16位定点数
	ValueTypeFp16 = ValueType("fp16")

	// ValueTypeFp32 32位定点数
	ValueTypeFp32 = ValueType("fp32")

	// ValueTypeFp64 64位定点数
	ValueTypeFp64 = ValueType("fp64")
)

func (vt ValueType) ToRawType() (RtdbType, string) {
	switch vt {
	case ValueTypeBool:
		return RtdbTypeBool, "bool"
	case ValueTypeUint8:
		return RtdbTypeUint8, "uint8"
	case ValueTypeInt8:
		return RtdbTypeInt8, "int8"
	case ValueTypeChar:
		return RtdbTypeChar, "char"
	case ValueTypeUint16:
		return RtdbTypeUint16, "uint16"
	case ValueTypeInt16:
		return RtdbTypeInt16, "int16"
	case ValueTypeUint32:
		return RtdbTypeUint32, "uint32"
	case ValueTypeInt32:
		return RtdbTypeInt32, "int32"
	case ValueTypeInt64:
		return RtdbTypeInt64, "int64"
	case ValueTypeFloat16:
		return RtdbTypeReal16, "float16"
	case ValueTypeFloat32:
		return RtdbTypeReal32, "float32"
	case ValueTypeFloat64:
		return RtdbTypeReal64, "float64"
	case ValueTypeCoor:
		return RtdbTypeCoor, "coor"
	case ValueTypeString:
		return RtdbTypeString, "string"
	case ValueTypeBlob:
		return RtdbTypeBlob, "blob"
	case ValueTypeDatetime:
		return RtdbTypeDatetime, "datetime"
	case ValueTypeFp16:
		return RtdbTypeFp16, "fp16"
	case ValueTypeFp32:
		return RtdbTypeFp32, "fp32"
	case ValueTypeFp64:
		return RtdbTypeFp64, "fp64"
	default:
		return RtdbTypeNamedT, string(vt)
	}
}

func FromRawType(typ RtdbType, namedTypeName string) ValueType {
	switch typ {
	case RtdbTypeBool:
		return ValueTypeBool
	case RtdbTypeUint8:
		return ValueTypeUint8
	case RtdbTypeInt8:
		return ValueTypeInt8
	case RtdbTypeChar:
		return ValueTypeChar
	case RtdbTypeUint16:
		return ValueTypeUint16
	case RtdbTypeInt16:
		return ValueTypeInt16
	case RtdbTypeUint32:
		return ValueTypeUint32
	case RtdbTypeInt32:
		return ValueTypeInt32
	case RtdbTypeInt64:
		return ValueTypeInt64
	case RtdbTypeReal16:
		return ValueTypeFloat16
	case RtdbTypeReal32:
		return ValueTypeFloat32
	case RtdbTypeReal64:
		return ValueTypeFloat64
	case RtdbTypeCoor:
		return ValueTypeCoor
	case RtdbTypeString:
		return ValueTypeString
	case RtdbTypeBlob:
		return ValueTypeBlob
	case RtdbTypeNamedT:
		return ValueType(namedTypeName)
	case RtdbTypeDatetime:
		return ValueTypeDatetime
	case RtdbTypeFp16:
		return ValueTypeFp16
	case RtdbTypeFp32:
		return ValueTypeFp32
	case RtdbTypeFp64:
		return ValueTypeFp64
	default:
		panic("分支不可达")
	}
}

// PointClass 点类型
type PointClass int32

const (
	// PointBase 基本点
	PointBase = PointClass(RtdbClassBase)

	// PointScan 采集点
	PointScan = PointClass(RtdbClassBase | RtdbClassScan)

	// PointCalc 计算点
	PointCalc = PointClass(RtdbClassBase | RtdbClassCalc)

	// PointScanCalc 计算采集点
	PointScanCalc = PointClass(RtdbClassBase | RtdbClassScan | RtdbClassCalc)
)

func (p PointClass) Desc() string {
	switch p {
	case PointBase:
		return "基本点"
	case PointScan:
		return "采集点"
	case PointCalc:
		return "计算点"
	case PointScanCalc:
		return "采集计算点"
	default:
		return "未知点类型"
	}
}

// IsScan 是否为采集点
func (p PointClass) IsScan() bool {
	return p&PointScan != 0
}

// IsCalc 是否为计算点
func (p PointClass) IsCalc() bool {
	return p&PointCalc != 0
}

// PointInfo 点属性
type PointInfo struct {
	// 核心配置
	ID        PointID       // 标签点ID
	TableID   TableID       // 当前标签点所属表ID
	Name      string        // 标签点名称
	ValueType ValueType     // 数值类型
	Class     PointClass    // 标签点类别
	Precision RtdbPrecision // 时间戳精度

	// 基本点配置
	Desc           string     // 标签点描述
	Unit           string     // 工程单位
	Archive        Switch     // 是否存档
	Digits         int16      // 数值位数
	Shutdown       Switch     // 停机状态字
	LowLimit       float32    // 量程下限
	HighLimit      float32    // 量程上限
	Step           Switch     // 是否阶跃
	Typical        float32    // 典型值
	Compress       Switch     // 是否压缩
	CompDev        float32    // 压缩偏差
	CompDevPercent float32    // 压缩偏差百分比
	CompTimeMax    int32      // 最大压缩间隔
	CompTimeMin    int32      // 最短压缩间隔
	ExcDev         float32    // 例外偏差
	ExcDevPercent  float32    // 例外偏差百分比
	ExcTimeMax     int32      // 最大例外间隔
	ExcTimeMin     int32      // 最短例外间隔
	Mirror         RtdbMirror // 镜像收发控制
	Summary        Switch     // 统计加速

	// 采集点配置，仅采集点有效
	Source     string                         // 数据源
	Scan       Switch                         // 是否采集
	Instrument string                         // 设备标签
	Locations  [RtdbConstLocationsSize]int32  // 共包含五个设备位址
	UserInts   [RtdbConstUserintSize]int32    // 共包含两个自定义整数
	UserReals  [RtdbConstUserrealSize]float32 // 共包含两个自定义单精度浮点数

	// 计算点配置, 仅计算点有效
	Equation string       // 实时方程式
	Trigger  RtdbTrigger  // 计算触发机制
	TimeCopy RtdbTimeCopy // 计算结果时间戳参考
	Period   int32        // 触发周期

	// 只读属性
	NamedType   NamedType    // 自定义类型结构, 仅自定义类型有效
	TableDotTag string       // 标签点全名，格式为“表名称.标签点名称”
	ChangeDate  DateTimeType // 标签点属性最后一次被修改的时间
	Changer     string       // 标签点属性最后一次被修改的用户名
	CreateDate  DateTimeType // 标签点被创建的时间
	Creator     string       // 标签点创建者的用户名
}

// NewPointInfo 新建标签点属性, 备注：只需填写必要属性，其他都是默认，需要时可自行设置
//
// input:
//   - name 点名
//   - tableId 表ID
//   - valueType 数值类型
//   - class 点类型
//   - precision 点时间戳精度
//   - unit 点单位
//   - desc 点描述
func NewPointInfo(name string, tableId TableID, valueType ValueType, class PointClass, precision RtdbPrecision, unit, desc string) *PointInfo {
	return &PointInfo{
		Name:           name,
		ValueType:      valueType,
		TableID:        tableId,
		Class:          class,
		Unit:           unit,
		Desc:           desc,
		Archive:        ON,
		Digits:         -5,
		Shutdown:       OFF,
		LowLimit:       0,
		HighLimit:      100,
		Step:           OFF,
		Typical:        50,
		Compress:       OFF,
		CompDev:        1,
		CompDevPercent: 0,
		CompTimeMax:    28800,
		CompTimeMin:    0,
		ExcDev:         0.5,
		ExcDevPercent:  0,
		ExcTimeMax:     600,
		ExcTimeMin:     0,
		Mirror:         RtdbMirrorPointOff,
		Summary:        OFF,
		Precision:      precision,
	}
}

// SetLimit 设置量程上下限
//
// input:
//   - lowLimit 量程上限
//   - highLimit 量程下限
//   - typical 典型值(默认值)
func (p *PointInfo) SetLimit(lowLimit float32, highLimit float32, typical float32) {
	p.LowLimit = lowLimit
	p.HighLimit = highLimit
	p.Typical = typical
}

// SetStoreDisplay 设置存储显示
//
// input:
//   - archive 是否存档
//   - digits 数值显示位数
//   - shutdown 是否停机补写
//   - step 是否阶跃
//   - mirror 镜像配置
//   - summary 是否开启统计加速
func (p *PointInfo) SetStoreDisplay(archive Switch, digits int16, shutdown Switch, step Switch, mirror RtdbMirror, summary Switch) {
	p.Archive = archive
	p.Digits = digits
	p.Shutdown = shutdown
	p.Step = step
	p.Mirror = mirror
	p.Summary = summary
}

// SetCompress 设置压缩
//
// input:
//   - compress 是否压缩
//   - compDev 压缩偏差
//   - compDevPercent 压缩偏差百分比
//   - compTimeMax 最大压缩间隔
//   - compTimeMin 最小压缩间隔
func (p *PointInfo) SetCompress(compress Switch, compDev float32, compDevPercent float32, compTimeMax int32, compTimeMin int32) {
	p.Compress = compress
	p.CompDev = compDev
	p.CompDevPercent = compDevPercent
	p.CompTimeMax = compTimeMax
	p.CompTimeMin = compTimeMin
}

// SetException 设置例外偏差
//
// input:
//   - excDev 例外偏差
//   - excDevPercent 例外偏差百分比
//   - excTimeMax 最大例外间隔
//   - excTimeMin 最短例外间隔
func (p *PointInfo) SetException(excDev float32, excDevPercent float32, excTimeMax int32, excTimeMin int32) {
	p.ExcDev = excDev
	p.ExcDevPercent = excDevPercent
	p.ExcTimeMax = excTimeMax
	p.ExcTimeMin = excTimeMin
}

// SetScan 设置采集点属性
//
// input:
//   - source 数据源
//   - scan 是否采集
//   - instrument 设备标签
//   - locations 共包含五个设备地址
//   - userInts 共包含两个自定义整数
//   - userReals 共包含两个自定义单精度浮点数
func (p *PointInfo) SetScan(
	source string, scan Switch, instrument string, locations [RtdbConstLocationsSize]int32,
	userInts [RtdbConstUserintSize]int32, userReals [RtdbConstUserrealSize]float32,
) {
	p.Class |= PointClass(RtdbClassScan)
	p.Source = source
	p.Scan = scan
	p.Instrument = instrument
	p.Locations = locations
	p.UserInts = userInts
	p.UserReals = userReals
}

// SetCalc 设置计算点
//
// input:
//   - equation 实时方程式
//   - trigger 计算触发机制
//   - timeCopy 计算结果时间戳参考
//   - period 触发周期
func (p *PointInfo) SetCalc(equation string, trigger RtdbTrigger, timeCopy RtdbTimeCopy, period int32) {
	p.Class |= PointClass(RtdbClassCalc)
	p.Equation = equation
	p.Trigger = trigger
	p.TimeCopy = timeCopy
	p.Period = period
}

// NewTVQ 新建TVQ
//
// input:
//   - timestamp: 时间戳
//   - value: 数值， 这个是泛型的，value类型需和p.Type相同
//   - quality: 质量码
//
// output:
//   - TVQ(tvq) TVQ值
func (p *PointInfo) NewTVQ(timestamp time.Time, value any, quality Quality) TVQ {
	rawType, name := p.ValueType.ToRawType()
	switch rawType {
	case RtdbTypeBool:
		return NewTvqBool(timestamp, value.(bool), quality)
	case RtdbTypeUint8:
		return NewTvqUint8(timestamp, value.(uint8), quality)
	case RtdbTypeInt8:
		return NewTvqInt8(timestamp, value.(int8), quality)
	case RtdbTypeChar:
		return NewTvqChar(timestamp, value.(byte), quality)
	case RtdbTypeUint16:
		return NewTvqUint16(timestamp, value.(uint16), quality)
	case RtdbTypeInt16:
		return NewTvqInt16(timestamp, value.(int16), quality)
	case RtdbTypeUint32:
		return NewTvqUint32(timestamp, value.(uint32), quality)
	case RtdbTypeInt32:
		return NewTvqInt32(timestamp, value.(int32), quality)
	case RtdbTypeInt64:
		return NewTvqInt64(timestamp, value.(int64), quality)
	case RtdbTypeReal16:
		return NewTvqFloat16(timestamp, value.(float32), quality)
	case RtdbTypeReal32:
		return NewTvqFloat32(timestamp, value.(float32), quality)
	case RtdbTypeReal64:
		return NewTvqFloat64(timestamp, value.(float64), quality)
	case RtdbTypeCoor:
		coordinates := value.(Coordinates)
		return NewTvqCoordinates(timestamp, coordinates.X, coordinates.Y, quality)
	case RtdbTypeString:
		return NewTvqString(timestamp, value.(string), quality)
	case RtdbTypeBlob:
		return NewTvqBlob(timestamp, value.([]byte), quality)
	case RtdbTypeDatetime:
		return NewTvqDatetime(timestamp, value.(string), quality)
	case RtdbTypeFp16:
		return NewTvqFp16(timestamp, value.(float32), quality)
	case RtdbTypeFp32:
		return NewTvqFp32(timestamp, value.(float32), quality)
	case RtdbTypeFp64:
		return NewTvqFp64(timestamp, value.(float64), quality)
	default:
		return NewTvqNamed(timestamp, ValueType(name), value.([]byte), quality)
	}
}

// NewPTVQ 新建PTVQ
//
// input:
//   - timestamp: 时间戳
//   - value: 数值， 这个是泛型的，value类型需和p.Type相同
//   - quality: 质量码
//
// output:
//   - PTVQ(tvq) PTVQ值
func (p *PointInfo) NewPTVQ(timestamp time.Time, value any, quality Quality) PTVQ {
	tvq := p.NewTVQ(timestamp, value, quality)
	return NewPTVQ(p, tvq)
}

// NewNowTVQ 新建TVQ, 使用当前时间戳
//
// input:
//   - value: 数值， 这个是泛型的，value类型需和p.Type相同
//   - quality: 质量码
//
// output:
//   - TVQ(tvq) TVQ值
func (p *PointInfo) NewNowTVQ(value any, quality Quality) TVQ {
	return p.NewTVQ(time.Now(), value, quality)
}

// PointInfoToRaw 点信息转换为Raw点属性表
func PointInfoToRaw(info *PointInfo) (*RtdbPoint, *RtdbScan, *RtdbCalc, string) {
	rtdbType, namedTypeName := info.ValueType.ToRawType()
	milliSecond := int8(0)
	if info.Precision != RtdbPrecisionSecond {
		milliSecond = 1
	}
	base := &RtdbPoint{
		ID:             info.ID,
		Tag:            info.Name,
		Type:           rtdbType,
		Table:          info.TableID,
		Desc:           info.Desc,
		Unit:           info.Unit,
		Archive:        info.Archive,
		Digits:         info.Digits,
		Shutdown:       info.Shutdown,
		LowLimit:       info.LowLimit,
		HighLimit:      info.HighLimit,
		Step:           info.Step,
		Typical:        info.Typical,
		Compress:       info.Compress,
		CompDev:        info.CompDev,
		CompDevPercent: info.CompDevPercent,
		CompTimeMax:    info.CompTimeMax,
		CompTimeMin:    info.CompTimeMin,
		ExcDev:         info.ExcDev,
		ExcDevPercent:  info.ExcDevPercent,
		ExcTimeMin:     info.ExcTimeMin,
		ExcTimeMax:     info.ExcTimeMax,
		Class:          RtdbClass(info.Class),
		Mirror:         info.Mirror,
		Summary:        info.Summary,
		Precision:      info.Precision,
		MilliSecond:    milliSecond,
	}
	scan := (*RtdbScan)(nil)
	if info.Class.IsScan() {
		scan = &RtdbScan{
			Source:     info.Source,
			Scan:       info.Scan,
			Instrument: info.Instrument,
			Locations:  info.Locations,
			UserInts:   info.UserInts,
			UserReals:  info.UserReals,
		}
	}
	calc := (*RtdbCalc)(nil)
	if info.Class.IsCalc() {
		calc = &RtdbCalc{
			Equation: info.Equation,
			Trigger:  info.Trigger,
			TimeCopy: info.TimeCopy,
			Period:   info.Period,
		}
	}
	return base, scan, calc, namedTypeName
}

// PointInfoFromRaw 点属性表转换为点信息
func PointInfoFromRaw(handle ConnectHandle, base *RtdbPoint, scan *RtdbScan, calc *RtdbCalc, isRecycled bool) (*PointInfo, error) {
	typ := (*NamedType)(nil)
	if base.Type == RtdbTypeNamedT {
		if !isRecycled {
			names, counts, rtes, rte := RawRtdbbGetNamedTypeNamesPropertyWarp(handle, []PointID{base.ID})
			if !RteIsOk(rte) {
				return nil, rte.GoError()
			}
			if !RteIsOk(rtes[0]) {
				return nil, rte.GoError()
			}
			fields, tLen, desc, rte := RawRtdbbGetNamedTypeWarp(handle, names[0], counts[0])
			typ = &NamedType{Name: names[0], Fields: fields, Desc: desc, Length: tLen}
		} else {
			names, counts, rtes, rte := RawRtdbbGetRecycledNamedTypeNamesPropertyWarp(handle, []PointID{base.ID})
			if !RteIsOk(rte) {
				return nil, rte.GoError()
			}
			if !RteIsOk(rtes[0]) {
				return nil, rte.GoError()
			}
			fields, tLen, desc, rte := RawRtdbbGetNamedTypeWarp(handle, names[0], counts[0])
			typ = &NamedType{Name: names[0], Fields: fields, Desc: desc, Length: tLen}
		}
	}
	namedTypeName := ""
	if typ != nil {
		namedTypeName = typ.Name
	}
	info := &PointInfo{
		ID:             base.ID,
		TableID:        base.Table,
		Name:           base.Tag,
		ValueType:      FromRawType(base.Type, namedTypeName),
		Class:          PointClass(base.Class),
		Precision:      base.Precision,
		Desc:           base.Desc,
		Unit:           base.Unit,
		Archive:        base.Archive,
		Digits:         base.Digits,
		Shutdown:       base.Shutdown,
		LowLimit:       base.LowLimit,
		HighLimit:      base.HighLimit,
		Step:           base.Step,
		Typical:        base.Typical,
		Compress:       base.Compress,
		CompDev:        base.CompDev,
		CompDevPercent: base.CompDevPercent,
		CompTimeMax:    base.CompTimeMax,
		CompTimeMin:    base.CompTimeMin,
		ExcDev:         base.ExcDev,
		ExcDevPercent:  base.ExcDevPercent,
		ExcTimeMax:     base.ExcTimeMax,
		ExcTimeMin:     base.ExcTimeMin,
		Mirror:         base.Mirror,
		Summary:        base.Summary,
		TableDotTag:    base.TableDotTag,
		ChangeDate:     base.ChangeDate,
		Changer:        base.Changer,
		CreateDate:     base.CreateDate,
		Creator:        base.Creator,
	}
	if typ != nil {
		info.NamedType = *typ
	}
	if scan != nil {
		info.Source = scan.Source
		info.Scan = scan.Scan
		info.Instrument = scan.Instrument
		info.Locations = scan.Locations
		info.UserInts = scan.UserInts
		info.UserReals = scan.UserReals
	}
	if calc != nil {
		info.Equation = calc.Equation
		info.Trigger = calc.Trigger
		info.TimeCopy = calc.TimeCopy
		info.Period = calc.Period
	}

	return info, nil
}

type PointInfoField string

const (
	// PointInfoFieldName 标签点名称
	PointInfoFieldName = PointInfoField("name")

	// PointInfoFieldClass 标签点类别
	PointInfoFieldClass = PointInfoField("class")

	// PointInfoFieldDesc 标签点描述
	PointInfoFieldDesc = PointInfoField("desc")

	// PointInfoFieldUnit 标签点单位
	PointInfoFieldUnit = PointInfoField("unit")

	// PointInfoFieldArchive 是否存档
	PointInfoFieldArchive = PointInfoField("archive")

	// PointInfoFieldDigits 数值位数
	PointInfoFieldDigits = PointInfoField("digits")

	// PointInfoFieldShutdown 停机状态字
	PointInfoFieldShutdown = PointInfoField("shutdown")

	// PointInfoFieldLowLimit 量程下限
	PointInfoFieldLowLimit = PointInfoField("low_limit")

	// PointInfoFieldHighLimit 量程上限
	PointInfoFieldHighLimit = PointInfoField("high_limit")

	// PointInfoFieldStep 是否阶跃
	PointInfoFieldStep = PointInfoField("step")

	// PointInfoFieldTypical 典型值
	PointInfoFieldTypical = PointInfoField("typical")

	// PointInfoFieldCompress 是否压缩
	PointInfoFieldCompress = PointInfoField("compress")

	// PointInfoFieldCompDev 压缩偏差
	PointInfoFieldCompDev = PointInfoField("comp_dev")

	// PointInfoFieldCompDevPercent 压缩偏差百分比
	PointInfoFieldCompDevPercent = PointInfoField("comp_dev_percent")

	// PointInfoFieldCompTimeMax 最大压缩间隔
	PointInfoFieldCompTimeMax = PointInfoField("comp_time_max")

	// PointInfoFieldCompTimeMin 最小压缩间隔
	PointInfoFieldCompTimeMin = PointInfoField("comp_time_min")

	// PointInfoFieldExcDev 例外偏差
	PointInfoFieldExcDev = PointInfoField("exc_dev")

	// PointInfoFieldExcDevPercent 例外偏差百分比
	PointInfoFieldExcDevPercent = PointInfoField("exc_dev_percent")

	// PointInfoFieldExcTimeMax 最大例外间隔
	PointInfoFieldExcTimeMax = PointInfoField("exc_time_max")

	// PointInfoFieldExcTimeMin 最小例外间隔
	PointInfoFieldExcTimeMin = PointInfoField("exc_time_min")

	// PointInfoFieldMirror 镜像收发控制
	PointInfoFieldMirror = PointInfoField("mirror")

	// PointInfoFieldSummary 统计加速
	PointInfoFieldSummary = PointInfoField("summary")

	// PointInfoFieldSource 数据源
	PointInfoFieldSource = PointInfoField("source")

	// PointInfoFieldScan 是否采集
	PointInfoFieldScan = PointInfoField("scan")

	// PointInfoFieldInstrument 设备标签
	PointInfoFieldInstrument = PointInfoField("instrument")

	// PointInfoFieldLocations 共包含五个设备位址
	PointInfoFieldLocations = PointInfoField("locations")

	// PointInfoFieldUserInts 共包含两个自定义整数
	PointInfoFieldUserInts = PointInfoField("user_ints")

	// PointInfoFieldUserReals 共包含两个自定义单精度浮点数
	PointInfoFieldUserReals = PointInfoField("user_reals")

	// PointInfoFieldEquation 实时方程式
	PointInfoFieldEquation = PointInfoField("equation")

	// PointInfoFieldTrigger 计算触发机制
	PointInfoFieldTrigger = PointInfoField("trigger")

	// PointInfoFieldTimeCopy 计算结果时间戳参考
	PointInfoFieldTimeCopy = PointInfoField("time_copy")

	// PointInfoFieldPeriod 触发周期
	PointInfoFieldPeriod = PointInfoField("period")
)

// Coordinates 坐标类型Value
type Coordinates struct {
	X float32
	Y float32
}

// AnyValue 任意Value, 理想状态下这是一个内联语意(union)，但由于Golang的语法限制，因此这里只能是一个结构
type AnyValue struct {
	IntValue         int64
	FloatValue       float64
	CoordinatesValue Coordinates
	StringValue      string
	BytesValue       []byte
}

// TVQ Timestamp + Value + Quality, 单个点值, 兼容所有类型
type TVQ struct {
	Timestamp time.Time // 时间戳
	Type      ValueType // 数值类型
	Value     AnyValue  // 数值
	Quality   Quality   // 质量码
}

// NewTvqBool 新建一个bool类型的TVQ
func NewTvqBool(timestamp time.Time, value bool, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeBool,
		Value:     AnyValue{IntValue: BoolToInt64(value)},
		Quality:   quality,
	}
}

// NewTvqUint8 新建一个uint8类型的TVQ
func NewTvqUint8(timestamp time.Time, value uint8, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeUint8,
		Value:     AnyValue{IntValue: int64(value)},
		Quality:   quality,
	}
}

// NewTvqInt8 新建一个int8类型的TVQ
func NewTvqInt8(timestamp time.Time, value int8, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeInt8,
		Value:     AnyValue{IntValue: int64(value)},
		Quality:   quality,
	}
}

// NewTvqChar 新建一个char类型的TVQ
func NewTvqChar(timestamp time.Time, value byte, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeChar,
		Value:     AnyValue{IntValue: int64(value)},
		Quality:   quality,
	}
}

// NewTvqUint16 新建一个uint16类型的TVQ
func NewTvqUint16(timestamp time.Time, value uint16, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeUint16,
		Value:     AnyValue{IntValue: int64(value)},
		Quality:   quality,
	}
}

// NewTvqInt16 新建一个int16类型的TVQ
func NewTvqInt16(timestamp time.Time, value int16, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeInt16,
		Value:     AnyValue{IntValue: int64(value)},
		Quality:   quality,
	}
}

// NewTvqUint32 新建一个uint32类型的TVQ
func NewTvqUint32(timestamp time.Time, value uint32, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeUint32,
		Value:     AnyValue{IntValue: int64(value)},
		Quality:   quality,
	}
}

// NewTvqInt32 新建一个int32类型的TVQ
func NewTvqInt32(timestamp time.Time, value int32, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeInt32,
		Value:     AnyValue{IntValue: int64(value)},
		Quality:   quality,
	}
}

// NewTvqInt64 新建一个int64类型的TVQ
func NewTvqInt64(timestamp time.Time, value int64, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeInt64,
		Value:     AnyValue{IntValue: value},
		Quality:   quality,
	}
}

// NewTvqFloat16 新建一个float16类型的TVQ (由于golang中没有float16，所以这里用的是float32)
func NewTvqFloat16(timestamp time.Time, value float32, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeFloat16,
		Value:     AnyValue{FloatValue: float64(value)},
		Quality:   quality,
	}
}

// NewTvqFloat32 新建一个float32类型的TVQ
func NewTvqFloat32(timestamp time.Time, value float32, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeFloat32,
		Value:     AnyValue{FloatValue: float64(value)},
		Quality:   quality,
	}
}

// NewTvqFloat64 新建一个float64类型的TVQ
func NewTvqFloat64(timestamp time.Time, value float64, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeFloat64,
		Value:     AnyValue{FloatValue: value},
		Quality:   quality,
	}
}

// NewTvqCoordinates 新建一个坐标(x,y)类型的TVQ
func NewTvqCoordinates(timestamp time.Time, x, y float32, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeCoor,
		Value:     AnyValue{CoordinatesValue: Coordinates{X: x, Y: y}},
		Quality:   quality,
	}
}

// NewTvqString 新建一个string类型的TVQ
func NewTvqString(timestamp time.Time, str string, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeString,
		Value:     AnyValue{StringValue: str},
		Quality:   quality,
	}
}

// NewTvqBlob 新建一个blob([]byte)类型的TVQ
func NewTvqBlob(timestamp time.Time, blob []byte, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeBlob,
		Value:     AnyValue{BytesValue: blob},
		Quality:   quality,
	}
}

// NewTvqDatetime 新建一个时间类型的TVQ
func NewTvqDatetime(timestamp time.Time, datetime string, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeDatetime,
		Value:     AnyValue{StringValue: datetime},
		Quality:   quality,
	}
}

// NewTvqFp16 新建一个FP16的TVQ (定点数还没对外开放, 不建议使用)
func NewTvqFp16(timestamp time.Time, value float32, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeFp16,
		Value:     AnyValue{FloatValue: float64(value)},
		Quality:   quality,
	}
}

// NewTvqFp32 新建一个FP32的TVQ (定点数还没对外开放, 不建议使用)
func NewTvqFp32(timestamp time.Time, value float32, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeFp32,
		Value:     AnyValue{FloatValue: float64(value)},
		Quality:   quality,
	}
}

// NewTvqFp64 新建一个FP64的TVQ (定点数还没对外开放, 不建议使用)
func NewTvqFp64(timestamp time.Time, value float64, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      ValueTypeFp64,
		Value:     AnyValue{FloatValue: value},
		Quality:   quality,
	}
}

// NewTvqNamed 新建一个自定义类型的TVQ
func NewTvqNamed(timestamp time.Time, valueType ValueType, data []byte, quality Quality) TVQ {
	return TVQ{
		Timestamp: timestamp,
		Type:      valueType,
		Value:     AnyValue{BytesValue: data},
		Quality:   quality,
	}
}

// GetRtdbTimestamp 获取时间戳
func (v *TVQ) GetRtdbTimestamp() (TimestampType, SubtimeType) {
	return GoTimeToRtdbTimestamp(v.Timestamp)
}

// GoTimeToRtdbTimestamp time.Time 转换成 (TimestampType, SubtimeType)
func GoTimeToRtdbTimestamp(timestamp time.Time) (TimestampType, SubtimeType) {
	datetime := TimestampType(timestamp.Unix())
	subtime := SubtimeType(timestamp.Nanosecond())
	return datetime, subtime
}

// RtdbTimestampToGoTime (TimestampType, SubtimeType) 转换成 time.Time
func RtdbTimestampToGoTime(datetime TimestampType, subtime SubtimeType) time.Time {
	return time.Unix(int64(datetime), int64(subtime))
}

// GetRtdbType 获取数值类型
func (v *TVQ) GetRtdbType() (RtdbType, string) {
	return v.Type.ToRawType()
}

// GetRtdbInt 获取Int类型数值
func (v *TVQ) GetRtdbInt() int64 {
	return v.Value.IntValue
}

// GetRtdbFloat 获取Float类型数值
func (v *TVQ) GetRtdbFloat() float64 {
	return v.Value.FloatValue
}

// GetRtdbCoordinates 获取坐标类型数值
func (v *TVQ) GetRtdbCoordinates() Coordinates {
	return v.Value.CoordinatesValue
}

// GetRtdbString 获取string类型数值
func (v *TVQ) GetRtdbString() string {
	return v.Value.StringValue
}

// GetRtdbBlob 获取[]byte类型数值
func (v *TVQ) GetRtdbBlob() []byte {
	return v.Value.BytesValue
}

// GetRtdbStringBlob 获取string/blob类型数值
func (v *TVQ) GetRtdbStringBlob(osType RtdbOsType) ([]byte, error) {
	var data []byte
	rtdbType, _ := v.Type.ToRawType()
	if rtdbType == RtdbTypeString {
		str := v.GetRtdbString()
		if osType == RtdbOsLinux {
			data = []byte(str)
		} else if osType == RtdbOsWindows {
			buf, err := StringToGBKBytes(str)
			if err != nil {
				return nil, err
			}
			data = buf
		} else {
			panic("分支不可达, 未知的OsType")
		}
	} else if rtdbType == RtdbTypeBlob {
		data = v.GetRtdbBlob()
	} else {
		panic("分支不可达")
	}
	return data, nil
}

// GetRtdbDatetime 获取日期类型数值
func (v *TVQ) GetRtdbDatetime() string {
	return v.Value.StringValue
}

// GetRtdbNamedObj 获取自定义类型数值
func (v *TVQ) GetRtdbNamedObj() []byte {
	return v.Value.BytesValue
}

// GetRtdbQuality 获取质量码
func (v *TVQ) GetRtdbQuality() Quality {
	return v.Quality
}

// PTVQ PointInfo + TVQ
type PTVQ struct {
	PointInfo *PointInfo // 点信息
	TVQ       TVQ        // timestamp + value + quality
}

// NewPTVQ 新建PTVQ
//
// input:
//   - info 点信息
//   - tvq timestamp + value + quality
//
// output:
//   - PTVQ(ptvq) PTVQ值
func NewPTVQ(info *PointInfo, tvq TVQ) PTVQ {
	return PTVQ{PointInfo: info, TVQ: tvq}
}

// PTVQs PointInfo + []TVQ
type PTVQs struct {
	PointInfo *PointInfo // 点信息
	TVQs      []TVQ      // timestamp + value + quality
}

// NewPTVQs 新建PTVQs
//
// input:
//   - info 点信息
//   - tvqs []TVQ
//
// output:
//   - PTVQs(ptvqs) PTVQs值
func NewPTVQs(info *PointInfo, tvqs []TVQ) PTVQs {
	return PTVQs{info, tvqs}
}

// ArchivedBaseInfo 存档文件基本信息
type ArchivedBaseInfo struct {
	Path        string              // 存档文件路径
	File        string              // 存档文件名称
	State       RtdbArchiveState    // 存档文件状态
	HeadPage    RtdbHeaderPage      // 存档首部信息
	PerfInfo    RtdbArchivePerfData // 存档性能计数
	SumPerfInfo RtdbArchivePerfData // 存档性能技术合计
}

////////////////////////////////////////////////
//////////////////上面是一些结构//////////////////
////////////////////摆烂的分隔线/////////////////
/////////////////下面是RtdbConnect函数///////////
////////////////////////////////////////////////

type RtdbConnect struct {
	HostIp                  string         // 服务端名称
	Port                    int32          // 服务端端口
	UserName                string         // 用户名
	Password                string         // 密码
	ConnectHandle           ConnectHandle  // 连接句柄
	Priv                    PrivGroup      // 用户权限
	SyncInfos               []RtdbSyncInfo // 元数据信息
	SocketHandles           []SocketHandle // 套接字句柄
	ServerOsType            RtdbOsType     // 服务端操作系统类型
	StringBlobMaxLen        int32          // 最大支持String/Blob长度
	SubscribeTagsConn       *RtdbConnect   // 订阅标签点属性更新专用连接
	SubscribeTagsName       string         // 订阅标签点属性channel名称
	SubscribeTagsParam      unsafe.Pointer // 订阅标签点属性channel(c指针)
	SubscribeSnapshotsConn  *RtdbConnect   // 订阅快照专用连接
	SubscribeSnapshotsName  string         // 订阅快照channel名称
	SubscribeSnapshotsParam unsafe.Pointer // 订阅快照channel(c指针)
	Location                *time.Location // 时区
}

// Login 登录数据库
//
// input:
//   - hostIp 主机IP
//   - port 端口
//   - userName 用户名
//   - password 密码
//
// output:
//   - RtdbConnect(conn) 返回数据库连接
func Login(hostIp string, port int32, userName string, password string, serverPrecision RtdbPrecision) (*RtdbConnect, error) {
	rtn := RtdbConnect{
		HostIp:   hostIp,
		Port:     port,
		UserName: userName,
		Password: password,
	}

	// 设置默认时间戳精度为纳秒
	rte := RawRtdbSetOptionWarp(RtdbApiOptionDefaultPrecision, 3)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}

	// 设置服务端时间戳精度
	rte = RawRtdbSetOptionWarp(RtdbApiOptionServerPrecision, int32(serverPrecision))
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}

	// 连接数据库
	cHandle, rte := RawRtdbConnectWarp(rtn.HostIp, rtn.Port)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.ConnectHandle = cHandle

	// 登录数据库
	priv, rte := RawRtdbLoginWarp(rtn.ConnectHandle, rtn.UserName, rtn.Password)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.Priv = priv

	// 获取元信息
	infos, errs, rte := RawRtdbbGetMetaSyncInfoWarp(rtn.ConnectHandle, 0)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	for _, rte := range errs {
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
	}
	rtn.SyncInfos = infos

	// 获取套接字句柄
	for i := range infos {
		sHandle, rte := RawRtdbGetOwnConnectionWarp(rtn.ConnectHandle, int32(i+1))
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		rtn.SocketHandles = append(rtn.SocketHandles, sHandle)
	}

	// 获取服务器操作系统类型
	osType, rte := RawRtdbOsType(rtn.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.ServerOsType = osType

	// 获取String/Blob最大长度
	maxLen, rte := RawRtdbGetMaxBlobLenWarp(rtn.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	rtn.StringBlobMaxLen = maxLen

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	rtn.Location = loc

	return &rtn, nil
}

// SetLocation 设置客户端时区，默认为亚洲东八
func (c *RtdbConnect) SetLocation(locName string) error {
	loc, err := time.LoadLocation(locName)
	if err != nil {
		return nil
	}
	c.Location = loc
	return nil
}

// Logout 登出数据库
func (c *RtdbConnect) Logout() error {
	rte := RawRtdbDisconnectWarp(c.ConnectHandle)
	return rte.GoError()
}

// GetClientVersion 获取客户端版本
//
// output:
//   - ApiVersion(version) 客户端版本
func (c *RtdbConnect) GetClientVersion() (*ApiVersion, error) {
	version, rte := RawRtdbGetApiVersionWarp()
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &version, rte.GoError()
}

// SetClientOption 设置客户端参数
//
// input:
//   - option: 客户端参数选项
//   - value: 客户端参数值
func (c *RtdbConnect) SetClientOption(option RtdbApiOption, value int32) error {
	rte := RawRtdbSetOptionWarp(option, value)
	return rte.GoError()
}

// GetServerOption 获取服务端参数
//
// input:
//   - param 服务端参数选项
//
// output:
//   - ServerOption(option) 服务端参数值
func (c *RtdbConnect) GetServerOption(param RtdbParam) (*ServerOption, error) {
	if param.IsStringParam() {
		opt, rte := RawRtdbGetDbInfo1Warp(c.ConnectHandle, param)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return &ServerOption{StringOption: opt, IsString: true}, nil
	} else {
		opt, rte := RawRtdbGetDbInfo2Warp(c.ConnectHandle, param)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return &ServerOption{IntOption: opt, IsString: false}, nil
	}
}

// SetServerOption 设置服务端参数
//
// input:
//   - param 服务端参数选项
//   - option 服务端参数值
func (c *RtdbConnect) SetServerOption(param RtdbParam, option ServerOption) error {
	if param.IsStringParam() {
		strOpt, err := option.GetString()
		if err != nil {
			return err
		}
		rte := RawRtdbSetDbInfo1Warp(c.ConnectHandle, param, strOpt)
		return rte.GoError()
	} else {
		intOpt, err := option.GetInt()
		if err != nil {
			return err
		}
		rte := RawRtdbSetDbInfo2Warp(c.ConnectHandle, param, intOpt)
		return rte.GoError()
	}
}

// GetSocketInfos 获取服务端SocketInfo列表，单机服务端返回一个SocketInfo列表，双活服务端返回两个SocketInfo列表
//
// output:
//   - [][]SocketInfo(infos) Socket信息列表
func (c *RtdbConnect) GetSocketInfos() ([][]SocketInfo, error) {
	if len(c.SyncInfos) == 1 { /* 单机,返回一个Socket列表 */
		count, rte := RawRtdbConnectionCountWarp(c.ConnectHandle, 0)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		sockets, rte := RawRtdbGetConnectionsWarp(c.ConnectHandle, 0, count)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}

		infos := make([]SocketInfo, 0)
		for _, socket := range sockets {
			info, err := getSocketInfo(c.ConnectHandle, 0, socket)
			if err != nil {
				return nil, err
			}
			infos = append(infos, *info)
		}
		return [][]SocketInfo{infos}, nil
	} else { /* 双活,返回两个Socket列表 */
		count1, rte := RawRtdbConnectionCountWarp(c.ConnectHandle, 1)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		sockets1, rte := RawRtdbGetConnectionsWarp(c.ConnectHandle, 1, count1)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		infos1 := make([]SocketInfo, 0)
		for _, socket := range sockets1 {
			info, err := getSocketInfo(c.ConnectHandle, 1, socket)
			if err != nil {
				return nil, err
			}
			infos1 = append(infos1, *info)
		}

		count2, rte := RawRtdbConnectionCountWarp(c.ConnectHandle, 2)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		sockets2, rte := RawRtdbGetConnectionsWarp(c.ConnectHandle, 2, count2)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		infos2 := make([]SocketInfo, 0)
		for _, socket := range sockets2 {
			info, err := getSocketInfo(c.ConnectHandle, 2, socket)
			if err != nil {
				return nil, err
			}
			infos2 = append(infos2, *info)
		}

		return [][]SocketInfo{infos1, infos2}, nil
	}
}

// GetOwnSocketInfo 获取当前连接的SocketInfo，单机服务端返回一个SocketInfo，双活服务端返回两个SocketInfo
//
// output:
//   - []Socket Socket信息
func (c *RtdbConnect) GetOwnSocketInfo() ([]SocketInfo, error) {
	if len(c.SyncInfos) == 1 { /* 单机,返回一个Socket句柄 */
		socket, rte := RawRtdbGetOwnConnectionWarp(c.ConnectHandle, 0)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		info, err := getSocketInfo(c.ConnectHandle, 0, socket)
		if err != nil {
			return nil, err
		}
		return []SocketInfo{*info}, nil
	} else { /* 双活,返回两个Socket句柄 */
		socket1, rte := RawRtdbGetOwnConnectionWarp(c.ConnectHandle, 1)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		info1, err := getSocketInfo(c.ConnectHandle, 1, socket1)
		if err != nil {
			return nil, err
		}
		socket2, rte := RawRtdbGetOwnConnectionWarp(c.ConnectHandle, 2)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		info2, err := getSocketInfo(c.ConnectHandle, 2, socket2)
		if err != nil {
			return nil, err
		}
		return []SocketInfo{*info1, *info2}, nil
	}
}

// SetSocketTimeout 设置Socket超时时间
//
// input:
//   - info Socket信息结构
//   - timeout 超时时间
func (c *RtdbConnect) SetSocketTimeout(info SocketInfo, timeout DateTimeType) error {
	rte := RawRtdbSetTimeoutWarp(c.ConnectHandle, info.SocketHandle, timeout)
	return rte.GoError()
}

// KillSocket 断开Socket
//
// input:
//   - info Socket信息结构
func (c *RtdbConnect) KillSocket(info SocketInfo) error {
	rte := RawRtdbKillConnectionWarp(c.ConnectHandle, info.SocketHandle)
	return rte.GoError()
}

// AddIpBlackList 添加IP黑名单项
//
// input:
//   - address 阻止连接段地址
//   - mask 阻止连接段子网掩码
//   - desc 阻止连接段的说明
func (c *RtdbConnect) AddIpBlackList(address string, mask string, desc string) error {
	rte := RawRtdbAddBlacklistWarp(c.ConnectHandle, address, mask, desc)
	return rte.GoError()
}

// UpdateIpBlackList 更新连接黑名单项
//
// input:
//   - oldAddr 原黑名单地址
//   - oldMask 原黑名单掩码
//   - newAddr 新黑名单地址
//   - newMask 新黑名单掩码
//   - newDesc 新黑名单描述
func (c *RtdbConnect) UpdateIpBlackList(oldAddr string, oldMask string, newAddr string, newMask string, newDesc string) error {
	rte := RawRtdbUpdateBlacklistWarp(c.ConnectHandle, oldAddr, oldMask, newAddr, newMask, newDesc)
	return rte.GoError()
}

// DeleteIpBlackList 删除连接黑名单项
//
// input:
//   - addr 黑名单地址
//   - mask 黑名单掩码
func (c *RtdbConnect) DeleteIpBlackList(addr string, mask string) error {
	rte := RawRtdbRemoveBlacklistWarp(c.ConnectHandle, addr, mask)
	return rte.GoError()
}

// GetIpBlackLists 获得连接黑名单列表
//
// output:
//   - []BlackList(lists) 连接黑名单列表
func (c *RtdbConnect) GetIpBlackLists() ([]BlackList, error) {
	lists, rte := RawRtdbGetBlacklistWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return lists, nil
}

// AddIpWhiteList 添加连接白名单
//
// input:
//   - addr 连接白名单地址
//   - mask 连接白名单掩码
//   - desc 连接白名单描述
//   - priv 连接白名单权限
func (c *RtdbConnect) AddIpWhiteList(addr string, mask string, desc string, priv PrivGroup) error {
	rte := RawRtdbAddAuthorizationWarp(c.ConnectHandle, addr, mask, desc, priv)
	return rte.GoError()
}

// UpdateIpWhiteList 更新连接白名单
//
// input:
//   - oldAddr 原连接白名单地址
//   - oldMask 原连接白名单掩码
//   - newAddr 新连接白名单地址
//   - newMask 新连接白名单掩码
//   - newDesc 新连接白名单描述
//   - newPriv 新连接白名单权限
func (c *RtdbConnect) UpdateIpWhiteList(oldAddr string, oldMask string, newAddr string, newMask string, newDesc string, newPriv PrivGroup) error {
	rte := RawRtdbUpdateAuthorizationWarp(c.ConnectHandle, oldAddr, oldMask, newAddr, newMask, newDesc, newPriv)
	return rte.GoError()
}

// DeleteIpWhiteList 删除白名单
//
// input:
//   - addr 连接白名单地址
//   - mask 连接白名单掩码
func (c *RtdbConnect) DeleteIpWhiteList(addr string, mask string) error {
	rte := RawRtdbRemoveAuthorizationWarp(c.ConnectHandle, addr, mask)
	return rte.GoError()
}

// GetIpWhiteLists 获取连接白名单列表
//
// output:
//   - []AuthorizationsList(lists)
func (c *RtdbConnect) GetIpWhiteLists() ([]AuthorizationsList, error) {
	lists, rte := RawRtdbGetAuthorizationsWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return lists, nil
}

// UpdatePassword 修改用户密码
//
// input:
//   - user 用户名
//   - password 用户密码
func (c *RtdbConnect) UpdatePassword(user string, password string) error {
	rte := RawRtdbChangePasswordWarp(c.ConnectHandle, user, password)
	return rte.GoError()
}

// UpdateOwnPassword 修改自己的密码
//
// input:
//   - oldPwd 旧密码
//   - newPwd 新密码
func (c *RtdbConnect) UpdateOwnPassword(oldPwd string, newPwd string) error {
	rte := RawRtdbChangeMyPasswordWarp(c.ConnectHandle, oldPwd, newPwd)
	return rte.GoError()
}

// GetPriv 获取连接权限
//
// output:
//   - PrivGroup(priv) 用户权限
func (c *RtdbConnect) GetPriv() (*PrivGroup, error) {
	priv, rte := RawRtdbGetPrivWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &priv, nil
}

// SetPriv 设置连接权限
//
// input:
//   - user 用户名
//   - priv 用户权限
func (c *RtdbConnect) SetPriv(user string, priv PrivGroup) error {
	rte := RawRtdbChangePrivWarp(c.ConnectHandle, user, priv)
	if RteIsOk(rte) && c.UserName == user {
		c.Priv = priv
	}
	return rte.GoError()
}

// AddUser 添加用户
//
// input:
//   - user 用户名
//   - password 用户密码
//   - priv 用户权限
func (c *RtdbConnect) AddUser(user string, password string, priv PrivGroup) error {
	rte := RawRtdbAddUserWarp(c.ConnectHandle, user, password, priv)
	return rte.GoError()
}

// DeleteUser 删除用户
//
// input:
//   - user 用户名
func (c *RtdbConnect) DeleteUser(user string) error {
	rte := RawRtdbRemoveUserWarp(c.ConnectHandle, user)
	return rte.GoError()
}

// LockUser 锁定用户
//
// input:
//   - user 用户名
//   - lock 是否锁定
func (c *RtdbConnect) LockUser(user string, lock Switch) error {
	rte := RawRtdbLockUserWarp(c.ConnectHandle, user, lock)
	return rte.GoError()
}

// GetUsers 获取用户列表
//
// output:
//   - []RtdbUserInfo(users) 用户列表
func (c *RtdbConnect) GetUsers() ([]RtdbUserInfo, error) {
	users, rte := RawRtdbGetUsersWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return users, nil
}

// AddNamedType 创建自定义类型
//
// input:
//   - name 自定义类型名称
//   - fields 自定义类型字段列表
//   - desc 自定义类型描述
func (c *RtdbConnect) AddNamedType(name string, desc string, fields ...RtdbDataTypeField) error {
	rte := RawRtdbbCreateNamedTypeWarp(c.ConnectHandle, name, desc, fields...)
	return rte.GoError()
}

// DeleteNamedType 删除自定义类型
//
// input:
//   - name 自定义类型的名称
func (c *RtdbConnect) DeleteNamedType(name string) error {
	rte := RawRtdbbRemoveNamedTypeWarp(c.ConnectHandle, name)
	return rte.GoError()
}

// GetNamedType 获取自定义类型
//
// output:
//   - NamedType 自定义类型
func (c *RtdbConnect) GetNamedType(name string) (*NamedType, error) {
	types, err := c.GetNamedTypes()
	if err != nil {
		return nil, err
	}

	for _, typ := range types {
		if typ.Name == name {
			return &typ, nil
		}
	}

	return nil, errors.New("未知自定义类型")
}

// GetNamedTypes 获取自定义类型列表
//
// output:
//   - []NamedType(types) 自定义类型列表
func (c *RtdbConnect) GetNamedTypes() ([]NamedType, error) {
	count, rte := RawRtdbbGetNamedTypesCountWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	names, fieldCounts, rte := RawRtdbbGetAllNamedTypesWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}

	types := make([]NamedType, count)
	for i := 0; i < len(names); i++ {
		fields, length, desc, rte := RawRtdbbGetNamedTypeWarp(c.ConnectHandle, names[i], fieldCounts[i])
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		types = append(types, NamedType{
			Name:   names[i],
			Fields: fields,
			Length: length,
			Desc:   desc,
		})
	}
	return types, nil
}

// UpdateNamedType 修改自定义类型
//
// input:
//   - name 自定义类型的名称
//   - modifyName 要修改的 自定义类型名称
//   - modifyDesc 要修改的 自定义类型的描述
//   - modifyFields 要修改的 字段名称<->字段描述
func (c *RtdbConnect) UpdateNamedType(name string, modifyName *string, modifyDesc *string, modifyFields map[string]string) error {
	fieldNames := make([]string, 0)
	fieldDescs := make([]string, 0)
	for name, desc := range modifyFields {
		fieldNames = append(fieldNames, name)
		fieldDescs = append(fieldDescs, desc)
	}
	rte := RawRtdbbModifyNamedTypeWarp(c.ConnectHandle, name, modifyName, modifyDesc, fieldNames, fieldDescs)
	return rte.GoError()
}

// ServerHostTime 服务端主机时间
func (c *RtdbConnect) ServerHostTime() (*time.Time, error) {
	datetime, rte := RawRtdbHostTime64Warp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	hostTime := time.Unix(int64(datetime), 0)
	return &hostTime, nil
}

// DurationToString 时间段转字符串, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
//
//	?y    ?年, 1年 = 365日
//	?m    ?月, 1月 = 30 日
//	?d    ?日
//	?h    ?小时
//	?n    ?分钟
//	?s    ?秒
//
// input:
//   - duration 时间段
//
// output:
//   - string(字符串格式时间段)
func (c *RtdbConnect) DurationToString(duration time.Duration) (string, error) {
	durationStr, rte := RawRtdbFormatTimespanWarp(int32(duration.Seconds()))
	if !RteIsOk(rte) {
		return "", rte.GoError()
	}
	return durationStr, nil
}

// StringToDuration 字符串转时间段, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
//
//	?y    ?年, 1年 = 365日
//	?m    ?月, 1月 = 30 日
//	?d    ?日
//	?h    ?小时
//	?n    ?分钟
//	?s    ?秒
//
// input:
//   - strDuration 字符串类型时间段
//
// output:
//   - time.Duration(duration) 时间段
func (c *RtdbConnect) StringToDuration(strDuration string) (time.Duration, error) {
	duration, rte := RawRtdbParseTimespanWarp(strDuration)
	if !RteIsOk(rte) {
		return 0, rte.GoError()
	}
	return time.Second * time.Duration(duration), nil
}

// StringToTime 字符串转时间戳
//
//	其中 base_time 表示基本时间，有三种形式：
//	1. 时间字符串，如 "2010-1-1" 及 "2010-1-1 8:00:00"；
//	2. 时间代码，表示客户端相对时间；
//	可用的时间代码及含义如下：
//	td             当天零点
//	yd             昨天零点
//	tm             明天零点
//	mon            本周一零点
//	tue            本周二零点
//	wed            本周三零点
//	thu            本周四零点
//	fri            本周五零点
//	sat            本周六零点
//	sun            本周日零点
//	3. 星号('*')，表示客户端当前时间。
//	offset_time 是可选的，可以出现多次，
//	可用的时间偏移代码及含义如下：
//	[+|-] ?y            偏移?年, 1年 = 365日
//	[+|-] ?m            偏移?月, 1月 = 30 日
//	[+|-] ?d            偏移?日
//	[+|-] ?h            偏移?小时
//	[+|-] ?n            偏移?分钟
//	[+|-] ?s            偏移?秒
//	[+|-] ?ms           偏移?毫秒
//	例如："*-1d" 表示当前时刻减去24小时。
//
// input:
//   - strTime 字符串类型时间戳
//
// output:
//   - time.Time(timestamp) 时间戳
func (c *RtdbConnect) StringToTime(strTime string) (*time.Time, error) {
	datetime, subtime, rte := RawRtdbParseTimeWarp(strTime)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	goTime := time.Unix(int64(datetime), int64(subtime))
	return &goTime, nil
}

// GetQualityDesc 获取质量码说明
func (c *RtdbConnect) GetQualityDesc(qualities []Quality) ([]string, error) {
	descs, rte := RawRtdbFormatQualityWarp(c.ConnectHandle, qualities)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return descs, nil
}

// GetDriveLetterList 获取盘符列表, windows平台是C、D、E、F这些盘符，linux平台是 / 盘符
//
// output:
//   - []string(litters) 盘符列表
func (c *RtdbConnect) GetDriveLetterList() ([]string, error) {
	letters, rte := RawRtdbGetLogicalDriversWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return letters, nil
}

// GetDirItemList 获取目录项列表
//
// input:
//   - dir 目录路径
//
// output:
//   - []DirItem(items) 目录项列表
func (c *RtdbConnect) GetDirItemList(dir string) ([]DirItem, error) {
	rte := RawRtdbOpenPathWarp(c.ConnectHandle, dir)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	defer func() {
		_ = RawRtdbClosePathWarp(c.ConnectHandle)
	}()

	items := make([]DirItem, 0)
	for {
		item, rte := RawRtdbReadPath64Warp(c.ConnectHandle)
		if !RteIsOk(rte) {
			if errors.Is(rte, RteBatchEnd) {
				break
			} else {
				return nil, rte.GoError()
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// CreateDir 创建目录
//
// input:
//   - path 目录路径
func (c *RtdbConnect) CreateDir(path string) error {
	rte := RawRtdbMkdirWarp(c.ConnectHandle, path)
	return rte.GoError()
}

// ReadFile 读取文件
//
// input:
//   - path 文件路径
//
// output:
//   - []byte(data) 文件内容
func (c *RtdbConnect) ReadFile(path string) ([]byte, error) {
	size, rte := RawRtdbGetFileSizeWarp(c.ConnectHandle, path)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	if size > MaxFileSize {
		return nil, errors.New("当前文件大小超出允许读取长度")
	}

	buf := bytes.NewBuffer(nil)
	for i := 0; i < int(size); i += MaxBlockSize {
		data, rte := RawRtdbReadFileWarp(c.ConnectHandle, path, int64(i*MaxBlockSize), MaxBlockSize)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		_, err := buf.Write(data)
		if err != nil {
			return nil, fmt.Errorf("写入缓存失败：%v", err)
		}
	}
	return buf.Bytes(), nil
}

// CreateTable 创建表
//
// input:
//   - name 表名
//   - desc 表描述
//
// output:
//   - RtdbTable(table) 返回表
func (c *RtdbConnect) CreateTable(name string, desc string) (*RtdbTable, error) {
	table, rte := RawRtdbbAppendTableWarp(c.ConnectHandle, name, desc)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &table, nil
}

// DeleteTable 删除表
//
// input:
//   - id 表ID
func (c *RtdbConnect) DeleteTable(id TableID) error {
	rte := RawRtdbbRemoveTableByIdWarp(c.ConnectHandle, id)
	return rte.GoError()
}

// GetTable
//
// input:
//   - id 获取表
func (c *RtdbConnect) GetTable(id TableID) (*RtdbTable, error) {
	table, rte := RawRtdbbGetTablePropertyByIdWarp(c.ConnectHandle, id)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return &table, nil
}

// GetTables 获取表列表
//
// output:
//   - []RtdbTable(tables) 表列表
func (c *RtdbConnect) GetTables() ([]RtdbTable, error) {
	count, rte := RawRtdbbTablesCountWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	ids, rte := RawRtdbbGetTablesWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	tables := make([]RtdbTable, 0)
	for _, id := range ids {
		table, rte := RawRtdbbGetTablePropertyByIdWarp(c.ConnectHandle, id)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		tables = append(tables, table)
	}
	return tables, nil
}

// UpdateTableName 更新表名
// input:
//   - id 表ID
//   - name 表名
func (c *RtdbConnect) UpdateTableName(id TableID, name string) error {
	rte := RawRtdbbUpdateTableNameWarp(c.ConnectHandle, id, name)
	return rte.GoError()
}

// UpdateTableDesc 更新表描述
//
// input:
//   - id 表ID
//   - desc 表描述
func (c *RtdbConnect) UpdateTableDesc(id TableID, desc string) error {
	rte := RawRtdbbUpdateTableDescByIdWarp(c.ConnectHandle, id, desc)
	return rte.GoError()
}

// AddPoint 创建点
//
// input:
//   - info 输入点信息
//
// output:
//   - PointInfo(info) 输出点信息
func (c *RtdbConnect) AddPoint(info *PointInfo) (*PointInfo, error) {
	base, scan, calc, tName := PointInfoToRaw(info)
	if base.Type == RtdbTypeNamedT {
		if tName == "" {
			return nil, errors.New("点数值类型为RtdbTypeNamedT, 此时NamedTypeName不能为空")
		}
		_, err := c.GetNamedType(tName)
		if err != nil {
			return nil, err
		}
		base, scan, rte := RawRtdbbInsertNamedTypePointWarp(c.ConnectHandle, base, scan, tName)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return PointInfoFromRaw(c.ConnectHandle, base, scan, nil, false)
	} else {
		base, scan, calc, rte := RawRtdbbInsertMaxPointWarp(c.ConnectHandle, base, scan, calc)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		return PointInfoFromRaw(c.ConnectHandle, base, scan, calc, false)
	}
}

// DeletePoint 删除点
//
// input:
//   - id 点ID
func (c *RtdbConnect) DeletePoint(id PointID) error {
	rte := RawRtdbbRemovePointByIdWarp(c.ConnectHandle, id)
	return rte.GoError()
}

// UpdatePoint 更新点
//
// input:
//   - id 点ID
//   - fields 需要更新的字段
func (c *RtdbConnect) UpdatePoint(id PointID, fields map[PointInfoField]any) error {
	pointInfo, err := c.GetPoint(id)
	if err != nil {
		return err
	}
	for k, v := range fields {
		switch k {
		case PointInfoFieldName:
			name, ok := v.(string)
			if !ok {
				return errors.New("Name应该为String类型")
			}
			pointInfo.Name = name
		case PointInfoFieldClass:
			class, ok := v.(PointClass)
			if !ok {
				return errors.New("Class应为PointClass类型")
			}
			pointInfo.Class = class
		case PointInfoFieldDesc:
			desc, ok := v.(string)
			if !ok {
				return errors.New("Desc应该为String类型")
			}
			pointInfo.Desc = desc
		case PointInfoFieldUnit:
			unit, ok := v.(string)
			if !ok {
				return errors.New("Unit应该为string类型")
			}
			pointInfo.Unit = unit
		case PointInfoFieldArchive:
			archive, ok := v.(Switch)
			if !ok {
				return errors.New("Archive应该为Switch类型")
			}
			pointInfo.Archive = archive
		case PointInfoFieldDigits:
			digits, ok := v.(int16)
			if !ok {
				return errors.New("Digits应该为Int16类型")
			}
			pointInfo.Digits = digits
		case PointInfoFieldShutdown:
			shutdown, ok := v.(Switch)
			if !ok {
				return errors.New("Shutdown应该为Switch类型")
			}
			pointInfo.Shutdown = shutdown
		case PointInfoFieldLowLimit:
			limit, ok := v.(float32)
			if !ok {
				return errors.New("LowLimit应该为float32类型")
			}
			pointInfo.LowLimit = limit
		case PointInfoFieldHighLimit:
			limit, ok := v.(float32)
			if !ok {
				return errors.New("HighLimit应该为float32类型")
			}
			pointInfo.HighLimit = limit
		case PointInfoFieldStep:
			step, ok := v.(Switch)
			if !ok {
				return errors.New("Step应该为Switch类型")
			}
			pointInfo.Step = step
		case PointInfoFieldTypical:
			typical, ok := v.(float32)
			if !ok {
				return errors.New("典型值应该为float32类型")
			}
			pointInfo.Typical = typical
		case PointInfoFieldCompress:
			compress, ok := v.(Switch)
			if !ok {
				return errors.New("Compress应该为Switch类型")
			}
			pointInfo.Compress = compress
		case PointInfoFieldCompDev:
			compDev, ok := v.(float32)
			if !ok {
				return errors.New("CompDev应该为float32类型")
			}
			pointInfo.CompDev = compDev
		case PointInfoFieldCompDevPercent:
			compDevPercent, ok := v.(float32)
			if !ok {
				return errors.New("CompDevPercent应该为float32类型")
			}
			pointInfo.CompDevPercent = compDevPercent
		case PointInfoFieldCompTimeMax:
			compTimeMax, ok := v.(int32)
			if !ok {
				return errors.New("CompTimeMax应为int32类型")
			}
			pointInfo.CompTimeMax = compTimeMax
		case PointInfoFieldCompTimeMin:
			compTimeMin, ok := v.(int32)
			if !ok {
				return errors.New("CompTimeMin应为int32类型")
			}
			pointInfo.CompTimeMin = compTimeMin
		case PointInfoFieldExcDev:
			excDev, ok := v.(float32)
			if !ok {
				return errors.New("ExcDev应为float32类型")
			}
			pointInfo.ExcDev = excDev
		case PointInfoFieldExcDevPercent:
			excDevPercent, ok := v.(float32)
			if !ok {
				return errors.New("ExcDevPercent应为float32类型")
			}
			pointInfo.ExcDevPercent = excDevPercent
		case PointInfoFieldExcTimeMax:
			excTimeMax, ok := v.(int32)
			if !ok {
				return errors.New("ExcTimeMax应为int32类型")
			}
			pointInfo.ExcTimeMax = excTimeMax
		case PointInfoFieldExcTimeMin:
			excTimeMin, ok := v.(int32)
			if !ok {
				return errors.New("ExcTimeMin应为int32类型")
			}
			pointInfo.ExcTimeMin = excTimeMin
		case PointInfoFieldMirror:
			mirror, ok := v.(RtdbMirror)
			if !ok {
				return errors.New("Mirror应该为RtdbMirror类型")
			}
			pointInfo.Mirror = mirror
		case PointInfoFieldSummary:
			summary, ok := v.(Switch)
			if !ok {
				return errors.New("Summary应该为Switch类型")
			}
			pointInfo.Summary = summary
		case PointInfoFieldSource:
			source, ok := v.(string)
			if !ok {
				return errors.New("Source应该为String类型")
			}
			pointInfo.Source = source
		case PointInfoFieldScan:
			scan, ok := v.(Switch)
			if !ok {
				return errors.New("scan应该为Switch类型")
			}
			pointInfo.Scan = scan
		case PointInfoFieldInstrument:
			instrument, ok := v.(string)
			if !ok {
				return errors.New("instrument应该为String类型")
			}
			pointInfo.Instrument = instrument
		case PointInfoFieldLocations:
			locations, ok := v.([RtdbConstLocationsSize]int32)
			if !ok {
				return errors.New("locations应该为[5]int32类型")
			}
			pointInfo.Locations = locations
		case PointInfoFieldUserInts:
			userInts, ok := v.([RtdbConstUserintSize]int32)
			if !ok {
				return errors.New("userInfos应该为[2]int32类型")
			}
			pointInfo.UserInts = userInts
		case PointInfoFieldUserReals:
			userReals, ok := v.([RtdbConstUserrealSize]float32)
			if !ok {
				return errors.New("userReals应该为[2]float32类型")
			}
			pointInfo.UserReals = userReals
		case PointInfoFieldEquation:
			equation, ok := v.(string)
			if !ok {
				return errors.New("equation应该为String类型")
			}
			pointInfo.Equation = equation
		case PointInfoFieldTrigger:
			trigger, ok := v.(RtdbTrigger)
			if !ok {
				return errors.New("trigger应该为RtdbTrigger类型")
			}
			pointInfo.Trigger = trigger
		case PointInfoFieldTimeCopy:
			timeCopy, ok := v.(RtdbTimeCopy)
			if !ok {
				return errors.New("timeCopy应该为RtdbTimeCopy类型")
			}
			pointInfo.TimeCopy = timeCopy
		case PointInfoFieldPeriod:
			period, ok := v.(int32)
			if !ok {
				return errors.New("period应该为int32类型")
			}
			pointInfo.Period = period
		default:
			return errors.New("未知的Field")
		}
	}
	base, scan, calc, _ := PointInfoToRaw(pointInfo)
	rte := RawRtdbbUpdateMaxPointPropertyWarp(c.ConnectHandle, base, scan, calc)
	return rte.GoError()
}

// GetPoints 批量获取标签点
//
// input:
//   - ids 标签点列表
//
// output:
//   - []PointInfo(infos) 标签点属性列表
func (c *RtdbConnect) GetPoints(ids []PointID) ([]*PointInfo, []error, error) {
	bases, scans, calcs, rtes, rte := RawRtdbbGetMaxPointsPropertyWarp(c.ConnectHandle, ids)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}
	errs := RtdbErrorListToErrorList(rtes)
	infos := make([]*PointInfo, 0)
	for i := 0; i < len(ids); i++ {
		info, err := PointInfoFromRaw(c.ConnectHandle, &bases[i], &scans[i], &calcs[i], false)
		if err != nil {
			errs[i] = err
		}
		infos = append(infos, info)
	}
	return infos, errs, nil
}

// GetPoint 获取点
//
// input:
//   - id 点ID
//
// output:
//   - PointInfo(info) 返回点信息
func (c *RtdbConnect) GetPoint(id PointID) (*PointInfo, error) {
	bases, scans, calcs, rtes, rte := RawRtdbbGetMaxPointsPropertyWarp(c.ConnectHandle, []PointID{id})
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	for _, rte := range rtes {
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
	}
	return PointInfoFromRaw(c.ConnectHandle, &bases[0], &scans[0], &calcs[0], false)
}

// FindPoints 根据 表名.点名 搜索标签点
//
// input:
//   - tableDotPoints 点全名， 表名.点名
//
// output:
//   - []*PointInfo(infos) 点信息列表
//   - []error 报错信息
func (c *RtdbConnect) FindPoints(tableDotPoints []string) ([]*PointInfo, []error, error) {
	ids, _, _, _, _, rte := RawRtdbbFindPointsExWarp(c.ConnectHandle, tableDotPoints)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}
	return c.GetPoints(ids)
}

// MovePoint 移动点到指定表
//
// input:
//   - id 点ID
//   - tableName 表名称
func (c *RtdbConnect) MovePoint(id PointID, tableName string) error {
	rte := RawRtdbbMovePointByIdWarp(c.ConnectHandle, id, tableName)
	return rte.GoError()
}

// SearchPoint 分页搜索点
//
// input:
//   - start 开始ID
//   - count 最多返回PointInfo个数
//   - tagMask 标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - tableMask 标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - source 数据源集合，字符串中的每个字符均表示一个数据源，空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//   - unit 标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//   - desc 标签点描述的子集，描述中包含该参数的标签点均满足条件，空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//   - instrument 标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//   - typeMask 标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
//   - classOfMask 标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
//   - timeUnitMask 标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
//   - otherTypeMask 使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
//   - otherTypeMaskValue 字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略, 当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
//   - mode 搜索结果排序模式
//   - 备注：多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
//
// output:
//   - int32(count) 点总数
//   - []*PointInfo(infos) 点信息列表
func (c *RtdbConnect) SearchPoint(start int32, count int32, tagMask, tableMask, source, unit, desc, instrument, typeMask string, classOfMask RtdbType, timeUnitMask RtdbPrecision, otherTypeMask RtdbSearch, otherTypeMaskValue string, model RtdbSortFlag) (int32, []*PointInfo, []error, error) {
	count, rte := RawRtdbbSearchPointsCountWarp(c.ConnectHandle, tagMask, tableMask, source, unit, desc, instrument, typeMask, classOfMask, timeUnitMask, otherTypeMask, otherTypeMaskValue)
	if !RteIsOk(rte) {
		return 0, nil, nil, rte.GoError()
	}
	ids, rte := RawRtdbbSearchExWarp(c.ConnectHandle, count, tagMask, tableMask, source, unit, desc, instrument, typeMask, classOfMask, timeUnitMask, otherTypeMask, otherTypeMaskValue, model)
	if !RteIsOk(rte) {
		return 0, nil, nil, rte.GoError()
	}
	ids = SafeSlice(ids, start, count)
	infos, errs, err := c.GetPoints(ids)
	if err != nil {
		return 0, nil, nil, err
	}
	return count, infos, errs, nil
}

// ClearRecycler 清空回收站
func (c *RtdbConnect) ClearRecycler() error {
	rte := RawRtdbbClearRecyclerWarp(c.ConnectHandle)
	return rte.GoError()
}

// GetRecycledPoints 分段获取回收站中的点
//
// input:
//   - start 开始ID
//   - count 一次最多返回点个数
//
// output:
//   - int32(count) 回收站中的总点数
//   - []*PointInfo(infos) 点信息列表
//   - []error(errs) 获取点信息时的错误列表
func (c *RtdbConnect) GetRecycledPoints(start int32, count int32) (int32, []*PointInfo, []error, error) {
	count, rte := RawRtdbbGetRecycledPointsCountWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return 0, nil, nil, rte.GoError()
	}
	ids, rte := RawRtdbbGetRecycledPointsWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return 0, nil, nil, rte.GoError()
	}
	ids = SafeSlice(ids, start, count)
	infos := make([]*PointInfo, 0)
	errs := make([]error, 0)
	for _, id := range ids {
		base, scan, calc, rte := RawRtdbbGetRecycledMaxPointPropertyWarp(c.ConnectHandle, id)
		info, _ := PointInfoFromRaw(c.ConnectHandle, base, scan, calc, true)
		infos = append(infos, info)
		if !RteIsOk(rte) {
			errs = append(errs, rte.GoError())
		} else {
			errs = append(errs, nil)
		}
	}
	return count, infos, errs, nil
}

// RecoverPoint 从回收站中恢复点到某个表
//
// input:
//   - tableID 点恢复到这个表
//   - pointID 需要恢复的点
func (c *RtdbConnect) RecoverPoint(tableId TableID, pointId PointID) error {
	rte := RawRtdbbRecoverPointWarp(c.ConnectHandle, tableId, pointId)
	return rte.GoError()
}

// PurgePoint 从回收站中清除点
//
// input:
//   - id 点ID
func (c *RtdbConnect) PurgePoint(id PointID) error {
	rte := RawRtdbbPurgePointWarp(c.ConnectHandle, id)
	return rte.GoError()
}

// SearchRecycledPoint 从回收站中搜索点
//
// input:
//   - start 分页开始位置
//   - count 分页获取个数
//   - tagMask 标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - tableMask 标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
//   - source 数据源集合，字符串中的每个字符均表示一个数据源，空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
//   - unit 标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
//   - desc 标签点描述的子集，描述中包含该参数的标签点均满足条件，空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
//   - instrument 标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
//   - mode 搜索结果排序模式
//
// output:
//   - int32(count) 回收站中的总点数
//   - []*PointInfo(infos) 点信息列表
//   - []error(errs) 获取点信息时的错误列表
func (c *RtdbConnect) SearchRecycledPoint(start int32, count int32, tagMask, tableMask, source, unit, desc, instrument string, mode RtdbSortFlag) (int32, []*PointInfo, []error, error) {
	maxCount, rte := RawRtdbbGetRecycledPointsCountWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return 0, nil, nil, rte.GoError()
	}
	ids, rte := RawRtdbbSearchRecycledPointsInBatchesWarp(c.ConnectHandle, start, maxCount, tagMask, tableMask, source, unit, desc, instrument, mode)
	if !RteIsOk(rte) {
		return 0, nil, nil, rte.GoError()
	}
	maxCount = int32(len(ids))
	rtnIds := SafeSlice(ids, start, count)

	infos := make([]*PointInfo, 0)
	errs := make([]error, 0)
	for _, id := range rtnIds {
		base, scan, calc, rte := RawRtdbbGetRecycledMaxPointPropertyWarp(c.ConnectHandle, id)
		info, _ := PointInfoFromRaw(c.ConnectHandle, base, scan, calc, true)
		infos = append(infos, info)
		if !RteIsOk(rte) {
			errs = append(errs, rte.GoError())
		} else {
			errs = append(errs, nil)
		}
	}

	return maxCount, infos, errs, nil
}

// GetPointCountFromValueType 获取某个数据类型的点个数 (可以是内置类型，也可以是自定义类型)
//
// input:
//   - valueType 数值类型
//
// output:
//   - int32(count) 该数值类型对应的点数量
func (c *RtdbConnect) GetPointCountFromValueType(valueType ValueType) (int32, error) {
	rtdbType, name := valueType.ToRawType()
	if rtdbType == RtdbTypeNamedT {
		count, rte := RawRtdbbGetNamedTypePointsCountWarp(c.ConnectHandle, name)
		if !RteIsOk(rte) {
			return 0, rte.GoError()
		}
		return count, nil
	} else {
		count, rte := RawRtdbbGetBaseTypePointsCountWarp(c.ConnectHandle, rtdbType)
		if !RteIsOk(rte) {
			return 0, rte.GoError()
		}
		return count, nil
	}
}

// WriteValue 写入值
//
// input:
//   - info 标签点基本信息
//   - fix 是否覆盖写入，备注：只有Int、Float和Coor支持
//   - tvq 时间戳+数值+质量码
func (c *RtdbConnect) WriteValue(info *PointInfo, fix bool, tvq TVQ) error {
	errs, err := c.WriteValues(info, fix, []TVQ{tvq})
	if errs[0] != nil {
		return errs[0]
	}
	return err
}

// WriteValues 批量写入值
//
// input:
//   - info 标签点信息
//   - fix 是否覆盖写入
//   - tvqs TVQ列表
//
// output:
//   - []error(errs) 错误列表
func (c *RtdbConnect) WriteValues(info *PointInfo, fix bool, tvqs []TVQ) ([]error, error) {
	ptvqs := make([]PTVQ, 0)
	for _, tvq := range tvqs {
		ptvqs = append(ptvqs, NewPTVQ(info, tvq))
	}
	return c.WriteSection(fix, ptvqs)
}

// WriteSection 写断面(批量写入多个Point，每个Point写入一个TVQ)
//
// input:
//   - fix 是否可以覆盖写入，只对 整数、浮点数、坐标 生效
//   - ptvqs PTVQ值数组, 备注：p是可以重复的，表示一个point中写入多条数值
func (c *RtdbConnect) WriteSection(fix bool, ptvqs []PTVQ) ([]error, error) {
	rtnRtes := make([]RtdbError, len(ptvqs))
	sort.Slice(ptvqs, func(i, j int) bool {
		if ptvqs[i].TVQ.Timestamp.Compare(ptvqs[j].TVQ.Timestamp) > 0 {
			return true
		} else {
			return false
		}
	})

	// 数值 int&float
	numberIds := make([]PointID, 0)
	numberDatetimes := make([]TimestampType, 0)
	numberSubtimes := make([]SubtimeType, 0)
	numberValues := make([]float64, 0)
	numberStates := make([]int64, 0)
	numberQualities := make([]Quality, 0)
	numberIdx := make([]int, 0)

	// 坐标
	coorIds := make([]PointID, 0)
	coorDatetimes := make([]TimestampType, 0)
	coorSubtimes := make([]SubtimeType, 0)
	coorXs := make([]float32, 0)
	coorYs := make([]float32, 0)
	coorQualities := make([]Quality, 0)
	coorIdx := make([]int, 0)

	// String｜Blob
	bIds := make([]PointID, 0)
	bDatetimes := make([]TimestampType, 0)
	bSubtimes := make([]SubtimeType, 0)
	bDatas := make([][]byte, 0)
	bQualities := make([]Quality, 0)
	bIdx := make([]int, 0)

	// named 自定义类型
	namedIds := make([]PointID, 0)
	namedDatetimes := make([]TimestampType, 0)
	namedSubtimes := make([]SubtimeType, 0)
	namedDatas := make([][]byte, 0)
	namedQualities := make([]Quality, 0)
	namedIdx := make([]int, 0)

	// datetime 日期
	dtIds := make([]PointID, 0)
	dtDatetimes := make([]TimestampType, 0)
	dtSubtimes := make([]SubtimeType, 0)
	dtDates := make([]string, 0)
	dtQualities := make([]Quality, 0)
	dtIdx := make([]int, 0)

	for i, ptvq := range ptvqs {
		rtdbType, _ := ptvq.PointInfo.ValueType.ToRawType()
		switch rtdbType {
		case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
			numberIds = append(numberIds, ptvq.PointInfo.ID)
			datetime, subtime := ptvq.TVQ.GetRtdbTimestamp()
			numberDatetimes = append(numberDatetimes, datetime)
			numberSubtimes = append(numberSubtimes, subtime)
			numberQualities = append(numberQualities, ptvq.TVQ.GetRtdbQuality())
			switch rtdbType {
			case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64:
				numberValues = append(numberValues, 0)
				numberStates = append(numberStates, ptvq.TVQ.GetRtdbInt())
			case RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
				numberValues = append(numberValues, ptvq.TVQ.GetRtdbFloat())
				numberStates = append(numberStates, 0)
			}
			numberIdx = append(numberIdx, i)
		case RtdbTypeCoor:
			coorIds = append(coorIds, ptvq.PointInfo.ID)
			datetime, subtime := ptvq.TVQ.GetRtdbTimestamp()
			coorDatetimes = append(coorDatetimes, datetime)
			coorSubtimes = append(coorSubtimes, subtime)
			coorQualities = append(coorQualities, ptvq.TVQ.GetRtdbQuality())
			xy := ptvq.TVQ.GetRtdbCoordinates()
			coorXs = append(coorXs, xy.X)
			coorYs = append(coorYs, xy.Y)
			coorIdx = append(coorIdx, i)
		case RtdbTypeString, RtdbTypeBlob:
			bIds = append(bIds, ptvq.PointInfo.ID)
			datetime, subtime := ptvq.TVQ.GetRtdbTimestamp()
			bDatetimes = append(bDatetimes, datetime)
			bSubtimes = append(bSubtimes, subtime)
			bQualities = append(bQualities, ptvq.TVQ.GetRtdbQuality())
			data, err := ptvq.TVQ.GetRtdbStringBlob(c.ServerOsType)
			if err != nil {
				return nil, err
			}
			bDatas = append(bDatas, data)
			bIdx = append(bIdx, i)
		case RtdbTypeNamedT:
			namedIds = append(namedIds, ptvq.PointInfo.ID)
			datetime, subtime := ptvq.TVQ.GetRtdbTimestamp()
			namedDatetimes = append(namedDatetimes, datetime)
			namedSubtimes = append(namedSubtimes, subtime)
			namedQualities = append(namedQualities, ptvq.TVQ.GetRtdbQuality())
			namedDatas = append(namedDatas, ptvq.TVQ.GetRtdbNamedObj())
			namedIdx = append(namedIdx, i)
		case RtdbTypeDatetime:
			dtIds = append(dtIds, ptvq.PointInfo.ID)
			datetime, subtime := ptvq.TVQ.GetRtdbTimestamp()
			dtDatetimes = append(dtDatetimes, datetime)
			dtSubtimes = append(dtSubtimes, subtime)
			dtQualities = append(dtQualities, ptvq.TVQ.GetRtdbQuality())
			dtDates = append(dtDates, ptvq.TVQ.GetRtdbDatetime())
			dtIdx = append(dtIdx, i)
		}
	}

	if len(numberIds) != 0 {
		rtes := make([]RtdbError, 0)
		rte := RtdbError(0)
		if fix {
			rtes, rte = RawRtdbsFixSnapshots64Warp(c.ConnectHandle, numberIds, numberDatetimes, numberSubtimes, numberValues, numberStates, numberQualities)
		} else {
			rtes, rte = RawRtdbsPutSnapshots64Warp(c.ConnectHandle, numberIds, numberDatetimes, numberSubtimes, numberValues, numberStates, numberQualities)
		}
		fmt.Println("s1111:", rtes)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		aIndex := make([]int, 0)
		aIds := make([]PointID, 0)
		aDatetimes := make([]TimestampType, 0)
		aSubtimes := make([]SubtimeType, 0)
		aValues := make([]float64, 0)
		aStates := make([]int64, 0)
		aQualities := make([]Quality, 0)
		for i, e := range rtes {
			if errors.Is(e, RteTimestampEarlierThanSnapshot) {
				aIndex = append(aIndex, i)
				aIds = append(aIds, numberIds[i])
				aDatetimes = append(aDatetimes, numberDatetimes[i])
				aSubtimes = append(aSubtimes, numberSubtimes[i])
				aValues = append(aValues, numberValues[i])
				aStates = append(aStates, numberStates[i])
				aQualities = append(aQualities, numberQualities[i])
			}
		}
		if len(aIds) != 0 {
			aRtes, aRte := RawRtdbhPutArchivedValues64Warp(c.ConnectHandle, aIds, aDatetimes, aSubtimes, aValues, aStates, aQualities)
			if !RteIsOk(aRte) {
				return nil, aRte.GoError()
			}
			for i, e := range aRtes {
				rtes[aIndex[i]] = e
			}
			for i, e := range rtes {
				rtnRtes[numberIdx[i]] = e
				fmt.Println("1111111:", e)
			}
		}
	}

	if len(coorIds) != 0 {
		rtes := make([]RtdbError, 0)
		rte := RtdbError(0)
		if fix {
			rtes, rte = RawRtdbsFixCoorSnapshots64Warp(c.ConnectHandle, coorIds, coorDatetimes, coorSubtimes, coorXs, coorYs, coorQualities)
		} else {
			rtes, rte = RawRtdbsPutCoorSnapshots64Warp(c.ConnectHandle, coorIds, coorDatetimes, coorSubtimes, coorXs, coorYs, coorQualities)
		}
		fmt.Println(coorDatetimes, coorSubtimes)
		fmt.Println("s22222:", rtes)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		aIndex := make([]int, 0)
		aIds := make([]PointID, 0)
		aDatetimes := make([]TimestampType, 0)
		aSubtimes := make([]SubtimeType, 0)
		aXs := make([]float32, 0)
		aYs := make([]float32, 0)
		aQualities := make([]Quality, 0)
		for i, e := range rtes {
			if errors.Is(e, RteTimestampEarlierThanSnapshot) {
				aIndex = append(aIndex, i)
				aIds = append(aIds, coorIds[i])
				aDatetimes = append(aDatetimes, coorDatetimes[i])
				aSubtimes = append(aSubtimes, coorSubtimes[i])
				aXs = append(aXs, coorXs[i])
				aYs = append(aYs, coorYs[i])
				aQualities = append(aQualities, coorQualities[i])
			}
		}
		if len(aIds) != 0 {
			aRtes, aRte := RawRtdbhPutArchivedCoorValues64Warp(c.ConnectHandle, aIds, aDatetimes, aSubtimes, aXs, aYs, aQualities)
			if !RteIsOk(aRte) {
				return nil, aRte.GoError()
			}
			for i, e := range aRtes {
				rtes[aIndex[i]] = e
			}
			for i, e := range rtes {
				rtnRtes[coorIdx[i]] = e
				fmt.Println("22222:", e)
			}
		}
	}

	if len(bIds) != 0 {
		rtes, rte := RawRtdbsPutBlobSnapshots64Warp(c.ConnectHandle, bIds, bDatetimes, bSubtimes, bDatas, bQualities)
		fmt.Println("s3333333:", rtes)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		aIndex := make([]int, 0)
		aIds := make([]PointID, 0)
		aDatetimes := make([]TimestampType, 0)
		aSubtimes := make([]SubtimeType, 0)
		aDatas := make([][]byte, 0)
		aQualities := make([]Quality, 0)
		for i, e := range rtes {
			if errors.Is(e, RteTimestampEarlierThanSnapshot) {
				aIndex = append(aIndex, i)
				aIds = append(aIds, bIds[i])
				aDatetimes = append(aDatetimes, bDatetimes[i])
				aSubtimes = append(aSubtimes, bSubtimes[i])
				aDatas = append(aDatas, bDatas[i])
				aQualities = append(aQualities, bQualities[i])
			}
		}
		if len(aIds) != 0 {
			aRtes, aRte := RawRtdbhPutArchivedBlobValues64Warp(c.ConnectHandle, aIds, aDatetimes, aSubtimes, aDatas, aQualities)
			if !RteIsOk(aRte) {
				return nil, aRte.GoError()
			}
			for i, e := range aRtes {
				rtes[aIndex[i]] = e
			}
			for i, e := range rtes {
				rtnRtes[bIdx[i]] = e
				fmt.Println("33333333:", e)
			}
		}
	}

	if len(namedIds) != 0 {
		rtes, rte := RawRtdbsPutNamedTypeSnapshots64Warp(c.ConnectHandle, namedIds, namedDatetimes, namedSubtimes, namedDatas, namedQualities)
		fmt.Println("s4444444:", rtes)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		aIndex := make([]int, 0)
		aIds := make([]PointID, 0)
		aDatetimes := make([]TimestampType, 0)
		aSubtimes := make([]SubtimeType, 0)
		aDatas := make([][]byte, 0)
		aQualities := make([]Quality, 0)
		for i, e := range rtes {
			if errors.Is(e, RteTimestampEarlierThanSnapshot) {
				aIndex = append(aIndex, i)
				aIds = append(aIds, namedIds[i])
				aDatetimes = append(aDatetimes, namedDatetimes[i])
				aSubtimes = append(aSubtimes, namedSubtimes[i])
				aDatas = append(aDatas, namedDatas[i])
				aQualities = append(aQualities, namedQualities[i])
			}
		}
		if len(aIds) != 0 {
			aRtes, aRte := RawRtdbhPutArchivedNamedTypeValues64Warp(c.ConnectHandle, aIds, aDatetimes, aSubtimes, aDatas, aQualities)
			if !RteIsOk(aRte) {
				return nil, aRte.GoError()
			}
			for i, e := range aRtes {
				rtes[aIndex[i]] = e
			}
			for i, e := range rtes {
				rtnRtes[namedIdx[i]] = e
				fmt.Println("44444444:", e)
			}
		}
	}

	if len(dtIds) != 0 {
		rtes, rte := RawRtdbsPutDatetimeSnapshots64Warp(c.ConnectHandle, dtIds, dtDatetimes, dtSubtimes, dtDates, dtQualities)
		fmt.Println("s55555:", rtes)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		aIndex := make([]int, 0)
		aIds := make([]PointID, 0)
		aDatetimes := make([]TimestampType, 0)
		aSubtimes := make([]SubtimeType, 0)
		aDates := make([]string, 0)
		aQualities := make([]Quality, 0)
		for i, e := range rtes {
			if errors.Is(e, RteTimestampEarlierThanSnapshot) {
				aIndex = append(aIndex, i)
				aIds = append(aIds, dtIds[i])
				aDatetimes = append(aDatetimes, dtDatetimes[i])
				aSubtimes = append(aSubtimes, dtSubtimes[i])
				aDates = append(aDates, dtDates[i])
				aQualities = append(aQualities, dtQualities[i])
			}
		}
		if len(aIds) != 0 {
			aRtes, aRte := RawRtdbhPutArchivedDatetimeValues64Warp(c.ConnectHandle, aIds, aDatetimes, aSubtimes, aDates, aQualities)
			if !RteIsOk(aRte) {
				return nil, aRte.GoError()
			}
			for i, e := range aRtes {
				rtes[aIndex[i]] = e
			}
			for i, e := range rtes {
				rtnRtes[dtIdx[i]] = e
				fmt.Println("55555555:", e)
			}
		}
	}
	fmt.Println(rtnRtes)
	return RtdbErrorListToErrorList(rtnRtes), nil
}

// ReadLasts 批量读取实时快照值(当前标签点最后一个写入的TVQ)
//
// input:
//   - infos 标签点信息列表
//
// output:
//   - []PTVQ(ptvqs) 实时数值列表
//   - []error(errs) 错误列表
func (c *RtdbConnect) ReadLasts(infos []*PointInfo) ([]PTVQ, []error, error) {
	rtnRtes := make([]RtdbError, len(infos))
	rtnPTVQs := make([]PTVQ, len(infos))

	// 数值 int&float
	numberIds := make([]PointID, 0)
	numberIdx := make([]int, 0)

	// 坐标
	coorIds := make([]PointID, 0)
	coorIdx := make([]int, 0)

	// String｜Blob
	bIds := make([]PointID, 0)
	bIdx := make([]int, 0)

	// named 自定义类型
	namedIds := make([]PointID, 0)
	namedIdx := make([]int, 0)

	// datetime 日期
	dtIds := make([]PointID, 0)
	dtIdx := make([]int, 0)

	for i, info := range infos {
		rtdbType, _ := info.ValueType.ToRawType()
		switch rtdbType {
		case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
			numberIds = append(numberIds, info.ID)
			numberIdx = append(numberIdx, i)
		case RtdbTypeCoor:
			coorIds = append(coorIds, info.ID)
			coorIdx = append(coorIdx, i)
		case RtdbTypeBlob, RtdbTypeString:
			bIds = append(bIds, info.ID)
			bIdx = append(bIdx, i)
		case RtdbTypeNamedT:
			namedIds = append(namedIds, info.ID)
			namedIdx = append(namedIdx, i)
		case RtdbTypeDatetime:
			dtIds = append(dtIds, info.ID)
			dtIdx = append(dtIdx, i)
		}
	}

	if len(numberIds) != 0 {
		dt, ms, values, states, qualities, rtes, rte := RawRtdbsGetSnapshots64Warp(c.ConnectHandle, numberIds)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		for i := 0; i < len(numberIds); i++ {
			idx := numberIdx[i]
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			rtnRtes[idx] = rtes[i]
			rtdbType, _ := infos[idx].ValueType.ToRawType()
			switch rtdbType {
			case RtdbTypeBool:
				tvq := NewTvqBool(ts, Int64ToBool(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeUint8:
				tvq := NewTvqUint8(ts, uint8(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeInt8:
				tvq := NewTvqInt8(ts, int8(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeChar:
				tvq := NewTvqChar(ts, byte(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeUint16:
				tvq := NewTvqUint16(ts, uint16(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeInt16:
				tvq := NewTvqInt16(ts, int16(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeUint32:
				tvq := NewTvqUint32(ts, uint32(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeInt32:
				tvq := NewTvqInt32(ts, int32(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeInt64:
				tvq := NewTvqInt64(ts, int64(states[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeReal16:
				tvq := NewTvqFloat16(ts, float32(values[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeReal32:
				tvq := NewTvqFloat32(ts, float32(values[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeReal64:
				tvq := NewTvqFloat64(ts, values[i], q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeFp16:
				tvq := NewTvqFp16(ts, float32(values[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeFp32:
				tvq := NewTvqFp32(ts, float32(values[i]), q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeFp64:
				tvq := NewTvqFp64(ts, values[i], q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			}
		}
	}

	if len(coorIds) != 0 {
		dt, ms, xs, ys, qualities, rtes, rte := RawRtdbsGetCoorSnapshots64Warp(c.ConnectHandle, coorIds)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		for i := 0; i < len(coorIds); i++ {
			idx := coorIdx[i]
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			rtnRtes[idx] = rtes[i]
			tvq := NewTvqCoordinates(ts, xs[i], ys[i], q)
			rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
		}
	}

	if len(bIds) != 0 {
		dt, ms, blobs, qualities, rtes, rte := RawRtdbsGetBlobSnapshots64Warp(c.ConnectHandle, bIds, c.StringBlobMaxLen)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		for i := 0; i < len(bIds); i++ {
			idx := bIdx[i]
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			rtnRtes[idx] = rtes[i]
			rtdbType, _ := infos[idx].ValueType.ToRawType()
			switch rtdbType {
			case RtdbTypeBlob:
				tvq := NewTvqBlob(ts, blobs[i], q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			case RtdbTypeString:
				str := string(blobs[i])
				if c.ServerOsType == RtdbOsWindows {
					s, err := GBKBytesToString(blobs[i])
					if err != nil {
						return nil, nil, err
					}
					str = s
				}
				tvq := NewTvqString(ts, str, q)
				rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
			}
		}
	}

	if len(namedIds) != 0 {
		lens := make([]int32, 0)
		for i := 0; i < len(namedIds); i++ {
			idx := namedIdx[i]
			lens = append(lens, infos[idx].NamedType.Length)
		}

		dt, ms, datas, qualities, rtes, rte := RawRtdbsGetNamedTypeSnapshots64Warp(c.ConnectHandle, namedIds, lens)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		for i := 0; i < len(namedIds); i++ {
			idx := namedIdx[i]
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			rtnRtes[idx] = rtes[i]
			tvq := NewTvqNamed(ts, infos[idx].ValueType, datas[i], q)
			rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
		}
	}

	if len(dtIds) != 0 {
		dt, ms, datetimes, qualities, rtes, rte := RawRtdbsGetDatetimeSnapshots64Warp(c.ConnectHandle, dtIds, 0)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		for i := 0; i < len(dtIds); i++ {
			idx := dtIdx[i]
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			rtnRtes[idx] = rtes[i]
			tvq := NewTvqDatetime(ts, datetimes[i], q)
			rtnPTVQs[idx] = NewPTVQ(infos[idx], tvq)
		}
	}

	return rtnPTVQs, RtdbErrorListToErrorList(rtnRtes), nil
}

// ReadLast 读取实时快照值(当前标签点最后一个写入的TVQ)
//
// input:
//   - info: 点信息
//
// output:
//   - PTVQ(ptvq): 实时数值
func (c *RtdbConnect) ReadLast(info *PointInfo) (PTVQ, error) {
	ptvqs, errs, err := c.ReadLasts([]*PointInfo{info})
	if err != nil {
		return PTVQ{}, err
	}
	if errs[0] != nil {
		return PTVQ{}, errs[0]
	}
	return ptvqs[0], nil
}

// ReadValue 读取单个TVQ
//
// input:
//   - info 点信息
//   - mode 读取模式
//   - timestamp 时间戳
//
// output:
//   - PTVQ(ptvq) 输出TVQ点值
func (c *RtdbConnect) ReadValue(info *PointInfo, mode RtdbHisMode, timestamp time.Time) (PTVQ, error) {
	rtdbType, _ := info.ValueType.ToRawType()
	datetime, subtime := GoTimeToRtdbTimestamp(timestamp)
	switch rtdbType {
	case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
		dt, ms, value, state, quality, rte := RawRtdbhGetSingleValue64Warp(c.ConnectHandle, info.ID, mode, datetime, subtime)
		if !RteIsOk(rte) {
			return PTVQ{}, rte.GoError()
		}
		ts := RtdbTimestampToGoTime(dt, ms)
		switch rtdbType {
		case RtdbTypeBool:
			return NewPTVQ(info, NewTvqBool(ts, Int64ToBool(state), quality)), nil
		case RtdbTypeUint8:
			return NewPTVQ(info, NewTvqUint8(ts, uint8(state), quality)), nil
		case RtdbTypeInt8:
			return NewPTVQ(info, NewTvqInt8(ts, int8(state), quality)), nil
		case RtdbTypeChar:
			return NewPTVQ(info, NewTvqChar(ts, byte(state), quality)), nil
		case RtdbTypeUint16:
			return NewPTVQ(info, NewTvqUint16(ts, uint16(state), quality)), nil
		case RtdbTypeInt16:
			return NewPTVQ(info, NewTvqInt16(ts, int16(state), quality)), nil
		case RtdbTypeUint32:
			return NewPTVQ(info, NewTvqUint32(ts, uint32(state), quality)), nil
		case RtdbTypeInt32:
			return NewPTVQ(info, NewTvqInt32(ts, int32(state), quality)), nil
		case RtdbTypeInt64:
			return NewPTVQ(info, NewTvqInt64(ts, state, quality)), nil
		case RtdbTypeReal16:
			return NewPTVQ(info, NewTvqFloat16(ts, float32(value), quality)), nil
		case RtdbTypeReal32:
			return NewPTVQ(info, NewTvqFloat32(ts, float32(value), quality)), nil
		case RtdbTypeReal64:
			return NewPTVQ(info, NewTvqFloat64(ts, value, quality)), nil
		case RtdbTypeFp16:
			return NewPTVQ(info, NewTvqFp16(ts, float32(value), quality)), nil
		case RtdbTypeFp32:
			return NewPTVQ(info, NewTvqFp32(ts, float32(value), quality)), nil
		case RtdbTypeFp64:
			return NewPTVQ(info, NewTvqFp64(ts, value, quality)), nil
		}
	case RtdbTypeCoor:
		dt, ms, x, y, quality, rte := RawRtdbhGetSingleCoorValue64Warp(c.ConnectHandle, info.ID, mode, datetime, subtime)
		if !RteIsOk(rte) {
			return PTVQ{}, rte.GoError()
		}
		ts := RtdbTimestampToGoTime(dt, ms)
		return NewPTVQ(info, NewTvqCoordinates(ts, x, y, quality)), nil
	case RtdbTypeString, RtdbTypeBlob:
		dt, ms, blob, quality, rte := RawRtdbhGetSingleBlobValue64Warp(c.ConnectHandle, info.ID, mode, datetime, subtime, c.StringBlobMaxLen)
		if !RteIsOk(rte) {
			return PTVQ{}, rte.GoError()
		}
		ts := RtdbTimestampToGoTime(dt, ms)
		switch rtdbType {
		case RtdbTypeString:
			str := string(blob)
			if c.ServerOsType == RtdbOsWindows {
				s, err := GBKBytesToString(blob)
				if err != nil {
					return PTVQ{}, err
				}
				str = s
			}
			return NewPTVQ(info, NewTvqString(ts, str, quality)), nil
		case RtdbTypeBlob:
			return NewPTVQ(info, NewTvqBlob(ts, blob, quality)), nil
		}
	case RtdbTypeNamedT:
		dt, ms, data, quality, rte := RawRtdbhGetSingleNamedTypeValue64Warp(c.ConnectHandle, info.ID, mode, datetime, subtime, info.NamedType.Length)
		if !RteIsOk(rte) {
			return PTVQ{}, rte.GoError()
		}
		ts := RtdbTimestampToGoTime(dt, ms)
		return NewPTVQ(info, NewTvqNamed(ts, info.ValueType, data, quality)), nil
	case RtdbTypeDatetime:
		dt, ms, data, quality, rte := RawRtdbhGetSingleDatetimeValue64Warp(c.ConnectHandle, info.ID, mode, datetime, subtime, 0)
		if !RteIsOk(rte) {
			return PTVQ{}, rte.GoError()
		}
		ts := RtdbTimestampToGoTime(dt, ms)
		return NewPTVQ(info, NewTvqString(ts, string(data), quality)), nil
	}
	return PTVQ{}, errors.New("未知的数据类型")
}

// ReadRange 读取某个时间段内的TVQ
//
// input:
//   - info 标签点信息
//   - start 开始时间
//   - end 结束时间
//
// output:
//   - PTVQs(ptvqs) 点值列表
func (c *RtdbConnect) ReadRange(info *PointInfo, start time.Time, end time.Time) (PTVQs, error) {
	rtdbType, _ := info.ValueType.ToRawType()
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)
	maxCount, rte := RawRtdbhArchivedValuesCount64Warp(c.ConnectHandle, info.ID, datetime1, subtime1, datetime2, subtime2)
	if !RteIsOk(rte) {
		return PTVQs{}, rte.GoError()
	}
	tvqs := make([]TVQ, 0)
	switch rtdbType {
	case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
		dt, ms, values, states, qualities, rte := RawRtdbhGetArchivedValues64Warp(c.ConnectHandle, info.ID, maxCount, datetime1, subtime1, datetime2, subtime2)
		if !RteIsOk(rte) {
			return PTVQs{}, rte.GoError()
		}
		for i := 0; i < len(dt); i++ {
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			switch rtdbType {
			case RtdbTypeBool:
				tvqs = append(tvqs, NewTvqBool(ts, Int64ToBool(states[i]), q))
			case RtdbTypeUint8:
				tvqs = append(tvqs, NewTvqUint8(ts, uint8(states[i]), q))
			case RtdbTypeInt8:
				tvqs = append(tvqs, NewTvqInt8(ts, int8(states[i]), q))
			case RtdbTypeChar:
				tvqs = append(tvqs, NewTvqChar(ts, byte(states[i]), q))
			case RtdbTypeUint16:
				tvqs = append(tvqs, NewTvqUint16(ts, uint16(states[i]), q))
			case RtdbTypeInt16:
				tvqs = append(tvqs, NewTvqInt16(ts, int16(states[i]), q))
			case RtdbTypeUint32:
				tvqs = append(tvqs, NewTvqUint32(ts, uint32(states[i]), q))
			case RtdbTypeInt32:
				tvqs = append(tvqs, NewTvqInt32(ts, int32(states[i]), q))
			case RtdbTypeInt64:
				tvqs = append(tvqs, NewTvqInt64(ts, states[i], q))
			case RtdbTypeReal16:
				tvqs = append(tvqs, NewTvqFloat16(ts, float32(values[i]), q))
			case RtdbTypeReal32:
				tvqs = append(tvqs, NewTvqFloat32(ts, float32(values[i]), q))
			case RtdbTypeReal64:
				tvqs = append(tvqs, NewTvqFloat64(ts, values[i], q))
			case RtdbTypeFp16:
				tvqs = append(tvqs, NewTvqFp16(ts, float32(values[i]), q))
			case RtdbTypeFp32:
				tvqs = append(tvqs, NewTvqFp32(ts, float32(values[i]), q))
			case RtdbTypeFp64:
				tvqs = append(tvqs, NewTvqFp64(ts, values[i], q))
			}
		}
	case RtdbTypeCoor:
		dt, ms, xs, ys, qualities, rte := RawRtdbhGetArchivedCoorValues64Warp(c.ConnectHandle, info.ID, maxCount, datetime1, subtime1, datetime2, subtime2)
		if !RteIsOk(rte) {
			return PTVQs{}, rte.GoError()
		}
		for i := 0; i < len(dt); i++ {
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			tvqs = append(tvqs, NewTvqCoordinates(ts, xs[i], ys[i], q))
		}
	case RtdbTypeBlob, RtdbTypeString:
		dt, ms, blobs, qualities, rte := RawRtdbhGetArchivedBlobValues64Warp(c.ConnectHandle, info.ID, maxCount, datetime1, subtime1, datetime2, subtime2, c.StringBlobMaxLen)
		if !RteIsOk(rte) {
			return PTVQs{}, rte.GoError()
		}
		for i := 0; i < len(dt); i++ {
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			switch rtdbType {
			case RtdbTypeBlob:
				tvqs = append(tvqs, NewTvqBlob(ts, blobs[i], q))
			case RtdbTypeString:
				str := string(blobs[i])
				if c.ServerOsType == RtdbOsWindows {
					s, err := GBKBytesToString(blobs[i])
					if err != nil {
						return PTVQs{}, rte.GoError()
					}
					str = s
				}
				tvqs = append(tvqs, NewTvqString(ts, str, q))
			}
		}
	case RtdbTypeNamedT:
		dt, ms, datas, qualities, rte := RawRtdbhGetArchivedNamedTypeValues64Warp(c.ConnectHandle, info.ID, maxCount, datetime1, subtime1, datetime2, subtime2, info.NamedType.Length)
		if !RteIsOk(rte) {
			return PTVQs{}, rte.GoError()
		}
		for i := 0; i < len(dt); i++ {
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			tvqs = append(tvqs, NewTvqNamed(ts, info.ValueType, datas[i], q))
		}
	case RtdbTypeDatetime:
		dt, ms, datetimes, qualities, rte := RawRtdbhGetArchivedDatetimeValues64Warp(c.ConnectHandle, info.ID, maxCount, datetime1, subtime1, datetime2, subtime2, 0)
		if !RteIsOk(rte) {
			return PTVQs{}, rte.GoError()
		}
		for i := 0; i < len(dt); i++ {
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			tvqs = append(tvqs, NewTvqDatetime(ts, string(datetimes[i]), q))
		}
	}

	return NewPTVQs(info, tvqs), nil
}

// ReadPlot 读取用于绘图的TVQ
//
// input:
//   - info 标签点信息
//   - interval 像素个数，比如1024像素的屏幕，interval建议少于1024
//   - start 开始时间
//   - end 结束时间
//
// output:
//   - PTVQs(ptvqs) 点值列表
func (c *RtdbConnect) ReadPlot(info *PointInfo, interval int32, start time.Time, end time.Time) (PTVQs, error) {
	rtdbType, _ := info.ValueType.ToRawType()
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)
	dt, ms, values, states, qualities, rte := RawRtdbhGetPlotValues64Warp(c.ConnectHandle, info.ID, interval, datetime1, subtime1, datetime2, subtime2)
	if !RteIsOk(rte) {
		return PTVQs{}, rte.GoError()
	}
	tvqs := make([]TVQ, 0)
	for i := 0; i < len(dt); i++ {
		ts := RtdbTimestampToGoTime(dt[i], ms[i])
		q := qualities[i]
		switch rtdbType {
		case RtdbTypeBool:
			tvqs = append(tvqs, NewTvqBool(ts, Int64ToBool(states[i]), q))
		case RtdbTypeUint8:
			tvqs = append(tvqs, NewTvqUint8(ts, uint8(states[i]), q))
		case RtdbTypeInt8:
			tvqs = append(tvqs, NewTvqInt8(ts, int8(states[i]), q))
		case RtdbTypeChar:
			tvqs = append(tvqs, NewTvqChar(ts, byte(states[i]), q))
		case RtdbTypeUint16:
			tvqs = append(tvqs, NewTvqUint16(ts, uint16(states[i]), q))
		case RtdbTypeInt16:
			tvqs = append(tvqs, NewTvqInt16(ts, int16(states[i]), q))
		case RtdbTypeUint32:
			tvqs = append(tvqs, NewTvqUint32(ts, uint32(states[i]), q))
		case RtdbTypeInt32:
			tvqs = append(tvqs, NewTvqInt32(ts, int32(states[i]), q))
		case RtdbTypeInt64:
			tvqs = append(tvqs, NewTvqInt64(ts, states[i], q))
		case RtdbTypeReal16:
			tvqs = append(tvqs, NewTvqFloat16(ts, float32(values[i]), q))
		case RtdbTypeReal32:
			tvqs = append(tvqs, NewTvqFloat32(ts, float32(values[i]), q))
		case RtdbTypeReal64:
			tvqs = append(tvqs, NewTvqFloat64(ts, values[i], q))
		case RtdbTypeFp16:
			tvqs = append(tvqs, NewTvqFp16(ts, float32(values[i]), q))
		case RtdbTypeFp32:
			tvqs = append(tvqs, NewTvqFp32(ts, float32(values[i]), q))
		case RtdbTypeFp64:
			tvqs = append(tvqs, NewTvqFp64(ts, values[i], q))
		}
	}
	return NewPTVQs(info, tvqs), nil
}

// ReadTimed 获取差值，每个差值都需要指定一个确定的时间戳
//
// input:
//   - info 标签点信息
//   - timestamps 时间戳列表
//
// output:
//   - PTVQs(ptvqs) 点值列表
func (c *RtdbConnect) ReadTimed(info *PointInfo, timestamps []time.Time) (PTVQs, error) {
	rtdbType, _ := info.ValueType.ToRawType()
	datetimes := make([]TimestampType, 0)
	subtimes := make([]SubtimeType, 0)
	for _, ts := range timestamps {
		dt, ms := GoTimeToRtdbTimestamp(ts)
		datetimes = append(datetimes, dt)
		subtimes = append(subtimes, ms)
	}
	tvqs := make([]TVQ, 0)
	switch rtdbType {
	case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
		values, states, qualities, rte := RawRtdbhGetTimedValues64Warp(c.ConnectHandle, info.ID, datetimes, subtimes)
		if !RteIsOk(rte) {
			return PTVQs{}, rte.GoError()
		}
		for i := 0; i < len(timestamps); i++ {
			ts := timestamps[i]
			q := qualities[i]
			switch rtdbType {
			case RtdbTypeBool:
				tvqs = append(tvqs, NewTvqBool(ts, Int64ToBool(states[i]), q))
			case RtdbTypeUint8:
				tvqs = append(tvqs, NewTvqUint8(ts, uint8(states[i]), q))
			case RtdbTypeInt8:
				tvqs = append(tvqs, NewTvqInt8(ts, int8(states[i]), q))
			case RtdbTypeChar:
				tvqs = append(tvqs, NewTvqChar(ts, byte(states[i]), q))
			case RtdbTypeUint16:
				tvqs = append(tvqs, NewTvqUint16(ts, uint16(states[i]), q))
			case RtdbTypeInt16:
				tvqs = append(tvqs, NewTvqInt16(ts, int16(states[i]), q))
			case RtdbTypeUint32:
				tvqs = append(tvqs, NewTvqUint32(ts, uint32(states[i]), q))
			case RtdbTypeInt32:
				tvqs = append(tvqs, NewTvqInt32(ts, int32(states[i]), q))
			case RtdbTypeInt64:
				tvqs = append(tvqs, NewTvqInt64(ts, states[i], q))
			case RtdbTypeReal16:
				tvqs = append(tvqs, NewTvqFloat16(ts, float32(values[i]), q))
			case RtdbTypeReal32:
				tvqs = append(tvqs, NewTvqFloat32(ts, float32(values[i]), q))
			case RtdbTypeReal64:
				tvqs = append(tvqs, NewTvqFloat64(ts, values[i], q))
			case RtdbTypeFp16:
				tvqs = append(tvqs, NewTvqFp16(ts, float32(values[i]), q))
			case RtdbTypeFp32:
				tvqs = append(tvqs, NewTvqFp32(ts, float32(values[i]), q))
			case RtdbTypeFp64:
				tvqs = append(tvqs, NewTvqFp64(ts, values[i], q))
			}
		}
		return NewPTVQs(info, tvqs), nil
	case RtdbTypeCoor:
		xs, ys, qualities, rte := RawRtdbhGetTimedCoorValues64Warp(c.ConnectHandle, info.ID, datetimes, subtimes)
		if !RteIsOk(rte) {
			return PTVQs{}, rte.GoError()
		}
		for i := 0; i < len(timestamps); i++ {
			ts := timestamps[i]
			q := qualities[i]
			tvqs = append(tvqs, NewTvqCoordinates(ts, xs[i], ys[i], q))
		}
		return NewPTVQs(info, tvqs), nil
	default:
		return PTVQs{}, errors.New("不支持的数据类型")
	}
}

// ReadInterpo 获取差值, 会自动将start、end等分成count个时间戳，然后取这些时间戳的差值
//
// input:
//   - info 标签点信息
//   - count 要查询差值的数量
//   - start 开始时间
//   - end 结束时间
//   - filter 过滤条件, 可为""
//
// output:
//   - PTVQs(ptvqs) 点值列表
func (c *RtdbConnect) ReadInterpo(info *PointInfo, count int32, start time.Time, end time.Time, filter string) (PTVQs, error) {
	rtdbType, _ := info.ValueType.ToRawType()
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)

	var dt []TimestampType
	var ms []SubtimeType
	var values []float64
	var states []int64
	var qualities []Quality
	var rte RtdbError

	if filter == "" {
		dt, ms, values, states, qualities, rte = RawRtdbhGetInterpoValues64Warp(c.ConnectHandle, info.ID, count, datetime1, subtime1, datetime2, subtime2)
	} else {
		dt, ms, values, states, qualities, rte = RawRtdbhGetInterpoValuesFilt64Warp(c.ConnectHandle, info.ID, filter, count, datetime1, subtime1, datetime2, subtime2)
	}
	if !RteIsOk(rte) {
		return PTVQs{}, rte.GoError()
	}

	switch rtdbType {
	case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
		tvqs := make([]TVQ, 0)
		for i := 0; i < len(dt); i++ {
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			switch rtdbType {
			case RtdbTypeBool:
				tvqs = append(tvqs, NewTvqBool(ts, Int64ToBool(states[i]), q))
			case RtdbTypeUint8:
				tvqs = append(tvqs, NewTvqUint8(ts, uint8(states[i]), q))
			case RtdbTypeInt8:
				tvqs = append(tvqs, NewTvqInt8(ts, int8(states[i]), q))
			case RtdbTypeChar:
				tvqs = append(tvqs, NewTvqChar(ts, byte(states[i]), q))
			case RtdbTypeUint16:
				tvqs = append(tvqs, NewTvqUint16(ts, uint16(states[i]), q))
			case RtdbTypeInt16:
				tvqs = append(tvqs, NewTvqInt16(ts, int16(states[i]), q))
			case RtdbTypeUint32:
				tvqs = append(tvqs, NewTvqUint32(ts, uint32(states[i]), q))
			case RtdbTypeInt32:
				tvqs = append(tvqs, NewTvqInt32(ts, int32(states[i]), q))
			case RtdbTypeInt64:
				tvqs = append(tvqs, NewTvqInt64(ts, states[i], q))
			case RtdbTypeReal16:
				tvqs = append(tvqs, NewTvqFloat16(ts, float32(values[i]), q))
			case RtdbTypeReal32:
				tvqs = append(tvqs, NewTvqFloat32(ts, float32(values[i]), q))
			case RtdbTypeReal64:
				tvqs = append(tvqs, NewTvqFloat64(ts, values[i], q))
			case RtdbTypeFp16:
				tvqs = append(tvqs, NewTvqFp16(ts, float32(values[i]), q))
			case RtdbTypeFp32:
				tvqs = append(tvqs, NewTvqFp32(ts, float32(values[i]), q))
			case RtdbTypeFp64:
				tvqs = append(tvqs, NewTvqFp64(ts, values[i], q))
			}
		}
		return NewPTVQs(info, tvqs), nil
	default:
		return PTVQs{}, errors.New("不支持的数据类型")
	}
}

// ReadInterval 读取从start开始的等间隔差值
//
// input:
//   - info 标签点信息
//   - filter 过滤条件, 可为""
//   - start 开始时间
//   - interval 时间间隔
//   - count 要读取差值的个数
//
// output:
//   - PTVQs(ptvqs) 点值列表
func (c *RtdbConnect) ReadInterval(info *PointInfo, filter string, start time.Time, interval time.Duration, count int32) (PTVQs, error) {
	rtdbType, _ := info.ValueType.ToRawType()
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)

	var dt []TimestampType
	var ms []SubtimeType
	var values []float64
	var states []int64
	var qualities []Quality
	var rte RtdbError

	if filter == "" {
		dt, ms, values, states, qualities, rte = RawRtdbhGetIntervalValues64Warp(c.ConnectHandle, info.ID, interval, count, datetime1, subtime1)
	} else {
		dt, ms, values, states, qualities, rte = RawRtdbhGetIntervalValuesFilt64Warp(c.ConnectHandle, info.ID, filter, interval, count, datetime1, subtime1)
	}

	if !RteIsOk(rte) {
		return PTVQs{}, rte.GoError()
	}

	switch rtdbType {
	case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
		tvqs := make([]TVQ, 0)
		for i := 0; i < len(dt); i++ {
			ts := RtdbTimestampToGoTime(dt[i], ms[i])
			q := qualities[i]
			switch rtdbType {
			case RtdbTypeBool:
				tvqs = append(tvqs, NewTvqBool(ts, Int64ToBool(states[i]), q))
			case RtdbTypeUint8:
				tvqs = append(tvqs, NewTvqUint8(ts, uint8(states[i]), q))
			case RtdbTypeInt8:
				tvqs = append(tvqs, NewTvqInt8(ts, int8(states[i]), q))
			case RtdbTypeChar:
				tvqs = append(tvqs, NewTvqChar(ts, byte(states[i]), q))
			case RtdbTypeUint16:
				tvqs = append(tvqs, NewTvqUint16(ts, uint16(states[i]), q))
			case RtdbTypeInt16:
				tvqs = append(tvqs, NewTvqInt16(ts, int16(states[i]), q))
			case RtdbTypeUint32:
				tvqs = append(tvqs, NewTvqUint32(ts, uint32(states[i]), q))
			case RtdbTypeInt32:
				tvqs = append(tvqs, NewTvqInt32(ts, int32(states[i]), q))
			case RtdbTypeInt64:
				tvqs = append(tvqs, NewTvqInt64(ts, states[i], q))
			case RtdbTypeReal16:
				tvqs = append(tvqs, NewTvqFloat16(ts, float32(values[i]), q))
			case RtdbTypeReal32:
				tvqs = append(tvqs, NewTvqFloat32(ts, float32(values[i]), q))
			case RtdbTypeReal64:
				tvqs = append(tvqs, NewTvqFloat64(ts, values[i], q))
			case RtdbTypeFp16:
				tvqs = append(tvqs, NewTvqFp16(ts, float32(values[i]), q))
			case RtdbTypeFp32:
				tvqs = append(tvqs, NewTvqFp32(ts, float32(values[i]), q))
			case RtdbTypeFp64:
				tvqs = append(tvqs, NewTvqFp64(ts, values[i], q))
			}
		}
		return NewPTVQs(info, tvqs), nil
	default:
		return PTVQs{}, errors.New("不支持的数据类型")
	}
}

// ReadSection 读取断面
//
// input:
//   - infos 标签点信息列表
//   - mode 读取规则
//   - timestamp 时间戳
//
// output:
//   - []PTVQ(ptvqs) 点值列表
//   - []error(errs) 错误列表
func (c *RtdbConnect) ReadSection(infos []*PointInfo, mode RtdbHisMode, timestamp time.Time) ([]PTVQ, []error, error) {
	numberInfos := make([]*PointInfo, 0)
	numberIds := make([]PointID, 0)
	numberIdx := make([]int, 0)

	otherInfos := make([]*PointInfo, 0)
	otherIdx := make([]int, 0)

	for i, info := range infos {
		rtdbType, _ := info.ValueType.ToRawType()
		switch rtdbType {
		case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
			numberInfos = append(numberInfos, info)
			numberIds = append(numberIds, info.ID)
			numberIdx = append(numberIdx, i)
		default:
			otherInfos = append(otherInfos, info)
			otherIdx = append(otherIdx, i)
		}
	}

	datetime, subtime := GoTimeToRtdbTimestamp(timestamp)

	numberPtvqs := make([]PTVQ, 0)
	numberRtes := make([]RtdbError, 0)
	if len(numberInfos) != 0 {
		dt, ms, values, states, qualities, rtes, rte := RawRtdbhGetCrossSectionValues64Warp(c.ConnectHandle, numberIds, mode, datetime, subtime)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		numberRtes = rtes

		for i, info := range numberInfos {
			rtdbType, _ := info.ValueType.ToRawType()
			switch rtdbType {
			case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
				ts := RtdbTimestampToGoTime(dt[i], ms[i])
				q := qualities[i]
				switch rtdbType {
				case RtdbTypeBool:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqBool(ts, Int64ToBool(states[i]), q)))
				case RtdbTypeUint8:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqUint8(ts, uint8(states[i]), q)))
				case RtdbTypeInt8:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqInt8(ts, int8(states[i]), q)))
				case RtdbTypeChar:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqChar(ts, byte(states[i]), q)))
				case RtdbTypeUint16:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqUint16(ts, uint16(states[i]), q)))
				case RtdbTypeInt16:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqInt16(ts, int16(states[i]), q)))
				case RtdbTypeUint32:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqUint32(ts, uint32(states[i]), q)))
				case RtdbTypeInt32:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqInt32(ts, int32(states[i]), q)))
				case RtdbTypeInt64:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqInt64(ts, states[i], q)))
				case RtdbTypeReal16:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqFloat16(ts, float32(values[i]), q)))
				case RtdbTypeReal32:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqFloat32(ts, float32(values[i]), q)))
				case RtdbTypeReal64:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqFloat64(ts, values[i], q)))
				case RtdbTypeFp16:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqFp16(ts, float32(values[i]), q)))
				case RtdbTypeFp32:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqFp32(ts, float32(values[i]), q)))
				case RtdbTypeFp64:
					numberPtvqs = append(numberPtvqs, NewPTVQ(info, NewTvqFp64(ts, values[i], q)))
				}
			}
		}
	}

	otherPtvqs := make([]PTVQ, 0)
	otherRtes := make([]error, 0)
	for _, info := range otherInfos {
		ptvq, rte := c.ReadValue(info, mode, timestamp)
		otherPtvqs = append(otherPtvqs, ptvq)
		otherRtes = append(otherRtes, rte)
	}

	rtnPtvqs := make([]PTVQ, len(infos))
	rtnRtes := make([]error, len(infos))
	for i := 0; i < len(numberInfos); i++ {
		idx := numberIdx[i]
		rtnPtvqs[idx] = numberPtvqs[i]
		rtnRtes[idx] = numberRtes[i].GoError()
	}
	for i := 0; i < len(otherInfos); i++ {
		idx := otherIdx[i]
		rtnPtvqs[idx] = otherPtvqs[i]
		rtnRtes[idx] = otherRtes[i]
	}

	return rtnPtvqs, rtnRtes, nil
}

// ReadSummary 获取统计值 (从start到end的统计值)
//
// input:
//   - info 标签点信息
//   - filter 过滤条件, 可为"“
//   - start 开始时间
//   - end 结束时间
func (c *RtdbConnect) ReadSummary(info *PointInfo, filter string, start time.Time, end time.Time) (*RtdbSummaryData, error) {
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)
	var summary *RtdbSummaryData
	var rte RtdbError
	if filter == "" {
		summary, rte = RawRtdbhSummaryDataWarp(c.ConnectHandle, info.ID, datetime1, subtime1, datetime2, subtime2)
	} else {
		summary, rte = RawRtdbhSummaryDataFiltWarp(c.ConnectHandle, info.ID, filter, datetime1, subtime1, datetime2, subtime2)
	}
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return summary, nil
}

// ReadBatchesSummary 获取等间隔统计值 (在start和end之间，按照interval作为时间间隔计算每个间隔的统计值)
//
// input:
//   - info 标签点信息
//   - filter 过滤条件, 可为""
//   - interval 时间间隔
//   - start 开始时间
//   - end 结束时间
func (c *RtdbConnect) ReadBatchesSummary(info *PointInfo, filter string, interval time.Duration, start time.Time, end time.Time) ([]RtdbSummaryData, []error, error) {
	maxCount := int32((end.Sub(start) / interval) + 1)
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)

	var summaryDatas []RtdbSummaryData
	var rtes []RtdbError
	var rte RtdbError
	if filter == "" {
		summaryDatas, rtes, rte = RawRtdbhSummaryDataInBatchesWarp(c.ConnectHandle, info.ID, maxCount, interval, datetime1, subtime1, datetime2, subtime2)
	} else {
		summaryDatas, rtes, rte = RawRtdbhSummaryDataFiltInBatchesWarp(c.ConnectHandle, info.ID, filter, maxCount, interval, datetime1, subtime1, datetime2, subtime2)
	}

	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}

	return summaryDatas, RtdbErrorListToErrorList(rtes), nil
}

// RemoveValue 删除点值
//
// input:
//   - info 标签点信息
//   - timestamp 时间戳
func (c *RtdbConnect) RemoveValue(info *PointInfo, timestamp time.Time) error {
	datetime, subtime := GoTimeToRtdbTimestamp(timestamp)
	rte := RawRtdbhRemoveValue64Warp(c.ConnectHandle, info.ID, datetime, subtime)
	return rte.GoError()
}

// RemoveRangeValues 批量删除点值(从start到end的范围)
//
// input:
//   - info 标签点信息
//   - start 开始时间
//   - end 结束时间
//
// output:
//   - int32(count) 删除的点值个数
func (c *RtdbConnect) RemoveRangeValues(info *PointInfo, start time.Time, end time.Time) (int32, error) {
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)
	count, rte := RawRtdbhRemoveValues64Warp(c.ConnectHandle, info.ID, datetime1, subtime1, datetime2, subtime2)
	if !RteIsOk(rte) {
		return 0, rte.GoError()
	}
	return count, nil
}

// UpdateValue
//
// input:
//   - info 标签点信息
//   - tvq 点值
func (c *RtdbConnect) UpdateValue(info *PointInfo, tvq TVQ) error {
	datetime, subtime := tvq.GetRtdbTimestamp()
	quality := tvq.GetRtdbQuality()
	rtdbType, _ := info.ValueType.ToRawType()
	switch rtdbType {
	case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64, RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
		var value float64
		var state int64
		switch rtdbType {
		case RtdbTypeBool, RtdbTypeUint8, RtdbTypeInt8, RtdbTypeChar, RtdbTypeUint16, RtdbTypeInt16, RtdbTypeUint32, RtdbTypeInt32, RtdbTypeInt64:
			state = tvq.GetRtdbInt()
		case RtdbTypeReal16, RtdbTypeReal32, RtdbTypeReal64, RtdbTypeFp16, RtdbTypeFp32, RtdbTypeFp64:
			value = tvq.GetRtdbFloat()
		}
		rte := RawRtdbhUpdateValue64Warp(c.ConnectHandle, info.ID, datetime, subtime, value, state, quality)
		return rte.GoError()
	case RtdbTypeCoor:
		coorValue := tvq.GetRtdbCoordinates()
		rte := RawRtdbhUpdateCoorValue64Warp(c.ConnectHandle, info.ID, datetime, subtime, coorValue.X, coorValue.Y, quality)
		return rte.GoError()
	default:
		return errors.New("不支持的数据类型")
	}
}

// FlushArchivedValues 刷新历史缓存(把标签点的缓存进行手动归档)
//
// input:
//   - info 标签点信息
//
// output:
//   - int32(count) 补历史点值个数
func (c *RtdbConnect) FlushArchivedValues(info *PointInfo) (int32, error) {
	count, rte := RawRtdbhFlushArchivedValuesWarp(c.ConnectHandle, info.ID)
	if !RteIsOk(rte) {
		return 0, rte.GoError()
	}
	return count, nil
}

// QueryBigJob 查询正在进行的大任务
//
// input:
//   - processName 服务名称(大任务引擎被分为多个服务，相当于是一个线程，每个服务处理一个特定种类的大任务)
//
// output:
//   - string(path) 大任务正在处理文件的路径
//   - BigJobName(bitJobName) 大任务名称
//   - RtdbError(state) 大任务的执行状态
//   - time.Time(timestamp) 大任务的时间戳
//   - float32(process) 大任务处理进度
func (c *RtdbConnect) QueryBigJob(processName RtdbProcess) (string, BigJobName, RtdbError, time.Time, float32, error) {
	dir, fileName, bigJobName, state, timestamp, process, rte := RawRtdbaQueryBigJob64Warp(c.ConnectHandle, processName)
	if !RteIsOk(rte) {
		return "", bigJobName, 0, time.Unix(0, 0), 0, rte.GoError()
	}
	return path.Join(dir, fileName), bigJobName, state, time.Unix(int64(timestamp), 0), process, nil
}

// CancelBigJob 取消正在执行的大任务
//
// input:
//   - processName 服务名称(大任务引擎被分为多个服务，相当于是一个线程，每个服务处理一个特定种类的大任务)
func (c *RtdbConnect) CancelBigJob(processName RtdbProcess) error {
	rte := RawRtdbaCancelBigJobWarp(c.ConnectHandle, processName)
	return rte.GoError()
}

// JobMessage 获取任务描述
//
// input:
//   - jobID 任务ID
//
// output:
//   - string(name) 任务名称
//   - string(desc) 任务描述
func (c *RtdbConnect) JobMessage(jobID int32) (string, string) {
	name, desc := RawRtdbJobMessageWarp(jobID)
	return name, desc
}

// ComputeHistory 重算｜补算 补历史值或者修改历史值之后，对应的计算点可以按需进行重算|补算，保证计算点的数值正确
//
// input:
//   - infos 标签点信息列表
//   - flag 不为 0 表示进行重算，删除时间范围内已经存在历史数据；为 0 表示补算，保留时间范围内已经存在历史数据，覆盖同时刻的计算值。
//   - start 开始时间
//   - end 结束时间
//
// output:
//   - []error(errs) 每个标签点的重算结果
func (c *RtdbConnect) ComputeHistory(infos []*PointInfo, flag int16, start time.Time, end time.Time) ([]error, error) {
	datetime1, subtime1 := GoTimeToRtdbTimestamp(start)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)
	ids := make([]PointID, 0)
	for _, info := range infos {
		ids = append(ids, info.ID)
	}
	rtes, rte := RawRtdbeComputeHistory64Warp(c.ConnectHandle, ids, flag, datetime1, subtime1, datetime2, subtime2)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return RtdbErrorListToErrorList(rtes), nil
}

// GetPointEquation 获取标签点方程式
//
// input:
//   - info 标签点信息
func (c *RtdbConnect) GetPointEquation(info *PointInfo) ([]byte, error) {
	equation, rte := RawRtdbbGetEquationByIdWarp(c.ConnectHandle, info.ID)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return equation, nil
}

// GetEquationGraph 获取标签点对应的方程式关联关系图(方程式本身是一个有向无环图，数据在这个图内的公式之间流转)
//
// input:
//   - info 标签点信息
//   - flag 拓扑关系类型
//
// output:
//   - []RtdbGraph(graph) 关系图
func (c *RtdbConnect) GetEquationGraph(info *PointInfo, flag RtdbGraphFlag) ([]RtdbGraph, error) {
	count, rte := RawRtdbeGetEquationGraphCountWarp(c.ConnectHandle, info.ID, flag)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	graph, rte := RawRtdbeGetEquationGraphDatasWarp(c.ConnectHandle, info.ID, flag, count)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return graph, nil
}

// GetPerfPointInfo 获取性能监控点的说明信息
//
// input:
//   - ids 性能监控点ID列表
//
// output:
//   - []RtdbPerfTagInfo(infos) 说明信息列表
//   - []error(errs) 错误列表
func (c *RtdbConnect) GetPerfPointInfo(ids []RtdbPerfTagID) ([]RtdbPerfTagInfo, []error, error) {
	infos, rtes, rte := RawRtdbpGetPerfTagsInfoWarp(c.ConnectHandle, ids)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}
	return infos, RtdbErrorListToErrorList(rtes), nil
}

// GetPerfPointValue 获取性能监控点的实时值
//
// input:
//   - ids []RtdbPerfTagID 性能监控点ID列表
//
// output:
//   - []TVQ 实时值列表
//   - []error 错误列表
func (c *RtdbConnect) GetPerfPointValue(ids []RtdbPerfTagID) ([]TVQ, []error, error) {
	infos, _, err := c.GetPerfPointInfo(ids)
	if err != nil {
		return nil, nil, err
	}
	datetimes, subtimes, values, states, qualities, rtes, rte := RawRtdbpGetPerfValues64Warp(c.ConnectHandle, ids)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}

	tvqs := make([]TVQ, 0)
	for i := 0; i < len(ids); i++ {
		ts := RtdbTimestampToGoTime(datetimes[i], subtimes[i])
		q := qualities[i]

		switch infos[i].Type {
		case RtdbTypeBool:
			tvqs = append(tvqs, NewTvqBool(ts, Int64ToBool(states[i]), q))
		case RtdbTypeUint8:
			tvqs = append(tvqs, NewTvqUint8(ts, uint8(states[i]), q))
		case RtdbTypeInt8:
			tvqs = append(tvqs, NewTvqInt8(ts, int8(states[i]), q))
		case RtdbTypeChar:
			tvqs = append(tvqs, NewTvqChar(ts, byte(states[i]), q))
		case RtdbTypeUint16:
			tvqs = append(tvqs, NewTvqUint16(ts, uint16(states[i]), q))
		case RtdbTypeInt16:
			tvqs = append(tvqs, NewTvqInt16(ts, int16(states[i]), q))
		case RtdbTypeUint32:
			tvqs = append(tvqs, NewTvqUint32(ts, uint32(states[i]), q))
		case RtdbTypeInt32:
			tvqs = append(tvqs, NewTvqInt32(ts, int32(states[i]), q))
		case RtdbTypeInt64:
			tvqs = append(tvqs, NewTvqInt64(ts, states[i], q))
		case RtdbTypeReal16:
			tvqs = append(tvqs, NewTvqFloat16(ts, float32(values[i]), q))
		case RtdbTypeReal32:
			tvqs = append(tvqs, NewTvqFloat32(ts, float32(values[i]), q))
		case RtdbTypeReal64:
			tvqs = append(tvqs, NewTvqFloat64(ts, values[i], q))
		case RtdbTypeFp16:
			tvqs = append(tvqs, NewTvqFp16(ts, float32(values[i]), q))
		case RtdbTypeFp32:
			tvqs = append(tvqs, NewTvqFp32(ts, float32(values[i]), q))
		case RtdbTypeFp64:
			tvqs = append(tvqs, NewTvqFp64(ts, values[i], q))
		default:
			return nil, nil, errors.New("不支持的数据类型，分支不可达")
		}
	}
	return tvqs, RtdbErrorListToErrorList(rtes), nil
}

// CreateRangedArchive 创建存档文件(按照时间范围创建)
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
//   - begin 开始时间
//   - end 结束时间
//   - mbSize 存档文件大小，单位MB
func (c *RtdbConnect) CreateRangedArchive(path string, file string, begin time.Time, end time.Time, mbSize int32) error {
	datetime1, subtime1 := GoTimeToRtdbTimestamp(begin)
	datetime2, subtime2 := GoTimeToRtdbTimestamp(end)
	if subtime1 != 0 || subtime2 != 0 {
		return errors.New("begin、end时间的亚秒部分应该为0")
	}
	rte := RawRtdbaCreateRangedArchive64Warp(c.ConnectHandle, path, file, datetime1, datetime2, mbSize)
	return rte.GoError()
}

// AppendArchive 存档文件入列(将存档文件插入到存档队列中, 队列中的存档文件按照start、end依次排列，注意时间段不能有交叠的部分)
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
//   - state 存档文件状态
func (c *RtdbConnect) AppendArchive(path string, file string, state RtdbArchiveState) error {
	rte := RawRtdbaAppendArchiveWarp(c.ConnectHandle, path, file, state)
	return rte.GoError()
}

// RemoveArchive 存档文件解列
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
func (c *RtdbConnect) RemoveArchive(path, file string) error {
	rte := RawRtdbaRemoveArchiveWarp(c.ConnectHandle, path, file)
	return rte.GoError()
}

// ShiftActived 切换活动文件
//
// remark:
//   - 当前活动文件被写满时该事务被启动，
//   - 改变当前活动文件的状态为普通状态，
//   - 在所有历史数据存档文件中寻找未被使用过的
//   - 插入到前活动文件的右侧并改为活动状态，
//   - 若找不到则将前活动文件右侧的文件改为活动状态，
//   - 并将active_archive_指向该文件。该事务进行过程中，
//   - 用锁保证所有读写操作都暂停等待该事务完成。
func (c *RtdbConnect) ShiftActived() error {
	rte := RawRtdbaShiftActivedWarp(c.ConnectHandle)
	return rte.GoError()
}

// GetArchives 获取存档文件的基本信息(返回全部存档文件的基本信息)
func (c *RtdbConnect) GetArchives() ([]ArchivedBaseInfo, []error, error) {
	count, rte := RawRtdbaGetArchivesCountWarp(c.ConnectHandle)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}
	paths, files, states, rte := RawRtdbaGetArchivesWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}
	_, _, headPages, rtes, rte := RawRtdbaGetArchivesInfoWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}

	_, _, perfInfos, sumPerfInfos, rtes, rte := RawRtdbaGetArchivesPerfDataWarp(c.ConnectHandle, count)
	if !RteIsOk(rte) {
		return nil, nil, rte.GoError()
	}

	infos := make([]ArchivedBaseInfo, 0)
	for i := 0; i < len(paths); i++ {
		infos = append(infos, ArchivedBaseInfo{
			Path:        paths[i],
			File:        files[i],
			State:       states[i],
			HeadPage:    headPages[i],
			PerfInfo:    perfInfos[i],
			SumPerfInfo: sumPerfInfos[i],
		})
	}
	return infos, RtdbErrorListToErrorList(rtes), nil
}

// UpdateArchive 更新存档文件选项
//
// input:
//   - handle 连接句柄
//   - path 文件所在目录路径，必须以"\"或"/"结尾。
//   - file 文件名。
//   - rated_capacity 文件额定大小，单位为 MB。
//   - ex_capacity 附属文件大小，单位为 MB。
//   - auto_merge 是否自动合并附属文件。
//   - auto_arrange 是否自动整理存档文件。
//   - 备注: rated_capacity 与 ex_capacity 参数可为 0，表示不修改对应的配置项。
func (c *RtdbConnect) UpdateArchive(path string, file string, ratedCapacity int32, exCapacity int32, autoMerge int16, autoArrange int16) error {
	rte := RawRtdbaUpdateArchiveWarp(c.ConnectHandle, path, file, ratedCapacity, exCapacity, autoMerge, autoArrange)
	return rte.GoError()
}

// ArrangeArchive 整理存档文件(会重新调整数据块分布，使存档文件更加紧凑)
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
func (c *RtdbConnect) ArrangeArchive(path string, file string) error {
	rte := RawRtdbaArrangeArchiveWarp(c.ConnectHandle, path, file)
	return rte.GoError()
}

// ReindexArchive 重建存档文件索引(用于进行数据恢复)
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
func (c *RtdbConnect) ReindexArchive(path string, file string) error {
	rte := RawRtdbaReindexArchiveWarp(c.ConnectHandle, path, file)
	return rte.GoError()
}

// ConvertIndex 存档文件索引转换格式(老版本索引格式转换为新版本，新版本索引更快)
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
func (c *RtdbConnect) ConvertIndex(path string, file string) error {
	rte := RawRtdbaConvertIndexWarp(c.ConnectHandle, path, file)
	return rte.GoError()
}

// BackupArchive 备份存档
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
//   - dist 备份存档路径
func (c *RtdbConnect) BackupArchive(path string, file string, dist string) error {
	rte := RawRtdbaBackupArchiveWarp(c.ConnectHandle, path, file, dist)
	return rte.GoError()
}

// MoveArchive 移动存档
//
// input:
//   - path 存档文件路径
//   - file 存档文件名称
//   - dist 移动存档路径
func (c *RtdbConnect) MoveArchive(path string, file string, dist string) error {
	rte := RawRtdbaMoveArchiveWarp(c.ConnectHandle, path, file, dist)
	return rte.GoError()
}

// CreateDatagram 创建数据流订阅
//
// input:
//   - port 端口
//   - remoteHost 对端IP地址
//
// output:
//   - DatagramHandle(handle) 数据流句柄
func (c *RtdbConnect) CreateDatagram(port int32, remoteHost string) (DatagramHandle, error) {
	datagramHandle, rte := RawRtdbCreateDatagramHandleWarp(port, remoteHost)
	if !RteIsOk(rte) {
		return DatagramHandle{}, rte.GoError()
	}
	return datagramHandle, nil
}

// RemoveDatagram 取消数据流订阅
//
// input:
//   - handle 数据流句柄
func (c *RtdbConnect) RemoveDatagram(handle DatagramHandle) error {
	rte := RawRtdbRemoveDatagramHandleWarp(handle)
	return rte.GoError()
}

// RecvDatagram 从订阅数据流中获取数据
//
// input:
//   - handle  数据流句柄
//   - cacheLen 缓存大小，会创建对应大小的缓存，用于接收数据流返回的数据
//   - remoteAddr 对端IP地址
//   - timeout 超时时间(单位秒)
//
// output:
//   - []byte(message) 返回数据流
func (c *RtdbConnect) RecvDatagram(handle DatagramHandle, cacheLen int32, remoteAddr string, timeout int32) ([]byte, error) {
	data, rte := RawRtdbRecvDatagramWarp(handle, cacheLen, remoteAddr, timeout)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return data, nil
}

// SubscribeTags 订阅标签点属性更新
func (c *RtdbConnect) SubscribeTags() (chan SubscribeTagsInfo, error) {
	if c.SubscribeTagsConn == nil {
		name := RandString(10)
		c.SubscribeTagsName = name
		con, err := Login(c.HostIp, c.Port, c.UserName, c.Password, RtdbPrecisionNano)
		if err != nil {
			return nil, err
		}
		c.SubscribeTagsConn = con
		ch := make(chan SubscribeTagsInfo, 1024)
		SubscribeTagsChannelLock.Lock()
		defer SubscribeTagsChannelLock.Unlock()
		SubscribeTagsChannelMap[name] = ch
		cName, rte := RawRtdbbSubscribeTagsExWarp(c.SubscribeTagsConn.ConnectHandle, RtdbSubscribeOptionAutoConn, name)
		if !RteIsOk(rte) {
			return nil, rte.GoError()
		}
		c.SubscribeTagsParam = cName
		return ch, nil
	} else {
		return nil, errors.New("已有标签点属性订阅，请不要重复订阅")
	}
}

// CancelSubscribeTags 取消订阅标签点属性更新
func (c *RtdbConnect) CancelSubscribeTags() error {
	if c.SubscribeTagsConn == nil && c.SubscribeTagsName == "" {
		return errors.New("无订阅")
	}
	rte := RawRtdbbCancelSubscribeTagsWarp(c.SubscribeTagsConn.ConnectHandle, c.SubscribeTagsParam)
	if !RteIsOk(rte) {
		return rte.GoError()
	}
	SubscribeTagsChannelLock.Lock()
	defer SubscribeTagsChannelLock.Unlock()
	ch, ok := SubscribeTagsChannelMap[c.SubscribeTagsName]
	if !ok {
		return errors.New("找不到订阅")
	}
	close(ch)
	delete(SubscribeTagsChannelMap, c.SubscribeTagsName)

	c.SubscribeTagsConn = nil
	c.SubscribeTagsName = ""
	c.SubscribeTagsParam = nil
	return nil
}

// SubscribeSnapshots 订阅快照，只要快照发生变化，就会触发订阅
func (c *RtdbConnect) SubscribeSnapshots(infos []*PointInfo) (chan SubscribeSnapshotsInfo, []error, error) {
	if c.SubscribeSnapshotsConn == nil {
		ids := make([]PointID, 0)
		for _, info := range infos {
			ids = append(ids, info.ID)
		}

		name := RandString(10)
		c.SubscribeSnapshotsName = name
		con, err := Login(c.HostIp, c.Port, c.UserName, c.Password, RtdbPrecisionNano)
		if err != nil {
			return nil, nil, err
		}
		c.SubscribeSnapshotsConn = con

		ch := make(chan SubscribeSnapshotsInfo, 1024)
		pointMap := make(map[PointID]*PointInfo)
		for _, info := range infos {
			pointMap[info.ID] = info
		}
		val := SubscribeSnapshotsPointsAndChannel{
			Ch:       ch,
			PointMap: pointMap,
		}
		SubscribeSnapshotsLock.Lock()
		defer SubscribeSnapshotsLock.Unlock()
		SubscribeSnapshotsMap[name] = &val

		cName, rtes, rte := RawRtdbsSubscribeSnapshotsEx64Warp(c.SubscribeSnapshotsConn.ConnectHandle, ids, RtdbSubscribeOptionAutoConn, name)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		c.SubscribeSnapshotsParam = cName
		return ch, RtdbErrorListToErrorList(rtes), nil
	} else {
		return nil, nil, errors.New("订阅快照失败，已存在订阅连接，不可重复订阅")
	}
}

// SubscribeDeltaSnapshots 订阅Delta快照，快照变化需要超过Delta，才会触发订阅，这样可以节约流量
func (c *RtdbConnect) SubscribeDeltaSnapshots(infos []*PointInfo, deltaValues []float64, deltaStates []int64) (chan SubscribeSnapshotsInfo, []error, error) {
	if c.SubscribeSnapshotsConn == nil {
		ids := make([]PointID, 0)
		for _, info := range infos {
			ids = append(ids, info.ID)
		}

		name := RandString(10)
		c.SubscribeSnapshotsName = name
		con, err := Login(c.HostIp, c.Port, c.UserName, c.Password, RtdbPrecisionNano)
		if err != nil {
			return nil, nil, err
		}
		c.SubscribeSnapshotsConn = con

		ch := make(chan SubscribeSnapshotsInfo, 1024)
		pointMap := make(map[PointID]*PointInfo)
		for _, info := range infos {
			pointMap[info.ID] = info
		}
		val := SubscribeSnapshotsPointsAndChannel{
			Ch:       ch,
			PointMap: pointMap,
		}
		SubscribeSnapshotsLock.Lock()
		defer SubscribeSnapshotsLock.Unlock()
		SubscribeSnapshotsMap[name] = &val

		cName, rtes, rte := RawRtdbsSubscribeDeltaSnapshots64Warp(c.SubscribeSnapshotsConn.ConnectHandle, ids, deltaValues, deltaStates, RtdbSubscribeOptionAutoConn, name)
		if !RteIsOk(rte) {
			return nil, nil, rte.GoError()
		}
		c.SubscribeSnapshotsParam = cName
		return ch, RtdbErrorListToErrorList(rtes), nil
	} else {
		return nil, nil, errors.New("订阅快照失败，已存在订阅连接，不可重复订阅")
	}
}

// ChangeSubscribeSnapshots 修改快照订阅设置，新增或删除标签点
func (c *RtdbConnect) ChangeSubscribeSnapshots(infos []*PointInfo, deltaValues []float64, deltaStates []int64, changedTypes []RtdbSubscribeChangeType) ([]error, error) {
	if c.SubscribeSnapshotsConn == nil {
		return nil, errors.New("当前没有快照订阅，无法修改")
	}
	SubscribeSnapshotsLock.Lock()
	defer SubscribeSnapshotsLock.Unlock()
	pointAndCh, ok := SubscribeSnapshotsMap[c.SubscribeSnapshotsName]
	if !ok {
		return nil, errors.New("找不到订阅")
	}
	ids := make([]PointID, 0)
	for i, info := range infos {
		ids = append(ids, info.ID)
		changedTyp := changedTypes[i]
		switch changedTyp {
		case RtdbSubscribeChangeTypeAdd:
			pointAndCh.PointMap[info.ID] = info
		case RtdbSubscribeChangeTypeRemove:
			delete(pointAndCh.PointMap, info.ID)
		}
	}
	rtes, rte := RawRtdbsChangeSubscribeSnapshotsWarp(c.SubscribeSnapshotsConn.ConnectHandle, ids, deltaValues, deltaStates, changedTypes)
	if !RteIsOk(rte) {
		return nil, rte.GoError()
	}
	return RtdbErrorListToErrorList(rtes), nil
}

// CancelSubscribeSnapshots 取消快照订阅
func (c *RtdbConnect) CancelSubscribeSnapshots() error {
	if c.SubscribeSnapshotsConn == nil && c.SubscribeSnapshotsName == "" {
		return errors.New("无订阅")
	}
	rte := RawRtdbsCancelSubscribeSnapshotsWarp(c.SubscribeSnapshotsConn.ConnectHandle, c.SubscribeSnapshotsParam)
	if !RteIsOk(rte) {
		return rte.GoError()
	}

	SubscribeSnapshotsLock.Lock()
	defer SubscribeSnapshotsLock.Unlock()
	pointAndCh, ok := SubscribeSnapshotsMap[c.SubscribeSnapshotsName]
	if !ok {
		return errors.New("找不到订阅")
	}
	close(pointAndCh.Ch)
	delete(SubscribeSnapshotsMap, c.SubscribeSnapshotsName)

	c.SubscribeSnapshotsConn = nil
	c.SubscribeSnapshotsParam = nil
	c.SubscribeSnapshotsName = ""

	return nil
}
