# 🐛 BUG FIX: Custom Theme bị ghi đè khi save Preset Theme

## 📋 Mô tả lỗi

Khi user thực hiện các bước sau:
1. Chọn theme Yoga → Chỉnh sửa → Save All → Chuyển sang Custom theme
2. Reload page → Chọn theme Dark → Save All  
3. Chọn lại Custom theme → **Custom theme bị mất, trở về default**

## 🔍 Nguyên nhân chính xác

### Vấn đề trong code cũ (dòng 233):

```javascript
// ❌ CODE CŨ - SAI
else {
    // Theme not modified → save as preset with header override
    savePromises.push(
        profileApi.updateProfile({
            theme_name: currentTheme,
            theme_config: null,
            header_config: JSON.stringify(updatedHeader),
            custom_theme_config: JSON.stringify({ customHeaderPresets: customHeaderPresetsSnapshot })
            // ☝️ LỖI Ở ĐÂY: Ghi đè toàn bộ custom_theme_config chỉ với customHeaderPresets
        }, token)
    );
}
```

### Luồng lỗi chi tiết:

**Bước 1: Tạo Custom Theme**
```
User: Chọn Yoga → Chỉnh sửa màu → Save All
DB: theme_name = 'custom'
    custom_theme_config = {
        pageBackground: '#red',
        cardBackground: '#blue',
        header: {...},
        customHeaderPresets: [...]
    }
✅ OK
```

**Bước 2: Chuyển sang Preset Theme**
```
User: Chọn Dark → Save All
Code cũ thực hiện:
    custom_theme_config = { customHeaderPresets: [...] }
    
DB: theme_name = 'dark'
    custom_theme_config = { customHeaderPresets: [...] }
    ❌ MẤT HẾT: pageBackground, cardBackground, header, etc.
```

**Bước 3: Quay lại Custom Theme**
```
User: Chọn Custom
Code load từ custom_theme_config:
    → Chỉ có customHeaderPresets
    → Không có theme config
    → Fallback về default theme
❌ Custom theme bị mất!
```

## ✅ Giải pháp

### Nguyên tắc:
**Khi save Preset theme, KHÔNG ĐƯỢC ghi đè `custom_theme_config`**

Thay vào đó:
1. Load `custom_theme_config` hiện tại từ database
2. MERGE `customHeaderPresets` vào config hiện có
3. Giữ nguyên toàn bộ custom theme config

### Code mới (đã fix):

```javascript
// ✅ CODE MỚI - ĐÚNG
else {
    // Get current custom_theme_config from database
    const currentProfile = await profileApi.getMyProfile(token);
    let existingCustomConfig = {};
    
    if (currentProfile?.custom_theme_config) {
        try {
            existingCustomConfig = typeof currentProfile.custom_theme_config === 'string'
                ? JSON.parse(currentProfile.custom_theme_config)
                : currentProfile.custom_theme_config;
        } catch (e) {
            console.warn('Failed to parse existing custom_theme_config:', e);
        }
    }
    
    // MERGE: Keep existing custom theme, only update customHeaderPresets
    const mergedCustomConfig = {
        ...existingCustomConfig,  // ← Giữ nguyên custom theme config
        customHeaderPresets: customHeaderPresetsSnapshot  // ← Chỉ update presets
    };
    
    savePromises.push(
        profileApi.updateProfile({
            theme_name: currentTheme,
            theme_config: null,
            header_config: JSON.stringify(updatedHeader),
            custom_theme_config: JSON.stringify(mergedCustomConfig)  // ← MERGE thay vì ghi đè
        }, token)
    );
}
```

## 🧪 Test Cases

### Test Case 1: Tạo Custom Theme rồi chuyển sang Preset
```
1. Chọn Yoga → Đổi màu đỏ → Save All → Custom theme
2. Reload page
3. Chọn Dark → Save All
4. Chọn Custom
✅ Expected: Custom theme vẫn giữ màu đỏ
✅ Actual: Custom theme giữ màu đỏ (FIXED)
```

### Test Case 2: Xóa Custom Header Preset trên Preset Theme
```
1. Đang ở theme Dark
2. Xóa một custom header preset → Save
3. Chọn Custom theme
✅ Expected: Custom theme không bị ảnh hưởng
✅ Actual: Custom theme không bị ảnh hưởng (FIXED)
```

### Test Case 3: Chỉnh sửa Header trên Preset Theme
```
1. Đang ở theme Yoga
2. Chỉnh header style → Save All
3. Chọn Custom theme
✅ Expected: Custom theme vẫn giữ nguyên
✅ Actual: Custom theme vẫn giữ nguyên (FIXED)
```

## 📝 Files đã sửa

1. `frontend/src/routes/dashboard/appearance/+page.svelte`
   - Function: `saveAllChanges()` (dòng ~230)
   - Function: `handleDeletePreset()` (dòng ~150)

## 🎯 Kết luận

**Root Cause**: Ghi đè `custom_theme_config` khi save preset theme

**Solution**: MERGE thay vì ghi đè, giữ nguyên custom theme config

**Impact**: 
- ✅ Custom theme không bị mất khi chuyển qua lại giữa preset và custom
- ✅ `customHeaderPresets` vẫn được lưu và sync đúng
- ✅ Không cần migration database
- ✅ Backward compatible

**Status**: ✅ FIXED
