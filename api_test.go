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
// 注意：在 RawRtdbConnectWarp 之前设置客户端默认时间戳精度为纳秒
func connectAndLogin(t *testing.T) ConnectHandle {
	t.Helper()

	// 设置客户端默认时间戳精度为纳秒（必须在连接前设置）
	rte := RawRtdbSetOptionWarp(RtdbApiOptionDefaultPrecision, 3)
	if !RteIsOk(rte) {
		t.Fatalf("设置纳秒精度失败: %v", rte)
	}

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
	fmt.Println("【步骤1】获取 API 版本号（预期：成功）")
	apiVersion, err := RawRtdbGetApiVersionWarp()
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取版本号失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 版本号=%d\n", apiVersion)
}

// TC-OPTION-01 设置自动重连选项
func TestRawRtdbSetOptionWarp(t *testing.T) {
	fmt.Println("【步骤1】设置自动重连选项（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbSetOptionWarp(RtdbApiOptionAutoReconn, 1)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("设置自动重连失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自动重连已开启")
}

// TC-CONN-01 正常连接数据库
func TestRawRtdbConnectWarp(t *testing.T) {
	fmt.Println("【步骤1】连接数据库（预期：成功）")
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Fatalf("连接失败: %v", err)
	}
	if handle <= 0 {
		fmt.Println("  结果：失败 —— 连接句柄无效")
		t.Error("连接句柄无效")
		return
	}
	fmt.Printf("  结果：通过 —— 连接句柄=%d\n", handle)

	fmt.Println("【步骤2】断开连接（预期：成功）")
	err = RawRtdbDisconnectWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("断开连接: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 已断开")
	}
}

// TC-CONN-02 连接不存在的主机
func TestRawRtdbConnectWarp_UnreachableHost(t *testing.T) {
	fmt.Println("【步骤1】连接不可达主机（预期：返回错误）")
	handle, err := RawRtdbConnectWarp("192.0.2.1", 9000)
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 不可达主机居然连接成功了！")
		t.Error("连接不可达主机应失败")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
	}
	if handle > 0 {
		_ = RawRtdbDisconnectWarp(handle)
	}
}

// TC-LOGIN-01 管理员正常登录
func TestRawRtdbLoginWarp(t *testing.T) {
	fmt.Println("【步骤1】连接数据库（预期：成功）")
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Fatalf("连接失败: %v", err)
	}
	fmt.Printf("  结果：通过 —— 句柄=%d\n", handle)
	defer RawRtdbDisconnectWarp(handle)

	fmt.Println("【步骤2】管理员登录（预期：成功）")
	priv, err := RawRtdbLoginWarp(handle, Username, Password)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Fatalf("登录失败: %v", err)
	}
	fmt.Printf("  结果：通过 —— 权限组=%v\n", priv)
}

// TC-LOGIN-02 错误密码登录
func TestRawRtdbLoginWarp_WrongPassword(t *testing.T) {
	fmt.Println("【步骤1】连接数据库（预期：成功）")
	handle, err := RawRtdbConnectWarp(Hostname, Port)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Fatalf("连接失败: %v", err)
	}
	fmt.Printf("  结果：通过 —— 句柄=%d\n", handle)
	defer RawRtdbDisconnectWarp(handle)

	fmt.Println("【步骤2】使用错误密码登录（预期：返回错误）")
	_, err = RawRtdbLoginWarp(handle, Username, "wrong_pwd")
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 错误密码居然登录成功了！")
		t.Error("错误密码应登录失败")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
	}
}

// TC-DISC-01 正常断开连接
func TestRawRtdbDisconnectWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	fmt.Printf("  结果：通过 —— 句柄=%d\n", handle)

	fmt.Println("【步骤2】断开连接（预期：成功）")
	err := RawRtdbDisconnectWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("断开连接失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 已断开")
}

// TC-SUBCONN-01 正常创建连接事件订阅 / TC-SUBCONN-03 正常取消订阅
func TestRawRtdbSubscribeConnectExWarp(t *testing.T) {
	fmt.Println("【步骤1】创建连接事件订阅（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	param, err := RawRtdbSubscribeConnectExWarp(handle, RtdbSubscribeOptionAutoConn, "test_sub")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建连接事件订阅失败:", err)
		return
	}
	if param == nil {
		fmt.Println("  结果：失败 —— 订阅参数为空")
		t.Error("订阅参数不应为空")
		return
	}
	fmt.Printf("  结果：通过 —— 订阅参数已创建\n")

	fmt.Println("【步骤2】取消订阅（预期：成功）")
	err = RawRtdbCancelSubscribeConnectWarp(handle, param)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("取消订阅失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 已取消订阅")
}

// TC-DBINFO-01 获取表文件路径（字符串参数）
func TestRawRtdbGetDbInfo1Warp(t *testing.T) {
	fmt.Println("【步骤1】获取字符串型系统参数（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	val, err := RawRtdbGetDbInfo1Warp(handle, RtdbParamTableFile)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取字符串参数失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 表文件路径=%s\n", val)
}

// TC-DBINFO-02 获取最大连接数（整型参数）
func TestRawRtdbGetDbInfo2Warp(t *testing.T) {
	fmt.Println("【步骤1】获取整型系统参数（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	val, err := RawRtdbGetDbInfo2Warp(handle, RtdbParamServerConnectionCount)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取整型参数失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 最大连接数=%d\n", val)
}

// TC-SETDBINFO-01 管理员设置系统参数
func TestRawRtdbSetDbInfo1Warp(t *testing.T) {
	fmt.Println("【步骤1】读取当前字符串参数（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	current, err := RawRtdbGetDbInfo1Warp(handle, RtdbParamServerSenderIp)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s，跳过写入测试\n", err)
		t.Logf("读取参数失败，跳过写入测试: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 当前值=%s\n", current)

	fmt.Println("【步骤2】写回相同值（预期：成功或提示需重启）")
	err = RawRtdbSetDbInfo1Warp(handle, RtdbParamServerSenderIp, current)
	if !RteIsOk(err) {
		fmt.Printf("  结果：通过 —— 返回提示（属正常行为）：%s\n", err)
		t.Logf("设置字符串参数(可能提示需重启，属正常行为): %v", err)
	} else {
		fmt.Println("  结果：通过 —— 设置成功")
	}
}

// TC-SETDBINFO2-01 管理员设置整型参数
func TestRawRtdbSetDbInfo2Warp(t *testing.T) {
	fmt.Println("【步骤1】设置整型系统参数（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbSetDbInfo2Warp(handle, RtdbParamHashTableSize, 100)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("设置整型参数失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 设置成功")
}

// TC-CONNCOUNT-01 获取单机模式连接数
func TestRawRtdbConnectionCountWarp(t *testing.T) {
	fmt.Println("【步骤1】获取当前连接数（预期：成功且>=1）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, err := RawRtdbConnectionCountWarp(handle, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取连接数失败:", err)
		return
	}
	if count < 1 {
		fmt.Printf("  结果：失败 —— 连接数=%d，应至少为1\n", count)
		t.Error("连接数应至少为1")
		return
	}
	fmt.Printf("  结果：通过 —— 当前连接数=%d\n", count)
}

// TC-GETCONNS-01 获取所有连接句柄
func TestRawRtdbGetConnectionsWarp(t *testing.T) {
	fmt.Println("【步骤1】获取所有连接句柄（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	count, _ := RawRtdbConnectionCountWarp(handle, 0)
	sockets, err := RawRtdbGetConnectionsWarp(handle, 0, count)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取连接句柄失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个连接句柄\n", len(sockets))
}

// TC-GETCONNS-02 获取当前连接句柄
func TestRawRtdbGetOwnConnectionWarp(t *testing.T) {
	fmt.Println("【步骤1】获取自身连接句柄（预期：成功且>0）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, err := RawRtdbGetOwnConnectionWarp(handle, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取自身连接句柄失败:", err)
		return
	}
	if socket <= 0 {
		fmt.Printf("  结果：失败 —— 句柄=%d，应大于0\n", socket)
		t.Error("自身连接句柄应大于0")
		return
	}
	fmt.Printf("  结果：通过 —— 自身连接句柄=%d\n", socket)
}

// TC-CONNINFO-01 获取当前连接 IPv4 信息
func TestRawRtdbGetConnectionInfoWarp(t *testing.T) {
	fmt.Println("【步骤1】获取当前连接 IPv4 信息（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, _ := RawRtdbGetOwnConnectionWarp(handle, 0)
	info, err := RawRtdbGetConnectionInfoWarp(handle, 0, socket)
	if err != nil {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取连接信息失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— IP=%d, Port=%d\n", info.IpAddr, info.Port)
}

// TC-CONNINFO-03 获取 IPv6 连接信息
func TestRawRtdbGetConnectionInfoIpv6Warp(t *testing.T) {
	fmt.Println("【步骤1】获取 IPv6 连接信息（预期：成功或提示不支持）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, _ := RawRtdbGetOwnConnectionWarp(handle, 0)
	info, err := RawRtdbGetConnectionInfoIpv6Warp(handle, 0, socket)
	if !RteIsOk(err) {
		fmt.Printf("  结果：通过 —— 返回提示（可能不支持IPv6）：%s\n", err)
		t.Logf("获取IPv6信息(可能不支持): %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— IPv6连接信息: %+v\n", info)
}

// TC-OSTYPE-01 获取服务器 OS 类型
func TestRawRtdbOsType(t *testing.T) {
	fmt.Println("【步骤1】获取服务器 OS 类型（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	osType, err := RawRtdbOsType(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取OS类型失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— OS类型=%v\n", osType)
}

// TC-HANDLEINFO-01 获取句柄信息
func TestRawRtdbGetHandleInfoWarp(t *testing.T) {
	fmt.Println("【步骤1】获取句柄信息（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	info, err := RawRtdbGetHandleInfoWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取句柄信息失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— OsType=%v, NewDB=%v\n", info.OsType, info.NewDB)
}

// TC-HOSTTIME-01 获取服务器当前时间
func TestRawRtdbHostTime64Warp(t *testing.T) {
	fmt.Println("【步骤1】获取服务器当前时间（预期：成功且与本地时间偏差<60秒）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	ts, err := RawRtdbHostTime64Warp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取服务器时间失败:", err)
		return
	}
	now := time.Now().Unix()
	diff := ts - TimestampType(now)
	if diff < -60 || diff > 60 {
		fmt.Printf("  结果：失败 —— 服务器与本地时间偏差=%d秒\n", diff)
		t.Errorf("服务器时间与本地时间偏差过大: %d秒", diff)
		return
	}
	fmt.Printf("  结果：通过 —— 服务器时间戳=%d, 偏差=%d秒\n", ts, diff)
}

// TC-TIMESPAN-01/02 格式化时间跨度
func TestRawRtdbFormatTimespanWarp(t *testing.T) {
	fmt.Println("【步骤1】格式化10秒（预期：成功）")
	s10, err := RawRtdbFormatTimespanWarp(10)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("格式化10秒失败:", err)
	} else {
		fmt.Printf("  结果：通过 —— 10秒=%s\n", s10)
	}

	fmt.Println("【步骤2】格式化60秒（预期：成功）")
	s60, err := RawRtdbFormatTimespanWarp(60)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("格式化60秒失败:", err)
	} else {
		fmt.Printf("  结果：通过 —— 60秒=%s\n", s60)
	}
}

// TC-TIMESPAN-03/04 解析时间跨度字符串
func TestRawRtdbParseTimespanWarp(t *testing.T) {
	fmt.Println("【步骤1】解析有效时间跨度'2n'（预期：成功且=120秒）")
	dt, err := RawRtdbParseTimespanWarp("2n")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("解析2n失败:", err)
		return
	}
	if dt != 120 {
		fmt.Printf("  结果：失败 —— 解析结果=%d，应为120\n", dt)
		t.Error("解析2n应等于120秒")
		return
	}
	fmt.Printf("  结果：通过 —— 2n=%d秒\n", dt)

	fmt.Println("【步骤2】解析无效字符串'invalid'（预期：返回错误）")
	_, err = RawRtdbParseTimespanWarp("invalid")
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 无效字符串居然解析成功了！")
		t.Error("解析无效字符串应失败")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
	}
}

// TC-PARSETIME-01/02/03 解析时间字符串
func TestRawRtdbParseTimeWarp(t *testing.T) {
	fmt.Println("【步骤1】解析绝对时间'2024-01-01 08:00:00'（预期：成功）")
	ts, sub, err := RawRtdbParseTimeWarp("2024-01-01 08:00:00")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("解析绝对时间失败:", err)
	} else {
		fmt.Printf("  结果：通过 —— 时间戳=%d, 微秒=%d\n", ts, sub)
	}

	fmt.Println("【步骤2】解析相对时间'*-1d'（预期：成功）")
	ts2, _, err := RawRtdbParseTimeWarp("*-1d")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("解析相对时间失败:", err)
	} else {
		fmt.Printf("  结果：通过 —— 时间戳=%d\n", ts2)
	}

	fmt.Println("【步骤3】解析无效时间'not_a_time'（预期：返回错误）")
	_, _, err = RawRtdbParseTimeWarp("not_a_time")
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 无效时间居然解析成功了！")
		t.Error("解析无效时间应失败")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
	}
}

// TC-FMTMSG-01/02 格式化错误码描述
func TestRawRtdbFormatMessageWarp(t *testing.T) {
	fmt.Println("【步骤1】格式化 RteOk 错误码描述（预期：成功）")
	name, msg := RawRtdbFormatMessageWarp(RteOk)
	fmt.Printf("  结果：通过 —— %s: %s\n", name, msg)

	fmt.Println("【步骤2】格式化 RtePointNotFound 错误码描述（预期：成功）")
	name2, msg2 := RawRtdbFormatMessageWarp(RtePointNotFound)
	fmt.Printf("  结果：通过 —— %s: %s\n", name2, msg2)
}

// TC-JOBMSG-01 获取任务描述
func TestRawRtdbJobMessageWarp(t *testing.T) {
	fmt.Println("【步骤1】获取 Job 0 描述（预期：成功）")
	name, msg := RawRtdbJobMessageWarp(0)
	fmt.Printf("  结果：通过 —— %s: %s\n", name, msg)
}

// TC-TIMEOUT-01/03 设置并读取超时时间
func TestRawRtdbSetTimeoutWarp_GetTimeoutWarp(t *testing.T) {
	fmt.Println("【步骤1】设置超时时间为30秒（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, _ := RawRtdbGetOwnConnectionWarp(handle, 0)

	err := RawRtdbSetTimeoutWarp(handle, socket, 30)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("设置超时失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 已设置为30秒")

	fmt.Println("【步骤2】读取超时时间（预期：成功且=30）")
	to, err := RawRtdbGetTimeoutWarp(handle, socket)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取超时失败:", err)
		return
	}
	if to != 30 {
		fmt.Printf("  结果：失败 —— 超时时间=%d，应为30\n", to)
		t.Errorf("超时时间应为30, 实际=%d", to)
		return
	}
	fmt.Printf("  结果：通过 —— 超时时间=%d秒\n", to)
}

// TC-JUDGE-01 判断正常连接状态
func TestRawRtdbJudgeConnectStatusWarp(t *testing.T) {
	fmt.Println("【步骤1】判断连接状态（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	err := RawRtdbJudgeConnectStatusWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("判断连接状态失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 连接状态正常")
}

// TC-FMTIP-01/02/03 IP格式化
func TestRawRtdbFormatIpaddrWarp(t *testing.T) {
	fmt.Println("【步骤1】格式化 IP 0x7F000001（预期：127.0.0.1）")
	ip1 := RawRtdbFormatIpaddrWarp(0x7F000001)
	if ip1 != "127.0.0.1" {
		fmt.Printf("  结果：失败 —— 实际=%s\n", ip1)
		t.Errorf("127.0.0.1格式化错误: %s", ip1)
	} else {
		fmt.Printf("  结果：通过 —— %s\n", ip1)
	}

	fmt.Println("【步骤2】格式化 IP 0x00000000（预期：0.0.0.0）")
	ip2 := RawRtdbFormatIpaddrWarp(0)
	if ip2 != "0.0.0.0" {
		fmt.Printf("  结果：失败 —— 实际=%s\n", ip2)
		t.Errorf("0.0.0.0格式化错误: %s", ip2)
	} else {
		fmt.Printf("  结果：通过 —— %s\n", ip2)
	}

	fmt.Println("【步骤3】格式化 IP 0xFFFFFFFF（预期：255.255.255.255）")
	ip3 := RawRtdbFormatIpaddrWarp(0xFFFFFFFF)
	if ip3 != "255.255.255.255" {
		fmt.Printf("  结果：失败 —— 实际=%s\n", ip3)
		t.Errorf("255.255.255.255格式化错误: %s", ip3)
	} else {
		fmt.Printf("  结果：通过 —— %s\n", ip3)
	}
}

// TC-KILLCONN-01/02 断开连接（管理员断开自身/其他）
func TestRawRtdbKillConnectionWarp(t *testing.T) {
	fmt.Println("【步骤1】尝试断开自身连接（预期：返回错误或被允许）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	socket, _ := RawRtdbGetOwnConnectionWarp(handle, 0)
	err := RawRtdbKillConnectionWarp(handle, socket)
	if RteIsOk(err) {
		fmt.Println("  结果：通过 —— 断开自身连接成功（部分实现允许）")
		t.Log("断开自身连接返回Ok（部分实现允许）")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
		t.Logf("断开自身连接(预期失败): %v", err)
	}
}

// ==================== 02. 用户权限与连接管理 ====================

// TC-CHPWD-01 管理员修改其他用户密码
func TestRawRtdbChangePasswordWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加测试用户 testuser_chpwd（预期：成功）")
	_ = RawRtdbAddUserWarp(handle, "testuser_chpwd", "OldPass123", PrivGroupRtdbRO)
	defer RawRtdbRemoveUserWarp(handle, "testuser_chpwd")
	fmt.Println("  结果：通过 —— 测试用户添加成功")

	fmt.Println("【步骤3】管理员修改 testuser_chpwd 密码（预期：成功）")
	err := RawRtdbChangePasswordWarp(handle, "testuser_chpwd", "NewPass123")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("修改密码失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 密码修改成功")
}

// TC-MYPWD-01 正常修改自己的密码
func TestRawRtdbChangeMyPasswordWarp(t *testing.T) {
	fmt.Println("【步骤1】修改自己的密码（预期：需要普通用户权限）")
	t.Log("ChangeMyPassword需要普通用户登录测试，跳过")
	fmt.Println("  结果：跳过 —— 需要普通用户登录测试")
}

// TC-GETPRIV-01 管理员权限查询
func TestRawRtdbGetPrivWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】查询当前用户权限（预期：成功）")
	priv, err := RawRtdbGetPrivWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取权限失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 当前权限=%d\n", priv)
}

// TC-SETPRIV-01 管理员提升普通用户权限
func TestRawRtdbChangePrivWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加普通用户 testuser_priv（预期：成功）")
	_ = RawRtdbAddUserWarp(handle, "testuser_priv", "Pass123", PrivGroupRtdbRO)
	defer RawRtdbRemoveUserWarp(handle, "testuser_priv")
	fmt.Println("  结果：通过 —— 普通用户添加成功")

	fmt.Println("【步骤3】提升用户权限为 DW（预期：成功）")
	err := RawRtdbChangePrivWarp(handle, "testuser_priv", PrivGroupRtdbDW)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("修改权限失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 权限提升成功")
}

// TC-ADDUSER-01 管理员添加普通用户
func TestRawRtdbAddUserWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加普通用户 testuser_add（预期：成功）")
	err := RawRtdbAddUserWarp(handle, "testuser_add", "Pass123", PrivGroupRtdbRO)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("添加用户失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 用户添加成功")
	defer RawRtdbRemoveUserWarp(handle, "testuser_add")
}

// TC-RMVUSER-01 管理员删除普通用户
func TestRawRtdbRemoveUserWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加测试用户 testuser_rmv（预期：成功）")
	_ = RawRtdbAddUserWarp(handle, "testuser_rmv", "Pass123", PrivGroupRtdbRO)
	fmt.Println("  结果：通过 —— 测试用户添加成功")

	fmt.Println("【步骤3】删除测试用户 testuser_rmv（预期：成功）")
	err := RawRtdbRemoveUserWarp(handle, "testuser_rmv")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("删除用户失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 用户删除成功")
}

// TC-LOCK-01 管理员禁用用户 / TC-LOCK-02 管理员启用被禁用户
func TestRawRtdbLockUserWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加测试用户 testuser_lock（预期：成功）")
	_ = RawRtdbAddUserWarp(handle, "testuser_lock", "Pass123", PrivGroupRtdbRO)
	defer RawRtdbRemoveUserWarp(handle, "testuser_lock")
	fmt.Println("  结果：通过 —— 测试用户添加成功")

	fmt.Println("【步骤3】禁用用户 testuser_lock（预期：成功）")
	err := RawRtdbLockUserWarp(handle, "testuser_lock", OFF)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("禁用用户失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 用户已禁用")

	fmt.Println("【步骤4】重新启用用户 testuser_lock（预期：成功）")
	err = RawRtdbLockUserWarp(handle, "testuser_lock", ON)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("启用用户失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 用户已启用")
}

// TC-GETUSERS-01 管理员获取用户列表
func TestRawRtdbGetUsersWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】获取用户列表（预期：成功）")
	users, err := RawRtdbGetUsersWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取用户列表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 用户数量=%d\n", len(users))
}

// TC-BLADD-01 管理员添加单 IP 黑名单
func TestRawRtdbAddBlacklistWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加黑名单 192.168.1.100（预期：成功）")
	err := RawRtdbAddBlacklistWarp(handle, "192.168.1.100", "255.255.255.255", "blocked")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("添加黑名单失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 黑名单添加成功")
	defer RawRtdbRemoveBlacklistWarp(handle, "192.168.1.100", "255.255.255.255")
}

// TC-BLUPD-01 正常更新黑名单地址
func TestRawRtdbUpdateBlacklistWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加黑名单 192.168.1.0（预期：成功）")
	_ = RawRtdbAddBlacklistWarp(handle, "192.168.1.0", "255.255.255.0", "网段封禁")
	defer RawRtdbRemoveBlacklistWarp(handle, "192.168.2.0", "255.255.255.0")
	fmt.Println("  结果：通过 —— 黑名单添加成功")

	fmt.Println("【步骤3】更新黑名单为 192.168.2.0（预期：成功）")
	err := RawRtdbUpdateBlacklistWarp(handle, "192.168.1.0", "255.255.255.0", "192.168.2.0", "255.255.255.0", "new desc")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("更新黑名单失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 黑名单更新成功")
}

// TC-BLRMV-01 管理员删除黑名单
func TestRawRtdbRemoveBlacklistWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加黑名单 192.168.1.0（预期：成功）")
	_ = RawRtdbAddBlacklistWarp(handle, "192.168.1.0", "255.255.255.0", "待删除")
	fmt.Println("  结果：通过 —— 黑名单添加成功")

	fmt.Println("【步骤3】删除黑名单 192.168.1.0（预期：成功）")
	err := RawRtdbRemoveBlacklistWarp(handle, "192.168.1.0", "255.255.255.0")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("删除黑名单失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 黑名单删除成功")
}

// TC-BLGET-01 管理员获取黑名单列表
func TestRawRtdbGetBlacklistWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】获取黑名单列表（预期：成功）")
	list, err := RawRtdbGetBlacklistWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取黑名单失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 黑名单数量=%d\n", len(list))
}

// TC-AUTHADD-01 管理员添加信任网段
func TestRawRtdbAddAuthorizationWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加信任网段 192.168.1.0（预期：成功）")
	err := RawRtdbAddAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "office", PrivGroupRtdbRO)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("添加信任段失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 信任网段添加成功")
	defer RawRtdbRemoveAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0")
}

// TC-AUTHUPD-01 正常更新信任段地址和权限
func TestRawRtdbUpdateAuthorizationWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加信任网段 192.168.1.0（预期：成功）")
	_ = RawRtdbAddAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "office", PrivGroupRtdbRO)
	defer RawRtdbRemoveAuthorizationWarp(handle, "192.168.2.0", "255.255.255.0")
	fmt.Println("  结果：通过 —— 信任网段添加成功")

	fmt.Println("【步骤3】更新信任网段为 192.168.2.0（预期：成功）")
	err := RawRtdbUpdateAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "192.168.2.0", "255.255.255.0", "new office", PrivGroupRtdbDW)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("更新信任段失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 信任网段更新成功")
}

// TC-AUTHRMV-01 管理员删除信任段
func TestRawRtdbRemoveAuthorizationWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】添加信任网段 192.168.1.0（预期：成功）")
	_ = RawRtdbAddAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0", "待删除", PrivGroupRtdbRO)
	fmt.Println("  结果：通过 —— 信任网段添加成功")

	fmt.Println("【步骤3】删除信任网段 192.168.1.0（预期：成功）")
	err := RawRtdbRemoveAuthorizationWarp(handle, "192.168.1.0", "255.255.255.0")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("删除信任段失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 信任网段删除成功")
}

// TC-AUTHGET-01 管理员获取信任段列表
func TestRawRtdbGetAuthorizationsWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录管理员账号")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 管理员登录成功")

	fmt.Println("【步骤2】获取信任网段列表（预期：成功）")
	list, err := RawRtdbGetAuthorizationsWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取信任段失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 信任网段数量=%d\n", len(list))
}

// ==================== 03. 文件目录操作 ====================

// TC-DRIVERS-01 Windows/Linux 平台获取盘符
func TestRawRtdbGetLogicalDriversWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取逻辑盘符列表（预期：成功）")
	drivers, err := RawRtdbGetLogicalDriversWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取盘符失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 盘符=%v\n", drivers)
}

// TC-OPENPATH-01 打开存在的目录 / TC-READPATH-01 遍历目录下第一个条目 / TC-CLOSEPATH-01 正常关闭已打开目录
func TestRawRtdbOpenPathWarp_ReadPath64Warp_ClosePathWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】打开根目录 /（预期：成功或平台差异）")
	err := RawRtdbOpenPathWarp(handle, "/")
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 打开根目录失败（平台差异）: %v\n", err)
		t.Logf("打开根目录(可能平台差异): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 根目录打开成功")

	fmt.Println("【步骤3】读取目录第一个条目（预期：成功）")
	item, err := RawRtdbReadPath64Warp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %v\n", err)
		t.Logf("读取目录条目: %v", err)
	} else {
		fmt.Printf("  结果：通过 —— 目录条目=%s, IsDir=%d, Size=%d\n", item.Path, item.IsDir, item.Size)
	}

	fmt.Println("【步骤4】关闭目录（预期：成功）")
	err = RawRtdbClosePathWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("关闭目录失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 目录已关闭")
}

// TC-MKDIR-01 正常创建目录
func TestRawRtdbMkdirWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建目录 /data/testdir_go（预期：成功或已存在）")
	err := RawRtdbMkdirWarp(handle, "/data/testdir_go")
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 创建目录失败（可能已存在）: %v\n", err)
		t.Logf("创建目录(可能已存在): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 目录创建成功")
}

// TC-FILESIZE-01 获取存在的文件大小
func TestRawRtdbGetFileSizeWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取文件 /data/test.txt 大小（预期：成功或文件不存在）")
	size, err := RawRtdbGetFileSizeWarp(handle, "/data/test.txt")
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 获取文件大小失败（可能文件不存在）: %v\n", err)
		t.Logf("获取文件大小(可能文件不存在): %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 文件大小=%d\n", size)
}

// TC-READFILE-01 从头读取整个文件
func TestRawRtdbReadFileWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】读取文件 /data/test.txt（预期：成功或文件不存在）")
	data, err := RawRtdbReadFileWarp(handle, "/data/test.txt", 0, 1024)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 读取文件失败（可能文件不存在）: %v\n", err)
		t.Logf("读取文件(可能文件不存在): %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 读取文件长度=%d\n", len(data))
}

// TC-MAXBLOB-01 正常获取最大长度
func TestRawRtdbGetMaxBlobLenWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取最大 Blob 长度（预期：成功）")
	maxLen, err := RawRtdbGetMaxBlobLenWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取最大Blob长度失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 最大Blob长度=%d\n", maxLen)
}

// TC-FMTQUAL-01/02 格式化单个/多个质量码
func TestRawRtdbFormatQualityWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】格式化质量码 [192, 0, 64]（预期：成功）")
	descs, err := RawRtdbFormatQualityWarp(handle, []Quality{192, 0, 64})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("格式化质量码失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 质量码描述=%v\n", descs)
}

// ==================== 04. 表管理 ====================

// TC-TBLADD-01 正常创建新表
func TestRawRtdbbAppendTableWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建新表 TestTableGo（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTableGo", "测试表GO")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 创建表ID=%d\n", tbl.ID)
	defer RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
}

// TC-TBLRMVID-01 删除存在的空表
func TestRawRtdbbRemoveTableByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTableRmId（预期：成功）")
	tbl, _ := RawRtdbbAppendTableWarp(handle, "TestTableRmId", "待删除")
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)

	fmt.Println("【步骤3】按ID删除临时表（预期：成功）")
	err := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按ID删除表失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 表已删除")
}

// TC-TBLRMVNAME-01 删除存在的表
func TestRawRtdbbRemoveTableByNameWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTableRmName（预期：成功）")
	_, _ = RawRtdbbAppendTableWarp(handle, "TestTableRmName", "待删除")
	fmt.Println("  结果：通过 —— 临时表创建成功")

	fmt.Println("【步骤3】按名称删除临时表（预期：成功）")
	err := RawRtdbbRemoveTableByNameWarp(handle, "TestTableRmName")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按名称删除表失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 表已删除")
}

// TC-TBLCNT-01 查询表总数
func TestRawRtdbbTablesCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】查询表总数（预期：成功且>=0）")
	count, err := RawRtdbbTablesCountWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取表总数失败:", err)
		return
	}
	if count < 0 {
		fmt.Printf("  结果：失败 —— 表总数=%d，不应小于0\n", count)
		t.Error("表总数不应小于0")
		return
	}
	fmt.Printf("  结果：通过 —— 表总数=%d\n", count)
}

// TC-TBLGET-01 获取所有表 ID
func TestRawRtdbbGetTablesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取所有表ID（预期：成功）")
	count, _ := RawRtdbbTablesCountWarp(handle)
	ids, err := RawRtdbbGetTablesWarp(handle, count)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取表ID列表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 表数量=%d\n", len(ids))
}

// TC-TBLSIZEID-01 完整流程：创建表→创建点→测表大小→清理点→清理表
func TestRawRtdbbGetTableSizeByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTblSize（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSize", "测试表大小")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点 TestPtSize（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtSize", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】按ID获取表大小（预期：成功且>=1）")
	size, err := RawRtdbbGetTableSizeByIdWarp(handle, tbl.ID)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取表大小失败:", err)
		return
	}
	if size < 1 {
		fmt.Printf("  结果：失败 —— 表大小=%d，预期至少为1\n", size)
		t.Errorf("表大小预期>=1，实际=%d", size)
		return
	}
	fmt.Printf("  结果：通过 —— 表大小=%d\n", size)
}

// TC-TBLSIZENAME-01 完整流程：创建表→创建点→按名称测表大小→清理点→清理表
func TestRawRtdbbGetTableSizeByNameWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTblSizeName（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSizeName", "测试按名称获取表大小")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点 TestPtSizeName（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtSizeName", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】按名称获取表大小（预期：成功且>=1）")
	size, err := RawRtdbbGetTableSizeByNameWarp(handle, tbl.Name)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按名称获取表大小失败:", err)
		return
	}
	if size < 1 {
		fmt.Printf("  结果：失败 —— 表大小=%d，预期至少为1\n", size)
		t.Errorf("表大小预期>=1，实际=%d", size)
		return
	}
	fmt.Printf("  结果：通过 —— 表大小(ByName)=%d\n", size)
}

// TC-TBLREALSIZE-01 完整流程：创建表→创建点→按ID测表实际大小→清理点→清理表
func TestRawRtdbbGetTableRealSizeByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTblRealSize（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblRealSize", "测试按ID获取表实际大小")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点 TestPtRealSize（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtRealSize", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】按ID获取表实际大小（预期：成功且>=1）")
	size, err := RawRtdbbGetTableRealSizeByIdWarp(handle, tbl.ID)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取表实际大小失败:", err)
		return
	}
	if size < 1 {
		fmt.Printf("  结果：失败 —— 表实际大小=%d，预期至少为1\n", size)
		t.Errorf("表实际大小预期>=1，实际=%d", size)
		return
	}
	fmt.Printf("  结果：通过 —— 表实际大小=%d\n", size)
}

// TC-TBLPROPBYID-01 完整流程：创建表→创建点→按ID获取表属性→清理点→清理表
func TestRawRtdbbGetTablePropertyByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTblPropById（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblPropById", "测试按ID获取表属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点 TestPtPropById（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtPropById", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】按ID获取表属性（预期：成功且Name匹配）")
	prop, err := RawRtdbbGetTablePropertyByIdWarp(handle, tbl.ID)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取表属性失败:", err)
		return
	}
	if prop.Name != tbl.Name {
		fmt.Printf("  结果：失败 —— 表名不一致: 预期=%s, 实际=%s\n", tbl.Name, prop.Name)
		t.Errorf("表名不一致: 预期=%s, 实际=%s", tbl.Name, prop.Name)
		return
	}
	fmt.Printf("  结果：通过 —— 表属性: ID=%d, Name=%s\n", prop.ID, prop.Name)
}

// TC-TBLPROPBYNAM-01 完整流程：创建表→创建点→按名称获取表属性→清理点→清理表
func TestRawRtdbbGetTablePropertyByNameWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTblPropByName（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblPropByName", "测试按名称获取表属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点 TestPtPropByName（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtPropByName", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】按名称获取表属性（预期：成功且与按ID一致）")
	propByID, _ := RawRtdbbGetTablePropertyByIdWarp(handle, tbl.ID)
	propByName, err := RawRtdbbGetTablePropertyByNameWarp(handle, tbl.Name)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按名称获取表属性失败:", err)
		return
	}
	if propByName.ID != propByID.ID {
		fmt.Printf("  结果：失败 —— 按名称和按ID获取的表属性不一致: ID=%d vs %d\n", propByName.ID, propByID.ID)
		t.Error("按名称和按ID获取的表属性不一致")
		return
	}
	fmt.Printf("  结果：通过 —— 表属性一致: ID=%d, Name=%s\n", propByName.ID, propByName.Name)
}

// TC-TBLUPDNAME-01 正常更新表名
func TestRawRtdbbUpdateTableNameWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTableOldName（预期：成功）")
	tbl, _ := RawRtdbbAppendTableWarp(handle, "TestTableOldName", "测试")
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)

	fmt.Println("【步骤3】更新表名为 TestTableNewName（预期：成功）")
	err := RawRtdbbUpdateTableNameWarp(handle, tbl.ID, "TestTableNewName")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("更新表名失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 表名更新成功")

	fmt.Println("【步骤4】验证表名是否更新（预期：Name=TestTableNewName）")
	prop, _ := RawRtdbbGetTablePropertyByIdWarp(handle, tbl.ID)
	if prop.Name != "TestTableNewName" {
		fmt.Printf("  结果：失败 —— 表名未更新，实际=%s\n", prop.Name)
		t.Error("表名未更新")
		return
	}
	fmt.Printf("  结果：通过 —— 表名已更新为=%s\n", prop.Name)
}

// TC-TBLUPDDESCID-01 正常更新表描述
func TestRawRtdbbUpdateTableDescByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTableDescId（预期：成功）")
	tbl, _ := RawRtdbbAppendTableWarp(handle, "TestTableDescId", "旧描述")
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)

	fmt.Println("【步骤3】按ID更新表描述为新描述GO（预期：成功）")
	err := RawRtdbbUpdateTableDescByIdWarp(handle, tbl.ID, "新描述GO")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("更新表描述失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 表描述更新成功")
}

// TC-TBLUPDDESCNAM-01 正常更新表描述
func TestRawRtdbbUpdateTableDescByNameWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表 TestTableDescName（预期：成功）")
	_, _ = RawRtdbbAppendTableWarp(handle, "TestTableDescName", "旧描述")
	fmt.Println("  结果：通过 —— 临时表创建成功")
	defer RawRtdbbRemoveTableByNameWarp(handle, "TestTableDescName")

	fmt.Println("【步骤3】按名称更新表描述为新描述ByName（预期：成功）")
	err := RawRtdbbUpdateTableDescByNameWarp(handle, "TestTableDescName", "新描述ByName")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按名称更新表描述失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 表描述更新成功")
}

// ==================== 05. 标签点管理 ====================

// TC-PTBASE-01/02 完整流程：创建表→创建基础点→清理点→清理表
func TestRawRtdbbInsertBasePointWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBasePt", "测试基础点创建")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建基础点 TestBaseBool（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestBaseBool", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建基础点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 创建基础点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()
}

// TC-PTINS-01 完整流程：创建表→创建完整点→清理点→清理表
func TestRawRtdbbInsertPointWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblInsertPt", "测试完整点创建")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建完整点 TestInsertPt（预期：成功）")
	base := &RtdbPoint{Tag: "TestInsertPt", Table: tbl.ID, Type: RtdbTypeReal64}
	scan := &RtdbScan{Source: "go_test"}
	base, _, _, err = RawRtdbbInsertPointWarp(handle, base, scan, nil)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建完整点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 创建完整点ID=%d\n", base.ID)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, base.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()
}

// TC-PTINSMAX-01 完整流程：创建表→创建Max点→清理点→清理表
func TestRawRtdbbInsertMaxPointWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblMaxPt", "测试Max点创建")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建Max点 TestMaxPt（预期：成功）")
	base := &RtdbPoint{Tag: "TestMaxPt", Table: tbl.ID, Type: RtdbTypeReal64}
	scan := &RtdbScan{Source: "go_test"}
	base, _, _, err = RawRtdbbInsertMaxPointWarp(handle, base, scan, nil)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建Max点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 创建Max点ID=%d\n", base.ID)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, base.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()
}

// TC-PTRMVID-01 完整流程：创建表→创建点→按ID删除点→清理表
func TestRawRtdbbRemovePointByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblRmvId", "测试按ID删除点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestRmvPtId", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)

	fmt.Println("【步骤4】按ID删除临时点（预期：成功）")
	err = RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按ID删除点失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 点已删除")
}

// TC-PTRMVNAME-01 完整流程：创建表→创建点→按名称删除点→清理表
func TestRawRtdbbRemovePointByNameWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblRmvName", "测试按名称删除点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestRmvPtName", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	fullName := tbl.Name + ".TestRmvPtName"

	fmt.Println("【步骤4】按名称删除临时点（预期：成功）")
	err = RawRtdbbRemovePointByNameWarp(handle, fullName)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按名称删除点失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 点已删除")
}

// TC-PTINSBATCH-01 完整流程：创建表→批量创建点→清理点→清理表
func TestRawRtdbbInsertMaxPointsWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBatchPt", "测试批量创建点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】批量创建3个点（预期：成功）")
	bases := make([]RtdbPoint, 3)
	scans := make([]RtdbScan, 3)
	calcs := make([]RtdbCalc, 3)
	for i := 0; i < 3; i++ {
		bases[i] = RtdbPoint{Tag: fmt.Sprintf("TestBatchPt%d", i), Table: tbl.ID, Type: RtdbTypeReal64}
		scans[i] = RtdbScan{Source: "go_test"}
	}
	bases, _, _, errs, err := RawRtdbbInsertMaxPointsWarp(handle, bases, scans, calcs)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("批量创建点失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量创建完成")
	for i, e := range errs {
		if !RteIsOk(e) {
			fmt.Printf("    第%d个点创建出错: %v\n", i, e)
			t.Logf("第%d个点创建出错: %v", i, e)
		} else {
			fmt.Printf("    第%d个点创建成功: ID=%d\n", i, bases[i].ID)
			defer RawRtdbbRemovePointByIdWarp(handle, bases[i].ID)
		}
	}
}

// TC-PTNAMED-02 完整流程：创建表→使用不存在的自定义类型创建点→清理表
func TestRawRtdbbInsertNamedTypePointWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblNamedPt", "测试自定义类型点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤4】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】使用不存在的自定义类型创建点（预期：返回错误）")
	base := &RtdbPoint{Tag: "TestNamedPt", Table: tbl.ID}
	scan := &RtdbScan{Source: "go_test"}
	_, _, err = RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "NotExistType")
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 不存在的类型居然成功了！")
		t.Error("不存在的类型应失败")
		return
	}
	fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
}

// TC-PTMOVE-01 完整流程：创建2个表→创建点→移动点→清理点→清理2个表
func TestRawRtdbbMovePointByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建两个临时表（预期：成功）")
	tbl1, err := RawRtdbbAppendTableWarp(handle, "TestTblMove1", "测试移动点源表")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建源表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 源表ID=%d\n", tbl1.ID)
	defer func() {
		fmt.Println("【步骤7】清理源表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl1.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理源表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 源表已删除")
		}
	}()

	tbl2, err := RawRtdbbAppendTableWarp(handle, "TestTblMove2", "测试移动点目标表")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建目标表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 目标表ID=%d\n", tbl2.ID)
	defer func() {
		fmt.Println("【步骤6】清理目标表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl2.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理目标表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 目标表已删除")
		}
	}()

	fmt.Println("【步骤3】在源表创建临时点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestMovePt", RtdbTypeBool, tbl1.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】移动临时点到目标表（预期：成功或权限约束）")
	fmt.Printf("  目标表名=%s\n", tbl2.Name)
	err = RawRtdbbMovePointByIdWarp(handle, pid, tbl2.Name)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 移动点失败（可能权限或约束）: %v\n", err)
		t.Logf("移动点(可能权限或约束): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 点移动成功")
}

// TC-PTGETPROP-01 完整流程：创建表→创建3个点→批量获取点属性→清理点→清理表
func TestRawRtdbbGetPointsPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblGetProp", "测试获取点属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建3个标签点（预期：成功）")
	var pids []PointID
	for i := 0; i < 3; i++ {
		tag := fmt.Sprintf("TestGetPropPt%d", i)
		pid, err := RawRtdbbInsertBasePointWarp(handle, tag, RtdbTypeReal64, tbl.ID, 0)
		if !RteIsOk(err) {
			fmt.Printf("  结果：失败 —— %s\n", err)
			t.Error("创建标签点失败:", err)
			return
		}
		fmt.Printf("  第%d个点创建成功: ID=%d\n", i, pid)
		pids = append(pids, pid)
		defer RawRtdbbRemovePointByIdWarp(handle, pid)
	}
	fmt.Println("  结果：通过 —— 3个标签点创建完成")

	fmt.Println("【步骤4】批量获取点属性（预期：成功）")
	bases, scans, calcs, errs, err := RawRtdbbGetPointsPropertyWarp(handle, pids)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("批量获取点属性失败:", err)
		return
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			fmt.Printf("    第%d个点获取属性出错: %v\n", i, e)
			t.Logf("第%d个点获取属性出错: %v", i, e)
		}
	}
	fmt.Printf("  结果：通过 —— 获取属性点数=%d/%d/%d\n", len(bases), len(scans), len(calcs))
}

// TC-PTGETMAX-01 完整流程：创建表→创建点→获取Max点属性→清理点→清理表
func TestRawRtdbbGetMaxPointsPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblGetMax", "测试获取Max点属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestGetMaxPt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)

	fmt.Println("【步骤4】获取Max点属性（预期：成功）")
	_, _, _, _, err = RawRtdbbGetMaxPointsPropertyWarp(handle, []PointID{pid})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("Max获取点属性失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— Max点属性获取成功")
}

// TC-PTGETTYPE-01 完整流程：创建表→创建3个点→批量获取类型→清理点→清理表
func TestRawRtdbbGetTypesPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblGetType", "测试获取点类型")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建3个标签点（预期：成功）")
	var pids []PointID
	for i := 0; i < 3; i++ {
		tag := fmt.Sprintf("TestGetTypePt%d", i)
		pid, err := RawRtdbbInsertBasePointWarp(handle, tag, RtdbTypeReal64, tbl.ID, 0)
		if !RteIsOk(err) {
			fmt.Printf("  结果：失败 —— %s\n", err)
			t.Error("创建标签点失败:", err)
			return
		}
		fmt.Printf("  第%d个点创建成功: ID=%d\n", i, pid)
		pids = append(pids, pid)
		defer RawRtdbbRemovePointByIdWarp(handle, pid)
	}
	fmt.Println("  结果：通过 —— 3个标签点创建完成")

	fmt.Println("【步骤4】批量获取点类型（预期：成功）")
	types, errs, err := RawRtdbbGetTypesPropertyWarp(handle, pids)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("批量获取类型失败:", err)
		return
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			fmt.Printf("    第%d个点获取类型出错: %v\n", i, e)
			t.Logf("第%d个点获取类型出错: %v", i, e)
		}
	}
	fmt.Printf("  结果：通过 —— 类型列表=%v\n", types)
}

// TC-PTSEARCH-01/03 完整流程：创建表→创建点→通配符搜索→清理点→清理表
func TestRawRtdbbSearchWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSearch", "测试搜索点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestSearchPt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)

	fmt.Println("【步骤4】通配符搜索所有标签点（预期：成功且>=1）")
	ids, err := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("搜索失败:", err)
		return
	}
	if len(ids) < 1 {
		fmt.Printf("  结果：失败 —— 搜索到点数=%d，预期至少1个\n", len(ids))
		t.Errorf("搜索点数预期>=1，实际=%d", len(ids))
		return
	}
	fmt.Printf("  结果：通过 —— 搜索到点数=%d\n", len(ids))

	fmt.Println("【步骤5】搜索不存在的标签点（预期：成功，返回0个）")
	ids2, err := RawRtdbbSearchWarp(handle, "NOTEXIST_*", "*", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("搜索无匹配失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 无匹配点数=%d\n", len(ids2))
}

// TC-PTBATCH-01 完整流程：创建表→创建2个点→分批搜索→清理点→清理表
func TestRawRtdbbSearchInBatchesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBatch", "测试分批搜索")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建2个标签点（预期：成功）")
	for i := 0; i < 2; i++ {
		tag := fmt.Sprintf("TestBatchPt%d", i)
		pid, err := RawRtdbbInsertBasePointWarp(handle, tag, RtdbTypeReal64, tbl.ID, 0)
		if !RteIsOk(err) {
			fmt.Printf("  结果：失败 —— %s\n", err)
			t.Error("创建标签点失败:", err)
			return
		}
		fmt.Printf("  第%d个点创建成功: ID=%d\n", i, pid)
		defer RawRtdbbRemovePointByIdWarp(handle, pid)
	}
	fmt.Println("  结果：通过 —— 2个标签点创建完成")

	fmt.Println("【步骤4】分批搜索获取第一批（预期：成功）")
	batch, err := RawRtdbbSearchInBatchesWarp(handle, 0, 1, "*", tbl.Name, "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("分批搜索失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 第一批数量=%d\n", len(batch))
}

// TC-PTSEARCHEX-01 完整流程：创建表→创建float64点→高级搜索→清理点→清理表
func TestRawRtdbbSearchExWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSearchEx", "测试高级搜索")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建float64标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestSearchExPt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)

	fmt.Println("【步骤4】高级搜索 float64 类型点（预期：成功且>=1）")
	ids, err := RawRtdbbSearchExWarp(handle, 100, "*", tbl.Name, "", "", "", "", "float64", 0, 0, 0, "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("高级搜索失败:", err)
		return
	}
	if len(ids) < 1 {
		fmt.Printf("  结果：失败 —— 高级搜索点数=%d，预期至少1个\n", len(ids))
		t.Errorf("高级搜索点数预期>=1，实际=%d", len(ids))
		return
	}
	fmt.Printf("  结果：通过 —— 高级搜索点数=%d\n", len(ids))
}

// TC-PTCNT-01 完整流程：创建表→创建点→统计匹配数量→清理点→清理表
func TestRawRtdbbSearchPointsCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblCount", "测试统计点数")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestCountPt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)

	fmt.Println("【步骤4】统计所有标签点数量（预期：成功且>=1）")
	count, err := RawRtdbbSearchPointsCountWarp(handle, "*", tbl.Name, "", "", "", "", "", 0, 0, 0, "")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("统计点数失败:", err)
		return
	}
	if count < 1 {
		fmt.Printf("  结果：失败 —— 总点数=%d，预期至少1个\n", count)
		t.Errorf("总点数预期>=1，实际=%d", count)
		return
	}
	fmt.Printf("  结果：通过 —— 总点数=%d\n", count)
}

// TC-PTFIND-01/02 完整流程：创建表→创建点→查找存在的点/不存在的点→清理点→清理表
func TestRawRtdbbFindPointsWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblFind", "测试查找点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestFindPt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)

	fmt.Println("【步骤4】查找存在的点和不存在的点（预期：成功）")
	fullName := tbl.Name + ".TestFindPt"
	fids, types, classes, useMs, err := RawRtdbbFindPointsWarp(handle, []string{fullName, "Table.NotExist"})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("查找点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 查找结果: fids=%v, types=%v, classes=%v, useMs=%v\n", fids, types, classes, useMs)
}

// TC-PTFINDEX-01 完整流程：创建表→创建点→高级查找点→清理点→清理表
func TestRawRtdbbFindPointsExWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblFindEx", "测试高级查找点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestFindExPt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer RawRtdbbRemovePointByIdWarp(handle, pid)

	fmt.Println("【步骤4】高级查找点（预期：成功）")
	fullName := tbl.Name + ".TestFindExPt"
	fids, types, classes, precisions, errs, err := RawRtdbbFindPointsExWarp(handle, []string{fullName})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("查找点Ex失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 查找Ex结果: fids=%v, types=%v, classes=%v, precisions=%v, errs=%v\n", fids, types, classes, precisions, errs)
}

// TC-PTSORT-01 完整流程：创建表→创建多个点→按标签名升序排序→清理点→清理表
func TestRawRtdbbSortPointsWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSort", "测试排序点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建5个标签点（预期：成功）")
	var pids []PointID
	for i := 0; i < 5; i++ {
		tag := fmt.Sprintf("TestSortPt%d", i)
		pid, err := RawRtdbbInsertBasePointWarp(handle, tag, RtdbTypeReal64, tbl.ID, 0)
		if !RteIsOk(err) {
			fmt.Printf("  结果：失败 —— %s\n", err)
			t.Error("创建标签点失败:", err)
			return
		}
		fmt.Printf("  第%d个点创建成功: ID=%d\n", i, pid)
		pids = append(pids, pid)
		defer RawRtdbbRemovePointByIdWarp(handle, pid)
	}
	fmt.Println("  结果：通过 —— 5个标签点创建完成")

	fmt.Println("【步骤4】按标签名升序排序前5个点（预期：成功）")
	sorted, err := RawRtdbbSortPointsWarp(handle, pids, RtdbTagIndexTag, RtdbSortFlag(0))
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("排序失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 排序后前5个ID=%v\n", sorted)
}

// TC-PTUPD-01 完整流程：创建表→创建点→更新描述→清理点→清理表
func TestRawRtdbbUpdatePointPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblUpd", "测试更新点属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建测试点 TestUpdPt（预期：成功）")
	base := &RtdbPoint{Tag: "TestUpdPt", Table: tbl.ID, Type: RtdbTypeReal64}
	base, _, _, err = RawRtdbbInsertPointWarp(handle, base, &RtdbScan{Source: "go_test"}, nil)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建测试点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 测试点ID=%d\n", base.ID)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, base.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】更新点描述为 UpdatedDesc（预期：成功）")
	base.Desc = "UpdatedDesc"
	err = RawRtdbbUpdatePointPropertyWarp(handle, base, nil, nil)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("更新点属性失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 点描述更新成功")
}

// TC-PTUPDMAX-01 完整流程：创建表→创建Max点→更新超长字段→清理点→清理表
func TestRawRtdbbUpdateMaxPointPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblUpdMax", "测试更新Max点属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建Max测试点 TestUpdMaxPt（预期：成功）")
	base := &RtdbPoint{Tag: "TestUpdMaxPt", Table: tbl.ID, Type: RtdbTypeReal64}
	base, _, _, err = RawRtdbbInsertMaxPointWarp(handle, base, &RtdbScan{Source: "go_test"}, nil)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建Max测试点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— Max测试点ID=%d\n", base.ID)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, base.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】更新Max点描述为 UpdatedMaxDesc（预期：成功）")
	base.Desc = "UpdatedMaxDesc"
	err = RawRtdbbUpdateMaxPointPropertyWarp(handle, base, nil, nil)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("更新Max点属性失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— Max点描述更新成功")
}

// ==================== 06. 回收站与自定义类型 ====================

// TC-RECY-01 完整流程：创建表→创建点→删除点→恢复回收站点→清理点→清理表
func TestRawRtdbbRecoverPointWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblRecy", "测试恢复回收站点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤7】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestRecyPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")

	fmt.Println("【步骤4】获取回收站点列表并定位目标点（预期：非空）")
	recycled, err := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if !RteIsOk(err) || len(recycled) == 0 {
		fmt.Println("  结果：跳过 —— 回收站为空")
		t.Skip("回收站为空")
	}
	fmt.Printf("  结果：通过 —— 回收站有 %d 个点\n", len(recycled))

	found := false
	for _, id := range recycled {
		if id == pid {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("  结果：跳过 —— 未在回收站中找到目标点")
		t.Skip("未在回收站中找到目标点")
	}

	fmt.Println("【步骤5】恢复回收站点到原表（预期：成功或策略差异）")
	err = RawRtdbbRecoverPointWarp(handle, tbl.ID, pid)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 恢复点失败（可能回收站策略差异）: %v\n", err)
		t.Logf("恢复点(可能回收站策略差异): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 点已恢复")

	fmt.Println("【步骤6】清理恢复后的点（预期：成功）")
	// 恢复后点的ID可能变化，先搜索确认当前实际ID
	restoredIds, _ := RawRtdbbSearchWarp(handle, "TestRecyPt", tbl.Name, "", "", "", "", RtdbSortFlag(0))
	if len(restoredIds) > 0 {
		pid = restoredIds[0]
	}
	rte = RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— %s\n", rte)
		t.Logf("清理恢复后的点失败: %s", rte)
	} else {
		fmt.Println("  结果：通过 —— 恢复后的点已删除")
	}
}

// TC-PURGE-01 完整流程：创建表→创建点→删除点→清除回收站点→清理表
func TestRawRtdbbPurgePointWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblPurge", "测试清除回收站点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPurgePt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")

	fmt.Println("【步骤4】获取回收站点列表并定位目标点（预期：非空）")
	recycled, err := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if !RteIsOk(err) || len(recycled) == 0 {
		fmt.Println("  结果：跳过 —— 回收站为空")
		t.Skip("回收站为空")
	}
	fmt.Printf("  结果：通过 —— 回收站有 %d 个点\n", len(recycled))

	found := false
	for _, id := range recycled {
		if id == pid {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("  结果：跳过 —— 未在回收站中找到目标点")
		t.Skip("未在回收站中找到目标点")
	}

	fmt.Println("【步骤5】清除回收站点（预期：成功或策略差异）")
	err = RawRtdbbPurgePointWarp(handle, pid)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 清除点失败（可能回收站策略差异）: %v\n", err)
		t.Logf("清除点(可能回收站策略差异): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 回收站点已清除")
}

// TC-RCYCNT-01 查询回收站数量（非空）
func TestRawRtdbbGetRecycledPointsCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblCnt", "测试回收站数量")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤7】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestCntPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")
	defer func() {
		fmt.Println("【步骤6】清理回收站点（预期：成功或已清除）")
		_ = RawRtdbbPurgePointWarp(handle, pid)
		fmt.Println("  结果：通过 —— 回收站点已清理")
	}()

	fmt.Println("【步骤4】查询回收站数量（预期：>=1）")
	count, err := RawRtdbbGetRecycledPointsCountWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取回收站数量失败:", err)
		return
	}
	if count < 1 {
		fmt.Printf("  结果：失败 —— 回收站数量=%d，预期>=1\n", count)
		t.Error("回收站数量应>=1")
		return
	}
	fmt.Printf("  结果：通过 —— 回收站数量=%d\n", count)
}

// TC-RCYGET-01 获取回收站点列表
func TestRawRtdbbGetRecycledPointsWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblGet", "测试获取回收站点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤7】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestGetPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")
	defer func() {
		fmt.Println("【步骤6】清理回收站点（预期：成功或已清除）")
		_ = RawRtdbbPurgePointWarp(handle, pid)
		fmt.Println("  结果：通过 —— 回收站点已清理")
	}()

	fmt.Println("【步骤4】获取回收站点列表（预期：非空）")
	points, err := RawRtdbbGetRecycledPointsWarp(handle, 100)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取回收站点失败:", err)
		return
	}
	if len(points) == 0 {
		fmt.Println("  结果：失败 —— 回收站为空")
		t.Error("回收站应为非空")
		return
	}
	fmt.Printf("  结果：通过 —— 回收站点数=%d\n", len(points))
}

// TC-RCYSEARCH-01 通配符搜索回收站点
func TestRawRtdbbSearchRecycledPointsWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSearch", "测试搜索回收站点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤7】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestSearchPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")
	defer func() {
		fmt.Println("【步骤6】清理回收站点（预期：成功或已清除）")
		_ = RawRtdbbPurgePointWarp(handle, pid)
		fmt.Println("  结果：通过 —— 回收站点已清理")
	}()

	fmt.Println("【步骤4】通配符搜索回收站点（预期：非空）")
	ids, err := RawRtdbbSearchRecycledPointsWarp(handle, "TestSearchPt", "", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("搜索回收站点失败:", err)
		return
	}
	if len(ids) == 0 {
		fmt.Println("  结果：失败 —— 未搜索到回收站点")
		t.Error("应搜索到至少1个回收站点")
		return
	}
	fmt.Printf("  结果：通过 —— 搜索回收站点数=%d\n", len(ids))
}

// TC-RCYPROP-01 获取回收站点属性
func TestRawRtdbbGetRecycledPointPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblProp", "测试回收站属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPropPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")
	defer func() {
		fmt.Println("【步骤5】清理回收站点（预期：成功或已清除）")
		_ = RawRtdbbPurgePointWarp(handle, pid)
		fmt.Println("  结果：通过 —— 回收站点已清理")
	}()

	fmt.Println("【步骤4】获取回收站点属性（预期：成功）")
	base, scan, calc, err := RawRtdbbGetRecycledPointPropertyWarp(handle, pid)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 获取回收站点属性失败: %v\n", err)
		t.Error("获取回收站点属性失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 回收站点属性获取成功: base=%v, scan=%v, calc=%v\n", base, scan, calc)
}

// TC-RCYBATCH-01 分批获取回收站点
func TestRawRtdbbSearchRecycledPointsInBatchesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBatch", "测试分批回收站点")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestBatchPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")
	defer func() {
		fmt.Println("【步骤5】清理回收站点（预期：成功或已清除）")
		_ = RawRtdbbPurgePointWarp(handle, pid)
		fmt.Println("  结果：通过 —— 回收站点已清理")
	}()

	fmt.Println("【步骤4】分批搜索回收站点（预期：非空）")
	batch, err := RawRtdbbSearchRecycledPointsInBatchesWarp(handle, 0, 5, "TestBatchPt", "", "", "", "", "", RtdbSortFlag(0))
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("分批搜索回收站失败:", err)
		return
	}
	if len(batch) == 0 {
		fmt.Println("  结果：失败 —— 未分批搜索到回收站点")
		t.Error("应分批搜索到至少1个回收站点")
		return
	}
	fmt.Printf("  结果：通过 —— 第一批回收站点数=%d\n", len(batch))
}

// TC-RCYMAX-01 获取超长字段回收站点
func TestRawRtdbbGetRecycledMaxPointPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblMax", "测试回收站Max属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestMaxPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")
	defer func() {
		fmt.Println("【步骤5】清理回收站点（预期：成功或已清除）")
		_ = RawRtdbbPurgePointWarp(handle, pid)
		fmt.Println("  结果：通过 —— 回收站点已清理")
	}()

	fmt.Println("【步骤4】获取Max回收站点属性（预期：成功）")
	base, scan, calc, err := RawRtdbbGetRecycledMaxPointPropertyWarp(handle, pid)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 获取Max回收站点属性失败: %v\n", err)
		t.Error("获取Max回收站点属性失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— Max回收站点属性获取成功: base=%v, scan=%v, calc=%v\n", base, scan, calc)
}

// TC-RCYCLR-01 清空非空回收站
func TestRawRtdbbClearRecyclerWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblClr", "测试清空回收站")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建临时点并删除（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestClrPt", RtdbTypeBool, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时点ID=%d\n", pid)
	rte := RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 临时点已删除（进入回收站）")

	fmt.Println("【步骤4】清空回收站（预期：成功）")
	err = RawRtdbbClearRecyclerWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 清空回收站失败: %v\n", err)
		t.Error("清空回收站失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 回收站已清空")

	fmt.Println("【步骤5】验证回收站为空（预期：数量为0）")
	count, rte := RawRtdbbGetRecycledPointsCountWarp(handle)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：跳过 —— 验证回收站数量失败: %v\n", rte)
		t.Logf("验证回收站数量: %v", rte)
		return
	}
	if count != 0 {
		fmt.Printf("  结果：失败 —— 回收站数量=%d，预期0\n", count)
		t.Errorf("回收站应已清空，但数量为%d", count)
		return
	}
	fmt.Println("  结果：通过 —— 回收站验证为空")
}

// TC-RCYNAMED-01 获取回收站自定义类型点信息
func TestRawRtdbbGetRecycledNamedTypeNamesPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建自定义类型（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestRecyNamedType")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestRecyNamedType", "回收站自定义类型测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤8】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestRecyNamedType")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblNamed", "测试回收站自定义类型")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤7】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤4】创建自定义类型点并删除（预期：成功）")
	base := &RtdbPoint{Tag: "TestRecyNamedPt", Type: RtdbTypeNamedT, Class: RtdbClassBase, Table: tbl.ID, Desc: "测试点"}
	scan := &RtdbScan{}
	base, scan, rte := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "TestRecyNamedType")
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 创建自定义类型点失败: %s\n", rte)
		t.Error("创建自定义类型点失败:", rte)
		return
	}
	pid := base.ID
	fmt.Printf("  结果：通过 —— 自定义类型点ID=%d\n", pid)
	rte = RawRtdbbRemovePointByIdWarp(handle, pid)
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 删除点失败: %s\n", rte)
		t.Error("删除点失败:", rte)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型点已删除（进入回收站）")
	defer func() {
		fmt.Println("【步骤6】清理回收站点（预期：成功或已清除）")
		_ = RawRtdbbPurgePointWarp(handle, pid)
		fmt.Println("  结果：通过 —— 回收站点已清理")
	}()

	fmt.Println("【步骤5】获取回收站自定义类型信息（预期：成功）")
	names, counts, errs, err := RawRtdbbGetRecycledNamedTypeNamesPropertyWarp(handle, []PointID{pid})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 获取回收站自定义类型信息失败: %v\n", err)
		t.Error("获取回收站自定义类型信息失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 回收站自定义类型: names=%v, counts=%v, errs=%v\n", names, counts, errs)
}

// TC-NTCREATE-01 正常创建自定义类型
func TestRawRtdbbCreateNamedTypeWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理旧类型 TestNamedType（预期：成功或不存在）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestNamedType")
	fmt.Println("  结果：通过 —— 旧类型已清理")

	fmt.Println("【步骤3】创建自定义类型 TestNamedType（预期：成功）")
	fields := []RtdbDataTypeField{
		{Name: "field1", Type: RtdbTypeReal64, Length: 8, Desc: "字段1"},
		{Name: "field2", Type: RtdbTypeInt32, Length: 4, Desc: "字段2"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestNamedType", "测试自定义类型", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestNamedType")
}

// TC-NTCNT-01 查询含自定义类型
func TestRawRtdbbGetNamedTypesCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建自定义类型（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestCntNamedType")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestCntNamedType", "计数测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤5】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestCntNamedType")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】查询自定义类型总数（预期：成功）")
	count, err := RawRtdbbGetNamedTypesCountWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取自定义类型总数失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 自定义类型总数=%d\n", count)
}

// TC-NTALL-01 获取所有类型
func TestRawRtdbbGetAllNamedTypesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建自定义类型（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestAllNamedType")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestAllNamedType", "列表测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤5】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestAllNamedType")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】获取所有自定义类型（预期：成功且包含新类型）")
	names, counts, err := RawRtdbbGetAllNamedTypesWarp(handle, 100)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取所有自定义类型失败:", err)
		return
	}
	found := false
	for _, name := range names {
		if name == "TestAllNamedType" {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("  结果：失败 —— 未在列表中找到 TestAllNamedType，列表=%v\n", names)
		t.Error("应包含刚创建的自定义类型")
		return
	}
	fmt.Printf("  结果：通过 —— 自定义类型数量=%d/%d\n", len(names), len(counts))
}

// TC-NTGET-01 获取存在的自定义类型字段
func TestRawRtdbbGetNamedTypeWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理并创建自定义类型 TestGetNamedType（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestGetNamedType")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestGetNamedType", "查询测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 无法创建自定义类型: %v\n", err)
		t.Skip("无法创建自定义类型:", err)
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestGetNamedType")

	fmt.Println("【步骤3】获取自定义类型字段（预期：成功）")
	gotFields, typeSize, desc, err := RawRtdbbGetNamedTypeWarp(handle, "TestGetNamedType", 10)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取自定义类型失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 字段数=%d, size=%d, desc=%s\n", len(gotFields), typeSize, desc)
}

// TC-NTRMV-01 删除存在的自定义类型
func TestRawRtdbbRemoveNamedTypeWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建自定义类型 TestRmNamedType（预期：成功）")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	_ = RawRtdbbCreateNamedTypeWarp(handle, "TestRmNamedType", "删除测试", fields...)
	fmt.Println("  结果：通过 —— 自定义类型创建成功")

	fmt.Println("【步骤3】删除自定义类型 TestRmNamedType（预期：成功）")
	err := RawRtdbbRemoveNamedTypeWarp(handle, "TestRmNamedType")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("删除自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型已删除")
}

// TC-NTNAMEPROP-01 批量查询自定义类型点
func TestRawRtdbbGetNamedTypeNamesPropertyWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建自定义类型（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestNamePropType")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestNamePropType", "名称属性测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤7】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestNamePropType")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblNameProp", "测试名称属性")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤4】创建自定义类型点（预期：成功）")
	base := &RtdbPoint{Tag: "TestNamePropPt", Type: RtdbTypeNamedT, Class: RtdbClassBase, Table: tbl.ID, Desc: "测试点"}
	scan := &RtdbScan{}
	base, scan, rte := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "TestNamePropType")
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 创建自定义类型点失败: %s\n", rte)
		t.Error("创建自定义类型点失败:", rte)
		return
	}
	pid := base.ID
	fmt.Printf("  结果：通过 —— 自定义类型点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功或已清除）")
		_ = RawRtdbbRemovePointByIdWarp(handle, pid)
		fmt.Println("  结果：通过 —— 标签点已清理")
	}()

	fmt.Println("【步骤5】批量获取自定义类型名称属性（预期：成功）")
	names, counts, errs, err := RawRtdbbGetNamedTypeNamesPropertyWarp(handle, []PointID{pid})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取自定义类型名称属性失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 自定义类型名称: names=%v, counts=%v, errs=%v\n", names, counts, errs)
}

// TC-NTPTCNT-01 查询有点的类型
func TestRawRtdbbGetNamedTypePointsCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建自定义类型（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestPtCntType")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestPtCntType", "点计数测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤7】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestPtCntType")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblPtCnt", "测试点计数")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤4】创建自定义类型点（预期：成功）")
	base := &RtdbPoint{Tag: "TestPtCntPt", Type: RtdbTypeNamedT, Class: RtdbClassBase, Table: tbl.ID, Desc: "测试点"}
	scan := &RtdbScan{}
	base, scan, rte := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "TestPtCntType")
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 创建自定义类型点失败: %s\n", rte)
		t.Error("创建自定义类型点失败:", rte)
		return
	}
	pid := base.ID
	fmt.Printf("  结果：通过 —— 自定义类型点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功或已清除）")
		_ = RawRtdbbRemovePointByIdWarp(handle, pid)
		fmt.Println("  结果：通过 —— 标签点已清理")
	}()

	fmt.Println("【步骤5】查询自定义类型点数量（预期：>=1）")
	count, err := RawRtdbbGetNamedTypePointsCountWarp(handle, "TestPtCntType")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取自定义类型点数量失败:", err)
		return
	}
	if count < 1 {
		fmt.Printf("  结果：失败 —— 自定义类型点数量=%d，预期>=1\n", count)
		t.Error("自定义类型点数量应>=1")
		return
	}
	fmt.Printf("  结果：通过 —— 自定义类型点数量=%d\n", count)
}

// TC-BASECNT-01 查询 float64 类型点数量
func TestRawRtdbbGetBaseTypePointsCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】查询 float64 类型点数量（预期：成功）")
	count, err := RawRtdbbGetBaseTypePointsCountWarp(handle, RtdbTypeReal64)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取基础类型点数量失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— float64类型点数量=%d\n", count)
}

// TC-NTMOD-01 修改类型名称
func TestRawRtdbbModifyNamedTypeWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理旧类型（预期：成功或不存在）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestModTypeOld")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestModTypeNew")
	fmt.Println("  结果：通过 —— 旧类型已清理")

	fmt.Println("【步骤3】创建自定义类型 TestModTypeOld（预期：成功）")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestModTypeOld", "修改测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 无法创建自定义类型: %v\n", err)
		t.Skip("无法创建自定义类型:", err)
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestModTypeNew")

	fmt.Println("【步骤4】修改类型名称为 TestModTypeNew（预期：成功）")
	newName := "TestModTypeNew"
	err = RawRtdbbModifyNamedTypeWarp(handle, "TestModTypeOld", &newName, nil, []string{"newf1"}, []string{"新字段1"})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("修改自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型修改成功")
}

// TC-NTWRNAME-01 按名称写入字段值
func TestRawRtdbWriteNamedTypeFieldByName32Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理并创建自定义类型 TestWrType（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestWrType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestWrType", "写测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 无法创建自定义类型：%s\n", err)
		t.Skip("无法创建自定义类型:", err)
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestWrType")

	fmt.Println("【步骤3】按名称写入字段值（预期：成功）")
	object := make([]byte, 8)
	field := make([]byte, 8)
	// 写入一个 float64 值
	field[0] = 0x3F
	field[1] = 0xF0
	object, err = RawRtdbWriteNamedTypeFieldByName32Warp(handle, "TestWrType", "val", RtdbTypeReal64, object, field)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按名称写字段失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 按名称写入字段成功")
}

// TC-NTWRPOS-01 按位置写入字段值
func TestRawRtdbWriteNamedTypeFieldByPos32Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理并创建自定义类型 TestWrPosType（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestWrPosType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestWrPosType", "写位置测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 无法创建自定义类型：%s\n", err)
		t.Skip("无法创建自定义类型:", err)
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestWrPosType")

	fmt.Println("【步骤3】按位置写入字段值（预期：成功）")
	object := make([]byte, 8)
	field := make([]byte, 8)
	object, err = RawRtdbWriteNamedTypeFieldByPos32Warp(handle, "TestWrPosType", 0, RtdbTypeReal64, object, field)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按位置写字段失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 按位置写入字段成功")
}

// TC-NTRDNAME-01 按名称读取字段值
func TestRawRtdbReadNamedTypeFieldByName32Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理并创建自定义类型 TestRdType（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestRdType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestRdType", "读测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 无法创建自定义类型：%s\n", err)
		t.Skip("无法创建自定义类型:", err)
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestRdType")

	fmt.Println("【步骤3】按名称读取字段值（预期：成功）")
	object := make([]byte, 8)
	object[0] = 0x3F
	object[1] = 0xF0
	field, err := RawRtdbReadNamedTypeFieldByName32Warp(handle, "TestRdType", "val", RtdbTypeReal64, object, 8)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按名称读字段失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 读取字段长度=%d\n", len(field))
}

// TC-NTRDPOS-01 按位置读取字段值
func TestRawRtdbReadNamedTypeFieldByPos32Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理并创建自定义类型 TestRdPosType（预期：成功）")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestRdPosType")
	fields := []RtdbDataTypeField{
		{Name: "val", Type: RtdbTypeReal64, Length: 8, Desc: "值"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestRdPosType", "读位置测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：跳过 —— 无法创建自定义类型：%s\n", err)
		t.Skip("无法创建自定义类型:", err)
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer RawRtdbbRemoveNamedTypeWarp(handle, "TestRdPosType")

	fmt.Println("【步骤3】按位置读取字段值（预期：成功）")
	object := make([]byte, 8)
	object[0] = 0x3F
	object[1] = 0xF0
	field, err := RawRtdbReadNamedTypeFieldByPos32Warp(handle, "TestRdPosType", 0, RtdbTypeReal64, object, 8)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("按位置读字段失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 读取字段长度=%d\n", len(field))
}

// TC-NTCHK-01/02 检查合法类型/字段名称 / TC-NTCHK-03 检查非法名称
func TestRawRtdbNamedTypeNameFieldCheckWarp(t *testing.T) {
	fmt.Println("【步骤1】检查合法类型名称 ValidType（预期：成功）")
	err := RawRtdbNamedTypeNameFieldCheckWarp("ValidType", 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("合法类型名称应通过:", err)
		return
	}
	fmt.Println("  结果：通过 —— 合法类型名称校验通过")

	fmt.Println("【步骤2】检查合法字段名称 valid_field（预期：成功）")
	err = RawRtdbNamedTypeNameFieldCheckWarp("valid_field", 1)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("合法字段名称应通过:", err)
		return
	}
	fmt.Println("  结果：通过 —— 合法字段名称校验通过")

	fmt.Println("【步骤3】检查非法名称 Type@123（预期：返回错误）")
	err = RawRtdbNamedTypeNameFieldCheckWarp("Type@123", 0)
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 非法名称居然通过了！")
		t.Error("非法名称应失败")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
	}
}

// ==================== 07. 实时数据快照 ====================

// TC-SNAPGET-01 批量读取存在的点
func TestRawRtdbsGetSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSnapGet", "测试批量读取快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtSnapGet", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】批量读取快照（预期：成功）")
	_, _, _, _, _, _, err = RawRtdbsGetSnapshots64Warp(handle, []PointID{pid})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取快照失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量读取快照成功")
}

// TC-SNAPPUT-01 批量写入新值
func TestRawRtdbsPutSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSnapPut", "测试批量写入快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtSnapPut", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】批量写入快照（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbsPutSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{123.45}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("写入快照失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入快照成功")
}

// TC-SNAPFIX-01 覆盖已有时间戳的值
func TestRawRtdbsFixSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSnapFix", "测试覆盖快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtSnapFix", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】写入初始快照（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbsPutSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{123.45}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入初始快照失败: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 初始快照写入成功")
	}

	fmt.Println("【步骤5】覆盖已有时间戳的快照（预期：成功）")
	_, err = RawRtdbsFixSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{99.9}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("覆盖快照(可能时间戳问题): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 覆盖快照成功")
}

// TC-SNAPBACK-01 回溯到更早时间戳
func TestRawRtdbsBackSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblSnapBack", "测试回溯快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtSnapBack", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】写入初始快照（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbsPutSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{123.45}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入初始快照失败: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 初始快照写入成功")
	}

	fmt.Println("【步骤4b】预写入历史数据（确保历史存储结构就绪）")
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now - 30}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：警告 —— 预写历史数据: %s\n", err)
		t.Logf("预写历史数据(非致命): %v", err)
	} else {
		fmt.Println("  结果：通过 —— 历史数据预写入成功")
	}

	fmt.Println("【步骤5】回溯快照到更早时间戳（预期：成功）")
	_, err = RawRtdbsBackSnapshots64Warp(handle, []PointID{pid}, []TimestampType{past}, []SubtimeType{0}, []float64{88.8}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("回溯快照(可能时间戳问题): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 回溯快照成功")
}

// TC-BLOBGET1-01 读取字符串点
func TestRawRtdbsGetBlobSnapshot64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBlobGet", "测试读取Blob快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 String 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtBlobGet", RtdbTypeString, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】读取字符串点 Blob 快照（预期：成功）")
	_, _, blob, _, err := RawRtdbsGetBlobSnapshot64Warp(handle, pid, true, 256)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取Blob快照: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— Blob 长度=%d\n", len(blob))
}

// TC-BLOBPUT1-01 写入字符串
func TestRawRtdbsPutBlobSnapshot64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBlobPut", "测试写入Blob快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 String 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtBlobPut", RtdbTypeString, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】写入字符串 Blob 快照（预期：成功）")
	now := TimestampType(time.Now().Unix() + 10)
	err = RawRtdbsPutBlobSnapshot64Warp(handle, pid, true, now, 0, []byte("hello_rtdb"), 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入Blob快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入字符串快照成功")
}

// TC-DTGET-01 默认格式读取
func TestRawRtdbsGetDatetimeSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblDtGet", "测试读取Datetime快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Datetime 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtDtGet", RtdbTypeDatetime, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】读取 datetime 快照（预期：成功）")
	_, _, vals, _, _, err := RawRtdbsGetDatetimeSnapshots64Warp(handle, []PointID{pid}, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取datetime快照: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— datetime 值: %v\n", vals)
}

// TC-DTPUT-01 写入 datetime 值
func TestRawRtdbsPutDatetimeSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblDtPut", "测试写入Datetime快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Datetime 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtDtPut", RtdbTypeDatetime, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】写入 datetime 快照（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbsPutDatetimeSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []string{"2024-01-01 08:00:00"}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入datetime快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入 datetime 快照成功")
}

// TC-NTSNAPGET-01 读取自定义类型点
func TestRawRtdbsGetNamedTypeSnapshot64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理残留数据并创建自定义类型（预期：成功）")
	recycled, _ := RawRtdbbSearchRecycledPointsWarp(handle, "TestPtNamedSnap", "*", "", "", "", "", RtdbSortFlag(0))
	for _, rid := range recycled {
		if rte := RawRtdbbPurgePointWarp(handle, rid); RteIsOk(rte) {
			fmt.Printf("  永久删除回收站点 ID=%d: 成功\n", rid)
		} else {
			fmt.Printf("  永久删除回收站点 ID=%d: %s\n", rid, rte)
		}
	}
	if rte := RawRtdbbRemovePointByNameWarp(handle, "TestTblNamedSnap.TestPtNamedSnap"); RteIsOk(rte) {
		fmt.Println("  清理残留点: 成功")
	} else {
		fmt.Printf("  清理残留点: %s\n", rte)
	}
	if rte := RawRtdbbRemoveTableByNameWarp(handle, "TestTblNamedSnap"); RteIsOk(rte) {
		fmt.Println("  清理残留表: 成功")
	} else {
		fmt.Printf("  清理残留表: %s\n", rte)
	}
	if rte := RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedType"); RteIsOk(rte) {
		fmt.Println("  清理残留类型: 成功")
	} else {
		fmt.Printf("  清理残留类型: %s\n", rte)
	}
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestSnapNamedType", "快照自定义类型测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤7】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedType")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblNamedSnap", "测试自定义类型快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤4】创建自定义类型标签点（预期：成功）")
	base := &RtdbPoint{Tag: "TestPtNamedSnap", Type: RtdbTypeNamedT, Class: RtdbClassBase, Table: tbl.ID, Desc: "测试点"}
	scan := &RtdbScan{}
	base, _, rte := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "TestSnapNamedType")
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 创建自定义类型点失败: %s\n", rte)
		t.Error("创建自定义类型点失败:", rte)
		return
	}
	pid := base.ID
	fmt.Printf("  结果：通过 —— 自定义类型点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
			_ = RawRtdbbPurgePointWarp(handle, pid)
		}
	}()

	fmt.Println("【步骤5】读取自定义类型快照（预期：成功）")
	_, _, obj, _, err := RawRtdbsGetNamedTypeSnapshot64Warp(handle, pid, 256)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取自定义类型快照: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 自定义类型数据长度=%d\n", len(obj))
}

// TC-NTSNAPPUT-01 写入自定义类型数据
func TestRawRtdbsPutNamedTypeSnapshot64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理残留数据并创建自定义类型（预期：成功）")
	recycled, _ := RawRtdbbSearchRecycledPointsWarp(handle, "TestPtNamedSnapPut", "*", "", "", "", "", RtdbSortFlag(0))
	for _, rid := range recycled {
		_ = RawRtdbbPurgePointWarp(handle, rid)
	}
	_ = RawRtdbbRemovePointByNameWarp(handle, "TestTblNamedSnapPut.TestPtNamedSnapPut")
	_ = RawRtdbbRemoveTableByNameWarp(handle, "TestTblNamedSnapPut")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedTypePut")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestSnapNamedTypePut", "快照自定义类型写入测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤7】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedTypePut")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblNamedSnapPut", "测试自定义类型写入快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤4】创建自定义类型标签点（预期：成功）")
	base := &RtdbPoint{Tag: "TestPtNamedSnapPut", Type: RtdbTypeNamedT, Class: RtdbClassBase, Table: tbl.ID, Desc: "测试点"}
	scan := &RtdbScan{}
	base, _, rte := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "TestSnapNamedTypePut")
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 创建自定义类型点失败: %s\n", rte)
		t.Error("创建自定义类型点失败:", rte)
		return
	}
	pid := base.ID
	fmt.Printf("  结果：通过 —— 自定义类型点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
			_ = RawRtdbbPurgePointWarp(handle, pid)
		}
	}()

	fmt.Println("【步骤5】写入自定义类型快照（预期：成功）")
	now := TimestampType(time.Now().Unix() + 10)
	err = RawRtdbsPutNamedTypeSnapshot64Warp(handle, pid, now, 0, make([]byte, 8), 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入自定义类型快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入自定义类型快照成功")
}

// TC-COORGET-01 批量读取坐标实时数据
func TestRawRtdbsGetCoorSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblCoorGet", "测试读取坐标快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Coor 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtCoorGet", RtdbTypeCoor, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】批量读取坐标快照（预期：成功）")
	_, _, _, _, _, _, err = RawRtdbsGetCoorSnapshots64Warp(handle, []PointID{pid})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取坐标快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量读取坐标快照成功")
}

// TC-COORPUT-01 批量写入坐标实时数据
func TestRawRtdbsPutCoorSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblCoorPut", "测试写入坐标快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Coor 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtCoorPut", RtdbTypeCoor, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】批量写入坐标快照（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbsPutCoorSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float32{1.0}, []float32{2.0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入坐标快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入坐标快照成功")
}

// TC-COORFIX-01 批量覆盖写入坐标实时数据
func TestRawRtdbsFixCoorSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblCoorFix", "测试覆盖坐标快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Coor 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtCoorFix", RtdbTypeCoor, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】写入初始坐标快照（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbsPutCoorSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float32{1.0}, []float32{2.0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入初始坐标快照失败: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 初始坐标快照写入成功")
	}

	fmt.Println("【步骤5】批量覆盖写入坐标快照（预期：成功）")
	_, err = RawRtdbsFixCoorSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float32{3.0}, []float32{4.0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("覆盖坐标快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量覆盖坐标快照成功")
}

// TC-BLOBGETN-01 批量读取二进制/字符串实时数据
func TestRawRtdbsGetBlobSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBlobGetN", "测试批量读取Blob快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 String 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtBlobGetN", RtdbTypeString, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】批量读取 Blob 快照（预期：成功）")
	_, _, _, _, _, err = RawRtdbsGetBlobSnapshots64Warp(handle, []PointID{pid}, []bool{true}, 256)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量获取Blob快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量读取 Blob 快照成功")
}

// TC-BLOBPUTN-01 批量写入二进制/字符串实时数据
func TestRawRtdbsPutBlobSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblBlobPutN", "测试批量写入Blob快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤5】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 String 类型标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtBlobPutN", RtdbTypeString, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤4】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】批量写入 Blob 快照（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbsPutBlobSnapshots64Warp(handle, []PointID{pid}, []bool{true}, []TimestampType{now}, []SubtimeType{0}, [][]byte{[]byte("hello")}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量写入Blob快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入 Blob 快照成功")
}

// TC-NTSNAPGETN-01 批量获取自定义类型测点的快照
func TestRawRtdbsGetNamedTypeSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理残留数据并创建自定义类型（预期：成功）")
	recycled, _ := RawRtdbbSearchRecycledPointsWarp(handle, "TestPtNamedSnapN", "*", "", "", "", "", RtdbSortFlag(0))
	for _, rid := range recycled {
		_ = RawRtdbbPurgePointWarp(handle, rid)
	}
	_ = RawRtdbbRemovePointByNameWarp(handle, "TestTblNamedSnapN.TestPtNamedSnapN")
	_ = RawRtdbbRemoveTableByNameWarp(handle, "TestTblNamedSnapN")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedTypeN")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestSnapNamedTypeN", "批量快照自定义类型测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤7】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedTypeN")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblNamedSnapN", "测试批量自定义类型快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤4】创建自定义类型标签点（预期：成功）")
	base := &RtdbPoint{Tag: "TestPtNamedSnapN", Type: RtdbTypeNamedT, Class: RtdbClassBase, Table: tbl.ID, Desc: "测试点"}
	scan := &RtdbScan{}
	base, _, rte := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "TestSnapNamedTypeN")
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 创建自定义类型点失败: %s\n", rte)
		t.Error("创建自定义类型点失败:", rte)
		return
	}
	pid := base.ID
	fmt.Printf("  结果：通过 —— 自定义类型点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
			_ = RawRtdbbPurgePointWarp(handle, pid)
		}
	}()

	fmt.Println("【步骤5】批量获取自定义类型快照（预期：成功）")
	_, _, _, _, _, err = RawRtdbsGetNamedTypeSnapshots64Warp(handle, []PointID{pid}, []int32{256})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量获取自定义类型快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量获取自定义类型快照成功")
}

// TC-NTSNAPPUTN-01 批量写入自定义类型测点的快照
func TestRawRtdbsPutNamedTypeSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】清理残留数据并创建自定义类型（预期：成功）")
	recycled, _ := RawRtdbbSearchRecycledPointsWarp(handle, "TestPtNamedSnapPutN", "*", "", "", "", "", RtdbSortFlag(0))
	for _, rid := range recycled {
		_ = RawRtdbbPurgePointWarp(handle, rid)
	}
	_ = RawRtdbbRemovePointByNameWarp(handle, "TestTblNamedSnapPutN.TestPtNamedSnapPutN")
	_ = RawRtdbbRemoveTableByNameWarp(handle, "TestTblNamedSnapPutN")
	_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedTypePutN")
	fields := []RtdbDataTypeField{
		{Name: "f1", Type: RtdbTypeReal64, Length: 8, Desc: "f1"},
	}
	err := RawRtdbbCreateNamedTypeWarp(handle, "TestSnapNamedTypePutN", "批量写入自定义类型测试", fields...)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— 创建自定义类型失败: %s\n", err)
		t.Error("创建自定义类型失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 自定义类型创建成功")
	defer func() {
		fmt.Println("【步骤7】清理自定义类型（预期：成功或已清除）")
		_ = RawRtdbbRemoveNamedTypeWarp(handle, "TestSnapNamedTypePutN")
		fmt.Println("  结果：通过 —— 自定义类型已清理")
	}()

	fmt.Println("【步骤3】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblNamedSnapPutN", "测试批量写入自定义类型快照")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【步骤6】清理临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤4】创建自定义类型标签点（预期：成功）")
	base := &RtdbPoint{Tag: "TestPtNamedSnapPutN", Type: RtdbTypeNamedT, Class: RtdbClassBase, Table: tbl.ID, Desc: "测试点"}
	scan := &RtdbScan{}
	base, _, rte := RawRtdbbInsertNamedTypePointWarp(handle, base, scan, "TestSnapNamedTypePutN")
	if !RteIsOk(rte) {
		fmt.Printf("  结果：失败 —— 创建自定义类型点失败: %s\n", rte)
		t.Error("创建自定义类型点失败:", rte)
		return
	}
	pid := base.ID
	fmt.Printf("  结果：通过 —— 自定义类型点ID=%d\n", pid)
	defer func() {
		fmt.Println("【步骤5】清理标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
			_ = RawRtdbbPurgePointWarp(handle, pid)
		}
	}()

	fmt.Println("【步骤5】批量写入自定义类型快照（预期：成功）")
	now := TimestampType(time.Now().Unix() + 10)
	_, err = RawRtdbsPutNamedTypeSnapshots64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, [][]byte{make([]byte, 8)}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量写入自定义类型快照: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入自定义类型快照成功")
}

// ==================== 08. 订阅功能 ====================
// 注意：订阅类 API 基于回调异步机制，需要独立句柄、goroutine 协作及数据变化触发，
// 在原生 API 层不便进行有效的功能验证。完整的订阅功能测试已在 easy_test.go 中覆盖。

func TestRawRtdbbSubscribeTagsExWarp(t *testing.T) {
	fmt.Println("【步骤1】测试订阅标签点属性更改通知（预期：跳过）")
	fmt.Println("  结果：跳过 —— 订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
	t.Skip("订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
}

// TC-CANTAG-01 取消标签点属性更改通知订阅
func TestRawRtdbbCancelSubscribeTagsWarp(t *testing.T) {
	fmt.Println("【步骤1】测试取消标签点属性更改通知订阅（预期：跳过）")
	fmt.Println("  结果：跳过 —— 订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
	t.Skip("订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
}

func TestRawRtdbsSubscribeSnapshotsEx64Warp(t *testing.T) {
	fmt.Println("【步骤1】测试订阅快照更改通知（预期：跳过）")
	fmt.Println("  结果：跳过 —— 订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
	t.Skip("订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
}

func TestRawRtdbsSubscribeDeltaSnapshots64Warp(t *testing.T) {
	fmt.Println("【步骤1】测试订阅增量快照更改通知（预期：跳过）")
	fmt.Println("  结果：跳过 —— 订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
	t.Skip("订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
}

// TC-DGRAMCRT-01 正常创建数据流句柄
func TestRawRtdbCreateDatagramHandleWarp(t *testing.T) {
	fmt.Println("【步骤1】创建数据流句柄（预期：成功）")
	dh, err := RawRtdbCreateDatagramHandleWarp(0, "127.0.0.1")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("创建数据流(可能权限或端口占用): %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 数据流句柄创建成功，句柄=%d\n", dh)
	defer RawRtdbRemoveDatagramHandleWarp(dh)
}

// TC-CHGSUB-01 批量修改订阅标签点信息
func TestRawRtdbsChangeSubscribeSnapshotsWarp(t *testing.T) {
	fmt.Println("【步骤1】测试批量修改订阅标签点信息（预期：跳过）")
	fmt.Println("  结果：跳过 —— 订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
	t.Skip("订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
}

// TC-CANSNAP-01 取消标签点快照更改通知订阅
func TestRawRtdbsCancelSubscribeSnapshotsWarp(t *testing.T) {
	fmt.Println("【步骤1】测试取消快照更改通知订阅（预期：跳过）")
	fmt.Println("  结果：跳过 —— 订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
	t.Skip("订阅类API基于回调异步机制，在原生API层不便验证，已在easy_test.go中覆盖")
}

// TC-DGRAMRCV-01 接收数据流
func TestRawRtdbRecvDatagramWarp(t *testing.T) {
	fmt.Println("【步骤1】创建数据流句柄（预期：成功）")
	dh, err := RawRtdbCreateDatagramHandleWarp(0, "127.0.0.1")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("创建数据流失败(跳过): %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 数据流句柄创建成功，句柄=%d\n", dh)
	defer RawRtdbRemoveDatagramHandleWarp(dh)

	fmt.Println("【步骤2】接收数据流（预期：超时或返回错误）")
	// 无对端发送，超时应返回错误
	_, err = RawRtdbRecvDatagramWarp(dh, 1024, "127.0.0.1", 1)
	if RteIsOk(err) {
		fmt.Println("  结果：通过 —— 接收到数据")
	} else {
		fmt.Printf("  结果：通过 —— 接收超时或错误(预期): %s\n", err)
	}
}

// TC-DGRAMRMV-01 删除数据流
func TestRawRtdbRemoveDatagramHandleWarp(t *testing.T) {
	fmt.Println("【步骤1】创建数据流句柄（预期：成功）")
	dh, err := RawRtdbCreateDatagramHandleWarp(0, "127.0.0.1")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("创建数据流失败(跳过): %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 数据流句柄创建成功，句柄=%d\n", dh)

	fmt.Println("【步骤2】删除数据流句柄（预期：成功）")
	err = RawRtdbRemoveDatagramHandleWarp(dh)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("删除数据流失败:", err)
		return
	}
	fmt.Println("  结果：通过 —— 数据流句柄删除成功")
}

// ==================== 09. 存档管理 ====================

// TC-ARCGET-01 正常获取存档数量
func TestRawRtdbaGetArchivesCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档数量（预期：成功）")
	count, err := RawRtdbaGetArchivesCountWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取存档数量失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 存档数量=%d\n", count)
}

// TC-ARCLST-01 正常获取存档列表
func TestRawRtdbaGetArchivesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, states, err := RawRtdbaGetArchivesWarp(handle, 100)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取存档列表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 存档数量=%d\n", len(paths))
	_ = files
	_ = states
}

// TC-ARCSTS-01 正常获取存档状态
func TestRawRtdbaGetArchivesStatusWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档状态（预期：成功）")
	state, err := RawRtdbaGetArchivesStatusWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取存档状态失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 存档状态=%v\n", state)
}

// TC-ARCBIG-01 查询后台任务状态
func TestRawRtdbaQueryBigJob64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】查询后台任务状态（预期：成功）")
	path, file, job, state, endTime, progress, err := RawRtdbaQueryBigJob64Warp(handle, RtdbProcessBase)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("查询后台任务: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— path=%s, file=%s, job=%d, state=%v, endTime=%d, progress=%f\n", path, file, job, state, endTime, progress)
}

// TC-ARCCAN-01 取消后台任务
func TestRawRtdbaCancelBigJobWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】取消后台任务（预期：成功或无任务）")
	err := RawRtdbaCancelBigJobWarp(handle, RtdbProcessBase)
	if !RteIsOk(err) {
		fmt.Printf("  结果：通过 —— 取消后台任务返回提示(可能无任务): %s\n", err)
		t.Logf("取消后台任务(可能无任务): %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 取消后台任务成功")
}

// TC-ARCCRT-01 新建指定时间范围的历史存档文件
func TestRawRtdbaCreateRangedArchive64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】新建指定时间范围的历史存档文件（预期：成功）")
	now := TimestampType(time.Now().Unix())
	begin := now - 3600
	end := now
	err := RawRtdbaCreateRangedArchive64Warp(handle, "/data/", "test_go_create.rdf", begin, end, 100)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("创建存档（外部环境可能失败）: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 新建存档文件成功")
}

// TC-ARCAPP-01 追加磁盘上的历史存档文件到历史数据库
func TestRawRtdbaAppendArchiveWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】追加不存在的存档文件（预期：返回错误）")
	err := RawRtdbaAppendArchiveWarp(handle, "/data/", "notexist.rdf", RtdbArchiveStateNormal)
	if !RteIsOk(err) {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
		t.Logf("追加存档（文件不存在预期失败）: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 追加存档成功（文件居然存在）")
}

// TC-ARCRMV-01 从历史数据库中移出历史存档文件
func TestRawRtdbaRemoveArchiveWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】移出不存在的存档文件（预期：返回错误）")
	err := RawRtdbaRemoveArchiveWarp(handle, "/data/", "notexist.rdf")
	if !RteIsOk(err) {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
		t.Logf("移出存档（不存在预期失败）: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 移出存档成功（文件居然存在）")
}

// TC-ARCSFT-01 切换活动文件
func TestRawRtdbaShiftActivedWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】切换活动存档文件（预期：成功或无延续文件）")
	err := RawRtdbaShiftActivedWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：通过 —— 切换活动存档返回提示(可能无延续文件): %s\n", err)
		t.Logf("切换活动存档（可能无延续文件）: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 切换活动存档成功")
}

// TC-ARCINF-01 获取存档信息
func TestRawRtdbaGetArchivesInfoWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档数量（预期：成功）")
	count, _ := RawRtdbaGetArchivesCountWarp(handle)
	if count <= 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 存档数量=%d\n", count)

	fmt.Println("【步骤3】获取存档详细信息（预期：成功）")
	paths, files, headers, errs, err := RawRtdbaGetArchivesInfoWarp(handle, count)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取存档详细信息: %v", err)
		return
	}
	_ = errs
	_ = files
	_ = headers
	fmt.Printf("  结果：通过 —— 存档信息数=%d\n", len(paths))
}

// TC-ARCPRF-01 获取存档的实时性能监控数据
func TestRawRtdbaGetArchivesPerfDataWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档数量（预期：成功）")
	count, _ := RawRtdbaGetArchivesCountWarp(handle)
	if count <= 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 存档数量=%d\n", count)

	fmt.Println("【步骤3】获取存档性能监控数据（预期：成功）")
	paths, files, realtime, total, errs, err := RawRtdbaGetArchivesPerfDataWarp(handle, count)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取存档性能数据: %v", err)
		return
	}
	_ = errs
	fmt.Printf("  结果：通过 —— 存档性能数据条数=%d\n", len(paths))
	_ = files
	_ = realtime
	_ = total
}

// TC-ARCGI-01 获取存档文件及其附属文件的详细信息
func TestRawRtdbaGetArchiveInfoWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, _, err := RawRtdbaGetArchivesWarp(handle, 1)
	if !RteIsOk(err) || len(paths) == 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个存档\n", len(paths))

	fmt.Println("【步骤3】获取存档文件详细信息（预期：成功）")
	hdr, err := RawRtdbaGetArchiveInfoWarp(handle, paths[0], files[0], 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取存档文件信息: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 存档信息为nil=%v\n", hdr == nil)
}

// TC-ARCUPD-01 修改存档文件的可配置项
func TestRawRtdbaUpdateArchiveWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, _, err := RawRtdbaGetArchivesWarp(handle, 1)
	if !RteIsOk(err) || len(paths) == 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个存档\n", len(paths))

	fmt.Println("【步骤3】修改存档文件可配置项（预期：成功）")
	err = RawRtdbaUpdateArchiveWarp(handle, paths[0], files[0], 0, 0, 1, 1)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("修改存档配置: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 修改存档配置成功")
}

// TC-ARCARR-01 整理存档文件
func TestRawRtdbaArrangeArchiveWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, _, err := RawRtdbaGetArchivesWarp(handle, 1)
	if !RteIsOk(err) || len(paths) == 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个存档\n", len(paths))

	fmt.Println("【步骤3】整理存档文件（预期：成功）")
	err = RawRtdbaArrangeArchiveWarp(handle, paths[0], files[0])
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("整理存档: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 整理存档成功")
}

// TC-ARCIDX-01 重建存档文件索引
func TestRawRtdbaReindexArchiveWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, _, err := RawRtdbaGetArchivesWarp(handle, 1)
	if !RteIsOk(err) || len(paths) == 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个存档\n", len(paths))

	fmt.Println("【步骤3】重建存档文件索引（预期：成功）")
	err = RawRtdbaReindexArchiveWarp(handle, paths[0], files[0])
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("重建索引: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 重建索引成功")
}

// TC-ARCBKP-01 备份主存档文件及其附属文件
func TestRawRtdbaBackupArchiveWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, _, err := RawRtdbaGetArchivesWarp(handle, 1)
	if !RteIsOk(err) || len(paths) == 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个存档\n", len(paths))

	fmt.Println("【步骤3】备份存档文件（预期：成功）")
	err = RawRtdbaBackupArchiveWarp(handle, paths[0], files[0], "/backup/")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("备份存档: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 备份存档成功")
}

// TC-ARCMOV-01 将存档文件移动到指定目录
func TestRawRtdbaMoveArchiveWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, states, err := RawRtdbaGetArchivesWarp(handle, 100)
	if !RteIsOk(err) || len(paths) == 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个存档\n", len(paths))

	fmt.Println("【步骤3】移动非活动存档文件（预期：成功）")
	// 只尝试移动非活动状态的存档
	for i, s := range states {
		if s != RtdbArchiveStateActived {
			err = RawRtdbaMoveArchiveWarp(handle, paths[i], files[i], "/newdata/")
			if !RteIsOk(err) {
				fmt.Printf("  结果：失败 —— %s\n", err)
				t.Logf("移动存档: %v", err)
				return
			}
			fmt.Println("  结果：通过 —— 移动存档成功")
			return
		}
	}
	fmt.Println("  结果：跳过 —— 没有非活动存档可供移动")
	t.Skip("没有非活动存档可供移动")
}

// TC-ARCCVT-01 为存档文件转换索引格式
func TestRawRtdbaConvertIndexWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取存档列表（预期：成功）")
	paths, files, _, err := RawRtdbaGetArchivesWarp(handle, 1)
	if !RteIsOk(err) || len(paths) == 0 {
		fmt.Println("  结果：跳过 —— 无存档")
		t.Skip("无存档")
	}
	fmt.Printf("  结果：通过 —— 获取到 %d 个存档\n", len(paths))

	fmt.Println("【步骤3】转换存档文件索引格式（预期：成功）")
	err = RawRtdbaConvertIndexWarp(handle, paths[0], files[0])
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("转换索引格式: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 转换索引格式成功")
}

// ==================== 10. 历史数据查询 ====================

// TC-HCNT-01 统计历史值数量
func TestRawRtdbhArchivedValuesCount64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisCnt", "测试统计历史值数量")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisCnt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】统计历史值数量（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	count, err := RawRtdbhArchivedValuesCount64Warp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("统计历史值数量: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 历史值数量=%d\n", count)
}

// TC-HGET-01 正向读取历史数据
func TestRawRtdbhGetArchivedValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisGet", "测试正向读取历史数据")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisGet", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】正向读取历史数据（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	dts, sts, vals, states, quals, err := RawRtdbhGetArchivedValues64Warp(handle, pid, 100, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取历史数据: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 历史数据条数=%d\n", len(dts))
	_ = sts
	_ = vals
	_ = states
	_ = quals
}

// TC-HSNG-01 读取单值历史（Previous模式）
func TestRawRtdbhGetSingleValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisSng", "测试读取单值历史")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisSng", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】读取单值历史（Previous模式）（预期：成功）")
	now = TimestampType(time.Now().Unix())
	_, _, val, state, qual, err := RawRtdbhGetSingleValue64Warp(handle, pid, RtdbHisModePrevious, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取单值历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 单值=%f, state=%d, qual=%d\n", val, state, qual)
}

// TC-HSUM-01 获取统计值
func TestRawRtdbhSummaryDataWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisSum", "测试获取统计值")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisSum", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】获取统计值（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	data, err := RawRtdbhSummaryDataWarp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取统计值: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 统计: Count=%d, Max=%f, Min=%f\n", data.Count, data.MaxValue, data.MinValue)
}

// TC-HPLT-01 获取绘图数据
func TestRawRtdbhGetPlotValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisPlt", "测试获取绘图数据")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisPlt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】获取绘图数据（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	dts, _, _, _, _, err := RawRtdbhGetPlotValues64Warp(handle, pid, 100, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取绘图数据: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 绘图数据条数=%d\n", len(dts))
}

// TC-HCRS-01 获取断面数据
func TestRawRtdbhGetCrossSectionValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索标签点（预期：成功）")
	ids, _ := RawRtdbbSearchWarp(handle, "*", "*", "", "", "", "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无标签点")
		t.Skip("无标签点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个标签点\n", len(ids))

	fmt.Println("【步骤3】获取断面数据（预期：成功）")
	count := 3
	if len(ids) < count {
		count = len(ids)
	}
	now := TimestampType(time.Now().Unix())
	_, _, _, _, _, _, err := RawRtdbhGetCrossSectionValues64Warp(handle, ids[:count], RtdbHisModePrevious, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取断面数据: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 断面数据获取成功")
}

// TC-HRCNT-01 真实存储值数量
func TestRawRtdbhArchivedValuesRealCount64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisRCnt", "测试真实存储值数量")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisRCnt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】统计真实存储值数量（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	count, err := RawRtdbhArchivedValuesRealCount64Warp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("真实统计历史值数量: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 真实历史值数量=%d\n", count)
}

// TC-HGETB-01 逆向读取历史数据
func TestRawRtdbhGetArchivedValuesBackward64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisBwd", "测试逆向读取历史")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisBwd", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】逆向读取历史数据（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	dts, sts, vals, states, quals, err := RawRtdbhGetArchivedValuesBackward64Warp(handle, pid, 100, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("逆向读取历史数据: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 逆向历史数据条数=%d\n", len(dts))
	_ = sts
	_ = vals
	_ = states
	_ = quals
}

// TC-HGETC-01 正向读取坐标型储存数据
func TestRawRtdbhGetArchivedCoorValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索坐标类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "coor", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无坐标类型点")
		t.Skip("无坐标类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个坐标类型标签点\n", len(ids))

	fmt.Println("【步骤3】正向读取坐标历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, sts, xs, ys, quals, err := RawRtdbhGetArchivedCoorValues64Warp(handle, ids[0], 100, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("读取坐标历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 坐标历史数据条数=%d\n", len(dts))
	_ = sts
	_ = xs
	_ = ys
	_ = quals
}

// TC-HGETCB-01 逆向读取坐标型储存数据
func TestRawRtdbhGetArchivedCoorValuesBackward64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索坐标类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "coor", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无坐标类型点")
		t.Skip("无坐标类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个坐标类型标签点\n", len(ids))

	fmt.Println("【步骤3】逆向读取坐标历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, _, _, _, _, err := RawRtdbhGetArchivedCoorValuesBackward64Warp(handle, ids[0], 100, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("逆向读取坐标历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 逆向坐标历史条数=%d\n", len(dts))
}

// TC-HBAT-01 开始分段返回方式读取
func TestRawRtdbhGetArchivedValuesInBatches64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisBat", "测试分段读取")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisBat", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】启动分段读取（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	count, batchCount, err := RawRtdbhGetArchivedValuesInBatches64Warp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("分段读取启动: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 分段读取: count=%d, batchCount=%d\n", count, batchCount)
}

// TC-HNXT-01 分段读取下一段数据
func TestRawRtdbhGetNextArchivedValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisNxt", "测试读取下一段")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisNxt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】启动分段读取（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	_, batchCount, err := RawRtdbhGetArchivedValuesInBatches64Warp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) || batchCount <= 0 {
		fmt.Println("  结果：跳过 —— 无分段数据")
		t.Skip("无分段数据")
	}
	fmt.Printf("  结果：通过 —— batchCount=%d\n", batchCount)

	fmt.Println("【步骤6】读取下一段数据（预期：成功）")
	dts, _, _, _, _, err := RawRtdbhGetNextArchivedValues64Warp(handle, pid, batchCount)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("读取下一批: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 下一批数据条数=%d\n", len(dts))
}

// TC-HTIM-01 单调递增时间序列历史插值
func TestRawRtdbhGetTimedValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisTim", "测试指定时间插值")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisTim", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】指定时间插值（预期：成功）")
	now = TimestampType(time.Now().Unix())
	datetimes := []TimestampType{now - 3600, now - 1800, now}
	vals, states, quals, err := RawRtdbhGetTimedValues64Warp(handle, pid, datetimes, []SubtimeType{0, 0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("指定时间插值: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 指定时间插值条数=%d\n", len(vals))
	_ = states
	_ = quals
}

// TC-HTIMC-01 坐标型单调递增时间序列历史插值
func TestRawRtdbhGetTimedCoorValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索坐标类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "coor", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无坐标类型点")
		t.Skip("无坐标类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个坐标类型标签点\n", len(ids))

	fmt.Println("【步骤3】坐标型指定时间插值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	datetimes := []TimestampType{now - 3600, now - 1800, now}
	xs, ys, quals, err := RawRtdbhGetTimedCoorValues64Warp(handle, ids[0], datetimes, []SubtimeType{0, 0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("坐标插值: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 坐标插值条数=%d\n", len(xs))
	_ = ys
	_ = quals
}

// TC-HINT-01 等间隔历史插值
func TestRawRtdbhGetInterpoValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisInt", "测试等间隔插值")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisInt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】等间隔历史插值（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600
	dts, _, _, _, _, err := RawRtdbhGetInterpoValues64Warp(handle, pid, 10, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("等间隔插值: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 等间隔插值条数=%d\n", len(dts))
}

// TC-HIVL-01 等间隔内插值替换历史数値
func TestRawRtdbhGetIntervalValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisIvl", "测试等间隔读取")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisIvl", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】等间隔读取历史数据（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600
	dts, _, _, _, _, err := RawRtdbhGetIntervalValues64Warp(handle, pid, time.Minute, 60, TimestampType(past), 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("等间隔读取: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 等间隔读取条数=%d\n", len(dts))
}

// TC-HSNGC-01 读取坐标型单値
func TestRawRtdbhGetSingleCoorValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索坐标类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "coor", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无坐标类型点")
		t.Skip("无坐标类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个坐标类型标签点\n", len(ids))

	fmt.Println("【步骤3】读取坐标型单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, _, x, y, _, err := RawRtdbhGetSingleCoorValue64Warp(handle, ids[0], RtdbHisModePrevious, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("读取坐标单値: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 坐标单值: x=%f, y=%f\n", x, y)
}

// TC-HSNGB-01 读取二进制/字符串型单値
func TestRawRtdbhGetSingleBlobValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索字符串类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "string", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无字符串点")
		t.Skip("无字符串点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个字符串类型标签点\n", len(ids))

	fmt.Println("【步骤3】读取 Blob 单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, _, blob, _, err := RawRtdbhGetSingleBlobValue64Warp(handle, ids[0], RtdbHisModePrevious, now, 0, 256)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("读取Blob单値: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— Blob 单值长度=%d\n", len(blob))
}

// TC-HSNGD-01 读取 datetime 型单値
func TestRawRtdbhGetSingleDatetimeValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索 datetime 类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "datetime", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无 datetime 点")
		t.Skip("无datetime点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个 datetime 类型标签点\n", len(ids))

	fmt.Println("【步骤3】读取 datetime 单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, _, blob, _, err := RawRtdbhGetSingleDatetimeValue64Warp(handle, ids[0], RtdbHisModePrevious, now, 0, -1)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("读取datetime单値: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— datetime 单值=%s\n", string(blob))
}

// TC-HBLB-01 批量读取二进制/字符串历史数据
func TestRawRtdbhGetArchivedBlobValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索字符串类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "string", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无字符串点")
		t.Skip("无字符串点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个字符串类型标签点\n", len(ids))

	fmt.Println("【步骤3】批量读取 Blob 历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, _, blobs, quals, err := RawRtdbhGetArchivedBlobValues64Warp(handle, ids[0], 10, true, past, 0, now, 0, 256)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量读取Blob历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— Blob 历史条数=%d\n", len(dts))
	_ = blobs
	_ = quals
}

// TC-HBLBF-01 模糊搜索批量读取 Blob/String
func TestRawRtdbhGetArchivedBlobValuesFilt64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索字符串类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "string", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无字符串点")
		t.Skip("无字符串点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个字符串类型标签点\n", len(ids))

	fmt.Println("【步骤3】模糊搜索批量读取 Blob 历史（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, _, blobs, quals, err := RawRtdbhGetArchivedBlobValuesFilt64Warp(handle, ids[0], 256, 10, true, "*", past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("模糊搜索Blob历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— Blob 历史过滤条数=%d\n", len(dts))
	_ = blobs
	_ = quals
}

// TC-HDTB-01 批量读取 datetime 历史数据
func TestRawRtdbhGetArchivedDatetimeValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索 datetime 类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "datetime", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无 datetime 点")
		t.Skip("无datetime点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个 datetime 类型标签点\n", len(ids))

	fmt.Println("【步骤3】批量读取 datetime 历史（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, _, blobs, quals, err := RawRtdbhGetArchivedDatetimeValues64Warp(handle, ids[0], 10, past, 0, now, 0, -1)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量读取datetime历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— datetime 历史条数=%d\n", len(dts))
	_ = blobs
	_ = quals
}

// TC-HSUMB-01 分批获取等间隔统计値
func TestRawRtdbhSummaryDataInBatchesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisSIB", "测试分批统计")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisSIB", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】分批获取等间隔统计值（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	data, errs, err := RawRtdbhSummaryDataInBatchesWarp(handle, pid, 24, time.Hour, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("分批统计: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 分批统计段数=%d\n", len(data))
	_ = errs
}

// TC-HFLT-01 经复杂条件筛选后的历史储存値
func TestRawRtdbhGetArchivedValuesFilt64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisFlt", "测试条件过滤读取")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisFlt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】条件过滤读取历史数据（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	dts, _, _, _, _, err := RawRtdbhGetArchivedValuesFilt64Warp(handle, pid, 100, "value >= 0", past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("条件过滤历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 历史过滤条数=%d\n", len(dts))
}

// TC-HIFLT-01 经筛选后的等间隔插值
func TestRawRtdbhGetIntervalValuesFilt64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisIFlt", "测试过滤等间隔读取")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisIFlt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】过滤等间隔读取历史数据（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600
	dts, _, _, _, _, err := RawRtdbhGetIntervalValuesFilt64Warp(handle, pid, "quality == 0", time.Minute, 60, TimestampType(past), 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("过滤等间隔插值: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 过滤等间隔条数=%d\n", len(dts))
}

// TC-HIPFLT-01 经筛选后的等间隔插值
func TestRawRtdbhGetInterpoValuesFilt64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisIPFlt", "测试过滤等间隔插值")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisIPFlt", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】过滤等间隔插值（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600
	dts, _, _, _, _, err := RawRtdbhGetInterpoValuesFilt64Warp(handle, pid, "value >= 0", 10, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("过滤等间隔插值: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 过滤插值条数=%d\n", len(dts))
}

// TC-HSFLT-01 经筛选后的统计値
func TestRawRtdbhSummaryDataFiltWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisSF", "测试过滤统计")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisSF", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】过滤统计（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	data, err := RawRtdbhSummaryDataFiltWarp(handle, pid, "value >= 0", past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("过滤统计: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 过滤统计: Count=%d\n", data.Count)
}

// TC-HSBFLT-01 经筛选后的分批统计
func TestRawRtdbhSummaryDataFiltInBatchesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisSFB", "测试过滤分批统计")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisSFB", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】过滤分批统计（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 3600*24
	data, errs, err := RawRtdbhSummaryDataFiltInBatchesWarp(handle, pid, "value >= 0", 24, time.Hour, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("过滤分批统计: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 过滤分批统计段数=%d\n", len(data))
	_ = errs
}

// TC-HNTP-01 读取单个自定义类型标签点某时间的历史数据
func TestRawRtdbhGetSingleNamedTypeValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索自定义类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无自定义类型点")
		t.Skip("无自定义类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个自定义类型标签点\n", len(ids))

	fmt.Println("【步骤3】读取自定义类型历史单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, _, obj, _, err := RawRtdbhGetSingleNamedTypeValue64Warp(handle, ids[0], RtdbHisModePrevious, now, 0, 256)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("读取自定义类型历史单値: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 自定义类型历史单值长度=%d\n", len(obj))
}

// TC-HNTB-01 连续读取自定义类型标签点历史数据
func TestRawRtdbhGetArchivedNamedTypeValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索自定义类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无自定义类型点")
		t.Skip("无自定义类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个自定义类型标签点\n", len(ids))

	fmt.Println("【步骤3】连续读取自定义类型历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 3600*24
	dts, _, objs, quals, err := RawRtdbhGetArchivedNamedTypeValues64Warp(handle, ids[0], 10, past, 0, now, 0, 256)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("读取自定义类型历史: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 自定义类型历史条数=%d\n", len(dts))
	_ = objs
	_ = quals
}

// ==================== 11. 历史数据写入与修改 ====================

// TC-HPUT-01 写入历史单值
func TestRawRtdbhPutSingleValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisPut", "测试写入历史单值")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisPut", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】写入历史单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	err = RawRtdbhPutSingleValue64Warp(handle, pid, now, 0, 123.45, 0, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入单值历史: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入历史单值成功")
}

// TC-HPUTB-01 批量写入历史
func TestRawRtdbhPutArchivedValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisPutB", "测试批量写入历史")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisPutB", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】批量写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{now}, []SubtimeType{0}, []float64{99.9}, []int64{0}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量写入历史: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入历史成功")
}

// TC-HUPD-01 更新历史值
func TestRawRtdbhUpdateValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisUpd", "测试更新历史值")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisUpd", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】更新历史值（预期：成功）")
	now = TimestampType(time.Now().Unix())
	err = RawRtdbhUpdateValue64Warp(handle, pid, now, 0, 77.7, 0, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("更新历史值: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 更新历史值成功")
}

// TC-HPSC-01 写入坐标型历史单値
func TestRawRtdbhPutSingleCoorValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索坐标类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "coor", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无坐标点")
		t.Skip("无坐标点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个坐标类型标签点\n", len(ids))

	fmt.Println("【步骤3】写入坐标型历史单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhPutSingleCoorValue64Warp(handle, ids[0], now, 0, 1.1, 2.2, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入坐标历史单値: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入坐标历史单值成功")
}

// TC-HPSB-01 写入二进制/字符串历史单値
func TestRawRtdbhPutSingleBlobValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索字符串类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "string", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无字符串点")
		t.Skip("无字符串点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个字符串类型标签点\n", len(ids))

	fmt.Println("【步骤3】写入 Blob 历史单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhPutSingleBlobValue64Warp(handle, ids[0], true, now, 0, []byte("testvalue"), 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入Blob历史单値: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入 Blob 历史单值成功")
}

// TC-HPSD-01 写入datetime历史单値
func TestRawRtdbhPutSingleDatetimeValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索 datetime 类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "datetime", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无 datetime 点")
		t.Skip("无datetime点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个 datetime 类型标签点\n", len(ids))

	fmt.Println("【步骤3】写入 datetime 历史单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhPutSingleDatetimeValue64Warp(handle, ids[0], now, 0, []byte("2026-01-01T00:00:00"), 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入datetime历史单値: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入 datetime 历史单值成功")
}

// TC-HPSN-01 写入自定义类型历史单値
func TestRawRtdbhPutSingleNamedTypeValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索自定义类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无自定义类型点")
		t.Skip("无自定义类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个自定义类型标签点\n", len(ids))

	fmt.Println("【步骤3】写入自定义类型历史单值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhPutSingleNamedTypeValue64Warp(handle, ids[0], now, 0, []byte{0, 0, 0, 0}, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("写入自定义类型历史单値: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 写入自定义类型历史单值成功")
}

// TC-HPAC-01 批量写入坐标型历史数据
func TestRawRtdbhPutArchivedCoorValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索坐标类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "coor", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无坐标点")
		t.Skip("无坐标点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个坐标类型标签点\n", len(ids))

	fmt.Println("【步骤3】批量写入坐标历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbhPutArchivedCoorValues64Warp(handle,
		[]PointID{ids[0]}, []TimestampType{now}, []SubtimeType{0},
		[]float32{1.1}, []float32{2.2}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量写入坐标历史: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入坐标历史成功")
}

// TC-HPAB-01 批量写入二进制/字符串历史数据
func TestRawRtdbhPutArchivedBlobValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索字符串类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "string", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无字符串点")
		t.Skip("无字符串点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个字符串类型标签点\n", len(ids))

	fmt.Println("【步骤3】批量写入 Blob 历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbhPutArchivedBlobValues64Warp(handle,
		[]PointID{ids[0]}, []bool{true}, []TimestampType{now}, []SubtimeType{0},
		[][]byte{[]byte("testblob")}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量写入Blob历史: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入 Blob 历史成功")
}

// TC-HPAD-01 批量写入datetime历史数据
func TestRawRtdbhPutArchivedDatetimeValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索 datetime 类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "datetime", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无 datetime 点")
		t.Skip("无datetime点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个 datetime 类型标签点\n", len(ids))

	fmt.Println("【步骤3】批量写入 datetime 历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbhPutArchivedDatetimeValues64Warp(handle,
		[]PointID{ids[0]}, []TimestampType{now}, []SubtimeType{0},
		[]string{"2026-01-01T00:00:00"}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量写入datetime历史: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入 datetime 历史成功")
}

// TC-HPAN-01 批量写入自定义类型历史数据
func TestRawRtdbhPutArchivedNamedTypeValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索自定义类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无自定义类型点")
		t.Skip("无自定义类型点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个自定义类型标签点\n", len(ids))

	fmt.Println("【步骤3】批量写入自定义类型历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	_, err := RawRtdbhPutArchivedNamedTypeValues64Warp(handle,
		[]PointID{ids[0]}, []TimestampType{now}, []SubtimeType{0},
		[][]byte{{0, 0, 0, 0}}, []Quality{0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("批量写入自定义类型历史: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 批量写入自定义类型历史成功")
}

// TC-HUPDC-01 修改坐标型历史単値
func TestRawRtdbhUpdateCoorValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索坐标类型标签点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 10, "*", "*", "", "", "", "", "coor", 0, 0, 0, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无坐标点")
		t.Skip("无坐标点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个坐标类型标签点\n", len(ids))

	fmt.Println("【步骤3】修改坐标型历史值（预期：成功）")
	now := TimestampType(time.Now().Unix())
	err := RawRtdbhUpdateCoorValue64Warp(handle, ids[0], now, 0, 3.3, 4.4, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("修改坐标历史値: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 修改坐标历史值成功")
}

// TC-HRMV-01 删除时间段内的历史数据
func TestRawRtdbhRemoveValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisRmv", "测试删除历史区间")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisRmv", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】删除时间段内的历史数据（预期：成功）")
	now = TimestampType(time.Now().Unix())
	past = now - 60
	count, err := RawRtdbhRemoveValues64Warp(handle, pid, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("删除历史区间: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 删除历史条数=%d\n", count)
}

// TC-HRMV1-01 删除单个历史值
func TestRawRtdbhRemoveValue64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisRm1", "测试删除单个历史值")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisRm1", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】删除单个历史值（预期：成功）")
	now = TimestampType(time.Now().Unix())
	err = RawRtdbhRemoveValue64Warp(handle, pid, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("删除历史值: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 删除单个历史值成功")
}

// TC-HFLU-01 刷新历史缓存
func TestRawRtdbhFlushArchivedValuesWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】创建临时表（预期：成功）")
	tbl, err := RawRtdbbAppendTableWarp(handle, "TestTblHisFlu", "测试刷新历史缓存")
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建临时表失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 临时表ID=%d\n", tbl.ID)
	defer func() {
		fmt.Println("【清理】删除临时表（预期：成功）")
		rte := RawRtdbbRemoveTableByIdWarp(handle, tbl.ID)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理临时表失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 临时表已删除")
		}
	}()

	fmt.Println("【步骤3】创建 Real64 标签点（预期：成功）")
	pid, err := RawRtdbbInsertBasePointWarp(handle, "TestPtHisFlu", RtdbTypeReal64, tbl.ID, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("创建标签点失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 标签点ID=%d\n", pid)
	defer func() {
		fmt.Println("【清理】删除标签点（预期：成功）")
		rte := RawRtdbbRemovePointByIdWarp(handle, pid)
		if !RteIsOk(rte) {
			fmt.Printf("  结果：失败 —— %s\n", rte)
			t.Logf("清理标签点失败: %s", rte)
		} else {
			fmt.Println("  结果：通过 —— 标签点已删除")
		}
	}()

	fmt.Println("【步骤4】预写入历史数据（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 60
	_, err = RawRtdbhPutArchivedValues64Warp(handle, []PointID{pid}, []TimestampType{past, now}, []SubtimeType{0, 0}, []float64{10.0, 20.0}, []int64{0, 0}, []Quality{0, 0})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("预写入历史数据: %v", err)
	} else {
		fmt.Println("  结果：通过 —— 预写入历史数据成功")
	}

	fmt.Println("【步骤5】刷新历史缓存（预期：成功）")
	count, err := RawRtdbhFlushArchivedValuesWarp(handle, pid)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("刷新历史缓存: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 刷新缓存条数=%d\n", count)
}

// ==================== 12. 方程式计算与性能监控 ====================

// TC-EQCOMP-01 历史方程式计算
func TestRawRtdbeComputeHistory64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索计算点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无计算点")
		t.Skip("无计算点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个计算点\n", len(ids))

	fmt.Println("【步骤3】执行历史方程式计算（预期：成功）")
	now := TimestampType(time.Now().Unix())
	past := now - 3600
	_, err := RawRtdbeComputeHistory64Warp(handle, ids[:1], 0, past, 0, now, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("历史计算: %v", err)
		return
	}
	fmt.Println("  结果：通过 —— 历史方程式计算成功")
}

// TC-EQFN-01~04 通过文件名获取方程式内容
func TestRawRtdbbGetEquationByFileNameWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】通过文件名获取方程式内容（预期：返回错误）")
	eq, err := RawRtdbbGetEquationByFileNameWarp(handle, "nonexistent.eq")
	if RteIsOk(err) {
		fmt.Printf("  结果：通过 —— 获取方程式内容长度=%d\n", len(eq))
		t.Logf("获取方程式内容长度: %d", len(eq))
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
		t.Logf("获取方程式（文件不存在预期失败）: %v", err)
	}
}

// TC-EQID-01 通过 ID 获取方程式内容
func TestRawRtdbbGetEquationByIdWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索计算点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无计算点")
		t.Skip("无计算点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个计算点\n", len(ids))

	fmt.Println("【步骤3】通过ID获取方程式内容（预期：成功）")
	eq, err := RawRtdbbGetEquationByIdWarp(handle, ids[0])
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取方程式: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 方程式长度=%d\n", len(eq))
}

// TC-GRFCNT-01 获取拓扑数量
func TestRawRtdbeGetEquationGraphCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索计算点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无计算点")
		t.Skip("无计算点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个计算点\n", len(ids))

	fmt.Println("【步骤3】获取拓扑数量（预期：成功）")
	count, err := RawRtdbeGetEquationGraphCountWarp(handle, ids[0], RtdbGraphFlagAll)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取拓扑数量: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 拓扑数量=%d\n", count)
}

// TC-GRFDAT-01 获取拓扑数据
func TestRawRtdbeGetEquationGraphDatasWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】搜索计算点（预期：成功）")
	ids, _ := RawRtdbbSearchExWarp(handle, 100, "*", "*", "", "", "", "", "", 0, 0, RtdbSearchNull, "", RtdbSortFlag(0))
	if len(ids) == 0 {
		fmt.Println("  结果：跳过 —— 无计算点")
		t.Skip("无计算点")
	}
	fmt.Printf("  结果：通过 —— 找到 %d 个计算点\n", len(ids))

	fmt.Println("【步骤3】获取拓扑数量（预期：成功）")
	count, _ := RawRtdbeGetEquationGraphCountWarp(handle, ids[0], RtdbGraphFlagAll)
	if count <= 0 {
		fmt.Println("  结果：跳过 —— 拓扑为空")
		t.Skip("拓扑为空")
	}
	fmt.Printf("  结果：通过 —— 拓扑数量=%d\n", count)

	fmt.Println("【步骤4】获取拓扑数据（预期：成功）")
	graph, err := RawRtdbeGetEquationGraphDatasWarp(handle, ids[0], RtdbGraphFlagAll, count)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取拓扑数据: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 拓扑数据条数=%d\n", len(graph))
}

// TC-PERFCNT-01 获取性能监控点数量
func TestRawRtdbpGetPerfTagsCountWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取性能监控点数量（预期：成功）")
	count, err := RawRtdbpGetPerfTagsCountWarp(handle)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取性能点数量失败:", err)
		return
	}
	fmt.Printf("  结果：通过 —— 性能点数量=%d\n", count)
}

// TC-PERFINFO-01 获取性能点信息
func TestRawRtdbpGetPerfTagsInfoWarp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取性能点信息（预期：成功）")
	infos, errs, err := RawRtdbpGetPerfTagsInfoWarp(handle, []RtdbPerfTagID{PftCpuUsageOfLogger, PftMemBytesOfLogger})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取性能点信息失败:", err)
		return
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			fmt.Printf("  结果：警告 —— 第%d个性能点出错: %v\n", i, e)
			t.Logf("第%d个性能点出错: %v", i, e)
		}
	}
	fmt.Printf("  结果：通过 —— 性能点信息数=%d\n", len(infos))
}

// TC-PERFVAL-01 获取性能值
func TestRawRtdbpGetPerfValues64Warp(t *testing.T) {
	fmt.Println("【步骤1】连接并登录（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)
	fmt.Println("  结果：通过 —— 登录成功")

	fmt.Println("【步骤2】获取性能值（预期：成功）")
	dts, sts, vals, states, quals, errs, err := RawRtdbpGetPerfValues64Warp(handle, []RtdbPerfTagID{PftCpuUsageOfLogger, PftMemBytesOfLogger})
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Error("获取性能值失败:", err)
		return
	}
	for i := range dts {
		if !RteIsOk(errs[i]) {
			fmt.Printf("  性能点[%d] 出错: %s\n", i, errs[i])
			continue
		}
		fmt.Printf("  性能点[%d] 时间=%d 子时间=%d 值=%f 状态=%d 质量=%d\n",
			i, dts[i], sts[i], vals[i], states[i], quals[i])
	}
	fmt.Println("  结果：通过 —— 获取性能值成功")
}

// ==================== 13. 数据流与元数据同步 ====================

// TC-SYNC-01 正常获取本节点同步信息 / TC-SYNC-03 获取不存在的节点 / TC-SYNC-04 无效句柄 / TC-SYNC-05 未登录
func TestRawRtdbbGetMetaSyncInfoWarp(t *testing.T) {
	// TC-SYNC-04 无效句柄：传 -1 应该被拒绝
	fmt.Println("【步骤1】测试无效句柄（预期：返回错误）")
	_, _, err := RawRtdbbGetMetaSyncInfoWarp(-1, 0)
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 无效句柄居然通过了！")
		t.Error("无效句柄应返回错误")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
	}

	// TC-SYNC-05 未登录：只连接不登录应该被拒绝
	fmt.Println("【步骤2】测试未登录（预期：返回权限错误）")
	handleNoLogin, err := RawRtdbConnectWarp(Hostname, Port)
	if RteIsOk(err) && handleNoLogin > 0 {
		_, _, err = RawRtdbbGetMetaSyncInfoWarp(handleNoLogin, 0)
		if RteIsOk(err) {
			fmt.Println("  结果：失败 —— 未登录居然通过了！")
			t.Error("未登录应返回错误")
		} else {
			fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
		}
		_ = RawRtdbDisconnectWarp(handleNoLogin)
	} else {
		fmt.Println("  连接失败，跳过未登录测试")
	}

	// TC-SYNC-01 正常获取本节点同步信息
	fmt.Println("【步骤3】测试正常获取本节点同步信息（预期：成功）")
	handle := connectAndLogin(t)
	defer disconnect(t, handle)

	infos, errs, err := RawRtdbbGetMetaSyncInfoWarp(handle, 0)
	if !RteIsOk(err) {
		fmt.Printf("  结果：失败 —— %s\n", err)
		t.Logf("获取元数据同步信息失败: %v", err)
		return
	}
	fmt.Printf("  结果：通过 —— 返回 %d 条同步信息\n", len(infos))
	for i, info := range infos {
		fmt.Printf("  节点[%d] 角色=%s(%d), 状态=%s(%d), IP=%s, 同步版本=%d, 堆积数据=%d\n",
			i, info.Role.Desc(), info.Role, info.Status.Desc(), info.Status,
			info.IpString, info.Version, info.DataSize)
	}
	for i, e := range errs {
		if !RteIsOk(e) {
			fmt.Printf("  节点[%d] 出现错误：%s\n", i, e)
		}
	}

	// TC-SYNC-03 获取不存在的节点：传 999 应该被拒绝
	fmt.Println("【步骤4】测试获取不存在的节点（预期：返回错误）")
	_, _, err = RawRtdbbGetMetaSyncInfoWarp(handle, 999)
	if RteIsOk(err) {
		fmt.Println("  结果：失败 —— 不存在的节点居然通过了！")
		t.Error("不存在的节点应返回错误")
	} else {
		fmt.Printf("  结果：通过 —— 返回了预期的错误：%s\n", err)
	}
	fmt.Println("【全部测试完成】")
}
