package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	var deleteFlag bool
	flag.BoolVar(&deleteFlag, "d", false, "删除指定版本")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("用法:")
		fmt.Println("  发布版本: publish_version <版本号>")
		fmt.Println("  删除版本: publish_version -d <版本号>")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Println("  publish_version v4.0.15_0.2.0")
		fmt.Println("  publish_version -d v4.0.15_0.2.0")
		os.Exit(1)
	}

	version := args[0]

	// 自动获取项目根目录（publish_version.go 位于 tools/ 子目录下）
	rootDir := getProjectRoot()

	if deleteFlag {
		deleteVersion(rootDir, version)
	} else {
		publishVersion(rootDir, version)
	}
}

// publishVersion 发布版本
func publishVersion(rootDir, version string) {
	fmt.Printf("========== 发布版本: %s ==========\n", version)

	// 步骤1：调用 gen_api_list 更新 API 文档
	fmt.Println("[1/4] 更新 API 文档 ...")
	if err := runGenAPIList(rootDir); err != nil {
		fmt.Println("更新 API 文档失败:", err)
		os.Exit(1)
	}
	fmt.Println("API 文档更新完成")

	// 步骤2：如果有变更，自动提交
	fmt.Println("[2/4] 检查并提交文档变更 ...")
	if hasChanges(rootDir) {
		if err := gitAdd(rootDir); err != nil {
			fmt.Println("git add 失败:", err)
			os.Exit(1)
		}
		if err := gitCommit(rootDir, fmt.Sprintf("chore: update api docs for %s", version)); err != nil {
			fmt.Println("git commit 失败:", err)
			os.Exit(1)
		}
		fmt.Println("文档变更已提交")
	} else {
		fmt.Println("没有需要提交的变更")
	}

	// 步骤3：本地打 tag
	fmt.Println("[3/4] 本地打 tag ...")
	if err := gitTag(rootDir, version); err != nil {
		fmt.Println("本地打 tag 失败:", err)
		os.Exit(1)
	}
	fmt.Printf("本地 tag %s 创建成功\n", version)

	// 步骤4：推送 tag 到远程
	fmt.Println("[4/4] 推送 tag 到远程 ...")
	if err := gitPushTag(rootDir, version); err != nil {
		fmt.Println("推送 tag 到远程失败:", err)
		os.Exit(1)
	}
	fmt.Printf("远程 tag %s 推送成功\n", version)

	fmt.Println("")
	fmt.Printf("版本 %s 发布完成！\n", version)
}

// deleteVersion 删除版本
func deleteVersion(rootDir, version string) {
	fmt.Printf("========== 删除版本: %s ==========\n", version)

	// 步骤1：删除本地 tag
	fmt.Println("[1/2] 删除本地 tag ...")
	if err := gitDeleteLocalTag(rootDir, version); err != nil {
		// 如果本地 tag 不存在，只打印警告，继续执行
		if strings.Contains(err.Error(), "does not exist") {
			fmt.Printf("警告: 本地 tag %s 不存在\n", version)
		} else {
			fmt.Println("删除本地 tag 失败:", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("本地 tag %s 已删除\n", version)
	}

	// 步骤2：删除远程 tag
	fmt.Println("[2/2] 删除远程 tag ...")
	if err := gitDeleteRemoteTag(rootDir, version); err != nil {
		fmt.Println("删除远程 tag 失败:", err)
		os.Exit(1)
	}
	fmt.Printf("远程 tag %s 已删除\n", version)

	fmt.Println("")
	fmt.Printf("版本 %s 删除完成！\n", version)
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 获取当前可执行文件或源码所在目录
	_, filename, _, _ := runtime.Caller(0)
	// 当前文件在 tools/ 目录下，上上级目录即为项目根目录
	// runtime.Caller(0) 返回的是源码文件路径，其所在目录是 tools/
	// 但 go run 时工作目录可能不同，所以通过文件位置向上推导
	toolsDir := filepath.Dir(filename)
	return filepath.Dir(toolsDir)
}

// runGenAPIList 调用 gen_api_list.go 更新文档
func runGenAPIList(rootDir string) error {
	toolsDir := filepath.Join(rootDir, "tools")
	cmd := exec.Command("go", "run", "gen_api_list.go")
	cmd.Dir = toolsDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// hasChanges 检查工作区是否有未提交的变更
func hasChanges(rootDir string) bool {
	cmd := exec.Command("git", "status", "--short")
	cmd.Dir = rootDir
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}

// gitAdd 执行 git add .
func gitAdd(rootDir string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitCommit 执行 git commit
func gitCommit(rootDir, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitTag 本地创建 tag
func gitTag(rootDir, version string) error {
	cmd := exec.Command("git", "tag", version)
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitPushTag 推送 tag 到远程
func gitPushTag(rootDir, version string) error {
	cmd := exec.Command("git", "push", "origin", version)
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitDeleteLocalTag 删除本地 tag
func gitDeleteLocalTag(rootDir, version string) error {
	cmd := exec.Command("git", "tag", "-d", version)
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// gitDeleteRemoteTag 删除远程 tag
func gitDeleteRemoteTag(rootDir, version string) error {
	cmd := exec.Command("git", "push", "--delete", "origin", version)
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
