@echo off

:: Find and kill the Go backend process
taskkill /F /IM main.exe /T

:: /F: Forcefully terminate the process
:: /IM: Specify the image name (backend.exe in this case)
:: /T: Terminate the specified process and any child processes