@echo off
setlocal enabledelayedexpansion

echo ============================================
echo 开始自动化构建流程
echo ============================================

REM 步骤1：删除旧文件
if exist webook (
    del webook
    echo ✅ 已删除旧文件 webook
) else (
    echo ℹ️ 未找到 webook 文件，跳过删除
)

REM 步骤2：执行 go mod tidy
echo 正在执行 go mod tidy...
go mod tidy
if !errorlevel! neq 0 (
    echo ❌ go mod tidy 执行失败
    goto restore
)
echo ✅ go mod tidy 完成

REM 步骤3：设置环境变量为Linux
echo 正在设置 GOOS=linux...
go env -w GOOS=linux
if !errorlevel! neq 0 (
    echo ❌ 环境变量设置失败
    goto restore
)
echo ✅ 环境变量设置完成

REM 步骤4：构建Go程序
echo 正在构建Go程序...
go build -tags=k8s -o webook .
if !errorlevel! neq 0 (
    echo ❌ Go程序构建失败
    goto restore
)
echo ✅ Go程序构建成功

REM 步骤5：删除旧Docker镜像
echo 正在删除旧Docker镜像...
docker rmi -f guanjian104/webook:v0.0.1
if !errorlevel! neq 0 (
    echo ℹ️ 未找到旧镜像或删除失败（可能是首次构建）
)
echo ✅ 旧镜像清理完成

REM 步骤6：构建新Docker镜像
echo 正在构建Docker镜像...
docker build -t guanjian104/webook:v0.0.1 .
if !errorlevel! neq 0 (
    echo ❌ Docker镜像构建失败
    goto restore
)
echo ✅ Docker镜像构建成功

REM 步骤7：恢复环境变量
:restore
echo 正在恢复 GOOS=windows...
go env -w GOOS=windows
if !errorlevel! neq 0 (
    echo ⚠️ 环境变量恢复失败
) else (
    echo ✅ 环境变量恢复完成
)

echo ============================================
echo 自动化构建流程完成！
echo ============================================
timeout /t 3 >nul
endlocal