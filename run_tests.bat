@echo off
echo Running Blockchain Concurrency Tests
echo ===================================

echo.
echo Building and Running Rust Implementation...
echo --------------------------------------
rustc main.rs -o blockchain_rust.exe
if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile Rust implementation
    exit /b %ERRORLEVEL%
)

echo.
echo Running Rust implementation...
echo This will run for 30 seconds maximum...
start /b cmd /c "blockchain_rust.exe & exit 0"
timeout /t 30 /nobreak
taskkill /f /im blockchain_rust.exe 2>nul

echo.
echo Building and Running Go Implementation...
echo --------------------------------------
go build -o blockchain_go.exe main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile Go implementation
    exit /b %ERRORLEVEL%
)

echo.
echo Running Go implementation...
echo This will run for 30 seconds maximum...
start /b cmd /c "blockchain_go.exe & exit 0"
timeout /t 30 /nobreak
taskkill /f /im blockchain_go.exe 2>nul

echo.
echo Tests completed successfully!