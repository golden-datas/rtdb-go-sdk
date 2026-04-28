#ifndef __RTDB_API_H__
#define __RTDB_API_H__

#include "rtdb.h"
#include "rtdb_error.h"

/**
* \defgroup config API 配置接口
* @{
*/

/**
* \brief   取得 rtdbapi 库的版本号
* \param [out]  major   主版本号
* \param [out]  minor   次版本号
* \param [out]  beta    发布版本号
* \return rtdb_error
* \remark 如果返回的版本号与 rtdb.h 中定义的不匹配(RTDB_API_XXX_VERSION)，则应用程序使用了错误的库。
*      应输出一条错误信息并退出，否则可能在调用某些 api 时会导致崩溃
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_api_version(
  rtdb_int32 *major,
  rtdb_int32 *minor,
  rtdb_int32 *beta
  );


/**
* \brief 配置 api 行为参数，参见枚举 \ref RTDB_API_OPTION
* \param [in] type  选项类别
* \param [in] value 选项值
* \return rtdb_error
* \remark 选项设置后在下一次调用 api 时才生效
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_set_option(
  rtdb_int32 type,
  rtdb_int32 value
  );

/**@}*/


/*备注： rtdb_datagram_handle* handle , 本身是 rtdb_datagram_handle& handle , go无法编译修改成 *指针 的*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_create_datagram_handle(
  rtdb_int32 port,
  const char* remotehost,
  rtdb_datagram_handle* handle
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_remove_datagram_handle(
  rtdb_datagram_handle handle
  );

/*备注： rtdb_int32* message_len 本身是  rtdb_int32& message_len, go无法编译因此修改成 *指针 */
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_recv_datagram(
  char* message,
  rtdb_int32* message_len,
  rtdb_datagram_handle handle,
  char* remote_addr,
  rtdb_int32 timeout GAPI_DEFAULT_VALUE(-1)
  );


/**
* \defgroup server 网络服务接口
* @{
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_subscribe_connect_ex(
  rtdb_int32 handle,
  rtdb_uint32 options,
  void* param,
  rtdb_connect_event_ex callback
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_cancel_subscribe_connect(
  rtdb_int32 handle
  );

/**
* \brief 建立同 RTDB 数据库的网络连接
* \param [in] hostname     RTDB 数据平台服务器的网络地址或机器名
* \param [in] port         连接断开，缺省值 6327
* \param [out]  handle  连接句柄
* \return rtdb_error
* \remark 在调用所有的接口函数之前，必须先调用本函数建立同Rtdb服务器的连接
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_connect(
  const char *hostname,
  rtdb_int32 port,
  rtdb_int32 *handle
  );

/**
* \brief 获取 RTDB 服务器当前连接个数
* \param [in] handle   连接句柄 参见 \ref rtdb_connect
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [out]  count 返回当前主机的连接个数
* \return rtdb_error
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_connection_count(
  rtdb_int32 handle,
  rtdb_int32 node_number,
  rtdb_int32 *count
  );

/**
* \brief 列出 RTDB 服务器的所有连接句柄
* \param [in] handle       连接句柄
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [out] sockets    整形数组，所有连接的套接字句柄
* \param [in,out]  count   输入时表示sockets的长度，输出时表示返回的连接个数
* \return rtdb_error
* \remark 用户须保证分配给 sockets 的空间与 count 相符。如果输入的 count 小于输出的 count，则只返回部分连接
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_connections(
  rtdb_int32 handle,
  rtdb_int32 node_number,
  rtdb_int32 *sockets,
  rtdb_int32 *count
  );
/**
* 命名：rtdb_get_own_connection
* 功能：获取当前连接的socket句柄
* 参数：
* \param [in] handle       连接句柄
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [out] sockets    整形数组，所有连接的套接字句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_own_connection(
	rtdb_int32 handle,
	rtdb_int32 node_number,
	rtdb_int32* socket
	);

/**
* \brief 获取 RTDB 服务器指定连接的信息
* \param [in] handle          连接句柄，参见 \ref rtdb_connect
* \param [in] node_number   双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP
* \param [in] socket          指定的连接
* \param [out] info          与连接相关的信息，参见 \ref RTDB_HOST_CONNECT_INFO
* \return rtdb_error
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_connection_info(
  rtdb_int32 handle,
  rtdb_int32 node_number,
  rtdb_int32 socket,
  RTDB_HOST_CONNECT_INFO *info
  );

/**
* \brief 获取 RTDB 服务器指定连接的ipv6版本
* \param [in] handle          连接句柄，参见 \ref rtdb_connect
* \param [in] node_number     双活模式下，指定节点编号，1为rtdb_connect中第1个IP，2为rtdb_connect中第2个IP，双活模式仅支持ipv4
* \param [in] socket          指定的连接
* \param [out] info           与连接相关的信息，参见 \ref RTDB_HOST_CONNECT_INFO_IPV6
* \return rtdb_error
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_connection_info_ipv6(
  rtdb_int32 handle,
  rtdb_int32 node_number,
  rtdb_int32 socket,
  RTDB_HOST_CONNECT_INFO_IPV6* info
);

/**
* \brief 断开同 RTDB 数据平台的连接
* \param handle  连接句柄
* \return rtdb_error
* \remark 完成对 RTDB 的访问后调用本函数断开连接。连接一旦断开，则需要重新连接后才能调用其他的接口函数。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_disconnect(
  rtdb_int32 handle
  );

/**
* \brief 以有效帐户登录
* \param handle          连接句柄
* \param user            登录帐户
* \param password        帐户口令
* \param [out] priv     账户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_login(
  rtdb_int32 handle,
  const char *user,
  const char *password,
  rtdb_int32 *priv
  );

/**
* \brief 获取连接句柄所连接的服务器操作系统类型
* \param     handle          连接句柄
* \param     ostype   操作系统类型 枚举 \ref RTDB_OS_TYPE 的值之一
* \return    rtdb_error
* \remark 如句柄未链接任何服务器，返回RTDB_OS_INVALID(当前支持操作系统类型：windows、linux)。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_linked_ostype(
  rtdb_int32 handle,
  RTDB_OS_TYPE* ostype
  );

/**
* \brief 获取连接句柄所连接的服务器相关信息
* \param     handle          连接句柄
* \param     handle_info   服务器相关信息
* \return    rtdb_error
* \remark 如句柄未链接任何服务器，返回RTDB_OS_INVALID(当前支持操作系统类型：windows、linux)。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_handle_info(
    rtdb_int32 handle,
    RTDB_HANDLE_INFO* info
);

/**
* \brief 修改用户帐户口令
* \param handle    连接句柄
* \param user      已有帐户
* \param password  帐户新口令
* \return rtdb_error
* \remark 只有系统管理员可以修改其它用户的密码
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_change_password(
  rtdb_int32 handle,
  const char *user,
  const char *password
  );

/**
* \brief 用户修改自己帐户口令
* \param handle  连接句柄
* \param old_pwd 帐户原口令
* \param new_pwd 帐户新口令
* \return rtdb_error
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_change_my_password(
  rtdb_int32 handle,
  const char *old_pwd,
  const char *new_pwd
  );


/**
* \brief 获取连接权限
* \param handle          连接句柄
* \param [out] priv  帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 如果还未登陆或不在服务器信任连接中，对应权限为-1，表示没有任何权限
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_priv(
  rtdb_int32 handle,
  rtdb_int32 *priv
  );


/**
* \brief 修改用户帐户权限
* \param handle  连接句柄
* \param user    已有帐户
* \param priv    帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 只有管理员有修改权限
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_change_priv(
  rtdb_int32 handle,
  const char *user,
  rtdb_int32 priv
  );

/**
* \brief 添加用户帐户
* \param handle    连接句柄
* \param user      帐户
* \param password  帐户初始口令
* \param priv      帐户权限， 枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 只有管理员有添加用户权限
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_add_user(
  rtdb_int32 handle,
  const char *user,
  const char *password,
  rtdb_int32 priv
  );

/**
* \brief 删除用户帐户
* \param handle  连接句柄
* \param user    帐户
* \return rtdb_error
* \remark 只有管理员有删除用户权限
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_remove_user(
  rtdb_int32 handle,
  const char *user
  );

/**
* \brief 启用或禁用用户
* \param     handle    连接句柄
* \param     user      字符串，输入，帐户名
* \param     lock      布尔，输入，是否禁用
* \return    rtdb_error
* \remark 只有管理员有启用禁用权限
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_lock_user(
  rtdb_int32 handle,
  const char *user,
  rtdb_int8 lock
  );


/**
* \brief 获得所有用户
* \param handle          连接句柄
* \param [in,out]  count 输入时表示 users、privs 的长度，即用户个数；输出时表示成功返回的用户信息个数
* \param [out] users     字符串指针数组，用户名称
* \param [out] privs    整型数组，用户权限，枚举 \ref RTDB_PRIV_GROUP 的值之一
* \return rtdb_error
* \remark 用户须保证分配给 users, privs 的空间与 count 相符，如果输入的 count 小于总的用户数，则只返回部分用户信息。且每个指针指向的字符串缓冲区尺寸不小于 \ref RTDB_USER_SIZE。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_users(
  rtdb_int32 handle,
  rtdb_int32 *count,
  RTDB_USER_INFO *infos
  );


/**
* \brief 添加连接黑名单项
* \param handle  连接句柄
* \param [in] addr    阻止连接段地址
* \param [in] mask    阻止连接段子网掩码
* \param [in] desc    阻止连接段的说明，超过 511 字符将被截断
* \return rtdb_error
* \remark addr 和 mask 进行与运算形成一个子网，
* 来自该子网范围内的连接都将被阻止，黑名单的优先级高于信任连接，
* 如果有连接既位于黑名单中，也位于信任连接中，则它将先被阻止。
* 有效的子网掩码的所有 1 位于 0 左侧，例如："255.255.254.0"。
* 当全部为 1 时，表示该子网中只有 addr 一个地址；但不能全部为 0。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_add_blacklist(
  rtdb_int32 handle,
  const char *addr,
  const char *mask,
  const char *desc
  );


/**
* \brief 更新连接连接黑名单项
* \param handle    连接句柄
* \param addr      原阻止连接段地址
* \param mask      原阻止连接段子网掩码
* \param addr_new  新的阻止连接段地址
* \param mask_new  新的阻止连接段子网掩码
* \param desc      新的阻止连接段的说明，超过 511 字符将被截断
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_update_blacklist(
  rtdb_int32 handle,
  const char *addr,
  const char *mask,
  const char *addr_new,
  const char *mask_new,
  const char *desc
  );

/**
* \brief 删除连接黑名单项
* \param handle  连接句柄
* \param addr    阻止连接段地址
* \param mask    阻止连接段子网掩码
* \remark 只有 addr 与 mask 完全相同才视为同一个阻止连接段
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_remove_blacklist(
  rtdb_int32 handle,
  const char *addr,
  const char *mask
  );

/**
* \brief 获得连接黑名单
* \param handle          连接句柄
* \param addrs           字符串指针数组，输出，阻止连接段地址列表
* \param masks           字符串指针数组，输出，阻止连接段子网掩码列表
* \param descs           字符串指针数组，输出，阻止连接段的说明。
* \param [in,out]  count 整型，输入/输出，用户个数
* \remark 用户须保证分配给 addrs, masks, descs 的空间与 count 相符，
*      如果输入的 count 小于输出的 count，则只返回部分阻止连接段，
*      addrs, masks 中每个字符串指针所指缓冲区尺寸不得小于 32 字节，
*      descs 中每个字符串指针所指缓冲区尺寸不得小于 512 字节。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_blacklist(
  rtdb_int32 handle,
  char* const* addrs,
  char* const* masks,
  char* const* descs,
  rtdb_int32 *count
  );


/**
* \brief 添加信任连接段
* \param handle  连接句柄
* \param addr    字符串，输入，信任连接段地址
* \param mask    字符串，输入，信任连接段子网掩码。
* \param priv    整数，输入，信任连接段拥有的用户权限。
* \param desc    字符串，输入，信任连接段的说明，超过 511 字符将被截断。
* \remark addr 和 mask 进行与运算形成一个子网，
*        来自该子网范围内的连接都被视为可信任的，
*        可以不用登录 (rtdb_login)，就直接调用 API，
*        它所拥有的权限在 priv 中指定。
*        有效的子网掩码的所有 1 位于 0 左侧，
*        例如："255.255.254.0"。当全部为 1 时，
*        表示该子网中只有 addr 一个地址；
*        但不能全部为 0。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_add_authorization(
  rtdb_int32 handle,
  const char *addr,
  const char *mask,
  rtdb_int32 priv,
  const char *desc
  );


/**
* \brief 更新信任连接段
* \param handle    连接句柄
* \param addr      字符串，输入，原信任连接段地址。
* \param mask      字符串，输入，原信任连接段子网掩码。
* \param addr_new  字符串，输入，新的信任连接段地址。
* \param mask_new  字符串，输入，新的信任连接段子网掩码。
* \param priv      整数，输入，新的信任连接段拥有的用户权限。
* \param desc      字符串，输入，新的信任连接段的说明，超过 511 字符将被截断。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_update_authorization(
  rtdb_int32 handle,
  const char *addr,
  const char *mask,
  const char *addr_new,
  const char *mask_new,
  rtdb_int32 priv,
  const char *desc
  );

/**
* \brief 删除信任连接段
* \param handle  连接句柄
* \param addr    字符串，输入，信任连接段地址。
* \param mask    字符串，输入，信任连接段子网掩码。
* \remark 只有 addr 与 mask 完全相同才视为同一个信任连接段
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_remove_authorization(
  rtdb_int32 handle,
  const char *addr,
  const char *mask
  );

/**
* \brief 获得所有信任连接段
* \param handle          连接句柄
* \param addrs           字符串指针数组，输出，信任连接段地址列表
* \param masks           字符串指针数组，输出，信任连接段子网掩码列表
* \param [in,out]  privs 整型数组，输出，信任连接段权限列表
* \param descs           字符串指针数组，输出，信任连接段的说明。
* \param [in,out]  count 整型，输入/输出，用户个数
* \remark 用户须保证分配给 addrs, masks, privs, descs 的空间与 count 相符，
*        如果输入的 count 小于输出的 count，则只返回部分信任连接段，
*        addrs, masks 中每个字符串指针所指缓冲区尺寸不得小于 32 字节，
*        descs 中每个字符串指针所指缓冲区尺寸不得小于 512 字节。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_authorizations(
  rtdb_int32 handle,
  char* const* addrs,
  char* const* masks,
  rtdb_int32 *privs,
  char* const* descs,
  rtdb_int32 *count
  );

/**
* \brief 获取 RTDB 服务器当前UTC时间
*
* \param handle       连接句柄
* \param hosttime     整型，输出，Rtdb服务器的当前UTC时间，
*                     表示距离1970年1月1日08:00:00的秒数。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_host_time(
  rtdb_int32 handle,
  rtdb_int32 *hosttime
  );

/**
* \brief 获取 RTDB 服务器当前UTC时间
*
* \param handle       连接句柄
* \param hosttime     整型，输出，Rtdb服务器的当前UTC时间，
*                     表示距离1970年1月1日08:00:00的秒数。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_host_time64(
    rtdb_int32 handle,
    rtdb_timestamp_type* hosttime
);

/**
* \brief 根据时间跨度值生成时间格式字符串
*
* \param str          字符串，输出，时间格式字符串，形如:
*                     "1d" 表示时间跨度为24小时。
*                     具体含义参见 rtdb_parse_timespan 注释。
* \param timespan     整型，输入，要处理的时间跨度秒数。
* \remark 字符串缓冲区大小不应小于 32 字节。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_format_timespan(
  char *str,
  rtdb_int32 timespan
  );

/**
* \brief 根据时间格式字符串解析时间跨度值
*
* \param str          字符串，输入，时间格式字符串，规则如下：
*                     [time_span]
*
*                     时间跨度部分可以出现多次，
*                     可用的时间跨度代码及含义如下：
*                     ?y            ?年, 1年 = 365日
*                     ?m            ?月, 1月 = 30 日
*                     ?d            ?日
*                     ?h            ?小时
*                     ?n            ?分钟
*                     ?s            ?秒
*                     例如："1d" 表示时间跨度为24小时。
*
* \param timespan     整型，输出，返回解析得到的时间跨度秒数。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_parse_timespan(
  const char *str,
  rtdb_int32 *timespan
  );

/**
* \brief 根据时间格式字符串解析时间值
*
* \param str          字符串，输入，时间格式字符串，规则如下：
*                     base_time [+|- offset_time]
*
*                     其中 base_time 表示基本时间，有三种形式：
*                     1. 时间字符串，如 "2010-1-1" 及 "2010-1-1 8:00:00"；
*                     2. 时间代码，表示客户端相对时间；
*                     可用的时间代码及含义如下：
*                     td             当天零点
*                     yd             昨天零点
*                     tm             明天零点
*                     mon            本周一零点
*                     tue            本周二零点
*                     wed            本周三零点
*                     thu            本周四零点
*                     fri            本周五零点
*                     sat            本周六零点
*                     sun            本周日零点
*                     3. 星号('*')，表示客户端当前时间。
*
*                     offset_time 是可选的，可以出现多次，
*                     可用的时间偏移代码及含义如下：
*                     [+|-] ?y            偏移?年, 1年 = 365日
*                     [+|-] ?m            偏移?月, 1月 = 30 日
*                     [+|-] ?d            偏移?日
*                     [+|-] ?h            偏移?小时
*                     [+|-] ?n            偏移?分钟
*                     [+|-] ?s            偏移?秒
*                     [+|-] ?ms           偏移?毫秒
*                     例如："*-1d" 表示当前时刻减去24小时。
*
* \param datetime     整型，输出，返回解析得到的时间值。
* \param ms           短整型，输出，返回解析得到的时间毫秒值。
*  备注：ms 可以为空指针，相应的毫秒信息将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_parse_time(
  const char *str,
  rtdb_int64 *datetime,
  rtdb_int16 *ms
  );

/**
* \brief 获取 Rtdb API 调用返回值的简短描述
*
* \param ecode        无符号整型，输入，Rtdb API调用后的返回值，详见rtdb_error.h头文件
* \param message      字符串，输出，返回错误码简短描述
* \param name         字符串，输出，返回错误码宏名称
* \param size         整型，输入，message 参数的字节长度
* \remark 用户须保证分配给 message， name 的空间与 size 相符,
*      name 或 message 可以为空指针，对应的信息将不再返回。
*/
RTDBAPI
void
RTDBAPI_CALLRULE
rtdb_format_message(
  rtdb_error ecode,
  char *message,
  char *name,
  rtdb_int32 size
  );

/**
* \brief 获取任务的简短描述
*
* \param job_id       整型，输入，RTDB_HOST_CONNECT_INFO::job 字段所表示的最近任务的描述
* \param desc         字符串，输出，返回任务描述
* \param name         字符串，输出，返回任务名称
* \param size         整型，输入，desc、name 参数的字节长度
* \remark 用户须保证分配给 desc、name 的空间与 size 相符，
*      name 或 message 可以为空指针，对应的信息将不再返回。
*/
RTDBAPI
void
RTDBAPI_CALLRULE
rtdb_job_message(
  rtdb_int32 job_id,
  char *desc,
  char *name,
  rtdb_int32 size
  );

/**
* \brief 设置连接超时时间
*
* \param handle   连接句柄
* \param socket   整型，输入，要设置超时时间的连接
* \param timeout  整型，输入，超时时间，单位为秒，0 表示始终保持
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_set_timeout(
  rtdb_int32 handle,
  rtdb_int32 socket,
  rtdb_int32 timeout
  );

/**
* \brief 获得连接超时时间
*
* \param handle   连接句柄
* \param socket   整型，输入，要获取超时时间的连接
* \param timeout  整型，输出，超时时间，单位为秒，0 表示始终保持
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_timeout(
  rtdb_int32 handle,
  rtdb_int32 socket,
  rtdb_int32 *timeout
  );

/**
* \brief 断开已知连接
*
* \param handle    连接句柄
* \param socket    整型，输入，要断开的连接
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_kill_connection(
  rtdb_int32 handle,
  rtdb_int32 socket
  );

/**
* \brief 获得字符串型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
* \param str       字符串型，输出，存放取得的字符串参数值。
* \param size      整型，输入，字符串缓冲区尺寸。
* \remark 本接口只接受 [RTDB_PARAM_STR_FIRST, RTDB_PARAM_STR_LAST) 范围之内参数索引。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_db_info1(
  rtdb_int32 handle,
  rtdb_int32 index,
  char *str,
  rtdb_int32 size
  );

/**
* \brief 获得整型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
* \param value     无符号整型，输出，存放取得的整型参数值。
* \remark 本接口只接受 [RTDB_PARAM_INT_FIRST, RTDB_PARAM_INT_LAST) 范围之内参数索引。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_db_info2(
  rtdb_int32 handle,
  rtdb_int32 index,
  rtdb_uint32 *value
  );

/**
* \brief 设置字符串型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要设置的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
*                  其中，仅以下列出的枚举值可用：
*                  RTDB_PARAM_AUTO_BACKUP_PATH,
*                  RTDB_PARAM_SERVER_SENDER_IP,
* \param str       字符串型，输入，新的参数值。
* \remark 如果修改了启动参数，将返回 RtE_DATABASE_NEED_RESTART 提示码。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_set_db_info1(
  rtdb_int32 handle,
  rtdb_int32 index,
  const char *str
  );

/**
* \brief 设置整型数据库系统参数
*
* \param handle    连接句柄
* \param index     整型，输入，要取得的参数索引，参见枚举 RTDB_DB_PARAM_INDEX。
*                  其中，仅以下列出的枚举值可用：
*                  RTDB_PARAM_SERVER_IPC_SIZE,
*                  RTDB_PARAM_EQUATION_IPC_SIZE,
*                  RTDB_PARAM_HASH_TABLE_SIZE,
*                  RTDB_PARAM_TAG_DELETE_TIMES,
*                  RTDB_PARAM_SERVER_PORT,
*                  RTDB_PARAM_SERVER_SENDER_PORT,
*                  RTDB_PARAM_SERVER_RECEIVER_PORT,
*                  RTDB_PARAM_SERVER_MODE,
*                  RTDB_PARAM_ARV_PAGES_NUMBER,
*                  RTDB_PARAM_ARVEX_PAGES_NUMBER,
*                  RTDB_PARAM_EXCEPTION_AT_SERVER,
*                  RTDB_PARAM_EX_ARCHIVE_SIZE,
*                  RTDB_PARAM_ARCHIVE_BATCH_SIZE,
*                  RTDB_PARAM_ARV_ASYNC_QUEUE_SLOWER_DOOR,
*                  RTDB_PARAM_ARV_ASYNC_QUEUE_NORMAL_DOOR,
*                  RTDB_PARAM_INDEX_ALWAYS_IN_MEMORY,
*                  RTDB_PARAM_DISK_MIN_REST_SIZE,
*                  RTDB_PARAM_DELAY_OF_AUTO_MERGE_OR_ARRANGE,
*                  RTDB_PARAM_START_OF_AUTO_MERGE_OR_ARRANGE,
*                  RTDB_PARAM_STOP_OF_AUTO_MERGE_OR_ARRANGE,
*                  RTDB_PARAM_START_OF_AUTO_BACKUP,
*                  RTDB_PARAM_STOP_OF_AUTO_BACKUP,
*                  RTDB_PARAM_MAX_LATENCY_OF_SNAPSHOT,
*                  RTDB_PARAM_PAGE_ALLOCATOR_RESERVE_SIZE,
* \param value     无符号整型，输入，新的参数值。
* \remark 如果修改了启动参数，将返回 RtE_DATABASE_NEED_RESTART 提示码。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_set_db_info2(
  rtdb_int32 handle,
  rtdb_int32 index,
  rtdb_uint32 value
  );

/**
* \brief 获得逻辑盘符
*
* \param handle     连接句柄
* \param drivers    字符数组，输出，
*                   返回逻辑盘符组成的字符串，每个盘符占一个字符。
* \remark drivers 的内存空间由用户负责维护，长度应不小于 32。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_logical_drivers(
  rtdb_int32 handle,
  char *drivers
  );

/**
* \brief 打开目录以便遍历其中的文件和子目录。
*
* \param handle       连接句柄
* \param dir          字符串，输入，要打开的目录
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_open_path(
  rtdb_int32 handle,
  const char *dir
  );

/**
* \brief 读取目录中的文件或子目录
*
* \param handle      连接句柄
* \param path        字符数组，输出，返回的文件、子目录全路径
* \param is_dir      短整数，输出，返回 1 为目录，0 为文件
* \param atime       整数，输出，为文件时，返回访问时间
* \param ctime       整数，输出，为文件时，返回建立时间
* \param mtime       整数，输出，为文件时，返回修改时间
* \param size        64 位整数，输出，为文件时，返回文件大小
* \remark path 的内存空间由用户负责维护，尺寸应不小于 RTDB_MAX_PATH。
*      当返回值为 RtE_BATCH_END 时表示目录下所有子目录和文件已经遍历完毕。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_path(
  rtdb_int32 handle,
  char *path,
  rtdb_int16 *is_dir,
  rtdb_int32 *atime,
  rtdb_int32 *ctime,
  rtdb_int32 *mtime,
  rtdb_int64 *size
  );

/**
* \brief 读取目录中的文件或子目录
*
* \param handle      连接句柄
* \param path        字符数组，输出，返回的文件、子目录全路径
* \param is_dir      短整数，输出，返回 1 为目录，0 为文件
* \param atime       整数，输出，为文件时，返回访问时间
* \param ctime       整数，输出，为文件时，返回建立时间
* \param mtime       整数，输出，为文件时，返回修改时间
* \param size        64 位整数，输出，为文件时，返回文件大小
* \remark path 的内存空间由用户负责维护，尺寸应不小于 RTDB_MAX_PATH。
*      当返回值为 RtE_BATCH_END 时表示目录下所有子目录和文件已经遍历完毕。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_path64(
    rtdb_int32 handle,
    char* path,
    rtdb_int16* is_dir,
    rtdb_timestamp_type* atime,
    rtdb_timestamp_type* ctime,
    rtdb_timestamp_type* mtime,
    rtdb_int64* size
);

/**
*
* \brief 关闭当前遍历的目录
*
* \param handle      连接句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_close_path(
  rtdb_int32 handle
  );

/**
*
* \brief 建立目录
*
* \param handle       连接句柄
* \param dir          字符串，输入，新建目录的全路径
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_mkdir(
  rtdb_int32 handle,
  const char *dir
  );

/**
*
* \brief 获得指定服务器端文件的大小
*
* \param handle     连接句柄
* \param file       字符串，输入，文件名
* \param size       64 位整数，输出，文件大小
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_file_size(
  rtdb_int32 handle,
  const char *file,
  rtdb_int64 *size
  );

/**
*
* \brief 读取服务器端指定文件的内容
*
* \param handle       连接句柄
* \param file         字符串，输入，要读取内容的文件名
* \param content      字符数组，输出，文件内容
* \param pos          64 位整型，输入，读取文件的起始位置
* \param size         整型，输入/输出，
*                     输入时表示要读取文件内容的字节大小；
*                     输出时表示实际读取的字节数
* \remark 用户须保证分配给 content 的空间与 size 相符。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_file(
  rtdb_int32 handle,
  const char *file,
  char *content,
  rtdb_int64 pos,
  rtdb_int64 *size
  );


/**
*
* \brief 取得数据库允许的blob与str类型测点的最大长度
*
* \param handle       连接句柄
* \param len          整形，输出参数，代表数据库允许的blob、str类型测点的最大长度
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_get_max_blob_len(
  rtdb_int32 handle,
  rtdb_int32 *len
  );

/**
*
* \brief 取得质量码对应的定义
*
* \param handle       连接句柄
* \param count        质量码个数，输入参数，
* \param qualities    质量码，输入参数
* \param definitions  输出参数，0~255为RTDB质量码（参加rtdb.h文件），256~511为OPC质量码，大于511为用户自定义质量码
* \param lens         输出参数，每个定义对应的长度
* \remark OPC质量码把8位分3部分定义：XX XXXX XX，对应着：品质位域 分状态位域 限定位域
* 品质位域：00（Bad）01（Uncertain）10（N/A）11（Good）
* 分状态位域：不同品质位域对应各自的分状态位域
* 限定位域：00（Not limited）01（Low limited）10（high limited）11（Constant）
* 三个域之间用逗号隔开，输出到definitions参数中，前面有有RTDB，OPC或者USER标识，说明标签点类别
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_format_quality(
  rtdb_int32 handle,
  rtdb_int32 *count,
  rtdb_int16 *qualities,
  rtdb_byte **definitions,
  rtdb_int32 *lens
  );

/*
* 命名：rtdb_write_named_type_field_by_name32
* 功能：按名称填充自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_name]   自定义类型中需要填充的字段的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数
*      [field_type]   field_name字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入/输出参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      需要填充的字段数值的缓冲区,输入参数
*     [field_len]    自定义类型中字段数值的缓冲区中数据的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_write_named_type_field_by_name32(
  rtdb_int32 handle,
  const char* type_name,
  const char* field_name,
  rtdb_int32 field_type,
  void* object,
  rtdb_length_type object_len,
  const void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_write_named_type_field_by_pos32
* 功能：按位置填充自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_pos]    自定义类型中需要填充的字段的位置，指字段在所有字段中的位置，从0开始，输入参数
*      [field_type]   field_pos位置所在字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入/输出参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      需要填充的字段数值的缓冲区,输入参数
*     [field_len]    自定义类型中字段数值的缓冲区中数据的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_write_named_type_field_by_pos32(
  rtdb_int32 handle,
  const char* type_name,
  rtdb_int32 field_pos,
  rtdb_int32 field_type,
  void* object,
  rtdb_length_type object_len,
  const void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_read_named_type_field_by_name32
* 功能：按名称提取自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_name]   自定义类型中需要提取的字段的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数
*      [field_type]   field_name字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      被读取的字段的数值的缓冲区,输入/输出参数
*     [field_len]    field字段数值缓冲区的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_named_type_field_by_name32(
  rtdb_int32 handle,
  const char* type_name,
  const char* field_name,
  rtdb_int32 field_type,
  const void* object,
  rtdb_length_type object_len,
  void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_read_named_type_field_by_pos32
* 功能：按位置提取自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_pos]    自定义类型中需要提取的字段的位置，指字段在所有字段中的位置，从0开始，输入参数
*      [field_type]   field_pos位置所在字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      被读取的字段的数值的缓冲区,输入/输出参数
*     [field_len]    field字段数值缓冲区的长度,输入参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_named_type_field_by_pos32(
  rtdb_int32 handle,
  const char* type_name,
  rtdb_int32 field_pos,
  rtdb_int32 field_type,
  const void* object,
  rtdb_length_type object_len,
  void* field,
  rtdb_length_type field_len
  );

/*
* 命名：rtdb_named_type_name_field_check
* 功能：检查自定义类型名称及字段命名是否符合规则；
* 规则：1. 只允许使用26个英文字母,数字0-9，下划线；
*       2. 必须以字母作为首字母；
*       3. 大小写不敏感。
* 参数：
*      [check_name]   需要检查的名称
*      [flag]         标志0--类型名称，其他 -- 字段名称，暂不启用
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_named_type_name_field_check(
  const char* check_name,
  rtdb_byte flag GAPI_DEFAULT_VALUE(0)
  );

/**
*
* \brief 判断连接是否可用
* \param handle   连接句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdb_judge_connect_status(rtdb_int32 handle,
    rtdb_int8* change_connection GAPI_DEFAULT_VALUE(0),
    char* current_ip_addr GAPI_DEFAULT_VALUE(0),
    rtdb_int32 size GAPI_DEFAULT_VALUE(0));

/**
* 命名：rtdb_format_ipaddr
* 功能：将整形IP转换为字符串形式的IP
* 参数：
*      [ip]        无符号整型，输入，整形的IP地址
*      [ip_addr]      字符串，输出，字符串IP地址缓冲区
*      [size]         整型，输入，ip_addr 参数的字节长度
* 备注：用户须保证分配给 ip_addr 的空间与 size 相符
*/
RTDBAPI
void
RTDBAPI_CALLRULE
rtdb_format_ipaddr(rtdb_uint32 ip, char* ip_addr, rtdb_int32 size);

/**@}*/

/**
* \defgroup base 基本标签点信息接口
* @{
*/

/**
* 命名：rtdbb_get_equation_by_file_name
* 功能：根据文件名获取方程式
* 参数：
*      [handle]   连接句柄
*      [file_name] 输入，字符串，方程式路径
*      [equation]  输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
*
*备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_equation_by_file_name(
  rtdb_int32 handle,
  const char* file_name,
  char equation[RTDB_MAX_EQUATION_SIZE]
  );

/**
* 命名：rtdbb_get_equation_by_id
* 功能：根ID径获取方程式
* 参数：
*      [handle]   连接句柄
*      [id]				输入，整型，方程式ID
*      [equation] 输出，返回的方程式长度最长为RTDB_MAX_EQUATION_SIZE-1
*
*备注：用户调用时为equation分配的空间不得小于RTDB_MAX_EQUATION_SIZE
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_equation_by_id(
  rtdb_int32 handle,
  rtdb_int32 id,
  char equation[RTDB_MAX_EQUATION_SIZE]
  );


/**
*
* \brief 添加新表
*
* \param handle   连接句柄
* \param field    RTDB_TABLE 结构，输入/输出，表信息。
*                 在输入时，type、name、desc 字段有效；
*                 输出时，id 字段由系统自动分配并返回给用户。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_append_table(
  rtdb_int32 handle,
  RTDB_TABLE *field
  );

/**
*
* \brief 取得标签点表总数
*
* \param handle   连接句柄
* \param count    整型，输出，标签点表总数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_tables_count(
  rtdb_int32 handle,
  rtdb_int32 *count
  );

/**
*
* \brief 取得所有标签点表的ID
*
* \param handle   连接句柄
* \param ids      整型数组，输出，标签点表的id
* \param count    整型，输入/输出，
*                 输入表示 ids 的长度，输出表示标签点表个数
* \remark 用户须保证分配给 ids 的空间与 count 相符
*      如果输入的 count 小于输出的 count，则只返回部分表id
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_tables(
  rtdb_int32 handle,
  rtdb_int32 *ids,
  rtdb_int32 *count
  );

/**
*
* \brief 根据表 id 获取表中包含的标签点数量
*
* \param handle   连接句柄
* \param id       整型，输入，表ID
* \param size     整型，输出，表中标签点数量
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_table_size_by_id(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *size
  );

/**
*
* \brief 根据表名称获取表中包含的标签点数量
*
* \param handle   连接句柄
* \param name     字符串，输入，表名称
* \param size     整型，输出，表中标签点数量
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_table_size_by_name(
  rtdb_int32 handle,
  const char *name,
  rtdb_int32 *size
  );

/**
*
* \brief 根据表 id 获取表中实际包含的标签点数量
*
* \param handle   连接句柄
* \param id       整型，输入，表ID
* \param size     整型，输出，表中标签点数量
* 注意：通过此API获取标签点数量，然后搜索此表中的标签点得到的数量可能会不一致，这是由于服务内部批量建点采取了异步的方式。
*       一般情况下请使用rtdbb_get_table_size_by_id来获取表中的标签点数量。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_table_real_size_by_id(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *size
  );

/**
*
* \brief 根据标签点表 id 获取表属性
*
* \param handle 连接句柄
* \param field  RTDB_TABLE 结构，输入/输出，标签点表属性，
*               输入时指定 id 字段，输出时返回 type、name、desc 字段。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_table_property_by_id(
  rtdb_int32 handle,
  RTDB_TABLE *field
  );

/**
*
* \brief 根据表名获取标签点表属性
*
* \param handle 连接句柄
* \param field  RTDB_TABLE 结构，输入/输出，标签点表属性
*               输入时指定 name 字段，输出时返回 id、type、desc 字段。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_table_property_by_name(
  rtdb_int32 handle,
  RTDB_TABLE *field
  );

/**
*
* \brief 使用完整的属性集来创建单个标签点
*
* \param handle 连接句柄
* \param base RTDB_POINT 结构，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
* \param scan RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
* \param calc RTDB_CALC_POINT 结构，输入，计算标签点扩展属性集。
* \remark 如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_insert_point(
  rtdb_int32 handle,
  RTDB_POINT *base,
  RTDB_SCAN_POINT *scan,
  RTDB_CALC_POINT *calc
  );


/**
* 命名：rtdbb_insert_max_point
* 功能：使用最大长度的完整属性集来创建单个标签点
* 参数：
*      [handle] 连接句柄
*      [base] RTDB_POINT 结构，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
*      [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
*      [calc] RTDB_MAX_CALC_POINT 结构，输入，计算标签点扩展属性集。
* 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_insert_max_point(
  rtdb_int32 handle,
  RTDB_POINT *base,
  RTDB_SCAN_POINT *scan,
  RTDB_MAX_CALC_POINT *calc
  );

/**
* 命名：rtdbb_insert_max_points
* 功能：使用最大长度的完整属性集来批量创建标签点
* 参数：
*      [handle] 连接句柄
*	   [count] count, 输入，base/scan/calc数组的长度；输出，成功的个数
*      [bases] RTDB_POINT 数组，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
*      [scans] RTDB_SCAN_POINT 数组，输入，采集标签点扩展属性集。
*      [calcs] RTDB_MAX_CALC_POINT 数组，输入，计算标签点扩展属性集。
*	   [errors] rtdb_error数组，输出，对应每个标签点的结果
* 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_insert_max_points(
  rtdb_int32 handle,
  rtdb_int32* count,
  RTDB_POINT* bases,
  RTDB_SCAN_POINT* scans,
  RTDB_MAX_CALC_POINT* calcs,
  rtdb_error* errors
);

/**
*
* 功能  使用最小的属性集来创建单个标签点
*
* \param handle     连接句柄
* \param tag        字符串，输入，标签点名称
* \param type       整型，输入，标签点数据类型，取值 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、
*                   RTDB_CHAR、RTDB_UINT16、RTDB_UINT32、RTDB_INT32、RTDB_INT64、
*                   RTDB_REAL16、RTDB_REAL32、RTDB_REAL64、RTDB_COOR、RTDB_STRING、RTDB_BLOB 之一。
* \param table_id   整型，输入，标签点所属表 id
* \param use_ms     短整型，输入，标签点时间戳精度，0 为秒；1 为纳秒。
* \param point_id   整型，输出，标签点 id
* \remark 标签点的其余属性将取默认值。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_insert_base_point(
  rtdb_int32 handle,
  const char *tag,
  rtdb_int32 type,
  rtdb_int32 table_id,
  rtdb_int16 use_ms,
  rtdb_int32 *point_id
  );

/**
* 命名：rtdbb_insert_named_type_point
* 功能：使用完整的属性集来创建单个自定义数据类型标签点
* 参数：
*      [handle] 连接句柄
*      [base] RTDB_POINT 结构，输入/输出，
*      输入除 id, createdate, creator, changedate, changer 字段外的其它字段，输出 id 字段。
*      [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
*      [name] 字符串，输入，自定义数据类型的名字。
* 备注：如果新建的标签点没有对应的扩展属性集，可置为空指针。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_insert_named_type_point(
  rtdb_int32 handle,
  RTDB_POINT *base,
  RTDB_SCAN_POINT *scan,
  const char* name
  );

/**
*
* \brief 根据 id 删除单个标签点
*
* \param handle 连接句柄
* \param id     整型，输入，标签点标识
* \remark 通过本接口删除的标签点为可回收标签点，
*        可以通过 rtdbb_recover_point 接口恢复。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_remove_point_by_id(
  rtdb_int32 handle,
  rtdb_int32 id
  );

/**
*
* \brief 根据标签点全名删除单个标签点
* \param handle        连接句柄
* \param table_dot_tag  字符串，输入，标签点全名称："表名.标签点名"
* \remark 通过本接口删除的标签点为可回收标签点，
*        可以通过 rtdbb_recover_point 接口恢复。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_remove_point_by_name(
  rtdb_int32 handle,
  const char *table_dot_tag
  );

/**
* 命名：rtdbb_move_point_by_id
* 功能：根据 id 移动单个标签点到其他表
* 参数：
*      [handle] 连接句柄
*      [id]     整型，输入，标签点标识
*	   [dest_table_name] 字符串，输入，移动的目标表名称
* 备注：通过本接口移动标签点后不改变标签点的id，且快照
*       和历史数据都不受影响
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_move_point_by_id(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* dest_table_name
);

/**
*
* \brief 批量获取标签点属性
*
* \param handle 连接句柄
* \param count  整数，输入，表示标签点个数。
* \param base   RTDB_POINT 结构数组，输入/输出，标签点基本属性列表，
*                 输入时，id 字段指定需要得到属性的标签点，输出时，其它字段返回标签点属性值。
* \param scan   RTDB_SCAN_POINT 结构数组，输出，采集标签点扩展属性列表
* \param calc   RTDB_CALC_POINT 结构数组，输出，计算标签点扩展属性列表
* \param errors 无符号整型数组，输出，获取标签属性的返回值列表，参考rtdb_error.h
* \remark 用户须保证分配给 base、scan、calc、errors 的空间与 count 相符，
*        扩展属性集 scan、calc 可为空指针，此时将不返回对应的扩展属性集。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_points_property(
  rtdb_int32 handle,
  rtdb_int32 count,
  RTDB_POINT *base,
  RTDB_SCAN_POINT *scan,
  RTDB_CALC_POINT *calc,
  rtdb_error *errors
  );


/**
* 命名：rtdbb_get_max_points_property
* 功能：按最大长度批量获取标签点属性
* 参数：
*        [handle] 连接句柄
*        [count]  整数，输入，表示标签点个数。
*        [base]   RTDB_POINT 结构数组，输入/输出，标签点基本属性列表，
*                 输入时，id 字段指定需要得到属性的标签点，输出时，其它字段返回标签点属性值。
*        [scan]   RTDB_SCAN_POINT 结构数组，输出，采集标签点扩展属性列表
*        [calc]   RTDB_MAX_CALC_POINT 结构数组，输出，计算标签点扩展属性列表
*        [errors] 无符号整型数组，输出，获取标签属性的返回值列表，参考rtdb_error.h
* 备注：用户须保证分配给 base、scan、calc、errors 的空间与 count 相符，
*        扩展属性集 scan、calc 可为空指针，此时将不返回对应的扩展属性集。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_max_points_property(
  rtdb_int32 handle,
  rtdb_int32 count,
  RTDB_POINT *base,
  RTDB_SCAN_POINT *scan,
  RTDB_MAX_CALC_POINT *calc,
  rtdb_error *errors
  );



/**
*
* \brief 搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_search(
  rtdb_int32 handle,
  const char *tagmask,
  const char *tablemask,
  const char *source,
  const char *unit,
  const char *desc,
  const char *instrument,
  rtdb_int32 mode,
  rtdb_int32 *ids,
  rtdb_int32 *count
  );

/**
*
* \brief 分批继续搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param start         整型，输入，搜索起始位置。
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*"。
*        当搜索到的标签点数比提供的要小时，表示这是最后一批符合条件的标签点 (即全部搜索完毕),
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_search_in_batches(
  rtdb_int32 handle,
  rtdb_int32 start,
  const char *tagmask,
  const char *tablemask,
  const char *source,
  const char *unit,
  const char *desc,
  const char *instrument,
  rtdb_int32 mode,
  rtdb_int32 *ids,
  rtdb_int32 *count
  );


/**
*
* \brief 搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param typemask      字符串，输入参数，标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,
*                        内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
* \param classofmask   整型，输入参数，标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，
*                        当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
* \param timeunitmask  整型，输入参数，标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，
*                        当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
* \param othertypemask 整型，输入参数，使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，
*                        当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，
*                        此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
* \param othertypemaskvalue
*                        字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，
*                        如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略,
*                        当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，
*                        当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_search_ex(
  rtdb_int32 handle,
  const char *tagmask,
  const char *tablemask,
  const char *source,
  const char *unit,
  const char *desc,
  const char *instrument,
  const char *typemask,
  rtdb_int32 classofmask,
  rtdb_int32 timeunitmask,
  rtdb_int32 othertypemask,
  const char *othertypemaskvalue,
  rtdb_int32 mode,
  rtdb_int32 *ids,
  rtdb_int32 *count
  );


/**
*
* \brief 搜索符合条件的标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE，支持多个搜索条件，以空格分隔。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param typemask      字符串，输入参数，标签点类型名称。缺省设置为空，长度不得超过 RTDB_TYPE_NAME_SIZE,
*                        内置的普通数据类型可以使用 bool、uint8、datetime等字符串表示，不区分大小写，支持模糊搜索。
* \param classofmask   整型，输入参数，标签点的类别，缺省设置为-1，表示可以是任意类型的标签点，
*                        当使用标签点类型作为搜索条件时，必须是RTDB_CLASS枚举中的一项或者多项的组合。
* \param timeunitmask  整型，输入参数，标签点的时间戳精度，缺省设置为-1，表示可以是任意时间戳精度，
*                        当使用此时间戳精度作为搜索条件时，timeunitmask的值可以为0或1，0表示时间戳精度为秒，1表示纳秒
* \param othertypemask 整型，输入参数，使用其他标签点属性作为搜索条件，缺省设置为0，表示不作为搜索条件，
*                        当使用此参数作为搜索条件时，othertypemaskvalue作为对应的搜索值，
*                        此参数的取值可以参考rtdb.h文件中的RTDB_SEARCH_MASK。
* \param othertypemaskvalue
*                        字符串，输入参数，当使用其他标签点属性作为搜索条件时，此参数作为对应的搜索值，缺省设置为0，表示不作为搜索条件，
*                        如果othertypemask的值为0，或者RTDB_SEARCH_NULL，则此参数被忽略,
*                        当othertypemask对应的标签点属性为数值类型时，此搜索值只支持相等判断，
*                        当othertypemask对应的标签点属性为字符串类型时，此搜索值支持模糊搜索。
* \param count         整型，输出，表示搜索到的标签点个数
* \remark  各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、tablemask 为空指针，则表示使用缺省设置"*",
*        多个搜索条件可以通过空格分隔，比如"demo_*1 demo_*2"，会将满足demo_*1或者demo_*2条件的标签点搜索出来。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_search_points_count(
  rtdb_int32 handle,
  const char *tagmask,
  const char *tablemask,
  const char *source,
  const char *unit,
  const char *desc,
  const char *instrument,
  const char *typemask,
  rtdb_int32 classofmask,
  rtdb_int32 timeunitmask,
  rtdb_int32 othertypemask,
  const char *othertypemaskvalue,
  rtdb_int32 *count
  );

/**
* 命名：rtdbb_remove_table_by_id
* \brief 根据表 id 删除表及表中标签点
*
* \param handle        连接句柄
* \param id            整型，输入，表 id
* \remark 删除的表不可恢复，删除的标签点可以通过 rtdbb_recover_point 接口恢复。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_remove_table_by_id(
  rtdb_int32 handle,
  rtdb_int32 id
  );

/**
*
* \brief 根据表名删除表及表中标签点
*
* \param handle        连接句柄
* \param name          字符串，输入，表名称
* \remark 删除的表不可恢复，删除的标签点可以通过 rtdbb_recover_point 接口恢复。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_remove_table_by_name(
  rtdb_int32 handle,
  const char *name
  );

/**
*
* \brief 更新单个标签点属性
*
* \param handle        连接句柄
* \param base RTDB_POINT 结构，输入，基本标签点属性集。
* \param scan RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
* \param calc RTDB_CALC_POINT 结构，输入，计算标签点扩展属性集。
* \remark 标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，
*      changedate、changer、createdate、creator 字段由系统维护，其余字段均可修改，
*      包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_update_point_property(
  rtdb_int32 handle,
  const RTDB_POINT *base,
  const RTDB_SCAN_POINT *scan,
  const RTDB_CALC_POINT *calc
  );

/**
* 命名：rtdbb_update_max_point_property
* 功能：按最大长度更新单个标签点属性
* 参数：
*        [handle]        连接句柄
*        [base] RTDB_POINT 结构，输入，基本标签点属性集。
*        [scan] RTDB_SCAN_POINT 结构，输入，采集标签点扩展属性集。
*        [calc] RTDB_MAX_CALC_POINT 结构，输入，计算标签点扩展属性集。
* 备注：标签点由 base 参数的 id 字段指定，其中 id、table、type、millisecond 字段不能修改，
*      changedate、changer、createdate、creator 字段由系统维护，其余字段均可修改，
*      包括 classof 字段。输入参数中 scan、calc 可为空指针，对应的扩展属性将保持不变。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_update_max_point_property(
  rtdb_int32 handle,
  const RTDB_POINT *base,
  const RTDB_SCAN_POINT *scan,
  const RTDB_MAX_CALC_POINT *calc
  );


/**
*
* \brief 根据 "表名.标签点名" 格式批量获取标签点标识
*
* \param handle           连接句柄
* \param count            整数，输入/输出，输入时表示标签点个数
*                           (即table_dot_tags、ids、types、classof、use_ms 的长度)，
*                           输出时表示找到的标签点个数
* \param table_dot_tags   字符串指针数组，输入，"表名.标签点名" 列表
* \param ids              整型数组，输出，标签点标识列表, 返回 0 表示未找到
* \param types            整型数组，输出，标签点数据类型
* \param classof          整型数组，输出，标签点类别
* \param use_ms           短整型数组，输出，时间戳精度，
*                           返回 1 表示时间戳精度为纳秒， 为 0 表示为秒。
* \remark 用户须保证分配给 table_dot_tags、ids、types、classof、use_ms 的空间与count相符，
*        其中 types、classof、use_ms 可为空指针，对应的字段将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_find_points(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const char* const* table_dot_tags,
  rtdb_int32 *ids,
  rtdb_int32 *types,
  rtdb_int32 *classof,
  rtdb_int16 *use_ms
  );

/**
*
* \brief 根据 "表名.标签点名" 格式批量获取标签点标识
*
* \param handle           连接句柄
* \param count            整数，输入/输出，输入时表示标签点个数
*                           (即table_dot_tags、ids、types、classof、use_ms 的长度)，
*                           输出时表示找到的标签点个数
* \param table_dot_tags   字符串指针数组，输入，"表名.标签点名" 列表
* \param ids              整型数组，输出，标签点标识列表, 返回 0 表示未找到
* \param types            整型数组，输出，标签点数据类型
* \param classof          整型数组，输出，标签点类别
* \param precisions       数组，输出，时间戳精度，
*                           0表示秒，1表示毫秒，2表示微秒，3纳秒。
* \param errors           无符号整型数组，输出，表示每个标签点的查询结果的错误码
* \remark 用户须保证分配给 table_dot_tags、ids、types、classof、precisions、errors 的空间与count相符，
*        其中 types、classof、precisions、errors 可为空指针，对应的字段将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_find_points_ex(
    rtdb_int32 handle,
    rtdb_int32* count,
    const char* const* table_dot_tags,
    rtdb_int32* ids,
    rtdb_int32* types,
    rtdb_int32* classof,
    rtdb_precision_type* precisions,
    rtdb_error* errors
);

/**
*
* \brief 根据标签属性字段对标签点标识进行排序
*
* \param handle           连接句柄
* \param count            整数，输入，表示标签点个数, 即 ids 的长度
* \param ids              整型数组，输入，标签点标识列表
* \param index            整型，输入，属性字段枚举，参见 RTDB_TAG_FIELD_INDEX，
*                           将根据该字段对 ID 进行排序。
* \param flag             整型，输入，标志位组合，参见 RTDB_TAG_SORT_FLAG 枚举，其中
*                           RTDB_SORT_FLAG_DESCEND             表示降序排序，不设置表示升序排列；
*                           RTDB_SORT_FLAG_CASE_SENSITIVE      表示进行字符串类型字段比较时大小写敏感，不设置表示不区分大小写；
*                           RTDB_SORT_FLAG_RECYCLED            表示对可回收标签进行排序，不设置表示对正常标签排序，
*                           不同的标志位可通过"或"运算连接在一起，
*                           当对可回收标签排序时，以下字段索引不可使用：
*                               RTDB_TAG_INDEX_TIMESTAMP
*                               RTDB_TAG_INDEX_VALUE
*                               RTDB_TAG_INDEX_QUALITY
* \remark 用户须保证分配给 ids 的空间与 count 相符, 如果 ID 指定的标签并不存在，
*        或标签不具备要求排序的字段 (如对非计算点进行方程式排序)，它们将被放置在数组的尾部。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_sort_points(
  rtdb_int32 handle,
  rtdb_int32 count,
  rtdb_int32 *ids,
  rtdb_int32 index,
  rtdb_int32 flag
  );

/**
*
* \brief 根据表 ID 更新表名称。
*
* \param handle    连接句柄
* \param tab_id    整型，输入，要修改表的标识
* \param name      字符串，输入，新的标签点表名称。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_update_table_name(
  rtdb_int32 handle,
  rtdb_int32 tab_id,
  const char *name
  );

/**
*
* \brief 根据表 ID 更新表描述。
*
* \param handle    连接句柄
* \param tab_id    整型，输入，要修改表的标识
* \param desc      字符串，输入，新的表描述。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_update_table_desc_by_id(
  rtdb_int32 handle,
  rtdb_int32 tab_id,
  const char *desc
  );

/**
*
* \brief 根据表名称更新表描述。
*
* \param handle    连接句柄
* \param name      字符串，输入，要修改表的名称。
* \param desc      字符串，输入，新的表描述。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_update_table_desc_by_name(
  rtdb_int32 handle,
  const char *name,
  const char *desc
  );

/**
*
* \brief 恢复已删除标签点
*
* \param handle    连接句柄
* \param table_id  整型，输入，要将标签点恢复到的表标识
* \param point_id  整型，输入，待恢复的标签点标识
* 备注: 本接口只对可回收标签点(通过接口rtdbb_remove_point_by_id/rtdbb_remove_point_by_tag)有效，
*        对正常的标签点没有作用。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_recover_point(
  rtdb_int32 handle,
  rtdb_int32 table_id,
  rtdb_int32 point_id
  );

/**
*
* \brief 清除标签点
*
* \param handle    连接句柄
* \param id        整数，输入，要清除的标签点标识
* 备注: 本接口仅对可回收标签点(通过接口rtdbb_remove_point_by_id/rtdbb_remove_point_by_name)有效，
*      对正常的标签点没有作用。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_purge_point(
  rtdb_int32 handle,
  rtdb_int32 id
  );


/**
*
* \brief 获取可回收标签点数量
*
* \param handle    连接句柄
* \param count     整型，输出，可回收标签点的数量
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_recycled_points_count(
  rtdb_int32 handle,
  rtdb_int32 *count
  );

/**
*
* \brief 获取可回收标签点 id 列表
*
* \param handle    连接句柄
* \param ids       整型数组，输出，可回收标签点 id
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids 的长度，
*                    输出时表示成功获取标签点的个数。
* \remark 用户须保证 ids 的长度与 count 一致
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_recycled_points(
  rtdb_int32 handle,
  rtdb_int32 *ids,
  rtdb_int32 *count
  );

/**
* 命名：rtdbb_search_recycled_points
* 功能：搜索符合条件的可回收标签点，使用标签点名时支持通配符
* 参数：
*        [handle]        连接句柄
*        [tagmask]       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
*        [tablemask]     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
*        [source]        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
*        [unit]          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
*        [desc]          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
*        [instrument]    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
*        [mode]          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
*        [ids]           整型数组，输出，返回搜索到的标签点标识列表
*        [count]         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* 备注：用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、fullmask 为空指针，则表示使用缺省设置"*"
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_search_recycled_points(
  rtdb_int32 handle,
  const char *tagmask,
  const char *fullmask,
  const char *source,
  const char *unit,
  const char *desc,
  const char *instrument,
  rtdb_int32 mode,
  rtdb_int32 *ids,
  rtdb_int32 *count
  );

/**
*
* \brief 分批搜索符合条件的可回收标签点，使用标签点名时支持通配符
*
* \param handle        连接句柄
* \param start         整型，输入，搜索的起始位置。
* \param tagmask       字符串，输入，标签点名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
* \param tablemask     字符串，输入，标签点表名称掩码，支持"*"和"?"通配符，缺省设置为"*"，长度不得超过 RTDB_TAG_SIZE。
* \param source        字符串，输入，数据源集合，字符串中的每个字符均表示一个数据源，
*                        空字符串表示不用数据源作搜索条件，缺省设置为空，长度不得超过 RTDB_DESC_SIZE。
* \param unit          字符串，输入，标签点工程单位的子集，工程单位中包含该参数的标签点均满足条件，
*                        空字符串表示不用工程单位作搜索条件，缺省设置为空，长度不得超过 RTDB_UNIT_SIZE。
* \param desc          字符串，输入，标签点描述的子集，描述中包含该参数的标签点均满足条件，
*                        空字符串表示不用描述作搜索条件，缺省设置为空，长度不得超过 RTDB_SOURCE_SIZE。
* \param instrument    字符串，输入参数，标签点设备名称。缺省设置为空，长度不得超过 RTDB_INSTRUMENT_SIZE。
* \param mode          整型，RTDB_SORT_BY_TABLE、RTDB_SORT_BY_TAG、RTDB_SORT_BY_ID 之一，
*                        搜索结果的排序模式，输入，缺省值为RTDB_SORT_BY_TABLE
* \param ids           整型数组，输出，返回搜索到的标签点标识列表
* \param count         整型，输入/输出，输入时表示 ids 的长度，输出时表示搜索到的标签点个数
* \remark 用户须保证分配给 ids 的空间与 count 相符，各参数中包含的搜索条件之间的关系为"与"的关系，
*        用包含通配符的标签点名称作搜索条件时，如果第一个字符不是通配符(如"ai67*")，会得到最快的搜索速度。
*        如果 tagmask、fullmask 为空指针，则表示使用缺省设置"*"
*        当搜索到的标签点数比提供的要小时，表示这是最后一批符合条件的标签点 (即全部搜索完毕)。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_search_recycled_points_in_batches(
  rtdb_int32 handle,
  rtdb_int32 start,
  const char *tagmask,
  const char *fullmask,
  const char *source,
  const char *unit,
  const char *desc,
  const char *instrument,
  rtdb_int32 mode,
  rtdb_int32 *ids,
  rtdb_int32 *count
  );

/**
*
* \brief 获取可回收标签点的属性
*
* \param handle   连接句柄
* \param base     RTDB_POINT 结构，输入/输出，标签点基本属性。
输入时，由 id 字段指定要取得的可回收标签点。
* \param scan     RTDB_SCAN_POINT 结构，输出，标签点采集扩展属性
* \param calc     RTDB_CALC_POINT 结构，输出，标签点计算扩展属性
* \remark scan、calc 可为空指针，对应的扩展信息将不返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_recycled_point_property(
  rtdb_int32 handle,
  RTDB_POINT *base,
  RTDB_SCAN_POINT *scan,
  RTDB_CALC_POINT *calc
  );

/**
* 命名：rtdbb_get_recycled_max_point_property
* 功能：按最大长度获取可回收标签点的属性
* 参数：
*        [handle]   连接句柄
*        [base]     RTDB_POINT 结构，输入/输出，标签点基本属性。
                    输入时，由 id 字段指定要取得的可回收标签点。
*        [scan]     RTDB_SCAN_POINT 结构，输出，标签点采集扩展属性
*        [calc]     RTDB_MAX_CALC_POINT 结构，输出，标签点计算扩展属性
* 备注：scan、calc 可为空指针，对应的扩展信息将不返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_recycled_max_point_property(
  rtdb_int32 handle,
  RTDB_POINT* base,
  RTDB_SCAN_POINT* scan,
  RTDB_MAX_CALC_POINT* calc
  );


/**
*
* \brief 清空标签点回收站
*
* \param handle   连接句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_clear_recycler(
  rtdb_int32 handle
  );

/**
* 命名：rtdbb_subscribe_tags_ex
* 功能：标签点属性更改通知订阅
* 参数：
*        [handle]    连接句柄
*        [options]   整型，输入，订阅选项，参见枚举RTDB_OPTION
*                    RTDB_O_AUTOCONN 订阅客户端与数据库服务器网络中断后自动重连并订阅
*        [param]     输入，用户参数，
*                    作为rtdbb_tags_change_ex的param参数
*        [callback]  rtdbb_tags_change_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                    当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                    作为event_type取值调用回掉函数后退出订阅。
*                    当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                    作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                    网络中断期间回掉函数调用频率为最少3秒
*                    event_type参数值含义如下：
*                      RTDB_E_DATA        标签点属性发生更改
*                      RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                      RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                    handle 产生订阅回掉的连接句柄，调用rtdbb_subscribe_tags_ex时的handle参数
*                    param  用户自定义参数，调用rtdbb_subscribe_tags_ex时的param参数
*                    count  event_type为RTDB_E_DATA时表示ids的数量
*                           event_type为其它值时，count值为0
*                    ids    event_type为RTDB_E_DATA时表示属性更改的标签点ID，数量由count指定
*                           event_type为其它值时，ids值为NULL
*                    what   event_type为RTDB_E_DATA时表示属性变更原因，参考RTDB_TAG_CHANGE_REASON
*                           event_type为其它值时，what时值为0
* 备注：用于订阅测点的连接句柄必需是独立的，不能再用来调用其它 api，
*       否则返回 RtE_OTHER_SDK_DOING 错误。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_subscribe_tags_ex(
  rtdb_int32 handle,
  rtdb_uint32 options,
  void* param,
  rtdbb_tags_change_event_ex callback
  );

/**
*
* \brief 取消标签点属性更改通知订阅
*
* \param handle    连接句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_cancel_subscribe_tags(
  rtdb_int32 handle
  );


/**
* 命名：rtdbb_create_named_type
* 功能：创建自定义类型
* 参数：
*        [handle]      连接句柄，输入参数
*        [name]        自定义类型的名称，类型的唯一标示,不能重复，长度不能超过RTDB_TYPE_NAME_SIZE，输入参数
*        [field_count]    自定义类型中包含的字段的个数,输入参数
*        [fields]      自定义类型中包含的字段的属性，RTDB_DATA_TYPE_FIELD结构的数组，个数与field_count相等，输入参数
*              RTDB_DATA_TYPE_FIELD中的length只对type为str或blob类型的数据有效。其他类型忽略
* 备注：自定义类型的大小必须要小于数据页大小(小于数据页大小的2/3，即需要合理定义字段的个数及每个字段的长度)。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_create_named_type(
  rtdb_int32 handle,
  const char* name,
  rtdb_int32 field_count,
  const RTDB_DATA_TYPE_FIELD* fields,
  char desc[RTDB_DESC_SIZE]
  );

/**
* 命名：rtdbb_get_named_types_count
* 功能：获取所有的自定义类型的总数
* 参数：
*        [handle]      连接句柄，输入参数
*        [count]      返回所有的自定义类型的总数，输入/输出参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_named_types_count(
  rtdb_int32 handle,
  rtdb_int32* count
  );

/**
* 命名：rtdbb_get_all_named_types
* 功能：获取所有的自定义类型
* 参数：
*        [handle]      连接句柄，输入参数
*        [count]      返回所有的自定义类型的总数，输入/输出参数，输入:为name,field_counts数组的长度，输出:获取的实际自定义类型的个数
*        [name]        返回所有的自定义类型的名称的数组，每个自定义类型的名称的长度不超过RTDB_TYPE_NAME_SIZE，输入/输出参数
*              输入：name数组长度要等于count.输出：实际获取的自定义类型名称的数组
*        [field_counts]    返回所有的自定义类型所包含字段个数的数组，输入/输出参数
*              输入：field_counts数组长度要等于count。输出:实际每个自定义类型所包含的字段的个数的数组
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_all_named_types(
  rtdb_int32 handle,
  rtdb_int32* count,
  char* name[RTDB_TYPE_NAME_SIZE],
  rtdb_int32* field_counts
  );

/**
* 命名：rtdbb_get_named_type
* 功能：获取自定义类型的所有字段
* 参数：
*        [handle]         连接句柄，输入参数
*        [name]           自定义类型的名称，输入参数
*        [field_count]    返回name指定的自定义类型的字段个数，输入/输出参数
*                         输入：指定fields数组长度.输出：实际的name自定义类型的字段的个数
*        [fields]         返回由name所指定的自定义类型所包含字段RTDB_DATA_TYPE_FIELD结构的数组，输入/输出参数
*                         输入：fields数组长度要等于count。输出:RTDB_DATA_TYPE_FIELD结构的数组
*        [type_size]      所有自定义类型fields结构中长度字段的累加和，输出参数
*        [desc]           自定义类型的描述，输出参数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_named_type(
  rtdb_int32 handle,
  const char* name,
  rtdb_int32* field_count,
  RTDB_DATA_TYPE_FIELD* fields,
  rtdb_int32* type_size,
  char desc[RTDB_DESC_SIZE]
  );

/**
* 命名：rtdbb_remove_named_type
* 功能：删除自定义类型
* 参数：
*        [handle]      连接句柄，输入参数
*        [name]        自定义类型的名称，输入参数
*        [reserved]      保留字段,暂时不用
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_remove_named_type(
  rtdb_int32 handle,
  const char* name,
  rtdb_int32 reserved GAPI_DEFAULT_VALUE(0)
  );

/**
* 命名：rtdbb_get_named_type_names_property
* 功能：根据标签点id查询标签点所对应的自定义类型的名字和字段总数
* 参数：
*        [handle]           连接句柄
*        [count]            输入/输出，标签点个数，
*                           输入时表示 ids、named_type_names、field_counts、errors 的长度，
*                           输出时表示成功获取自定义类型名字的标签点个数
*        [ids]              整型数组，输入，标签点标识列表
*        [named_type_names] 字符串数组，输出，标签点自定义类型的名字
*        [field_counts]     整型数组，输出，标签点自定义类型的字段个数
*        [errors]           无符号整型数组，输出，获取自定义类型名字的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、named_type_names、field_counts、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_NAMED_T 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_named_type_names_property(
  rtdb_int32 handle,
  rtdb_int32 *count,
  rtdb_int32 *ids,
  char* const *named_type_names,
  rtdb_int32 *field_counts,
  rtdb_error *errors
  );

/**
* 命名：rtdbb_get_recycled_named_type_names_property
* 功能：根据回收站标签点id查询标签点所对应的自定义类型的名字和字段总数
* 参数：
*        [handle]           连接句柄
*        [count]            输入/输出，标签点个数，
*                           输入时表示 ids、named_type_names、field_counts、errors 的长度，
*                           输出时表示成功获取自定义类型名字的标签点个数
*        [ids]              整型数组，输入，回收站标签点标识列表
*        [named_type_names] 字符串数组，输出，标签点自定义类型的名字
*        [field_counts]     整型数组，输出，标签点自定义类型的字段个数
*        [errors]           无符号整型数组，输出，获取自定义类型名字的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、named_type_names、field_counts、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_NAMED_T 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_recycled_named_type_names_property(
  rtdb_int32 handle,
  rtdb_int32 *count,
  rtdb_int32 *ids,
  char* const *named_type_names,
  rtdb_int32 *field_counts,
  rtdb_error *errors
  );

/**
* 命名：rtdbb_get_named_type_points_count
* 功能：获取该自定义类型的所有标签点个数
* 参数：
*        [handle]           连接句柄，输入参数
*        [name]             自定义类型的名称，输入参数
*        [points_count]     返回name指定的自定义类型的标签点个数，输入参数
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_named_type_points_count(
  rtdb_int32 handle,
  const char* name,
  rtdb_int32 *points_count);


/**
*
* \brief 获取该内置的基本类型的所有标签点个数
*
* \param handle           整型，输入参数，连接句柄
* \param type             整型，输入参数，内置的基本类型，参数的值可以是除RTDB_NAME_T以外的所有RTDB_TYPE枚举值
* \param points_count     整型，输入参数，返回type指定的内置基本类型的标签点个数
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_base_type_points_count(
  rtdb_int32 handle,
  rtdb_int32 type,
  rtdb_int32 *points_count
  );

/**
* 命名：rtdbb_modify_named_type
* 功能：修改自定义类型名称,描述,字段名称,字段描述
* 参数：
*        [handle]             连接句柄，输入参数
*        [name]               自定义类型的名称，输入参数
*        [modify_name]        要修改的自定义类型名称，输入参数
*        [modify_desc]        要修改的自定义类型的描述，输入参数
*        [modify_field_name]  要修改的自定义类型字段的名称，输入参数
*        [modify_field_desc]  要修改的自定义类型字段的描述，输入参数
*        [field_count]        自定义类型字段的个数，输入参数
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_modify_named_type(
  rtdb_int32 handle,
  const char* name,
  const char* modify_name,
  const char* modify_desc,
  const char* modify_field_name[RTDB_TYPE_NAME_SIZE],
  const char* modify_field_desc[RTDB_DESC_SIZE],
  rtdb_int32 field_count);

/**
*
* \brief 获取元数据同步信息
*
* \param handle           整型，输入参数，连接句柄
* \param node_number      整型，输入参数，双活节点id，1表示第一个节点，2表示第二个节点。0表示所有节点
* \param count            整型，输入参数，sync_infos参数的数量
*                              输出参数，输出实际获取到的sync_infos的个数
* \param sync_infos       RTDB_SYNC_INFO数组，输出参数，输出实际获取到的同步信息
* \param errors           rtdb_error数组，输出参数，输出对应节点的错误信息
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_meta_sync_info(
    rtdb_int32 handle,
    rtdb_int32 node_number,
    rtdb_int32* count,
    RTDB_SYNC_INFO* sync_infos,
    rtdb_error* errors);

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_types_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* types,
    rtdb_error* errors
);

/**@}*/


/**
* \defgroup snapshot 实时数据接口
* @{
*/

/**
*
* \brief 批量读取开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param values    双精度浮点型数组，输出，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
* \param states    64 位整型数组，输出，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*                    但它们的时间戳必需是递增的。
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float64 *values,
  const rtdb_int64 *states,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入开关量、模拟量快照数值
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*                    但它们的时间戳必需是递增的。
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
*        其余情况下会调用 rtdbs_put_snapshots()
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_fix_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float64 *values,
  const rtdb_int64 *states,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_fix_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量回溯快照
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表，同一个标签点标识可以出现多次，
*
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param values    双精度浮点型数组，输入，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的快照值；否则忽略
* \param states    64 位整型数组，输入，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的快照值；否则忽略
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*        本接口对数据类型为 RTDB_COOR、RTDB_STRING、RTDB_BLOB 的标签点无效。
* 功能说明：
*       批量将标签点的快照值vtmq改成传入的vtmq，如果传入的时间戳早于当前快照，会删除传入时间戳到当前快照的历史存储值。
*       如果传入的时间戳等于或者晚于当前快照，什么也不做。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_back_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float64 *values,
  const rtdb_int64 *states,
  const rtdb_int16 *qualities,
  rtdb_error *errors
);

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_back_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量读取坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param x         单精度浮点型数组，输出，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输出，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_coor_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float32 *x,
  rtdb_float32 *y,
  rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_coor_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_coor_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float32 *x,
  const rtdb_float32 *y,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_coor_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float32* x,
    const rtdb_float32* y,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量写入坐标实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识列表
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param x         单精度浮点型数组，输入，实时浮点型横坐标数值列表
* \param y         单精度浮点型数组，输入，实时浮点型纵坐标数值列表
* \param qualities 短整型数组，输入，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，写入实时坐标数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*        仅当输入时间戳与当前快照时间戳完全相等时，会替换当前快照的值和质量；
*        其余情况下会调用 rtdbs_put_coor_snapshots()
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_fix_coor_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float32 *x,
  const rtdb_float32 *y,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_fix_coor_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float32* x,
    const rtdb_float32* y,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blob      字节型数组，输出，实时二进制/字符串数值
* \param len       短整型，输出，二进制/字符串数值长度
* \param quality   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshot32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *datetime,
  rtdb_time_type *ms,
  rtdb_byte *blob,
  rtdb_length_type *len,
  rtdb_int16 *quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_byte* blob,
    rtdb_length_type* len,
    rtdb_int16* quality
);

/**
*
* \brief 批量读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blobs     字节型指针数组，输出，实时二进制/字符串数值
* \param lens      短整型数组，输入/输出，二进制/字符串数值长度，
*                    输入时表示对应的 blobs 指针指向的缓冲区长度，
*                    输出时表示实际得到的 blob 长度，如果 blob 的长度大于缓冲区长度，会被截断。
* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshots32(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_byte* const* blobs,
  rtdb_length_type *lens,
  rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_byte* const* blobs,
    rtdb_length_type* lens,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，实时二进制/字符串数值
* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshot32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  const rtdb_byte *blob,
  rtdb_length_type len,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const rtdb_byte* blob,
    rtdb_length_type len,
    rtdb_int16 quality
);

/**
*
* \brief 批量写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param blobs     字节型指针数组，输入，实时二进制/字符串数值
* \param lens      短整型数组，输入，二进制/字符串数值长度，
*                    表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshots32(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_byte* const* blobs,
  const rtdb_length_type *lens,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_byte* const* blobs,
    const rtdb_length_type* lens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量读取datetime类型标签点实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param dtvalues  字节型指针数组，输出，实时datetime数值
* \param dtlens    短整型数组，输入/输出，datetime数值长度，
*                    输入时表示对应的 dtvalues 指针指向的缓冲区长度，
* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \param type      短整型，输入，所有标签点的显示类型，如“yyyy-mm-dd hh:mm:ss.000”的type为1，默认类型1，
*                    “yyyy/mm/dd hh:mm:ss.000”的type为2
*                    如果不传type，则按照标签点属性显示，否则按照type类型显示
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_datetime_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_byte* const* dtvalues,
  rtdb_int16 *dtlens,
  rtdb_int16 *qualities,
  rtdb_error *errors,
  rtdb_int16 type GAPI_DEFAULT_VALUE(-1)
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_datetime_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_byte* const* dtvalues,
    rtdb_length_type* dtlens,
    rtdb_int16* qualities,
    rtdb_error* errors,
    rtdb_int16 type
);

/**
*
* \brief 批量插入datetime类型标签点数据
*
* \param handle      连接句柄
* \param count       整型，输入/输出，标签点个数，
*                      输入时表示 ids、datetimes、ms、dtvalues、dtlens、qualities、errors的长度，
*                      输出时表示成功写入的标签点个数
* \param ids         整型数组，输入，标签点标识
* \param datetimes   整型数组，输入，实时值时间列表
*                      表示距离1970年1月1日08:00:00的秒数
* \param ms          短整型数组，输入，实时数值时间列表，
*                      对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param dtvalues    字节型指针数组，输入，datetime标签点的值
* \param dtlens      短整型数组，输入，数值长度
* \param qualities   短整型数组，输入，实时数值品质，，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors      无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 被接口只对数据类型 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_datetime_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_byte* const* dtvalues,
  const rtdb_int16 *dtlens,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_datetime_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_byte* const* dtvalues,
    const rtdb_length_type* dtlens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量标签点快照改变的通知订阅
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param options        订阅选项
*                           RTDB_O_AUTOCONN 自动重连
* \param param          用户自定义参数
* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                       当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                       作为event_type取值调用回掉函数后退出订阅。
*                       当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                       作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                       网络中断期间回掉函数调用频率为最少3秒
*                       event_type参数值含义如下：
*                         RTDB_E_DATA        标签点快照改变
*                         RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                         RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                         RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
*                       handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
*                       param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
*                       count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
*                              event_type为其它值时，count值为0
*                       ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
*                              event_type为其它值时，ids值为NULL
*                       datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
*                                 event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
*                                 event_type为其它值时，datetimes值为NULL
*                       ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
*                              event_type为其它值时，ms值为NULL
*                       values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
*                              event_type为其它值时，values值为NULL
*                       states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
*                              event_type为其它值时，states值为NULL
*                       qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
*                              event_type为其它值时，qualities值为NULL
*                       errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
*                              event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
*                              event_type为其它值时，errors值为NULL
* \param errors         无符号整型数组，输出，
*                           写入实时数据的返回值列表，参考rtdb_error.h
* \remark   用户须保证 ids、errors 的长度与 count 一致。
*        用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*        否则返回 RtE_OTHER_SDK_DOING 错误。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_subscribe_snapshots_ex(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  rtdb_uint32 options,
  void* param,
  rtdbs_snaps_event_ex callback,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_subscribe_snapshots_ex64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_uint32 options,
    void* param,
    rtdbs_snaps_event_ex64 callback,
    rtdb_error* errors
);

/**
*
* \brief 批量标签点快照改变的通知订阅
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
* \param options        订阅选项
*                           RTDB_O_AUTOCONN 自动重连
* \param param          用户自定义参数
* \param callback       rtdbs_snaps_event_ex 类型回调接口，输入，当回掉函数返回非RtE_OK时退出订阅
*                         当未设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                         作为event_type取值调用回掉函数后退出订阅。
*                         当设置options为RTDB_O_AUTOCONN时，订阅断开后使用RTDB_E_DISCONNECT
*                         作为event_type取值调用回掉函数直到连接恢复或回掉函数返回非RtE_OK，
*                         网络中断期间回掉函数调用频率为最少3秒
*                         event_type参数值含义如下：
*                           RTDB_E_DATA        标签点快照改变
*                           RTDB_E_DISCONNECT  订阅客户端与数据库网络断开
*                           RTDB_E_RECOVERY    订阅客户端与数据库网络及订阅恢复
*                           RTDB_E_CHANGED     客户端修改订阅标签点信息，即通过rtdbs_change_subscribe_snapshots修改订阅信息的结果
*                         handle 产生订阅回掉的连接句柄，调用rtdbs_subscribe_snapshots_ex时的handle参数
*                         param  用户自定义参数，调用rtdbs_subscribe_snapshots_ex时的param参数
*                         count  event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示ids，datetimes,values等的数量
*                                event_type为其它值时，count值为0
*                         ids    event_type为RTDB_E_DATA和RTDB_E_CHANGED时表示快照改变的标签点ID，数量由count指定
*                                event_type为其它值时，ids值为NULL
*                         datetimes event_type为RTDB_E_DATA时表示快照时间，数量由count指定
*                                   event_type为RTDB_E_CHANGED时表示changed_types，即通过rtdbs_change_subscribe_snapshots传入的changed_types
*                                   event_type为其它值时，datetimes值为NULL
*                         ms     event_type为RTDB_E_DATA时表示快照的毫秒，数量由count指定
*                                event_type为其它值时，ms值为NULL
*                         values event_type为RTDB_E_DATA时表示浮点数据类型快照值，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示delta_values，即通过rtdbs_change_subscribe_snapshots传入的delta_values
*                                event_type为其它值时，values值为NULL
*                         states event_type为RTDB_E_DATA时表示整形数据类型快照值，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示delta_states，通过rtdbs_change_subscribe_snapshots传入的delta_states
*                                event_type为其它值时，states值为NULL
*                         qualities event_type为RTDB_E_DATA时表示快照质量码，数量由count指定
*                                event_type为其它值时，qualities值为NULL
*                         errors event_type为RTDB_E_DATA时表示快照错误码，数量由count指定
*                                event_type为RTDB_E_CHANGED时，表示修改结果对应的错误码，数量由count指定
*                                event_type为其它值时，errors值为NULL
* \param errors         无符号整型数组，输出，
*                           写入实时数据的返回值列表，参考rtdb_error.h
* \remark delta_values和delta_states可以为空指针，表示不设置容差值。 只有两个参数都不为空时，设置容差值才会生效。
            用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致
*           用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*           否则返回 RtE_OTHER_SDK_DOING 错误。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_subscribe_delta_snapshots(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_float64* delta_values,
  const rtdb_int64* delta_states,
  rtdb_uint32 options,
  void* param,
  rtdbs_snaps_event_ex callback,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_subscribe_delta_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_float64* delta_values,
    const rtdb_int64* delta_states,
    rtdb_uint32 options,
    void* param,
    rtdbs_snaps_event_ex64 callback,
    rtdb_error* errors
);

/**
*
* \brief 批量修改订阅标签点信息
*
* \param handle         连接句柄
* \param count          整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                           输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids            整型数组，输入，标签点标识列表。
* \param delta_values   double型数组，输入，订阅浮点类型标签点的容差值，变化超过设置的容差值才会推送
* \param delta_values   整型数组，输入，订阅整型标签点的容差值，变化超过设置的容差值才会推送
* \param changed_types  整型数组，输入，修改类型，参考RTDB_SUBSCRIBE_CHANGE_TYPE
* \param errors         异步调用，保留参数，暂时不启用
* \remark   用户须保证 ids、delta_values、delta_states、errors 的长度与 count 一致。
*               可以同时添加、修改、删除订阅的标签点信息，
*               delta_values和delta_states，可以为空指针，为空，则表示不设置容差值，即写入新数据即推送
*               只有delta_values和delta_states都不为空时，设置的容差值才有效。
*               用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*               否则返回 RtE_OTHER_SDK_DOING 错误。
*               此方法是异步方法，当网络中断等异常情况时，会通过方法的返回值返回错误，参考rtdb_error.h。
*               当方法返回值为RtE_OK时，表示已经成功发送给数据库，但是并没有等待修改结果。
*               数据库的修改结果，会异步通知给api的回调函数，通过rtdbs_snaps_event_ex的RTDB_E_CHANGED事件通知修改结果
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_change_subscribe_snapshots(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_float64* delta_values,
    const rtdb_int64* delta_states,
    const rtdb_int32* changed_types,
    rtdb_error* errors
);

/**
*
* \brief 取消标签点快照更改通知订阅
*
* \param handle    连接句柄
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_cancel_subscribe_snapshots(
  rtdb_int32 handle
  );

/**
* 命名：rtdbs_get_named_type_snapshot32
* 功能：获取自定义类型测点的单个快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [object]    字节型数组，输出，实时自定义类型标签点的数值
*        [length]    短整型，输入/输出，自定义类型标签点的数值长度
*        [quality]   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshot32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *datetime,
  rtdb_time_type *ms,
  void* object,
  rtdb_length_type *length,
  rtdb_int16 *quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    void* object,
    rtdb_length_type* length,
    rtdb_int16* quality
);

/**
* 命名：rtdbs_get_named_type_snapshots32
* 功能：批量获取自定义类型测点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [objects]   指针数组，输出，自定义类型标签点数值
*        [lengths]   短整型数组，输入/输出，自定义类型标签点数值长度，
*                    输入时表示对应的 objects 指针指向的缓冲区长度，
*                    输出时表示实际得到的 objects 长度，如果 objects 的长度大于缓冲区长度，会被截断。
*        [qualities] 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshots32(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  void* const* objects,
  rtdb_length_type *lengths,
  rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    void* const* objects,
    rtdb_length_type* lengths,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
* 命名：rtdbs_put_named_type_snapshot32
* 功能：写入单个自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
*        [object]    void类型数组，输入，自定义类型标签点数值
*        [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
*        [quality]   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshot32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  const void *object,
  rtdb_length_type length,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshot64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const void* object,
    rtdb_length_type length,
    rtdb_int16 quality
);

/**
* 命名：rtdbs_put_named_type_snapshots32
* 功能：批量写入自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
*        [objects]   void类型指针数组，输入，自定义类型标签点数值
*        [lengths]   短整型数组，输入，自定义类型标签点数值长度，
*                    表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
*        [qualities] 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshots32(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const void* const* objects,
  const rtdb_length_type *lengths,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshots64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const void* const* objects,
    const rtdb_length_type* lengths,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**@}*/


/**
* \defgroup archive 存档文件接口
* @{
*/

/**
*
* \brief 获取存档文件数量
*
* \param handle    连接句柄
* \param count     整型，输出，存档文件数量
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_count(
  rtdb_int32 handle,
  rtdb_int32 *count
  );

/**
*
* \brief 新建指定时间范围的历史存档文件并插入到历史数据库
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
* \param begin      整数，输入，起始时间，距离1970年1月1日08:00:00的秒数
* \param end        整数，输入，终止时间，距离1970年1月1日08:00:00的秒数
* \param mb_size    整型，输入，文件兆字节大小，单位为 MB。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_create_ranged_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  rtdb_int32 begin,
  rtdb_int32 end,
  rtdb_int32 mb_size
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_create_ranged_archive64(
    rtdb_int32 handle,
    const char* path,
    const char* file,
    rtdb_timestamp_type begin,
    rtdb_timestamp_type end,
    rtdb_int32 mb_size
);


/**
*
* \brief 追加磁盘上的历史存档文件到历史数据库。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名，后缀名应为.rdf。
* \param state      整型，输入，取值 RTDB_ACTIVED_ARCHIVE、RTDB_NORMAL_ARCHIVE、
*                     RTDB_READONLY_ARCHIVE 之一，表示文件状态
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_append_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  rtdb_int32 state
  );

/**
*
* \brief 从历史数据库中移出历史存档文件。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_remove_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );


/**
*
* \brief 切换活动文件
*
* \param handle     连接句柄
* \remark 当前活动文件被写满时该事务被启动，
*        改变当前活动文件的状态为普通状态，
*        在所有历史数据存档文件中寻找未被使用过的
*        插入到前活动文件的右侧并改为活动状态，
*        若找不到则将前活动文件右侧的文件改为活动状态，
*        并将active_archive_指向该文件。该事务进行过程中，
*        用锁保证所有读写操作都暂停等待该事务完成。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_shift_actived(
  rtdb_int32 handle
  );


/**
* 命名：rtdba_get_archives
* 功能：获取存档文件的路径、名称、状态和最早允许写入时间。
* 参数：
*        [handle]          连接句柄
*        [paths]            字符串数组，输出，存档文件的目录路径，长度至少为 RTDB_PATH_SIZE。
*        [files]            字符串数组，输出，存档文件的名称，长度至少为 RTDB_FILE_NAME_SIZE。
*        [states]           整型数组，输出，取值 RTDB_INVALID_ARCHIVE、RTDB_ACTIVED_ARCHIVE、
*                          RTDB_NORMAL_ARCHIVE、RTDB_READONLY_ARCHIVE 之一，表示文件状态
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives(
  rtdb_int32 handle,
  rtdb_int32* count,
  rtdb_path_string* paths,
  rtdb_filename_string* files,
  rtdb_int32 *states
  );


RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_info(
  rtdb_int32 handle,
  rtdb_int32* count,
  const rtdb_path_string* const paths,
  const rtdb_filename_string* const files,
  RTDB_HEADER_PAGE *infos,
  rtdb_error* errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_perf_data(
  rtdb_int32 handle,
  rtdb_int32* count,
  const rtdb_path_string* const paths,
  const rtdb_filename_string* const files,
  RTDB_ARCHIVE_PERF_DATA* real_time_datas,
  RTDB_ARCHIVE_PERF_DATA* total_datas,
  rtdb_error* errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archives_status(rtdb_int32 handle, rtdb_error* status);

/**
*
* \brief 获取存档文件及其附属文件的详细信息。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
* \param file_id    整型，输入，附属文件标识，0 表示获取主文件信息。
* \param info       RTDB_HEADER_PAGE 结构，输出，存档文件信息
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_archive_info(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  rtdb_int32 file_id,
  RTDB_HEADER_PAGE *info
  );


/**
*
* \brief 修改存档文件的可配置项。
*
* \param handle         连接句柄
* \param path           字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file           字符串，输入，文件名。
* \param rated_capacity 整型，输入，文件额定大小，单位为 MB。
* \param ex_capacity    整型，输入，附属文件大小，单位为 MB。
* \param auto_merge     短整型，输入，是否自动合并附属文件。
* \param auto_arrange   短整型，输入，是否自动整理存档文件。
* 备注: rated_capacity 与 ex_capacity 参数可为 0，表示不修改对应的配置项。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_update_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  rtdb_int32 rated_capacity,
  rtdb_int32 ex_capacity,
  rtdb_int16 auto_merge,
  rtdb_int16 auto_arrange
  );

/**
*
* \brief 整理存档文件，将同一标签点的数据块存放在一起以提高查询效率。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_arrange_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );

/**
*
* \brief 为存档文件重新生成索引，用于恢复数据。
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_reindex_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );

/**
*
* \brief 备份主存档文件及其附属文件到指定路径
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名。
* \param dest       字符串，输入，备份目录路径，必须以"\"或"/"结尾。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_backup_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  const char *dest
  );

/**
* 命名：rtdba_move_archive
* 功能：将存档文件移动到指定目录
* 参数：
*        [handle]     连接句柄
*        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
*        [file]       字符串，输入，文件名。
*        [dest]       字符串，输入，移动目录路径，必须以"\"或"/"结尾。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_move_archive(
  rtdb_int32 handle,
  const char *path,
  const char *file,
  const char *dest
  );

/**
* 命名：rtdba_reindex_archive
* 功能：为存档文件转换索引格式。
* 参数：
*        [handle]     连接句柄
*        [path]       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
*        [file]       字符串，输入，文件名。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_convert_index(
  rtdb_int32 handle,
  const char *path,
  const char *file
  );


/**
* 命名：rtdba_query_big_job
* \brief 查询进程正在执行的后台任务类型、状态和进度
*
* \param handle     连接句柄
* \param process    所查询的进程代号，进程的标识参见枚举 RTDB_PROCESS_NAME,
*                     RTDB_PROCESS_HISTORIAN: 历史服务进程，具有以下任务类型：
*                         RTDB_MERGE: 合并附属文件到主文件;
*                         RTDB_ARRANGE: 整理存档文件;
*                         RTDB_REINDEX: 重建索引;
*                         RTDB_BACKUP: 备份;
*                         RTDB_REACTIVE: 激活为活动存档;
*                     RTDB_PROCESS_EQUATION: 方程式服务进程，具有以下任务类型：
*                         RTDB_COMPUTE: 历史计算;
*                     RTDB_PROCESS_BASE: 标签信息服务进程，具有以下任务类型：
*                         RTDB_UPDATE_TABLE: 修改表名称;
*                         RTDB_REMOVE_TABLE: 删除表;
* \param path       字符串，输出，长度至少为 RTDB_PATH_SIZE，
*                     对以下任务，这个字段表示存档文件所在目录路径：
*                         RTDB_MERGE
*                         RTDB_ARRANGE
*                         RTDB_REINDEX
*                         RTDB_BACKUP
*                         RTDB_REACTIVE
*                     对于以下任务，这个字段表示原来的表名：
*                         RTDB_UPDATE_TABLE
*                         RTDB_REMOVE_TABLE
*                     对于其它任务不可用。
* \param file       字符串，输出，长度至少为 RTDB_FILE_NAME_SIZE，
*                     对以下任务，这个字段表示存档文件名：
*                         RTDB_MERGE
*                         RTDB_ARRANGE
*                         RTDB_REINDEX
*                         RTDB_BACKUP
*                         RTDB_REACTIVE
*                     对于以下任务，这个字段表示修改后的表名：
*                          RTDB_UPDATE_TABLE
*                     对于其它任务不可用。
* \param job        短整型，输出，任务的标识参见枚举 RTDB_BIG_JOB_NAME。
* \param state      整型，输出，任务的执行状态，参考 rtdb_error.h
* \param end_time   整型，输出，任务的完成时间。
* \param progress   单精度浮点型，输出，任务的进度百分比。
* \remark path 及 file 参数可传空指针，对应的信息将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_query_big_job(
  rtdb_int32 handle,
  rtdb_int32 process,
  char *path,
  char *file,
  rtdb_int16 *job,
  rtdb_int32 *state,
  rtdb_int32 *end_time,
  rtdb_float32 *progress
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_query_big_job64(
    rtdb_int32 handle,
    rtdb_int32 process,
    char* path,
    char* file,
    rtdb_int16* job,
    rtdb_int32* state,
    rtdb_timestamp_type* end_time,
    rtdb_float32* progress
);

/**
* 命名：rtdba_cancel_big_job
* 功能：取消进程正在执行的后台任务
* 参数：
*        [handle]     连接句柄
*        [process]    所查询的进程代号，进程的标识参见枚举 RTDB_PROCESS_NAME,
*                     RTDB_PROCESS_HISTORIAN: 历史服务进程，具有以下任务类型：
*                         RTDB_MERGE: 合并附属文件到主文件;
*                         RTDB_ARRANGE: 整理存档文件;
*                         RTDB_REINDEX: 重建索引;
*                         RTDB_BACKUP: 备份;
*                         RTDB_REACTIVE: 激活为活动存档;
*                     RTDB_PROCESS_EQUATION: 方程式服务进程，具有以下任务类型：
*                         RTDB_COMPUTE: 历史计算;
*                     RTDB_PROCESS_BASE: 标签信息服务进程，具有以下任务类型：
*                         RTDB_UPDATE_TABLE: 修改表名称;
*                         RTDB_REMOVE_TABLE: 删除表;
* 备注：path 及 file 参数可传空指针，对应的信息将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdba_cancel_big_job(rtdb_int32 handle, rtdb_int32 process);



/**@}*/



/**
* \defgroup historian 历史数据接口
* @{
*/

/**
*
* \brief 获取单个标签点在一段时间范围内的存储值数量.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整型，输出，返回上述时间范围内的存储值数量
* \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_archived_values_count(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  rtdb_int32 *count
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_archived_values_count64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count
);

/**
*
* \brief 获取单个标签点在一段时间范围内的真实的存储值数量.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整型，输出，返回上述时间范围内的存储值数量
* \remark 由 datetime1、ms1 形成的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_archived_values_real_count(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  rtdb_int32 *count
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_archived_values_real_count64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count
);

/**
*
* \brief 读取单个标签点一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，历史浮点型数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
* \param states        64 位整数数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);


/**
*
* \brief 逆向读取单个标签点一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，历史浮点型数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
* \param states        64 位整数数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_backward(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_backward64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点一段时间内的坐标型储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、x、y、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
* \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_coor_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float32 *x,
  rtdb_float32 *y,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_coor_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities
);

/**
*
* \brief 逆向读取单个标签点一段时间内的坐标型储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、x、y、qualities 的长度；
*                        输出时返回实际得到的数值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param x             单精度浮点型数组，输出，浮点型横坐标历史数值列表
* \param y             单精度浮点型数组，输出，浮点型纵坐标历史数值列表
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_coor_values_backward(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float32 *x,
  rtdb_float32 *y,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_coor_values_backward64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities
);


/**
*
* \brief 开始以分段返回方式读取一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整型，输出，返回上述时间范围内的存储值数量
* \param batch_count   整型，输出，每次分段返回的长度，用于继续调用 rtdbh_get_next_archived_values 接口
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_in_batches(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  rtdb_int32 *count,
  rtdb_int32 *batch_count
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_in_batches64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count,
    rtdb_int32* batch_count
);

/**
*
* \brief 分段读取一段时间内的储存数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整形，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度；
*                        输出时表示实际得到的存储值个数。
* \param datetimes     整型数组，输出，历史数值时间列表,
*                        表示距离1970年1月1日08:00:00的秒数
* \param ms            短整型数组，输出，历史数值时间列表，
*                        对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param values        双精度浮点型数组，输出，历史浮点型数值列表，
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史存储值；否则为 0
* \param states        64 位整型数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
*        且 count 不能小于 rtdbh_get_archived_values_in_batches 接口中返回的 batch_count 的值，
*        当返回 RtE_BATCH_END 表示全部数据获取完毕。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_next_archived_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_next_archived_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点的单调递增时间序列历史插值。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度。
* \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
*                        为距离1970年1月1日08:00:00的秒数
* \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
*                        表示需要的单调递增时间对应的纳秒值；否则忽略。
* \param values        双精度浮点型数组，输出，历史浮点型数值列表，
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的历史插值；否则为 0
* \param states        64 位整型数组，输出，历史整型数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 相符，
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_timed_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 count,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_timed_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 count,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个坐标标签点的单调递增时间序列历史插值。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入，表示 datetimes、ms、x、y、qualities 的长度。
* \param datetimes     整型数组，输入，表示需要的单调递增时间列表，
*                        为距离1970年1月1日08:00:00的秒数
* \param ms            短整型数组，输入，对于时间精度为纳秒的标签点，
*                        表示需要的单调递增时间对应的纳秒值；否则忽略。
* \param x             单精度浮点型数组，输出，浮点型横坐标历史插值数值列表
* \param y             单精度浮点型数组，输出，浮点型纵坐标历史插值数值列表
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、x、y、qualities 的长度与 count 相符，
*        本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_timed_coor_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 count,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  rtdb_float32 *x,
  rtdb_float32 *y,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_timed_coor_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 count,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点一段时间内等间隔历史插值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数；输出时返回实际得到的插值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interpo_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interpo_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时刻之后一定数量的等间隔内插值替换的历史数值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param interval      整型，输入，插值时间间隔，单位为纳秒
* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素用于存放起始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interval_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int64 interval,
  rtdb_int32 count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interval_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int64 interval,
    rtdb_int32 count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时间的历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*                        RTDB_INTER 取指定时间的内插值数据。
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param value         双精度浮点数，输出，浮点型历史数值
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
* \param state         64 位整数，输出，整型历史数值，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 mode,
  rtdb_int32 *datetime,
  rtdb_time_type *ms,
  rtdb_float64 *value,
  rtdb_int64 *state,
  rtdb_int16 *quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_float64* value,
    rtdb_int64* state,
    rtdb_int16* quality
);

/**
*
* \brief 读取单个标签点某个时间的坐标型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*                        RTDB_INTER 取指定时间的内插值数据。
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param x             单精度浮点型，输出，横坐标历史数值
* \param y             单精度浮点型，输出，纵坐标历史数值
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_coor_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 mode,
  rtdb_int32 *datetime,
  rtdb_time_type *ms,
  rtdb_float32 *x,
  rtdb_float32 *y,
  rtdb_int16 *quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_coor_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_float32* x,
    rtdb_float32* y,
    rtdb_int16* quality
);

/**
*
* \brief 读取单个标签点某个时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param blob          字节型数组，输出，二进制/字符串历史值
* \param len           短整型，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_blob_value32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 mode,
  rtdb_int32 *datetime,
  rtdb_time_type *ms,
  rtdb_byte *blob,
  rtdb_length_type *len,
  rtdb_int16 *quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_blob_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_byte* blob,
    rtdb_length_type* len,
    rtdb_int16* quality
);

/**
*
* \brief 读取单个标签点一段时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param count         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
* \param datetime1     整型，输入，表示开始时间秒数；
* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetime2     整型，输入,表示结束时间秒数；
* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
* \param lens          短整型数组，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
*                        当blobs为空指针时，表示只获取每条数据的长度，此时会忽略输入的lens
* \param blobs         字节型数组，输出，二进制/字符串历史值。可以设置为空指针，表示只获取每条数据的长度
* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_blob_values32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_length_type *lens,
  rtdb_byte* const* blobs,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_blob_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_length_type* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities
);

/**
*
* \brief 读取并模糊搜索单个标签点一段时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param count         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
* \param datetime1     整型，输入，表示开始时间秒数；
* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetime2     整型，输入,表示结束时间秒数；
* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param filter        字符串，输入，支持通配符的模糊搜索字符串，多个模糊搜索的条件通过空格分隔，只针对string类型有效
*                        当filter为空指针时，表示不进行过滤,
*                        限制最大长度为RTDB_EQUATION_SIZE-1，超过此长度会返回错误
* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
* \param lens          短整型数组，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
*                        当blobs为空指针时，表示只获取每条数据的长度，此时会忽略输入的lens
* \param blobs         字节型数组，输出，二进制/字符串历史值。可以设置为空指针，表示只获取每条数据的长度
* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_blob_values_filt(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_int32 datetime1,
    rtdb_time_type ms1,
    rtdb_int32 datetime2,
    rtdb_time_type ms2,
    const char* filter,
    rtdb_int32* datetimes,
    rtdb_time_type* ms,
    rtdb_length_type* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities
);

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_blob_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    const char* filter,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_length_type* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时间的datetime历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param dtblob          字节型数组，输出，datetime历史值
* \param dtlen           短整型，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的datetime数据长度。
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param type           短整型 datetime字符串的格式类型，默认为-1
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_datetime_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 mode,
  rtdb_int32 *datetime,
  rtdb_time_type *ms,
  rtdb_byte *dtblob,
  rtdb_int16 *dtlen,
  rtdb_int16 *quality,
  rtdb_int16 type GAPI_DEFAULT_VALUE(-1)
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_datetime_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    rtdb_byte* dtblob,
    rtdb_length_type* dtlen,
    rtdb_int16* quality,
    rtdb_int16 type
);

/**
*
* \brief 读取单个标签点一段时间的时间类型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param count         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
* \param datetime1     整型，输入，表示开始时间秒数；
* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetime2     整型，输入,表示结束时间秒数；
* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
* \param dtlens          短整型数组，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的时间数据长度。
* \param dtvalues         字节型数组，输出，时间历史值
* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param type          短整型，输入，“yyyy-mm-dd hh:mm:ss.000”的type为1， 同样默认输入格式也为 “yyyy-mm-dd hh:mm:ss.000”
*                       “yyyy/mm/dd hh:mm:ss.000”的type为2
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_datetime_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_int16 *lens,
  rtdb_byte* const* blobs,
  rtdb_int16 *qualities,
  rtdb_int16 type GAPI_DEFAULT_VALUE(-1)
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_datetime_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_length_type* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities,
    rtdb_int16 type
);


/**
*
* \brief 写入批量标签点批量时间型历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、dtlens、dtvalues、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param dtvalues      字节型指针数组，输入，实时时间数值
* \param dtlens        短整型数组，输入，时间数值长度，
*                        表示对应的 dtvalues 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、dtlens、dtvalues、qualities、errors 的长度与 count 一致，
*        本接口仅对数据类型为 RTDB_DATETIME 的标签点有效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_datetime_values(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_byte* const* dtvalues,
  const rtdb_int16 *dtlens,
  const rtdb_int16 *qualities,
  rtdb_error* errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_datetime_values64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_byte* const* dtvalues,
    const rtdb_length_type* dtlens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 获取单个标签点一段时间内的统计值。
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param datetime1         整型，输入/输出，输入时表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回最大值的时间秒数。
* \param ms1               短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            表示起始时间对应的纳秒，输出时表示最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetime2         整型，输入/输出，输入时表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回最小值的时间秒数。
* \param ms2               短整型，如果 id 指定的标签点时间精度为纳秒，
*                            表示结束时间对应的纳秒，输出时表示最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_value         双精度浮点型，输出，表示统计时间段内的最大数值。
* \param min_value         双精度浮点型，输出，表示统计时间段内的最小数值。
* \param total_value       双精度浮点型，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avg          双精度浮点型，输出，表示统计时间段内的算术平均值。
* \param power_avg         双精度浮点型，输出，表示统计时间段内的加权平均值。
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *datetime1,
  rtdb_time_type *ms1,
  rtdb_int32 *datetime2,
  rtdb_time_type *ms2,
  rtdb_float64 *max_value,
  rtdb_float64 *min_value,
  rtdb_float64 *total_value,
  rtdb_float64 *calc_avg,
  rtdb_float64 *power_avg
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_data
);

/**
* 命名：rtdbh_summary_in_batches
* \brief 分批获取单一标签点一段时间内的统计值
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
*                            max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
*                            即分段的个数；输出时表示成功取得统计值的分段个数。
* \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
*                            如果为纳秒点，输入时间必须大于1纳秒，如果为秒级点，则必须大于1000000000纳秒。
* \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回各个分段对应的最大值的时间秒数。
* \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示起始时间对应的纳秒，
*                            输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回各个分段对应的最小值的时间秒数。
* \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示结束时间对应的纳秒，
*                            输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
* \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
* \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
* \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
* \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
* \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_in_batches(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count,
  rtdb_int64 interval,
  rtdb_int32 *datetimes1,
  rtdb_time_type *ms1,
  rtdb_int32 *datetimes2,
  rtdb_time_type *ms2,
  rtdb_float64 *max_values,
  rtdb_float64 *min_values,
  rtdb_float64 *total_values,
  rtdb_float64 *calc_avgs,
  rtdb_float64 *power_avgs,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data_in_batches(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_int64 interval,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_datas,
    rtdb_error* errors
);

/**
*
* \brief 获取单个标签点一段时间内用于绘图的历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param interval      整型，输入，时间区间数量，单位为个，
*                        一般会使用绘图的横轴(时间轴)所用屏幕像素数，
*                        该功能将起始至结束时间等分为 interval 个区间，
*                        并返回每个区间的第一个和最后一个数值、最大和最小数值、一条异常数值；
*                        故参数 count 有可能输出五倍于 interval 的历史值个数，
*                        所以推荐输入的 count 至少是 interval 的五倍。
* \param count         整型，输入/输出，输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要获取的最大历史值个数，输出时返回实际得到的历史值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
* \param states        64 位整数数组，输出，整型历史值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
* \param qualities     短整型数组，输出，历史值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，用以存放起始及结束时间。
*        第一个元素形成的时间可以大于最后一个元素形成的时间，
*        此时第一个元素表示结束时间，最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_plot_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 interval,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_plot_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 interval,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取批量标签点在某一时间的历史断面数据
*
* \param handle        连接句柄
* \param ids           整型数组，输入，标签点标识列表
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT、RTDB_INTER 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*                        RTDB_INTER 取指定时间的内插值数据。
* \param count         整型，输入，表示 ids、datetimes、ms、values、states、qualities 的长度，即标签点个数。
* \param datetimes     整型数组，输入/输出，输入时表示对应标签点的历史数值时间秒数，
*                        输出时表示根据 mode 实际寻找到的数值时间秒数。
* \param ms            短整型数组，输入/输出，对于时间精度为纳秒的标签点，
*                        输入时表示历史数值时间纳秒数，存放相应的纳秒值，
*                        输出时表示根据 mode 实际寻找到的数值时间纳秒数；否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史值；否则为 0
* \param states        64 位整数数组，输出，整型历史值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史值；否则为 0
* \param qualities     短整型数组，输出，历史值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，读取历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities 的长度与 count 一致，
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_cross_section_values(
  rtdb_int32 handle,
  const rtdb_int32 *ids,
  rtdb_int32 mode,
  rtdb_int32 count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_cross_section_values64(
    rtdb_int32 handle,
    const rtdb_int32* ids,
    rtdb_int32 mode,
    rtdb_int32 count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
* 命名：rtdbh_get_archived_values_filt
* 功能：读取单个标签点在一段时间内经复杂条件筛选后的历史储存值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                        长度不得超过 RTDB_EQUATION_SIZE，为 0 则不进行条件筛选。
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的数值个数；输出时返回实际得到的数值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史存储值；否则为 0
* \param states        64 位整数数组，输出，整型历史数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史存储值；否则为 0
* \param qualities     短整型数组，输出，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_filt(
  rtdb_int32 handle,
  rtdb_int32 id,
  const char *filter,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 读取单个标签点某个时刻之后经复杂条件筛选后一定数量的等间隔内插值替换的历史数值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                        长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param interval      整型，输入，插值时间间隔，单位为纳秒
* \param count         整型，输入，表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数。
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素用于表示起始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interval_values_filt(
  rtdb_int32 handle,
  rtdb_int32 id,
  const char *filter,
  rtdb_int64 interval,
  rtdb_int32 count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interval_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int64 interval,
    rtdb_int32 count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点一段时间内经复杂条件筛选后的等间隔插值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param filter        字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                        长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param count         整型，输入/输出，
*                        输入时表示 datetimes、ms、values、states、qualities 的长度，
*                        即需要的插值个数；输出时返回实际得到的插值个数
* \param datetimes     整型数组，输入/输出，
*                        输入时第一个元素表示起始时间秒数，
*                        最后一个元素表示结束时间秒数，如果为 0，表示直到数据的最后时间；
*                        输出时表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时第一个元素表示起始时间纳秒，
*                        最后一个元素表示结束时间纳秒；
*                        输出时表示对应的历史数值时间纳秒。
*                        否则忽略输入，输出时为 0。
* \param values        双精度浮点数数组，输出，浮点型历史插值数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放相应的历史插值；否则为 0
* \param states        64 位整数数组，输出，整型历史插值数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放相应的历史插值；否则为 0
* \param qualities     短整型数组，输出，历史插值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 用户须保证 datetimes、ms、values、states、qualities 的长度与 count 一致，
*        在输入时，datetimes、ms 中至少应有一个元素，第一个元素形成的时间可以
*        大于最后一个元素形成的时间，此时第一个元素表示结束时间，
*        最后一个元素表示开始时间。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interpo_values_filt(
  rtdb_int32 handle,
  rtdb_int32 id,
  const char *filter,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  rtdb_float64 *values,
  rtdb_int64 *states,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_interpo_values_filt64(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities
);

/**
*
* \brief 获取单个标签点一段时间内经复杂条件筛选后的统计值
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                            长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param datetime1         整型，输入/输出，输入时表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回最大值的时间秒数。
* \param ms1               短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            表示起始时间对应的纳秒，输出时表示最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetime2         整型，输入/输出，输入时表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回最小值的时间秒数。
* \param ms2               短整型，如果 id 指定的标签点时间精度为纳秒，
*                            表示结束时间对应的纳秒，输出时表示最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_value         双精度浮点型，输出，表示统计时间段内的最大数值。
* \param min_value         双精度浮点型，输出，表示统计时间段内的最小数值。
* \param total_value       双精度浮点型，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avg          双精度浮点型，输出，表示统计时间段内的算术平均值。
* \param power_avg         双精度浮点型，输出，表示统计时间段内的加权平均值。
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_filt(
  rtdb_int32 handle,
  rtdb_int32 id,
  const char *filter,
  rtdb_int32 *datetime1,
  rtdb_time_type *ms1,
  rtdb_int32 *datetime2,
  rtdb_time_type *ms2,
  rtdb_float64 *max_value,
  rtdb_float64 *min_value,
  rtdb_float64 *total_value,
  rtdb_float64 *calc_avg,
  rtdb_float64 *power_avg
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data_filt(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_data
);

/**
* 命名：rtdbh_summary_filt_in_batches
* 功能：分批获取单一标签点一段时间内经复杂条件筛选后的统计值
*
* \param handle            连接句柄
* \param id                整型，输入，标签点标识
* \param filter            字符串，输入，由算术、逻辑运算符组成的复杂条件表达式，
*                            长度不得超过 RTDB_EQUATION_SIZE，长度为 0 则不进行条件筛选。
* \param count             整形，输入/输出，输入时表示 datatimes1、ms1、datatimes2、ms2、
*                            max_values、min_values、total_values、calc_avgs、power_avgs、errors 的长度，
*                            即分段的个数；输出时表示成功取得统计值的分段个数。
* \param interval          64 位整型，输入，分段时间间隔，单位为纳秒。
* \param datetimes1        整型数组，输入/输出，输入时第一个元素表示起始时间秒数。
*                            如果为 0，表示从存档中最早时间的数据开始进行统计。
*                            输出时返回各个分段对应的最大值的时间秒数。
* \param ms1               短整型数组，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示起始时间对应的纳秒，
*                            输出时返回各个分段对应的最大值的时间纳秒数；否则忽略，返回值为 0
* \param datetimes2        整型数组，输入/输出，输入时第一个元素表示结束时间秒数。
*                            如果为 0，表示统计到存档中最近时间的数据为止。
*                            输出时返回各个分段对应的最小值的时间秒数。
* \param ms2               短整型数组，如果 id 指定的标签点时间精度为纳秒，
*                            第一个元素表示结束时间对应的纳秒，
*                            输出时返回各个分段对应的最小值的时间纳秒数；否则忽略，返回值为 0
* \param max_values        双精度浮点型数组，输出，表示统计时间段内的最大数值。
* \param min_values        双精度浮点型数组，输出，表示统计时间段内的最小数值。
* \param total_values      双精度浮点型数组，输出，表示统计时间段内的累计值，结果的单位为标签点的工程单位。
* \param calc_avgs         双精度浮点型数组，输出，表示统计时间段内的算术平均值。
* \param power_avgs        双精度浮点型数组，输出，表示统计时间段内的加权平均值。
* \param errors            无符号整型数组，输出，表示各个分段取得统计值的返回值。
* \remark 由 datetimes1[0]、ms1[0] 表示的时间可以大于 datetimes2[0]、ms2[0] 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*        如果输出的最大值或最小值的时间戳秒值为 0，
*        则表明仅有累计值和加权平均值输出有效，其余统计结果无效。
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_filt_in_batches(
  rtdb_int32 handle,
  rtdb_int32 id,
  const char *filter,
  rtdb_int32 *count,
  rtdb_int64 interval,
  rtdb_int32 *datetimes1,
  rtdb_time_type *ms1,
  rtdb_int32 *datetimes2,
  rtdb_time_type *ms2,
  rtdb_float64 *max_values,
  rtdb_float64 *min_values,
  rtdb_float64 *total_values,
  rtdb_float64 *calc_avgs,
  rtdb_float64 *power_avgs,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_summary_data_filt_in_batches(
    rtdb_int32 handle,
    rtdb_int32 id,
    const char* filter,
    rtdb_int32* count,
    rtdb_int64 interval,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    RTDB_SUMMARY_DATA* summary_datas,
    rtdb_error* errors
);

/**
*
* \brief 修改单个标签点某一时间的历史存储值.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param value         双精度浮点数，输入，浮点型历史数值
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放新的历史值；否则忽略
* \param state         64 位整数，输入，整型历史数值，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放新的历史值；否则忽略
* \param quality       短整型，输入，新的历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_update_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  rtdb_float64 value,
  rtdb_int64 state,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_update_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float64 value,
    rtdb_int64 state,
    rtdb_int16 quality
);

/**
*
* \brief 修改单个标签点某一时间的历史存储值.
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param x             单精度浮点型，输入，新的横坐标历史数值
* \param y             单精度浮点型，输入，新的纵坐标历史数值
* \param quality       短整型，输入，新的历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口仅对数据类型为 RTDB_COOR 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_update_coor_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  rtdb_float32 x,
  rtdb_float32 y,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_update_coor_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float32 x,
    rtdb_float32 y,
    rtdb_int16 quality
);


/**
*
* \brief 删除单个标签点某个时间的历史存储值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_remove_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_remove_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime
);

/**
*
* \brief 删除单个标签点一段时间内的历史存储值
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime1     整型，输入，表示起始时间秒数。如果为 0，表示从存档中最早时间的数据开始读取
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示读取直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param count         整形，输出，表示删除的历史值个数
* \remark 由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_remove_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  rtdb_int32 *count
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_remove_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_int32* count
);

/**
*
* \brief 写入单个标签点在某一时间的历史数据。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param value         双精度浮点数，输入，浮点型历史数值
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，存放历史值；否则忽略
* \param state         64 位整数，输入，整型历史数值，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，存放历史值；否则忽略
* \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  rtdb_float64 value,
  rtdb_int64 state,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float64 value,
    rtdb_int64 state,
    rtdb_int16 quality
);

/**
*
* \brief 写入单个标签点在某一时间的坐标型历史数据。
*
* \param handle              连接句柄
* \param id            整型，输入，标签点标识
* \param datetime      整型，输入，时间秒数
* \param ms            短整型，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；否则忽略。
* \param x             单精度浮点型，输入，横坐标历史数值
* \param y             单精度浮点型，输入，纵坐标历史数值
* \param quality       短整型，输入，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_coor_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  rtdb_float32 x,
  rtdb_float32 y,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_coor_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    rtdb_float32 x,
    rtdb_float32 y,
    rtdb_int16 quality
);

/**
*
* \brief 写入单个二进制/字符串标签点在某一时间的历史数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，历史二进制/字符串数值
* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_blob_value32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  const rtdb_byte *blob,
  rtdb_length_type len,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_blob_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const rtdb_byte* blob,
    rtdb_length_type len,
    rtdb_int16 quality
);

/**
*
* \brief 写入批量标签点批量历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、values、states、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识，同一个标签点标识可以出现多次，
*                        但它们的时间戳必需是递增的。
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param values        双精度浮点数数组，输入，浮点型历史数值列表
*                        对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，表示相应的历史存储值；否则忽略
* \param states        64 位整数数组，输入，整型历史数值列表，
*                        对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                        RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，表示相应的历史存储值；否则忽略
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致，
*        本接口对数据类型为 RTDB_COOR、RTDB_BLOB、RTDB_STRING 的标签点无效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_values(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float64 *values,
  const rtdb_int64 *states,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_values64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float64* values,
    const rtdb_int64* states,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 写入批量标签点批量坐标型历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、x、y、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param x             单精度浮点型数组，输入，浮点型横坐标历史数值列表
* \param y             单精度浮点型数组，输入，浮点型纵坐标历史数值列表
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、x、y、qualities、errors 的长度与 count 一致，
*        本接口仅对数据类型为 RTDB_COOR 的标签点有效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_coor_values(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_float32 *x,
  const rtdb_float32 *y,
  const rtdb_int16 *qualities,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_coor_values64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_float32* x,
    const rtdb_float32* y,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 写入单个datetime标签点在某一时间的历史数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，历史datetime数值
* \param len       短整型，输入，datetime数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_DATETIME 的标签点有效。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_datetime_value(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  const rtdb_byte *blob,
  rtdb_int16 len,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_datetime_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const rtdb_byte* blob,
    rtdb_length_type len,
    rtdb_int16 quality
);

/**
*
* \brief 写入批量标签点批量字符串型历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param blobs         字节型指针数组，输入，实时二进制/字符串数值
* \param lens          短整型数组，输入，二进制/字符串数值长度，
*                        表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、lens、blobs、qualities、errors 的长度与 count 一致，
*        本接口仅对数据类型为 RTDB_STRING、RTDB_BLOB 的标签点有效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_blob_values32(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const rtdb_byte* const* blobs,
  const rtdb_length_type *lens,
  const rtdb_int16 *qualities,
  rtdb_error* errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_blob_values64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const rtdb_byte* const* blobs,
    const rtdb_length_type* lens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 将标签点未写满的补历史缓存页写入存档文件中。
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识。
* \param count         整型，输出，缓存页中数据个数。
* \remark 补历史缓存页写满后会自动写入存档文件中，不满的历史缓存页也会写入文件，
*      但会有一个时间延迟，在此期间此段数据可能查询不到，为了及时看到补历史的结果，
*      应在结束补历史后调用本接口。
*      count 参数可为空指针，对应的信息将不再返回。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_flush_archived_values(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 *count
  );

/**
* 命名：rtdbh_get_single_named_type_value32
* 功能：读取单个自定义类型标签点某个时间的历史数据
* 参数：
*        [handle]        连接句柄
*        [id]            整型，输入，标签点标识
*        [mode]          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*        [datetime]      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
*        [ms]            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
*        [object]        void数组，输出，自定义类型标签点历史值
*        [length]        短整型，输入/输出，输入时表示 object 的长度，
*                        输出时表示实际获取的自定义类型标签点数据长度。
*        [quality]       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_named_type_value32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 mode,
  rtdb_int32 *datetime,
  rtdb_time_type *ms,
  void *object,
  rtdb_length_type *length,
  rtdb_int16 *quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_named_type_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_timestamp_type* datetime,
    rtdb_subtime_type* subtime,
    void* object,
    rtdb_length_type* length,
    rtdb_int16* quality
);

/**
* 命名：rtdbh_get_archived_named_type_values32
* 功能：连续读取自定义类型标签点的历史数据
* 参数：
*        [handle]        连接句柄
*        [id]            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*        [datetime1]     整型，输入，表示开始时间秒数；
*        [ms1]           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
*        [datetime2]     整型，输入,表示结束时间秒数；
*        [ms2]           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
*        [length]        短整型数组，输入，输入时表示 objects 的长度，
*        [count]         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
*        [datetimes]     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
*        [ms]            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
*        [objects]       void类型数组，输出，自定义类型标签点历史值
*        [qualities]     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_named_type_values32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  rtdb_length_type length,
  rtdb_int32 *count,
  rtdb_int32 *datetimes,
  rtdb_time_type *ms,
  void* const* objects,
  rtdb_int16 *qualities
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_named_type_values64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    rtdb_length_type length,
    rtdb_int32* count,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    void* const* objects,
    rtdb_int16* qualities
);

/**
* 命名：rtdbh_put_single_named_type_value32
* 功能：写入自定义类型标签点的单个历史事件
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
*        [object]    void数组，输入，历史自定义类型标签点数值
*        [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
*        [quality]   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_named_type_value32(
  rtdb_int32 handle,
  rtdb_int32 id,
  rtdb_int32 datetime,
  rtdb_time_type ms,
  const void *object,
  rtdb_length_type length,
  rtdb_int16 quality
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_named_type_value64(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_timestamp_type datetime,
    rtdb_subtime_type subtime,
    const void* object,
    rtdb_length_type length,
    rtdb_int16 quality
);

/**
* 命名：rtdbh_put_archived_named_type_values32
* 功能：批量补写自定义类型标签点的历史事件
* 参数：
*        [handle]        连接句柄
*        [count]         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
*        [ids]           整型数组，输入，标签点标识
*        [datetimes]     整型数组，输入，表示对应的历史数值时间秒数。
*        [ms]            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
*        [objects]       void类型指针数组，输入，自定义类型标签点数值
*        [lengths]       短整型数组，输入，自定义类型标签点数值长度，
*                        表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
*        [qualities]     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、datetimes、ms、lens、objects、qualities、errors 的长度与 count 一致，
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_named_type_values32(
  rtdb_int32 handle,
  rtdb_int32 *count,
  const rtdb_int32 *ids,
  const rtdb_int32 *datetimes,
  const rtdb_time_type *ms,
  const void* const* objects,
  const rtdb_length_type *lengths,
  const rtdb_int16 *qualities,
  rtdb_error* errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_named_type_values64(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_timestamp_type* datetimes,
    const rtdb_subtime_type* subtimes,
    const void* const* objects,
    const rtdb_length_type* lengths,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**@}*/


/**
* \defgroup equation 方程式计算接口
* @{
*/

/**
*
* \brief 重算或补算批量计算标签点历史数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、errors 的长度，
*                        即标签点个数；输出时返回成功开始计算的标签点个数
* \param flag          短整型，输入，不为 0 表示进行重算，删除时间范围内已经存在历史数据；
*                        为 0 表示补算，保留时间范围内已经存在历史数据，覆盖同时刻的计算值。
* \param datetime1     整型，输入，表示起始时间秒数。
* \param ms1           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示起始时间对应的纳秒；否则忽略
* \param datetime2     整型，输入，表示结束时间秒数。如果为 0，表示计算直至存档中数据的最后时间
* \param ms2           短整型，输入，如果 id 指定的标签点时间精度为纳秒，表示结束时间对应的纳秒；否则忽略
* \param ids           整型数组，输入，标签点标识
* \param errors        无符号整型数组，输出，计算历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、errors 的长度与 count 一致，本接口仅对带有计算扩展属性的标签点有效。
*        由 datetime1、ms1 表示的时间可以大于 datetime2、ms2 表示的时间，
*        此时前者表示结束时间，后者表示起始时间。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_compute_history(
  rtdb_int32 handle,
  rtdb_int32 *count,
  rtdb_int16 flag,
  rtdb_int32 datetime1,
  rtdb_time_type ms1,
  rtdb_int32 datetime2,
  rtdb_time_type ms2,
  const rtdb_int32 *ids,
  rtdb_error *errors
  );

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_compute_history64(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int16 flag,
    rtdb_timestamp_type datetime1,
    rtdb_subtime_type subtime1,
    rtdb_timestamp_type datetime2,
    rtdb_subtime_type subtime2,
    const rtdb_int32* ids,
    rtdb_error* errors
);

/**
* 命名：rtdbe_get_equation_graph_count
* 功能：根据标签点 id 获取相关联方程式键值对数量
* 参数：
*      [handle]   连接句柄
*      [id]       整型，输入，标签点标识
*      [flag]     枚举，输入，获取的拓扑图的关系
*      [count]    整型，输入，拓扑图键值对数量
* 备注：键值对为数据结构，存储方程式涉及到的各标签点ID、及其父ID等
*		具体参考rtdbe_get_equation_graph_datas
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_get_equation_graph_count(
  rtdb_int32 handle,
  rtdb_int32 id,
  RTDB_GRAPH_FLAG flag,
  rtdb_int32 *count
  );

/**
* 命名：rtdbe_get_equation_graph_datas
* 功能：根据标签点 id 获取相关联方程式键值对数据
* 参数：
*      [handle]   连接句柄
*      [id]       整型，输入，标签点标识
*      [flag]     枚举，输入，获取的拓扑图的关系
*      [count]    整型，输出
                    输入时，表示拓扑图键值对数量
                    输出时，表示实际获取到的拓扑图键值对数量
*      [graph]    输出，GOLDE_GRAPH数据结构，拓扑图键值对信息
* 备注：键值对为数据结构，存储方程式涉及到的各标签点ID、及其父ID等
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbe_get_equation_graph_datas(
  rtdb_int32 handle,
  rtdb_int32 id,
  RTDB_GRAPH_FLAG flag,
  rtdb_int32 *count,
  RTDB_GRAPH *graph
  );

/**@}*/


/**
* \defgroup perf 查询性能计数信息接口
* @{
*/

/**
* 命名：rtdbp_get_perf_tags_count
* 功能：获取Perf服务中支持的性能计数点的数量
* 参数：
*      [handle]   连接句柄
*      [count]    整型，输出，表示实际获取到的Perf服务中支持的性能计数点的数量
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbp_get_perf_tags_count(
  rtdb_int32 handle,
  int* count);


/**
* 命名：rtdbp_get_perf_tags_info
* 功能：根据性能计数点ID获取相关的性能计数点信息
* 参数：
*      [handle]   连接句柄
*      [count]    整型，输入，输出
                    输入时，表示想要获取的性能计数点信息的数量，也表示tags_info，errors等的长度
                    输出时，表示实际获取到的性能计数点信息的数量
       [errors] 无符号整型数组，输出，获取性能计数点信息的返回值列表，参考rtdb_error.h
* 备注：用户须保证分配给 tags_info，errors 的空间与 count 相符
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbp_get_perf_tags_info(
  rtdb_int32 handle,
  rtdb_int32* count,
  RTDB_PERF_TAG_INFO* tags_info,
  rtdb_error* errors);


/**
* 命名：rtdbp_get_perf_values
* 功能：批量读取性能计数点的当前快照数值
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，性能点个数，
*                    输入时表示 perf_ids、datetimes、ms、values、states、qualities、errors 的长度，
*                    输出时表示成功获取实时值的性能计数点个数
*        [perf_ids]  整型数组，输入，性能计数点标识列表，参考RTDB_PERF_TAG_ID
*        [datetimes] 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [values]    双精度浮点型数组，输出，实时浮点型数值列表，
*                    对于数据类型为 RTDB_REAL16、RTDB_REAL32、RTDB_REAL64 的标签点，返回相应的快照值；否则为 0
*        [states]    64 位整型数组，输出，实时整型数值列表，
*                    对于数据类型为 RTDB_BOOL、RTDB_UINT8、RTDB_INT8、RTDB_CHAR、RTDB_UINT16、RTDB_INT16、
*                    RTDB_UINT32、RTDB_INT32、RTDB_INT64 的标签点，返回相应的快照值；否则为 0
*        [qualities] 短整型数组，输出，实时数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、datetimes、ms、values、states、qualities、errors 的长度与 count 一致。
*/
RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbp_get_perf_values(
  rtdb_int32 handle,
  rtdb_int32* count,
  int* perf_ids,
  int* datetimes,
  rtdb_time_type* ms,
  rtdb_float64* values,
  rtdb_int64* states,
  rtdb_int16* qualities,
  rtdb_error* errors);

RTDBAPI
rtdb_error
RTDBAPI_CALLRULE
rtdbp_get_perf_values64(
    rtdb_int32 handle,
    rtdb_int32* count,
    int* perf_ids,
    rtdb_timestamp_type* datetimes,
    rtdb_subtime_type* subtimes,
    rtdb_float64* values,
    rtdb_int64* states,
    rtdb_int16* qualities,
    rtdb_error* errors);


////////////////////////////////////////// deprecated ////////////////////////////////////////////
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdb_subscribe_connect_ex)
rtdb_error
RTDBAPI_CALLRULE
rtdb_subscribe_connect(
    rtdb_int32 handle,
    rtdb_connect_event callback
);

/*
* 命名：rtdb_write_named_type_field_by_name
* 功能：按名称填充自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_name]   自定义类型中需要填充的字段的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数
*      [field_type]   field_name字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入/输出参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      需要填充的字段数值的缓冲区,输入参数
*     [field_len]    自定义类型中字段数值的缓冲区中数据的长度,输入参数
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdb_write_named_type_field_by_name32)
rtdb_error
RTDBAPI_CALLRULE
rtdb_write_named_type_field_by_name(
    rtdb_int32 handle,
    const char* type_name,
    const char* field_name,
    rtdb_int32 field_type,
    void* object,
    rtdb_int16 object_len,
    const void* field,
    rtdb_int16 field_len
);

/*
* 命名：rtdb_write_named_type_field_by_pos
* 功能：按位置填充自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_pos]    自定义类型中需要填充的字段的位置，指字段在所有字段中的位置，从0开始，输入参数
*      [field_type]   field_pos位置所在字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入/输出参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      需要填充的字段数值的缓冲区,输入参数
*     [field_len]    自定义类型中字段数值的缓冲区中数据的长度,输入参数
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdb_write_named_type_field_by_pos32)
rtdb_error
RTDBAPI_CALLRULE
rtdb_write_named_type_field_by_pos(
    rtdb_int32 handle,
    const char* type_name,
    rtdb_int32 field_pos,
    rtdb_int32 field_type,
    void* object,
    rtdb_int16 object_len,
    const void* field,
    rtdb_int16 field_len
);

/*
* 命名：rtdb_read_named_type_field_by_name
* 功能：按名称提取自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_name]   自定义类型中需要提取的字段的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数
*      [field_type]   field_name字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      被读取的字段的数值的缓冲区,输入/输出参数
*     [field_len]    field字段数值缓冲区的长度,输入参数
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdb_read_named_type_field_by_name32)
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_named_type_field_by_name(
    rtdb_int32 handle,
    const char* type_name,
    const char* field_name,
    rtdb_int32 field_type,
    const void* object,
    rtdb_int16 object_len,
    void* field,
    rtdb_int16 field_len
);

/*
* 命名：rtdb_read_named_type_field_by_pos
* 功能：按位置提取自定义类型数值中字段的内容
* 参数：
*      [handle]       连接句柄
*      [type_name]    自定义类型的名称，名称长度不能超过RTDB_TYPE_NAME_SIZE的长度，输入参数，
*      [field_pos]    自定义类型中需要提取的字段的位置，指字段在所有字段中的位置，从0开始，输入参数
*      [field_type]   field_pos位置所在字段的类型，RTDB_TYPE所支持的基础类型，输入参数
*      [object]       自定义类型数值的缓冲区,输入参数
*     [object_len]   object缓冲区的长度,输入参数
*     [field]      被读取的字段的数值的缓冲区,输入/输出参数
*     [field_len]    field字段数值缓冲区的长度,输入参数
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdb_read_named_type_field_by_pos32)
rtdb_error
RTDBAPI_CALLRULE
rtdb_read_named_type_field_by_pos(
    rtdb_int32 handle,
    const char* type_name,
    rtdb_int32 field_pos,
    rtdb_int32 field_type,
    const void* object,
    rtdb_int16 object_len,
    void* field,
    rtdb_int16 field_len
);

/**
*
* \brief 批量标签点属性更改通知订阅
*
* \param handle    连接句柄
* \param callback  rtdbb_tags_change_event 类型回调接口，输入，
*                    被订阅的标签点的属性值发生变化或被删除时将调用该接口通知用户,
*                    参见 rtdb.h 中原型的定义
* \remark 用于订阅测点的连接句柄必需是独立的，不能再用来调用其它 api，
*        否则返回 RtE_OTHER_SDK_DOING 错误。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbb_subscribe_tags_ex)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_subscribe_tags(
    rtdb_int32 handle,
    rtdbb_tags_change_event callback
);

/**
*
* \brief 读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blob      字节型数组，输出，实时二进制/字符串数值
* \param len       短整型，输出，二进制/字符串数值长度
* \param quality   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_get_blob_snapshot64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshot(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* datetime,
    rtdb_time_type* ms,
    rtdb_byte* blob,
    rtdb_int16* len,
    rtdb_int16* quality
);

/**
*
* \brief 批量读取二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
* \param blobs     字节型指针数组，输出，实时二进制/字符串数值
* \param lens      短整型数组，输入/输出，二进制/字符串数值长度，
*                    输入时表示对应的 blobs 指针指向的缓冲区长度，
*                    输出时表示实际得到的 blob 长度，如果 blob 的长度大于缓冲区长度，会被截断。
* \param qualities 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_get_blob_snapshots64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_blob_snapshots(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_int32* datetimes,
    rtdb_time_type* ms,
    rtdb_byte* const* blobs,
    rtdb_int16* lens,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，实时二进制/字符串数值
* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_put_blob_snapshot64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshot(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 datetime,
    rtdb_time_type ms,
    const rtdb_byte* blob,
    rtdb_int16 len,
    rtdb_int16 quality
);

/**
*
* \brief 批量写入二进制/字符串实时数据
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、blobs、lens、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
* \param ids       整型数组，输入，标签点标识
* \param datetimes 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
* \param blobs     字节型指针数组，输入，实时二进制/字符串数值
* \param lens      短整型数组，输入，二进制/字符串数值长度，
*                    表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_put_blob_snapshots64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_blob_snapshots(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_int32* datetimes,
    const rtdb_time_type* ms,
    const rtdb_byte* const* blobs,
    const rtdb_int16* lens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 批量标签点快照改变的通知订阅
*
* \param handle    连接句柄
* \param count     整型，输入/输出，标签点个数，输入时表示 ids、errors 的长度，
*                    输出时表示成功订阅的标签点个数，不得超过 RTDB_MAX_SUBSCRIBE_SNAPSHOTS。
* \param ids       整型数组，输入，标签点标识列表。
* \param callback  rtdbs_snaps_event 类型回调接口，输入，
*                    被订阅的标签点的快照值发生变化时将调用该接口，
*                    参见 rtdb.h 中原型的定义。
* \param errors    无符号整型数组，输出，
*                    写入实时数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、errors 的长度与 count 一致。
*        用于订阅快照的连接句柄必需是独立的，不能再用来调用其它 api，
*        否则返回 RtE_OTHER_SDK_DOING 错误。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_subscribe_snapshots_ex64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_subscribe_snapshots(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdbs_snaps_event callback,
    rtdb_error* errors
);

/**
* 命名：rtdbs_get_named_type_snapshot
* 功能：获取自定义类型测点的单个快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [object]    字节型数组，输出，实时自定义类型标签点的数值
*        [length]    短整型，输入/输出，自定义类型标签点的数值长度
*        [quality]   短整型，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_get_named_type_snapshot64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshot(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* datetime,
    rtdb_time_type* ms,
    void* object,
    rtdb_int16* length,
    rtdb_int16* quality
);

/**
* 命名：rtdbs_get_named_type_snapshots
* 功能：批量获取自定义类型测点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功获取实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输出，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输出，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，返回相应的纳秒值；否则为 0
*        [objects]   指针数组，输出，自定义类型标签点数值
*        [lengths]   短整型数组，输入/输出，自定义类型标签点数值长度，
*                    输入时表示对应的 objects 指针指向的缓冲区长度，
*                    输出时表示实际得到的 objects 长度，如果 objects 的长度大于缓冲区长度，会被截断。
*        [qualities] 短整型数组，输出，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_get_named_type_snapshots64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_get_named_type_snapshots(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    rtdb_int32* datetimes,
    rtdb_time_type* ms,
    void* const* objects,
    rtdb_int16* lengths,
    rtdb_int16* qualities,
    rtdb_error* errors
);

/**
* 命名：rtdbs_put_named_type_snapshot
* 功能：写入单个自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
*        [object]    void类型数组，输入，自定义类型标签点数值
*        [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
*        [quality]   短整型，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_put_named_type_snapshot64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshot(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 datetime,
    rtdb_time_type ms,
    const void* object,
    rtdb_int16 length,
    rtdb_int16 quality
);

/**
* 命名：rtdbs_put_named_type_snapshots
* 功能：批量写入自定义类型标签点的快照
* 参数：
*        [handle]    连接句柄
*        [count]     整型，输入/输出，标签点个数，
*                    输入时表示 ids、datetimes、ms、objects、lengths、qualities、errors 的长度，
*                    输出时表示成功写入实时值的标签点个数
*        [ids]       整型数组，输入，标签点标识
*        [datetimes] 整型数组，输入，实时数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型数组，输入，实时数值时间列表，
*                    对于时间精度为纳秒的标签点，表示相应的纳秒值；否则忽略
*        [objects]   void类型指针数组，输入，自定义类型标签点数值
*        [lengths]   短整型数组，输入，自定义类型标签点数值长度，
*                    表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
*        [qualities] 短整型数组，输入，实时数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]    无符号整型数组，输出，读取实时数据的返回值列表，参考rtdb_error.h
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbs_put_named_type_snapshots64)
rtdb_error
RTDBAPI_CALLRULE
rtdbs_put_named_type_snapshots(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_int32* datetimes,
    const rtdb_time_type* ms,
    const void* const* objects,
    const rtdb_int16* lengths,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
*
* \brief 获取首个存档文件的路径、名称、状态和最早允许写入时间。
*
* \param handle          连接句柄
* \param path            字符数组，输出，首个存档文件的目录路径，长度至少为 RTDB_PATH_SIZE。
* \param file            字符数组，输出，首个存档文件的名称，长度至少为 RTDB_FILE_NAME_SIZE。
* \param state           整型，输出，取值 RTDB_INVALID_ARCHIVE、RTDB_ACTIVED_ARCHIVE、
*                          RTDB_NORMAL_ARCHIVE、RTDB_READONLY_ARCHIVE 之一，表示文件状态
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdba_get_archives)
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_first_archive(
    rtdb_int32 handle,
    char* path,
    char* file,
    rtdb_int32* state
);

/**
*
* \brief 获取下一个存档文件的路径、名称、状态和最早允许写入时间。
*
* \param handle         连接句柄
* \param path           字符数组，输入/输出，
*                         输入由调用 rtdba_get_first_archive 或
*                         上次调用 rtdba_get_next_archive 返回的文件目录路径，
*                         输出下一个存档文件的目录路径，长度至少为 RTDB_PATH_SIZE。
* \param file           字符数组，输入/输出，
*                         输入由调用 rtdba_get_first_archive 或
*                         上次调用 rtdba_get_next_archive 返回的文件名，
*                         输出下一个存档文件的名称，长度至少为 RTDB_FILE_NAME_SIZE。
* \param state          整型，输出，取值 RTDB_INVALID_ARCHIVE、RTDB_ACTIVED_ARCHIVE、
*                         RTDB_NORMAL_ARCHIVE、RTDB_READONLY_ARCHIVE 之一，表示文件状态
* \remark 当 path 返回内容为 "END" 时表示全部存档文件已经遍历完毕。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdba_get_archives)
rtdb_error
RTDBAPI_CALLRULE
rtdba_get_next_archive(
    rtdb_int32 handle,
    char* path,
    char* file,
    rtdb_int32* state
);

/**
*
* \brief 读取单个标签点某个时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
* \param mode          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param datetime      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
* \param blob          字节型数组，输出，二进制/字符串历史值
* \param len           短整型，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
* \param quality       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_get_single_blob_value64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_blob_value(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_int32* datetime,
    rtdb_time_type* ms,
    rtdb_byte* blob,
    rtdb_int16* len,
    rtdb_int16* quality
);

/**
*
* \brief 读取单个标签点一段时间的二进制/字符串型历史数据
*
* \param handle        连接句柄
* \param id            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
* \param count         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
* \param datetime1     整型，输入，表示开始时间秒数；
* \param ms1           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetime2     整型，输入,表示结束时间秒数；
* \param ms2           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
* \param datetimes     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
* \param ms            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
* \param lens          短整型数组，输入/输出，输入时表示 blob 的长度，
*                        输出时表示实际获取的二进制/字符串数据长度。
*                        当blobs为空指针时，表示只获取每条数据的长度，此时会忽略输入的lens
* \param blobs         字节型数组，输出，二进制/字符串历史值。可以设置为空指针，表示只获取每条数据的长度
* \param qualities     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_get_archived_blob_values64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_blob_values(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32* count,
    rtdb_int32 datetime1,
    rtdb_time_type ms1,
    rtdb_int32 datetime2,
    rtdb_time_type ms2,
    rtdb_int32* datetimes,
    rtdb_time_type* ms,
    rtdb_int16* lens,
    rtdb_byte* const* blobs,
    rtdb_int16* qualities
);

/**
*
* \brief 写入单个二进制/字符串标签点在某一时间的历史数据
*
* \param handle    连接句柄
* \param id        整型，输入，标签点标识
* \param datetime  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
* \param ms        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
* \param blob      字节型数组，输入，历史二进制/字符串数值
* \param len       短整型，输入，二进制/字符串数值长度，超过一个页大小数据将被截断。
* \param quality   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
* \remark 本接口只对数据类型为 RTDB_BLOB、RTDB_STRING 的标签点有效。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_put_single_blob_value64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_blob_value(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 datetime,
    rtdb_time_type ms,
    const rtdb_byte* blob,
    rtdb_int16 len,
    rtdb_int16 quality
);

/**
*
* \brief 写入批量标签点批量字符串型历史存储数据
*
* \param handle        连接句柄
* \param count         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
* \param ids           整型数组，输入，标签点标识
* \param datetimes     整型数组，输入，表示对应的历史数值时间秒数。
* \param ms            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
* \param blobs         字节型指针数组，输入，实时二进制/字符串数值
* \param lens          短整型数组，输入，二进制/字符串数值长度，
*                        表示对应的 blobs 指针指向的缓冲区长度，超过一个页大小数据将被截断。
* \param qualities     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
* \param errors        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* \remark 用户须保证 ids、datetimes、ms、lens、blobs、qualities、errors 的长度与 count 一致，
*        本接口仅对数据类型为 RTDB_STRING、RTDB_BLOB 的标签点有效。
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_put_archived_blob_values64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_blob_values(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_int32* datetimes,
    const rtdb_time_type* ms,
    const rtdb_byte* const* blobs,
    const rtdb_int16* lens,
    const rtdb_int16* qualities,
    rtdb_error* errors
);

/**
* 命名：rtdbh_get_single_named_type_value
* 功能：读取单个自定义类型标签点某个时间的历史数据
* 参数：
*        [handle]        连接句柄
*        [id]            整型，输入，标签点标识
*        [mode]          整型，输入，取值 RTDB_NEXT、RTDB_PREVIOUS、RTDB_EXACT 之一：
*                        RTDB_NEXT 寻找下一个最近的数据；
*                        RTDB_PREVIOUS 寻找上一个最近的数据；
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*        [datetime]      整型，输入/输出，输入时表示时间秒数；
*                        输出时表示实际取得的历史数值对应的时间秒数。
*        [ms]            短整型，输入/输出，如果 id 指定的标签点时间精度为纳秒，
*                        则输入时表示时间纳秒数；输出时表示实际取得的历史数值时间纳秒数。
*                        否则忽略输入，输出时为 0。
*        [object]        void数组，输出，自定义类型标签点历史值
*        [length]        短整型，输入/输出，输入时表示 object 的长度，
*                        输出时表示实际获取的自定义类型标签点数据长度。
*        [quality]       短整型，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_get_single_named_type_value64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_single_named_type_value(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 mode,
    rtdb_int32* datetime,
    rtdb_time_type* ms,
    void* object,
    rtdb_int16* length,
    rtdb_int16* quality
);

/**
* 命名：rtdbh_get_archived_named_type_values
* 功能：连续读取自定义类型标签点的历史数据
* 参数：
*        [handle]        连接句柄
*        [id]            整型，输入，标签点标识
*                        RTDB_EXACT 取指定时间的数据，如果没有则返回错误 RtE_DATA_NOT_FOUND；
*        [datetime1]     整型，输入，表示开始时间秒数；
*        [ms1]           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
*        [datetime2]     整型，输入,表示结束时间秒数；
*        [ms2]           短整型，输入，指定的标签点时间精度为纳秒，
*                        表示时间纳秒数；
*        [length]        短整型数组，输入，输入时表示 objects 的长度，
*        [count]         整型，输入/输出，输入表示想要查询多少数据
*                        输出表示实际查到多少数据
*        [datetimes]     整型数组，输出，表示实际取得的历史数值对应的时间秒数。
*        [ms]            短整型，输出，如果 id 指定的标签点时间精度为纳秒，
*                        表示实际取得的历史数值时间纳秒数。
*        [objects]       void类型数组，输出，自定义类型标签点历史值
*        [qualities]     短整型数组，输出，历史值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_get_archived_named_type_values64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_get_archived_named_type_values(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 datetime1,
    rtdb_time_type ms1,
    rtdb_int32 datetime2,
    rtdb_time_type ms2,
    rtdb_int16 length,
    rtdb_int32* count,
    rtdb_int32* datetimes,
    rtdb_time_type* ms,
    void* const* objects,
    rtdb_int16* qualities
);

/**
* 命名：rtdbh_put_single_named_type_value
* 功能：写入自定义类型标签点的单个历史事件
* 参数：
*        [handle]    连接句柄
*        [id]        整型，输入，标签点标识
*        [datetime]  整型，输入，数值时间列表,
*                    表示距离1970年1月1日08:00:00的秒数
*        [ms]        短整型，输入，历史数值时间，
*                    对于时间精度为纳秒的标签点，存放相应的纳秒值；否则忽略
*        [object]    void数组，输入，历史自定义类型标签点数值
*        [length]    短整型，输入，自定义类型标签点数值长度，超过一个页大小数据将被截断。
*        [quality]   短整型，输入，历史数值品质，数据库预定义的品质参见枚举 RTDB_QUALITY
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_put_single_named_type_value64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_single_named_type_value(
    rtdb_int32 handle,
    rtdb_int32 id,
    rtdb_int32 datetime,
    rtdb_time_type ms,
    const void* object,
    rtdb_int16 length,
    rtdb_int16 quality
);

/**
* 命名：rtdbh_put_archived_named_type_values
* 功能：批量补写自定义类型标签点的历史事件
* 参数：
*        [handle]        连接句柄
*        [count]         整型，输入/输出，
*                        输入时表示 ids、datetimes、ms、lens、blobs、qualities、errors 的长度，
*                        即历史值个数；输出时返回实际写入的数值个数
*        [ids]           整型数组，输入，标签点标识
*        [datetimes]     整型数组，输入，表示对应的历史数值时间秒数。
*        [ms]            短整型数组，输入，如果 id 指定的标签点时间精度为纳秒，
*                        表示对应的历史数值时间纳秒；否则忽略。
*        [objects]       void类型指针数组，输入，自定义类型标签点数值
*        [lengths]       短整型数组，输入，自定义类型标签点数值长度，
*                        表示对应的 objects 指针指向的缓冲区长度，超过一个页大小数据将被截断。
*        [qualities]     短整型数组，输入，历史数值品质列表，数据库预定义的品质参见枚举 RTDB_QUALITY
*        [errors]        无符号整型数组，输出，写入历史数据的返回值列表，参考rtdb_error.h
* 备注：用户须保证 ids、datetimes、ms、lens、objects、qualities、errors 的长度与 count 一致，
*        如果 datetimes、ms 标识的数据已经存在，其值将被替换。
*/
RTDBAPI
RTDBAPI_DEPRECATED_FOR(deprecated, rtdbh_put_archived_named_type_values64)
rtdb_error
RTDBAPI_CALLRULE
rtdbh_put_archived_named_type_values(
    rtdb_int32 handle,
    rtdb_int32* count,
    const rtdb_int32* ids,
    const rtdb_int32* datetimes,
    const rtdb_time_type* ms,
    const void* const* objects,
    const rtdb_int16* lengths,
    const rtdb_int16* qualities,
    rtdb_error* errors
);


////////////////////////////////////////// deprecated and not support ////////////////////////////////////////////
/**
* 命名：rtdbb_query_load_memory
* 功能：查询标签点属性加载到内存中的当前配置
* 参数：
*        [handle]   连接句柄
*        [load_memory_flag]     输出，加载到内存中的当前配置
*/
RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_query_load_memory(rtdb_int32 handle,
    rtdb_int32* load_memory_flag);

/**
*
* \brief 设置标签点属性加载到内存中的配置
*
* \param handle   连接句柄
* \param load_memory_flag     输入/输出，输入时，表示要设置的加载到内存中的配置
*                                          输出时，表示成功加载到内存中的配置
* \remark 如果某些属性列成功加载到内存，某些属性列分配内存失败，就会返回RtE_OUT_OF_MEMORY，
*       此时还要根据暑促的load_memory_flag来判断是不是全部分配失败。
*/
RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_set_load_memory(
    rtdb_int32 handle,
    rtdb_int32* load_memory_flag);

/**
*
* \brief 查询根据配置加载到内存中的标签点属性所占用的内存
*
* \param handle   连接句柄
* \param load_memory_flag     输入，计算需要多少内存的配置
* \param need_memory_size     输出，需要占用内存的字节数，单位是byte
*/
RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_query_need_memory_size(rtdb_int32 handle, rtdb_int32* load_memory_flag, rtdb_int64* need_memory_size);

/**
*
* \brief 新建历史存档文件并追加到历史数据库
*
* \param handle     连接句柄
* \param path       字符串，输入，文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，文件名，文件后缀名应为.rdf。
* \param mb_size    整型，输入，文件兆字节大小，单位为 MB。
*/
RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdba_create_archive(
    rtdb_int32 handle,
    const char* path,
    const char* file,
    rtdb_int32 mb_size
);

/**
*
* \brief 激活指定存档文件为活动存档文件
*
* \param handle     连接句柄
* \param path       字符串，输入，存档文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，存档文件名。
*/
RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdba_reactive_archive(
    rtdb_int32 handle,
    const char* path,
    const char* file
);

/**
*
* \brief 合并附属文件到所属主文件
*
* \param handle     连接句柄
* \param path       字符串，输入，主文件所在目录路径，必须以"\"或"/"结尾。
* \param file       字符串，输入，主文件名。
*/
RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdba_merge_archive(
    rtdb_int32 handle,
    const char* path,
    const char* file
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_tags_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* tags,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_ms_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* ms,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_compdevs_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* compdevs,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_compmaxs_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* compmaxs,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_compmins_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* compmins,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_excdevs_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* excdevs,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_excmaxs_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* excmaxs,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_excmins_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* excmins,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_classofs_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_uint8* classofs,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_tables_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_uint16* tables,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_summarys_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* summarys,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_mirrors_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* mirrors,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_compress_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* compress,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_steps_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* steps,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_shutdowns_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* shutdowns,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_archives_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* archives,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_table_dot_tags_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* table_dot_tags,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_descs_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* descs,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_units_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* units,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_changers_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* changers,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_creators_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* creators,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_lowlimits_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* lowlimits,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_highlimits_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* highlimits,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_typicals_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* typicals,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_changedates_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* changedates,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_createdates_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* createdates,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_digits_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int16* digits,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_compdevpercents_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* compdevpercents,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_excdevpercents_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* excdevpercents,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_sources_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* sources,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_scans_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* scans,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_instruments_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* instruments,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_location1s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* location1s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_location2s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* location2s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_location3s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* location3s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_location4s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* location4s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_location5s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* location5s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_userint1s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* userint1s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_userint2s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* userint2s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_userreal1s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* userreal1s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_userreal2s_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_float32* userreal2s,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_equations_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    char* const* equations,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_triggers_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* triggers,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_timecopys_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_byte* timecopys,
    rtdb_error* errors
);

RTDBAPI
RTDBAPI_DEPRECATED(removed)
rtdb_error
RTDBAPI_CALLRULE
rtdbb_get_periods_property(
    rtdb_int32 handle,
    rtdb_int32* count,
    rtdb_int32* ids,
    rtdb_int32* periods,
    rtdb_error* errors
);

#endif  // __RTDB_API_H__
