@echo off
echo Fixing database permissions...
echo.

REM Check if PostgreSQL CLI is available
where psql >nul 2>&1
if errorlevel 1 (
    echo Using docker exec...
    docker exec -i linkbio_postgres psql -U postgres -d linkbio -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO linkbio;"
    docker exec -i linkbio_postgres psql -U postgres -d linkbio -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO linkbio;"
    docker exec -i linkbio_postgres psql -U postgres -d linkbio -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO linkbio;"
    docker exec -i linkbio_postgres psql -U postgres -d linkbio -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO linkbio;"
) else (
    echo Using local psql...
    psql postgres://postgres:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO linkbio;"
    psql postgres://postgres:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO linkbio;"
    psql postgres://postgres:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO linkbio;"
    psql postgres://postgres:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO linkbio;"
)

echo.
echo Permissions fixed! Try registering again.
pause
