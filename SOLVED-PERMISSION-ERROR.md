# ĐÃ SỬA: Lỗi "permission denied for table users"

## Nguyên nhân thực sự

Backend đang kết nối vào **PostgreSQL local** (đã cài trên máy) thay vì **Docker container**!

Cả hai đều chạy trên port 5432:
- PostgreSQL local: `C:\Program Files\PostgreSQL\18\bin\postgres.exe` (PID 4568)
- Docker PostgreSQL: Container `linkbio_postgres` (PID 8092)

Backend kết nối vào PostgreSQL local, nơi:
- User `linkbio` không tồn tại hoặc không có quyền
- Database `linkbio` có thể không tồn tại hoặc có cấu trúc khác
- Dẫn đến lỗi "permission denied"

## Cách đã sửa

Dừng PostgreSQL local service:
```powershell
Stop-Service -Name postgresql-x64-18 -Force
```

Sau đó restart backend để kết nối vào Docker container.

## Cách tránh lỗi này

### Option 1: Tắt PostgreSQL local khỏi startup (Khuyến nghị)
```powershell
Set-Service -Name postgresql-x64-18 -StartupType Disabled
```

### Option 2: Thay đổi port của Docker container
Sửa `docker-compose.yml`:
```yaml
services:
  postgres:
    ports:
      - "5433:5432"  # Thay vì 5432:5432
```

Và sửa `backend/.env`:
```
DATABASE_URL=postgres://linkbio:linkbio_dev_password@localhost:5433/linkbio?sslmode=disable
```

### Option 3: Kiểm tra trước khi chạy
```powershell
# Kiểm tra process nào đang dùng port 5432
netstat -ano | Select-String ":5432"
Get-Process -Id <PID> | Select-Object ProcessName, Path
```

## Kết quả

✅ Đăng ký tài khoản đã hoạt động!
✅ Console logs chi tiết giúp debug dễ dàng
✅ Backend kết nối đúng Docker PostgreSQL

## Test lại

Mở http://localhost:5173/auth/register và đăng ký tài khoản mới.

Logs sẽ hiển thị:
```
[AuthHandler] Received registration request
[AuthService] Starting registration
[UserRepository] Creating user
[UserRepository] User inserted successfully
[AuthHandler] Registration successful
```
