# Theme Architecture - Custom Theme Management

## Overview
Hệ thống quản lý theme được refactor để tách biệt hoàn toàn giữa **Preset Themes** và **Custom Theme**, đảm bảo không bao giờ mất dữ liệu custom theme khi user chuyển đổi giữa các theme.

## Database Schema

### Profiles Table
```sql
- theme_name: VARCHAR(50)           -- Tên theme hiện tại ('yoga', 'dark', 'custom', etc.)
- theme_config: JSONB               -- Config của preset theme (deprecated, giữ để backward compatibility)
- custom_theme_config: JSONB       -- Config của custom theme (NEW)
```

## Logic Flow

### 1. Khi User Chỉnh Sửa Preset Theme
```
User chọn "Yoga" → Chỉnh card background → Auto-switch sang "Custom"
├─ UI: currentTheme = 'custom'
├─ Memory: globalTheme.update({ cardBackground: '#ff0000' })
└─ Chưa save vào DB
```

### 2. Khi Save All
```typescript
if (currentTheme === 'custom') {
    // Lưu vào custom_theme_config
    await profileApi.updateProfile({ 
        theme_name: 'custom',
        custom_theme_config: JSON.stringify(globalTheme.getCurrent())
    });
} else {
    // Lưu preset name, GIỮ NGUYÊN custom_theme_config
    await profileApi.updateProfile({ 
        theme_name: currentTheme
    });
}
```

### 3. Khi Load Theme (Page Load)
```typescript
if (theme_name === 'custom') {
    // Load từ custom_theme_config
    globalTheme.loadFromJSON(custom_theme_config);
} else {
    // Load preset từ themePresets
    globalTheme.setPreset(theme_name);
}
```

### 4. Khi Chọn Theme
```typescript
// Chọn Preset
selectTheme('yoga') → globalTheme.setPreset('yoga') → Mark pending changes

// Chọn Custom
selectTheme('custom') → Load từ custom_theme_config → Reset pending changes
```

## Key Benefits

✅ **Tách biệt rõ ràng**: Preset và Custom không xung đột  
✅ **Không mất data**: Custom config luôn được giữ trong `custom_theme_config`  
✅ **Logic đơn giản**: Không cần xóa/null field  
✅ **Backward compatible**: Giữ `theme_config` cho compatibility  

## Migration

```sql
-- Migration 024
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS custom_theme_config JSONB;
```

## Files Changed

### Backend
- `backend/migrations/024_add_custom_theme_config.sql` - Migration file
- `backend/repository/models.go` - Thêm field `CustomThemeConfig`
- `backend/repository/profile_repository.go` - Update queries để handle field mới

### Frontend
- `frontend/src/routes/dashboard/appearance/+page.svelte` - Update logic save/load theme

## Testing Scenarios

### Scenario 1: Tạo Custom Theme
1. Chọn theme "Yoga"
2. Thay đổi card background → màu đỏ
3. UI auto-switch sang "Custom"
4. Click "Save All"
5. ✅ DB: `theme_name='custom'`, `custom_theme_config={...màu đỏ...}`

### Scenario 2: Chuyển sang Preset
1. Tiếp tục từ Scenario 1
2. Chọn theme "Dark"
3. Click "Save All"
4. ✅ DB: `theme_name='dark'`, `custom_theme_config` VẪN GIỮ NGUYÊN

### Scenario 3: Quay lại Custom
1. Tiếp tục từ Scenario 2
2. Chọn theme "Custom"
3. ✅ UI hiển thị đúng màu đỏ (load từ `custom_theme_config`)

### Scenario 4: Reload Page
1. Refresh trang
2. ✅ Theme hiển thị đúng (Dark hoặc Custom tùy lần save cuối)
3. ✅ Custom config không bị mất

## Notes

- `theme_config` field được giữ lại để backward compatibility nhưng không còn được sử dụng
- `custom_theme_config` là single source of truth cho custom theme
- Khi user chưa bao giờ tạo custom theme, `custom_theme_config` sẽ là NULL
- Khi chọn Custom mà `custom_theme_config` NULL → fallback về default preset
