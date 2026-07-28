# 可用 API 列表

## Easy API 列表

1. Login  登录数据库
1. SetLocation  设置客户端时区，默认为亚洲东八
1. Logout  登出数据库
1. GetClientVersion  获取客户端版本
1. GetHandleInfo  获取连接句柄所连接的服务器相关信息
1. SetClientOption  设置客户端参数
1. GetServerOption  获取服务端参数
1. SetServerOption  设置服务端参数
1. GetSocketInfos  获取服务端SocketInfo列表，单机服务端返回一个SocketInfo列表，双活服务端返回两个SocketInfo列表
1. GetOwnSocketInfo  获取当前连接的SocketInfo，单机服务端返回一个SocketInfo，双活服务端返回两个SocketInfo
1. SetSocketTimeout  设置Socket超时时间
1. KillSocket  断开Socket
1. AddIpBlackList  添加IP黑名单项
1. UpdateIpBlackList  更新连接黑名单项
1. DeleteIpBlackList  删除连接黑名单项
1. GetIpBlackLists  获得连接黑名单列表
1. AddIpWhiteList  添加连接白名单
1. UpdateIpWhiteList  更新连接白名单
1. DeleteIpWhiteList  删除白名单
1. GetIpWhiteLists  获取连接白名单列表
1. UpdatePassword  修改用户密码
1. UpdateOwnPassword  修改自己的密码
1. GetPriv  获取连接权限
1. SetPriv  设置连接权限
1. AddUser  添加用户
1. DeleteUser  删除用户
1. LockUser  锁定用户
1. GetUsers  获取用户列表
1. AddNamedType  创建自定义类型
1. DeleteNamedType  删除自定义类型
1. GetNamedType  获取自定义类型
1. GetNamedTypes  获取自定义类型列表
1. UpdateNamedType  修改自定义类型
1. WriteNamedTypeFieldByName  备注！：这五个函数相对复杂，用处不大，在Easy端保留，暂不对外开放
1. WriteNamedTypeFieldByPos  按位置填充自定义类型数值中字段的内容
1. ReadNamedTypeFieldByName  按名称提取自定义类型数值中字段的内容
1. ReadNamedTypeFieldByPos  按位置提取自定义类型数值中字段的内容
1. NamedTypeNameFieldCheck  检查自定义类型名称及字段命名是否符合规则
1. ServerHostTime  服务端主机时间
1. DurationToString  时间段转字符串, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
1. StringToDuration  字符串转时间段, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
1. StringToTime  字符串转时间戳
1. GetQualityDesc  获取质量码说明
1. GetDriveLetterList  获取盘符列表, windows平台是C、D、E、F这些盘符，linux平台是 / 盘符
1. GetDirItemList  获取目录项列表
1. CreateDir  创建目录
1. ReadFile  读取文件
1. CreateTable  创建表
1. DeleteTable  删除表
1. GetTable  获取表
1. GetTables  获取表列表
1. UpdateTableName  更新表名
1. UpdateTableDesc  更新表描述
1. AddPoint  创建点
1. AddPoints  批量创建点
1. DeletePoint  删除点
1. UpdatePoint  更新点
1. GetPoints  批量获取标签点
1. GetPoint  获取点
1. GetTypesProperty  批量获取标签点的数据类型
1. FindPoints  根据 表名.点名 搜索标签点
1. MovePoint  移动点到指定表
1. SearchPoint  分页搜索点
1. ClearRecycler  清空回收站
1. GetRecycledPoints  分段获取回收站中的点
1. RecoverPoint  从回收站中恢复点到某个表
1. PurgePoint  从回收站中清除点
1. SearchRecycledPoint  从回收站中搜索点
1. GetPointCountFromValueType  获取某个数据类型的点个数 (可以是内置类型，也可以是自定义类型)
1. WriteValue  写入值
1. WriteValues  批量写入值
1. WriteSection  写断面(批量写入多个Point，每个Point写入一个TVQ)
1. ReadLasts  批量读取实时快照值(当前标签点最后一个写入的TVQ)
1. ReadLast  读取实时快照值(当前标签点最后一个写入的TVQ)
1. hisModeForSingleRead  将组合查询模式拆分为底层历史读取接口支持的基础模式与未命中时的回退模式。
1. ReadValue  读取单个TVQ
1. ReadRange  读取某个时间段内的TVQ
1. ReadPlot  读取用于绘图的TVQ
1. ReadTimed  获取差值，每个差值都需要指定一个确定的时间戳
1. ReadInterpo  获取差值, 会自动将start、end等分成count个时间戳，然后取这些时间戳的差值
1. ReadInterval  读取从start开始的等间隔差值
1. ReadSection  读取断面
1. ReadSummary  获取统计值 (从start到end的统计值)
1. ReadBatchesSummary  获取等间隔统计值 (在start和end之间，按照interval作为时间间隔计算每个间隔的统计值)
1. ReadBatches  开始以分段返回方式读取一段时间内的储存数据
1. ReadNext  分段读取一段时间内的储存数据
1. RemoveValue  删除点值
1. RemoveRangeValues  批量删除点值(从start到end的范围)
1. UpdateValue  更新点值
1. FlushArchivedValues  刷新历史缓存(把标签点的缓存进行手动归档)
1. QueryBigJob  查询正在进行的大任务
1. CancelBigJob  取消正在执行的大任务
1. JobMessage  获取任务描述
1. ComputeHistory  重算｜补算 补历史值或者修改历史值之后，对应的计算点可以按需进行重算|补算，保证计算点的数值正确
1. GetPointEquation  获取标签点方程式
1. GetEquationGraph  获取标签点对应的方程式关联关系图(方程式本身是一个有向无环图，数据在这个图内的公式之间流转)
1. GetPerfPointInfo  获取性能监控点的说明信息
1. GetPerfPointValue  获取性能监控点的实时值
1. CreateRangedArchive  创建存档文件(按照时间范围创建)
1. AppendArchive  存档文件入列(将存档文件插入到存档队列中, 队列中的存档文件按照start、end依次排列，注意时间段不能有交叠的部分)
1. RemoveArchive  存档文件解列
1. ShiftActived  切换活动文件
1. GetArchives  获取存档文件的基本信息(返回全部存档文件的基本信息)
1. UpdateArchive  更新存档文件选项
1. ArrangeArchive  整理存档文件(会重新调整数据块分布，使存档文件更加紧凑)
1. ReindexArchive  重建存档文件索引(用于进行数据恢复)
1. ConvertIndex  存档文件索引转换格式(老版本索引格式转换为新版本，新版本索引更快)
1. BackupArchive  备份存档
1. MoveArchive  移动存档
1. CreateDatagram  创建数据流订阅
1. RemoveDatagram  取消数据流订阅
1. RecvDatagram  从订阅数据流中获取数据
1. SubscribeTags  订阅标签点属性更新
1. CancelSubscribeTags  取消订阅标签点属性更新
1. SubscribeSnapshots  订阅快照，只要快照发生变化，就会触发订阅
1. SubscribeDeltaSnapshots  订阅Delta快照，快照变化需要超过Delta，才会触发订阅，这样可以节约流量
1. ChangeSubscribeSnapshots  修改快照订阅设置，新增或删除标签点
1. CancelSubscribeSnapshots  取消快照订阅
1. SubscribeConnectEvents  订阅API调用连接事件
1. CancelSubscribeConnectEvents  取消订阅API调用连接事件

## Raw API 列表

1. RawRtdbGetApiVersionWarp  返回 ApiVersion 版本号
1. RawRtdbSetOptionWarp  配置 API库 的行为参数，详见 RtdbApiOption 枚举
1. RawRtdbCreateDatagramHandleWarp  创建数据流
1. RawRtdbRemoveDatagramHandleWarp  删除数据流
1. RawRtdbRecvDatagramWarp  接收数据流
1. RawRtdbConnectWarp  建立同 RTDB 数据库的网络连接, 注意这里只是创建连接，并没有进行用户登陆
1. RawRtdbLoginWarp  以有效帐户登录
1. RawRtdbDisconnectWarp  断开同 RTDB 数据平台的连接
1. RawRtdbSubscribeConnectExWarp  创建API调用订阅连接
1. RawRtdbCancelSubscribeConnectWarp  关闭API调用订阅链接
1. RawRtdbGetDbInfo1Warp  获得字符串型数据库系统参数
1. RawRtdbGetDbInfo2Warp  获得整型数据库系统参数
1. RawRtdbSetDbInfo1Warp  设置字符串型数据库系统参数
1. RawRtdbSetDbInfo2Warp  设置整型数据库系统参数
1. RawRtdbConnectionCountWarp  获取 RTDB 服务器当前连接个数
1. RawRtdbGetConnectionsWarp  列出 RTDB 服务器的所有socket连接句柄, 注意这里指的是socket连接，区分于ConnectHandle
1. RawRtdbGetOwnConnectionWarp  获取当前连接的socket句柄
1. RawRtdbGetConnectionInfoWarp  获取 RTDB 服务器指定连接的信息
1. RawRtdbGetConnectionInfoIpv6Warp  获取 RTDB 服务器指定连接的ipv6版本
1. RawRtdbOsType  获取连接句柄所连接的服务器操作系统类型
1. RawRtdbGetHandleInfoWarp  获取连接句柄所连接的服务器相关信息
1. RawRtdbChangePasswordWarp  修改用户帐户口令
1. RawRtdbChangeMyPasswordWarp  用户修改自己帐户口令
1. RawRtdbGetPrivWarp  获取连接权限
1. RawRtdbChangePrivWarp  修改用户帐户权限, 只有管理员有修改权限
1. RawRtdbAddUserWarp  添加用户帐户
1. RawRtdbRemoveUserWarp  删除用户帐户
1. RawRtdbLockUserWarp  启用或禁用用户, 只有管理员有启用禁用权限
1. RawRtdbGetUsersWarp  获得所有用户
1. RawRtdbAddBlacklistWarp  添加连接黑名单项
1. RawRtdbUpdateBlacklistWarp  更新连接连接黑名单项
1. RawRtdbRemoveBlacklistWarp  删除连接黑名单项
1. RawRtdbGetBlacklistWarp  获得连接黑名单
1. RawRtdbAddAuthorizationWarp  添加信任连接段
1. RawRtdbUpdateAuthorizationWarp  更新信任连接段
1. RawRtdbRemoveAuthorizationWarp  删除信任连接段
1. RawRtdbGetAuthorizationsWarp  获得所有信任连接段
1. RawRtdbHostTime64Warp  RawRtdbHostTimeWarp 获取 RTDB 服务器当前UTC时间
1. RawRtdbFormatTimespanWarp  根据时间跨度值生成时间格式字符串, 如：输入10， 输出10s, 输入60，输出1n
1. RawRtdbParseTimespanWarp  根据时间格式字符串解析时间跨度值, 如：输入2n，输出120，表示2分钟
1. RawRtdbParseTimeWarp  根据时间格式字符串解析时间值
1. RawRtdbFormatMessageWarp  获取 Rtdb API 调用返回值的简短描述(错误码对应的Desc)
1. RawRtdbJobMessageWarp  获取任务的简短描述
1. RawRtdbSetTimeoutWarp  设置连接超时时间
1. RawRtdbGetTimeoutWarp  获得连接超时时间
1. RawRtdbKillConnectionWarp  断开已知连接
1. RawRtdbGetLogicalDriversWarp  获得逻辑盘符
1. RawRtdbOpenPathWarp  打开目录以便遍历其中的文件和子目录。
1. RawRtdbReadPath64Warp  RawRtdbReadPathWarp 读取目录中的文件或子目录
1. RawRtdbClosePathWarp  关闭当前遍历的目录
1. RawRtdbMkdirWarp  建立目录
1. RawRtdbGetFileSizeWarp  获得指定服务器端文件的大小
1. RawRtdbReadFileWarp  读取服务器端指定文件的内容
1. RawRtdbGetMaxBlobLenWarp  取得数据库允许的blob与str类型测点的最大长度
1. RawRtdbFormatQualityWarp  取得质量码对应的定义
1. RawRtdbJudgeConnectStatusWarp  判断连接是否可用
1. RawRtdbFormatIpaddrWarp  将整形IP转换为字符串形式的IP
1. RawRtdbbAppendTableWarp  添加新表
1. RawRtdbbRemoveTableByIdWarp  根据表 id 删除表及表中标签点
1. RawRtdbbRemoveTableByNameWarp  根据表名删除表及表中标签点
1. RawRtdbbTablesCountWarp  取得标签点表总数
1. RawRtdbbGetTablesWarp  取得所有标签点表的ID
1. RawRtdbbGetTableSizeByIdWarp  根据表 id 获取表中包含的标签点数量(大概数量, 包含被标记删除的点)
1. RawRtdbbGetTableSizeByNameWarp  根据表名称获取表中包含的标签点数量(大概数量, 包含被标记删除的点)
1. RawRtdbbGetTableRealSizeByIdWarp  根据表 id 获取表中实际包含的标签点数量(实际数量, 不含被删除的点)
1. RawRtdbbGetTablePropertyByIdWarp  根据标签点表 id 获取表属性
1. RawRtdbbGetTablePropertyByNameWarp  根据表名获取标签点表属性
1. RawRtdbbInsertPointWarp  使用完整的属性集来创建单个标签点
1. RawRtdbbInsertMaxPointWarp  使用最大长度的完整属性集来创建单个标签点
1. RawRtdbbRemovePointByIdWarp  根据 id 删除单个标签点
1. RawRtdbbRemovePointByNameWarp  根据标签点全名删除单个标签点
1. RawRtdbbInsertMaxPointsWarp  使用最大长度的完整属性集来批量创建标签点
1. RawRtdbbInsertBasePointWarp  使用最小的属性集来创建单个标签点
1. RawRtdbbInsertNamedTypePointWarp  使用完整的属性集来创建单个自定义数据类型标签点
1. RawRtdbbMovePointByIdWarp  根据 id 移动单个标签点到其他表
1. RawRtdbbGetPointsPropertyWarp  批量获取标签点属性
1. RawRtdbbGetMaxPointsPropertyWarp  按最大长度批量获取标签点属性
1. RawRtdbbGetTypesPropertyWarp  批量获取标签点的数据类型
1. RawRtdbbSearchWarp  搜索符合条件的标签点，使用标签点名时支持通配符
1. RawRtdbbSearchInBatchesWarp  分批继续搜索符合条件的标签点，使用标签点名时支持通配符
1. RawRtdbbSearchExWarp  搜索符合条件的标签点，使用标签点名时支持通配符
1. RawRtdbbSearchPointsCountWarp  搜索符合条件的标签点，获取标签点数，使用标签点名时支持通配符
1. RawRtdbbUpdatePointPropertyWarp  更新单个标签点属性
1. RawRtdbbUpdateMaxPointPropertyWarp  按最大长度更新单个标签点属性
1. RawRtdbbFindPointsWarp  根据 "表名.标签点名" 格式批量获取标签点标识
1. RawRtdbbFindPointsExWarp  根据 "表名.标签点名" 格式批量获取标签点标识
1. RawRtdbbSortPointsWarp  根据标签属性字段对标签点标识进行排序
1. RawRtdbbUpdateTableNameWarp  根据表 ID 更新表名称。
1. RawRtdbbUpdateTableDescByIdWarp  根据表 ID 更新表描述。
1. RawRtdbbUpdateTableDescByNameWarp  根据表名称更新表描述。
1. RawRtdbbRecoverPointWarp  恢复已删除标签点
1. RawRtdbbPurgePointWarp  清除标签点
1. RawRtdbbGetRecycledPointsCountWarp  获取可回收标签点数量
1. RawRtdbbGetRecycledPointsWarp  获取可回收标签点 id 列表
1. RawRtdbbSearchRecycledPointsWarp  搜索符合条件的可回收标签点，使用标签点名时支持通配符
1. RawRtdbbGetRecycledPointPropertyWarp  获取可回收标签点的属性
1. RawRtdbbSearchRecycledPointsInBatchesWarp  分批搜索符合条件的可回收标签点，使用标签点名时支持通配符
1. RawRtdbbGetRecycledMaxPointPropertyWarp  按最大长度获取可回收标签点的属性
1. RawRtdbbClearRecyclerWarp  清空标签点回收站
1. RawRtdbbSubscribeTagsExWarp  标签点属性更改通知订阅
1. RawRtdbbCancelSubscribeTagsWarp  取消标签点属性更改通知订阅
1. RawRtdbbCreateNamedTypeWarp  创建自定义类型
1. RawRtdbbGetNamedTypesCountWarp  获取所有的自定义类型的总数
1. RawRtdbbGetAllNamedTypesWarp  获取所有的自定义类型
1. RawRtdbbGetNamedTypeWarp  获取自定义类型的所有字段
1. RawRtdbbRemoveNamedTypeWarp  删除自定义类型
1. RawRtdbbGetNamedTypeNamesPropertyWarp  根据标签点id查询标签点所对应的自定义类型的名字和字段总数
1. RawRtdbbGetRecycledNamedTypeNamesPropertyWarp  根据回收站标签点id查询标签点所对应的自定义类型的名字和字段总数
1. RawRtdbbGetNamedTypePointsCountWarp  获取该自定义类型的所有标签点个数
1. RawRtdbbGetBaseTypePointsCountWarp  获取该内置的基本类型的所有标签点个数
1. RawRtdbbModifyNamedTypeWarp  修改自定义类型名称,描述,字段名称,字段描述
1. RawRtdbWriteNamedTypeFieldByName32Warp  按名称填充自定义类型数值中字段的内容
1. RawRtdbWriteNamedTypeFieldByPos32Warp  按位置填充自定义类型数值中字段的内容
1. RawRtdbReadNamedTypeFieldByName32Warp  按名称提取自定义类型数值中字段的内容
1. RawRtdbReadNamedTypeFieldByPos32Warp  按位置提取自定义类型数值中字段的内容
1. RawRtdbNamedTypeNameFieldCheckWarp  检查自定义类型名称及字段命名是否符合规则
1. RawRtdbbGetMetaSyncInfoWarp  获取元数据同步信息
1. RawRtdbsGetSnapshots64Warp  批量读取开关量、模拟量快照数值
1. RawRtdbsPutSnapshots64Warp  批量写入开关量、模拟量快照数值
1. RawRtdbsFixSnapshots64Warp  批量覆盖写入开关量、模拟量快照数值 (时间戳相同的时候更新原有值)
1. RawRtdbsBackSnapshots64Warp  批量回溯快照, 批量将标签点的快照值vtmq改成传入的vtmq，如果传入的时间戳早于当前快照，会删除传入时间戳到当前快照的历史存储值。如果传入的时间戳等于或者晚于当前快照，什么也不做。
1. RawRtdbsGetCoorSnapshots64Warp  批量读取坐标实时数据
1. RawRtdbsPutCoorSnapshots64Warp  批量写入坐标实时数据
1. RawRtdbsFixCoorSnapshots64Warp  批量写入坐标实时数据(修复写入)
1. RawRtdbsGetBlobSnapshot64Warp  读取二进制/字符串实时数据
1. RawRtdbsGetBlobSnapshots64Warp  批量读取二进制/字符串实时数据
1. RawRtdbsPutBlobSnapshot64Warp  写入二进制/字符串实时数据
1. RawRtdbsPutBlobSnapshots64Warp  批量写入二进制/字符串实时数据
1. RawRtdbsGetDatetimeSnapshots64Warp  批量读取datetime类型标签点实时数据
1. RawRtdbsPutDatetimeSnapshots64Warp  批量插入datetime类型标签点数据
1. RawRtdbsSubscribeSnapshotsEx64Warp  批量标签点快照改变的通知订阅
1. RawRtdbsSubscribeDeltaSnapshots64Warp  批量标签点快照改变的通知订阅 (增量订阅，指的是数值超出一定变化后才会触发，可减少网络流量占用)
1. RawRtdbsChangeSubscribeSnapshotsWarp  批量修改订阅标签点信息
1. RawRtdbsCancelSubscribeSnapshotsWarp  取消标签点快照更改通知订阅
1. RawRtdbsGetNamedTypeSnapshot64Warp  获取自定义类型测点的单个快照
1. RawRtdbsGetNamedTypeSnapshots64Warp  批量获取自定义类型测点的快照
1. RawRtdbsPutNamedTypeSnapshot64Warp  写入单个自定义类型标签点的快照
1. RawRtdbsPutNamedTypeSnapshots64Warp  批量写入自定义类型标签点的快照
1. RawRtdbaGetArchivesCountWarp  获取存档文件数量
1. RawRtdbaCreateRangedArchive64Warp  新建指定时间范围的历史存档文件并插入到历史数据库
1. RawRtdbaAppendArchiveWarp  追加磁盘上的历史存档文件到历史数据库。
1. RawRtdbaRemoveArchiveWarp  从历史数据库中移出历史存档文件。
1. RawRtdbaShiftActivedWarp  切换活动文件
1. RawRtdbaGetArchivesWarp  获取存档文件的路径、名称、状态和最早允许写入时间。
1. RawRtdbaGetArchivesInfoWarp  获取存档信息
1. RawRtdbaGetArchivesPerfDataWarp  获取存档的实时信息(可用于分析存档处理的性能, 存档性能监控数据)
1. RawRtdbaGetArchivesStatusWarp  获取存档状态
1. RawRtdbaGetArchiveInfoWarp  获取存档文件及其附属文件的详细信息。
1. RawRtdbaUpdateArchiveWarp  修改存档文件的可配置项。
1. RawRtdbaArrangeArchiveWarp  整理存档文件，将同一标签点的数据块存放在一起以提高查询效率。
1. RawRtdbaReindexArchiveWarp  为存档文件重新生成索引，用于恢复数据。
1. RawRtdbaBackupArchiveWarp  备份主存档文件及其附属文件到指定路径
1. RawRtdbaMoveArchiveWarp  将存档文件移动到指定目录
1. RawRtdbaConvertIndexWarp  为存档文件转换索引格式。
1. RawRtdbaQueryBigJob64Warp  查询进程正在执行的后台任务类型、状态和进度
1. RawRtdbaCancelBigJobWarp  取消进程正在执行的后台任务
1. RawRtdbhArchivedValuesCount64Warp  获取单个标签点在一段时间范围内的存储值数量. (这个比较快，只需要读取索引，用于估算，无法处理标记删除)
1. RawRtdbhArchivedValuesRealCount64Warp  获取单个标签点在一段时间范围内的真实的存储值数量. (这个比较慢，需要全量读取)
1. RawRtdbhGetArchivedValues64Warp  读取单个标签点一段时间内的储存数据
1. RawRtdbhGetArchivedValuesBackward64Warp  逆向读取单个标签点一段时间内的储存数据
1. RawRtdbhGetArchivedCoorValues64Warp  读取单个标签点一段时间内的坐标型储存数据
1. RawRtdbhGetArchivedCoorValuesBackward64Warp  逆向读取单个标签点一段时间内的坐标型储存数据
1. RawRtdbhGetArchivedValuesInBatches64Warp  开始以分段返回方式读取一段时间内的储存数据
1. RawRtdbhGetNextArchivedValues64Warp  分段读取一段时间内的储存数据
1. RawRtdbhGetTimedValues64Warp  获取单个标签点的单调递增时间序列历史插值。
1. RawRtdbhGetTimedCoorValues64Warp  获取单个坐标标签点的单调递增时间序列历史插值。
1. RawRtdbhGetInterpoValues64Warp  获取单个标签点一段时间内等间隔历史插值
1. RawRtdbhGetIntervalValues64Warp  读取单个标签点某个时刻之后一定数量的等间隔内插值替换的历史数值
1. RawRtdbhGetSingleValue64Warp  读取单个标签点某个时间的历史数据
1. RawRtdbhGetSingleCoorValue64Warp  读取单个标签点某个时间的坐标型历史数据
1. RawRtdbhGetSingleBlobValue64Warp  读取单个标签点某个时间的二进制/字符串型历史数据
1. RawRtdbhGetArchivedBlobValues64Warp  读取单个标签点一段时间的二进制/字符串型历史数据
1. RawRtdbhGetArchivedBlobValuesFilt64Warp  读取并模糊搜索单个标签点一段时间的二进制/字符串型历史数据
1. RawRtdbhGetSingleDatetimeValue64Warp  读取单个标签点某个时间的datetime历史数据
1. RawRtdbhGetArchivedDatetimeValues64Warp  读取单个标签点一段时间的时间类型历史数据
1. RawRtdbhPutArchivedDatetimeValues64Warp  写入批量标签点批量时间型历史存储数据
1. RawRtdbhSummaryDataWarp  获取单个标签点一段时间内的统计值。
1. RawRtdbhSummaryDataInBatchesWarp  分批获取单一标签点一段时间内的统计值
1. RawRtdbhGetPlotValues64Warp  获取单个标签点一段时间内用于绘图的历史数据
1. RawRtdbhGetCrossSectionValues64Warp  获取批量标签点在某一时间的历史断面数据
1. RawRtdbhGetArchivedValuesFilt64Warp  读取单个标签点在一段时间内经复杂条件筛选后的历史储存值
1. RawRtdbhGetIntervalValuesFilt64Warp  读取单个标签点某个时刻之后经复杂条件筛选后一定数量的等间隔内插值替换的历史数值
1. RawRtdbhGetInterpoValuesFilt64Warp  获取单个标签点一段时间内经复杂条件筛选后的等间隔插值
1. RawRtdbhSummaryDataFiltWarp  获取单个标签点一段时间内经复杂条件筛选后的统计值
1. RawRtdbhSummaryDataFiltInBatchesWarp  分批获取单一标签点一段时间内经复杂条件筛选后的统计值
1. RawRtdbhUpdateValue64Warp  修改单个标签点某一时间的历史存储值.
1. RawRtdbhUpdateCoorValue64Warp  修改单个标签点某一时间的历史存储值(坐标类型)
1. RawRtdbhRemoveValue64Warp  删除单个标签点某个时间的历史存储值
1. RawRtdbhRemoveValues64Warp  删除单个标签点一段时间内的历史存储值
1. RawRtdbhPutSingleValue64Warp  写入单个标签点在某一时间的历史数据。
1. RawRtdbhPutSingleCoorValue64Warp  写入单个标签点在某一时间的坐标型历史数据。
1. RawRtdbhPutSingleBlobValue64Warp  写入单个二进制/字符串标签点在某一时间的历史数据
1. RawRtdbhPutArchivedValues64Warp  写入批量标签点批量历史存储数据
1. RawRtdbhPutArchivedCoorValues64Warp  写入批量标签点批量坐标型历史存储数据
1. RawRtdbhPutSingleDatetimeValue64Warp  写入单个datetime标签点在某一时间的历史数据
1. RawRtdbhPutArchivedBlobValues64Warp  写入批量标签点批量字符串型历史存储数据
1. RawRtdbhFlushArchivedValuesWarp  将标签点未写满的补历史缓存页写入存档文件中。
1. RawRtdbhGetSingleNamedTypeValue64Warp  读取单个自定义类型标签点某个时间的历史数据
1. RawRtdbhGetArchivedNamedTypeValues64Warp  连续读取自定义类型标签点的历史数据
1. RawRtdbhPutSingleNamedTypeValue64Warp  写入自定义类型标签点的单个历史事件
1. RawRtdbhPutArchivedNamedTypeValues64Warp  批量补写自定义类型标签点的历史事件
1. RawRtdbeComputeHistory64Warp  重算或补算批量计算标签点历史数据
1. RawRtdbbGetEquationByFileNameWarp  根据文件名获取方程式
1. RawRtdbbGetEquationByIdWarp  根ID径获取方程式
1. RawRtdbeGetEquationGraphCountWarp  根据标签点 id 获取相关联方程式键值对数量
1. RawRtdbeGetEquationGraphDatasWarp  根据标签点 id 获取相关联方程式键值对数据
1. RawRtdbpGetPerfTagsCountWarp  获取Perf服务中支持的性能计数点的数量
1. RawRtdbpGetPerfTagsInfoWarp  根据性能计数点ID获取相关的性能计数点信息
1. RawRtdbpGetPerfValues64Warp  批量读取性能计数点的当前快照数值
