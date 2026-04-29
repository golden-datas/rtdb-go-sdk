# rtdb_api

## 注意

由于本库使用了CGO技术，需要编译一定量的C代码，因此编译时间较长，尤其是在Windows环境下，请耐心等待

## 编译环境

* 在Linux平台中，需安装Golang编译器、GCC工具链
* 在Windows平台中，需安装Golang编译器、MicrosoftVisualC++BuildTools工具链(可通过VisualStudio2022安装)

## 版本

格式：SO库版本号_Go库版本号

* v4.0.11_0.1.0: 第一个试用版本
* v4.0.15_0.2.0: 新增了一些Raw函数

## 一些基本概念

* API库：这里特指连接数据库的so库，这个so库是用C写的，是负责和数据库进行通信的客户端，本包封装了这个so库的接口，使用FFI+Warp技术进行的。
* CGO: 由于本库使用了部分C语言，因此编译的时候需要开启CGO，否则会导致无法编译，开启CGO命令:```go env -w CGO_ENABLED=1```

## 分层

```text
+---------------------+
|       C Header      | // 原版的C库函数
+---------------------+
|        api.h        | // 在C库函数的基础上，封装了一层Wrapper函数，用于动态加载C库函数的指针
+---------------------+
|       api.go        | // 在api.h的基础上，将Wrapper函数封装成Go函数
+---------------------+
|       easy.go       | // 在api.go的基础上，对原始API函数进一步简化，达到简单易用的效果
+---------------------+
```

## 代码结构

* cinclude: C代码的.h部分，里面包含了一些必要的C头文件
* clibrary: C代码的(.so/.dll)部分，里面包含了跨平台的动态库(linux_amd64、linux_arm64、windows_amd64)
* api.go: 基于C代码封装的原始API，函数名均以Raw开头，由于是基于C原始代码1比1封装，因此缺乏对象化，相对难用但功能全性能高
* api_test.go: api.go中封装函数的代码示例
* easy.go: 基于api.go进行二次封装的代码，更加简单易用，符合Golang语言风格，推荐使用easy封装的代码，更加简洁明了
* easy_test.go: easy.go中封装函数的代码示例

## 注意

尽量避免使用Raw开头的函数，此为原始C函数的Go封装，属于中间层代码，但是由于他的全面性和标准性，这里还是进行了保留并且对外提供调用方式，有极致性能需求的情况下可酌情调用。

## 最简调用示例

```go
package main

import (
	"fmt"
	"time"

	"github.com/kkbase/rtdb_api"
)

const (
	ip       = "127.0.0.1"
	port     = int32(6327)
	username = "sa"
	password = "golden"
)

func main() {
	// 登录数据库
	conn, err := rtdb_api.Login(ip, port, username, password, rtdb_api.RtdbPrecisionNano)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		// 登出数据库
		_ = conn.Logout()
	}()

	// 获取 Client端 版本
	fmt.Println(conn.GetClientVersion())

	// 创建表
	table, err := conn.CreateTable("example_table", "example table desc")
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		// 删除表
		_ = conn.DeleteTable(table.ID)
	}()

	// 创建点
	info := rtdb_api.NewPointInfo("example_point", table.ID, rtdb_api.ValueTypeFloat64, rtdb_api.PointBase, rtdb_api.RtdbPrecisionNano, "", "")
	info.SetLimit(-100, 100, 0)
	info, err = conn.AddPoint(info)
	if err != nil {
		fmt.Println("添加点失败: ", err)
		return
	}

	// 删除点
	defer func() { _ = conn.DeletePoint(info.ID) }()

	// 写入数据
	n := 10
	for i := 0; i < n; i++ {
		err := conn.WriteValue(info, false, info.NewNowTVQ(float64(i), rtdb_api.Quality(0)))
		if err != nil {
			fmt.Println("写入数据失败：", err)
			return
		}
		if i != n-1 {
			time.Sleep(time.Second)
		}
	}

	// 获取实时值
	ptvq, err := conn.ReadLast(info)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ptvq)
}
```

# go.mod引用版本说明

引用本库的时候需注意版本号，请使用最新版本, 可执行 ```git ls-remote --tags git@github.com:kkbase/rtdb_api.git```
命令行查看最新版本号

```text
require github.com/kkbase/rtdb_api v0.1.0
// 注意这里的 v0.1.0 ， 需替换成最新版本
```

## Easy可用函数列表

1. Login 登录数据库
2. (RtdbConnect)Logout 登出数据库
3. (RtdbConnect)SetLocation 设置客户端时区，不设置默认为亚洲东八
4. (RtdbConnect)GetClientVersion 获取客户端版本
5. (RtdbConnect)SetClientOption 设置客户端参数
6. (RtdbConnect)GetServerOption 获取服务端参数
7. (RtdbConnect)SetServerOption 设置服务端参数
8. (RtdbConnect)GetSocketInfos 获取服务端SocketInfo列表，单机服务端返回一个SocketInfo列表，双活服务端返回两个SocketInfo列表
9. (RtdbConnect)GetOwnSocketInfo 获取当前连接的SocketInfo，单机服务端返回一个SocketInfo，双活服务端返回两个SocketInfo
10. (RtdbConnect)SetSocketTimeout 设置Socket超时时间
11. (RtdbConnect)KillSocket 断开Socket
12. (RtdbConnect)AddIpBlackList 添加IP黑名单项
13. (RtdbConnect)UpdateIpBlackList 更新连接黑名单项
14. (RtdbConnect)DeleteIpBlackList 删除连接黑名单项
15. (RtdbConnect)GetIpBlackLists 获得连接黑名单列表
16. (RtdbConnect)AddIpWhiteList 添加连接白名单
17. (RtdbConnect)UpdateIpWhiteList 更新连接白名单
18. (RtdbConnect)DeleteIpWhiteList 删除白名单
19. (RtdbConnect)GetIpWhiteLists 获取连接白名单列表
20. (RtdbConnect)UpdatePassword 修改用户密码
21. (RtdbConnect)UpdateOwnPassword 修改自己的密码
22. (RtdbConnect)GetPriv 获取连接权限
23. (RtdbConnect)SetPriv 设置连接权限
24. (RtdbConnect)AddUser 添加用户
25. (RtdbConnect)DeleteUser 删除用户
26. (RtdbConnect)LockUser 锁定用户
27. (RtdbConnect)GetUsers 获取用户列表
28. (RtdbConnect)AddNamedType 创建自定义类型
29. (RtdbConnect)DeleteNamedType 删除自定义类型
30. (RtdbConnect)GetNamedType 获取自定义类型
31. (RtdbConnect)GetNamedTypes 获取自定义类型列表
32. (RtdbConnect)UpdateNamedType 修改自定义类型
33. (RtdbConnect)ServerHostTime 服务端主机时间
34. (RtdbConnect)DurationToString 时间段转字符串, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
35. (RtdbConnect)StringToDuration 字符串转时间段, 这个是服务端的时间段字符串格式，和通用时间段字符串有区别, 具体如下：
36. (RtdbConnect)StringToTime 字符串转时间戳
37. (RtdbConnect)GetQualityDesc 获取质量码说明
38. (RtdbConnect)GetDriveLetterList 获取盘符列表, windows平台是C、D、E、F这些盘符，linux平台是 / 盘符
39. (RtdbConnect)GetDirItemList 获取目录项列表
40. (RtdbConnect)CreateDir 创建目录
41. (RtdbConnect)ReadFile 读取文件
42. (RtdbConnect)CreateTable 创建表
43. (RtdbConnect)DeleteTable 删除表
44. (RtdbConnect)GetTable 获取表
45. (RtdbConnect)GetTables 获取表列表
46. (RtdbConnect)UpdateTableName 更新表名
47. (RtdbConnect)UpdateTableDesc 更新表描述
48. (RtdbConnect)AddPoint 创建点
49. (RtdbConnect)AddPoints 批量创建点
50. (RtdbConnect)DeletePoint 删除点
51. (RtdbConnect)UpdatePoint 更新点
52. (RtdbConnect)GetPoints 批量获取标签点
53. (RtdbConnect)GetPoint 获取点
54. (RtdbConnect)FindPoints 根据 表名.点名 搜索标签点
55. (RtdbConnect)MovePoint 移动点到指定表
56. (RtdbConnect)SearchPoint 分页搜索点
57. (RtdbConnect)ClearRecycler 清空回收站
58. (RtdbConnect)GetRecycledPoints 分段获取回收站中的点
59. (RtdbConnect)RecoverPoint 从回收站中恢复点到某个表
60. (RtdbConnect)PurgePoint 从回收站中清除点
61. (RtdbConnect)SearchRecycledPoint 从回收站中搜索点
62. (RtdbConnect)GetPointCountFromValueType 获取某个数据类型的点个数 (可以是内置类型，也可以是自定义类型)
63. (RtdbConnect)WriteValue 写入值
64. (RtdbConnect)WriteValues 批量写入值
65. (RtdbConnect)WriteSection 写断面(批量写入多个Point，每个Point写入一个TVQ)
66. (RtdbConnect)ReadLasts 批量读取实时快照值(当前标签点最后一个写入的TVQ)
67. (RtdbConnect)ReadLast 读取实时快照值(当前标签点最后一个写入的TVQ)
68. (RtdbConnect)ReadValue 读取单个TVQ
69. (RtdbConnect)ReadRange 读取某个时间段内的TVQ
70. (RtdbConnect)ReadPlot 读取用于绘图的TVQ
71. (RtdbConnect)ReadTimed 获取差值，每个差值都需要指定一个确定的时间戳
72. (RtdbConnect)ReadInterpo 获取差值, 会自动将start、end等分成count个时间戳，然后取这些时间戳的差值
73. (RtdbConnect)ReadInterval 读取从start开始的等间隔差值
74. (RtdbConnect)ReadSection 读取断面
75. (RtdbConnect)ReadSummary 获取统计值 (从start到end的统计值)
76. (RtdbConnect)ReadBatchesSummary 获取等间隔统计值 (在start和end之间，按照interval作为时间间隔计算每个间隔的统计值)
77. (RtdbConnect)RemoveValue 删除点值
78. (RtdbConnect)RemoveRangeValues 批量删除点值(从start到end的范围)
79. (RtdbConnect)UpdateValue 更新点值
80. (RtdbConnect)FlushArchivedValues 刷新历史缓存(把标签点的缓存进行手动归档)
81. (RtdbConnect)QueryBigJob 查询正在进行的大任务
82. (RtdbConnect)CancelBigJob 取消正在执行的大任务
83. (RtdbConnect)JobMessage 获取任务描述
84. (RtdbConnect)ComputeHistory 重算｜补算 补历史值或者修改历史值之后，对应的计算点可以按需进行重算|补算，保证计算点的数值正确
85. (RtdbConnect)GetPointEquation 获取标签点方程式
86. (RtdbConnect)GetEquationGraph 获取标签点对应的方程式关联关系图(
    方程式本身是一个有向无环图，数据在这个图内的公式之间流转)
87. (RtdbConnect)GetPerfPointInfo 获取性能监控点的说明信息
88. (RtdbConnect)GetPerfPointValue 获取性能监控点的实时值
89. (RtdbConnect)CreateRangedArchive 创建存档文件(按照时间范围创建)
90. (RtdbConnect)AppendArchive 存档文件入列(将存档文件插入到存档队列中,
    队列中的存档文件按照start、end依次排列，注意时间段不能有交叠的部分)
91. (RtdbConnect)RemoveArchive 存档文件解列
92. (RtdbConnect)ShiftActived 切换活动文件
93. (RtdbConnect)GetArchives 获取存档文件的基本信息(返回全部存档文件的基本信息)
94. (RtdbConnect)UpdateArchive 更新存档文件选项
95. (RtdbConnect)ArrangeArchive 整理存档文件(会重新调整数据块分布，使存档文件更加紧凑)
96. (RtdbConnect)ReindexArchive 重建存档文件索引(用于进行数据恢复)
97. (RtdbConnect)ConvertIndex 存档文件索引转换格式(老版本索引格式转换为新版本，新版本索引更快)
98. (RtdbConnect)BackupArchive 备份存档
99. (RtdbConnect)MoveArchive 移动存档
100. (RtdbConnect)CreateDatagram 创建数据流订阅
101. (RtdbConnect)RemoveDatagram 取消数据流订阅
102. (RtdbConnect)RecvDatagram 从订阅数据流中获取数据
103. (RtdbConnect)SubscribeTags 订阅标签点属性更新
104. (RtdbConnect)CancelSubscribeTags 取消订阅标签点属性更新
105. (RtdbConnect)SubscribeSnapshots 订阅快照，只要快照发生变化，就会触发订阅
106. (RtdbConnect)SubscribeDeltaSnapshots 订阅Delta快照，快照变化需要超过Delta，才会触发订阅，这样可以节约流量
107. (RtdbConnect)ChangeSubscribeSnapshots 修改快照订阅设置，新增或删除标签点
108. (RtdbConnect)CancelSubscribeSnapshots 取消快照订阅
109. (RtdbConnect)ReadBatches 开始以分段返回方式读取一段时间内的储存数据
110. (RtdbConnect)ReadNext 分段读取一段时间内的储存数据

# Raw可用函数列表

1. RawRtdbGetApiVersionWarp 返回 ApiVersion 版本号
2. RawRtdbSetOptionWarp 配置 API库 的行为参数，详见 RtdbApiOption 枚举
3. RawRtdbCreateDatagramHandleWarp 创建数据流
4. RawRtdbRemoveDatagramHandleWarp 删除数据流
5. RawRtdbRecvDatagramWarp 接收数据流
6. RawRtdbConnectWarp 建立同 RTDB 数据库的网络连接, 注意这里只是创建连接，并没有进行用户登陆
7. RawRtdbLoginWarp 以有效帐户登录
8. RawRtdbDisconnectWarp 断开同 RTDB 数据平台的连接
9. RawRtdbGetDbInfo1Warp 获得字符串型数据库系统参数
10. RawRtdbGetDbInfo2Warp 获得整型数据库系统参数
11. RawRtdbSetDbInfo1Warp 设置字符串型数据库系统参数
12. RawRtdbSetDbInfo2Warp 设置整型数据库系统参数
13. RawRtdbConnectionCountWarp 获取 RTDB 服务器当前连接个数
14. RawRtdbGetConnectionsWarp 列出 RTDB 服务器的所有socket连接句柄, 注意这里指的是socket连接，区分于ConnectHandle
15. RawRtdbGetOwnConnectionWarp 获取当前连接的socket句柄
16. RawRtdbGetConnectionInfoIpv6Warp 获取 RTDB 服务器指定连接的ipv6版本
17. RawRtdbOsType 获取连接句柄所连接的服务器操作系统类型
18. RawRtdbChangePasswordWarp 修改用户帐户口令
19. RawRtdbChangeMyPasswordWarp 用户修改自己帐户口令
20. RawRtdbGetPrivWarp 获取连接权限
21. RawRtdbChangePrivWarp 修改用户帐户权限, 只有管理员有修改权限
22. RawRtdbAddUserWarp 添加用户帐户
23. RawRtdbRemoveUserWarp 删除用户帐户
24. RawRtdbLockUserWarp 启用或禁用用户, 只有管理员有启用禁用权限
25. RawRtdbGetUsersWarp 获得所有用户
26. RawRtdbAddBlacklistWarp 添加连接黑名单项
27. RawRtdbUpdateBlacklistWarp 更新连接连接黑名单项
28. RawRtdbRemoveBlacklistWarp 删除连接黑名单项
29. RawRtdbGetBlacklistWarp 获得连接黑名单
30. RawRtdbAddAuthorizationWarp 添加信任连接段
31. RawRtdbUpdateAuthorizationWarp 更新信任连接段
32. RawRtdbRemoveAuthorizationWarp 删除信任连接段
33. RawRtdbGetAuthorizationsWarp 获得所有信任连接段
34. RawRtdbHostTime64Warp 获取 RTDB 服务器当前UTC时间
35. RawRtdbFormatTimespanWarp 根据时间跨度值生成时间格式字符串, 如：输入10， 输出10s, 输入60，输出1n
36. RawRtdbParseTimespanWarp 根据时间格式字符串解析时间跨度值, 如：输入2n，输出120，表示2分钟
37. RawRtdbParseTimeWarp 根据时间格式字符串解析时间值
38. RawRtdbFormatMessageWarp 获取 Rtdb API 调用返回值的简短描述(错误码对应的Desc)
39. RawRtdbJobMessageWarp 获取任务的简短描述
40. RawRtdbSetTimeoutWarp 设置连接超时时间
41. RawRtdbGetTimeoutWarp 获得连接超时时间
42. RawRtdbKillConnectionWarp 断开已知连接
43. RawRtdbGetLogicalDriversWarp 获得逻辑盘符
44. RawRtdbOpenPathWarp 打开目录以便遍历其中的文件和子目录。
45. RawRtdbReadPath64Warp 读取目录中的文件或子目录
46. RawRtdbClosePathWarp 关闭当前遍历的目录
47. RawRtdbMkdirWarp 建立目录
48. RawRtdbGetFileSizeWarp 获得指定服务器端文件的大小
49. RawRtdbReadFileWarp 读取服务器端指定文件的内容
50. RawRtdbGetMaxBlobLenWarp 取得数据库允许的blob与str类型测点的最大长度
51. RawRtdbFormatQualityWarp 取得质量码对应的定义
52. RawRtdbJudgeConnectStatusWarp 判断连接是否可用
53. RawRtdbbAppendTableWarp 添加新表
54. RawRtdbbRemoveTableByIdWarp 根据表 id 删除表及表中标签点
55. RawRtdbbRemoveTableByNameWarp 根据表名删除表及表中标签点
56. RawRtdbbTablesCountWarp 取得标签点表总数
57. RawRtdbbGetTablesWarp 取得所有标签点表的ID
58. RawRtdbbGetTableSizeByIdWarp 根据表 id 获取表中包含的标签点数量(大概数量, 包含被标记删除的点)
59. RawRtdbbGetTableSizeByNameWarp 根据表名称获取表中包含的标签点数量(大概数量, 包含被标记删除的点)
60. RawRtdbbGetTableRealSizeByIdWarp 根据表 id 获取表中实际包含的标签点数量(实际数量, 不含被删除的点)
61. RawRtdbbGetTablePropertyByIdWarp 根据标签点表 id 获取表属性
62. RawRtdbbGetTablePropertyByNameWarp 根据表名获取标签点表属性
63. RawRtdbbInsertMaxPointWarp 使用最大长度的完整属性集来创建单个标签点
64. RawRtdbbRemovePointByIdWarp 根据 id 删除单个标签点
65. RawRtdbbRemovePointByNameWarp 根据标签点全名删除单个标签点
66. RawRtdbbInsertMaxPointsWarp 使用最大长度的完整属性集来批量创建标签点
67. RawRtdbbInsertNamedTypePointWarp 使用完整的属性集来创建单个自定义数据类型标签点
68. RawRtdbbMovePointByIdWarp 根据 id 移动单个标签点到其他表
69. RawRtdbbGetMaxPointsPropertyWarp 按最大长度批量获取标签点属性
70. RawRtdbbSearchInBatchesWarp 分批继续搜索符合条件的标签点，使用标签点名时支持通配符
71. RawRtdbbSearchExWarp 搜索符合条件的标签点，使用标签点名时支持通配符
72. RawRtdbbSearchPointsCountWarp 搜索符合条件的标签点，获取标签点数，使用标签点名时支持通配符
73. RawRtdbbUpdateMaxPointPropertyWarp 按最大长度更新单个标签点属性
74. RawRtdbbFindPointsExWarp 根据 "表名.标签点名" 格式批量获取标签点标识
75. RawRtdbbSortPointsWarp 根据标签属性字段对标签点标识进行排序
76. RawRtdbbUpdateTableNameWarp 根据表 ID 更新表名称。
77. RawRtdbbUpdateTableDescByIdWarp 根据表 ID 更新表描述。
78. RawRtdbbUpdateTableDescByNameWarp 根据表名称更新表描述。
79. RawRtdbbRecoverPointWarp 恢复已删除标签点
80. RawRtdbbPurgePointWarp 清除标签点
81. RawRtdbbGetRecycledPointsCountWarp 获取可回收标签点数量
82. RawRtdbbGetRecycledPointsWarp 获取可回收标签点 id 列表
83. RawRtdbbSearchRecycledPointsInBatchesWarp 分批搜索符合条件的可回收标签点，使用标签点名时支持通配符
84. RawRtdbbGetRecycledMaxPointPropertyWarp 按最大长度获取可回收标签点的属性
85. RawRtdbbClearRecyclerWarp 清空标签点回收站
86. RawRtdbbSubscribeTagsExWarp 标签点属性更改通知订阅
87. RawRtdbbCancelSubscribeTagsWarp 取消标签点属性更改通知订阅
88. RawRtdbbCreateNamedTypeWarp 创建自定义类型
89. RawRtdbbGetNamedTypesCountWarp 获取所有的自定义类型的总数
90. RawRtdbbGetAllNamedTypesWarp 获取所有的自定义类型
91. RawRtdbbGetNamedTypeWarp 获取自定义类型的所有字段
92. RawRtdbbRemoveNamedTypeWarp 删除自定义类型
93. RawRtdbbGetNamedTypeNamesPropertyWarp 根据标签点id查询标签点所对应的自定义类型的名字和字段总数
94. RawRtdbbGetRecycledNamedTypeNamesPropertyWarp 根据回收站标签点id查询标签点所对应的自定义类型的名字和字段总数
95. RawRtdbbGetNamedTypePointsCountWarp 获取该自定义类型的所有标签点个数
96. RawRtdbbGetBaseTypePointsCountWarp 获取该内置的基本类型的所有标签点个数
97. RawRtdbbModifyNamedTypeWarp 修改自定义类型名称,描述,字段名称,字段描述
98. RawRtdbbGetMetaSyncInfoWarp 获取元数据同步信息
99. RawRtdbsGetSnapshots64Warp 批量读取开关量、模拟量快照数值
100. RawRtdbsPutSnapshots64Warp 批量写入开关量、模拟量快照数值
101. RawRtdbsFixSnapshots64Warp 批量覆盖写入开关量、模拟量快照数值 (时间戳相同的时候更新原有值)
102. RawRtdbsBackSnapshots64Warp 批量回溯快照,
     批量将标签点的快照值vtmq改成传入的vtmq，如果传入的时间戳早于当前快照，会删除传入时间戳到当前快照的历史存储值。如果传入的时间戳等于或者晚于当前快照，什么也不做。
103. RawRtdbsGetCoorSnapshots64Warp 批量读取坐标实时数据
104. RawRtdbsPutCoorSnapshots64Warp 批量写入坐标实时数据
105. RawRtdbsFixCoorSnapshots64Warp 批量写入坐标实时数据(修复写入)
106. RawRtdbsGetBlobSnapshot64Warp 读取二进制/字符串实时数据
107. RawRtdbsGetBlobSnapshots64Warp 批量读取二进制/字符串实时数据
108. RawRtdbsPutBlobSnapshot64Warp 写入二进制/字符串实时数据
109. RawRtdbsPutBlobSnapshots64Warp 批量写入二进制/字符串实时数据
110. RawRtdbsGetDatetimeSnapshots64Warp 批量读取datetime类型标签点实时数据
111. RawRtdbsPutDatetimeSnapshots64Warp 批量插入datetime类型标签点数据
112. RawRtdbsSubscribeSnapshotsEx64Warp 批量标签点快照改变的通知订阅
113. RawRtdbsSubscribeDeltaSnapshots64Warp 批量标签点快照改变的通知订阅 (
     增量订阅，指的是数值超出一定变化后才会触发，可减少网络流量占用)
114. RawRtdbsChangeSubscribeSnapshotsWarp 批量修改订阅标签点信息
115. RawRtdbsCancelSubscribeSnapshotsWarp 取消标签点快照更改通知订阅
116. RawRtdbsGetNamedTypeSnapshot64Warp 获取自定义类型测点的单个快照
117. RawRtdbsGetNamedTypeSnapshots64Warp 批量获取自定义类型测点的快照
118. RawRtdbsPutNamedTypeSnapshot64Warp 写入单个自定义类型标签点的快照
119. RawRtdbsPutNamedTypeSnapshots64Warp 批量写入自定义类型标签点的快照
120. RawRtdbaGetArchivesCountWarp 获取存档文件数量
121. RawRtdbaCreateRangedArchive64Warp 新建指定时间范围的历史存档文件并插入到历史数据库
122. RawRtdbaAppendArchiveWarp 追加磁盘上的历史存档文件到历史数据库。
123. RawRtdbaRemoveArchiveWarp 从历史数据库中移出历史存档文件。
124. RawRtdbaShiftActivedWarp 切换活动文件
125. RawRtdbaGetArchivesWarp 获取存档文件的路径、名称、状态和最早允许写入时间。
126. RawRtdbaGetArchivesInfoWarp 获取存档信息
127. RawRtdbaGetArchivesPerfDataWarp 获取存档的实时信息(可用于分析存档处理的性能, 存档性能监控数据)
128. RawRtdbaGetArchivesStatusWarp 获取存档状态
129. RawRtdbaGetArchiveInfoWarp 获取存档文件及其附属文件的详细信息。
130. RawRtdbaUpdateArchiveWarp 修改存档文件的可配置项。
131. RawRtdbaArrangeArchiveWarp 整理存档文件，将同一标签点的数据块存放在一起以提高查询效率。
132. RawRtdbaReindexArchiveWarp 为存档文件重新生成索引，用于恢复数据。
133. RawRtdbaBackupArchiveWarp 备份主存档文件及其附属文件到指定路径
134. RawRtdbaMoveArchiveWarp 将存档文件移动到指定目录
135. RawRtdbaConvertIndexWarp 为存档文件转换索引格式。
136. RawRtdbaQueryBigJob64Warp 查询进程正在执行的后台任务类型、状态和进度
137. RawRtdbaCancelBigJobWarp 取消进程正在执行的后台任务
138. RawRtdbhArchivedValuesCount64Warp 获取单个标签点在一段时间范围内的存储值数量. (
     这个比较快，只需要读取索引，用于估算，无法处理标记删除)
139. RawRtdbhArchivedValuesRealCount64Warp 获取单个标签点在一段时间范围内的真实的存储值数量. (这个比较慢，需要全量读取)
140. RawRtdbhGetArchivedValues64Warp 读取单个标签点一段时间内的储存数据
141. RawRtdbhGetArchivedValuesBackward64Warp 逆向读取单个标签点一段时间内的储存数据
142. RawRtdbhGetArchivedCoorValues64Warp 读取单个标签点一段时间内的坐标型储存数据
143. RawRtdbhGetArchivedCoorValuesBackward64Warp 逆向读取单个标签点一段时间内的坐标型储存数据
144. RawRtdbhGetArchivedValuesInBatches64Warp 开始以分段返回方式读取一段时间内的储存数据
145. RawRtdbhGetNextArchivedValues64Warp 分段读取一段时间内的储存数据
146. RawRtdbhGetTimedValues64Warp 获取单个标签点的单调递增时间序列历史插值。
147. RawRtdbhGetTimedCoorValues64Warp 获取单个坐标标签点的单调递增时间序列历史插值。
148. RawRtdbhGetInterpoValues64Warp 获取单个标签点一段时间内等间隔历史插值
149. RawRtdbhGetIntervalValues64Warp 读取单个标签点某个时刻之后一定数量的等间隔内插值替换的历史数值
150. RawRtdbhGetSingleValue64Warp 读取单个标签点某个时间的历史数据
151. RawRtdbhGetSingleCoorValue64Warp 读取单个标签点某个时间的坐标型历史数据
152. RawRtdbhGetSingleBlobValue64Warp 读取单个标签点某个时间的二进制/字符串型历史数据
153. RawRtdbhGetArchivedBlobValues64Warp 读取单个标签点一段时间的二进制/字符串型历史数据
154. RawRtdbhGetArchivedBlobValuesFilt64Warp 读取并模糊搜索单个标签点一段时间的二进制/字符串型历史数据
155. RawRtdbhGetSingleDatetimeValue64Warp 读取单个标签点某个时间的datetime历史数据
156. RawRtdbhGetArchivedDatetimeValues64Warp 读取单个标签点一段时间的时间类型历史数据
157. RawRtdbhPutArchivedDatetimeValues64Warp 写入批量标签点批量时间型历史存储数据
158. RawRtdbhSummaryDataWarp 获取单个标签点一段时间内的统计值。
159. RawRtdbhSummaryDataInBatchesWarp 分批获取单一标签点一段时间内的统计值
160. RawRtdbhGetPlotValues64Warp 获取单个标签点一段时间内用于绘图的历史数据
161. RawRtdbhGetCrossSectionValues64Warp 获取批量标签点在某一时间的历史断面数据
162. RawRtdbhGetArchivedValuesFilt64Warp 读取单个标签点在一段时间内经复杂条件筛选后的历史储存值
163. RawRtdbhGetIntervalValuesFilt64Warp 读取单个标签点某个时刻之后经复杂条件筛选后一定数量的等间隔内插值替换的历史数值
164. RawRtdbhGetInterpoValuesFilt64Warp 获取单个标签点一段时间内经复杂条件筛选后的等间隔插值
165. RawRtdbhSummaryDataFiltWarp 获取单个标签点一段时间内经复杂条件筛选后的统计值
166. RawRtdbhSummaryDataFiltInBatchesWarp 分批获取单一标签点一段时间内经复杂条件筛选后的统计值
167. RawRtdbhUpdateValue64Warp 修改单个标签点某一时间的历史存储值.
168. RawRtdbhUpdateCoorValue64Warp 修改单个标签点某一时间的历史存储值(坐标类型)
169. RawRtdbhRemoveValue64Warp 删除单个标签点某个时间的历史存储值
170. RawRtdbhRemoveValues64Warp 删除单个标签点一段时间内的历史存储值
171. RawRtdbhPutSingleValue64Warp 写入单个标签点在某一时间的历史数据。
172. RawRtdbhPutSingleCoorValue64Warp 写入单个标签点在某一时间的坐标型历史数据。
173. RawRtdbhPutSingleBlobValue64Warp 写入单个二进制/字符串标签点在某一时间的历史数据
174. RawRtdbhPutArchivedValues64Warp 写入批量标签点批量历史存储数据
175. RawRtdbhPutArchivedCoorValues64Warp 写入批量标签点批量坐标型历史存储数据
176. RawRtdbhPutSingleDatetimeValue64Warp 写入单个datetime标签点在某一时间的历史数据
177. RawRtdbhPutArchivedBlobValues64Warp 写入批量标签点批量字符串型历史存储数据
178. RawRtdbhFlushArchivedValuesWarp 将标签点未写满的补历史缓存页写入存档文件中。
179. RawRtdbhGetSingleNamedTypeValue64Warp 读取单个自定义类型标签点某个时间的历史数据
180. RawRtdbhGetArchivedNamedTypeValues64Warp 连续读取自定义类型标签点的历史数据
181. RawRtdbhPutSingleNamedTypeValue64Warp 写入自定义类型标签点的单个历史事件
182. RawRtdbhPutArchivedNamedTypeValues64Warp 批量补写自定义类型标签点的历史事件
183. RawRtdbeComputeHistory64Warp 重算或补算批量计算标签点历史数据
184. RawRtdbbGetEquationByFileNameWarp 根据文件名获取方程式
185. RawRtdbbGetEquationByIdWarp 根ID径获取方程式
186. RawRtdbeGetEquationGraphCountWarp 根据标签点 id 获取相关联方程式键值对数量
187. RawRtdbeGetEquationGraphDatasWarp 根据标签点 id 获取相关联方程式键值对数据
188. RawRtdbpGetPerfTagsCountWarp 获取Perf服务中支持的性能计数点的数量
189. RawRtdbpGetPerfTagsInfoWarp 根据性能计数点ID获取相关的性能计数点信息
190. RawRtdbpGetPerfValues64Warp 批量读取性能计数点的当前快照数值
