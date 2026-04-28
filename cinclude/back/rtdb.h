#ifndef __RTDB_H__
#define __RTDB_H__

/*
#if defined(WIN32) || defined(_WIN32) || defined(_WIN64)
#ifdef RTDBAPI_EXPORTS
#  define RTDBAPI extern "C" __declspec(dllexport)
#else
#  define RTDBAPI extern "C" __declspec(dllimport)
#endif

#define RTDBAPI_CALLRULE _stdcall
#else
#define RTDBAPI extern "C" __attribute__ ((visibility ("default")))
#define RTDBAPI_CALLRULE
#endif
*/

#ifdef __cplusplus
#define RTDBAPI_EXTERN extern "C"
#else
#define RTDBAPI_EXTERN extern
#endif

#if defined(_WIN32) || defined(WIN32)
  #ifdef RTDBAPI_EXPORTS
    #define RTDBAPI RTDBAPI_EXTERN __declspec(dllexport)
  #else
    #define RTDBAPI RTDBAPI_EXTERN __declspec(dllimport)
  #endif
  #define RTDBAPI_CALLRULE _stdcall
#else
  #define RTDBAPI RTDBAPI_EXTERN __attribute__ ((visibility ("default")))
  #define RTDBAPI_CALLRULE
#endif

#ifdef __cplusplus
#define GAPI_DEFAULT_VALUE(v) = v
#else
#define GAPI_DEFAULT_VALUE(v)
#endif

#ifndef RTDBAPI_DISABLE_DEPRECATED
#ifdef __GNUC__
#define RTDBAPI_DEPRECATED(msg) __attribute__((__deprecated__("Since " #msg)))
#define RTDBAPI_DEPRECATED_FOR(msg, replacement) __attribute__((__deprecated__("Since " #msg "; use " #replacement)))
#elif defined(_MSC_VER)
#define RTDBAPI_DEPRECATED(msg) __declspec(deprecated("Since " # msg))
#define RTDBAPI_DEPRECATED_FOR(msg, replacement) __declspec(deprecated("Since " #msg "; use " #replacement))
#else
#define RTDBAPI_DEPRECATED(msg)
#define RTDBAPI_DEPRECATED_FOR(msg, replacement)
#endif
#else
#define RTDBAPI_DEPRECATED()
#define RTDBAPI_DEPRECATED_FOR(msg, replacement)
#endif


#ifndef For_Classof
#define Rtdb_Mark_Classof        0x0F
#define Rtdb_Len_Classof         4
#define Rtdb_Base_Datetime_Type  0x10
#endif

/*! \defgroup ddatatype 数据类型 */
/*! \defgroup denum 枚举 */
/*! \defgroup dstruct 结构体 */
/*! \defgroup dmacro 宏定义 */
/*! \defgroup dcallback 回调函数定义 */

/**
* \ingroup ddatatype
* \typedef struct _RTDB_COORDINATE RTDB_COORDINATE
* \brief 坐标类型
* \see _RTDB_COORDINATE
*/

/**
* \ingroup ddatatype
* \brief   坐标类型
*/
typedef struct _RTDB_COORDINATE
{
  float x;  //!< x 坐标
  float y;  //!< y 坐标
} RTDB_COORDINATE;

/**
* \ingroup ddatatype
* \typedef unsigned char rtdb_byte
* \brief   单字节数值
*/
typedef unsigned char       rtdb_byte;
/**
* \typedef char rtdb_int8
* \ingroup ddatatype
* \brief 8位整数
*/
typedef signed char         rtdb_int8;
/**
* \typedef unsigned char rtdb_uint8
* \ingroup ddatatype
* \brief 8位正整数
*/
typedef unsigned char       rtdb_uint8;
/**
* \typedef short rtdb_int16
* \ingroup ddatatype
* \brief 16位整数
*/
typedef short               rtdb_int16;
/**
* \typedef unsigned short rtdb_uint16
* \ingroup ddatatype
* \brief 16位正整数
*/
typedef unsigned short      rtdb_uint16;
/**
* \typedef int rtdb_int32
* \ingroup ddatatype
* \brief 32位整数
*/
typedef int                 rtdb_int32;
/**
* \typedef unsigned int rtdb_uint32
* \ingroup ddatatype
* \brief 32位正整数
*/
typedef unsigned int        rtdb_uint32;
#ifdef WIN32
/**
* \typedef long long rtdb_int64
* \ingroup ddatatype
* \brief 64位整数
*/
typedef long long           rtdb_int64;
/**
* \typedef unsigned long long rtdb_uint64
* \ingroup ddatatype
* \brief 64位正整数
*/
typedef unsigned long long  rtdb_uint64;
#else
/**
* \typedef long int rtdb_int64
* \ingroup ddatatype
* \brief 64位整数
*/
typedef long int            rtdb_int64;
/**
* \typedef unsigned long int rtdb_uint64
* \ingroup ddatatype
* \brief 64位正整数
*/
typedef unsigned long int   rtdb_uint64;
#endif // WIN32
/**
* \typedef short rtdb_float16
* \ingroup ddatatype
* \brief 16位浮点数
*/
typedef short               rtdb_float16;
/**
* \typedef float rtdb_float32
* \ingroup ddatatype
* \brief 32位浮点数
*/
typedef float               rtdb_float32;
/**
* \typedef double rtdb_float64
* \ingroup ddatatype
* \brief 64位浮点数
*/
typedef double              rtdb_float64;
/**
* \typedef RTDB_COORDINATE rtdb_coordinate
* \ingroup ddatatype
* \brief 二维坐标
*/
typedef RTDB_COORDINATE   rtdb_coordinate;
/**
* \typedef unsigned int rtdb_error
* \ingroup ddatatype
* \brief 错误数值
*/
typedef unsigned int        rtdb_error;
/**
* \typedef void* rtdb_datagram_handle
* \ingroup ddatatype
* \brief UDP连接句柄
*/
typedef void*               rtdb_datagram_handle;


typedef rtdb_int32		rtdb_time_type;	        //!< ms字段类型
#define RTDB_MS_PRECISION       1000		    //!< 时间转换精度：毫秒
#define RTDB_MS_MAX		        999		        //!< ms字段范围最大值

#define RTDB_MICRO_PRECISION    1000000	        //!< 时间转换精度：微秒
#define RTDB_MICRO_MAX		    999999		    //!< ms字段范围最大值

#define RTDB_NANO_PRECISION     1000000000	    //!< 时间转换精度：纳秒
#define RTDB_NANO_MAX		    999999999	    //!< ms字段范围最大值


typedef rtdb_int32 rtdb_length_type;            //!< string、blob、自定义数据的长度类型

typedef rtdb_uint32 rtdb_datetime_type;         //32位时间戳类型，秒级时间戳
typedef rtdb_int64 rtdb_timestamp_type;         //64位时间戳类型，秒级时间戳
typedef rtdb_int32 rtdb_subtime_type;           //时间戳，小于秒的部分，根据设置的全局时间戳精度，表示毫秒、微秒、纳秒的部分
typedef rtdb_int8 rtdb_precision_type;          //时间戳精度类型，0秒，1毫秒，2微秒，3纳秒


/**
* \ingroup denum
* \brief 0x0200 以内的系统定义质量码.
*/
typedef enum _RTDB_QUALITY
{
  RTDB_Q_GOOD = 0,              //!< 正常
  RTDB_Q_NODATA = 1,            //!< 无数据
  RTDB_Q_CREATED = 2,           //!< 创建
  RTDB_Q_SHUTDOWN = 3,          //!< 停机
  RTDB_Q_CALCOFF = 4,           //!< 计算停止
  RTDB_Q_BAD = 5,               //!< 坏点
  RTDB_Q_DIVBYZERO = 6,         //!< 被零除
  RTDB_Q_REMOVED = 7,           //!< 已被删除
  RTDB_Q_OPC = 256,             //!< 从0x0100至0x01FF为OPC质量码
  RTDB_Q_USER = 512             //!< 此质量码（含）之后为用户自定义
} RTDB_QUALITY;

/**
* \ingroup denum
* \brief 订阅选项
*/
typedef enum _RTDB_OPTION
{
  RTDB_O_AUTOCONN = 1,           //!< 自动重连
} RTDB_OPTION;

/**
* \ingroup denum
* \brief 订阅事件
*/
typedef enum _RTDB_EVENT_TYPE
{
  RTDB_E_DATA       = 0,        //!< 数据
  RTDB_E_DISCONNECT = 1,        //!< 连接断开
  RTDB_E_RECOVERY   = 2,        //!< 连接恢复
  RTDB_E_SWITCHING  = 3,        //!< 双活模式，快照订阅，开始切换连接
  RTDB_E_SWITCHED   = 4,        //!< 双活模式，快照订阅，切换连接完毕
  RTDB_E_CHANGED    = 5,        //!< 订阅信息发生变化
} RTDB_EVENT_TYPE;

/**
* \ingroup dmacro
* \def RTDB_QUALITY_MODIFY_UPDATE_FLAG(Q, U)
* \brief 为质量码增加或去除 UPDATE 标志，\b U 为 true 增加 UPDATE 标志，否则去除 UPDATE 标志.
* \param Q 质量码
* \param U 增加或者去除 UPDATE标志
*/
#define RTDB_QUALITY_MODIFY_UPDATE_FLAG(Q, U)   (U) ? (Q) |= 0x8000 : (Q) &= 0x7FFF;

/**
* \ingroup dmacro
* \def RTDB_QUALITY_IS_UPDATED(Q)
* \brief 根据质量码判定事件是否被修改过.
* \param Q 质量码
*/
#define RTDB_QUALITY_IS_UPDATED(Q)              ( (Q) & 0x8000 )

/**
* \ingroup dmacro
* \def RTDB_QUALITY_WITHOUT_FLAG(Q)
* \brief 获得不包含 UPDATE 标志的质量码.
* \param Q 质量码
*/
#define RTDB_QUALITY_WITHOUT_FLAG(Q)            ( (Q) & 0x7FFF )

/**
* \ingroup dmacro
* \def RTDB_QUALITY_OPC_GOOD
* \brief OPC协议中正常有效数据的质量.
*/
#define RTDB_QUALITY_OPC_GOOD                   0xC0

/**
* \ingroup dmacro
* \def RTDB_QUALITY_FROM_OPC(OPC_Q)
* \brief 从OPC协议获取的质量码，使用此宏转换为RTDB中对应的质量码.
* \param OPC_Q 从OPC协议获取的质量码
*/
#define RTDB_QUALITY_FROM_OPC(OPC_Q)            ( RTDB_QUALITY_OPC_GOOD == (OPC_Q) ? RTDB_Q_GOOD : (OPC_Q) & 0xFF | RTDB_Q_OPC )

/**
* \ingroup dmacro
* \def RTDB_QUALITY_IS_OPC(Q)
* \brief 判定RTDB中的质量码是否属于OPC协议范围.
* \param Q RTDB中的质量码
*/
#define RTDB_QUALITY_IS_OPC(Q)                  ( RTDB_QUALITY_WITHOUT_FLAG(Q) >= RTDB_Q_OPC && RTDB_QUALITY_WITHOUT_FLAG(Q) < RTDB_Q_USER )

/**
* \ingroup dmacro
* \def RTDB_QUALITY_TO_OPC(Q)
* \brief RTDB中属于OPC协议范围的质量码，使用此宏还原为OPC原始质量码.
* \param Q RTDB中属于OPC协议范围的质量码
*/
#define RTDB_QUALITY_TO_OPC(Q)                  ( RTDB_Q_GOOD == RTDB_QUALITY_WITHOUT_FLAG(Q) ? RTDB_QUALITY_OPC_GOOD : (Q) & 0xFF )

/**
* \ingroup dmacro
* \def RTDB_QUALITY_IS_VALID(Q)
* \brief 依据质量判定对应的事件是否正常有效.
* \param Q 质量码
*/
#define RTDB_QUALITY_IS_VALID(Q)                ( RTDB_Q_GOOD == RTDB_QUALITY_WITHOUT_FLAG(Q) || RTDB_Q_REMOVED < RTDB_QUALITY_WITHOUT_FLAG(Q))

/**
* \ingroup dmacro
* \def RTDB_QUALITY_IS_REMOVED(Q)
* \brief 依据质量判定对应的事件是否被删除.
* \param Q 质量码
*/
#define RTDB_QUALITY_IS_REMOVED(Q)              ( RTDB_Q_REMOVED == RTDB_QUALITY_WITHOUT_FLAG(Q) )

/**
* \ingroup dmacro
* \def RTDB_COMPRESSED_BLOCK_HEAD
* \brief 已压缩数据块固定块头 40字节
*/
#define RTDB_COMPRESSED_BLOCK_HEAD 40

/**
* \ingroup denum
* \brief 系统常数定义.
*/
typedef enum _RTDB_CONST
{
  RTDB_TAG_SIZE = 80,                                 //!< 标签点名称占用字节数。
  RTDB_DESC_SIZE = 100,                               //!< 标签点描述占用字节数。
  RTDB_UNIT_SIZE = 20,                                //!< 标签点单位占用字节数。
  RTDB_USER_SIZE = 20,                                //!< 用户名占用字节数。
  RTDB_SOURCE_SIZE = 256,                             //!< 标签点数据源占用字节数。
  RTDB_INSTRUMENT_SIZE = 50,                          //!< 标签点所属设备占用字节数。
  RTDB_LOCATIONS_SIZE = 5,                            //!< 采集标签点位址个数。
  RTDB_USERINT_SIZE = 2,                              //!< 采集标签点用户自定义整数个数。
  RTDB_USERREAL_SIZE = 2,                             //!< 采集标签点用户自定义浮点数个数。
  RTDB_EQUATION_SIZE = 2036,                          //!< 计算标签点方程式占用字节数。
  RTDB_TYPE_NAME_SIZE = 21,                           //!< 自定义类型名称占用字节数。
  RTDB_PACK_OF_SNAPSHOT = 0,                          //!< 事件快照备用字节空间。
  RTDB_PACK_OF_POINT = 4,                             //!< 标签点备用字节空间。
  RTDB_PACK_OF_BASE_POINT = 74,                       //!< 基本标签点备用字节空间。
  RTDB_PACK_OF_SCAN = 164,                            //!< 采集标签点备用字节空间。
  RTDB_PACK_OF_CALC = 0,                              //!< 计算标签点备用字节空间。
  RTDB_FILE_NAME_SIZE = 64,                           //!< 文件名字符串字节长度
  RTDB_PATH_SIZE = 1024 - 4 - RTDB_FILE_NAME_SIZE,    //!< 路径字符串字节长度
  RTDB_MAX_USER_COUNT = 100,                          //!< 最大用户个数
  RTDB_MAX_AUTH_COUNT = 100,                          //!< 最大信任连接段个数
  RTDB_MAX_BLACKLIST_LEN = 100,                       //!< 连接黑名单最大长度
  RTDB_MAX_SUBSCRIBE_SNAPSHOTS = 1000,                //!< 单个连接最大订阅快照数量
  RTDB_API_SERVER_DESCRIPTION_LEN = 512,              //!< API_SERVER中，decription描述字段的长度

  RTDB_MIN_EQUATION_SIZE = 480,                       //!< 缩减长度后的方程式占用字节数
  RTDB_PACK_OF_MIN_CALC = 19,                         //!< 缩减长度后的计算标签点备用字节空间
  RTDB_MAX_EQUATION_SIZE = 62 * 1024,                 //!< 最大长度的方程式占用字节数
  RTDB_PACK_OF_MAX_CALC = 64 * 1024 - RTDB_MAX_EQUATION_SIZE - 12,  //!< 最大长度的计算标签点备用字节空间
  RTDB_MAX_JSON_SIZE = 16 * 1024,                     //!< 允许的json字符串的最大长度
  RTDB_IPV6_ADDR_SIZE = 16,                           //!< ipv6地址空间大小

  RTDB_OUTPUT_PLUGIN_LIB_VERSION_LENGTH = 0x40,       /// 输入输出适配器插件库版本号长度 64  Bytes
  RTDB_OUTPUT_PLUGIN_LIB_NAME_LENGTH = 0x80,          /// 输入输出适配器插件库名长度    128 Bytes
  RTDB_OUTPUT_PLUGIN_ACTOR_NAME_LENGTH = 0x80,        /// 输入输出适配器执行实例长度    128 Bytes
  RTDB_OUTPUT_PLUGIN_NAME_LENGTH = 0xFF,              /// 输入输出适配器插件名长度      255 Bytes
  RTDB_OUTPUT_PLUGIN_DIR_LENGTH = 0xFF,               /// 输入输出适配器路径长度        255 Bytes
  RTDB_PER_OF_USEFUL_MEM_SIZE = 10,					  /// 历史数据缓存/补历史数据缓存/blob补历史数据缓存/内存大小上限占可用内存百分比值占用的字节数
} RTDB_CONST;


typedef char rtdb_tag_string[RTDB_TAG_SIZE];
typedef char rtdb_desc_string[RTDB_DESC_SIZE];
typedef char rtdb_unit_string[RTDB_USER_SIZE];
typedef char rtdb_source_string[RTDB_SOURCE_SIZE];
typedef char rtdb_instrument_string[RTDB_INSTRUMENT_SIZE];
typedef char rtdb_equation_string[RTDB_EQUATION_SIZE];
typedef char rtdb_type_name_string[RTDB_TYPE_NAME_SIZE];
typedef char rtdb_filename_string[RTDB_FILE_NAME_SIZE];
typedef char rtdb_path_string[RTDB_PATH_SIZE];
typedef char rtdb_min_equation_string[RTDB_MIN_EQUATION_SIZE];
typedef char rtdb_json_string[RTDB_MAX_JSON_SIZE];

/**
* \ingroup denum
* \brief 标签点数值类型，决定了标签点数值所占用的存储字节数。
*/
typedef enum _RTDB_TYPE
{
  RTDB_BOOL = 0,        //!< 布尔类型，0值或1值。
  RTDB_UINT8 = 1,       //!< 无符号8位整数，占用1字节。
  RTDB_INT8 = 2,        //!< 有符号8位整数，占用1字节。
  RTDB_CHAR = 3,        //!< 单字节字符，占用1字节。
  RTDB_UINT16 = 4,      //!< 无符号16位整数，占用2字节。
  RTDB_INT16 = 5,       //!< 有符号16位整数，占用2字节。
  RTDB_UINT32 = 6,      //!< 无符号32位整数，占用4字节。
  RTDB_INT32 = 7,       //!< 有符号32位整数，占用4字节。
  RTDB_INT64 = 8,       //!< 有符号64位整数，占用8字节。
  RTDB_REAL16 = 9,      //!< 16位浮点数，占用2字节。
  RTDB_REAL32 = 10,     //!< 32位单精度浮点数，占用4字节。
  RTDB_REAL64 = 11,     //!< 64位双精度浮点数，占用8字节。
  RTDB_COOR = 12,       //!< 二维坐标，具有x、y两个维度的浮点数，占用8字节。
  RTDB_STRING = 13,     //!< 字符串，长度不超过存储页面大小。
  RTDB_BLOB = 14,       //!< 二进制数据块，占用字节不超过存储页面大小。
  RTDB_NAMED_T = 15,    //!< 自定义类型，由用户创建时确定字节长度。
  RTDB_DATETIME = 16,   //!< 时间格式类型
  RTDB_FP16 = 17,       //!< 定点数
  RTDB_FP32 = 18,       //!< 定点数
  RTDB_FP64 = 19,       //!< 定点数
} RTDB_TYPE;

#define RTDB_TYPE_COUNT (RTDB_FP64 + 1)

/**
* \ingroup dmacro
* \def RTDB_TYPE_IS_NORMAL(TYPE)
* \brief 判断此数据类型是否是除 \b RTDB_STRING 、\b RTDB_BLOB、\b RTDB_NAMED_T 之外的数据类型
* \param TYPE  数据类型
* \see RTDB_TYPE
*/
#define RTDB_TYPE_IS_NORMAL(TYPE)    ( (TYPE == RTDB_STRING || TYPE == RTDB_BLOB || TYPE == RTDB_NAMED_T) ? false : true )

/**
* \ingroup denum
* \brief 标签点类别，决定了标签点具有哪些扩展属性。标签点可以同时具备多个类别，最多可以定义33个标签点类别。.
*/
typedef enum _RTDB_CLASS
{
  RTDB_BASE = 0,    //!< 基本标签点，所有类别标签点均在基本标签点的属性集上扩展自己的属性集。
  RTDB_SCAN = 1,    //!< 采集标签点。
  RTDB_CALC = 2,    //!< 计算标签点。
  RTDB_ALARM = 4,   //!< 报警标签点。
} RTDB_CLASS;

/**
* \ingroup dmacro
* \def RTDB_CLASS_IS_SCAN(CLASSOF)
* \brief 判定标签点是否采集标签点.
* \param CLASSOF 标签点类别
* \see RTDB_CLASS
*/
#define RTDB_CLASS_IS_SCAN(CLASSOF)    (CLASSOF & RTDB_SCAN)

/**
* \ingroup dmacro
* \def RTDB_CLASS_IS_CALC(CLASSOF)
* \brief 判定标签点是否计算标签点.
* \param CLASSOF 标签点类别
* \see RTDB_CLASS
*/
#define RTDB_CLASS_IS_CALC(CLASSOF)    (CLASSOF & RTDB_CALC)

/**
* \ingroup dmacro
* \def RTDB_CLASS_IS_ALARM(CLASSOF)
* \brief 判定标签点是否报警标签点.
* \param CLASSOF 标签点类别
* \see RTDB_CLASS
*/
#define RTDB_CLASS_IS_ALARM(CLASSOF)    (CLASSOF & RTDB_ALARM)

/**
* \ingroup denum
* \brief 计算标签点触发机制.
*/
typedef enum _RTDB_TRIGGER
{
  RTDB_NULL_TRIGGER,        //!< 无触发
  RTDB_EVENT_TRIGGER,       //!< 事件触发
  RTDB_TIMER_TRIGGER,       //!< 周期触发
  RTDB_FIXTIME_TRIGGER,     //!< 定时触发
} RTDB_TRIGGER;

/**
* \ingroup denum
* \brief 计算结果时间戳参考.
*/
typedef enum _RTDB_TIME_COPY
{
  RTDB_CALC_TIME,             //!< 采用计算时间
  RTDB_LATEST_TIME,           //!< 采用最晚标签点时间
  RTDB_EARLIEST_TIME,         //!< 采用最早标签点时间
} RTDB_TIME_COPY;

/**
* \ingroup denum
* \brief 标签点搜索结果排序方式.
*/
typedef enum _RTDB_SEARCH_SORT
{
  RTDB_SORT_BY_TABLE,     //!< 首先按所属表排序，同一个表内的标签点之间按标签点名称排序
  RTDB_SORT_BY_TAG,       //!< 以标签点名称排序
  RTDB_SORT_BY_ID,        //!< 以标签点ID排序
} RTDB_SEARCH_SORT;

/**
* \ingroup denum
* \brief 历史数据搜索方式.
*/
typedef enum _RTDB_HIS_MODE
{
  RTDB_NEXT,            //!< 寻找下一个最近的数据；
  RTDB_PREVIOUS,        //!< 寻找上一个最近的数据；
  RTDB_EXACT,           //!< 取指定时间的数据，如果没有则返回错误 \b RtE_DATA_NOT_FOUND；
  RTDB_INTER,           //!< 取指定时间的内插值数据;
  RTDB_EXACT_OR_NEXT,   //!< 取指定时间的数据，如果没有则取下一条数据。如果都没有数据则返回错误 \b RtE_DATA_NOT_FOUND;
  RTDB_EXACT_OR_PREV,   //!< 取指定时间的数据，如果没有则取上一条数据。如果都没有数据则返回错误 \b RtE_DATA_NOT_FOUND;
  RTDB_INTER_OR_NEXT,   //!< 取指定时间的内插值数据, 如果没有则取下一条数据。如果都没有数据则返回错误 \b RtE_DATA_NOT_FOUND;
} RTDB_HIS_MODE;

/**
* \ingroup denum
* \brief 搜索标签点所指定的属性集合
*/
typedef enum _RTDB_SEARCH_MASK
{
  RTDB_SEARCH_NULL,                               //!< 不使用任何标签点属性作为搜索条件
  RTDB_SEARCH_COMPDEV,                            //!< 使用压缩偏差作为搜索条件
  RTDB_SEARCH_COMPMAX,                            //!< 最大压缩间隔
  RTDB_SEARCH_COMPMIN,                            //!< 最小压缩间隔
  RTDB_SEARCH_EXCDEV,                             //!< 例外偏差
  RTDB_SEARCH_EXCMAX,                             //!< 最大例外间隔
  RTDB_SEARCH_EXCMIN,                             //!< 最小例外间隔
  RTDB_SEARCH_SUMMARY,                            //!< 是否加速
  RTDB_SEARCH_MIRROR,                             //!< 是否镜像
  RTDB_SEARCH_COMPRESS,                           //!< 是否压缩
  RTDB_SEARCH_STEP,                               //!< 是否阶跃
  RTDB_SEARCH_SHUTDOWN,                           //!< 是否停机补写
  RTDB_SEARCH_ARCHIVE,                            //!< 是否归档
  RTDB_SEARCH_CHANGER,                            //!< 修改人
  RTDB_SEARCH_CREATOR,                            //!< 创建人
  RTDB_SEARCH_LOWLIMIT,                           //!< 量程下限
  RTDB_SEARCH_HIGHLIMIT,                          //!< 量程上限
  RTDB_SEARCH_TYPICAL,                            //!< 典型值
  RTDB_SEARCH_CHANGEDATE,                         //!< 修改时间
  RTDB_SEARCH_CREATEDATE,                         //!< 创建时间
  RTDB_SEARCH_DIGITS,                             //!< 数值位数
  RTDB_SEARCH_COMPDEVPERCENT,                     //!< 压缩偏差百分比
  RTDB_SEARCH_EXCDEVPERCENT,                      //!< 例外偏差百分比
  RTDB_SEARCH_SCAN_BEGIN,                         //!< 辅助作用，不能作为搜索条件
  RTDB_SEARCH_SCAN,                               //!< 是否采集
  RTDB_SEARCH_LOCATION1,                          //!< 设备位址1
  RTDB_SEARCH_LOCATION2,                          //!< 设备位址2
  RTDB_SEARCH_LOCATION3,                          //!< 设备位址3
  RTDB_SEARCH_LOCATION4,                          //!< 设备位址4
  RTDB_SEARCH_LOCATION5,                          //!< 设备位址5
  RTDB_SEARCH_USERINT1,                           //!< 自定义整数1
  RTDB_SEARCH_USERINT2,                           //!< 自定义整数2
  RTDB_SEARCH_USERREAL1,                          //!< 自定义单精度浮点数1
  RTDB_SEARCH_USERREAL2,                          //!< 自定义单精度浮点数2
  RTDB_SEARCH_SCAN_END,                           //!< 辅助作用，不能作为搜索条件
  RTDB_SEARCH_CALC_BEGIN,                         //!< 辅助作用，不能作为搜索条件
  RTDB_SEARCH_EQUATION,                           //!< 方程式
  RTDB_SEARCH_TRIGGER,                            //!< 计算触发机制
  RTDB_SEARCH_TIMECOPY,                           //!< 计算结果时间戳参考
  RTDB_SEARCH_PERIOD,                             //!< 计算周期
  RTDB_SEARCH_CALC_END,                           //!< 辅助作用，不能作为搜索条件
} RTDB_SEARCH_MASK;

/**
* \ingroup denum
* \brief 用于设置API的工作模式的参数选项
* \see rtdb_set_option
*/
typedef enum _RTDB_API_OPTION
{
  RTDB_API_AUTO_RECONN,       //!< api 在连接中断后是否自动重连, 0 不重连；1 重连。默认为 0 不重连
  RTDB_API_CONN_TIMEOUT,      //!< api 连接超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
  RTDB_API_SEND_TIMEOUT,      //!< api 发送超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为1000
  RTDB_API_RECV_TIMEOUT,      //!< api 接收超时值设置（单位：毫秒）,0 阻塞模式，无限等待，默认为60000
  RTDB_API_USER_TIMEOUT,	  //!< api TCP_USER_TIMEOUT超时值设置（单位：毫秒），默认为10000，Linux内核2.6.37以上有效
  RTDB_API_DEFAULT_PRECISION, //!< api 默认的时间戳精度，当使用旧版相关的api，以及新版api中未设置时间戳精度时，则使用此默认时间戳精度。 默认为毫秒精度
  RTDB_API_SERVER_PRECISION,  //!< api 连接3.0数据库时，设置3.0数据库的时间戳精度，0表示毫秒精度，非0表示纳秒精度，默认为毫秒精度
} RTDB_API_OPTION;

/**
* \ingroup denum
* \brief 数据库对应的服务
*/
typedef enum _RTDB_PROCESS_NAME
{
  RTDB_PROCESS_FIRST = 1,                         //!< 计数作用
  RTDB_PROCESS_HISTORIAN = RTDB_PROCESS_FIRST,    //!< 历史服务
  RTDB_PROCESS_EQUATION,                          //!< 方程式服务
  RTDB_PROCESS_BASE,                              //!< 标签点服务
  RTDB_PROCESS_LAST,                              //!< 计数作用
} RTDB_PROCESS_NAME;

/**
* \ingroup denum
* \brief 数据库对应服务的大任务，每个服务最多同时执行一个大任务
*/
typedef enum _RTDB_BIG_JOB_NAME
{
  /// 历史数据服务
  RTDB_MERGE = 1,               //!< 合并附属文件到主文件
  RTDB_ARRANGE = 2,             //!< 整理存档文件，整理过程中会完成合并
  RTDB_REINDEX = 3,             //!< 重建索引
  RTDB_BACKUP = 4,              //!< 备份
  RTDB_REACTIVE = 5,            //!< 激活为活动存档
  RTDB_MOVE_ARCHIVE = 6,        //!< 移动存档文件
  RTDB_CONVERT_INDEX = 7,       //!< 转换存档文件索引类型
  RTDB_COMPRESS = 8,            //!< 压缩
  RTDB_MOVE_ARCHIVE_AUTO = 9,   //!< 自动移动存档文件（用于存档文件滚动存储）
  /// 方程式服务
  RTDB_COMPUTE = 11,            //!< 历史计算
  /// 标签点信息服务
  RTDB_UPDATE_TABLE = 21,       //!< 修改表名称
  RTDB_REMOVE_TABLE = 22,       //!< 删除表
} RTDB_BIG_JOB_NAME;

/**
* \ingroup ddatatype
* \typedef enum RTDB_POINT_MIRROR
* \brief 标签点镜像属性
* \see _RTDB_TABLE
*/
typedef enum _RTDB_POINT_MIRROR
{
	RTDB_POINT_OFF = 0,		    //!<镜像关闭
	RTDB_POINT_SEND_RECV = 1,   //!<镜像收发
	RTDB_POINT_RECV = 2,		//!<镜像接收
	RTDB_POINT_SEND = 3		    //!<镜像发送
} RTDB_POINT_MIRROR;

/**
* \ingroup ddatatype
* \typedef struct _RTDB_TABLE RTDB_TABLE
* \brief 标签点索引表属性集
* \see _RTDB_TABLE
*/

/**
* \ingroup dstruct
* \brief 标签点索引表属性集。.
*/
typedef struct _RTDB_TABLE
{
  /**
  * 表的唯一标识。
  * 从1开始，到上限为止。
  */
  int  id;
  /**
  * 表类型。
  * 暂时保留用途。
  */
  int  type;
  /** 表名称。
  *  命名规则：
  *  1、第一个字符必须是26个字母之一或数字0-9之一；
  *  2、不允许使用控制字符，比如换行符或制表符；
  *  3、不允许使用以下字符（'*'、'''、'?'、';'、'{'、'}'、'['、']'、'|'、'\'、'`'、'''、'"'、'.'）；
  *  4、字节长度不要超出 \b RTDB_TAG_SIZE，如果超出，系统会自动将后面的字符截断。
  */
  char name[RTDB_TAG_SIZE];
  /**
  *  表描述。
  *  缺省值：空字符串。
  *  字节长度不要超出 \b RTDB_DESC_SIZE，多余的部分会被截断。
  */
  char desc[RTDB_DESC_SIZE];
} RTDB_TABLE;

/**
* \ingroup ddatatype
* \typedef struct _RTDB_POINT RTDB_POINT
* \brief 基本标签点属性集
* \see _RTDB_POINT
*/
/**
* \ingroup dstruct
* \brief 基本标签点属性集。8字节对齐的条件下占用512字节.
*/
typedef struct _RTDB_POINT
{
  /**
  *  标签点名称。
  *  用于在表中唯一标识一个标签点；
  *  该属性允许被修改；
  *  命名规则：
  *  1、第一个字符必须是26个字母之一或数字0-9之一或者"_"、"%"；
  *  2、不允许使用控制字符，比如换行符或制表符；
  *  3、不允许使用以下字符（'*'、'''、'?'、';'、'('、')'、'{'、'}'、'['、']'、'|'、'\'、'`'、'''、'"'、'.'）；
  *  4、字节长度不要超出 \b RTDB_TAG_SIZE，如果超出，系统会自动将后面的字符截断。
  */
  char tag[RTDB_TAG_SIZE];
  /**
  *  全库唯一标识。
  *  只读属性，创建标签点时系统会自动为标签点分配的唯一标识，即标签点ID，标签点ID一经创建后永不更改。
  */
  int id;
  /**
  *  标签点的数值类型。
  *  只读属性，在创建标签点时指定。
  */
  int type;
  /**
  *  标签点所属表ID。
  */
  int table;
  /**
  *  有关标签点的描述性文字。
  *  字节长度不要超出 \b RTDB_DESC_SIZE，多余的部分会被截断。
  */
  char desc[RTDB_DESC_SIZE];
  /**
  *  工程单位。
  *  字节长度不要超出 \b RTDB_UNIT_SIZE，多余的部分会被截断。
  */
  char unit[RTDB_UNIT_SIZE];
  /**
  *  是否存档。
  *  缺省值：ON，1；
  *  ON或1表示存档，OFF或0表示不存档。
  */
  rtdb_byte archive;
  /**
  *  数值位数。
  *  缺省值：-5；
  *  范围：>=-20、<=10；
  *  用来控制数值的显示格式；
  *  如果为0或正数，表示数值的小数位数，如果为负数，表示数值的有效位数。
  */
  short digits;
  /**
  *  停机状态字（Shutdown）
  *  缺省值：0；
  *  定义该点在停机状态下是否补写停机状态值。
  *  1 表示补写；0 表示不补写。
  */
  rtdb_byte shutdown;
  /**
  *  量程下限。
  *  缺省值：0；
  *  单位：标签点工程单位。
  */
  float lowlimit;
  /**
  *  量程上限。
  *  缺省值：100；
  *  单位：标签点工程单位。
  */
  float highlimit;
  /**
  *  是否阶跃。
  *  缺省值：OFF，0；
  *  该属性决定了中间值的计算是用阶梯还是连续的内插值替换；
  *  缺省情况下该属性为OFF，即中间值的计算是用内插值替换；
  *  如果被设置为ON，则中间值的数值同前一个有记录的数值相同。
  *  在历史数据检索中，本设置可能被外部输入的阶跃开关覆盖。
  */
  rtdb_byte step;
  /**
  *  典型值。
  *  缺省值：50；
  *  大于等于量程下限，小于等于量程上限。
  */
  float typical;
  /**
  *  是否压缩。
  *  缺省值：ON，1；
  *  如果该属性被关闭（OFF，0），任何到达数据存储服务器Server的数据都会被提交到历史数据库；否则（ON，1），只有满足压缩条件的数据才会被提交到历史数据库。
  *  需要手工录入的标签点应该将该属性设置为OFF，0。
  */
  rtdb_byte compress;
  /**
  *  压缩偏差。
  *  单位：标签点工程单位；
  *  缺省值：1；
  *  当有新的数值被提交到数据存储服务器Server，如果从上一个被提交到历史数据库的数值开始有数值超出了压缩偏差外，则上一个被提交到数据存储服务器Server的数值被视为关键数值；
  *  该属性与[压缩偏差百分比（the percent of compression deviation）]属性含义相同，该属性与量程（high
  *  limit - low limit）的百分比即[压缩偏差百分比（the percent of compression
  *  deviation）]属性的值；
  *  对该属性的修改将导致对[压缩偏差百分比（the percent of compression
  *  deviation）]的修改，同样，对[压缩偏差百分比（the percent of compression
  *  deviation）]的修改也将修改该属性，如果两个同时被修改，[压缩偏差百分比（the percent of compression
  *  deviation）]具有更高的优先权。
  */
  float compdev;
  /**
  *  压缩偏差百分比。
  *  单位：百分比；
  *  \see compdev。
  */
  float compdevpercent;
  /**
  *  最大压缩间隔。
  *  单位：秒；
  *  缺省值：28800；
  *  如果某个数值与上一个被提交到历史数据库的数值的时间间隔大于或等于最大压缩间隔，无论是否满足压缩条件，该数值都应该被视为关键数值从而被提交到历史数据库的数据队列；
  *  数据库中两个标签点间的时间间隔有可能超出该属性值，因为数据存储服务器Server可能长时间接收不到提交的数据，而且任何系统绝对不会自己创造数据。
  */
  int comptimemax;
  /**
  *  最短压缩间隔。
  *  单位：秒；
  *  缺省值：0；
  *  如果某个数值与上一个被提交到历史数据库的数值的时间间隔小于最短压缩间隔，该数值会被忽略；
  *  该属性有降噪（suppress noise）的作用。
  */
  int comptimemin;
  /**
  *  例外偏差。
  *  单位：标签点工程单位；
  *  缺省值：0.5；
  *  如果某个数值与上一个被提交到数据存储服务器Server的数值的偏差大于该例外偏差（以数值的工程单位为准），则该数值被视为例外数值，应该被提交到数据存储服务器Server；
  *  建议例外偏差应该小于等于压缩偏差的一半；
  *  该属性与[例外偏差百分比（The Percent of Exception Deviation）]属性含义相同，该属性与量程（high
  *  limit - low limit）的百分比即[例外偏差百分比（The Percent of Exception
  *  Deviation）]属性的值；
  *  对该属性的修改将导致对[例外偏差百分比（The Percent of Exception
  *  Deviation）]的修改，同样，对[例外偏差百分比（The Percent of Exception
  *  Deviation）]的修改也将修改该属性，如果两个同时被修改，[例外偏差百分比（The Percent of Exception
  *  Deviation）]具有更高的优先权。
  */
  float excdev;
  /**
  *  例外偏差百分比。
  *  单位：百分比；
  *  \see excdev。
  */
  float excdevpercent;
  /**
  *  最大例外间隔。
  *  单位：秒；
  *  缺省值：600；
  *  如果某个数值与上一个被提交到数据存储服务器Server的数值的时间间隔大于或等于最大例外间隔，无论是否满足例外条件，该数值都应该被视为例外数值从而被提交到数据存储服务器Server；
  *  数据库中两个标签点间的时间间隔有可能超出该属性值，因为接口可能长时间采集不到数据，而且任何系统绝对不会自己创造数据；
  *  如果要关闭例外过滤，设置该属性为0即可。
  */
  int exctimemax;
  /**
  *  最短例外间隔。
  *  单位：秒；
  *  缺省值：0；
  *  如果某个数值与上一个被提交到数据存储服务器Server的数值的时间间隔小于最短例外间隔，无论是否满足例外条件，该数值都会被忽略；
  *  该属性有降噪（suppress noise）的作用。
  */
  int exctimemin;
  /**
  *  标签点类别。
  *  RTDB_CLASS类型的组合，最多可以扩展至32种类型的组合；
  *  所有类别标签点均继承自"基本"类型标签点。
  *  不同类别的标签点具有不同的属性集，"采集"类别的标签点具有"设备标签"、"位号"、"自定义整数"和"自定义浮点数"等扩展属性，"计算"类别的标签点具有"扩展描述"、"触发机制"等扩展属性。
  */
  unsigned int classof;
  /**
  *  标签点属性最后一次被修改的时间。
  */
  rtdb_datetime_type changedate;
  /**
  *  标签点属性最后一次被修改的用户名。
  */
  char changer[RTDB_USER_SIZE];
  /**
  *  标签点被创建的时间。
  */
  rtdb_datetime_type createdate;
  /**
  *  标签点创建者的用户名。
  */
  char creator[RTDB_USER_SIZE];
  /**
  *  镜像收发控制。
  *  默认值：关闭，0
  *  开启镜像收发，1，既接收，又发送
  *  开启镜像接收，2，只接收，不发送
  *  开启镜像发送，3，只发送，不接收
  */
  rtdb_byte mirror;
  /**
  *  时间戳精度。
  *  默认值：秒，0；
  *  用于设定标签点的历史值在存储中精确到"秒"（0）还是"毫秒/纳秒"（1）。
  *  标签点一经创建就不允许修改该属性。
  */
  rtdb_byte millisecond;
  /**
  *  采集点扩展属性集存储地址索引。
  */
  unsigned int scanindex;
  /**
  *  计算点扩展属性集存储地址索引。
  */
  unsigned int calcindex;
  /**
  *  报警点扩展属性集存储地址索引。
  */
  unsigned int alarmindex;
  /**
  *  标签点全名，格式为“表名称.标签点名称”。
  */
  char table_dot_tag[RTDB_TAG_SIZE + RTDB_TAG_SIZE];
  /**
  * 统计加速。
  * 默认值：关，0；
  * 用于设定是否生成标签点统计信息，从而加速历史数据统计过程。
  */
  rtdb_byte summary;
  /**
  *  标签点对应自定义类型id，只用标签点类别为自定义类型时，才有意义。
  */
  rtdb_uint16 named_type_id;

  // 时间戳精度，0秒、1毫秒、2微秒、3纳秒
  rtdb_precision_type precision;
  /**
  *  基本标签点备用字节。
  */
  rtdb_byte padding[RTDB_PACK_OF_POINT];
} RTDB_POINT;

/**
* \ingroup ddatatype
* \typedef struct _RTDB_SCAN_POINT RTDB_SCAN_POINT
* \brief 采集标签点扩展属性集
* \see _RTDB_SCAN_POINT
*/
/**
* \ingroup dstruct
* \brief 采集标签点扩展属性集。8字节对齐条件下占用512字节。.
*/
typedef struct _RTDB_SCAN_POINT
{
  /**
  *  全库唯一标识。0表示无效。
  */
  int id;
  /**
  *  数据源。
  *  缺省值：空（NULL）；
  *  将标签点同某些接口或某些模块相关联；
  *  每个数据源字符串只允许由26个字母（大小写敏感）和数字（0-9）组成，字节长度不要超出 \b RTDB_SOURCE_SIZE，多余的部分会被截断。
  */
  char source[RTDB_SOURCE_SIZE];
  /**
  *  是否采集。
  *  缺省值：ON，1；
  *  该属性可能会被某些接口用到，如果该属性被关闭（OFF，0），从接口传来的数据可能不会被报告到数据库。
  */
  rtdb_byte scan;
  /**
  *  设备标签。
  *  缺省值：空（NULL）；
  *  字节长度不要超出 \b RTDB_INSTRUMENT_SIZE，多余的部分会被截断。
  */
  char instrument[RTDB_INSTRUMENT_SIZE];
  /**
  *  共包含五个设备位址，缺省值全部为0。
  */
  int locations[RTDB_LOCATIONS_SIZE];
  /**
  *  共包含两个自定义整数，缺省值全部为0。
  */
  int userints[RTDB_USERINT_SIZE];
  /**
  *  共包含两个自定义单精度浮点数，缺省值全部为0。
  */
  float userreals[RTDB_USERREAL_SIZE];
  /**
  *  采集标签点备用字节。
  */
  rtdb_byte padding[RTDB_PACK_OF_SCAN];
} RTDB_SCAN_POINT;

/**
* \ingroup ddatatype
* \typedef struct _RTDB_CALC_POINT RTDB_CALC_POINT
* \brief 计算点扩展属性集
* \see _RTDB_CALC_POINT
*/
/**
* \ingroup dstruct
* \brief 计算点扩展属性集。8字节对齐条件下占用2048字节。.
*/
typedef struct _RTDB_CALC_POINT
{
  /**
  *  全库唯一标识。0表示无效。
  */
  int id;
  /**
  *  实时方程式。
  *  缺省值：空（NULL）；
  *  字节长度不要超出 \b RTDB_EQUATION_SIZE，长度超长的方程式将被拒绝设置入库，返回一个错误，避免错误的方程式进入系统，引发不安全因素。
  */
  char equation[RTDB_EQUATION_SIZE];
  /**
  *  计算触发机制。枚举值参见 \b RTDB_TRIGGER。
  *  仅对"计算"类别标签点起作用，用于设置实时方程式服务对单个计算点的计算触发采用"事件触发"还是"周期触发"，
  *  对于"周期触发"以"事件触发"作为其先决判断条件，如果"事件触发"不满足，则不进行"周期触发"。
  */
  rtdb_byte trigger;
  /**
  *  计算结果时间戳参考, 枚举值参见 \b RTDB_TIME_COPY
  *  0: 表示采用计算时间作为计算结果时间戳；
  *  1: 表示采用输入标签点中的最晚时间作为计算结果时间戳；
  *  2: 表示采用输入标签点中的最早时间作为计算结果时间戳。
  */
  rtdb_byte timecopy;
  /**
  *  对于“周期触发”的计算点，设定其计算周期，单位：秒。
  */
  int period;
  // 计算标签点备用字节。
  // rtdb_byte padding[RTDB_PACK_OF_CALC];
} RTDB_CALC_POINT;

/// 计算点扩展属性集。
/**
*  计算点扩展属性集。8字节对齐条件下占用512字节。
*/
typedef struct _RTDB_MIN_CALC_POINT
{
  /**
  *  全库唯一标识。0表示无效。
  */
  int id;
  /**
  *  实时方程式。
  *  缺省值：空（NULL）；
  *  字节长度不要超出480，如果字节长度超过480，请考虑调用RTDB_CALC_POINT或RTDB_MAX_CALC_POINT
  */
  char equation[RTDB_MIN_EQUATION_SIZE];
  /**
  *  计算触发机制。枚举值参见 RTDB_TRIGGER。
  *  仅对"计算"类别标签点起作用，用于设置实时方程式服务对单个计算点的计算触发采用"事件触发"还是"周期触发"，
  *  对于"周期触发"以"事件触发"作为其先决判断条件，如果"事件触发"不满足，则不进行"周期触发"。
  */
  rtdb_byte trigger;
  /**
  *  计算结果时间戳参考, 枚举值参见 RTDB_TIME_COPY
  *  0: 表示采用计算时间作为计算结果时间戳；
  *  1: 表示采用输入标签点中的最晚时间作为计算结果时间戳；
  *  2: 表示采用输入标签点中的最早时间作为计算结果时间戳。
  */
  rtdb_byte timecopy;
  /**
  *  对于“周期触发”的计算点，设定其计算周期，单位：秒。
  */
  int period;
  /**
  *  此方程式中保存的是否是方程式
  */
  rtdb_byte is_equation;
  // 计算标签点备用字节。
  rtdb_byte padding[RTDB_PACK_OF_MIN_CALC];
} RTDB_MIN_CALC_POINT;


/**
* \ingroup ddatatype
* \typedef struct _RTDB_MAX_CALC_POINT RTDB_MAX_CALC_POINT
* \brief 最大长度计算点扩展属性集
* \see _RTDB_MAX_CALC_POINT
*/
/**
* \ingroup dstruct
* \brief 最大长度计算点扩展属性集。8字节对齐条件下占用6K字节。.
*/
typedef struct _RTDB_MAX_CALC_POINT
{
  /**
  *  全库唯一标识。0表示无效。
  */
  int id;
  /**
  *  实时方程式。
  *  缺省值：空（NULL）；
  *  字节长度不要超出 \b RTDB_MAX_EQUATION_SIZE，长度超长的方程式将被拒绝设置入库，返回一个错误，避免错误的方程式进入系统，引发不安全因素。
  */
  char equation[RTDB_MAX_EQUATION_SIZE];
  /**
  *  计算触发机制。枚举值参见 \b RTDB_TRIGGER。
  *  仅对"计算"类别标签点起作用，用于设置实时方程式服务对单个计算点的计算触发采用"事件触发"还是"周期触发"，
  *  对于"周期触发"以"事件触发"作为其先决判断条件，如果"事件触发"不满足，则不进行"周期触发"。
  */
  rtdb_byte trigger;
  /**
  *  计算结果时间戳参考, 枚举值参见 \b RTDB_TIME_COPY
  *  0: 表示采用计算时间作为计算结果时间戳；
  *  1: 表示采用输入标签点中的最晚时间作为计算结果时间戳；
  *  2: 表示采用输入标签点中的最早时间作为计算结果时间戳。
  */
  rtdb_byte timecopy;
  /**
  *  对于“周期触发”的计算点，设定其计算周期，单位：秒。
  */
  int period;
  // 计算标签点备用字节。
  rtdb_byte padding[RTDB_PACK_OF_MAX_CALC];
} RTDB_MAX_CALC_POINT;

/// 标签点派别结构
typedef union _RTDB_TAG_FACTION
{
  unsigned int faction;
  struct
  {
    rtdb_byte type;
    rtdb_int8 precision;
    rtdb_byte reserved1;
    rtdb_byte reserved2;
  } factions;
} RTDB_TAG_FACTION;

/**
* \ingroup denum
* \brief 操作系统类型
* \see rtdb_get_linked_ostype
*/
typedef enum _RTDB_OS_TYPE
{
  RTDB_OS_WINDOWS,        //!< Windows 操作系统
  RTDB_OS_LINUX,          //!< Linux 操作系统
  RTDB_OS_INVALID = 50,   //!< 无效的操作系统
} RTDB_OS_TYPE;


/**
* \ingroup denum
* \brief 用户权限.
*/
typedef enum _RTDB_PRIV_GROUP
{
  RTDB_RO,      //!< 只读
  RTDB_DW,      //!< 数据记录
  RTDB_TA,      //!< 标签点表管理员
  RTDB_SA,      //!< 数据库管理员
} RTDB_PRIV_GROUP;

/**
* \ingroup denum
* \brief 标签点变更原因，用于标签点订阅
* \see rtdbb_tags_change_event rtdbb_subscribe_tags
*/
typedef enum _RTDB_TAG_CHANGE_REASON
{
  RTDB_TAG_CREATED = 1,  //!< 标签点被创建
  RTDB_TAG_UPDATED,      //!< 标签点属性被更新
  RTDB_TAG_REMOVED,      //!< 标签点被放入回收站
  RTDB_TAG_RECOVERD,     //!< 标签点被恢复
  RTDB_TAG_PURGED,       //!< 标签点被清除
  RTDB_TAB_UPDATED,      //!< 标签点表被重命名
  RTDB_TAB_REMOVED,      //!< 标签点表被删除
} RTDB_TAG_CHANGE_REASON;

/**
* \ingroup dmacro
* \def RTDB_ARVEX_MAX
* \brief 每个历史数据存档文件最多允许有99个附属文件.
*/
#define RTDB_ARVEX_MAX          99

/*
* \ingroup denum
* \brief 存档文件管理状态
*/
typedef enum _RTDB_ARCHIVE_MANAGE_TYPE
{
  GAMT_NORMAL = 0,			//!< 正常状态，不做处理
  GAMT_NOT_MANAGED = 1,     //!< 未被管理(已解列)
  GAMT_MANAGED = 2,			//!< 被管理状态
} RTDB_ARCHIVE_MANAGE_TYPE;

/*
* \ingroup denum
* \brief 存档文件压缩类型
*/
typedef enum _RTDB_ARCHIVE_COMPRESS_TYPE
{
    GACT_NORMAL = 0,				//!< 未压缩，定长数据块
    GACT_COMPRESS = 1,				//!< 无损压缩，不定长数据块
    GACT_COMPRESS_TWO_STAGE = 2,	//!< 无损压缩，不定长数据块，两阶段压缩
} RTDB_ARCHIVE_COMPRESS_TYPE;

/**
* \ingroup ddatatype
* \typedef struct _RTDB_HEADER_PAGE RTDB_HEADER_PAGE
* \brief 历史数据存档文件头部信息
* \see _RTDB_HEADER_PAGE
*/

#define RTDB_VER_CODE_SIZE 6

/**
* \ingroup dstruct
* \brief 历史数据存档文件头部信息.
*/
typedef struct _RTDB_HEADER_PAGE
{
  int db_ver;                           //!< 所属数据库版本
  int data_ver;                         //!< 数据格式变更版本
  rtdb_datetime_type begin;             //!< 数据起始时间
  rtdb_datetime_type end;               //!< 数据结束时间
  rtdb_int64 real_size;                 //!< 实际写入的字节数 用来统计实际使用率
  rtdb_datetime_type create_time;       //!< 创建时间
  rtdb_datetime_type modify_time;       //!< 修改时间
  rtdb_datetime_type merge_time;        //!< 上次合并时间
  rtdb_datetime_type arrange_time;      //!< 上次整理时间
  rtdb_datetime_type reindex_time;      //!< 上次重建索引时间
  rtdb_datetime_type backup_time;       //!< 上次备份时间
  rtdb_int64 rated_capacity;            //!< 创建时容量（额定容量）
  rtdb_int64 capacity;                  //!< 当前容量，文件内包含的总数据页数，不包括头页。
  rtdb_int64 size;                      //!< 实际使用量，已被占用的数据页数，不包括头页。
  rtdb_int64 ex_capacity;               //!< 附属文件容量
  rtdb_byte is_main;                    //!< 是主文件还是附属文件
  rtdb_byte page_size;                  //!< 单个页的字节尺寸，单位为KB。
  rtdb_byte id_or_count;                //!< 主文件在这里存放附属文件数量，附属文件在这里存放自身的ID。
  rtdb_byte auto_merge;                 //!< 启用自动合并
  rtdb_byte auto_arrange;               //!< 启用自动整理
  rtdb_byte merged;                     //!< 1：已合并，0：尚未合并或合并后又产生了新的附属文件。
  rtdb_byte arranged;                   //!< 1：已整理，0：尚未整理过或整理后内容发生变更。
  rtdb_byte index_type;                 //!< 索引类型，0 为红黑树，1为跳跃链表。默认为红黑树
  char file_name[RTDB_FILE_NAME_SIZE];  //!< 在这里存放自己的文件名。
  rtdb_uint32 crc32;                    //!< 以上内容的CRC32校验码，暂不启用。
  rtdb_byte index_in_mem;               //!< 索引是否加载到内存中
  rtdb_byte manage_type;                //!< 被管理类型，RTDB_ARCHIVE_MANAGE_TYPE， 0:兼容模式  1:未管理 2:管理中
  rtdb_byte compress_type;              //!< 压缩类型，RTDB_ARCHIVE_COMPRESS_TYPE， 0:未压缩 1:无损压缩 2:两阶段压缩
  rtdb_byte reserve[1];                 //!< 保留字节，用于内存对齐
  rtdb_error status;                    //!< 存档文件的当前状态，表示存档文件操作中遇到的异常错误码，正常为RtE_OK，可能遇到的异常有：RtE_INDEX_NOT_READY, RtE_CAN_NOT_CREATE_INDEX, RtE_NOT_ENOUGH_SPACE
  rtdb_byte reserve2[4];                //!< 保留字节，用于内存对齐
  rtdb_int64 used_size;                 //!< 实际使用字节数，不包括存档文件头
  rtdb_int64 block_count;               //!< 数据块数量
  rtdb_int64 del_block_size;            //!< 逻辑上删除的数据块字节数
  rtdb_int64 total_count;               //!< 数据总条数
  rtdb_uint16 big_page_size;            //!< 数据页大小，单位为KB
  char vercode[RTDB_VER_CODE_SIZE];     //!< 授权信息
  char reserve_1[48];
} RTDB_HEADER_PAGE;       //保证256字节

/**
* \ingroup dstruct
* \brief 用户信息.
*/
typedef struct _RTDB_USER_INFO
{
  char user[RTDB_USER_SIZE];
  rtdb_int32 length;
  rtdb_int32 privilege;
  rtdb_int8 islocked;
  char reserve_1[15];
} RTDB_USER_INFO;         // 44 bytes

/**
 * \ingroup dmacro
 * \def RTDB_MAX_PATH
 * \brief 系统支持的最大路径长度.
 */
#define RTDB_MAX_PATH                 2048

 /**
 * \ingroup dmacro
 * \def RTDB_MAX_HOSTNAME_SIZE
 * \brief 系统支持的最大主机名长度.
 */
#define RTDB_MAX_HOSTNAME_SIZE        1024

 /**
 * \ingroup ddatatype
 * \typedef struct _RTDB_HOST_CONNECT_INFO RTDB_HOST_CONNECT_INFO
 * \brief 连接到RTDB数据库服务器的连接信息
 * \see _RTDB_HOST_CONNECT_INFO
 */

 /**
 * \ingroup ddatatype
 * \typedef struct _RTDB_HOST_CONNECT_INFO *PRTDB_HOST_CONNECT_INFO
 * \brief 连接到RTDB数据库服务器的连接信息
 * \see _RTDB_HOST_CONNECT_INFO
 */

 /**
 * \ingroup dstruct
 * \brief 连接到RTDB数据库服务器的连接信息.
 */
typedef struct _RTDB_HOST_CONNECT_INFO
{
    rtdb_int32 ipaddr;                                            //!< 连接的客户端IP地址
    rtdb_uint16 port;                                             //!< 连接端口
    rtdb_int32 job;                                               //!< 连接最近处理的任务
    rtdb_datetime_type job_time;                                  //!< 最近处理任务的时间
    rtdb_datetime_type connect_time;                              //!< 客户端连接时间
    char client[RTDB_MAX_HOSTNAME_SIZE];                          //!< 连接的客户端主机名称
    char process[RTDB_PATH_SIZE + RTDB_FILE_NAME_SIZE];           //!< 连接的客户端程序名
    char user[RTDB_USER_SIZE];                                    //!< 登录的用户
    rtdb_int32 length;                                            //!< 记录用户名长度，用于加密传输
} RTDB_HOST_CONNECT_INFO, *PRTDB_HOST_CONNECT_INFO;

/**
* \ingroup dstruct
* \brief 连接到RTDB数据库服务器的连接信息, IPV6版本
*/
typedef struct _RTDB_HOST_CONNECT_INFO_IPV6
{
    rtdb_int32 ipaddr;                                            //!< 连接的客户端IP地址
    char ipaddr6[RTDB_IPV6_ADDR_SIZE];                            //!<ipv6地址
    rtdb_uint16 port;                                             //!< 连接端口
    rtdb_int32 job;                                               //!< 连接最近处理的任务
    rtdb_datetime_type job_time;                                  //!< 最近处理任务的时间
    rtdb_datetime_type connect_time;                              //!< 客户端连接时间
    char client[RTDB_MAX_HOSTNAME_SIZE];                          //!< 连接的客户端主机名称
    char process[RTDB_PATH_SIZE + RTDB_FILE_NAME_SIZE];           //!< 连接的客户端程序名
    char user[RTDB_USER_SIZE];                                    //!< 登录的用户
    rtdb_int32 length;                                            //!< 记录用户名长度，用于加密传输
} RTDB_HOST_CONNECT_INFO_IPV6, * PRTDB_HOST_CONNECT_INFO_IPV6;

/*
* \ingroup ddatatype
* \typedef struct _RTDB_GRAPH_DATA RTDB_GRAPH_DATA
* \brief 计算标签点方程式拓扑图键值对信息
* \see _RTDB_GRAPH_DATA
*/

/*
* \ingroup dstruct
* \brief 计算标签点方程式拓扑图键值对信息
*/
typedef struct _RTDB_GRAPH_DATA
{
  int id;                               //!< 标签点ID
  int parent_id;                        //!< 标签点父ID，即父ID方程式中引用了子ID做运算
  char tag[RTDB_TAG_SIZE];              //!< 标签点名称
  char error_msg[RTDB_DESC_SIZE];       //!< 无效标签点错误信息
} RTDB_GRAPH;

/**
* \ingroup denum
* \brief 标签点拓扑图类型.
*/
typedef enum _RTDB_GRAPH_FLAG
{
  RTDB_GRAPH_BEGIN = -1,
  RTDB_GRAPH_ALL,       //!< 任何有关联的标签的关系图
  RTDB_GRAPH_DIRECT,    //!< 有直接关系的关系图
  RTDB_GRAPH_END,
} RTDB_GRAPH_FLAG;

/**
* \ingroup denum
* \brief 历史存档文件状态.
*/
typedef enum _RTDB_ARCHIVE_STATE
{
  RTDB_INVALID_ARCHIVE, //!< 0:无效
  RTDB_ACTIVED_ARCHIVE, //!< 1:活动
  RTDB_NORMAL_ARCHIVE,  //!< 2:普通
  RTDB_READONLY_ARCHIVE //!< 3:只读
} RTDB_ARCHIVE_STATE;

/**
* \ingroup denum
* \brief 查询系统参数时对应的索引
*/
typedef enum _RTDB_DB_PARAM_INDEX
{
	/// string parameter.
	RTDB_PARAM_STR_FIRST = 0x0,
	RTDB_PARAM_TABLE_FILE = RTDB_PARAM_STR_FIRST,     // 标签点表文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_BASE_FILE,                             // 基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_SCAN_FILE,                             // 采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_CALC_FILE,                             // 计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_SNAP_FILE,                             // 标签点快照文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_LIC_FILE,                              // 协议文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_HIS_FILE,                              // 历史信息文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_LOG_DIR,                               // 服务器端日志文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_USER_FILE,                             // 用户权限信息文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_SERVER_FILE,                           // 网络服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_EQAUTION_FILE,                         // 方程式服务进程与其它进程交互所用的共享内存文件，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_ARV_PAGES_FILE,                        // 历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_ARVEX_PAGES_FILE,                      // 补历史数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_ARVEX_PAGES_BLOB_FILE,                 // 补历史数据blob、str缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_AUTH_FILE,                             // 信任连接段信息文件路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_RECYCLED_BASE_FILE,                    // 可回收基本标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_RECYCLED_SCAN_FILE,                    // 可回收采集标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_RECYCLED_CALC_FILE,                    // 可回收计算标签点文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_AUTO_BACKUP_PATH,                      // 自动备份目的地全路径，必须以“\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_SERVER_SENDER_IP,                      // 镜像发送地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE
	RTDB_PARAM_BLACKLIST_FILE,                        // 连接黑名单文件路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_DB_VERSION,                            // 数据库版本
	RTDB_PARAM_LIC_USER,                              // 授权单位
	RTDB_PARAM_LIC_TYPE,                              // 授权方式
	RTDB_PARAM_INDEX_DIR,                             // 索引文件存放目录，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_MIRROR_BUFFER_PATH,                    // 镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_MIRROR_EX_BUFFER_PATH,                 // 补写镜像缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_EQAUTION_PATH_FILE,                    // 方程式长度超过规定长度时进行保存的文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_TAGS_FILE,                             // 标签点关键属性文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_RECYCLED_SNAP_FILE,                    // 可回收标签点快照事件文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_SWAP_PAGE_FILE,					      // 历史数据交换页文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_PAGE_ALLOCATOR_FILE,				      /* 活动存档数据页分配器文件全路径，字符串最大长度为 RTDB_MAX_PATH
													    该系统配置项2.1版数据库在使用，3.0数据库已去掉，但为了保证系统选项索引号
													    与2.1一致，此处不能去掉。便于java sdk统一调用*/
	RTDB_PARAM_NAMED_TYPE_FILE,					      // 自定义类型配置信息全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_STRBLOB_MIRROR_PATH,				      // BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_STRBLOB_MIRROR_EX_PATH,			      // 补写BLOB/STRING镜像数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_BUFFER_DIR,						      // 临时数据缓存路径
	RTDB_PARAM_POOL_CACHE_FLIE,					      // 曲线池索引文件全路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_POOL_DATA_FILE_DIR,				      // 曲线池缓存文件目录，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_ARCHIVE_FILE_PATH,				      // 存档文件低速存储区路径，字符串最大长度为 RTDB_MAX_PATH
	RTDB_PARAM_LIC_VERSION_TYPE,                      // 授权版本
    RTDB_PARAM_AUTO_MOVE_PATH,                        // 自动移动目的地全路径，必须以“\”或“/”结束，字符串最大长度为 RTDB_MAX_PATH
    RTDB_PARAM_REPLICATION_BUFFER_PATH,				  // 双活：数据同步缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
    RTDB_PARAM_REPLICATION_EX_BUFFER_PATH,			  // 双活：数据同步补写数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
    RTDB_PARAM_STRBLOB_REPLICATION_BUFFER_PATH,	      // 双活：数据同步BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
    RTDB_PARAM_STRBLOB_REPLICATION_EX_BUFFER_PATH,	  // 双活：数据同步补写BLOB/STRING数据缓存文件全路径，字符串最大长度为 RTDB_MAX_PATH
    RTDB_PARAM_REPLICATION_GROUP_IP,                  // 双活：同步组地址，字符串最大长度为 RTDB_MAX_HOSTNAME_SIZE
    RTDB_PARAM_ARC_FILENAME_PREFIX_WHEN_USING_DATE,   // 是否归档文件使用日期作为文件名
    RTDB_PARAM_HOT_ARCHIVE_FILE_PATH,                 // 存档文件高速存储区路径，字符串最大长度为 RTDB_MAX_PATH
    RTDB_PARAM_STR_LAST,

	// int parameter.
	RTDB_PARAM_INT_FIRST = 0x1000,
	RTDB_PARAM_LIC_TABLES_COUNT = RTDB_PARAM_INT_FIRST, // 协议中限定的标签点表数量
	RTDB_PARAM_LIC_TAGS_COUNT,                        // 协议中限定的所有标签点数量
	RTDB_PARAM_LIC_SCAN_COUNT,                        // 协议中限定的采集标签点数量
	RTDB_PARAM_LIC_CALC_COUNT,                        // 协议中限定的计算标签点数量
	RTDB_PARAM_LIC_ARCHICVE_COUNT,                    // 协议中限定的历史存档文件数量
	RTDB_PARAM_SERVER_IPC_SIZE,                       // 网络服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）
	RTDB_PARAM_EQUATION_IPC_SIZE,                     // 方程式服务进程与其它进程进行交互所使用的共享内存池的字节尺寸（单位：B）
	RTDB_PARAM_HASH_TABLE_SIZE,                       // 标签点求余哈希表的尺寸
	RTDB_PARAM_TAG_DELETE_TIMES/*obsolete*/,          // 可整库删除标签点的次数
	RTDB_PARAM_SERVER_PORT,                           // 网络服务独立服务器端口
	RTDB_PARAM_SERVER_SENDER_PORT,                    // 网络服务镜像发送端口
	RTDB_PARAM_SERVER_RECEIVER_PORT,                  // 网络服务镜像接收端口
	RTDB_PARAM_SERVER_MODE,                           // 网络服务启动模式
	RTDB_PARAM_SERVER_CONNECTION_COUNT,               // 协议中限定网络服务连接并发数量
	RTDB_PARAM_ARV_PAGES_NUMBER,                      // 历史数据缓存中的页数量
	RTDB_PARAM_ARVEX_PAGES_NUMBER,                    // 补历史数据缓存中的页数量
	RTDB_PARAM_EXCEPTION_AT_SERVER,                   // 是否由服务器进行例外判定
	RTDB_PARAM_ARV_PAGE_RECYCLE_DELAY/*obsolete*/,    // 历史数据缓存页回收延时（毫秒）
	RTDB_PARAM_EX_ARCHIVE_SIZE,                       // 历史数据存档文件文件自动增长大小（单位：MB）
	RTDB_PARAM_ARCHIVE_BATCH_SIZE,                    // 历史存储值分段查询个数
	RTDB_PARAM_DATAFILE_PAGESIZE,                     // 系统数据文件页大小
	RTDB_PARAM_ARV_ASYNC_QUEUE_NORMAL_DOOR,           // 历史数据缓存队列中速归档区（单位：百分比）
	RTDB_PARAM_INDEX_ALWAYS_IN_MEMORY,                // 常驻内存的历史数据索引大小（单位：MB）
	RTDB_PARAM_DISK_MIN_REST_SIZE,                    // 最低可用磁盘空间（单位：MB）
	RTDB_PARAM_MIN_SIZE_OF_ARCHIVE,                   // 历史存档文件和附属文件的最小尺寸（单位：MB）
	RTDB_PARAM_DELAY_OF_AUTO_MERGE_OR_ARRANGE,        // 自动合并/整理最小延迟（单位：小时）
	RTDB_PARAM_START_OF_AUTO_MERGE_OR_ARRANGE,        // 自动合并/整理开始时间（单位：点钟）
	RTDB_PARAM_STOP_OF_AUTO_MERGE_OR_ARRANGE,         // 自动合并/整理停止时间（单位：点钟）
	RTDB_PARAM_START_OF_AUTO_BACKUP,                  // 自动备份开始时间（单位：点钟）
	RTDB_PARAM_STOP_OF_AUTO_BACKUP,                   // 自动备份停止时间（单位：点钟）
	RTDB_PARAM_MAX_LATENCY_OF_SNAPSHOT,               // 允许服务器时间之后多少小时内的数据进入快照（单位：小时）
	RTDB_PARAM_PAGE_ALLOCATOR_RESERVE_SIZE,/*obsolete*/// 活动页分配器预留大小（单位：KB）， 0 表示使用操作系统视图大小
	RTDB_PARAM_INCLUDE_SNAPSHOT_IN_QUERY,             // 决定取样本值和统计值时，快照是否应该出现在查询结果中
	RTDB_PARAM_LIC_BLOB_COUNT,                        // 协议中限定的字符串或BLOB类型标签点数量
	RTDB_PARAM_MIRROR_BUFFER_SIZE,                    // 镜像文件大小（单位：GB）
	RTDB_PARAM_BLOB_ARVEX_PAGES_NUMBER,               // blob、str补历史的默认缓存页数量
	RTDB_PARAM_MIRROR_EVENT_QUEUE_CAPACITY,           // 镜像缓存队列容量
	RTDB_PARAM_NOTIFY_NOT_ENOUGH_SPACE,               // 提示磁盘空间不足，一旦启用，设置为ON，则通过API返回大错误码，否则只记录日志
	RTDB_PARAM_ARCHIVE_FIXED_RANGE,                   // 历史数据存档文件的固定时间范围，默认为0表示不使用固定时间范围（单位：分钟）
	RTDB_PARAM_ONE_CLINET_MAX_CONNECTION_COUNT,       // 单个客户端允许的最大连接数，默认为0表示不限制
	RTDB_PARAM_ARV_PAGES_CAPACITY,                    // 历史数据缓存所占字节大小，单位：字节
	RTDB_PARAM_ARVEX_PAGES_CAPACITY,                  // 历史数据补写缓存所占字节大小，单位：字节
	RTDB_PARAM_BLOB_ARVEX_PAGES_CAPACITY,             // blob、string类型标签点历史数据补写缓存所占字节大小，单位：字节
	RTDB_PARAM_LOCKED_PAGES_MEM,                      // 指定分配给数据库用的内存大小，单位：MB
	RTDB_PARAM_LIC_RECYCLE_COUNT,                     // 协议中回收站的容量
	RTDB_PARAM_ARCHIVED_POLICY,                       // 快照数据和补写数据的归档策略
	RTDB_PARAM_NETWORK_ISOLATION_ACK_BYTE,            // 网络隔离装置ACK字节

	RTDB_PARAM_ENABLE_LOGGER,                         // 启用日志输出，0为不启用
	RTDB_PARAM_LOG_ENCODE,                            // 启用日志加密，0为不启用
	RTDB_PARAM_LOGIN_TRY,                             // 启用登录失败次数验证，0为不启用
	RTDB_PARAM_USER_LOG,                              // 启用用户详细日志，0为不启用
	RTDB_PARAM_COVER_WRITE_LOG,                       // 启用日志覆盖写功能，0为不启用

	RTDB_PARAM_LIC_NAMED_TYPE_COUNT,				  // 协议中限定的自定义类型标签点数量
	RTDB_PARAM_MIRROR_RECEIVER_THREADPOOL_SIZE,		  // 镜像接收线程数量
	RTDB_PARAM_SNAPSHOT_USE_ARCHIVE_INTERFACE,		  // 按照补历史流程归档快照数据页
	RTDB_PARAM_NO_ARCDATA_WRITE_LOG,				  // 归档无对应存档文件的数据时记录日志
	RTDB_PARAM_PUT_ARCHIVE_THREAD_NUM,				  // 补历史归档线程数
	RTDB_PARAM_ARVEX_DATA_ARCHIVED_THRESHOLD,		  // 单次补写数据归档阈值
	RTDB_PARAM_SNAPSHOT_FLUSH_BUFFER_DELAY,			  // 快照服务的共享缓存刷新到磁盘的周期
	RTDB_PARAM_DATA_SPEED,							  // 查询时使用加速统计
	RTDB_PARAM_USE_NEW_PLOT_ALGO,					  // 启用新的曲线算法
	RTDB_PARAM_QUERY_THREAD_POOL_SIZE,				  // 曲线查询线程池中线程数量
	RTDB_PARAM_ARCHIVED_VALUES,                       // 使用查询线程池查询历史数据
	RTDB_PARAM_ARCHIVED_VALUES_COUNT,                 // 使用查询线程池查询历史数据的条数

	RTDB_PARAM_POOL_USE_FLAG,						  // 启用曲线池
	RTDB_PARAM_POOL_OUT_LOG_FLAG,					  // 输出曲线池日志
	RTDB_PARAM_POOL_TIME_USE_POOL_FLAG,				  // 使用曲线池缓存计算插值
	RTDB_PARAM_POOL_MAX_POINT_COUNT,				  // 曲线池的标签点容量
	RTDB_PARAM_POOL_ONE_FILE_SAVE_POINT_COUNT,		  // 曲线池每个数据文件的标签点容量
	RTDB_PARAM_POOL_SAVE_MEMORY_SIZE,				  // 曲线缓存退出时临时缓冲区大小
	RTDB_PARAM_POOL_MIN_TIME_UNIT_SECONDS,			  // 曲线池缓存数据当前时间单位
	RTDB_PARAM_POOL_TIME_UNIT_VIEW_RATE,			  // 曲线池查询数据最小时间单位显示系数
	RTDB_PARAM_POOL_TIMER_INTERVAL_SECONDS,			  // 曲线池定时器刷新周期
	RTDB_PARAM_POOL_PERF_TIMER_INTERVAL_SECONDS,	  // 曲线池性能计算点刷新周期

	RTDB_PARAM_ARCHIVE_INIT_FILE_SIZE,				  // 存档文件初始大小
	RTDB_PARAM_ARCHIVE_INCREASE_MODE,				  // 存档文件增长模式
	RTDB_PARAM_ARCHIVE_INCREASE_SIZE,				  // 固定模式下文件增长大小
	RTDB_PARAM_ARCHIVE_INCREASE_PERCENT,			  // 百分比模式下增长百分比
	RTDB_PARAM_ALLOW_CONVERT_SKL_TO_RBT_INDEX,	      // 跳跃链表转换到红黑树
	RTDB_PARAM_EARLY_DATA_TIME,					      // 冷数据时间
	RTDB_PARAM_EARLY_INDEX_TIME,					  // 自动转换索引时间
	RTDB_PARAM_ARRANGE_RBT_TIME,					  // 整理存档文件时决定索引格式的时间轴
	RTDB_PARAM_ENABLE_BIG_DATA,						  // 将存档文件全部读取到内存中
	RTDB_PARAM_AUTO_ARRANGE_PERCENT,				  // 自动整理存档文件时的实际使用率
	RTDB_PARAM_EARLY_ARRANGE_TIME,					  // 自动整理存档文件的时间
	RTDB_PARAM_MIN_AUTO_ARRANGE_ARCFILE_PERCENT,	  // 自动整理存档文件时的最小使用率
	RTDB_PARAM_ARRANGE_ARC_WITH_MEMORY,			      // 在内存中整理存档文件
	RTDB_PARAM_ARAANGE_ARC_MAX_MEM_PERCENT,		      // 整理存档文件最大内存使用率
	RTDB_PARAM_MAX_DISK_SPACE_PERCENT,			      // 磁盘最大使用率
	RTDB_PARAM_USE_DISPATH,						      // windows 用 linux 已禁用,是否启用转发服务
	RTDB_PARAM_USE_SMART_PARAM,					      // windows 用 linux 已禁用,是否使用推荐参数
	RTDB_PARAM_SUBSCRIBE_SNAPSHOT_COUNT,              // 单连接快照事件订阅个数
	RTDB_PARAM_SUBSCRIBE_QUEUE_SIZE,                  // 订阅事件队列大小
	RTDB_PARAM_SUBSCRIBE_TIMEOUT,                     // 订阅事件超时时间

	RTDB_PARAM_MIRROR_COMPRESS_ONOFF,				  // 镜像报文压缩是否打开
	RTDB_PARAM_MIRROR_COMPRESS_TYPE,				  // 镜像报文压缩类型
	RTDB_PARAM_MIRROR_COMPRESS_MIN,					  // 镜像报文压缩最小值
	RTDB_PARAM_ARCHIVE_ROLL_TIME,					  // 存档文件滚动时间轴
	RTDB_PARAM_HANDLE_TIME_OUT,						  // 连接超时断开，单位：秒
    RTDB_PARAM_MOVE_ARV_TIME,					      // 移动存档文件时决定移动存档的时间轴
	RTDB_PARAM_USE_NEW_INTERP_ALGO,					  // 启用新的插值算法
	RTDB_PARAM_ENABLE_REPLICATION,                    // 启用双活，0为不启用，1为启用
	RTDB_PARAM_REPLICATION_GROUP_PORT,                // 双活：同步组端口
	RTDB_PARAM_REPLICATION_THREAD_SIZE,				  // 双活：同步线程数
	RTDB_PARAM_FORCE_ARCHIVE_INCOMPLETE_DATA_PAGE_DELAY, // 强制归档补历史缓存里面未满数据页的延迟时间
	RTDB_PARAM_ARCHIVE_ROLL_DISK_PERCENTAGE,          // 存档文件滚动存储空间百分比
    RTDB_PARAM_ENABLE_IPV6,                           // 启用ipv6设置
    RTDB_PARAM_ENABLE_USE_ARCHIVED_VALUE,             // 按条件获取历史值时，是否直接获取条件中点的历史值，0:获取插值，1:获取历史值
    RTDB_PARAM_TIMESTAMP_TYPE,                        // 获取服务器时间戳类型
    RTDB_PARAM_ARC_FILENAME_USING_DATE,				  // 是否归档文件使用日期作为文件名
    RTDB_PARAM_LOG_MAX_SPACE,				          // 日志文件占用的最大磁盘空间
    RTDB_PARAM_LOG_FILE_SIZE,						  // 单个日志文件大小
    RTDB_PARAM_IGNORE_TO_WRITE_NOARCBUFFER,           // 是否丢弃补历史数据
    RTDB_PARAM_ARCHIVES_COUNT_FOR_CALC,               // 统计存档文件平均大小的存档文件个数

	// !!!添加参数前必看!!!!
	// 出于兼容性考虑，如果新加INT类型参数
	// 一定要紧挨着 RTDB_PARAM_INT_LAST，
	// 在前一个位置加参数
	// 这样，客户端应用程序如GEM，没有高版本低版本参数错位问题，顶多是某一项功能不能用，而非功能大批量不能用。
	// 对于字符串类型的参数也是同理
	// 字符串类型的参数添加在 RTDB_PARAM_STR_LAST 之前且紧挨着RTDB_PARAM_STR_LAST
	RTDB_PARAM_INT_LAST,

	//exp int param
	RTDB_PARAM_EXP_INT_FIRST = 0x2000,
	RTDB_PARAM_MAX_BLOB_SIZE = RTDB_PARAM_EXP_INT_FIRST,   // blob、str类型数据在数据库中允许的最大长度
	RTDB_PARAM_EXP_INT_LAST,
} RTDB_DB_PARAM_INDEX;

//表名称枚举
typedef enum _RTDB_TABLE_ID
{
  RTDB_TABLE_BASE = 1,
  RTDB_TABLE_SCAN = 2,
  RTDB_TABLE_CALC = 4,
  RTDB_TABLE_ALARM = 8,
} RTDB_TABLE_ID;
#define RTDB_TABLE_ID_CONTAIN_BASE(TABLEID) (TABLEID & RTDB_TABLE_BASE)
#define RTDB_TABLE_ID_CONTAIN_SCAN(TABLEID) (TABLEID & RTDB_TABLE_SCAN)
#define RTDB_TABLE_ID_CONTAIN_CALC(TABLEID) (TABLEID & RTDB_TABLE_CALC)
#define RTDB_TABLE_ID_CONTAIN_ALARM(TABLEID) (TABLEID & RTDB_TABLE_ALARM)
//表字段信息
#define RTDB_TAG_FIELD_NAME_LENGTH  20
#define RTDB_TAG_FIELD_TYPE_LENGTH    12
typedef struct TAG_RTDB_TAG_FIELD
{
	RTDB_TABLE_ID table_id;                                   // 列类型
	rtdb_int16 column_index;                                  // 字段序号
	rtdb_int16 column_length;                                 // 字段长度
	char column_name[RTDB_TAG_FIELD_NAME_LENGTH];             // 字段名称
	char column_type[RTDB_TAG_FIELD_TYPE_LENGTH];             // 字段类型
} RTDB_TAG_FIELD;

typedef enum _RTDB_TAG_FIELD_INDEX
{
  RTDB_TAG_INDEX_BASE_FIRST = 0x0,
  RTDB_TAG_INDEX_TAG = RTDB_TAG_INDEX_BASE_FIRST,   //!< tag
  RTDB_TAG_INDEX_ID,                                //!< id
  RTDB_TAG_INDEX_TYPE,                              //!< type
  RTDB_TAG_INDEX_TABLE,                             //!< table
  RTDB_TAG_INDEX_DESC,                              //!< desc
  RTDB_TAG_INDEX_UNIT,                              //!< unit
  RTDB_TAG_INDEX_ARCHIVE,                           //!< archive
  RTDB_TAG_INDEX_DIGITS,                            //!< digits
  RTDB_TAG_INDEX_SHUTDOWN,                          //!< shutdown
  RTDB_TAG_INDEX_LOWLIMIT,                          //!< lowlimit
  RTDB_TAG_INDEX_HIGHLIMIT,                         //!< highlimit
  RTDB_TAG_INDEX_STEP,                              //!< step
  RTDB_TAG_INDEX_TYPICAL,                           //!< typical
  RTDB_TAG_INDEX_COMPRESS,                          //!< compress
  RTDB_TAG_INDEX_COMPDEV,                           //!< compdev
  RTDB_TAG_INDEX_COMPDEVPERCENT,                    //!< compdevpercent
  RTDB_TAG_INDEX_COMPTIMEMAX,                       //!< comptimemax
  RTDB_TAG_INDEX_COMPTIMEMIN,                       //!< comptimemin
  RTDB_TAG_INDEX_EXCDEV,                            //!< excdev
  RTDB_TAG_INDEX_EXCDEVPERCENT,                     //!< excdevpercent
  RTDB_TAG_INDEX_EXCTIMEMAX,                        //!< exctimemax
  RTDB_TAG_INDEX_EXCTIMEMIN,                        //!< exctimemin
  RTDB_TAG_INDEX_CLASSOF,                           //!< classof
  RTDB_TAG_INDEX_CHANGEDATE,                        //!< changedate
  RTDB_TAG_INDEX_CHANGER,                           //!< changer
  RTDB_TAG_INDEX_CREATEDATE,                        //!< createdate
  RTDB_TAG_INDEX_CREATOR,                           //!< creator
  RTDB_TAG_INDEX_MIRROR,                            //!< mirror
  RTDB_TAG_INDEX_MS,                                //!< precision
  RTDB_TAG_INDEX_FULLNAME,                          //!< table_dot_tag
  RTDB_TAG_INDEX_SUMMARY,                           //!< summary
  RTDB_TAG_INDEX_DATETIMEFORMAT,                    //!< datetimeformat
  RTDB_TAG_INDEX_BASE_LAST,

  RTDB_TAG_INDEX_SCAN_FIRST = 0x1000,
  RTDB_TAG_INDEX_SOURCE = RTDB_TAG_INDEX_SCAN_FIRST,    //!< source
  RTDB_TAG_INDEX_SCAN,                                  //!< scan
  RTDB_TAG_INDEX_INSTRUMENT,                            //!< instrument
  RTDB_TAG_INDEX_LOCATION1,                             //!< locations[0]
  RTDB_TAG_INDEX_LOCATION2,                             //!< locations[2]
  RTDB_TAG_INDEX_LOCATION3,                             //!< locations[3]
  RTDB_TAG_INDEX_LOCATION4,                             //!< locations[4]
  RTDB_TAG_INDEX_LOCATION5,                             //!< locations[5]
  RTDB_TAG_INDEX_USERINT1,                              //!< userints[0]
  RTDB_TAG_INDEX_USERINT2,                              //!< userints[1]
  RTDB_TAG_INDEX_USERREAL1,                             //!< userreals[0]
  RTDB_TAG_INDEX_USERREAL2,                             //!< userreals[1]
  RTDB_TAG_INDEX_SCAN_LAST,

  RTDB_TAG_INDEX_CALC_FIRST = 0x2000,
  RTDB_TAG_INDEX_EQUATION = RTDB_TAG_INDEX_CALC_FIRST,      //!< equation
  RTDB_TAG_INDEX_TRIGGER,                                   //!< trigger
  RTDB_TAG_INDEX_TIMECOPY,                                  //!< timecopy
  RTDB_TAG_INDEX_PERIOD,                                    //!< period
  RTDB_TAG_INDEX_CALC_LAST,

  RTDB_TAG_INDEX_SNAPSHOT_FIRST = 0x3000,
  RTDB_TAG_INDEX_TIMESTAMP = RTDB_TAG_INDEX_SNAPSHOT_FIRST,     //!< snapshot stamp (ms)
  RTDB_TAG_INDEX_VALUE,                                         //!< snapshot value
  RTDB_TAG_INDEX_QUALITY,                                       //!< snapshot quality
  RTDB_TAG_INDEX_SNAPSHOT_LAST,
} RTDB_TAG_FIELD_INDEX;

/**
* \ingroup denum
* \brief 标签点排序的标志
*/
typedef enum _RTDB_TAG_SORT_FLAG
{
  RTDB_SORT_FLAG_DESCEND = 0x0001,        //!< 降序
  RTDB_SORT_FLAG_CASE_SENSITIVE = 0x0002, //!< 大小写敏感
  RTDB_SORT_FLAG_RECYCLED = 0x0004,       //!< 用于回收站标签点排序
} RTDB_TAG_SORT_FLAG;

/**
* \ingroup denum
* \brief 性能计数点的ID
*/
typedef enum _RTDB_PERF_TAG_ID
{
  PFT_CPU_USAGE_OF_LOGGER,            //!< 日志服务CPU使用
  PFT_MEM_BYTES_OF_LOGGER,            //!< 日志服务内存
  PFT_VMEM_BYTES_OF_LOGGER,           //!< 日志服务虚拟内存
  PFT_READ_BYTES_OF_LOGGER,           //!< 日志服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_LOGGER,          //!< 日志服务 I/O 写入字节
  PFT_CPU_USAGE_OF_HISTORIAN,         //!< 历史数据服务CPU使用
  PFT_MEM_BYTES_OF_HISTORIAN,         //!< 历史数据服务内存
  PFT_VMEM_BYTES_OF_HISTORIAN,        //!< 历史数据服务虚拟内存
  PFT_READ_BYTES_OF_HISTORIAN,        //!< 历史数据服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_HISTORIAN,       //!< 历史数据服务 I/O 写入字节
  PFT_CPU_USAGE_OF_SNAPSHOT,          //!< 快照数据服务CPU使用
  PFT_MEM_BYTES_OF_SNAPSHOT,          //!< 快照数据服务内存
  PFT_VMEM_BYTES_OF_SNAPSHOT,         //!< 快照数据服务虚拟内存
  PFT_READ_BYTES_OF_SNAPSHOT,         //!< 快照数据服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_SNAPSHOT,        //!< 快照数据服务 I/O 写入字节
  PFT_CPU_USAGE_OF_EQUATION,          //!< 实时方程式服务CPU使用
  PFT_MEM_BYTES_OF_EQUATION,          //!< 实时方程式服务内存
  PFT_VMEM_BYTES_OF_EQUATION,         //!< 实时方程式服务虚拟内存
  PFT_READ_BYTES_OF_EQUATION,         //!< 实时方程式服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_EQUATION,        //!< 实时方程式服务 I/O 写入字节
  PFT_CPU_USAGE_OF_BASE,              //!< 标签点信息服务CPU使用
  PFT_MEM_BYTES_OF_BASE,              //!< 标签点信息服务内存
  PFT_VMEM_BYTES_OF_BASE,             //!< 标签点信息服务虚拟内存
  PFT_READ_BYTES_OF_BASE,             //!< 标签点信息服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_BASE,            //!< 标签点信息服务 I/O 写入字节
  PFT_CPU_USAGE_OF_SERVER,            //!< 网络服务CPU使用
  PFT_MEM_BYTES_OF_SERVER,            //!< 网络服务内存
  PFT_VMEM_BYTES_OF_SERVER,           //!< 网络服务虚拟内存
  PFT_READ_BYTES_OF_SERVER,           //!< 网络服务 I/O 读取字节
  PFT_WRITE_BYTES_OF_SERVER,          //!< 网络服务 I/O 写入字节
  PFT_ARV_ASYNC_QUEUE,                //!< 历史数据队列地址
  PFT_ARV_ASYNC_QUEUE_USAGE,          //!< 历史数据队列使用率
  PFT_ARVEX_ASYNC_QUEUE,              //!< 补历史数据队列地址
  PFT_ARVEX_ASYNC_QUEUE_USAGE,        //!< 补历史数据队列使用率
  PFT_EVENTS_INPUT_RATE,              //!< 普通事件入库速度（KB/秒）
  PFT_EVENTS_OUTPUT_RATE,             //!< 普通事件归档速度（KB/秒）
  PFT_FILL_IN_INPUT_RATE,             //!< 补历史事件入库速度（KB/秒）
  PFT_FILL_IN_OUTPUT_RATE,            //!< 补历史事件归档速度（KB/秒）
  PFT_ARV_CACHE_USAGE,                //!< 历史数据缓存使用率
  PFT_ARVEX_CACHE_USAGE,              //!< 补历史数据缓存使用率
  PFT_MIRROR_SNAPSHOTS_QUEUE,         //!< 快照数据的镜像队列地址
  PFT_MIRROR_SNAPSHOTS_QUEUE_USAGE,   //!< 快照数据的镜像队列使用率
  PFT_ARVEX_BLOB_ASYNC_QUEUE,         //!< str、blob补历史数据队列地址
  PFT_ARVEX_BLOB_ASYNC_QUEUE_USAGE,   //!< str、blob补历史数据队列使用率
  PFT_ARVEX_BLOB_CACHE_USAGE,         //!< str、blob补历史数据缓存使用率
  PFT_MIRROR_BUFFER_SIZE,             //!< 快照数据的镜像缓存文件
  PFT_CLUTTER_POOL_USAGE,             //!< 消息交换池利用率
  PFT_MAX_BLOCK_IN_CLUTTER_POOL,      //!< 消息交换池的最大可用额度
  PFT_ARV_ARCHIVED_TIME,              //!< 历史数据归档耗时
  PFT_ARVEX_ARCHIVED_TIME,            //!< 补历史数据归档耗时
  PFT_ARVEX_BLOB_ARCHIVED_TIME,       //!< str、blob补历史数据归档耗时
  PFT_ARV_ARCHIVED_PAGE_COUNT,        //!< 历史数据归档的数据页数量
  PFT_ARVEX_ARCHIVED_PAGE_COUNT,      //!< 补历史数据归档的数据页数量
  PFT_ARVEX_BLOB_ARCHIVED_PAGE_COUNT, //!< str、blob补历史数据归档的数据页数量
  PFT_MIRROR_ARV_VALUES_QUEUE,        //!< 补写历史数据的镜像队列地址
  PFT_MIRROR_ARV_VALUES_QUEUE_USAGE,  //!< 补写历史数据的镜像队列使用率
  PFT_MIRROR_ARV_BUFFER_SIZE,         //!< 补写历史数据的镜像缓存文件

  PFT_ARV_WRITE_COUNT,                //!< 历史数据归档写磁盘次数
  PFT_ARV_READ_COUNT,                 //!< 历史数据归档读磁盘次数
  PFT_ARV_WRITE_TIME,                 //!< 历史数据归档写磁盘时间
  PFT_ARV_READ_TIME,                  //!< 历史数据归档读磁盘时间
  PFT_ARV_INDEX_WRITE_COUNT,          //!< 历史数据归档写索引次数
  PFT_ARV_INDEX_READ_COUNT,           //!< 历史数据归档读索引次数
  PFT_ARV_INDEX_WRITE_TIME,           //!< 历史数据归档写索引时间
  PFT_ARV_INDEX_READ_TIME,            //!< 历史数据归档读索引时间
  PFT_ARV_ARC_LIST_LOCK_TIME,         //!< 历史数据归档列表锁时间
  PFT_ARV_ARC_LOCK_TIME,              //!< 历史数据归档文件锁时间
  PFT_ARV_INDEX_LOCK_TIME,            //!< 历史数据归档索引锁时间
  PFT_ARV_TOTAL_LOCK_TIME,            //!< 历史数据归档锁总时间
  PFT_ARV_WRITE_SIZE,                 //!< 历史数据归档写磁盘数据量
  PFT_ARV_READ_SIZE,                  //!< 历史数据归档读磁盘数据量
  PFT_ARV_WRITE_REAL_SIZE,            //!< 历史数据归档写磁盘有效数据量
  PFT_ARV_READ_REAL_SIZE,             //!< 历史数据归档读磁盘有效数据量
  PFT_ARVEX_WRITE_COUNT,              //!< 补历史数据归档写磁盘次数
  PFT_ARVEX_READ_COUNT,               //!< 补历史数据归档读磁盘次数
  PFT_ARVEX_WRITE_TIME,               //!< 补历史数据归档写磁盘时间
  PFT_ARVEX_READ_TIME,                //!< 补历史数据归档读磁盘时间
  PFT_ARVEX_INDEX_WRITE_COUNT,        //!< 补历史数据归档写索引次数
  PFT_ARVEX_INDEX_READ_COUNT,         //!< 补历史数据归档读索引次数
  PFT_ARVEX_INDEX_WRITE_TIME,         //!< 补历史数据归档写索引时间
  PFT_ARVEX_INDEX_READ_TIME,          //!< 补历史数据归档读索引时间
  PFT_ARVEX_ARC_LIST_LOCK_TIME,       //!< 补历史数据归档列表锁时间
  PFT_ARVEX_ARC_LOCK_TIME,            //!< 补历史数据归档文件锁时间
  PFT_ARVEX_INDEX_LOCK_TIME,          //!< 补历史数据归档索引锁时间
  PFT_ARVEX_TOTAL_LOCK_TIME,          //!< 补历史数据归档锁总时间
  PFT_ARVEX_WRITE_SIZE,               //!< 补历史数据归档写磁盘数据量
  PFT_ARVEX_READ_SIZE,                //!< 补历史数据归档读磁盘数据量
  PFT_ARVEX_WRITE_REAL_SIZE,          //!< 补历史数据归档写磁盘有效数据量
  PFT_ARVEX_READ_REAL_SIZE,           //!< 补历史数据归档读磁盘有效数据量

  PFT_PLOT_POOL_POINT_COUNT,          //!< 曲线缓存标签点数量
  PFT_PLOT_POOL_WEIGHTED_POINT_COUNT, //!< 曲线缓存权重点数量
  PFT_PLOT_POOL_TOTAL_MEM_SIZE,       //!< 曲线缓存总内存数
  PFT_PLOT_POOL_CACHED_HIT_PERCENT,   //!< 曲线缓存命中率

  PFT_OS_CPU_USAGE = 93,              //!< 数据库所在操作系统的CPU使用率
  PFT_OS_MEM_SIZE,                    //!< 数据库所在操作系统的物理内存大小，单位MB
  PFT_OS_MEM_USAGE,                   //!< 数据库所在操作系统的物理内存使用率

  PFT_QUERY_POOL_WAIT_TASKS_SIZE,     //!< 查询线程池中等待执行的任务数
  PFT_MIRROR_ENQUEUE,                 //!< 镜像每秒入队的数量，单位字节
  PFT_MIRROR_OUTQUEUE,                //!< 镜像每秒出对的数量，单位字节
  PFT_MIRROR_SEND_CPRS,               //!< 镜像每秒压缩的数量，单位字节
  PFT_MIRROR_RECV_CPRS,               //!< 镜像每秒收到的压缩数量，单位字节
  PFT_MIRROR_RECV_UNCPRS,             //!< 镜像每秒解压缩的数量，单位字节
  PFT_MIRROR_CPRS_SPAN,               //!< 镜像报文每秒的压缩耗时总和，单位毫秒
  PFT_MIRROR_COMPRESS_RATE,			  //!< 1秒内镜像的压缩率
  PFT_API_CPRS_RATE,                  //!< API报文压缩率
  PFT_SERVER_CPRS_RATE,               //!< Server报文压缩率

  PFT_TAG_SUBSCRIBE_CUSTOMER_COUNT,                 //!< 标签点信息订阅客户端数量
  PFT_TAG_SUBSCRIBE_SEND_EVENT_COUNT,               //!< 标签点信息订阅发送事件数量
  PFT_SNAP_SUBSCRIBE_CUSTOMER_COUNT,                //!< 快照信息订阅客户端数量
  PFT_SNAP_SUBSCRIBE_SEND_EVENT_COUNT,              //!< 快照信息订阅发送事件数量
  PFT_SNAP_SUBSCRIBE_POINT_COUNT,                   //!< 快照信息订阅标签点数量
  PFT_CONNECT_SUBSCRIBE_CUSTOMER_COUNT,             //!< API监视订阅客户端数量
  PFT_CONNECT_SUBSCRIBE_SEND_EVENT_COUNT,           //!< API监视订阅发送事件数量
  PFT_NAMED_TYPE_CREATE_SUBSCRIBE_CUSTOMER_COUNT,   //!< 创建自定义类型订阅客户端数量
  PFT_NAMED_TYPE_CREATE_SEND_EVENT_COUNT,           //!< 创建自定义类型订阅发送事件数量
  PFT_NAMED_TYPE_REMOVE_SUBSCRIBE_CUSTOMER_COUNT,   //!< 删除自定义类型订阅客户端数量
  PFT_NAMED_TYPE_REMOVE_SEND_EVENT_COUNT,           //!< 删除自定义类型订阅发送事件数量

  PFT_DOUBLE_ACTIVE_SYNC_SEND_COUNT,                //!< 双活同步每秒同步发送的数据量
  PFT_DOUBLE_ACTIVE_SYNC_RECEIVE_COUNT,             //!< 双活同步每秒同步接授的数据量
  PFT_IS_RECEIVING_NORMAL_DATA_FROM_PEER,           //!< 双活系统正在接收普通类型数据
  PFT_IS_RECEIVING_BLOB_DATA_FROM_PEER,             //!< 双活系统正在接收blob类型数据
  PFT_REPLICATOR_BUFFER_BLOCK_COUNT,                //!< 双活本地的同步历史缓存还有多少数据块
  PFT_REPLICATOR_EX_BUFFER_BLOCK_COUNT,             //!< 双活本地的同步补历史缓存还有多少数据块
  PFT_REPLICATOR_BLOB_BUFFER_BLOCK_COUNT,           //!< 双活本地的同步历史缓存(blob string 数据)还有多少数据块
  PFT_REPLICATOR_BLOB_EX_BUFFER_BLOCK_COUNT,        //!< 双活本地的同步补历史缓存(blob string 数据)还有多少数据块

  PFT_SNAPSHOT_PUT_RATE,                            //!< 每秒写入快照记录数，单位 条
  PFT_SNAPSHOT_GET_RATE,                            //!< 每秒读取快照记录数，单位 条
  PFT_HISTORIAN_PUT_RATE,                           //!< 每秒写入历史记录数，单位 条
  PFT_HISTORIAN_GET_RATE,                           //!< 每秒读取历史记录数，单位 条

  PFT_HISTORIAN_WRITE_RECORD_COUNT,                 //!< 每秒写入历史数据块数
  PFT_HISTORIAN_READ_RECORD_COUNT,                  //!< 每秒读取历史数据块数
  PFT_SERVER_NETWORK_READ_BYTES,                    //!< 网络服务网络 IO 每秒读取字节数
  PFT_SERVER_NETWORK_WRITE_BYTES,                   //!< 网络服务网络 IO 每秒写入字节数

  PFT_END,                            //!< 信息数量
} RTDB_PERF_TAG_ID;

/// 性能计数点的信息
typedef struct  _RTDB_PERF_TAG_INFO
{
  int perf_id;                      /// 性能计数点的ID 参考RTDB_PERF_TAG_ID
  char tag_name[RTDB_TAG_SIZE];     /// 性能计数点的名字
  char desc[RTDB_DESC_SIZE];        /// 性能计数点的描述
  char unit[RTDB_UNIT_SIZE];        /// 性能计数点数值的单位
  int type;                         /// 性能计数点的数值类型
} RTDB_PERF_TAG_INFO;

/**
* \ingroup ddatatype
* \typedef struct _RTDB_DATA_TYPE_FIELD RTDB_DATA_TYPE_FIELD
* \brief 组成自定义类型的字段定义
* \see _RTDB_DATA_TYPE_FIELD
*/

/**
* \ingroup dstruct
* \brief 组成自定义类型的字段定义.
*/
typedef struct _RTDB_DATA_TYPE_FIELD
{
  char         name[RTDB_TYPE_NAME_SIZE]; //!< 自定义类型的字段的名称，不要大于\b RTDB_TYPE_NAME_SIZE个字节
  rtdb_int32 type;                        //!< 字段的类型,只支持 \b RTDB_TYPE 里的类型，不支持struct，union等组合类型。
  rtdb_int32 length;                      //!< 字段类型的长度, \b RTDB_STRING、\b RTDB_BLOB等类型的具体的长度，基本类型本身的长度(基本类型可以忽略)，单位：字节。
  char         desc[RTDB_DESC_SIZE];      //!< 字段类型的描述，不要大于 \b RTDB_DESC_SIZE个字节
} RTDB_DATA_TYPE_FIELD, *PRTDB_DATA_TYPE_FIELD;

/**
* \ingroup denum
* \brief 将标签点属性加载到内存中的标志
*/
typedef enum _RTDB_TAG_LOAD_MEMORY_FLAG
{
  RTDB_LOAD_EMPTY_POINT = 0x0, //!< 什么也不加载

  RTDB_LOAD_TABLE_DOT_TAG = 0x1, //!< 点的全名称
  RTDB_LOAD_DESC = 0x2,          //!< 描述
  RTDB_LOAD_UNIT = 0x4,          //!< 单位
  RTDB_LOAD_CHANGER = 0x8,       //!< 修改人
  RTDB_LOAD_CREATOR = 0x10,      //!< 创建人
  RTDB_LOAD_LOWLIMIT = 0x20,     //!< 下限
  RTDB_LOAD_HIGHLIMIT = 0x40,    //!< 上限
  RTDB_LOAD_TYPICAL = 0x80,      //!< 典型值
  RTDB_LOAD_CHANGEDATE = 0x100,  //!< 修改日期
  RTDB_LOAD_CREATEDATE = 0x200,  //!< 创建日期
  RTDB_LOAD_DIGITS = 0x400,      //!< 数值位数

  RTDB_LOAD_COMPDEVPERCENT = 0x800, //!< 压缩偏差百分比
  RTDB_LOAD_EXCDEVPERCENT = 0x1000, //!< 例外偏差百分比

  RTDB_LOAD_SOURCE = 0x10000,      //!< 数据源
  RTDB_LOAD_SCAN = 0x20000,        //!< 是否采集
  RTDB_LOAD_INSTRUMENT = 0x40000,  //!< 设备标签
  RTDB_LOAD_LOCATION1 = 0x80000,   //!< 设备位址1
  RTDB_LOAD_LOCATION2 = 0x100000,  //!< 设备位址2
  RTDB_LOAD_LOCATION3 = 0x200000,  //!< 设备位址3
  RTDB_LOAD_LOCATION4 = 0x400000,  //!< 设备位址4
  RTDB_LOAD_LOCATION5 = 0x800000,  //!< 设备位址5
  RTDB_LOAD_USERINT1 = 0x1000000,  //!< 自定义整数1
  RTDB_LOAD_USERINT2 = 0x2000000,  //!< 自定义整数2
  RTDB_LOAD_USERREAL1 = 0x4000000, //!< 自定义浮点数1
  RTDB_LOAD_USERREAL2 = 0x8000000, //!< 自定义浮点数2

  RTDB_LOAD_BASE_POINT = RTDB_LOAD_TABLE_DOT_TAG | RTDB_LOAD_DESC | RTDB_LOAD_UNIT | RTDB_LOAD_CHANGER | RTDB_LOAD_CREATOR | RTDB_LOAD_LOWLIMIT | RTDB_LOAD_HIGHLIMIT | RTDB_LOAD_TYPICAL | RTDB_LOAD_CHANGEDATE | RTDB_LOAD_CREATEDATE | RTDB_LOAD_DIGITS | RTDB_LOAD_COMPDEVPERCENT | RTDB_LOAD_EXCDEVPERCENT,                          //!< base 属性合集
  RTDB_LOAD_SCAN_POINT = RTDB_LOAD_SOURCE | RTDB_LOAD_SCAN | RTDB_LOAD_INSTRUMENT | RTDB_LOAD_LOCATION1 | RTDB_LOAD_LOCATION2 | RTDB_LOAD_LOCATION3 | RTDB_LOAD_LOCATION4 | RTDB_LOAD_LOCATION5 | RTDB_LOAD_USERINT1 | RTDB_LOAD_USERINT2 | RTDB_LOAD_USERREAL1 | RTDB_LOAD_USERREAL2, //!< scan 属性合集
  RTDB_LOAD_ALL_POINT = RTDB_LOAD_BASE_POINT | RTDB_LOAD_SCAN_POINT,  //!< 所有属性合集
} RTDB_TAG_LOAD_MEMORY_FLAG;

#define RTDB_GET_FROM_FLAG(FLAG, BIT) (((FLAG) & (BIT)) ? 1 : 0)
#define RTDB_SET_FROM_FLAG(FLAG, BIT, VALUE) {if (VALUE) (FLAG) |= (BIT); else (FLAG) &= (~(BIT));}

#define RTDB_HAS_EMPTY_FROM_FLAG(FLAG) (((FLAG) & (RTDB_LOAD_EMPTY_POINT)) ? 1 : 0)
#define RTDB_HAS_BASE_FROM_FLAG(FLAG) (((FLAG) & (RTDB_LOAD_BASE_POINT)) ? 1 : 0)
#define RTDB_HAS_SCAN_FROM_FLAG(FLAG) (((FLAG) & (RTDB_LOAD_SCAN_POINT)) ? 1 : 0)
#define RTDB_HAS_ALL_FROM_FLAG(FLAG) (((FLAG) & (RTDB_LOAD_ALL_POINT)) ? 1 : 0)

/**
* \ingroup denum
* \brief 数据归档策略
*/
typedef enum _RTDB_ARCHIVED_POLICY
{
  RTDB_ARCHIVED_SNAPSHOT_FIRST,           //!< 快照数据优先归档
  RTDB_ARCHIVED_ARCHIVEX_FIRST,           //!< 补写数据优先归档
  RTDB_ARCHIVED_AUTO,                     //!< 自动判断快照数据和补写数据的优先级
  RTDB_ARCHIVED_PAUSE,                    //!< 暂停归档
} RTDB_ARCHIVED_POLICY;

/**
* \ingroup denum
* \brief API类别
*/
typedef enum _API_CATEGORY
{
  API_SERVER,    //!< 网络服务API
  API_BASE,      //!< 标签点服务API
  API_SNAPSHOT,  //!< 快照服务API
  API_HISTORIAN, //!< 历史服务API
  API_ARCHIVE,   //!< 存档文件API
  API_EQUATION,  //!< 方程式服务API
  API_LOGGER,    //!< 日志服务API
  API_PERF,      //!< 性能计数服务API
  API_DISPATCH,  //!< 转发服务API
  API_MEMORYDB,  //!< 内存库服务API
} API_CATEGORY;

/**
* \ingroup ddatatype
* \typedef struct _RTDB_CONNECT_EVENT RTDB_CONNECT_EVENT
* \brief 每个数据库连接的调用信息
* \see _RTDB_CONNECT_EVENT
*/

/**
* \ingroup dstruct
* \brief 每个数据库连接的调用信息.
*/
typedef struct _RTDB_CONNECT_EVENT
{
  rtdb_int32 msg_id;                      //!< 调用的方法ID
  rtdb_datetime_type begin_s;             //!< 调用的开始时间
  rtdb_int16 begin_ms;                    //!< 调用的开始时间的毫秒，恢复为毫秒定义，以适配其他GEM
  rtdb_int16 api_category;                //!< 调用的API类别
  rtdb_uint32 client_addr;                //!< 客户端的IP地址
  rtdb_int32 client_process_id;           //!< 客户端的进程ID
  rtdb_int32 client_thread_id;            //!< 客户端的线程ID
  rtdb_int32 server_thread_id;            //!< 服务端的线程ID
  rtdb_error ret_val;                     //!< 调用的方法返回值
  rtdb_float32 elapsed;                   //!< 调用的服务端耗时 单位毫秒
  rtdb_int32 pre_count;                   //!< 调用时传入的count
  rtdb_int32 post_count;                  //!< 调用时传出的count

  rtdb_uint32 write_count;                //!< 历史服务查询的写磁盘次数
  rtdb_uint32 read_count;                 //!< 历史服务查询的读磁盘次数
  rtdb_float32 write_time;                //!< 历史服务查询的写磁盘时间
  rtdb_float32 read_time;                 //!< 历史服务查询的读磁盘时间
  rtdb_uint32 index_write_count;          //!< 历史服务查询的写索引次数
  rtdb_uint32 index_read_count;           //!< 历史服务查询的读索引次数
  rtdb_float32 index_write_time;          //!< 历史服务查询的写索引时间
  rtdb_float32 index_read_time;           //!< 历史服务查询的读索引时间
  rtdb_float32 arc_list_lock_time;        //!< 历史服务查询的文件列表锁时间
  rtdb_float32 arc_lock_time;             //!< 历史服务查询的存档文件锁时间
  rtdb_float32 index_lock_time;           //!< 历史服务查询的索引锁时间
  rtdb_float32 total_lock_time;           //!< 历史服务查询的锁的总时间
  rtdb_float32 write_size;                //!< 历史服务查询的写入数据量，单位KB
  rtdb_float32 read_size;                 //!< 历史服务查询的读取数据量，单位KB
  rtdb_float32 write_real_size;           //!< 历史服务查询的写入有效数据量，单位KB
  rtdb_float32 read_real_size;            //!< 历史服务查询的读取有效数据量，单位KB

  rtdb_byte client_addr6[RTDB_IPV6_ADDR_SIZE]; //!< 客户端的ipv6地址
  char reserve_1[4];                        //!< 保留字节
}RTDB_CONNECT_EVENT;    //128字节

/**
* \ingroup ddatatype
* \typedef struct _RTDB_ARCHIVE_PERF_DATA RTDB_ARCHIVE_PERF_DATA
* \brief 每个存档文件的性能信息，重启服务后会重新统计
* \see _RTDB_ARCHIVE_PERF_DATA
*/

/**
* \ingroup dstruct
* \brief 每个存档文件的性能信息，重启服务后会重新统计
*/
typedef struct _RTDB_ARCHIVE_PERF_DATA
{
  rtdb_uint32 write_count;                //!< 历史服务查询的写磁盘次数
  rtdb_uint32 read_count;                 //!< 历史服务查询的读磁盘次数
  rtdb_float32 write_time;                //!< 历史服务查询的写磁盘时间
  rtdb_float32 read_time;                 //!< 历史服务查询的读磁盘时间
  rtdb_uint32 index_write_count;          //!< 历史服务查询的写索引次数
  rtdb_uint32 index_read_count;           //!< 历史服务查询的读索引次数
  rtdb_float32 index_write_time;          //!< 历史服务查询的写索引时间
  rtdb_float32 index_read_time;           //!< 历史服务查询的读索引时间
  rtdb_float32 arc_list_lock_time;        //!< 历史服务查询的文件列表锁时间
  rtdb_float32 arc_lock_time;             //!< 历史服务查询的存档文件锁时间
  rtdb_float32 index_lock_time;           //!< 历史服务查询的索引锁时间
  rtdb_float32 total_lock_time;           //!< 历史服务查询的锁的总时间
  rtdb_float32 write_size;                //!< 历史服务查询的写入数据量，单位KB
  rtdb_float32 read_size;                 //!< 历史服务查询的读取数据量，单位KB
  rtdb_float32 write_real_size;           //!< 历史服务查询的写入有效数据量，单位KB
  rtdb_float32 read_real_size;            //!< 历史服务查询的读取有效数据量，单位KB
}RTDB_ARCHIVE_PERF_DATA;      //64字节

/**
* \ingroup dcallback
* \brief 标签点属性更改通知订阅回调接口
* \param count     整型，输入，标签点个数，ids 的长度，
* \param ids       整型数组，输入，标签点被订阅且属性发生更改的标识列表
* \param what      整型，参考枚举 \b RTDB_TAG_CHANGE_REASON,
*                    表示引起变更的源类型。
* \see rtdbb_subscribe_tags
*/
typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbb_tags_change_event)(
  rtdb_int32 count,
  const rtdb_int32 *ids,
  rtdb_int32 what
  );

/**
* \ingroup dcallback
* \brief 标签点属性更改通知订阅回调接口(ex)
*/
typedef rtdbb_tags_change_event rtdbb_tags_change;

/**
* \ingroup dcallback
* \brief 标签点属性更改通知订阅回调接口(ex)
* \param event_type		无符号整型，输入，通知类型
* \param handle			整型，输入，产生通知的句柄
* \param param			void指针，输入，用户调用rtdbb_subscribe_tags_ex时param参数
* \param count		    整型，输入，event_type为RTDB_E_DATA时表示标签点个数，否则count为0
* \param ids			整型数组，输入，标签点被订阅且属性发生更改的标识列表
* \param what			整型，参考枚举 RTDB_TAG_CHANGE_REASON,
*		                表示引起变更的源类型。
* \see RTDB_EVENT_TYPE
*/
typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbb_tags_change_event_ex)(
  rtdb_uint32 event_type,
  rtdb_int32 handle,
  void* param,
  rtdb_int32 count,
  const rtdb_int32 *ids,
  rtdb_int32 what
  );

/**
* \ingroup dcallback
* \brief 标签点快照改变通知订阅回调接口
* \param count     整型，输入，
*                    表示 ids、datetimes、ms、values、states、qualities、errors的长度
* \param ids       整型数组，输入，标签点被订阅且快照发生改变的标识列表
* \param datatimes 整型数组，输入，实时数值时间列表，
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则为 0
* \param values    双精度浮点型数组，输入，实时浮点型数值列表
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则为 0
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则为 0
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbs_snaps_event)(
  rtdb_int32 count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float64 *values,
  const rtdb_int64 *states,
  const rtdb_int16 *qualities,
  const rtdb_error *errors
  );

/**
* \ingroup dcallback
* \brief 标签点快照改变通知订阅回调接口
*/
typedef rtdbs_snaps_event rtdbs_data_change;

/**
* \ingroup dcallback
* \brief 标签点快照改变通知订阅回调接口(ex)
* \param event_type  无符号整形，输入，通知类型
* \param handle      整型，输入，产生通知的句柄
* \param param       void指针，输入，用户调用rtdbs_subscribe_snapshots_ex时param参数
* \param count       整型，输入，
					   表示 ids、datetimes、ms、values、states、qualities、errors的长度
* \param ids         整型数组，输入，标签点被订阅且快照发生改变的标识列表
* \param datatimes   整型数组，输入，实时数值时间列表，

*                      表示距离1970年1月1日08:00:00的秒数
* \param ms          短整型数组，输入，实时数值时间列表，
*                     对于时间精度为纳秒的标签点，存放相应的纳秒值；否则为 0
* \param values      双精度浮点型数组，输入，实时浮点型数值列表

*                      对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则为 0
* \param states      64 位整型数组，输入，实时整型数值列表，
*                      对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                      RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则为 0
* \param qualities   短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors      无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
* \see RTDB_EVENT_TYPE

*/
typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbs_snaps_event_ex)(
  rtdb_uint32 event_type,
  rtdb_int32 handle,
  void* param,
  rtdb_int32 count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float64 *values,
  const rtdb_int64 *states,
  const rtdb_int16 *qualities,
  const rtdb_error *errors
  );

typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbs_snaps_event_ex64)(
    rtdb_uint32 event_type,
    rtdb_int32 handle,
    void* param,
    rtdb_int32 count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    const rtdb_error* errors
    );


/**
* 命名：rtdbh_data_playback
* 功能：标签点历史数据回放回调接口
* 参数：
*        [count]     整型，输入，
*                    表示 ids、datetimes、ms、values、states、qualities、errors的长度
*        [ids]       整型数组，输入，到达数据的标识列表
*        [datatimes] 整型数组，输入，到达数值时间列表，
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输入，到达数值时间列表，
*                    对于时间精度为毫秒的标签点，存放相应的毫秒值。
*        [values]    双精度浮点型数组，输入，到达的浮点型数值列表
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应数据。
*        [states]    64 位整型数组，输入，到达的整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应数据。
*        [x]         32 位浮点数，输入，二维坐标的x值，对于数据类型为 RTDB_COOR 的标签点，存放相应数据。
*        [y]         32 位浮点数，输入，二维坐标的y值，对于数据类型为 RTDB_COOR 的标签点，存放相应数据。
*        [qualities] 短整型数组，输入，到达的数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [error]     无符号整型，输出，如返回RtE_DATA_PLAYBACK_DONE则表明是最后一次回放，否则只会返回RtE_OK。
* 备注：本接口对数据类型为 RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbh_data_playback)(
                     rtdb_int32 count,
                     const rtdb_int32 *ids,
                     const rtdb_int32 *datetimes,
                     const rtdb_int16 *ms,
                     const rtdb_float64 *values,
                     const rtdb_int64 *states,
                     const rtdb_float32 *x,
                     const rtdb_float32 *y,
                     const rtdb_int16 *qualities,
                     const rtdb_error *error
                     );
/**
* 命名：rtdbh_data_playback
* 功能：标签点历史数据回放回调接口
* 参数：
*        [count]     整型，输入，
*                    表示 ids、datetimes、ms、values、states、qualities、errors的长度
*        [ids]       整型数组，输入，到达数据的标识列表
*        [datatimes] 整型数组，输入，到达数值时间列表，
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输入，到达数值时间列表，
*                    对于时间精度为毫秒的标签点，存放相应的毫秒值。
*        [values]    双精度浮点型数组，输入，到达的浮点型数值列表
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应数据。
*        [states]    64 位整型数组，输入，到达的整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应数据。
*        [x]         32 位浮点数，输入，二维坐标的x值，对于数据类型为 RTDB_COOR 的标签点，存放相应数据。
*        [y]         32 位浮点数，输入，二维坐标的y值，对于数据类型为 RTDB_COOR 的标签点，存放相应数据。
*        [qualities] 短整型数组，输入，到达的数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [error]     无符号整型，输出，如返回RtE_DATA_PLAYBACK_DONE则表明是最后一次回放，否则只会返回RtE_OK。
* 备注：本接口对数据类型为 RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbh_data_playback_ex)(
    rtdb_uint32 event_type,
    rtdb_int32 handle,
    void* param,
    rtdb_int32 count,
    const rtdb_int32* ids,
    const rtdb_int32* datetimes,
    const rtdb_int16* ms,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_float32* x,
    const rtdb_float32* y,
    const rtdb_int16* qualities,
    const rtdb_error* error
);

typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbb_named_type_remove_event_ex)(
  rtdb_uint32 event_type,
  rtdb_int32 handle,
  void* param,
  const char* name
  );

typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdbb_named_type_create_event_ex)(
  rtdb_uint32 event_type,
  rtdb_int32 handle,
  void* param,
  const char* name,
  long field_count,
  const RTDB_DATA_TYPE_FIELD* fields
  );

/**
* \ingroup dcallback
* \brief 数据库连接调用信息通知订阅回调接口
* \param count        连接调用events的个数
* \param events       连接调用信息
* \param pre_calls    连接调用时传入的参数信息
* \param post_calls   连接调用后传出的参数信息
*/
typedef
rtdb_error
(RTDBAPI_CALLRULE *rtdb_connect_event)(
  rtdb_int32 count,
  const RTDB_CONNECT_EVENT* const* events,
  const char* const* pre_calls,
  const char* const* post_calls
  );

/**
* \ingroup dcallback
* \brief 数据库连接调用信息通知订阅回调接口
* \param event_type   无符号整形，输入，通知类型
* \param handle       整型，输入，产生通知的句柄
* \param param        void指针，输入，用户调用rtdb_subscribe_connect_ex时param参数
* \param count        连接调用events的个数
* \param events       连接调用信息
* \param pre_calls    连接调用时传入的参数信息
* \param post_calls   连接调用后传出的参数信息
* \see RTDB_EVENT_TYPE
*/
typedef rtdb_error(RTDBAPI_CALLRULE *rtdb_connect_event_ex)(
  rtdb_uint32 event_type,
  rtdb_int32 handle,
  void* param,
  rtdb_int32 count,
  const RTDB_CONNECT_EVENT* const* events,
  const char* const* pre_calls,
  const char* const* post_calls
  );

/**
* \ingroup denum
* \brief 元数据同步角色
*/
typedef enum _RTDB_SYNC_ROLE
{
    RTDB_SYNC_ROLE_OFFLINE = 0,           //!< 离线
    RTDB_SYNC_ROLE_UNSYNCED = 1,          //!< 未同步
    RTDB_SYNC_ROLE_SYNCING = 2,           //!< 同步中
    RTDB_SYNC_ROLE_SLAVE = 3,             //!< 备库
    RTDB_SYNC_ROLE_MASTER = 4             //!< 主库
} RTDB_SYNC_ROLE;

/**
* \ingroup denum
* \brief 元数据同步状态
*/
typedef enum _RTDB_SYNC_STATUS
{
    RTDB_SYNC_STATUS_INIT = 0,            //!< 正常
    RTDB_SYNC_STATUS_START = 1,           //!< 启动同步
    RTDB_SYNC_STATUS_FILE = 2,            //!< 同步文件
    RTDB_SYNC_STATUS_CACHE = 3            //!< 同步缓存
} RTDB_SYNC_STATUS;

/**
* \ingroup dstruct
* \brief 节点的元数据同步信息
*/
typedef struct _RTDB_SYNC_INFO
{
    rtdb_int8 role;                 //!< 当前同步角色，参考 RTDB_SYNC_ROLE
    rtdb_int8 status;               //!< 当前同步状态，参考 RTDB_SYNC_STATUS
    char reserved[2];               //!< 保留字段
    rtdb_uint32 ip;                 //!< 当前节点ip
    rtdb_uint64 version;            //!< 当前节点的同步版本号
    rtdb_uint64 data_size;          //!< 同步数据堆积的数据量，单位字节
    char reserved2[40];             //!< 保留字段
} RTDB_SYNC_INFO; //64byte

/**
* \ingroup denum
* \brief 元数据同步状态
*/
typedef enum _RTDB_SUBSCRIBE_CHANGE_TYPE
{
    RTDB_SUBSCRIBE_ADD,           //!< 增加订阅
    RTDB_SUBSCRIBE_UPDATE,        //!< 更新订阅信息
    RTDB_SUBSCRIBE_REMOVE,        //!< 移除订阅
} RTDB_SUBSCRIBE_CHANGE_TYPE;

typedef struct _RTDB_SUMMARY_DATA
{
    rtdb_timestamp_type first_time;
    rtdb_timestamp_type last_time;
    rtdb_timestamp_type max_time;
    rtdb_timestamp_type min_time;
    rtdb_subtime_type first_subtime;
    rtdb_subtime_type last_subtime;
    rtdb_subtime_type max_subtime;
    rtdb_subtime_type min_subtime;
    rtdb_float64 first_value;
    rtdb_float64 last_value;
    rtdb_float64 max_value;
    rtdb_float64 min_value;
    rtdb_int16 first_quality;
    rtdb_int16 last_quality;
    rtdb_int16 max_quality;
    rtdb_int16 min_quality;
    rtdb_float64 power;
    rtdb_float64 power_avg;
    rtdb_float64 total;
    rtdb_float64 calc_avg;
    rtdb_int32 count;
    rtdb_int32 valid_count;

    char reserved[128];
} RTDB_SUMMARY_DATA; //256 byte


/**
* \ingroup denum
* \brief 时间戳精度
*/
typedef enum _RTDB_TIME_PRECISION_TYPE
{
    RTP_SECOND,     //!< 秒
    RTP_MILLI,      //!< 毫秒
    RTP_MICRO,      //!< 微秒
    RTP_NANO,       //!< 纳秒
} RTDB_TIME_PRECISION_TYPE;

typedef struct _RTDB_HANDLE_INFO
{
    rtdb_int8 os_type;      //!< 当前连接数据库的系统，参考 RTDB_OS_TYPE
    rtdb_int8 new_db;       //!< 当前连接数据库的版本，0表示旧版本，1表示新版本
    char reserved[62];
} RTDB_HANDLE_INFO; //64 byte

#ifdef WIN32
#ifdef _WIN64
#include <math.h>
double inline _rtdb_sqrt(double d)
{
  return sqrt(d);
}
#else
double inline __declspec (naked) __fastcall _rtdb_sqrt(double)
{
  _asm fld qword ptr [esp+4]
  _asm fsqrt
  _asm ret 8
}
#endif
#endif

#define RTDB_ABS(V1, V2)                          ( (V1) > (V2) ? ((V1) - (V2)) : ((V2) - (V1)) )
#define RTDB_MAX(V1, V2)                          ( (V1) > (V2) ? (V1) : (V2) )
#define RTDB_MIN(V1, V2)                          ( (V1) < (V2) ? (V1) : (V2) )
#define RTDB_TIME_LESS_THAN(S1, MS1, S2, MS2)     ( S1 < S2 ? true : ( S1 == S2 ? ( MS1 < MS2 ? true : false ) : false ) )
#define RTDB_TIME_EQUAL_TO(S1, MS1, S2, MS2)      ( S1 == S2 ? ( MS1 == MS2 ? true : false ) : false )
#define RTDB_TIME_GREATER_THAN(S1, MS1, S2, MS2)  ( S1 > S2 ? true : ( S1 == S2 ? ( MS1 > MS2 ? true : false ) : false ) )
#define RTDB_TIME_EQUAL_OR_LESS_THAN(S1, MS1, S2, MS2)     ( S1 < S2 ? true : ( S1 == S2 ? ( MS1 <= MS2 ? true : false ) : false ) )
#define RTDB_TIME_EQUAL_OR_GREATER_THAN(S1, MS1, S2, MS2)  ( S1 > S2 ? true : ( S1 == S2 ? ( MS1 >= MS2 ? true : false ) : false ) )
#define RTDB_MS_DELAY_BETWEEN(S1, MS1, S2, MS2)   ((unsigned long long) ((unsigned long long)(S1 - S2) * RTDB_MS_PRECISION + (unsigned long long)(MS1 - MS2)))
//以下两个时间操作宏，需要确保标签点类型为毫/纳秒点，如为秒点请勿使用下述宏进行计算
#define RTDB_MS_ADD_TIME(S1, MS1, MSES, S2, MS2)  { \
                                                      unsigned long long mses__ = MSES + MS1; \
                                                      S2 = S1 + static_cast<rtdb_int64>(mses__ / RTDB_MS_PRECISION); \
                                                      MS2 = static_cast<rtdb_time_type>(mses__ % RTDB_MS_PRECISION); }
#define RTDB_MS_SUB_TIME(S1, MS1, MSES, S2, MS2)  { \
                                                      rtdb_int64 subs__  = static_cast<rtdb_int64>(MSES / RTDB_MS_PRECISION);  \
                                                      rtdb_time_type subms__ = static_cast<rtdb_time_type>(MSES % RTDB_MS_PRECISION);\
                                                      S2 = (MS1 < subms__ ? S1 - subs__ - 1 : S1 - subs__); \
                                                      MS2 = (MS1 < subms__ ? MS1 + RTDB_MS_PRECISION - subms__ : MS1 - subms__);}

#ifdef WIN32
#define RTDB_DISTANCE(X1, Y1, X2, Y2)             _rtdb_sqrt( RTDB_ABS(X1, X2) * RTDB_ABS(X1, X2) + RTDB_ABS(Y1, Y2) * RTDB_ABS(Y1, Y2) )
#else
#define RTDB_DISTANCE(X1, Y1, X2, Y2)             sqrt( RTDB_ABS(X1, X2) * RTDB_ABS(X1, X2) + RTDB_ABS(Y1, Y2) * RTDB_ABS(Y1, Y2) )
#endif


/**
* \ingroup dmacro
* \def RTDB_API_MAJOR_VERSION
* \brief API的主版本号
*/
#define RTDB_API_MAJOR_VERSION      4
/**
* \ingroup dmacro
* \def RTDB_API_MINOR_VERSION
* \brief API的副版本号
*/
#define RTDB_API_MINOR_VERSION      0
/**
* \ingroup dmacro
* \def RTDB_API_BETA_VERSION
* \brief API的发行版本号
*/
#define RTDB_API_BETA_VERSION       11

#endif /* __RTDB_H__ */
