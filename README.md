# rtdb_api

## 备注
由于本库使用了CGO技术，需要编译一定量的C代码，因此编译时间较长，尤其是在Windows环境下，请耐心等待

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

## 可用函数列表
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
86. (RtdbConnect)GetEquationGraph 获取标签点对应的方程式关联关系图(方程式本身是一个有向无环图，数据在这个图内的公式之间流转)
87. (RtdbConnect)GetPerfPointInfo 获取性能监控点的说明信息
88. (RtdbConnect)GetPerfPointValue 获取性能监控点的实时值
89. (RtdbConnect)CreateRangedArchive 创建存档文件(按照时间范围创建)
90. (RtdbConnect)AppendArchive 存档文件入列(将存档文件插入到存档队列中, 队列中的存档文件按照start、end依次排列，注意时间段不能有交叠的部分)
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