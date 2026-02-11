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

var SubscribeTagsChannelMap sync.Map

type SubscribeTagsInfo struct {
	Name      string
	EventType int32
	Handle    int32
	What      int32
	Ids       []PointID
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
	name := C.GoString(param)
	goIds := ConvertCArrayToSlice(ids, count)
	pointIds := make([]PointID, 0)
	for _, id := range goIds {
		pointIds = append(pointIds, PointID(id))
	}
	info := SubscribeTagsInfo{
		Name:      name,
		EventType: int32(eventType),
		Handle:    int32(handle),
		What:      int32(what),
		Ids:       pointIds,
	}

	val, ok := SubscribeTagsChannelMap.Load(name)
	if !ok {
		return C.rtdb_error(0)
	}
	ch := val.(chan SubscribeTagsInfo)
	// 非阻塞发送，如果满了就丢弃
	select {
	case ch <- info:
	default:
	}

	SubscribeTagsChannelMap.Store(name, ch)

	return C.rtdb_error(0)
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
	return C.rtdb_error(0)
}
