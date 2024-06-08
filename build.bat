@REM 2024年 author：Mr.Fang
@echo off
chcp 65001

setlocal enabledelayedexpansion

rem 定义目标平台
set platforms=linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64

rem 定义输出文件名
set output=wxdown

rem 定义输出版本
set version=1.0.2

rem 遍历每个平台并构建
for %%p in (%platforms%) do (
    for /f "tokens=1,2 delims=/" %%i in ("%%p") do (
        set GOOS=%%i
        set GOARCH=%%j
        if "%%i" == "windows" (
            echo Build failed for %output%-%version%-%%i-%%j
            go build -ldflags "-X main.runMode=binary -X main.version=%version%" -o %output%-%version%-%%i-%%j/%output%.exe  main.go
        ) else (
            echo Build failed for %output%-%version%-%%i-%%j
            go build -ldflags "-X main.runMode=binary -X main.version=%version%" -o %output%-%version%-%%i-%%j/%output%  main.go
        )
    )
)

echo Build completed!
endlocal
exit /b 0
