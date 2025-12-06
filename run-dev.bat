@echo off
echo Starting Link in Bio Development Environment...
echo.

REM Check if Docker is running
docker ps >nul 2>&1
if errorlevel 1 (
    echo Error: Docker is not running. Please start Docker Desktop.
    pause
    exit /b 1
)

REM Start PostgreSQL
echo Starting PostgreSQL...
docker-compose up -d postgres
timeout /t 3 /nobreak >nul

REM Start Backend
echo Starting Backend Server...
start cmd /k "cd backend && go run main.go"

REM Wait a bit for backend to start
timeout /t 3 /nobreak >nul

REM Start Frontend
echo Starting Frontend Server...
start cmd /k "cd frontend && npm run dev"

echo.
echo ========================================
echo Development servers are starting...
echo ========================================
echo Frontend: http://localhost:5173
echo Backend:  http://localhost:3000
echo ========================================
echo.
echo Press any key to stop all servers...
pause >nul

REM Stop servers
taskkill /F /FI "WINDOWTITLE eq *go run main.go*" >nul 2>&1
taskkill /F /FI "WINDOWTITLE eq *npm run dev*" >nul 2>&1
docker-compose down
