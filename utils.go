package rtdb_api

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"math/rand"
	"runtime"
	"unsafe"
)

func CCharArrayToString(p *C.char, n int) string {
	b := C.GoBytes(unsafe.Pointer(p), C.int(n))
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func GoStringToCCharArray(s string, p *C.char, n int) {
	if p == nil || n <= 0 {
		return
	}

	dst := unsafe.Slice((*byte)(unsafe.Pointer(p)), n)

	b := []byte(s)
	if len(b) >= n {
		b = b[:n-1]
	}

	copy(dst, b)
	dst[len(b)] = 0

	// 可选：清零剩余空间
	for i := len(b) + 1; i < n; i++ {
		dst[i] = 0
	}
}

func RtdbErrorListToErrorList(errs []RtdbError) []error {
	rtn := make([]error, 0)
	for _, err := range errs {
		rtn = append(rtn, err.GoError())
	}
	return rtn
}

// SafeSlice 安全获取切片子集，自动处理越界问题
func SafeSlice[T any](slice []T, start, count int32) []T {
	// 处理空切片或无效参数
	if slice == nil || start < 0 || count <= 0 {
		return []T{}
	}

	// 转换为 int 类型便于操作
	s := int(start)
	c := int(count)
	length := len(slice)

	// 检查起始位置是否超出范围
	if s >= length {
		return []T{}
	}

	// 计算实际结束位置
	end := s + c
	if end > length {
		end = length
	}

	// 返回有效的子切片
	return slice[s:end]
}

// BoolToInt64 bool转换为Int
func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

// Int64ToBool Int转换为bool
func Int64ToBool(v int64) bool {
	if v == 0 {
		return false
	} else {
		return true
	}
}

// ClientIsWindows 客户端是Windows
func ClientIsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

// GBKBytesToString GBK格式的bytes，转换为UTF8 string
func GBKBytesToString(gbkBytes []byte) (string, error) {
	// 创建GBK解码器
	decoder := simplifiedchinese.GBK.NewDecoder()

	// 使用transform进行解码转换
	reader := transform.NewReader(bytes.NewReader(gbkBytes), decoder)

	// 读取解码后的数据
	utf8Bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("GBK转UTF-8失败: %w", err)
	}

	return string(utf8Bytes), nil
}

// StringToGBKBytes UTF8格式的string转换为GBK格式的bytes
func StringToGBKBytes(str string) ([]byte, error) {
	encoder := simplifiedchinese.GBK.NewEncoder()
	buf, n, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		return nil, errors.New("str转换成GBK格式[]byte报错：" + err.Error())
	}
	return buf[:n], nil
}

// StringInDB 字符串向数据库输入
func StringInDB(str string) string {
	if ClientIsWindows() {
		data, _ := StringToGBKBytes(str)
		return string(data)
	} else {
		return str
	}
}

// StringOutDB 数据库字符串输出
func StringOutDB(str string) string {
	if ClientIsWindows() {
		data, _ := GBKBytesToString([]byte(str))
		return data
	} else {
		return str
	}
}

// RandString 生成随机字符串
func RandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// SafeCopyToSlice C指针转换为Go切片
func SafeCopyToSlice[T any](ptr unsafe.Pointer, count int) []T {
	if ptr == nil || count <= 0 {
		return nil
	}

	// 创建新的Go切片
	result := make([]T, count)

	// 计算总字节数
	var zero T
	elementSize := int(unsafe.Sizeof(zero))
	totalSize := count * elementSize

	// 复制内存数据
	src := unsafe.Slice((*byte)(ptr), totalSize)
	dst := unsafe.Slice((*byte)(unsafe.Pointer(&result[0])), totalSize)
	copy(dst, src)

	return result
}
