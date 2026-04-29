@echo off
chcp 65001 >nul

REM RTDB API 发版脚本（Windows CMD）

set DELETE_FLAG=false
set VERSION=

REM 解析参数
if "%~1"=="-d" (
    set DELETE_FLAG=true
    set VERSION=%~2
) else (
    set VERSION=%~1
)

REM 参数校验
if "%VERSION%"=="" (
    echo 用法:
    echo   发布版本: publish.bat ^<版本号^>
    echo   删除版本: publish.bat -d ^<版本号^>
    echo.
    echo 示例:
    echo   publish.bat v4.0.15_0.2.0
    echo   publish.bat -d v4.0.15_0.2.0
    exit /b 1
)

REM 获取脚本所在目录（项目根目录）
cd /d "%~dp0"

REM 调用 Go 发版工具
if "%DELETE_FLAG%"=="true" (
    echo ========== 删除版本: %VERSION% ==========
    go run tools\publish_version.go -d %VERSION%
) else (
    echo ========== 发布版本: %VERSION% ==========
    go run tools\publish_version.go %VERSION%
)
