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
echo Running: go mod tidy completed
echo.

REM Step 3: Set GOOS to linux
echo [3/7] Setting GOOS=linux...
go env -w GOOS=linux
echo Running: GOOS set to linux
echo.

REM Step 4: Build Go application
echo [4/7] Building Go application (tags: %BUILD_TAGS%)...
go build -tags=%BUILD_TAGS% -o webook .

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
echo SUCCESS: Docker service is ready
echo.

REM Step 5: Remove old Docker image
echo [5/7] Removing old Docker image...
docker rmi -f "%IMAGE_NAME%:%IMAGE_VERSION%" 2>nul
echo Running: Old image cleanup completed
echo.

REM Step 6: Build new Docker image
echo [6/7] Building Docker image...
docker build -t "%IMAGE_NAME%:%IMAGE_VERSION%" .

REM Verify image exists
docker images --format "{{.Repository}}:{{.Tag}}" | findstr /C:"%IMAGE_NAME%:%IMAGE_VERSION%" >nul
echo Running: Docker image built
echo.

:restore
REM Step 7: Restore environment variables
echo [7/7] Restoring environment variables...
go env -w GOOS=windows
echo SUCCESS: GOOS restored to windows
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
echo SUCCESS: Docker image created

echo ============================================
echo Press any key to exit...
pause >nul
endlocal