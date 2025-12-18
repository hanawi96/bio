# Quick Fix: Card Text Color Issue

## Vấn đề
Database có giá trị mặc định `#10b981` cho `card_text_color`, khiến frontend nghĩ là custom value.

## Giải pháp

### Option 1: Chạy SQL trực tiếp (NHANH NHẤT)
```sql
-- Chạy trong database của bạn
ALTER TABLE links ALTER COLUMN card_text_color DROP NOT NULL;
UPDATE links SET card_text_color = NULL WHERE is_group = true;
ALTER TABLE links ALTER COLUMN card_text_color SET DEFAULT NULL;
```

### Option 2: Restart backend để migration tự chạy
1. Stop backend server
2. Start lại backend
3. Migration `027_fix_card_text_color_nullable.sql` sẽ tự động chạy

### Option 3: Chạy migration thủ công
```bash
# Nếu có migrate tool
cd backend
go run cmd/migrate/main.go
```

## Sau khi fix
1. Refresh trang Bio
2. Card text color sẽ không còn chấm xanh
3. Màu sẽ inherit từ theme (trắng)
4. Khi user custom màu → chấm xanh xuất hiện

## Verify
Check console log:
```
card_text_color: null  ← Đúng!
card_text_color_custom: false  ← Đúng!
```
