@echo off
echo ========================================
echo Running Thumbnail Migration
echo ========================================
echo.

echo Connecting to PostgreSQL...
psql -U linkbio -d linkbio -f backend/db/migrations/003_add_thumbnail_url.sql

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo Migration completed successfully!
    echo ========================================
    echo.
    echo Verifying column...
    psql -U linkbio -d linkbio -c "\d links"
) else (
    echo.
    echo ========================================
    echo Migration failed!
    echo ========================================
    echo Please check your PostgreSQL connection
)

echo.
pause
