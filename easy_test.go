package rtdb_api

import (
	"fmt"
	"path"
	"testing"
	"time"
)

// 用户登录/登出
func TestLoginLogout(t *testing.T) {
	// 登录
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	// 登出
	defer func() { _ = conn.Logout() }()

	fmt.Println(conn.SyncInfos, conn.StringBlobMaxLen)
}

// 获取客户端版本
func TestRtdbConnect_GetClientVersion(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取客户端版本
	version, err := conn.GetClientVersion()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(version)
}

// 设置客户端选项
func TestRtdbConnect_SetClientOption(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 设置客户端选项
	err = conn.SetClientOption(RtdbApiOptionAutoReconn, 0)
	if err != nil {
		t.Error(err)
	}
}

// 服务端选项
func TestRtdbConnect_GetSetServerOption(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取服务端选项
	opt, err := conn.GetServerOption(RtdbParamLockedPagesMem)
	if err != nil {
		t.Error("获取服务端选项失败", err)
		return
	}
	fmt.Println(opt.GetLiteralValue())

	// 设置服务端选项
	err = conn.SetServerOption(RtdbParamLockedPagesMem, NewServerOption(opt.GetLiteralValue()))
	if err != nil {
		t.Error("设置服务端选项失败", err)
		return
	}
}

// 获取当前用户的SocketInfo，获取所有用户的SocketInfo
func TestRtdbConnect_GetSocketInfo(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取自己的Socket信息
	ownInfo, err := conn.GetOwnSocketInfo()
	if err != nil {
		t.Error("获取自己的SocketInfo失败", err)
		return
	}
	fmt.Println(ownInfo)

	// 设置Socket超时时间
	err = conn.SetSocketTimeout(ownInfo[0], ownInfo[0].Timeout)
	if err != nil {
		t.Error("设置timeout失败", err)
		return
	}

	// 获取全部Socket信息
	allInfos, err := conn.GetSocketInfos()
	if err != nil {
		t.Error("获取所有SocketInfo失败", err)
		return
	}
	fmt.Println(allInfos)
}

// IP黑名单
func TestRtdbConnect_BlackList(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 添加黑名单
	err = conn.AddIpBlackList("192.168.123.123", "255.255.255.0", "add 123")
	if err != nil {
		t.Error("添加黑名单失败：", err)
		return
	}

	// 修改黑名单
	err = conn.UpdateIpBlackList("192.168.123.123", "255.255.255.0", "192.168.123.123", "255.255.255.0", "update 123")
	if err != nil {
		t.Error("修改黑名单失败：", err)
		return
	}

	// 获取黑名单
	bLists, err := conn.GetIpBlackLists()
	if err != nil {
		t.Error("获取黑名单失败：", err)
		return
	}
	bOk := false
	for _, b := range bLists {
		if b.Desc == "update 123" {
			bOk = true
			break
		}
	}
	if !bOk {
		t.Error("修改黑名单失败")
		return
	}

	err = conn.DeleteIpBlackList("192.168.123.123", "255.255.255.0")
	if err != nil {
		t.Error("删除黑名单失败：", err)
		return
	}
}

// IP白名单
func TestRtdbConnect_WhiteList(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 添加白名单
	err = conn.AddIpWhiteList("192.168.123.120", "255.255.255.0", "add 120", PrivGroupRtdbSA)
	if err != nil {
		t.Error("添加白名单失败：", err)
		return
	}

	// 修改白名单
	err = conn.UpdateIpWhiteList("192.168.123.120", "255.255.255.0", "192.168.123.120", "255.255.255.0", "update 120", PrivGroupRtdbSA)
	if err != nil {
		t.Error("修改白名单失败：", err)
		return
	}

	// 获取白名单
	wLists, err := conn.GetIpWhiteLists()
	if err != nil {
		t.Error("获取白名单失败：", err)
		return
	}
	wOk := false
	for _, w := range wLists {
		if w.Desc == "update 120" {
			wOk = true
			break
		}
	}
	if !wOk {
		t.Error("修改白名单失败")
		return
	}

	err = conn.DeleteIpWhiteList("192.168.123.120", "255.255.255.0")
	if err != nil {
		t.Error("删除白名单失败：", err)
		return
	}
}

// 用户
func TestRtdbConnect_User(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 添加用户
	err = conn.AddUser("test111", "122333", PrivGroupRtdbSA)
	if err != nil {
		t.Error("添加用户失败: ", err)
		return
	}

	// 修改用户密码
	err = conn.UpdatePassword("test111", "123123")
	if err != nil {
		t.Error("修改密码失败: ", err)
		return
	}

	// 验证密码是否修改成功
	conn2, err := Login(Hostname, Port, "test111", "123123", RtdbPrecisionNano)
	if err != nil {
		t.Error("登录用户失败", err)
		return
	}
	defer func() { _ = conn2.Logout() }()

	// 修改自己的密码
	err = conn2.UpdateOwnPassword("123123", "122333")
	if err != nil {
		t.Error("修改自己的密码失败：", err)
		return
	}

	// 获取连接权限
	priv, err := conn2.GetPriv()
	if err != nil {
		t.Error("获取权限失败：", err)
		return
	}
	if *priv != PrivGroupRtdbSA {
		t.Error("验证权限失败")
		return
	}

	// 设置连接权限
	err = conn2.SetPriv("test111", PrivGroupRtdbRO)
	if err != nil {
		t.Error("设置权限失败：", err)
		return
	}

	// 锁定用户
	err = conn.LockUser("test111", OFF)
	if err != nil {
		t.Error("锁定用户失败：", err)
		return
	}

	// 用户列表
	users, err := conn.GetUsers()
	if err != nil {
		t.Error("获取用户列表失败：", err)
		return
	}
	uOk := false
	// 验证
	for _, u := range users {
		if u.User == "test111" {
			uOk = true
			break
		}
	}
	if !uOk {
		t.Error("用户列表中不存在test111")
		return
	}

	// 删除用户
	err = conn.DeleteUser("test111")
	if err != nil {
		t.Error("删除用户失败：", err)
	}
}

// 自定义类型
func TestRtdbConnect_NamedType(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建自定义类型
	err = conn.AddNamedType(
		"abc",
		"abc desc",
		RtdbDataTypeField{
			Name:   "A",
			Type:   RtdbTypeReal64,
			Length: 0,
			Desc:   "A desc",
		}, RtdbDataTypeField{
			Name:   "B",
			Type:   RtdbTypeReal64,
			Length: 0,
			Desc:   "B desc",
		}, RtdbDataTypeField{
			Name:   "C",
			Type:   RtdbTypeReal64,
			Length: 0,
			Desc:   "C desc",
		})
	if err != nil {
		t.Error("添加自定义类型失败")
		return
	}

	// 删除自定义类型
	defer func() {
		err := conn.DeleteNamedType("abc")
		if err != nil {
			t.Error("删除自定义类型失败")
			return
		}
	}()

	// 获取自定义类型
	types, err := conn.GetNamedTypes()
	if err != nil {
		t.Error("获取列表失败")
		return
	}
	fmt.Println(types)

	// 更新自定义类型
	desc := "up abc desc"
	err = conn.UpdateNamedType("abc", nil, &desc, map[string]string{"A": "A up", "B": "B up", "C": "C up"})
	if err != nil {
		t.Error("更新列表失败")
		return
	}

	// 获取自定义类型
	typ, err := conn.GetNamedType("abc")
	if err != nil {
		t.Error("获取列表失败")
		return
	}
	fmt.Println(typ)
}

// 时间
func TestRtdbConnect_Time(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取服务端主机时间
	hostTime, err := conn.ServerHostTime()
	if err != nil {
		t.Error("获取服务端时间失败：", err)
		return
	}
	fmt.Println(hostTime)

	// 时间段转字符串
	dStr, err := conn.DurationToString(time.Second * 60)
	if err != nil {
		t.Error("时间段转换失败：", err)
		return
	}
	if dStr != "1n" {
		t.Error("不为1n")
		return
	}

	// 字符串转时间段
	duration, err := conn.StringToDuration("1n")
	if err != nil {
		t.Error("时间段转换失败：", err)
		return
	}
	if duration != time.Minute {
		t.Error("时间段失败")
		return
	}

	// 字符串转时间戳
	ts, err := conn.StringToTime("2010-1-1 8:00:00")
	if err != nil {
		t.Error("字符串转时间戳失败：", err)
		return
	}
	fmt.Println(ts)
}

// 质量
func TestRtdbConnect_Quality(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取质量码对应的说明
	descs, err := conn.GetQualityDesc([]Quality{1})
	if err != nil {
		t.Error("获取质量码失败: ", err)
		return
	}
	fmt.Println(descs)
}

// 磁盘
func TestRtdbConnect_Disk(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 获取盘符
	letters, err := conn.GetDriveLetterList()
	if err != nil {
		t.Error("获取盘符列表失败")
		return
	}
	fmt.Println(letters)

	// 读取指定文件夹的所有目录项
	items, err := conn.GetDirItemList(letters[0])
	if err != nil {
		t.Error("获取目录项失败:", err)
		return
	}
	fmt.Println(items)

	// 创建目录
	err = conn.CreateDir(path.Join(letters[0], "hello"))
	if err != nil {
		t.Error("创建目录失败：", err)
		return
	}

	// 读取文件
	data, err := conn.ReadFile("/etc/hosts")
	if err != nil {
		t.Error("读取文件失败：", err)
		return
	}
	fmt.Println(string(data))
}

// 表
func TestRtdbConnect_Table(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable("ttt", "ttt desc")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	// 删除表
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 更新表名
	err = conn.UpdateTableName(table.ID, "ttt2")
	if err != nil {
		t.Error("更新表名失败：", err)
		return
	}
	time.Sleep(time.Second)

	// 更新表描述
	err = conn.UpdateTableDesc(table.ID, "ttt2 desc")
	if err != nil {
		t.Error("更新表描述失败：", err)
		return
	}
	time.Sleep(time.Second)

	// 获取表列表
	tables, err := conn.GetTables()
	if err != nil {
		t.Error("获取表列表失败：", err)
		return
	}
	fmt.Println(tables)

	// 获取表
	table2, err := conn.GetTable(table.ID)
	if err != nil {
		t.Error("获取表失败：", err)
		return
	}
	fmt.Println(table2)
}

// 标签点
func TestRtdbConnect_Point(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable("ppp", "ppp desc")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	// 删除表
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 添加点
	info := NewPointInfo("aaa", table.ID, ValueTypeInt32, PointBase, RtdbPrecisionMicro, "", "")
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败: ", err)
		return
	}

	// 删除点
	defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	// 获取单个点
	pInfo2, err := conn.GetPoint(pInfo.ID)
	if err != nil {
		t.Error("获取点失败: ", err)
		return
	}
	fmt.Println(pInfo2)

	// 修改点
	err = conn.UpdatePoint(pInfo.ID, map[PointInfoField]any{
		PointInfoFieldDesc: "point desc ???",
	})
	if err != nil {
		t.Error("修改点失败：", err)
		return
	}

	// 获取多个点
	pInfos, _, err := conn.GetPoints([]PointID{pInfo.ID})
	if err != nil {
		t.Error("获取点失败: ", err)
		return
	}
	fmt.Println(pInfos)

	// 移动点
	table2, err := conn.CreateTable("pp2", "pp2 desc")
	if err != nil {
		t.Error("创建表2失败：", err)
		return
	}
	defer func() { _ = conn.DeleteTable(table2.ID) }()
	err = conn.MovePoint(pInfo.ID, table2.Name)
	if err != nil {
		t.Error("移动点失败：", err)
		return
	}
	time.Sleep(time.Second)

	// 查找点
	ps, _, err := conn.FindPoints([]string{"pp2.aaa"})
	if err != nil {
		t.Error("查找点失败")
		return
	}
	fmt.Println(ps)

	// 搜索点
	count, ps2, _, err := conn.SearchPoint(0, 100, "", "pp2", "", "", "", "", "", -1, -1, 0, "", RtdbSortFlagDescend)
	if err != nil {
		t.Error("搜索点失败：", err)
		return
	}
	fmt.Println("点总数：", count)
	fmt.Println(ps2)
}

// 回收站 (标签点删除后会先进入回收站，从回收站清楚后才算是彻底删除)
func TestRtdbConnect_Recycler(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 清空回收站
	err = conn.ClearRecycler()
	if err != nil {
		t.Error("清空回收站失败：", err)
		return
	}

	// 创建表
	table, err := conn.CreateTable("ppp", "ppp desc")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	// 删除表
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 添加点
	info := NewPointInfo("aaa", table.ID, ValueTypeInt32, PointBase, RtdbPrecisionMicro, "", "")
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败: ", err)
		return
	}

	// 删除点, 注意，此时回收站中应该有一个点的
	err = conn.DeletePoint(pInfo.ID)
	if err != nil {
		t.Error("删除点失败: ", err)
		return
	}

	// 分批获取回收站中的点
	rCount, infos, errs, err := conn.GetRecycledPoints(0, 1024)
	if err != nil {
		t.Error("获取点失败：", err)
		return
	}
	if rCount != 1 {
		t.Error("回收站中点数量不为1")
		return
	}
	if errs[0] != nil {
		t.Error("获取点信息失败：", errs[0])
		return
	}
	fmt.Println(infos)

	// 恢复点到表
	err = conn.RecoverPoint(table.ID, infos[0].ID)
	if err != nil {
		t.Error("恢复点失败:", err)
		return
	}

	// 查找已恢复的点
	infos, errs, err = conn.FindPoints([]string{"ppp.aaa"})
	if err != nil {
		t.Error("查找点失败：", err)
		return
	}
	if errs[0] != nil {
		t.Error("获取点信息失败：", errs[0])
		return
	}
	fmt.Println(infos)

	// 删除点, 此时回收站中点个数应该为1
	err = conn.DeletePoint(infos[0].ID)
	if err != nil {
		t.Error("删除点失败：", err)
		return
	}

	// 在回收站中搜索点
	rCount, infos, errs, err = conn.SearchRecycledPoint(0, 1024, "", "", "", "", "", "", RtdbSortFlagDescend)
	if err != nil {
		t.Error("查找点失败：", err)
		return
	}
	if errs[0] != nil {
		t.Error("获取点信息失败：", errs[0])
		return
	}
	if rCount != 1 {
		t.Error("回收站中的点应为1")
		return
	}
	fmt.Println(infos[0])

	// 从回收站中清除点，此时点会被彻底删除
	err = conn.PurgePoint(infos[0].ID)
	if err != nil {
		t.Error("从回收站中清除点失败：", err)
		return
	}
}

// 获取某个数值类型对应的点数量
func TestRtdbConnect_GetPointCountFromValueType(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	count, err := conn.GetPointCountFromValueType(ValueTypeInt32)
	if err != nil {
		t.Error("获取int32类型对应的count失败:", err)
		return
	}
	fmt.Println(count)
}

// 点值(TVQ)写入， 读取 - Float64类型
func TestRtdbConnect_ReadWriteValue(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable("ttt2", "ttt desc")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	// 删除表
	defer func() { _ = conn.DeleteTable(table.ID) }()

	fmt.Println(table.ID)

	// 添加Float64类型的点
	desc := "温度数据"
	fmt.Printf("Desc content: %s\n", desc)
	fmt.Printf("Desc bytes: %v\n", []byte(desc))
	fmt.Printf("Desc hex: %x\n", []byte(desc))
	fmt.Printf("Desc length: %d bytes\n", len(desc))

	// UTF-8 "温度数据" 应该是: e6 b8 a9 e5 ba a6 e6 95 b0 e6 8d ae (12字节)
	// GBK "温度数据" 应该是: ce c2 b6 c8 ca fd be dd (8字节)

	info := NewPointInfo("aaa", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "", desc)
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败: ", err)
		return
	}

	// 删除点
	defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	serverTime, err := conn.ServerHostTime()
	if err != nil {
		t.Error("获取系统时间失败：", err)
	}
	fmt.Println("server time:", serverTime.Format(time.RFC3339))
	fmt.Println("client time:", time.Now().Format(time.RFC3339))

	// 写入数据
	n := 10
	for i := 0; i < n; i++ {
		// 单条时间序列，写单个TVQ
		value := 25.0 + float64(i)*0.5
		err := conn.WriteValue(pInfo, false, pInfo.NewNowTVQ(value, Quality(0)))
		if err != nil {
			t.Error("写入数据失败：", err)
			return
		}
		if i != n-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// 单条时间序列，写多个TVQ
	now := time.Now()
	tvqs := []TVQ{
		pInfo.NewTVQ(now.Add(-2*time.Second), 30.0, Quality(0)),
		pInfo.NewTVQ(now.Add(-1*time.Second), 31.0, Quality(0)),
		pInfo.NewTVQ(now, 32.0, Quality(0)),
	}
	errs, err := conn.WriteValues(pInfo, false, tvqs)
	if err != nil {
		t.Error("WriteValues写入数据失败：", err)
		return
	}
	for i, e := range errs {
		if e != nil {
			t.Errorf("WriteValues第%d个数据写入失败: %v", i, e)
		}
	}

	// 写断面 (可以是单条序列，也可以是多条序列，更灵活)
	ptvqs := []PTVQ{
		NewPTVQ(pInfo, pInfo.NewTVQ(now.Add(-2*time.Second), 33.0, Quality(0))),
		NewPTVQ(pInfo, pInfo.NewTVQ(now.Add(-1*time.Second), 34.0, Quality(0))),
	}
	errs, err = conn.WriteSection(false, ptvqs)
	if err != nil {
		t.Error("WriteSection写入数据失败：", err)
		return
	}
	for i, e := range errs {
		if e != nil {
			t.Errorf("WriteSection第%d个数据写入失败: %v", i, e)
		}
	}

	// 刷新历史归档缓存，确保写入的历史数据落盘后再读取
	{
		_, err := conn.FlushArchivedValues(pInfo)
		if err != nil {
			t.Error("flush archived values err: ", err)
			return
		}
	}

	// 读取最新的实时值
	{
		ptvq, err := conn.ReadLast(pInfo)
		if err != nil {
			t.Error("read last err：", err)
		} else {
			fmt.Printf("ReadLast: Time=%v, Value=%v, Quality=%d\n",
				ptvq.TVQ.Timestamp.Format(time.RFC3339Nano),
				ptvq.TVQ.Value.FloatValue, ptvq.TVQ.Quality)
		}
	}

	// 读取单个TVQ
	{
		ptvq, err := conn.ReadValue(pInfo, RtdbHisModeExactOrPrev, time.Now())
		if err != nil {
			t.Error("read value err：", err)
		} else {
			fmt.Printf("ReadValue: Time=%v, Value=%v, Quality=%d\n",
				ptvq.TVQ.Timestamp.Format(time.RFC3339Nano),
				ptvq.TVQ.Value.FloatValue, ptvq.TVQ.Quality)
		}
	}

	// 读取Range
	{
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		ptvqs, err := conn.ReadRange(pInfo, startTime, endTime)
		if err != nil {
			t.Error("read range err：", err)
		} else {
			fmt.Printf("ReadRange: 读取到%d条数据\n", len(ptvqs.TVQs))
			for i, tvq := range ptvqs.TVQs {
				if i >= 5 {
					fmt.Println("  ... (仅显示前5条)")
					break
				}
				fmt.Printf("  [%d] Time=%v, Value=%v, Quality=%d\n",
					i, tvq.Timestamp.Format(time.RFC3339Nano),
					tvq.Value.FloatValue, tvq.Quality)
			}
		}
	}

	// 读取Plot （用于绘图的TVQ）
	{
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		ptvqs, err := conn.ReadPlot(pInfo, 100, startTime, endTime)
		if err != nil {
			t.Error("read plot err：", err)
		} else {
			fmt.Printf("ReadPlot: 读取到%d条数据(最大100条)\n", len(ptvqs.TVQs))
		}
	}

	// 读取差值（按照指定时间戳）
	{
		targetTimes := []time.Time{
			time.Now().Add(-30 * time.Second),
			time.Now().Add(-20 * time.Second),
			time.Now().Add(-10 * time.Second),
		}
		ptvqs, err := conn.ReadTimed(pInfo, targetTimes)
		if err != nil {
			t.Error("read timed err：", err)
		} else {
			fmt.Printf("ReadTimed: 读取到%d条数据\n", len(ptvqs.TVQs))
			for i, ptvq := range ptvqs.TVQs {
				fmt.Printf("  [%d] TargetTime=%v, ActualTime=%v, Value=%v, Quality=%d\n",
					i, targetTimes[i].Format(time.RFC3339Nano),
					ptvq.Timestamp.Format(time.RFC3339Nano),
					ptvq.Value.FloatValue, ptvq.Quality)
			}
		}
	}

	// 读取差值（start、end之间等分成count份）
	{
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		ptvqs, err := conn.ReadInterpo(pInfo, 10, startTime, endTime, "")
		if err != nil {
			t.Error("read interpo err：", err)
		} else {
			fmt.Printf("ReadInterpo: 读取到%d条数据\n", len(ptvqs.TVQs))
			for i, ptvq := range ptvqs.TVQs {
				fmt.Printf("  [%d] Time=%v, Value=%v, Quality=%d\n",
					i, ptvq.Timestamp.Format(time.RFC3339Nano),
					ptvq.Value.FloatValue, ptvq.Quality)
			}
		}
	}

	// 读取差值 (从start开始，每隔duration时间间隔读取一个差值，最多读取count个)
	{
		startTime := time.Now().Add(-1 * time.Hour)
		ptvqs, err := conn.ReadInterval(pInfo, "", startTime, 10*time.Minute, 6)
		if err != nil {
			t.Error("read interval err：", err)
		} else {
			fmt.Printf("ReadInterval: 读取到%d条数据\n", len(ptvqs.TVQs))
			for i, ptvq := range ptvqs.TVQs {
				fmt.Printf("  [%d] Time=%v, Value=%v, Quality=%d\n",
					i, ptvq.Timestamp.Format(time.RFC3339Nano),
					ptvq.Value.FloatValue, ptvq.Quality)
			}
		}
	}

	// 读取断面
	{
		ptvqs, errs, err := conn.ReadSection([]*PointInfo{pInfo}, RtdbHisModeExactOrPrev, time.Now())
		if err != nil {
			t.Error("read section err：", err)
		} else {
			for i, e := range errs {
				if e != nil {
					t.Errorf("ReadSection第%d个点读取失败: %v", i, e)
				}
			}
			fmt.Printf("ReadSection: 读取到%d个点的数据\n", len(ptvqs))
			for i, ptvq := range ptvqs {
				fmt.Printf("  [%d] PointID=%d, Time=%v, Value=%v, Quality=%d\n",
					i, ptvq.PointInfo.ID, ptvq.TVQ.Timestamp.Format(time.RFC3339Nano),
					ptvq.TVQ.Value.FloatValue, ptvq.TVQ.Quality)
			}
		}
	}

	// 读取统计值
	{
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		summary, err := conn.ReadSummary(pInfo, "", startTime, endTime)
		if err != nil {
			t.Error("read summary err：", err)
		} else {
			fmt.Printf("ReadSummary:\n")
			fmt.Printf("  Count=%d\n", summary.Count)
			fmt.Printf("  Min=%v (Time=%v)\n", summary.MinValue, summary.MinTime)
			fmt.Printf("  Max=%v (Time=%v)\n", summary.MaxValue, summary.MaxTime)
			fmt.Printf("  PowerAvg=%v\n", summary.PowerAvg)
		}
	}

	// 读取等间隔统计值
	{
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()
		summaryList, errs, err := conn.ReadBatchesSummary(pInfo, "", 15*time.Minute, startTime, endTime)
		if err != nil {
			t.Error("read batches summary err：", err)
		} else {
			for i, e := range errs {
				if e != nil {
					t.Errorf("ReadBatchesSummary第%d个统计段读取失败: %v", i, e)
				}
			}
			fmt.Printf("ReadBatchesSummary: 读取到%d个时间段的统计值\n", len(summaryList))
			for i, summary := range summaryList {
				fmt.Printf("  [%d] Count=%d, Min=%v, Max=%v, Avg=%v\n",
					i, summary.Count, summary.MinValue, summary.MaxValue, summary.PowerAvg)
			}
		}
	}

	// 批量删除点 - 删除历史数据范围
	{
		// 删除历史数据范围（1小时前到30分钟前）
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now().Add(-30 * time.Minute)
		count, err := conn.RemoveRangeValues(pInfo, startTime, endTime)
		if err != nil {
			t.Error("RemoveRangeValues删除失败:", err)
		} else {
			fmt.Printf("RemoveRangeValues成功: 删除了%d条历史数据\n", count)
		}
	}

	// 刷新数据页缓存 (就是把内存中的快照数据，刷新到历史中)
	{
		count, err := conn.FlushArchivedValues(pInfo)
		if err != nil {
			t.Error("flush archived values err: ", err)
		} else {
			fmt.Printf("FlushArchivedValues: 刷新了%d条数据\n", count)
		}
	}
}

// 修改单个点值测试 - 专门测试UpdateValue功能
func TestRtdbConnect_UpdateSingleValue(t *testing.T) {
	prefix := fmt.Sprintf("p9_update_%d_", time.Now().UnixMilli())
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable(prefix+"table", "update test table")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 添加Float64类型的点
	info := NewPointInfo(prefix+"point", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "°C", "更新测试点")
	info.SetLimit(0, 200, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败: ", err)
		return
	}
	defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	// 测试场景1: 写入单个历史值，然后更新它
	t.Run("UpdateSingleValue", func(t *testing.T) {
		// targetTime 作为历史时间点（早于快照时间，触发历史归档路径）
		targetTime := time.Now().Add(-5 * time.Second)
		originalValue := 50.0
		updatedValue := 75.5

		// 先清除 targetTime 附近旧历史数据，防止污染
		_, _ = conn.RemoveRangeValues(pInfo, targetTime.Add(-time.Second), targetTime.Add(time.Second))

		// 建立快照（快照时间 > targetTime），后续写入 targetTime 将触发 EarlierThanSnapshot
		err := conn.WriteValue(pInfo, false, pInfo.NewTVQ(time.Now(), 0.0, Quality(0)))
		if err != nil {
			t.Error("建立初始快照失败:", err)
			return
		}

		// 写入历史值（时间早于快照，走历史归档路径）
		err = conn.WriteValue(pInfo, false, pInfo.NewTVQ(targetTime, originalValue, Quality(0)))
		if err != nil {
			t.Error("写入原始历史值失败:", err)
			return
		}

		// 刷新历史归档缓存，确保数据落盘
		_, err = conn.FlushArchivedValues(pInfo)
		if err != nil {
			t.Error("刷新历史归档失败:", err)
			return
		}

		// 读取确认历史写入成功
		ptvq, err := conn.ReadValue(pInfo, RtdbHisModeExactOrPrev, targetTime)
		if err != nil {
			t.Error("读取原始历史值失败:", err)
			return
		}
		if ptvq.TVQ.Value.FloatValue != originalValue {
			t.Errorf("写入后读取值不匹配: 期望%v, 实际%v", originalValue, ptvq.TVQ.Value.FloatValue)
			return
		}

		// 更新历史值（修改Value，时间戳保持不变）
		err = conn.UpdateValue(pInfo, NewTvqFloat64(targetTime, updatedValue, Quality(0)))
		if err != nil {
			t.Error("UpdateValue失败:", err)
			return
		}

		// 验证更新成功
		ptvq, err = conn.ReadValue(pInfo, RtdbHisModeExactOrPrev, targetTime)
		if err != nil {
			t.Error("读取更新后的历史值失败:", err)
			return
		}
		if ptvq.TVQ.Value.FloatValue != updatedValue {
			t.Errorf("更新后的值不匹配: 期望%v, 实际%v", updatedValue, ptvq.TVQ.Value.FloatValue)
			return
		}

		fmt.Printf("UpdateSingleValue成功: 历史值已从%v更新为%v @ %v\n",
			originalValue, updatedValue, targetTime.Format(time.RFC3339))
	})

	// 测试场景2: 批量写入历史值，然后更新其中一个（非快照）
	t.Run("UpdateOneOfBatchValues", func(t *testing.T) {
		// 使用当前时间之前的时间，确保早于快照时间
		baseTime := time.Now().Add(-20 * time.Second)

		// 先清除该时间段的旧历史数据，防止污染
		_, _ = conn.RemoveRangeValues(pInfo, baseTime.Add(-time.Second), baseTime.Add(6*time.Second))

		// 建立快照（快照时间 > baseTime）
		err := conn.WriteValue(pInfo, false, pInfo.NewTVQ(time.Now(), 0.0, Quality(0)))
		if err != nil {
			t.Error("建立初始快照失败:", err)
			return
		}

		tvqs := []TVQ{
			pInfo.NewTVQ(baseTime, 10.0, Quality(0)),
			pInfo.NewTVQ(baseTime.Add(2*time.Second), 20.0, Quality(0)),
			pInfo.NewTVQ(baseTime.Add(4*time.Second), 30.0, Quality(0)),
		}

		// 批量写入历史值
		errs, err := conn.WriteValues(pInfo, false, tvqs)
		if err != nil {
			t.Error("批量写入失败:", err)
			return
		}
		for i, e := range errs {
			if e != nil {
				t.Errorf("第%d个历史值写入失败: %v", i, e)
				return
			}
		}

		// 刷新历史归档缓存
		_, err = conn.FlushArchivedValues(pInfo)
		if err != nil {
			t.Error("刷新历史归档失败:", err)
			return
		}

		// 更新中间的历史值（不是最新的）
		updateTime := baseTime.Add(2 * time.Second)
		newValue := 25.0
		err = conn.UpdateValue(pInfo, NewTvqFloat64(updateTime, newValue, Quality(0)))
		if err != nil {
			t.Error("更新中间历史值失败:", err)
			return
		}

		// 读取时间范围验证
		ptvqs, err := conn.ReadRange(pInfo, baseTime, baseTime.Add(5*time.Second))
		if err != nil {
			t.Error("读取历史范围失败:", err)
			return
		}

		// 验证更新的历史值正确
		updated := false
		for _, tvq := range ptvqs.TVQs {
			if tvq.Timestamp.Equal(updateTime) {
				if tvq.Value.FloatValue != newValue {
					t.Errorf("更新的历史值不匹配: 期望%v, 实际%v", newValue, tvq.Value.FloatValue)
					return
				}
				updated = true
			}
		}
		if !updated {
			t.Errorf("未找到更新时间点的历史记录: %v", updateTime)
			return
		}

		fmt.Printf("UpdateOneOfBatchValues成功，当前历史数据:\n")
		for i, tvq := range ptvqs.TVQs {
			marker := ""
			if tvq.Timestamp.Equal(updateTime) {
				marker = " <- 已更新"
			}
			fmt.Printf("  [%d] Time=%v, Value=%v%s\n",
				i, tvq.Timestamp.Format(time.RFC3339), tvq.Value.FloatValue, marker)
		}
	})

	// 测试场景3: 验证快照值不允许被 UpdateValue 修改（预期失败）
	// 快照是最新实时值，不在历史归档中；UpdateValue 只能修改历史归档
	t.Run("UpdateLatestValue_ExpectedFail", func(t *testing.T) {
		// 读出当前快照（前面场景已建立）
		ptvq, err := conn.ReadLast(pInfo)
		if err != nil {
			t.Error("ReadLast失败:", err)
			return
		}
		snapshotTime := ptvq.TVQ.Timestamp
		snapshotValue := ptvq.TVQ.Value.FloatValue
		newValue := snapshotValue + 50.0

		// 尝试用 UpdateValue 修改快照时间戳对应的历史值（预期失败：快照不在历史归档中）
		err = conn.UpdateValue(pInfo, NewTvqFloat64(snapshotTime, newValue, Quality(0)))
		if err != nil {
			// 这是预期的行为
			fmt.Printf("预期内的失败: 快照值不允许被UpdateValue修改 - %v\n", err)
		} else {
			t.Error("预期UpdateValue应该失败，但实际成功了")
			return
		}

		// 验证快照值未被修改
		ptvq, err = conn.ReadLast(pInfo)
		if err != nil {
			t.Error("ReadLast失败:", err)
			return
		}
		if ptvq.TVQ.Value.FloatValue == snapshotValue {
			fmt.Printf("验证成功: 快照值保持原值%v，未被修改\n", snapshotValue)
		} else {
			t.Errorf("验证失败: 期望%v, 实际%v", snapshotValue, ptvq.TVQ.Value.FloatValue)
		}
	})
}

// 快照订阅测试
func TestRtdbConnect_SubscribeSnapshot(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable("ttt3", "ttt desc")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	// 删除表
	defer func() { _ = conn.DeleteTable(table.ID) }()

	fmt.Println(table.ID)

	// 添加Float64类型的点
	desc := "温度数据"
	fmt.Printf("Desc content: %s\n", desc)
	fmt.Printf("Desc bytes: %v\n", []byte(desc))
	fmt.Printf("Desc hex: %x\n", []byte(desc))
	fmt.Printf("Desc length: %d bytes\n", len(desc))

	// UTF-8 "温度数据" 应该是: e6 b8 a9 e5 ba a6 e6 95 b0 e6 8d ae (12字节)
	// GBK "温度数据" 应该是: ce c2 b6 c8 ca fd be dd (8字节)

	info := NewPointInfo("aaa", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "", desc)
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败: ", err)
		return
	}

	// 删除点
	defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	serverTime, err := conn.ServerHostTime()
	if err != nil {
		t.Error("获取系统时间失败：", err)
	}
	fmt.Println("server time:", serverTime.Format(time.RFC3339))
	fmt.Println("client time:", time.Now().Format(time.RFC3339))

	snapChan, errs, err := conn.SubscribeSnapshots([]*PointInfo{pInfo})
	if err != nil {
		t.Error("订阅快照失败:", err)
		return
	}
	if len(errs) > 0 && errs[0] != nil {
		t.Error("订阅快照返回错误:", errs[0])
		return
	}
	closeChan := make(chan struct{})
	defer close(closeChan)
	defer func() {
		err := conn.CancelSubscribeSnapshots()
		if err != nil {
			t.Error("关闭订阅失败：", err)
		}
	}()
	go func() {
		for {
			select {
			case snap, ok := <-snapChan:
				if !ok {
					fmt.Println("快照订阅channel已关闭")
					return
				}
				fmt.Println(snap)
			case <-closeChan:
				return
			}
		}
	}()

	// 写入数据
	n := 10
	for i := 0; i < n; i++ {
		// 单条时间序列，写单个TVQ
		value := 25.0 + float64(i)*0.5
		err := conn.WriteValue(pInfo, false, pInfo.NewNowTVQ(value, Quality(0)))
		if err != nil {
			t.Error("写入数据失败：", err)
			return
		}
		if i != n-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// 测试订阅标签点属性更新
func TestRtdbConnect_SubscribeTags(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 订阅标签点属性更新
	tagsChan, err := conn.SubscribeTags()
	if err != nil {
		t.Error("订阅标签点属性失败:", err)
		return
	}

	// 确保取消订阅
	defer func() {
		err := conn.CancelSubscribeTags()
		if err != nil {
			t.Log("取消订阅标签点属性失败:", err)
		}
	}()

	// 启动goroutine接收订阅数据
	done := make(chan struct{})
	defer close(done)
	go func() {
		for {
			select {
			case tagInfo, ok := <-tagsChan:
				if !ok {
					fmt.Println("标签点属性订阅channel已关闭")
					return
				}
				fmt.Printf("收到标签点属性更新: Name=%s, EventType=%d, Ids=%v\n", tagInfo.Name, tagInfo.EventType, tagInfo.Ids)
			case <-done:
				return
			}
		}
	}()

	// 创建表和点来触发标签点属性变更
	table, err := conn.CreateTable("sub_tags_test", "订阅标签点属性测试")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 添加点
	info := NewPointInfo("test_point", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "", "测试点")
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败:", err)
		return
	}
	defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	// 修改点属性以触发订阅通知
	time.Sleep(100 * time.Millisecond)
	err = conn.UpdatePoint(pInfo.ID, map[PointInfoField]any{
		PointInfoFieldDesc: "修改后的描述",
	})
	if err != nil {
		t.Error("修改点属性失败:", err)
		return
	}

	// 等待一段时间接收订阅数据
	time.Sleep(500 * time.Millisecond)
	fmt.Println("标签点属性订阅测试完成")
}

// 测试Delta快照订阅
func TestRtdbConnect_SubscribeDeltaSnapshots(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable("delta_snap_test", "Delta快照订阅测试")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 添加Float64类型的点
	info := NewPointInfo("delta_point", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "", "Delta测试点")
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败:", err)
		return
	}
	defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	// 订阅Delta快照，设置deltaValue为5.0，deltaState为0
	deltaValues := []float64{5.0}
	deltaStates := []int64{0}
	snapChan, errs, err := conn.SubscribeDeltaSnapshots([]*PointInfo{pInfo}, deltaValues, deltaStates)
	if err != nil {
		t.Error("订阅Delta快照失败:", err)
		return
	}
	if len(errs) > 0 && errs[0] != nil {
		t.Error("订阅Delta快照返回错误:", errs[0])
		return
	}

	// 确保取消订阅
	defer func() {
		err := conn.CancelSubscribeSnapshots()
		if err != nil {
			t.Log("取消订阅快照失败:", err)
		}
	}()

	// 启动goroutine接收订阅数据
	done := make(chan struct{})
	defer close(done)
	receivedCount := 0
	go func() {
		for {
			select {
			case snap, ok := <-snapChan:
				if !ok {
					fmt.Println("Delta快照订阅channel已关闭")
					return
				}
				receivedCount++
				fmt.Printf("收到Delta快照更新 [%d]: Name=%s, PTVQs数量=%d\n",
					receivedCount, snap.Name, len(snap.PTVQs))
			case <-done:
				return
			}
		}
	}()

	// 写入数据，变化小于5.0，应该不会触发订阅
	fmt.Println("写入变化小于5.0的数据，应该不会触发订阅...")
	for i := 0; i < 3; i++ {
		value := 25.0 + float64(i)*1.0 // 每次变化1.0，小于5.0
		err := conn.WriteValue(pInfo, false, pInfo.NewNowTVQ(value, Quality(0)))
		if err != nil {
			t.Error("写入数据失败：", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 写入数据，变化大于5.0，应该触发订阅
	fmt.Println("写入变化大于5.0的数据，应该触发订阅...")
	for i := 0; i < 3; i++ {
		value := 30.0 + float64(i)*6.0 // 每次变化6.0，大于5.0
		err := conn.WriteValue(pInfo, false, pInfo.NewNowTVQ(value, Quality(0)))
		if err != nil {
			t.Error("写入数据失败：", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 等待一段时间接收订阅数据
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Delta快照订阅测试完成，共收到%d条订阅通知\n", receivedCount)
}

// 测试修改快照订阅
func TestRtdbConnect_ChangeSubscribeSnapshots(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable("change_snap_test", "修改快照订阅测试")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 添加两个点
	info1 := NewPointInfo("point1", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "", "测试点1")
	info1.SetLimit(-100, 100, 0)
	pInfo1, err := conn.AddPoint(info1)
	if err != nil {
		t.Error("添加点1失败:", err)
		return
	}
	defer func() { _ = conn.DeletePoint(pInfo1.ID) }()

	info2 := NewPointInfo("point2", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "", "测试点2")
	info2.SetLimit(-100, 100, 0)
	pInfo2, err := conn.AddPoint(info2)
	if err != nil {
		t.Error("添加点2失败:", err)
		return
	}
	defer func() { _ = conn.DeletePoint(pInfo2.ID) }()

	// 先订阅第一个点
	snapChan, errs, err := conn.SubscribeSnapshots([]*PointInfo{pInfo1})
	if err != nil {
		t.Error("订阅快照失败:", err)
		return
	}
	if len(errs) > 0 && errs[0] != nil {
		t.Error("订阅快照返回错误:", errs[0])
		return
	}

	// 确保取消订阅
	defer func() {
		err := conn.CancelSubscribeSnapshots()
		if err != nil {
			t.Log("取消订阅快照失败:", err)
		}
	}()

	// 启动goroutine接收订阅数据
	done := make(chan struct{})
	defer close(done)
	go func() {
		for {
			select {
			case snap, ok := <-snapChan:
				if !ok {
					fmt.Println("快照订阅channel已关闭")
					return
				}
				fmt.Printf("收到快照更新: Name=%s, PTVQs数量=%d\n", snap.Name, len(snap.PTVQs))
			case <-done:
				return
			}
		}
	}()

	// 写入第一个点的数据
	fmt.Println("写入point1的数据...")
	err = conn.WriteValue(pInfo1, false, pInfo1.NewNowTVQ(10.0, Quality(0)))
	if err != nil {
		t.Error("写入数据失败：", err)
		return
	}
	time.Sleep(100 * time.Millisecond)

	// 修改订阅，添加第二个点
	fmt.Println("修改订阅，添加point2...")
	errs, err = conn.ChangeSubscribeSnapshots(
		[]*PointInfo{pInfo2},
		nil,
		nil,
		[]RtdbSubscribeChangeType{RtdbSubscribeChangeTypeAdd},
	)
	if err != nil {
		t.Error("修改订阅失败:", err)
		return
	}
	for i, e := range errs {
		if e != nil {
			t.Errorf("修改订阅第%d个点失败: %v", i, e)
		}
	}

	// 写入第二个点的数据
	fmt.Println("写入point2的数据...")
	err = conn.WriteValue(pInfo2, false, pInfo2.NewNowTVQ(20.0, Quality(0)))
	if err != nil {
		t.Error("写入数据失败：", err)
		return
	}
	time.Sleep(100 * time.Millisecond)

	// 修改订阅，移除第一个点
	fmt.Println("修改订阅，移除point1...")
	errs, err = conn.ChangeSubscribeSnapshots(
		[]*PointInfo{pInfo1},
		nil,
		nil,
		[]RtdbSubscribeChangeType{RtdbSubscribeChangeTypeRemove},
	)
	if err != nil {
		t.Error("修改订阅失败:", err)
		return
	}

	// 等待一段时间
	time.Sleep(300 * time.Millisecond)
	fmt.Println("修改快照订阅测试完成")
}

// 测试数据流订阅（Datagram）
func TestRtdbConnect_Datagram(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建数据流订阅
	// 注意：这里使用本地地址作为示例，实际使用时需要根据RTDB服务器的配置调整
	remoteHost := "127.0.0.1"
	datagramHandle, err := conn.CreateDatagram(8000, remoteHost)
	if err != nil {
		t.Error("创建数据流订阅失败:", err)
		return
	}

	// 确保移除数据流订阅
	defer func() {
		err := conn.RemoveDatagram(datagramHandle)
		if err != nil {
			t.Log("移除数据流订阅失败:", err)
		}
	}()

	fmt.Printf("数据流订阅创建成功: Handle=%v\n", datagramHandle)

	// 尝试接收数据（设置较短的超时时间，因为可能没有数据）
	// 注意：这个测试主要是验证接口调用，实际数据接收需要服务器配合
	fmt.Println("尝试接收数据流（设置1秒超时）...")
	data, err := conn.RecvDatagram(datagramHandle, 1024, remoteHost, 1)
	if err != nil {
		// 超时或没有数据是预期的行为
		fmt.Printf("接收数据流结果: %v\n", err)
	} else {
		fmt.Printf("接收到数据流: %v\n", data)
	}

	fmt.Println("数据流订阅测试完成")
}

// 测试API调用连接事件订阅
func TestRtdbConnect_SubscribeConnectEvents(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password, RtdbPrecisionNano)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	table, err := conn.CreateTable("connect_event_test", "API调用事件订阅测试")
	if err != nil {
		t.Error("创建表失败：", err)
		return
	}
	defer func() { _ = conn.DeleteTable(table.ID) }()

	// 添加点
	info := NewPointInfo("event_point", table.ID, ValueTypeFloat64, PointBase, RtdbPrecisionNano, "", "事件测试点")
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败:", err)
		return
	}
	defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	// 订阅API调用连接事件
	eventChan, err := conn.SubscribeConnectEvents()
	if err != nil {
		t.Error("订阅API调用事件失败:", err)
		return
	}

	// 确保取消订阅
	defer func() {
		err := conn.CancelSubscribeConnectEvents()
		if err != nil {
			t.Log("取消订阅API调用事件失败:", err)
		}
	}()

	// 启动goroutine接收订阅数据
	done := make(chan struct{})
	defer close(done)
	receivedCount := 0
	go func() {
		for {
			select {
			case eventInfo, ok := <-eventChan:
				if !ok {
					fmt.Println("API调用事件订阅channel已关闭")
					return
				}
				receivedCount++
				fmt.Printf("收到API调用事件 [%d]: EventType=%d, Handle=%d, Events数量=%d, PreCalls数量=%d, PostCalls数量=%d\n",
					receivedCount, eventInfo.EventType, eventInfo.Handle, len(eventInfo.Events),
					len(eventInfo.PreCalls), len(eventInfo.PostCalls))
				for _, ev := range eventInfo.Events {
					fmt.Printf("  API调用详情: msg_id=%d, elapsed=%.2fms, ret_val=%d, client_process=%d, client_thread=%d\n",
						ev.MsgID, ev.Elapsed, ev.RetVal, ev.ClientProcessID, ev.ClientThreadID)
				}
			case <-done:
				return
			}
		}
	}()

	// 在主连接上执行API调用，触发订阅事件
	fmt.Println("执行API调用以触发事件...")
	for i := 0; i < 5; i++ {
		value := 20.0 + float64(i)*2.0
		err := conn.WriteValue(pInfo, false, pInfo.NewNowTVQ(value, Quality(0)))
		if err != nil {
			t.Error("写入数据失败：", err)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	// 再执行一个读取操作
	_, err = conn.ReadLast(pInfo)
	if err != nil {
		t.Error("读取实时值失败：", err)
		return
	}

	// 等待一段时间接收订阅数据
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("API调用事件订阅测试完成，共收到%d条订阅通知\n", receivedCount)
}
