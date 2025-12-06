@echo off
echo Running database migrations...
echo.

REM Check if PostgreSQL CLI is available
where psql >nul 2>&1
if errorlevel 1 (
    echo psql not found in PATH. Using docker exec instead...
    echo Running 001_init.sql...
    docker exec -i linkbio_postgres psql -U linkbio -d linkbio < backend\db\migrations\001_init.sql
    echo Running 002_grant_permissions.sql...
    docker exec -i linkbio_postgres psql -U linkbio -d linkbio < backend\db\migrations\002_grant_permissions.sql
) else (
    echo Using local psql...
    echo Running 001_init.sql...
    psql postgres://linkbio:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable -f backend\db\migrations\001_init.sql
    echo Running 002_grant_permissions.sql...
    psql postgres://linkbio:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable -f backend\db\migrations\002_grant_permissions.sql
)

if errorlevel 1 (
    echo.
    echo Error: Migration failed!
    echo Make sure PostgreSQL is running: docker-compose up -d postgres
    pause
    exit /b 1
)

echo.
echo Migration completed successfully!
pause
