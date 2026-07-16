@echo off
setlocal

REM Build and run twitch-miner-go (local development defaults)
REM Adds -no-lifecycle-notify to suppress start/stop/crash notifications.
REM Any additional flags are passed through.
REM Usage: _run-localdev.bat [flags]
REM Example: _run-localdev.bat -config configs -port 9090 -log-level debug

cd /d "%~dp0"
_run.bat -no-lifecycle-notify %*
