@echo off
echo Building for Windows...
cd ..\main
go build -o ..\cutwise\src-tauri\resources\main\main.exe -ldflags=-H=windowsgui
