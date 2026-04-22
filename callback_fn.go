package rtdb_api

// #cgo CFLAGS: -DPNG_DEBUG=1 -I./cinclude
// #cgo CXXFLAGS: -std=c++11
// #include <stdlib.h>
// #include "gofn.h"
import "C"
import (
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
