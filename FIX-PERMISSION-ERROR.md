# Hướng dẫn sửa lỗi "permission denied for table users"

## Nguyên nhân
Lỗi này xảy ra vì user `linkbio` trong PostgreSQL không có quyền INSERT/UPDATE/DELETE vào các bảng. Khi migration tạo bảng bằng user `postgres` (superuser), user `linkbio` không tự động có quyền truy cập.

## Cách sửa nhanh (Khuyến nghị)

### Bước 1: Chạy script fix permissions
```bash
fix-db-permissions.bat
```

Script này sẽ:
- Kết nối vào database với user `postgres` (superuser)
- Cấp tất cả quyền cho user `linkbio` trên các bảng hiện tại
- Cấp quyền mặc định cho các bảng tương lai

### Bước 2: Khởi động lại backend
```bash
# Dừng backend nếu đang chạy (Ctrl+C)
# Sau đó chạy lại
run-dev.bat
```

### Bước 3: Thử đăng ký lại
Mở http://localhost:5173/auth/register và thử đăng ký tài khoản mới.

## Cách sửa đầy đủ (Nếu muốn reset database)

### Bước 1: Dừng tất cả services
```bash
docker-compose down -v
```

### Bước 2: Khởi động lại database
```bash
docker-compose up -d postgres
```

### Bước 3: Chạy migrations (đã bao gồm fix permissions)
```bash
run-migrations.bat
```

### Bước 4: Khởi động lại toàn bộ
```bash
run-dev.bat
```

## Kiểm tra logs

Backend đã được thêm console.log chi tiết:
- `[AuthHandler]` - Log từ API handler
- `[AuthService]` - Log từ service layer
- `[UserRepository]` - Log từ database layer

Xem terminal backend để theo dõi quá trình đăng ký và phát hiện lỗi chính xác.

## Các lỗi khác có thể gặp

### Lỗi: "duplicate key value violates unique constraint"
- Email hoặc username đã tồn tại
- Thử với email/username khác

### Lỗi: "connection refused"
- Database chưa chạy
- Chạy: `docker-compose up -d postgres`

### Lỗi: "invalid password"
- Kiểm tra DATABASE_URL trong backend/.env
- Đảm bảo password khớp với docker-compose.yml
