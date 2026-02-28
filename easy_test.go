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
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	// 登出
	defer func() { _ = conn.Logout() }()

	fmt.Println(conn.SyncInfos, conn.StringBlobMaxLen)
}

// 获取客户端版本
func TestRtdbConnect_GetClientVersion(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn2, err := Login(Hostname, Port, "test111", "123123")
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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
	conn, err := Login(Hostname, Port, Username, Password)
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

// 点值(TVQ)写入
func TestRtdbConnect_Value(t *testing.T) {
	prefix := "p9_"
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	// 创建表
	// table, err := conn.CreateTable(prefix+"ppp", "ppp desc")
	// if err != nil {
	//	t.Error("创建表失败：", err)
	//	return
	// }
	// 删除表
	// defer func() { _ = conn.DeleteTable(table.ID) }()

	time.Sleep(time.Second)

	// 添加点
	info := NewPointInfo(prefix+"xxx", 1, ValueTypeCoor, PointBase, RtdbPrecisionNano, "", "")
	info.SetLimit(-100, 100, 0)
	pInfo, err := conn.AddPoint(info)
	if err != nil {
		t.Error("添加点失败: ", err)
		return
	}

	// 删除点
	// defer func() { _ = conn.DeletePoint(pInfo.ID) }()

	time.Sleep(time.Second)

	serverTime, err := conn.ServerHostTime()
	if err != nil {
		t.Error("获取系统时间失败：", err)
	}
	fmt.Println("server time:", serverTime.Format(time.RFC3339))
	fmt.Println("client time:", time.Now().Format(time.RFC3339))

	// 写入数据
	n := 100
	for i := 0; i < n; i++ {
		err := conn.WriteValue(pInfo, true, pInfo.NewNowTVQ(Coordinates{X: 1.0, Y: 2.0}, Quality(0)))
		if err != nil {
			t.Error("写入数据失败：", err)
			return
		}
		if i != n-1 {
			time.Sleep(1 * time.Second)
		}
	}
}

func TestRtdbConnect_ReadLast(t *testing.T) {
	conn, err := Login(Hostname, Port, Username, Password)
	if err != nil {
		t.Fatal("登录用户失败", err)
	}
	defer func() { _ = conn.Logout() }()

	infos, errs, err := conn.FindPoints([]string{"p7_ppp.p8_xxx"})
	if err != nil {
		t.Error("find point err：", err)
	}
	for _, err := range errs {
		if err != nil {
			t.Error("find point err2: ", err)
		}
	}
	info := infos[0]
	fmt.Println(info)

	ptvq, err := conn.ReadLast(info)
	if err != nil {
		t.Error("find point err：", err)
	}
	fmt.Println(ptvq)
}
