@echo off
echo Starting Blockchain Interface Server...
echo.
echo Please wait while the server starts...
echo.

REM Check if Node.js is installed
node --version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Error: Node.js is not installed or not in PATH.
    echo Please install Node.js from https://nodejs.org/
    echo.
    pause
    exit /b 1
)

echo Node.js is installed. Starting server...
echo.
echo The interface will be available at: http://localhost:3000/
echo.
echo Press Ctrl+C to stop the server when done.
echo.

node server.js

pause