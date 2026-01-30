@echo off
echo ========================================
echo Super Admin Dashboard - Restart Script
echo ========================================
echo.
echo Step 1: Checking for process on port 3004...
for /f "tokens=5" %%A in ('netstat -ano ^| findstr :3004') do (
    echo Found process ID: %%A
    echo Killing process...
    taskkill /PID %%A /F
)
echo.
echo Step 2: Clearing cache files...
if exist ".next" rmdir /s /q .next
echo Cleared .next directory
echo.
echo Step 3: Starting dev server...
echo.
call pnpm dev
pause
