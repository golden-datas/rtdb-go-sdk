package rtdb_api

// #cgo CFLAGS: -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include <stdlib.h>
// #include "gofn.h"
import "C"
import (
	"fmt"
	"net"
	"sync"
	"unsafe"
)

var SubscribeTagsChannelLock sync.Mutex
var SubscribeTagsChannelMap = make(map[string]chan SubscribeTagsInfo)

type SubscribeTagsInfo struct {
	Name      string    // 订阅名称（一个随机字符串，内部使用）
	EventType int32     // 事件类型
	Handle    int32     // 句柄
	What      int32     // 想要做什么
	Ids       []PointID // 更新Point的ID列表
}

//export goSubscribeTagsEx
func goSubscribeTagsEx(
	eventType C.rtdb_uint32,
	handle C.rtdb_int32,
	param unsafe.Pointer,
	count C.rtdb_int32,
	ids *C.rtdb_int32,
	what C.rtdb_int32,
) C.rtdb_error {
	name := C.GoString((*C.char)(param))
	cIds := SafeCopyToSlice[C.rtdb_int32](unsafe.Pointer(ids), int(count))
	pointIds := make([]PointID, 0)
	for _, id := range cIds {
		pointIds = append(pointIds, PointID(id))
	}
	info := SubscribeTagsInfo{
		Name:      name,
		EventType: int32(eventType),
		Handle:    int32(handle),
		What:      int32(what),
		Ids:       pointIds,
	}

	SubscribeTagsChannelLock.Lock()
	defer SubscribeTagsChannelLock.Unlock()
	ch, ok := SubscribeTagsChannelMap[name]
	if !ok {
		return C.rtdb_error(0)
	}
	// 非阻塞发送，如果满了就丢弃
	select {
	case ch <- info:
	default:
	}

	return C.rtdb_error(0)
}

var SubscribeSnapshotsMap = make(map[string]*SubscribeSnapshotsPointsAndChannel)
var SubscribeSnapshotsLock sync.Mutex

type SubscribeSnapshotsPointsAndChannel struct {
	PointMap map[PointID]*PointInfo
	Ch       chan SubscribeSnapshotsInfo
}

type SubscribeSnapshotsInfo struct {
	Name      string
	Handle    int32
	EventType uint32
	PTVQs     []PTVQ
	Errs      []error
}

//export goSnapsEventEx
func goSnapsEventEx(
	eventType C.rtdb_uint32,
	handle C.rtdb_int32,
	param unsafe.Pointer,
	count C.rtdb_int32,
	ids *C.rtdb_int32,
	datetimes *C.rtdb_timestamp_type,
	subtimes *C.rtdb_subtime_type,
	values *C.rtdb_float64,
	status *C.rtdb_int64,
	qualities *C.rtdb_int16,
	errors *C.rtdb_error,
) C.rtdb_error {
	name := C.GoString((*C.char)(param))
	cIds := SafeCopyToSlice[C.rtdb_int32](unsafe.Pointer(ids), int(count))
	pointIds := make([]PointID, 0)
	for _, id := range cIds {
		pointIds = append(pointIds, PointID(id))
	}
	cDatetimes := SafeCopyToSlice[C.rtdb_timestamp_type](unsafe.Pointer(datetimes), int(count))
	cSubtimes := SafeCopyToSlice[C.rtdb_subtime_type](unsafe.Pointer(subtimes), int(count))
	cValues := SafeCopyToSlice[C.rtdb_float64](unsafe.Pointer(values), int(count))
	cStates := SafeCopyToSlice[C.rtdb_int64](unsafe.Pointer(status), int(count))
	cQualities := SafeCopyToSlice[C.rtdb_int16](unsafe.Pointer(qualities), int(count))
	cErrs := SafeCopyToSlice[C.rtdb_error](unsafe.Pointer(errors), int(count))

	SubscribeSnapshotsLock.Lock()
	defer SubscribeSnapshotsLock.Unlock()
	pointsAndChannel, ok := SubscribeSnapshotsMap[name]
	if !ok {
		return C.rtdb_error(0)
	}
	ptvqs := make([]PTVQ, 0)
	for i := 0; i < int(count); i++ {
		ts := RtdbTimestampToGoTime(TimestampType(cDatetimes[i]), SubtimeType(cSubtimes[i]))
		q := Quality(cQualities[i])
		info, ok := pointsAndChannel.PointMap[pointIds[i]]
		if !ok {
			continue
		}
		rtdbType, _ := info.ValueType.ToRawType()
		switch rtdbType {
		case RtdbTypeBool:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqBool(ts, Int64ToBool(int64(cStates[i])), q)))
		case RtdbTypeUint8:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqUint8(ts, uint8(cStates[i]), q)))
		case RtdbTypeInt8:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqInt8(ts, int8(cStates[i]), q)))
		case RtdbTypeChar:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqChar(ts, byte(cStates[i]), q)))
		case RtdbTypeUint16:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqUint16(ts, uint16(cStates[i]), q)))
		case RtdbTypeInt16:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqInt16(ts, int16(cStates[i]), q)))
		case RtdbTypeUint32:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqUint32(ts, uint32(cStates[i]), q)))
		case RtdbTypeInt32:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqInt32(ts, int32(cStates[i]), q)))
		case RtdbTypeInt64:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqInt64(ts, int64(cStates[i]), q)))
		case RtdbTypeReal16:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqFloat16(ts, float32(cValues[i]), q)))
		case RtdbTypeReal32:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqFloat32(ts, float32(cValues[i]), q)))
		case RtdbTypeReal64:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqFloat64(ts, float64(cValues[i]), q)))
		case RtdbTypeFp16:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqFp16(ts, float32(cValues[i]), q)))
		case RtdbTypeFp32:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqFp32(ts, float32(cValues[i]), q)))
		case RtdbTypeFp64:
			ptvqs = append(ptvqs, NewPTVQ(info, NewTvqFp64(ts, float64(cValues[i]), q)))
		default:
			continue
		}
	}

	errs := make([]error, 0)
	for _, err := range cErrs {
		errs = append(errs, RtdbError(err).GoError())
	}

	info := SubscribeSnapshotsInfo{
		Name:      name,
		EventType: uint32(eventType),
		Handle:    int32(handle),
		PTVQs:     ptvqs,
		Errs:      errs,
	}

	select {
	case pointsAndChannel.Ch <- info:
	default:
	}

	return C.rtdb_error(0)
}

var SubscribeConnectEventLock sync.Mutex
var SubscribeConnectEventMap = make(map[string]chan SubscribeConnectEventInfo)

// RtdbConnectEvent 映射 C 结构 RTDB_CONNECT_EVENT
type RtdbConnectEvent struct {
	MsgID           int32
	MsgIdNameString string
	MsgIdDescString string
	BeginS          uint32
	BeginMs         int16
	ApiCategory     ApiCategory
	ClientAddr      uint32
	ClientProcessID int32
	ClientThreadID  int32
	ServerThreadID  int32
	RetVal          RtdbError
	Elapsed         float32
	PreCount        int32
	PostCount       int32
	WriteCount      uint32
	ReadCount       uint32
	WriteTime       float32
	ReadTime        float32
	IndexWriteCount uint32
	IndexReadCount  uint32
	IndexWriteTime  float32
	IndexReadTime   float32
	ArcListLockTime float32
	ArcLockTime     float32
	IndexLockTime   float32
	TotalLockTime   float32
	WriteSize       float32
	ReadSize        float32
	WriteRealSize   float32
	ReadRealSize    float32
	ClientAddr6     string // ipv6地址，16字节二进制转换后的可读字符串
	AddrString      string // 客户端地址可读字符串：优先IPv6，其次IPv4点分十进制
}

// SubscribeConnectEventInfo 回调事件信息
type SubscribeConnectEventInfo struct {
	Name      string
	EventType uint32
	Handle    int32
	Events    []RtdbConnectEvent
	PreCalls  []string
	PostCalls []string
}

//export goConnectEventEx
func goConnectEventEx(
	eventType C.rtdb_uint32,
	handle C.rtdb_int32,
	param unsafe.Pointer,
	count C.rtdb_int32,
	events **C.RTDB_CONNECT_EVENT,
	preCalls **C.char,
	postCalls **C.char,
) C.rtdb_error {
	name := C.GoString((*C.char)(param))

	goEvents := make([]RtdbConnectEvent, 0)
	if events != nil && count > 0 {
		eventPtrs := (*[1 << 30]*C.RTDB_CONNECT_EVENT)(unsafe.Pointer(events))[:count:count]
		for _, cEvent := range eventPtrs {
			if cEvent == nil {
				continue
			}
			event := RtdbConnectEvent{
				MsgID:           int32(cEvent.msg_id),
				BeginS:          uint32(cEvent.begin_s),
				BeginMs:         int16(cEvent.begin_ms),
				ApiCategory:     ApiCategory(cEvent.api_category),
				ClientAddr:      uint32(cEvent.client_addr),
				ClientProcessID: int32(cEvent.client_process_id),
				ClientThreadID:  int32(cEvent.client_thread_id),
				ServerThreadID:  int32(cEvent.server_thread_id),
				RetVal:          RtdbError(cEvent.ret_val),
				Elapsed:         float32(cEvent.elapsed),
				PreCount:        int32(cEvent.pre_count),
				PostCount:       int32(cEvent.post_count),
				WriteCount:      uint32(cEvent.write_count),
				ReadCount:       uint32(cEvent.read_count),
				WriteTime:       float32(cEvent.write_time),
				ReadTime:        float32(cEvent.read_time),
				IndexWriteCount: uint32(cEvent.index_write_count),
				IndexReadCount:  uint32(cEvent.index_read_count),
				IndexWriteTime:  float32(cEvent.index_write_time),
				IndexReadTime:   float32(cEvent.index_read_time),
				ArcListLockTime: float32(cEvent.arc_list_lock_time),
				ArcLockTime:     float32(cEvent.arc_lock_time),
				IndexLockTime:   float32(cEvent.index_lock_time),
				TotalLockTime:   float32(cEvent.total_lock_time),
				WriteSize:       float32(cEvent.write_size),
				ReadSize:        float32(cEvent.read_size),
				WriteRealSize:   float32(cEvent.write_real_size),
				ReadRealSize:    float32(cEvent.read_real_size),
			}
			// 通过 MsgID 获取任务名称和描述
			event.MsgIdNameString, event.MsgIdDescString = RawRtdbJobMessageWarp(int32(cEvent.msg_id))
			// client_addr6 是 16 字节原始二进制（struct in6_addr），需用 net.IP 转为可读字符串
			addr6Bytes := make([]byte, 16)
			src := unsafe.Slice((*byte)(unsafe.Pointer(&cEvent.client_addr6[0])), C.RTDB_IPV6_ADDR_SIZE)
			copy(addr6Bytes, src)
			event.ClientAddr6 = net.IP(addr6Bytes).String()
			// AddrString：IPv6非全零时优先显示IPv6，否则显示IPv4点分十进制
			ipv6Valid := event.ClientAddr6 != "" && event.ClientAddr6 != "::"
			if ipv6Valid {
				event.AddrString = event.ClientAddr6
			} else {
				event.AddrString = fmt.Sprintf("%d.%d.%d.%d",
					(event.ClientAddr>>24)&0xFF,
					(event.ClientAddr>>16)&0xFF,
					(event.ClientAddr>>8)&0xFF,
					event.ClientAddr&0xFF)
			}
			goEvents = append(goEvents, event)
		}
	}

	goPreCalls := make([]string, 0)
	if preCalls != nil && count > 0 {
		preCallPtrs := (*[1 << 30]*C.char)(unsafe.Pointer(preCalls))[:count:count]
		for _, cStr := range preCallPtrs {
			if cStr != nil {
				goPreCalls = append(goPreCalls, C.GoString(cStr))
			}
		}
	}

	goPostCalls := make([]string, 0)
	if postCalls != nil && count > 0 {
		postCallPtrs := (*[1 << 30]*C.char)(unsafe.Pointer(postCalls))[:count:count]
		for _, cStr := range postCallPtrs {
			if cStr != nil {
				goPostCalls = append(goPostCalls, C.GoString(cStr))
			}
		}
	}

	info := SubscribeConnectEventInfo{
		Name:      name,
		EventType: uint32(eventType),
		Handle:    int32(handle),
		Events:    goEvents,
		PreCalls:  goPreCalls,
		PostCalls: goPostCalls,
	}

	SubscribeConnectEventLock.Lock()
	defer SubscribeConnectEventLock.Unlock()
	ch, ok := SubscribeConnectEventMap[name]
	if !ok {
		return C.rtdb_error(0)
	}

	select {
	case ch <- info:
	default:
	}

	return C.rtdb_error(0)
}
