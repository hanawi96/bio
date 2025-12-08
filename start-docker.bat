@echo off
echo Starting Docker Desktop...
start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"

echo Waiting for Docker to start...
timeout /t 10 /nobreak > nul

:check
docker info > nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not ready yet, waiting...
    timeout /t 5 /nobreak > nul
    goto check
)

echo Docker is ready!
echo Starting database containers...
docker-compose up -d

echo.
echo Database is starting...
echo PostgreSQL: localhost:5432
echo Redis: localhost:6379
echo.
echo Use 'docker-compose logs -f' to view logs
pause
