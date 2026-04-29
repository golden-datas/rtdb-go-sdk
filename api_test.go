package rtdb_api

import (
	"fmt"
	"testing"
	"time"
)

const Hostname = "159.75.187.68"
const Port = 6327
const Username = "sa"
const Password = "golden"

// connectAndLogin 建立连接并以管理员身份登录
func connectAndLogin(t *testing.T) ConnectHandle {
	t.Helper()
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if !RteIsOk(err) {
		t.Fatalf("连接失败: %v", err)
	}
	_, err = RawRtdbLoginWarp(handle, Username, Password)
	if !RteIsOk(err) {
		t.Fatalf("登录失败: %v", err)
	}
	return handle
}

// disconnect 断开连接
func disconnect(t *testing.T, handle ConnectHandle) {
	t.Helper()
	err := RawRtdbDisconnectWarp(handle)
	if !RteIsOk(err) {
		t.Logf("断开连接返回: %v", err)
	}
}

func TestNULL(t *testing.T) {}

// ==================== 01. 连接与系统参数 ====================

// TC-APIVER-01 正常获取版本号
func TestRawRtdbGetApiVersionWarp(t *testing.T) {
	apiVersion, err := RawRtdbGetApiVersionWarp()
	if !RteIsOk(err) {
		t.Error("获取版本号失败:", err)
	}
	fmt.Println("库版本号:", apiVersion)
}

// TC-OPTION-01 设置自动重连选项
func TestRawRtdbSetOptionWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbSetOptionWarp(RtdbApiOptionAutoReconn, 1)
	if !RteIsOk(err) {
		t.Error("设置自动重连失败:", err)
	}
}

// TC-CONN-01 正常连接数据库
func TestRawRtdbConnectWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if !RteIsOk(err) {
		t.Fatalf("连接失败: %v", err)
	}
	if handle <= 0 {
		t.Error("连接句柄无效")
	}
	fmt.Println("连接句柄:", handle)
	err = RawRtdbDisconnectWarp(handle)
	if !RteIsOk(err) {
		t.Logf("断开连接: %v", err)
	}
}

// TC-CONN-02 连接不存在的主机
func TestRawRtdbConnectWarp_UnreachableHost(t *testing.T) {
	handle, err := RawRtdbConnectWarp("192.0.2.1", 9000)
	if RteIsOk(err) {
		t.Error("连接不可达主机应失败")
	}
	if handle > 0 {
		RawRtdbDisconnectWarp(handle)
	}
	fmt.Println("不可达主机连接结果:", err)
}

// TC-LOGIN-01 管理员正常登录
func TestRawRtdbLoginWarp(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if !RteIsOk(err) {
		t.Fatalf("连接失败: %v", err)
	}
	defer RawRtdbDisconnectWarp(handle)

	priv, err := RawRtdbLoginWarp(handle, Username, Password)
	if !RteIsOk(err) {
		t.Fatalf("登录失败: %v", err)
	}
	fmt.Println("登录权限:", priv)
}

// TC-LOGIN-02 错误密码登录
func TestRawRtdbLoginWarp_WrongPassword(t *testing.T) {
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if !RteIsOk(err) {
		t.Fatalf("连接失败: %v", err)
	}
	defer RawRtdbDisconnectWarp(handle)

	_, err = RawRtdbLoginWarp(handle, Username, "wrong_pwd")
	if RteIsOk(err) {
		t.Error("错误密码应登录失败")
	}
	fmt.Println("错误密码登录结果:", err)
}

// TC-DISC-01 正常断开连接
func TestRawRtdbDisconnectWarp(t *testing.T) {
	handle := connectAndLogin(t)
	err := RawRtdbDisconnectWarp(handle)
	if !RteIsOk(err) {
		t.Error("断开连接失败:", err)
	}
}

// TC-SUBCONN-01 正常创建连接事件订阅
func TestRawRtdbSubscribeConnectExWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	param, err := RawRtdbSubscribeConnectExWarp(handle, RtdbSubscribeOptionAutoConn, "test_sub")
	if !RteIsOk(err) {
		t.Error("创建连接事件订阅失败:", err)
	}
	if param == nil {
		t.Error("订阅参数不应为空")
	}
	fmt.Println("订阅参数:", param)

	// TC-SUBCONN-03 正常取消订阅
	err = RawRtdbCancelSubscribeConnectWarp(handle, param)
	if !RteIsOk(err) {
		t.Error("取消订阅失败:", err)
	}
}

// TC-DBINFO-01 获取表文件路径（字符串参数）
func TestRawRtdbGetDbInfo1Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	val, err := RawRtdbGetDbInfo1Warp(handle, RtdbParamTableFile)
	if !RteIsOk(err) {
		t.Error("获取字符串参数失败:", err)
	}
	fmt.Println("表文件路径:", val)
}

// TC-DBINFO-02 获取最大连接数（整型参数）
func TestRawRtdbGetDbInfo2Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	val, err := RawRtdbGetDbInfo2Warp(handle, RtdbParamServerConnectionCount)
	if !RteIsOk(err) {
		t.Error("获取整型参数失败:", err)
	}
	fmt.Println("最大连接数:", val)
}

// TC-SETDBINFO-01 管理员设置系统参数
func TestRawRtdbSetDbInfo1Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbSetDbInfo1Warp(handle, RtdbParamServerSenderIp, "127.0.0.1")
	if !RteIsOk(err) {
		t.Error("设置字符串参数失败:", err)
	}

	// 恢复
	RawRtdbSetDbInfo1Warp(handle, RtdbParamServerSenderIp, "")
}

// TC-SETDBINFO2-01 管理员设置整型参数
func TestRawRtdbSetDbInfo2Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbSetDbInfo2Warp(handle, RtdbParamHashTableSize, 100)
	if !RteIsOk(err) {
		t.Error("设置整型参数失败:", err)
	}
}

// TC-CONNCOUNT-01 获取单机模式连接数
func TestRawRtdbConnectionCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbConnectionCountWarp(handle, 0)
	if !RteIsOk(err) {
		t.Error("获取连接数失败:", err)
	}
	if count < 1 {
		t.Error("连接数应至少为1")
	}
	fmt.Println("连接数:", count)
}

// TC-GETCONNS-01 获取所有连接句柄
func TestRawRtdbGetConnectionsWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, _ := RawRtdbConnectionCountWarp(handle, 0)
	sockets, err := RawRtdbGetConnectionsWarp(handle, 0, count)
	if !RteIsOk(err) {
		t.Error("获取连接句柄失败:", err)
	}
	fmt.Println("连接句柄数量:", len(sockets))
}

// TC-GETCONNS-02 获取当前连接句柄
func TestRawRtdbGetOwnConnectionWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, err := RawRtdbGetOwnConnectionWarp(handle, 0)
	if !RteIsOk(err) {
		t.Error("获取自身连接句柄失败:", err)
	}
	if socket <= 0 {
		t.Error("自身连接句柄应大于0")
	}
	fmt.Println("自身连接句柄:", socket)
}

// TC-CONNINFO-01 获取当前连接 IPv4 信息
func TestRawRtdbGetConnectionInfoWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, _ := RawRtdbGetOwnConnectionWarp(handle, 0)
	info, err := RawRtdbGetConnectionInfoWarp(handle, 0, socket)
	if err != nil {
		t.Error("获取连接信息失败:", err)
	}
	fmt.Printf("连接信息: IP=%d, Port=%d\n", info.IpAddr, info.Port)
}

// TC-CONNINFO-03 获取 IPv6 连接信息
func TestRawRtdbGetConnectionInfoIpv6Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, _ := RawRtdbGetOwnConnectionWarp(handle, 0)
	info, err := RawRtdbGetConnectionInfoIpv6Warp(handle, 0, socket)
	if !RteIsOk(err) {
		t.Logf("获取IPv6信息(可能不支持): %v", err)
		return
	}
	fmt.Println("IPv6连接信息:", info)
}

// TC-OSTYPE-01 获取服务器 OS 类型
func TestRawRtdbOsType(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	osType, err := RawRtdbOsType(handle)
	if !RteIsOk(err) {
		t.Error("获取OS类型失败:", err)
	}
	fmt.Println("OS类型:", osType)
}

// TC-HANDLEINFO-01 获取句柄信息
func TestRawRtdbGetHandleInfoWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	info, err := RawRtdbGetHandleInfoWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取句柄信息失败:", err)
	}
	fmt.Printf("句柄信息: OsType=%v, NewDB=%v\n", info.OsType, info.NewDB)
}

// TC-HOSTTIME-01 获取服务器当前时间
func TestRawRtdbHostTime64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ts, err := RawRtdbHostTime64Warp(handle)
	if !RteIsOk(err) {
		t.Error("获取服务器时间失败:", err)
	}
	now := time.Now().Unix()
	diff := ts - TimestampType(now)
	if diff < -60 || diff > 60 {
		t.Errorf("服务器时间与本地时间偏差过大: %d秒", diff)
	}
	fmt.Println("服务器时间戳:", ts)
}

// TC-TIMESPAN-01/02 格式化时间跨度
func TestRawRtdbFormatTimespanWarp(t *testing.T) {
	s10, err := RawRtdbFormatTimespanWarp(10)
	if !RteIsOk(err) {
		t.Error("格式化10秒失败:", err)
	}
	fmt.Println("10秒格式:", s10)

	s60, err := RawRtdbFormatTimespanWarp(60)
	if !RteIsOk(err) {
		t.Error("格式化60秒失败:", err)
	}
	fmt.Println("60秒格式:", s60)
}

// TC-TIMESPAN-03/04 解析时间跨度字符串
func TestRawRtdbParseTimespanWarp(t *testing.T) {
	dt, err := RawRtdbParseTimespanWarp("2n")
	if !RteIsOk(err) {
		t.Error("解析2n失败:", err)
	}
	if dt != 120 {
		t.Error("解析2n应等于120秒")
	}

	_, err = RawRtdbParseTimespanWarp("invalid")
	if RteIsOk(err) {
		t.Error("解析无效字符串应失败")
	}
}

// TC-PARSETIME-01/02/03 解析时间字符串
func TestRawRtdbParseTimeWarp(t *testing.T) {
	ts, sub, err := RawRtdbParseTimeWarp("2024-01-01 08:00:00")
	if !RteIsOk(err) {
		t.Error("解析绝对时间失败:", err)
	}
	fmt.Println("解析时间:", ts, sub)

	ts2, _, err := RawRtdbParseTimeWarp("*-1d")
	if !RteIsOk(err) {
		t.Error("解析相对时间失败:", err)
	}
	fmt.Println("相对时间:", ts2)

	_, _, err = RawRtdbParseTimeWarp("not_a_time")
	if RteIsOk(err) {
		t.Error("解析无效时间应失败")
	}
}

// TC-FMTMSG-01/02 格式化错误码描述
func TestRawRtdbFormatMessageWarp(t *testing.T) {
	name, msg := RawRtdbFormatMessageWarp(RteOk)
	fmt.Println("RteOk描述:", name, msg)

	name2, msg2 := RawRtdbFormatMessageWarp(RtePointNotFound)
	fmt.Println("RtePointNotFound描述:", name2, msg2)
}

// TC-JOBMSG-01 获取任务描述
func TestRawRtdbJobMessageWarp(t *testing.T) {
	name, msg := RawRtdbJobMessageWarp(0)
	fmt.Println("Job 0描述:", name, msg)
}

// TC-TIMEOUT-01/03 设置并读取超时时间
func TestRawRtdbSetTimeoutWarp_GetTimeoutWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, _ := RawRtdbGetOwnConnectionWarp(handle, 0)

	err := RawRtdbSetTimeoutWarp(handle, socket, 30)
	if !RteIsOk(err) {
		t.Error("设置超时失败:", err)
	}

	to, err := RawRtdbGetTimeoutWarp(handle, socket)
	if !RteIsOk(err) {
		t.Error("获取超时失败:", err)
	}
	if to != 30 {
		t.Errorf("超时时间应为30, 实际=%d", to)
	}
	fmt.Println("超时时间:", to)
}

// TC-JUDGE-01 判断正常连接状态
func TestRawRtdbJudgeConnectStatusWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbJudgeConnectStatusWarp(handle)
	if !RteIsOk(err) {
		t.Error("判断连接状态失败:", err)
	}
}

// TC-FMTIP-01/02/03 IP格式化
func TestRawRtdbFormatIpaddrWarp(t *testing.T) {
	ip1 := RawRtdbFormatIpaddrWarp(0x7F000001)
	if ip1 != "127.0.0.1" {
		t.Errorf("127.0.0.1格式化错误: %s", ip1)
	}

	ip2 := RawRtdbFormatIpaddrWarp(0)
	if ip2 != "0.0.0.0" {
		t.Errorf("0.0.0.0格式化错误: %s", ip2)
	}

	ip3 := RawRtdbFormatIpaddrWarp(0xFFFFFFFF)
	if ip3 != "255.255.255.255" {
		t.Errorf("255.255.255.255格式化错误: %s", ip3)
	}
}

// ==================== 02. 用户权限与连接管理 ====================

func TestRawRtdbChangePasswordWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddUserWarp(handle, "testuser_chpwd", "OldPass123", PrivGroupRtdbRO)
	defer RawRtdbRemoveUserWarp(handle, "testuser_chpwd")
	err := RawRtdbChangePasswordWarp(handle, "testuser_chpwd", "NewPass123")
	if !RteIsOk(err) {
		t.Error("修改密码失败:", err)
	}
}

func TestRawRtdbChangeMyPasswordWarp(t *testing.T) {
	t.Log("ChangeMyPassword需要普通用户登录测试，跳过")
}

func TestRawRtdbGetPrivWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	priv, err := RawRtdbGetPrivWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取权限失败:", err)
	}
	fmt.Println("当前权限:", priv)
}

func TestRawRtdbChangePrivWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddUserWarp(handle, "testuser_priv", "Pass123", PrivGroupRtdbRO)
	defer RawRtdbRemoveUserWarp(handle, "testuser_priv")
	err := RawRtdbChangePrivWarp(handle, "testuser_priv", PrivGroupRtdbDW)
	if !RteIsOk(err) {
		t.Error("修改权限失败:", err)
	}
}

func TestRawRtdbAddUserWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	err := RawRtdbAddUserWarp(handle, "testuser_add", "Pass123", PrivGroupRtdbRO)
	if !RteIsOk(err) {
		t.Error("添加用户失败:", err)
	}
	defer RawRtdbRemoveUserWarp(handle, "testuser_add")
}

func TestRawRtdbRemoveUserWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddUserWarp(handle, "testuser_rmv", "Pass123", PrivGroupRtdbRO)
	err := RawRtdbRemoveUserWarp(handle, "testuser_rmv")
	if !RteIsOk(err) {
		t.Error("删除用户失败:", err)
	}
}

func TestRawRtdbLockUserWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddUserWarp(handle, "testuser_lock", "Pass123", PrivGroupRtdbRO)
	defer RawRtdbRemoveUserWarp(handle, "testuser_lock")
	err := RawRtdbLockUserWarp(handle, "testuser_lock", OFF)
	if !RteIsOk(err) {
		t.Error("禁用用户失败:", err)
	}
	err = RawRtdbLockUserWarp(handle, "testuser_lock", ON)
	if !RteIsOk(err) {
		t.Error("启用用户失败:", err)
	}
}

func TestRawRtdbGetUsersWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	users, err := RawRtdbGetUsersWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取用户列表失败:", err)
	}
	fmt.Println("用户数量:", len(users))
}

func TestRawRtdbAddBlacklistWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	err := RawRtdbAddBlacklistWarp(handle, "192.168.1.100", "255.255.255.255", "blocked")
	if !RteIsOk(err) {
		t.Error("添加黑名单失败:", err)
	}
	defer RawRtdbRemoveBlacklistWarp(handle, "192.168.1.100", "255.255.255.255")
}

func TestRawRtdbUpdateBlacklistWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddBlacklistWarp(handle, "192.168.1.0", "255.255.255.0", "网段封禁")
	defer RawRtdbRemoveBlacklistWarp(handle, "192.168.2.0", "255.255.255.0")
	err := RawRtdbUpdateBlacklistWarp(handle, "192.168.1.0", "255.255.255.0", "192.168.2.0", "255.255.255.0", "new desc")
	if !RteIsOk(err) {
		t.Error("更新黑名单失败:", err)
	}
}

func TestRawRtdbRemoveBlacklistWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddBlacklistWarp(handle, "192.168.1.0", "255.255.255.0", "待删除")
	err := RawRtdbRemoveBlacklistWarp(handle, "192.168.1.0", "255.255.255.0")
	if !RteIsOk(err) {
		t.Error("删除黑名单失败:", err)
	}
}

func TestRawRtdbGetBlacklistWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	list, err := RawRtdbGetBlacklistWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取黑名单失败:", err)
	}
	fmt.Println("黑名单数量:", len(list))
}

func TestRawRtdbAddAuthorizationWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	err := RawRtdbAddAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "office", PrivGroupRtdbRO)
	if !RteIsOk(err) {
		t.Error("添加信任段失败:", err)
	}
	defer RawRtdbRemoveAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0")
}

func TestRawRtdbUpdateAuthorizationWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "office", PrivGroupRtdbRO)
	defer RawRtdbRemoveAuthorizationWarp(handle, "192.168.2.0", "255.255.255.0")
	err := RawRtdbUpdateAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "192.168.2.0", "255.255.255.0", "new office", PrivGroupRtdbDW)
	if !RteIsOk(err) {
		t.Error("更新信任段失败:", err)
	}
}

func TestRawRtdbRemoveAuthorizationWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	_ = RawRtdbAddAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "待删除", PrivGroupRtdbRO)
	err := RawRtdbRemoveAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0")
	if !RteIsOk(err) {
		t.Error("删除信任段失败:", err)
	}
}

func TestRawRtdbGetAuthorizationsWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	list, err := RawRtdbGetAuthorizationsWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取信任段失败:", err)
	}
	fmt.Println("信任段数量:", len(list))
}

// ==================== 03. 文件目录操作 ====================

func TestRawRtdbGetLogicalDriversWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	drivers, err := RawRtdbGetLogicalDriversWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取盘符失败:", err)
	}
	fmt.Println("盘符:", drivers)
}

func TestRawRtdbOpenPathWarp_ReadPath64Warp_ClosePathWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	err := RawRtdbOpenPathWarp(handle, "/")
	if !RteIsOk(err) {
		t.Logf("打开根目录(可能平台差异): %v", err)
		return
	}
	item, err := RawRtdbReadPath64Warp(handle)
	if !RteIsOk(err) {
		t.Logf("读取目录条目: %v", err)
	} else {
		fmt.Printf("目录条目: %s, IsDir=%d, Size=%d\n", item.Path, item.IsDir, item.Size)
	}
	err = RawRtdbClosePathWarp(handle)
	if !RteIsOk(err) {
		t.Error("关闭目录失败:", err)
	}
}

func TestRawRtdbMkdirWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	err := RawRtdbMkdirWarp(handle, "/data/testdir_go")
	if !RteIsOk(err) {
		t.Logf("创建目录(可能已存在): %v", err)
	}
}

func TestRawRtdbGetFileSizeWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	size, err := RawRtdbGetFileSizeWarp(handle, "/data/test.txt")
	if !RteIsOk(err) {
		t.Logf("获取文件大小(可能文件不存在): %v", err)
		return
	}
	fmt.Println("文件大小:", size)
}

func TestRawRtdbReadFileWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	data, err := RawRtdbReadFileWarp(handle, "/data/test.txt", 0, 1024)
	if !RteIsOk(err) {
		t.Logf("读取文件(可能文件不存在): %v", err)
		return
	}
	fmt.Println("读取文件长度:", len(data))
}

func TestRawRtdbGetMaxBlobLenWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	maxLen, err := RawRtdbGetMaxBlobLenWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取最大Blob长度失败:", err)
	}
	fmt.Println("最大Blob长度:", maxLen)
}

func TestRawRtdbFormatQualityWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	descs, err := RawRtdbFormatQualityWarp(handle, []Quality{192, 0, 64})
	if !RteIsOk(err) {
		t.Error("格式化质量码失败:", err)
	}
	fmt.Println("质量码描述:", descs)
}

// ==================== 04. 表管理 ====================

func TestRawRtdbbAppendTableWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTableGo", "测试表GO")
	if !RteIsOk(err) {
		t.Error("创建表失败:", err)
	}
	fmt.Println("创建表ID:", tbl.ID)
	defer RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
}

func TestRawRtdbbRemoveTableByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tbl, _ := RawRtdbbAppendTableWarp(handle, "TestTableRmId", "待删除")
	err := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
	if !RteIsOk(err) {
		t.Error("按ID删除表失败:", err)
	}
}

func TestRawRtdbbRemoveTableByNameWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_, _ = RawRtdbbAppendTableWarp(handle, "TestTableRmName", "待删除")
	err := RawRtdbbRemoveTableByNameWarp(handle, "TestTableRmName")
	if !RteIsOk(err) {
		t.Error("按名称删除表失败:", err)
	}
}

func TestRawRtdbbTablesCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbbTablesCountWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取表总数失败:", err)
	}
	if count < 0 {
		t.Error("表总数不应小于0")
	}
	fmt.Println("表总数:", count)
}

func TestRawRtdbbGetTablesWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, _ := RawRtdbbTablesCountWarp(handle)
	ids, err := RawRtdbbGetTablesWarp(handle, count)
	if !RteIsOk(err) {
		t.Error("获取表ID列表失败:", err)
	}
	fmt.Println("表数量:", len(ids))
}

func TestRawRtdbbGetTableSizeByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbGetTablesWarp(handle, 100)
	if len(ids) == 0 {
		t.Skip("无表可测")
	}
	size, err := RawRtdbbGetTableSizeByIdWarp(handle, ids[0])
	if !RteIsOk(err) {
		t.Error("获取表大小失败:", err)
	}
	fmt.Println("表大小:", size)
}

func TestRawRtdbbGetTableSizeByNameWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbGetTablesWarp(handle, 100)
	if len(ids) == 0 {
		t.Skip("无表可测")
	}
	prop, _ := RawRtdbbGetTablePropertyByIdWarp(handle, ids[0])
	size, err := RawRtdbbGetTableSizeByNameWarp(handle, prop.Name)
	if !RteIsOk(err) {
		t.Error("按名称获取表大小失败:", err)
	}
	fmt.Println("表大小(ByName):", size)
}

func TestRawRtdbbGetTableRealSizeByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbGetTablesWarp(handle, 100)
	if len(ids) == 0 {
		t.Skip("无表可测")
	}
	size, err := RawRtdbbGetTableRealSizeByIdWarp(handle, ids[0])
	if !RteIsOk(err) {
		t.Error("获取表实际大小失败:", err)
	}
	fmt.Println("表实际大小:", size)
}

func TestRawRtdbbGetTablePropertyByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbGetTablesWarp(handle, 100)
	if len(ids) == 0 {
		t.Skip("无表可测")
	}
	prop, err := RawRtdbbGetTablePropertyByIdWarp(handle, ids[0])
	if !RteIsOk(err) {
		t.Error("获取表属性失败:", err)
	}
	fmt.Printf("表属性: ID=%d, Name=%s\n", prop.ID, prop.Name)
}

func TestRawRtdbbGetTablePropertyByNameWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbGetTablesWarp(handle, 100)
	if len(ids) == 0 {
		t.Skip("无表可测")
	}
	prop, _ := RawRtdbbGetTablePropertyByIdWarp(handle, ids[0])
	prop2, err := RawRtdbbGetTablePropertyByNameWarp(handle, prop.Name)
	if !RteIsOk(err) {
		t.Error("按名称获取表属性失败:", err)
	}
	if prop2.ID != prop.ID {
		t.Error("按名称和按ID获取的表属性不一致")
	}
}

func TestRawRtdbbUpdateTableNameWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tbl, _ := RawRtdbbAppendTableWarp(handle, "TestTableOldName", "测试")
	defer RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
	err := RawRtdbbUpdateTableNameWarp(handle, tbl.ID, "TestTableNewName")
	if !RteIsOk(err) {
		t.Error("更新表名失败:", err)
	}
	prop, _ := RawRtdbbGetTablePropertyByIdWarp(handle, tbl.ID)
	if prop.Name != "TestTableNewName" {
		t.Error("表名未更新")
	}
}

func TestRawRtdbbUpdateTableDescByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tbl, _ := RawRtdbbAppendTableWarp(handle, "TestTableDescId", "旧描述")
	defer RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
	err := RawRtdbbUpdateTableDescByIdWarp(handle, tbl.ID, "新描述GO")
	if !RteIsOk(err) {
		t.Error("更新表描述失败:", err)
	}
}

func TestRawRtdbbUpdateTableDescByNameWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_, _ = RawRtdbbAppendTableWarp(handle, "TestTableDescName", "旧描述")
	defer RawRtdbbRemoveTableByNameWarp(handle, "TestTableDescName")
	err := RawRtdbbUpdateTableDescByNameWarp(handle, "TestTableDescName", "新描述ByName")
	if !RteIsOk(err) {
		t.Error("按名称更新表描述失败:", err)
	}
}

// ==================== 05. 标签点管理 ====================

func getFirstTableID(t *testing.T, handle ConnectHandle) TableID {
	ids, err := RawRtdbbGetTablesWarp(handle, 100)
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("无可用表")
	}
	return ids[0]
}

func TestRawRtdbbInsertBasePointWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestBaseBool", RtdbTypeBool, tableID, 0)
	if !RteIsOk(err) {
		t.Error("创建基础点失败:", err)
	}
	fmt.Println("创建基础点ID:", pid)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)
}

func TestRawRtdbbInsertPointWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	base := &RtdbPoint{Tag: "TestInsertPt", Table: tableID, Type: RtdbTypeReal64}
	scan := &RtdbScan{Source: "go_test"}
	base, _, _, err := RawRtdbbInsertPointWarp(handle, base, scan, nil)
	if !RteIsOk(err) {
		t.Error("创建完整点失败:", err)
	}
	fmt.Println("创建完整点ID:", base.ID)
	defer RawRtdbbRemovePointByIdWarp(handle, base.ID)
}

func TestRawRtdbbInsertMaxPointWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	base := &RtdbPoint{Tag: "TestMaxPt", Table: tableID, Type: RtdbTypeReal64}
	scan := &RtdbScan{Source: "go_test"}
	base, _, _, err := RawRtdbbInsertMaxPointWarp(handle, base, scan, nil)
	if !RteIsOk(err) {
		t.Error("创建Max点失败:", err)
	}
	defer RawRtdbbRemovePointByIdWarp(handle, base.ID)
}

func TestRawRtdbbRemovePointByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	pid, _ := RawRtdbbInsertBasePointWarp(handle, "TestRmvPtId", RtdbTypeBool, tableID, 0)
	err := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(err) {
		t.Error("按ID删除点失败:", err)
	}
}

func TestRawRtdbbRemovePointByNameWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbGetTablesWarp(handle, 100)
	if len(ids) == 0 {
		t.Skip("无表")
	}
	prop, _ := RawRtdbbGetTablePropertyByIdWarp(handle, ids[0])
	_, _ = RawRtdbbInsertBasePointWarp(handle, "TestRmvPtName", RtdbTypeBool, ids[0], 0)
	fullName := prop.Name + ".TestRmvPtName"
	err := RawRtdbbRemovePointByNameWarp(handle, fullName)
	if !RteIsOk(err) {
		t.Error("按名称删除点失败:", err)
	}
}

func TestRawRtdbbInsertMaxPointsWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	bases := make([]RtdbPoint, 3)
	scans := make([]RtdbScan, 3)
	calcs := make([]RtdbCalc, 3)
	for i := 0; i < 3; i++ {
		bases[i] = RtdbPoint{Tag: fmt.Sprintf("TestBatchPt%d", i), Table: tableID, Type: RtdbTypeReal64}
		scans[i] = RtdbScan{Source: "go_test"}
	}
	bases, _, _, errs, err := RawRtdbbInsertMaxPointsWarp(handle, bases, scans, calcs)
	if !RteIsOk(err) {
		t.Error("批量创建点失败:", err)
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			t.Logf("第%d个点创建出错: %v", i, e)
		} else {
			defer RawRtdbbRemovePointByIdWarp(handle, bases[i].ID)
		}
	}
}

func TestRawRtdbbInsertNamedTypePointWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	base := &RtdbPoint{Tag: "TestNamedPt", Table: tableID}
	scan := &RtdbScan{Source: "go_test"}
	_, _, err := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "NotExistType")
	if RteIsOk(err) {
		t.Error("不存在的类型应失败")
	}
	fmt.Println("不存在的自定义类型结果:", err)
}

func TestRawRtdbbMovePointByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbGetTablesWarp(handle, 100)
	if len(ids) < 2 {
		t.Skip("需要至少2个表测试移动")
	}
	prop1, _ := RawRtdbbGetTablePropertyByIdWarp(handle, ids[0])
	prop2, _ := RawRtdbbGetTablePropertyByIdWarp(handle, ids[1])
	pid, _ := RawRtdbbInsertBasePointWarp(handle, "TestMovePt", RtdbTypeBool, ids[0], 0)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)
	err := RawRtdbbMovePointByIdWarp(handle, pid, prop2.Name)
	if !RteIsOk(err) {
		t.Logf("移动点(可能权限或约束): %v", err)
		return
	}
	_ = prop1
}

func TestRawRtdbbGetPointsPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("未找到标签点")
	}
	count := 3
	if len(ids) < count {
		count = len(ids)
	}
	bases, scans, calcs, errs, err := RawRtdbbGetPointsPropertyWarp(handle, ids[:count])
	if !RteIsOk(err) {
		t.Error("批量获取点属性失败:", err)
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			t.Logf("第%d个点获取属性出错: %v", i, e)
		}
	}
	fmt.Println("获取属性点数:", len(bases), len(scans), len(calcs))
}

func TestRawRtdbbGetMaxPointsPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("未找到标签点")
	}
	_, _, _, _, err = RawRtdbbGetMaxPointsPropertyWarp(handle, ids[:1])
	if !RteIsOk(err) {
		t.Error("Max获取点属性失败:", err)
	}
}

func TestRawRtdbbGetTypesPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("未找到标签点")
	}
	types, errs, err := RawRtdbbGetTypesPropertyWarp(handle, ids[:3])
	if !RteIsOk(err) {
		t.Error("批量获取类型失败:", err)
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			t.Logf("第%d个点获取类型出错: %v", i, e)
		}
	}
	fmt.Println("类型列表:", types)
}

func TestRawRtdbbSearchWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		t.Error("搜索失败:", err)
	}
	fmt.Println("搜索到点数:", len(ids))

	ids2, err := RawRtdbbSearchWarp(handle, "NOTEXIST_*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		t.Error("搜索无匹配失败:", err)
	}
	fmt.Println("无匹配点数:", len(ids2))
}

func TestRawRtdbbSearchInBatchesWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) <= 1 {
		t.Skip("点数不足分批测试")
	}
	batch, err := RawRtdbbSearchInBatchesWarp(handle, 0, 1, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		t.Error("分批搜索失败:", err)
	}
	fmt.Println("第一批数量:", len(batch))
}

func TestRawRtdbbSearchExWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "float64", 0, 0, 0, "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		t.Error("高级搜索失败:", err)
	}
	fmt.Println("高级搜索点数:", len(ids))
}

func TestRawRtdbbSearchPointsCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbbSearchPointsCountWarp(handle, "*", "*", "", "", "", "", "", 0, 0, 0, "")
	if !RteIsOk(err) {
		t.Error("统计点数失败:", err)
	}
	fmt.Println("总点数:", count)
}

func TestRawRtdbbFindPointsWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("无标签点")
	}
	bases, _, _, _, _ := RawRtdbbGetPointsPropertyWarp(handle, ids[:1])
	if len(bases) == 0 {
		t.Skip("无法获取点属性")
	}
	fullName := bases[0].TableDotTag
	fids, types, classes, useMs, err := RawRtdbbFindPointsWarp(handle, []string{fullName, "Table.NotExist"})
	if !RteIsOk(err) {
		t.Error("查找点失败:", err)
	}
	fmt.Println("查找结果:", fids, types, classes, useMs)
}

func TestRawRtdbbFindPointsExWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("无标签点")
	}
	bases, _, _, _, _ := RawRtdbbGetPointsPropertyWarp(handle, ids[:1])
	if len(bases) == 0 {
		t.Skip("无法获取点属性")
	}
	fullName := bases[0].TableDotTag
	fids, types, classes, precisions, errs, err := RawRtdbbFindPointsExWarp(handle, []string{fullName})
	if !RteIsOk(err) {
		t.Error("查找点Ex失败:", err)
	}
	fmt.Println("查找Ex结果:", fids, types, classes, precisions, errs)
}

func TestRawRtdbbSortPointsWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) < 2 {
		t.Skip("点数不足排序测试")
	}
	sorted, err := RawRtdbbSortPointsWarp(handle, ids[:5], RtdbTagIndexTag, RtdbSortFlag(0))
	if !RteIsOk(err) {
		t.Error("排序失败:", err)
	}
	fmt.Println("排序后前5个ID:", sorted)
}

func TestRawRtdbbUpdatePointPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	base := &RtdbPoint{Tag: "TestUpdPt", Table: tableID, Type: RtdbTypeReal64}
	base, _, _, _ = RawRtdbbInsertPointWarp(handle, base, &RtdbScan{Source: "go_test"}, nil)
	defer RawRtdbbRemovePointByIdWarp(handle, base.ID)

	base.Desc = "UpdatedDesc"
	err := RawRtdbbUpdatePointPropertyWarp(handle, base, nil, nil)
	if !RteIsOk(err) {
		t.Error("更新点属性失败:", err)
	}
}

func TestRawRtdbbUpdateMaxPointPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	base := &RtdbPoint{Tag: "TestUpdMaxPt", Table: tableID, Type: RtdbTypeReal64}
	base, _, _, _ = RawRtdbbInsertMaxPointWarp(handle, base, &RtdbScan{Source: "go_test"}, nil)
	defer RawRtdbbRemovePointByIdWarp(handle, base.ID)

	base.Desc = "UpdatedMaxDesc"
	err := RawRtdbbUpdateMaxPointPropertyWarp(handle, base, nil, nil)
	if !RteIsOk(err) {
		t.Error("更新Max点属性失败:", err)
	}
}

// ==================== 06. 回收站与自定义类型 ====================

func TestRawRtdbbRecoverPointWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	pid, _ := RawRtdbbInsertBasePointWarp(handle, "TestRecyPt", RtdbTypeBool, tableID, 0)
	RawRtdbbRemovePointByIdWarp(handle, pid)

	recycled, _ := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if len(recycled) == 0 {
		t.Skip("回收站为空")
	}
	err := RawRtdbbRecoverPointWarp(handle, tableID, recycled[0])
	if !RteIsOk(err) {
		t.Logf("恢复点(可能回收站策略差异): %v", err)
	}
}

func TestRawRtdbbPurgePointWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	tableID := getFirstTableID(t, handle)
	pid, _ := RawRtdbbInsertBasePointWarp(handle, "TestPurgePt", RtdbTypeBool, tableID, 0)
	RawRtdbbRemovePointByIdWarp(handle, pid)

	recycled, _ := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if len(recycled) == 0 {
		t.Skip("回收站为空")
	}
	err := RawRtdbbPurgePointWarp(handle, recycled[0])
	if !RteIsOk(err) {
		t.Logf("清除点(可能回收站策略差异): %v", err)
	}
}

func TestRawRtdbbGetRecycledPointsCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbbGetRecycledPointsCountWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取回收站数量失败:", err)
	}
	fmt.Println("回收站数量:", count)
}

func TestRawRtdbbGetRecycledPointsWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	points, err := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if !RteIsOk(err) {
		t.Error("获取回收站点失败:", err)
	}
	fmt.Println("回收站点数:", len(points))
}

func TestRawRtdbbSearchRecycledPointsWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchRecycledPointsWarp(handle, "*", "", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		t.Error("搜索回收站点失败:", err)
	}
	fmt.Println("搜索回收站点数:", len(ids))
}

func TestRawRtdbbGetRecycledPointPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	recycled, _ := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if len(recycled) == 0 {
		t.Skip("回收站为空")
	}
	base, scan, calc, err := RawRtdbbGetRecycledPointPropertyWarp(handle, recycled[0])
	if !RteIsOk(err) {
		t.Logf("获取回收站点属性: %v", err)
		return
	}
	fmt.Println("回收站点属性:", base, scan, calc)
}

func TestRawRtdbbSearchRecycledPointsInBatchesWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	batch, err := RawRtdbbSearchRecycledPointsInBatchesWarp(handle, 0, 5, "*", "", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		t.Error("分批搜索回收站失败:", err)
	}
	fmt.Println("第一批回收站点:", len(batch))
}

func TestRawRtdbbGetRecycledMaxPointPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	recycled, _ := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if len(recycled) == 0 {
		t.Skip("回收站为空")
	}
	base, scan, calc, err := RawRtdbbGetRecycledMaxPointPropertyWarp(handle, recycled[0])
	if !RteIsOk(err) {
		t.Logf("获取Max回收站点属性: %v", err)
		return
	}
	fmt.Println("Max回收站点属性:", base, scan, calc)
}

func TestRawRtdbbClearRecyclerWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbbClearRecyclerWarp(handle)
	if !RteIsOk(err) {
		t.Logf("清空回收站: %v", err)
	}
}

func TestRawRtdbbGetRecycledNamedTypeNamesPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	recycled, _ := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if len(recycled) == 0 {
		t.Skip("回收站为空")
	}
	names, counts, errs, err := RawRtdbbGetRecycledNamedTypeNamesPropertyWarp(handle, recycled[:1])
	if !RteIsOk(err) {
		t.Logf("获取回收站自定义类型信息: %v", err)
		return
	}
	fmt.Println("回收站自定义类型:", names, counts, errs)
}

func TestRawRtdbbCreateNamedTypeWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestNamedType")
	fields := []RtdbDataTypeField{
		{Name: "field1", Type: RtdbTypeReal64, Length: 8, Desc: "字段1"},
		{Name: "field2", Type: RtdbTypeInt32, Length: 4, Desc: "字段2"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestNamedType", "测试自定义类型", fields...)
	if !RteIsOk(err) {
		t.Error("创建自定义类型失败:", err)
	}
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestNamedType")
}

func TestRawRtdbbGetNamedTypesCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbbGetNamedTypesCountWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取自定义类型总数失败:", err)
	}
	fmt.Println("自定义类型总数:", count)
}

func TestRawRtdbbGetAllNamedTypesWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	names, counts, err := RawRtdbbGetAllNamedTypesWarp(handle, 100)
	if !RteIsOk(err) {
		t.Error("获取所有自定义类型失败:", err)
	}
	fmt.Println("自定义类型数量:", len(names), len(counts))
}

func TestRawRtdbbGetNamedTypeWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestGetNamedType")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestGetNamedType", "查询测试", fields...)
	if !RteIsOk(err) {
		t.Skip("无法创建自定义类型:", err)
	}
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestGetNamedType")

	gotFields, typeSize, desc, err := RawRtdbbGetNamedTypeWarp(handle, "TestGetNamedType", 10)
	if !RteIsOk(err) {
		t.Error("获取自定义类型失败:", err)
	}
	fmt.Println("自定义类型字段:", len(gotFields), "size:", typeSize, "desc:", desc)
}

func TestRawRtdbbRemoveNamedTypeWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	_ = RawRtdbbCreateNamedTypeWarp(handle, "TestRmNamedType", "删除测试", fields...)
	err := RawRtdbbRemoveNamedTypeWarp(handle, "TestRmNamedType")
	if !RteIsOk(err) {
		t.Error("删除自定义类型失败:", err)
	}
}

func TestRawRtdbbGetNamedTypeNamesPropertyWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("无标签点")
	}
	names, counts, errs, err := RawRtdbbGetNamedTypeNamesPropertyWarp(handle, ids[:3])
	if !RteIsOk(err) {
		t.Error("获取自定义类型名称属性失败:", err)
	}
	fmt.Println("自定义类型名称:", names, counts, errs)
}

func TestRawRtdbbGetNamedTypePointsCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbbGetNamedTypePointsCountWarp(handle, "TestNamedType")
	if !RteIsOk(err) {
		t.Logf("获取自定义类型点数量(可能类型不存在): %v", err)
		return
	}
	fmt.Println("自定义类型点数量:", count)
}

func TestRawRtdbbGetBaseTypePointsCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbbGetBaseTypePointsCountWarp(handle, RtdbTypeReal64)
	if !RteIsOk(err) {
		t.Error("获取基础类型点数量失败:", err)
	}
	fmt.Println("float64类型点数量:", count)
}

func TestRawRtdbbModifyNamedTypeWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestModTypeOld")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestModTypeNew")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestModTypeOld", "修改测试", fields...)
	if !RteIsOk(err) {
		t.Skip("无法创建自定义类型:", err)
	}
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestModTypeNew")

	newName := "TestModTypeNew"
	err = RawRtdbbModifyNamedTypeWarp(handle, "TestModTypeOld", &newName, nil, []string{"newf1"}, []string{"新字段1"})
	if !RteIsOk(err) {
		t.Error("修改自定义类型失败:", err)
	}
}

func TestRawRtdbWriteNamedTypeFieldByName32Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestWrType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestWrType", "写测试", fields...)
	if !RteIsOk(err) {
		t.Skip("无法创建自定义类型:", err)
	}
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestWrType")

	object := make([]byte, 8)
	field := make([]byte, 8)
	// 写入一个 float64 值
	field[0] = 0x3F
	field[1] = 0xF0
	object, err = RawRtdbWriteNamedTypeFieldByName32Warp(handle, "TestWrType", "val", RtdbTypeReal64, object, field)
	if !RteIsOk(err) {
		t.Error("按名称写字段失败:", err)
	}
}

func TestRawRtdbWriteNamedTypeFieldByPos32Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestWrPosType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestWrPosType", "写位置测试", fields...)
	if !RteIsOk(err) {
		t.Skip("无法创建自定义类型:", err)
	}
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestWrPosType")

	object := make([]byte, 8)
	field := make([]byte, 8)
	object, err = RawRtdbWriteNamedTypeFieldByPos32Warp(handle, "TestWrPosType", 0, RtdbTypeReal64, object, field)
	if !RteIsOk(err) {
		t.Error("按位置写字段失败:", err)
	}
}

func TestRawRtdbReadNamedTypeFieldByName32Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestRdType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestRdType", "读测试", fields...)
	if !RteIsOk(err) {
		t.Skip("无法创建自定义类型:", err)
	}
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestRdType")

	object := make([]byte, 8)
	object[0] = 0x3F
	object[1] = 0xF0
	field, err := RawRtdbReadNamedTypeFieldByName32Warp(handle, "TestRdType", "val", RtdbTypeReal64, object, 8)
	if !RteIsOk(err) {
		t.Error("按名称读字段失败:", err)
	}
	fmt.Println("读取字段长度:", len(field))
}

func TestRawRtdbReadNamedTypeFieldByPos32Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestRdPosType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestRdPosType", "读位置测试", fields...)
	if !RteIsOk(err) {
		t.Skip("无法创建自定义类型:", err)
	}
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestRdPosType")

	object := make([]byte, 8)
	object[0] = 0x3F
	object[1] = 0xF0
	field, err := RawRtdbReadNamedTypeFieldByPos32Warp(handle, "TestRdPosType", 0, RtdbTypeReal64, object, 8)
	if !RteIsOk(err) {
		t.Error("按位置读字段失败:", err)
	}
	fmt.Println("读取字段长度:", len(field))
}

func TestRawRtdbNamedTypeNameFieldCheckWarp(t *testing.T) {
	err := RawRtdbNamedTypeNameFieldCheckWarp("ValidType", 0)
	if !RteIsOk(err) {
		t.Error("合法类型名称应通过:", err)
	}

	err = RawRtdbNamedTypeNameFieldCheckWarp("valid_field", 1)
	if !RteIsOk(err) {
		t.Error("合法字段名称应通过:", err)
	}

	err = RawRtdbNamedTypeNameFieldCheckWarp("Type@123", 0)
	if RteIsOk(err) {
		t.Error("非法名称应失败")
	}
}

// ==================== 07. 实时数据快照 ====================

func getFirstPointID(t *testing.T, handle ConnectHandle) PointID {
	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) || len(ids) == 0 {
		t.Skip("无可用标签点")
	}
	return ids[0]
}

func TestRawRtdbsGetSnapshots64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	_, _, _, _, _, _, err := RawRtdbsGetSnapshots64Warp(handle, []PointID{pid})
	if !RteIsOk(err) {
		t.Error("获取快照失败:", err)
	}
}

func TestRawRtdbsPutSnapshots64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbsPutSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{123.45}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		t.Error("写入快照失败:", err)
	}
}

func TestRawRtdbsFixSnapshots64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbsFixSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{99.9}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		t.Logf("覆盖快照(可能时间戳问题): %v", err)
	}
}

func TestRawRtdbsBackSnapshots64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbsBackSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{88.8}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		t.Logf("回溯快照(可能时间戳问题): %v", err)
	}
}

func TestRawRtdbsGetBlobSnapshot64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "string", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无字符串点")
	}
	_, _, blob, _, err := RawRtdbsGetBlobSnapshot64Warp(handle, ids[0], true, 256)
	if !RteIsOk(err) {
		t.Logf("获取Blob快照: %v", err)
		return
	}
	fmt.Println("Blob长度:", len(blob))
}

func TestRawRtdbsPutBlobSnapshot64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "string", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无字符串点")
	}
	now := TimestampType(time.Now().Unix())
	err := RawRtdbsPutBlobSnapshot64Warp(handle, ids[0], true, now, 0, []byte("hello_rtdb"), 0)
	if !RteIsOk(err) {
		t.Logf("写入Blob快照: %v", err)
	}
}

func TestRawRtdbsGetDatetimeSnapshots64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "datetime", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无datetime点")
	}
	_, _, vals, _, _, err := RawRtdbsGetDatetimeSnapshots64Warp(handle, ids[:1], 0)
	if !RteIsOk(err) {
		t.Logf("获取datetime快照: %v", err)
		return
	}
	fmt.Println("datetime值:", vals)
}

func TestRawRtdbsPutDatetimeSnapshots64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "datetime", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无datetime点")
	}
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbsPutDatetimeSnapshots64Warp(handle, ids[:1], []TimestampType{now}, []SubtimeType{0}, []string{"2024-01-01 08:00:00"}, []Quality{0})
	if !RteIsOk(err) {
		t.Logf("写入datetime快照: %v", err)
	}
}

func TestRawRtdbsGetNamedTypeSnapshot64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无自定义类型点")
	}
	_, _, obj, _, err := RawRtdbsGetNamedTypeSnapshot64Warp(handle, ids[0], 256)
	if !RteIsOk(err) {
		t.Logf("获取自定义类型快照: %v", err)
		return
	}
	fmt.Println("自定义类型数据长度:", len(obj))
}

func TestRawRtdbsPutNamedTypeSnapshot64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无自定义类型点")
	}
	now := TimestampType(time.Now().Unix())
	err := RawRtdbsPutNamedTypeSnapshot64Warp(handle, ids[0], now, 0, make([]byte, 8), 0)
	if !RteIsOk(err) {
		t.Logf("写入自定义类型快照: %v", err)
	}
}

// ==================== 08. 订阅功能 ====================

func TestRawRtdbbSubscribeTagsExWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	param, err := RawRtdbbSubscribeTagsExWarp(handle, RtdbSubscribeOptionAutoConn, "tag_sub_test")
	if !RteIsOk(err) {
		t.Error("创建属性订阅失败:", err)
	}
	if param == nil {
		t.Error("订阅参数不应为空")
	}
	err = RawRtdbbCancelSubscribeTagsWarp(handle, param)
	if !RteIsOk(err) {
		t.Error("取消属性订阅失败:", err)
	}
}

func TestRawRtdbsSubscribeSnapshotsEx64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无标签点")
	}
	count := 3
	if len(ids) < count {
		count = len(ids)
	}
	param, errs, err := RawRtdbsSubscribeSnapshotsEx64Warp(handle, ids[:count], RtdbSubscribeOptionAutoConn, "snap_sub_test")
	if !RteIsOk(err) {
		t.Error("创建快照订阅失败:", err)
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			t.Logf("第%d个点订阅出错: %v", i, e)
		}
	}
	if param != nil {
		_ = RawRtdbsCancelSubscribeSnapshotsWarp(handle, param)
	}
}

func TestRawRtdbsSubscribeDeltaSnapshots64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无标签点")
	}
	param, errs, err := RawRtdbsSubscribeDeltaSnapshots64Warp(handle, ids[:1], []float64{1.0}, []int64{0}, RtdbSubscribeOptionAutoConn, "delta_sub_test")
	if !RteIsOk(err) {
		t.Error("创建增量订阅失败:", err)
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			t.Logf("第%d个点增量订阅出错: %v", i, e)
		}
	}
	if param != nil {
		_ = RawRtdbsCancelSubscribeSnapshotsWarp(handle, param)
	}
}

func TestRawRtdbCreateDatagramHandleWarp(t *testing.T) {
	dh, err := RawRtdbCreateDatagramHandleWarp(0, "127.0.0.1")
	if !RteIsOk(err) {
		t.Logf("创建数据流(可能权限或端口占用): %v", err)
		return
	}
	defer RawRtdbRemoveDatagramHandleWarp(dh)
	fmt.Println("数据流句柄创建成功")
}

// ==================== 09. 存档管理 ====================

func TestRawRtdbaGetArchivesCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbaGetArchivesCountWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取存档数量失败:", err)
	}
	fmt.Println("存档数量:", count)
}

func TestRawRtdbaGetArchivesWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	paths, files, states, err := RawRtdbaGetArchivesWarp(handle, 100)
	if !RteIsOk(err) {
		t.Error("获取存档列表失败:", err)
	}
	fmt.Println("存档数量:", len(paths), len(files), len(states))
}

func TestRawRtdbaGetArchivesStatusWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	state, err := RawRtdbaGetArchivesStatusWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取存档状态失败:", err)
	}
	fmt.Println("存档状态:", state)
}

func TestRawRtdbaQueryBigJob64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	path, file, job, state, endTime, progress, err := RawRtdbaQueryBigJob64Warp(handle, RtdbProcessBase)
	if !RteIsOk(err) {
		t.Logf("查询后台任务: %v", err)
		return
	}
	fmt.Printf("后台任务: path=%s, file=%s, job=%d, state=%v, endTime=%d, progress=%f\n", path, file, job, state, endTime, progress)
}

func TestRawRtdbaCancelBigJobWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbaCancelBigJobWarp(handle, RtdbProcessBase)
	if !RteIsOk(err) {
		t.Logf("取消后台任务(可能无任务): %v", err)
	}
}

// ==================== 10. 历史数据查询 ====================

func TestRawRtdbhArchivedValuesCount64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	count, err := RawRtdbhArchivedValuesCount64Warp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) {
		t.Logf("统计历史值数量: %v", err)
		return
	}
	fmt.Println("历史值数量:", count)
}

func TestRawRtdbhGetArchivedValues64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, sts, vals, states, quals, err := RawRtdbhGetArchivedValues64Warp(handle, pid, 100, past, 0, now, 0)
	if !RteIsOk(err) {
		t.Logf("获取历史数据: %v", err)
		return
	}
	fmt.Println("历史数据条数:", len(dts), len(sts), len(vals), len(states), len(quals))
}

func TestRawRtdbhGetSingleValue64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	_, _, val, state, qual, err := RawRtdbhGetSingleValue64Warp(handle, pid, RtdbHisModePrevious, now, 0)
	if !RteIsOk(err) {
		t.Logf("获取单值历史: %v", err)
		return
	}
	fmt.Println("单值:", val, state, qual)
}

func TestRawRtdbhSummaryDataWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	data, err := RawRtdbhSummaryDataWarp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) {
		t.Logf("获取统计值: %v", err)
		return
	}
	fmt.Printf("统计: Count=%d, Max=%f, Min=%f\n", data.Count, data.MaxValue, data.MinValue)
}

func TestRawRtdbhGetPlotValues64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, _, _, _, _, err := RawRtdbhGetPlotValues64Warp(handle, pid, 100, past, 0, now, 0)
	if !RteIsOk(err) {
		t.Logf("获取绘图数据: %v", err)
		return
	}
	fmt.Println("绘图数据条数:", len(dts))
}

func TestRawRtdbhGetCrossSectionValues64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无标签点")
	}
	count := 3
	if len(ids) < count {
		count = len(ids)
	}
	now := TimestampType(time.Now().Unix())
	_, _, _, _, _, _, err := RawRtdbhGetCrossSectionValues64Warp(handle, ids[:count], RtdbHisModePrevious, now, 0)
	if !RteIsOk(err) {
		t.Logf("获取断面数据: %v", err)
		return
	}
	fmt.Println("断面数据获取成功")
}

// ==================== 11. 历史数据写入与修改 ====================

func TestRawRtdbhPutSingleValue64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhPutSingleValue64Warp(handle, pid, now, 0, 123.45, 0, 0)
	if !RteIsOk(err) {
		t.Logf("写入单值历史: %v", err)
	}
}

func TestRawRtdbhPutArchivedValues64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{99.9}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		t.Logf("批量写入历史: %v", err)
	}
}

func TestRawRtdbhUpdateValue64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhUpdateValue64Warp(handle, pid, now, 0, 77.7, 0, 0)
	if !RteIsOk(err) {
		t.Logf("更新历史值: %v", err)
	}
}

func TestRawRtdbhRemoveValue64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhRemoveValue64Warp(handle, pid, now, 0)
	if !RteIsOk(err) {
		t.Logf("删除历史值: %v", err)
	}
}

func TestRawRtdbhFlushArchivedValuesWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	pid := getFirstPointID(t, handle)
	count, err := RawRtdbhFlushArchivedValuesWarp(handle, pid)
	if !RteIsOk(err) {
		t.Logf("刷新历史缓存: %v", err)
		return
	}
	fmt.Println("刷新缓存条数:", count)
}

// ==================== 12. 方程式计算与性能监控 ====================

func TestRawRtdbeComputeHistory64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无计算点")
	}
	now := TimestampType(time.Now().Unix())
	past := now - 3600
	_, err := RawRtdbeComputeHistory64Warp(handle, ids[:1], 0, past, 0, now, 0)
	if !RteIsOk(err) {
		t.Logf("历史计算: %v", err)
	}
}

func TestRawRtdbbGetEquationByIdWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无计算点")
	}
	eq, err := RawRtdbbGetEquationByIdWarp(handle, ids[0])
	if !RteIsOk(err) {
		t.Logf("获取方程式: %v", err)
		return
	}
	fmt.Println("方程式长度:", len(eq))
}

func TestRawRtdbeGetEquationGraphCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无计算点")
	}
	count, err := RawRtdbeGetEquationGraphCountWarp(handle, ids[0], RtdbGraphFlagAll)
	if !RteIsOk(err) {
		t.Logf("获取拓扑数量: %v", err)
		return
	}
	fmt.Println("拓扑数量:", count)
}

func TestRawRtdbeGetEquationGraphDatasWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		t.Skip("无计算点")
	}
	count, _ := RawRtdbeGetEquationGraphCountWarp(handle, ids[0], RtdbGraphFlagAll)
	if count <= 0 {
		t.Skip("拓扑为空")
	}
	graph, err := RawRtdbeGetEquationGraphDatasWarp(handle, ids[0], RtdbGraphFlagAll, count)
	if !RteIsOk(err) {
		t.Logf("获取拓扑数据: %v", err)
		return
	}
	fmt.Println("拓扑数据条数:", len(graph))
}

func TestRawRtdbpGetPerfTagsCountWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbpGetPerfTagsCountWarp(handle)
	if !RteIsOk(err) {
		t.Error("获取性能点数量失败:", err)
	}
	fmt.Println("性能点数量:", count)
}

func TestRawRtdbpGetPerfTagsInfoWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	infos, errs, err := RawRtdbpGetPerfTagsInfoWarp(handle, []RtdbPerfTagID{PftCpuUsageOfLogger, PftMemBytesOfLogger})
	if !RteIsOk(err) {
		t.Error("获取性能点信息失败:", err)
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			t.Logf("第%d个性能点出错: %v", i, e)
		}
	}
	fmt.Println("性能点信息数:", len(infos))
}

func TestRawRtdbpGetPerfValues64Warp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	_, _, _, _, _, _, err := RawRtdbpGetPerfValues64Warp(handle, []RtdbPerfTagID{PftCpuUsageOfLogger, PftMemBytesOfLogger})
	if !RteIsOk(err) {
		t.Error("获取性能值失败:", err)
	}
}

// ==================== 13. 数据流与元数据同步 ====================

func TestRawRtdbbGetMetaSyncInfoWarp(t *testing.T) {
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	infos, errs, err := RawRtdbbGetMetaSyncInfoWarp(handle, 0)
	if !RteIsOk(err) {
		t.Logf("获取元数据同步信息: %v", err)
		return
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			t.Logf("第%d个节点出错: %v", i, e)
		}
	}
	fmt.Println("同步信息数:", len(infos))
}
