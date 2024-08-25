@echo off
cd .\resources\main
echo Running backend... > backend.log
start /B main.exe >> backend.log 2>&1
echo Done. Check backend.log for details.
pause