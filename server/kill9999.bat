@echo off
for /f "tokens=2,5" %%i in ('netstat -ano') do (
 	if %%i == 0.0.0.0:9999 (if not %%j == 0 taskkill -PID %%j -F)
)
pause