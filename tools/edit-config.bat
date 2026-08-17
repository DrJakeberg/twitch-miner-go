@echo off
setlocal

REM Open the Config Editor
REM Usage: tools\edit-config.bat [--config DIR] [--port PORT] [--tui] [--no-browser]

set "SCRIPT_DIR=%~dp0"
set "PROJECT_DIR=%SCRIPT_DIR%\.."
for %%I in ("%PROJECT_DIR%") do set "PROJECT_DIR=%%~fI"
set "BINARY=%PROJECT_DIR%\config-editor.exe"

echo Building config-editor...
cd /d "%PROJECT_DIR%"
go build -o config-editor.exe ./cmd/config-editor
if %errorlevel% neq 0 (
    echo Build failed.
    pause
    exit /b 1
)

"%BINARY%" %*

echo.
echo Config editor exited with code %errorlevel%.
pause
