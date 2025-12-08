@echo off
title LinkBio - Quick Start
color 0A

echo ========================================
echo   LINKBIO - QUICK START
echo ========================================
echo.

echo [1/4] Starting Docker Desktop...
start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"
timeout /t 10 /nobreak > nul

:check_docker
docker info > nul 2>&1
if %errorlevel% neq 0 (
    echo Waiting for Docker...
    timeout /t 3 /nobreak > nul
    goto check_docker
)

echo [2/4] Starting Database...
docker-compose up -d
timeout /t 5 /nobreak > nul

echo [3/4] Starting Backend...
start "Backend Server" cmd /k "cd backend && go run main.go"
timeout /t 3 /nobreak > nul

echo [4/4] Starting Frontend...
start "Frontend Dev Server" cmd /k "cd frontend && npm run dev"

echo.
echo ========================================
echo   ALL SERVICES STARTED!
echo ========================================
echo.
echo Backend:  http://localhost:8080
echo Frontend: http://localhost:5173
echo Database: localhost:5432
echo.
echo Press any key to open browser...
pause > nul

start http://localhost:5173/dashboard/links

echo.
echo Services are running in separate windows.
echo Close those windows to stop the services.
pause
