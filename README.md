# rtdb_api

## 注意

由于本库使用了CGO技术，需要编译一定量的C代码，因此编译时间较长，尤其是在Windows环境下，请耐心等待

## 编译环境

* 在Linux平台中，需安装Golang编译器、GCC工具链
* 在Windows平台中，需安装Golang编译器、MicrosoftVisualC++BuildTools工具链(可通过VisualStudio2022安装)

## 版本

格式：Go库版本号+So库版本号

* v0.1.x-4.0.11: 基于SO库4.0.11，这个版本已经定版了，后续除了bug fix不会有其余API的增减
* v0.2.x-4.0.15: 新增了一些Raw函数, SO库升级到4.0.15，由于API变更所以记作小版本号0.2.0

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

	"github.com/golden-datas/rtdb-go-sdk"
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

## go.mod引用版本说明

引用本库的时候需注意版本号，请使用最新版本, 可执行 ```git ls-remote --tags git@github.com/golden-datas/rtdb-go-sdk.git```
命令行查看最新版本号

```text
require github.com/golden-datas/rtdb-go-sdk v0.2.4-4.0.15
// 注意这里的 v0.2.4-4.0.15 ， 需替换成最新版本
```

## 版本发布工具

项目根目录下提供了跨平台发版脚本，发布前会自动更新 API 文档并推送。

### Linux / macOS / Git Bash

```bash
# 发布版本
./publish.sh v0.2.4-4.0.15

# 删除版本
./publish.sh -d v0.2.4-4.0.15
```

### Windows PowerShell

```powershell
# 发布版本
.\publish.ps1 v0.2.4-4.0.15

# 删除版本
.\publish.ps1 -d v0.2.4-4.0.15
```

## 可用函数列表

详见 [可用API列表.md](可用API列表.md)

## API文档

详见 [APIDOC.md](APIDOC.md)
