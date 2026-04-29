#!/bin/bash

# RTDB API 发版脚本
# 支持 Linux / macOS / Windows Git Bash

set -e

DELETE_FLAG=false
VERSION=""

# 解析参数
if [ "$1" == "-d" ]; then
    DELETE_FLAG=true
    VERSION="$2"
else
    VERSION="$1"
fi

# 参数校验
if [ -z "$VERSION" ]; then
    echo "用法:"
    echo "  发布版本: ./publish.sh <版本号>"
    echo "  删除版本: ./publish.sh -d <版本号>"
    echo ""
    echo "示例:"
    echo "  ./publish.sh v4.0.15_0.2.0"
    echo "  ./publish.sh -d v4.0.15_0.2.0"
    exit 1
fi

# 获取脚本所在目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# 调用 Go 发版工具
if [ "$DELETE_FLAG" == true ]; then
    echo "========== 删除版本: $VERSION =========="
    go run tools/publish_version.go -d "$VERSION"
else
    echo "========== 发布版本: $VERSION =========="
    go run tools/publish_version.go "$VERSION"
fi
