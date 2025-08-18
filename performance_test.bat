@echo off
echo Blockchain Performance Comparison
echo ===============================

echo.
echo Building Rust Implementation...
rustc -O main.rs -o blockchain_rust.exe
if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile Rust implementation
    exit /b %ERRORLEVEL%
)

echo.
echo Building Go Implementation...
go build -o blockchain_go.exe main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile Go implementation
    exit /b %ERRORLEVEL%
)

echo.
echo Running Rust Performance Test...
echo -------------------------------
echo This will run for 30 seconds maximum...
set start_time=%time%
start /b cmd /c "blockchain_rust.exe & exit 0"
timeout /t 30 /nobreak
taskkill /f /im blockchain_rust.exe 2>nul
set end_time=%time%

echo.
echo Start time: %start_time%
echo End time: %end_time%

echo.
echo Running Go Performance Test...
echo -----------------------------
echo This will run for 30 seconds maximum...
set start_time=%time%
start /b cmd /c "blockchain_go.exe & exit 0"
timeout /t 30 /nobreak
taskkill /f /im blockchain_go.exe 2>nul
set end_time=%time%

echo.
echo Start time: %start_time%
echo End time: %end_time%

echo.
echo Performance tests completed!
echo.
echo Note: For more accurate performance measurements, consider using a dedicated
echo benchmarking tool that can provide detailed metrics on CPU usage, memory
echo consumption, and execution time.