@echo off
setlocal enabledelayedexpansion

REM ===== Configurable Variables =====
set "IMAGE_NAME=guanjian104/webook"
set IMAGE_VERSION=%1
if "%IMAGE_VERSION%"=="" set IMAGE_VERSION=v0.0.1
set "BUILD_TAGS=k8s"
set "PROJECT_ROOT=%cd%"
REM =================================

echo ============================================
echo Starting automated build process
echo Project path: %PROJECT_ROOT%
echo Docker image: %IMAGE_NAME%:%IMAGE_VERSION%
echo ============================================
echo.

REM Step 1: Delete old executable
echo [1/7] Deleting old binary...
if exist webook (
    del /f /q webook >nul 2>&1
    if exist webook (
        echo ERROR: Failed to delete webook file (file may be in use)
    ) else (
        echo SUCCESS: Old webook file deleted
    )
) else (
    echo INFO: webook file not found, skipping deletion
)
echo.

REM Step 2: Execute go mod tidy
echo [2/7] Running go mod tidy...
go mod tidy
if %errorlevel% neq 0 (
    echo ERROR: go mod tidy failed (error code: %errorlevel%)
    goto restore
)
echo SUCCESS: go mod tidy completed
echo.

REM Step 3: Set GOOS to linux
echo [3/7] Setting GOOS=linux...
go env -w GOOS=linux
if %errorlevel% neq 0 (
    echo ERROR: Failed to set GOOS=linux (error code: %errorlevel%)
    goto restore
)
echo SUCCESS: GOOS set to linux
echo.

REM Step 4: Build Go application
echo [4/7] Building Go application (tags: %BUILD_TAGS%)...
go build -tags=%BUILD_TAGS% -o webook .
if %errorlevel% neq 0 (
    echo ERROR: Go build failed (error code: %errorlevel%)
    goto restore
)

if exist webook (
    echo SUCCESS: Go application built
) else (
    echo ERROR: webook file not found after build
    goto restore
)
echo.

REM Check Docker service status
echo [Check] Verifying Docker service...
docker version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Docker service not running or inaccessible
    goto restore
)
echo SUCCESS: Docker service is ready
echo.

REM Step 5: Remove old Docker image
echo [5/7] Removing old Docker image...
docker rmi -f "%IMAGE_NAME%:%IMAGE_VERSION%" 2>nul
if %errorlevel% neq 0 (
    echo INFO: Old image not found or removal failed (may be first build)
)
echo SUCCESS: Old image cleanup completed
echo.

REM Step 6: Build new Docker image
echo [6/7] Building Docker image...
docker build -t "%IMAGE_NAME%:%IMAGE_VERSION%" .
if %errorlevel% neq 0 (
    echo ERROR: Docker build failed (error code: %errorlevel%)
    goto restore
)

REM Verify image exists
docker images --format "{{.Repository}}:{{.Tag}}" | findstr /C:"%IMAGE_NAME%:%IMAGE_VERSION%" >nul
if %errorlevel% neq 0 (
    echo ERROR: Image not found after build
    goto restore
)
echo SUCCESS: Docker image built
echo.

:restore
REM Step 7: Restore environment variables
echo [7/7] Restoring environment variables...
go env -w GOOS=windows
if %errorlevel% neq 0 (
    echo WARNING: Failed to restore GOOS=windows (error code: %errorlevel%)
) else (
    echo SUCCESS: GOOS restored to windows
)
echo.

:end
echo ============================================
echo Build summary:
echo Project path: %PROJECT_ROOT%
echo Start time: %time%
echo.

if exist webook (
    echo SUCCESS: Go build completed
) else (
    echo FAILED: Go build failed
)

docker images --format "{{.Repository}}:{{.Tag}}" | findstr /C:"%IMAGE_NAME%:%IMAGE_VERSION%" >nul
if %errorlevel% equ 0 (
    echo SUCCESS: Docker image created
) else (
    echo FAILED: Docker image creation failed
)

echo ============================================
echo Press any key to exit...
pause >nul
endlocal