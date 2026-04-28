/*
    C的回调函数声明，这些函数的实现都是用Go写的，位于callback_fn.go中
*/

#ifndef _GOFN_H_
#define _GOFN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "rtdb.h"
#include "rtdbapi.h"
#include "rtdb_error.h"

// rtdbb_subscribe_tags_ex_warp 中的回调函数，用于订阅 标签点属性变更
extern rtdb_error goSubscribeTagsEx(
    rtdb_uint32 event_type,
    rtdb_int32 handle,
    void* param,
    rtdb_int32 count,
    rtdb_int32 *ids,
    rtdb_int32 what);

// rtdbs_subscribe_snapshots_ex64_warp 中的回调函数，用于订阅 快照值变更
extern rtdb_error goSnapsEventEx(
    rtdb_uint32 event_type,
    rtdb_int32 handle,
    void* param,
    rtdb_int32 count,
    rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities,
    rtdb_error* errors);

// rtdb_subscribe_connect_ex_warp 中的回调函数，用于订阅 API调用连接事件
extern rtdb_error goConnectEventEx(
    rtdb_uint32 event_type,
    rtdb_int32 handle,
    void* param,
    rtdb_int32 count,
    RTDB_CONNECT_EVENT** events,
    char** pre_calls,
    char** post_calls);

#ifdef __cplusplus
}
#endif

#endif // _C_API_H_
