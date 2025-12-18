# Debug Guide: Theme không hoạt động sau khi đăng ký

## Vấn đề
Khi user mới đăng ký, chọn theme Yoga trong onboarding, sau đó vào trang Appearance chỉnh căn lề text → **KHÔNG hoạt động**.

## Nguyên nhân có thể
1. Theme không được lưu vào database khi onboarding
2. Theme được lưu nhưng không đúng format
3. Frontend load theme không đúng
4. Update theme không hoạt động vì thiếu data

## Các điểm đã thêm Console Log

### 1. Frontend - Onboarding (choose-theme)
File: `frontend/src/routes/onboarding/choose-theme/+page.svelte`

Khi click "Continue to Dashboard":
```
[ONBOARDING] Applying theme: yoga
[ONBOARDING] Theme data: {...}
[ONBOARDING] Apply theme response: {...}
```

### 2. Frontend - Appearance Page Load
File: `frontend/src/routes/dashboard/appearance/+page.svelte`

Khi vào trang Appearance:
```
[APPEARANCE] Profile data loaded: {...}
[APPEARANCE] theme_name: yoga
[APPEARANCE] theme_config: {...}
[APPEARANCE] custom_theme_config: null
[APPEARANCE] header_config: {...}
[APPEARANCE] Loaded theme: yoga category: cozy
[APPEARANCE] Current globalTheme: {...}
```

### 3. Frontend - Theme Loader
File: `frontend/src/lib/utils/themeLoader.ts`

Khi load theme:
```
[THEME_LOADER] Loading theme from profile: {...}
[THEME_LOADER] Theme name from profile: yoga
[THEME_LOADER] Loading preset theme: yoga
[THEME_LOADER] Theme loaded, current state: {...}
```

### 4. Frontend - Save Changes
File: `frontend/src/routes/dashboard/appearance/+page.svelte`

Khi click "Save All":
```
[APPEARANCE] Saving changes...
[APPEARANCE] Current theme: yoga
[APPEARANCE] Updated theme config: {...}
[APPEARANCE] Updated header config: {...}
[APPEARANCE] Saving as preset theme with header override
[APPEARANCE] Update payload: {...}
```

### 5. Backend - Apply Theme
File: `backend/service/profile_service.go`

Khi apply theme từ onboarding:
```
[APPLY_THEME] User: xxx, Theme: yoga
[APPLY_THEME] Theme config: {...}
[APPLY_THEME] Header config: {...}
[APPLY_THEME] Update data: {...}
[APPLY_THEME] Profile updated successfully
```

### 6. Backend - Profile Repository
File: `backend/repository/profile_repository.go`

Khi update profile:
```
[PROFILE_REPO] Updating profile for user: xxx
[PROFILE_REPO] Update data: {...}
[PROFILE_REPO] Theme config to save: {...}
[PROFILE_REPO] Header config to save: {...}
[PROFILE_REPO] Profile updated successfully
```

## Cách Debug

### Bước 1: Test Onboarding Flow
1. Đăng ký tài khoản mới
2. Chọn theme Yoga
3. Click "Continue to Dashboard"
4. **Mở Browser Console** (F12)
5. Xem logs:
   - `[ONBOARDING]` - Theme có được gửi đúng không?
   - `[APPLY_THEME]` - Backend có nhận được request không?
   - `[PROFILE_REPO]` - Database có được update không?

**Kiểm tra:**
- ✅ Nếu thấy tất cả logs → Theme đã được lưu
- ❌ Nếu thiếu logs → Tìm điểm bị lỗi

### Bước 2: Test Appearance Page Load
1. Sau khi onboarding xong, vào trang Appearance
2. **Mở Browser Console** (F12)
3. Xem logs:
   - `[APPEARANCE] Profile data loaded` - Profile có data không?
   - `[APPEARANCE] theme_name` - Theme name là gì?
   - `[APPEARANCE] theme_config` - Theme config có data không?
   - `[THEME_LOADER]` - Theme có được load đúng không?

**Kiểm tra:**
- ✅ Nếu `theme_name: "yoga"` và `theme_config: {...}` → Theme đã được lưu đúng
- ❌ Nếu `theme_name: null` hoặc `theme_config: null` → Theme không được lưu trong onboarding
- ❌ Nếu `theme_name: "yoga"` nhưng `theme_config: null` → ApplyTheme không lưu theme_config

### Bước 3: Test Update Theme
1. Ở trang Appearance, chỉnh căn lề text sang phải
2. Click "Save All"
3. **Mở Browser Console** (F12)
4. Xem logs:
   - `[APPEARANCE] Saving changes...`
   - `[APPEARANCE] Current theme` - Theme hiện tại là gì?
   - `[APPEARANCE] Updated theme config` - Config có thay đổi không?
   - `[PROFILE_REPO]` - Database có được update không?

**Kiểm tra:**
- ✅ Nếu thấy logs và không có error → Update thành công
- ❌ Nếu có error → Xem error message

### Bước 4: Kiểm tra Database
Nếu vẫn không rõ, check database trực tiếp:

```sql
-- Lấy profile của user vừa đăng ký
SELECT 
  theme_name,
  theme_config,
  custom_theme_config,
  header_config
FROM profiles
WHERE user_id = 'YOUR_USER_ID';
```

**Kết quả mong đợi sau onboarding:**
```
theme_name: "yoga"
theme_config: {"pageBackground": "#...", ...}  -- Phải có data
header_config: {"layout": "...", ...}          -- Phải có data
custom_theme_config: null                       -- Có thể null
```

**Nếu theme_config là NULL:**
→ Vấn đề ở ApplyTheme API, không lưu theme_config

**Nếu theme_name là NULL:**
→ Vấn đề ở onboarding flow, không gọi ApplyTheme API

## Các Scenario và Giải pháp

### Scenario 1: Theme không được lưu trong onboarding
**Triệu chứng:**
- Không thấy logs `[ONBOARDING]` hoặc `[APPLY_THEME]`
- Database: `theme_name` và `theme_config` đều NULL

**Nguyên nhân:**
- Onboarding không gọi `profileApi.applyTheme()`
- API call bị lỗi nhưng không hiển thị error

**Giải pháp:**
- Kiểm tra network tab xem có request `/profile/apply-theme` không
- Xem response có error không

### Scenario 2: Theme được lưu nhưng theme_config NULL
**Triệu chứng:**
- Thấy logs `[APPLY_THEME]` với theme config
- Database: `theme_name: "yoga"` nhưng `theme_config: NULL`

**Nguyên nhân:**
- Backend nhận được data nhưng không lưu vào database
- SQL query có vấn đề với COALESCE

**Giải pháp:**
- Xem logs `[PROFILE_REPO]` để xác nhận data được gửi đến repository
- Check SQL query trong `profile_repository.go`

### Scenario 3: Theme được lưu đúng nhưng không load được
**Triệu chứng:**
- Database có đầy đủ data
- Logs `[APPEARANCE]` show `theme_config: null` hoặc empty

**Nguyên nhân:**
- Frontend không parse JSON đúng
- `loadThemeFromProfile()` có bug

**Giải pháp:**
- Xem logs `[THEME_LOADER]` để xác nhận theme được load
- Check `globalTheme.setPreset()` có hoạt động không

### Scenario 4: Theme load được nhưng update không hoạt động
**Triệu chứng:**
- Theme hiển thị đúng ban đầu
- Chỉnh sửa không có effect

**Nguyên nhân:**
- Preview không sync với theme store
- Update API không được gọi

**Giải pháp:**
- Xem logs `[APPEARANCE] Saving changes...`
- Check `syncPreviewStylesFromTheme()` có được gọi không

## Kết luận

Sau khi chạy test với console logs, bạn sẽ biết chính xác:
1. Theme có được lưu trong onboarding không?
2. Theme có được load đúng trong Appearance không?
3. Update theme có hoạt động không?

Dựa vào logs, bạn sẽ tìm được nguyên nhân chính xác và fix đúng chỗ.
