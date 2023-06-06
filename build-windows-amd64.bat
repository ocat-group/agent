@echo off
set TARGET_DIR=target

cd /d %~dp0

echo Building your Go project...

set GOOS=windows&set GOARCH=amd64&go build -o %TARGET_DIR%\agent-windows.exe

if %errorlevel%==0 (
  echo Build successful!
) else (
  echo Build failed!
)

pause