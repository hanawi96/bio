@echo off
echo ========================================
echo PostgreSQL Management Tool
echo ========================================
echo.
echo 1. Check PostgreSQL status
echo 2. Stop local PostgreSQL (if running)
echo 3. Disable local PostgreSQL startup
echo 4. Enable local PostgreSQL startup
echo 5. Start Docker PostgreSQL
echo 6. Stop Docker PostgreSQL
echo 7. Check port 5432 usage
echo 8. Exit
echo.
set /p choice="Choose an option (1-8): "

if "%choice%"=="1" goto check_status
if "%choice%"=="2" goto stop_local
if "%choice%"=="3" goto disable_local
if "%choice%"=="4" goto enable_local
if "%choice%"=="5" goto start_docker
if "%choice%"=="6" goto stop_docker
if "%choice%"=="7" goto check_port
if "%choice%"=="8" goto end
goto invalid

:check_status
echo.
echo Checking PostgreSQL services...
powershell -Command "Get-Service -Name postgresql* | Select-Object Name, Status, StartType | Format-Table -AutoSize"
echo.
echo Checking Docker containers...
docker ps --filter "name=linkbio_postgres" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
goto end

:stop_local
echo.
echo Stopping local PostgreSQL...
powershell -Command "Stop-Service -Name postgresql-x64-18 -Force"
echo Local PostgreSQL stopped!
goto end

:disable_local
echo.
echo Disabling local PostgreSQL from startup...
powershell -Command "Set-Service -Name postgresql-x64-18 -StartupType Disabled"
powershell -Command "Get-Service -Name postgresql-x64-18 | Select-Object Name, Status, StartType"
echo.
echo Local PostgreSQL will NOT start automatically on boot!
goto end

:enable_local
echo.
echo Enabling local PostgreSQL startup...
powershell -Command "Set-Service -Name postgresql-x64-18 -StartupType Automatic"
powershell -Command "Get-Service -Name postgresql-x64-18 | Select-Object Name, Status, StartType"
echo.
echo Local PostgreSQL will start automatically on boot!
goto end

:start_docker
echo.
echo Starting Docker PostgreSQL...
docker-compose up -d postgres
echo.
echo Waiting for PostgreSQL to be ready...
timeout /t 3 /nobreak >nul
docker exec linkbio_postgres pg_isready -U linkbio
goto end

:stop_docker
echo.
echo Stopping Docker PostgreSQL...
docker-compose stop postgres
goto end

:check_port
echo.
echo Checking what's using port 5432...
powershell -Command "netstat -ano | Select-String ':5432'"
echo.
echo Process details:
powershell -Command "$ports = netstat -ano | Select-String ':5432' | ForEach-Object { ($_ -split '\s+')[-1] } | Select-Object -Unique; foreach ($pid in $ports) { if ($pid -match '^\d+$') { Get-Process -Id $pid -ErrorAction SilentlyContinue | Select-Object Id, ProcessName, Path } }"
goto end

:invalid
echo.
echo Invalid choice! Please choose 1-8.
goto end

:end
echo.
pause
