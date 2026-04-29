package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type FuncInfo struct {
	Name string
	Desc string
}

func main() {
	rootDir := `C:\Users\sailyond\Desktop\rtdb_api`
	easyPath := filepath.Join(rootDir, "easy.go")
	apiPath := filepath.Join(rootDir, "api.go")
	outputPath := filepath.Join(rootDir, "可用API列表.md")
	apiDocPath := filepath.Join(rootDir, "APIDOC.md")

	// ========== 步骤1：检测/安装 gomarkdoc ==========
	fmt.Println("步骤1：检测 API 文档生成工具 gomarkdoc ...")
	gomarkdocPath := ensureGomarkdoc()
	fmt.Println("gomarkdoc 路径:", gomarkdocPath)

	// ========== 步骤2：生成 APIDOC.md ==========
	fmt.Println("步骤2：生成 APIDOC.md ...")
	generateAPIDoc(gomarkdocPath, rootDir, apiDocPath)
	fmt.Println("APIDOC.md 生成成功:", apiDocPath)

	// ========== 步骤3：生成 可用API列表.md ==========
	fmt.Println("步骤3：生成 可用API列表.md ...")

	// 解析 Easy 函数
	easyFuncs := parseEasyFuncs(easyPath)

	// 解析 Raw 函数
	rawFuncs := parseRawFuncs(apiPath)

	// 生成 Markdown
	var sb strings.Builder
	sb.WriteString("# 可用 API 列表\n\n")

	sb.WriteString("## Easy API 列表\n\n")
	for _, f := range easyFuncs {
		sb.WriteString(fmt.Sprintf("1. %s  %s\n", f.Name, f.Desc))
	}

	sb.WriteString("\n## Raw API 列表\n\n")
	for _, f := range rawFuncs {
		sb.WriteString(fmt.Sprintf("1. %s  %s\n", f.Name, f.Desc))
	}

	err := os.WriteFile(outputPath, []byte(sb.String()), 0644)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		os.Exit(1)
	}
	fmt.Println("可用API列表.md 生成成功:", outputPath)
}

// ensureGomarkdoc 检测 gomarkdoc 是否已安装，没有则自动安装
func ensureGomarkdoc() string {
	// 先尝试在 PATH 中查找
	if path, err := exec.LookPath("gomarkdoc"); err == nil {
		return path
	}

	// 再尝试在 GOPATH/bin 中查找
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// 通过 go env GOPATH 获取
		out, err := exec.Command("go", "env", "GOPATH").Output()
		if err == nil {
			gopath = strings.TrimSpace(string(out))
		}
	}
	if gopath != "" {
		binPath := filepath.Join(gopath, "bin")
		exeName := "gomarkdoc"
		if runtime.GOOS == "windows" {
			exeName = "gomarkdoc.exe"
		}
		candidate := filepath.Join(binPath, exeName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// 未找到，自动安装
	fmt.Println("gomarkdoc 未找到，正在自动安装 ...")
	cmd := exec.Command("go", "install", "github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("安装 gomarkdoc 失败:", err)
		os.Exit(1)
	}

	// 安装完成后再次查找
	if path, err := exec.LookPath("gomarkdoc"); err == nil {
		return path
	}
	if gopath != "" {
		binPath := filepath.Join(gopath, "bin")
		exeName := "gomarkdoc"
		if runtime.GOOS == "windows" {
			exeName = "gomarkdoc.exe"
		}
		candidate := filepath.Join(binPath, exeName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	fmt.Println("安装后仍无法找到 gomarkdoc，请检查 GOPATH/bin 是否在 PATH 中")
	os.Exit(1)
	return ""
}

// generateAPIDoc 调用 gomarkdoc 生成 APIDOC.md
func generateAPIDoc(gomarkdocPath, rootDir, outputPath string) {
	cmd := exec.Command(gomarkdocPath, "-o", outputPath, ".")
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("生成 APIDOC.md 失败:", err)
		os.Exit(1)
	}
}

func parseEasyFuncs(path string) []FuncInfo {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const delimiter = "下面是RtdbConnect函数"
	foundDelimiter := false

	var results []FuncInfo
	var currentComment []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !foundDelimiter {
			if strings.Contains(line, delimiter) {
				foundDelimiter = true
			}
			continue
		}

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") {
			currentComment = append(currentComment, trimmed)
			continue
		}

		if strings.HasPrefix(trimmed, "func ") {
			funcName := extractFuncName(trimmed)
			if funcName != "" {
				desc := extractDesc(funcName, currentComment)
				results = append(results, FuncInfo{
					Name: funcName,
					Desc: desc,
				})
			}
			currentComment = nil
			continue
		}

		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "func ") {
			currentComment = nil
		}
	}

	return results
}

func parseRawFuncs(path string) []FuncInfo {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const delimiter = "下面是函数实现"
	foundDelimiter := false

	var results []FuncInfo
	var currentComment []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !foundDelimiter {
			if strings.Contains(line, delimiter) {
				foundDelimiter = true
			}
			continue
		}

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") {
			currentComment = append(currentComment, trimmed)
			continue
		}

		if strings.HasPrefix(trimmed, "func ") {
			funcName := extractFuncName(trimmed)
			if funcName != "" && strings.HasPrefix(funcName, "Raw") {
				desc := extractDesc(funcName, currentComment)
				results = append(results, FuncInfo{
					Name: funcName,
					Desc: desc,
				})
			}
			currentComment = nil
			continue
		}

		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "func ") {
			currentComment = nil
		}
	}

	return results
}

func extractFuncName(line string) string {
	re := regexp.MustCompile(`^func\s+(?:\([^)]*\)\s+)?(\w+)\s*`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func extractDesc(funcName string, comments []string) string {
	if len(comments) == 0 {
		return ""
	}
	first := comments[0]
	first = strings.TrimPrefix(first, "//")
	first = strings.TrimSpace(first)
	first = strings.TrimPrefix(first, funcName)
	first = strings.TrimSpace(first)
	return first
}
