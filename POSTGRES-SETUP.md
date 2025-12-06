# Hướng dẫn cấu hình PostgreSQL

## Vấn đề

Bạn có PostgreSQL cài local trên máy Windows và PostgreSQL trong Docker. Cả hai đều cố gắng chạy trên port 5432, gây xung đột.

## Giải pháp đã thực hiện ✅

PostgreSQL local đã được **DISABLED** khỏi startup. Nó sẽ không tự động chạy khi khởi động máy.

```powershell
# Kiểm tra trạng thái
Get-Service -Name postgresql-x64-18

# Kết quả:
# Name               Status StartType
# ----               ------ ---------
# postgresql-x64-18 Stopped  Disabled
```

## Sử dụng hàng ngày

### Khởi động project (chỉ dùng Docker)

```bash
# 1. Khởi động Docker PostgreSQL
docker-compose up -d postgres

# 2. Chạy migrations (nếu cần)
run-migrations.bat

# 3. Khởi động backend
cd backend
go run main.go

# 4. Khởi động frontend (terminal khác)
cd frontend
npm run dev
```

### Hoặc dùng script all-in-one

```bash
run-dev.bat
```

## Tool quản lý PostgreSQL

Chạy script tiện ích:
```bash
manage-postgres.bat
```

Menu options:
1. **Check status** - Xem trạng thái cả local và Docker PostgreSQL
2. **Stop local** - Dừng PostgreSQL local nếu đang chạy
3. **Disable local** - Tắt auto-start PostgreSQL local
4. **Enable local** - Bật lại auto-start (nếu cần dùng local)
5. **Start Docker** - Khởi động Docker PostgreSQL
6. **Stop Docker** - Dừng Docker PostgreSQL
7. **Check port** - Xem process nào đang dùng port 5432
8. **Exit**

## Khi nào cần PostgreSQL local?

Nếu bạn có project khác cần PostgreSQL local:

### Option 1: Bật tạm thời (không auto-start)
```powershell
Start-Service -Name postgresql-x64-18
# Dùng xong thì dừng
Stop-Service -Name postgresql-x64-18
```

### Option 2: Đổi port cho Docker (khuyến nghị)
Sửa `docker-compose.yml`:
```yaml
services:
  postgres:
    ports:
      - "5433:5432"  # Docker dùng port 5433
```

Sửa `backend/.env`:
```env
DATABASE_URL=postgres://linkbio:linkbio_dev_password@localhost:5433/linkbio?sslmode=disable
```

Sau đó:
```bash
docker-compose down
docker-compose up -d postgres
```

## Troubleshooting

### Lỗi: "port 5432 already in use"

```bash
# Kiểm tra process nào đang dùng
manage-postgres.bat
# Chọn option 7

# Hoặc dùng PowerShell
netstat -ano | Select-String ":5432"
```

### Lỗi: "connection refused"

```bash
# Kiểm tra Docker PostgreSQL có chạy không
docker ps --filter "name=linkbio_postgres"

# Nếu không chạy, khởi động
docker-compose up -d postgres
```

### Lỗi: "permission denied for table users"

Đây là lỗi cũ - đã fix bằng cách dừng PostgreSQL local. Nếu vẫn gặp:

1. Kiểm tra backend đang kết nối đúng database:
   - Xem logs backend khi khởi động
   - Phải thấy: `[Database] Current user: linkbio Database: linkbio`

2. Kiểm tra không có PostgreSQL local đang chạy:
   ```bash
   manage-postgres.bat
   # Chọn option 1
   ```

## Best Practices

✅ **DO:**
- Dùng Docker PostgreSQL cho development
- Dùng `manage-postgres.bat` để quản lý
- Kiểm tra port 5432 trước khi chạy project

❌ **DON'T:**
- Bật cả local và Docker PostgreSQL cùng lúc
- Quên check service status khi gặp lỗi connection
- Thay đổi DATABASE_URL mà không restart backend

## Tóm tắt

- ✅ PostgreSQL local: **DISABLED** (không auto-start)
- ✅ Docker PostgreSQL: Dùng cho project này
- ✅ Port 5432: Chỉ Docker PostgreSQL sử dụng
- ✅ Tool: `manage-postgres.bat` để quản lý dễ dàng
