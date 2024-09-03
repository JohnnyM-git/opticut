@echo off
echo Building for Windows...

:: Change to the directory containing your Go code
cd ..\main

:: Build the Go executable
go build -o ../cutwise/src-tauri/resources/main/main.exe -ldflags=-H=windowsgui

:: Copy settings.json to the resources directory
copy ..\main\globals\settings.json ..\cutwise\src-tauri\resources\globals\settings.json

:: Copy setup_db.sql to the resources directory
copy ..\main\internal\db\setup_db.sql ..\cutwise\src-tauri\resources\main\internal\db\setup_db.sql

cd ..\cutwise

pnpm tauri build


echo Build and file copying completed.
