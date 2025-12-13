# 🔌 Integration Guide

Hướng dẫn tích hợp theme system vào hệ thống LinkBio hiện tại.

## 1. Migrate ProfilePreview.svelte

### Before (Hardcoded):
```svelte
<a
  href={child.url}
  class="block transition-all"
  style="
    background-color: rgba(255, 255, 255, 1);
    border-radius: 12px;
    padding: 16px;
    box-shadow: 0px 4px 10px rgba(0,0,0,0.2);
  "
>
  {child.title}
</a>
```

### After (Theme-based):
```svelte
<script>
  import { useTheme } from '$lib/theme';
  const theme = useTheme();
</script>

<a
  href={child.url}
  class="block transition-all"
  style="
    background-color: var(--link-bg);
    border-radius: var(--link-radius);
    padding: var(--link-padding);
    box-shadow: var(--link-shadow);
    color: var(--link-text);
  "
>
  {child.title}
</a>
```

## 2. Migrate [username]/+page.svelte

### Before:
```svelte
<div class="bg-gradient-to-br from-purple-50 to-blue-50">
  <!-- content -->
</div>
```

### After:
```svelte
<script>
  import { ThemeProvider } from '$lib/theme';
</script>

<ThemeProvider>
  <div style="background-color: var(--color-surface);">
    <!-- content -->
  </div>
</ThemeProvider>
```

## 3. Load theme từ profile

### +page.ts (Load data):
```typescript
export async function load({ params }) {
  const data = await api.get(`/p/${params.username}`);
  
  return {
    profile: data.profile,
    links: data.links,
    blocks: data.blocks,
    themeConfig: data.profile.theme_config // Theme JSON từ DB
  };
}
```

### +page.svelte (Apply theme):
```svelte
<script lang="ts">
  import { ThemeProvider, useTheme } from '$lib/theme';
  import { onMount } from 'svelte';
  
  export let data;
  
  const theme = useTheme();
  
  onMount(() => {
    // Load theme từ profile
    if (data.themeConfig) {
      theme.loadFromJSON(data.themeConfig);
    }
  });
</script>

<ThemeProvider>
  <!-- Your content -->
</ThemeProvider>
```

## 4. Backend API endpoints

### GET /api/profile/theme
```go
// repository/profile_repository.go
func (r *ProfileRepository) GetThemeConfig(userID int) (string, error) {
    var themeConfig sql.NullString
    err := r.db.QueryRow(`
        SELECT theme_config FROM profiles WHERE user_id = $1
    `, userID).Scan(&themeConfig)
    
    if err != nil {
        return "", err
    }
    
    if !themeConfig.Valid {
        return "", nil // No theme config
    }
    
    return themeConfig.String, nil
}
```

### POST /api/profile/theme
```go
// api/profile_handler.go
func (h *ProfileHandler) UpdateTheme(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    var req struct {
        ThemeConfig string `json:"theme_config"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // Validate JSON
    var themeData map[string]interface{}
    if err := json.Unmarshal([]byte(req.ThemeConfig), &themeData); err != nil {
        c.JSON(400, gin.H{"error": "Invalid theme JSON"})
        return
    }
    
    // Save to DB
    err := h.profileRepo.UpdateThemeConfig(userID, req.ThemeConfig)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to update theme"})
        return
    }
    
    c.JSON(200, gin.H{"message": "Theme updated"})
}
```

## 5. Dashboard - Theme Editor

### Create theme editor component:
```svelte
<!-- frontend/src/lib/components/dashboard/ThemeEditor.svelte -->
<script lang="ts">
  import { useTheme, themePresets } from '$lib/theme';
  import { api } from '$lib/api/client';
  
  const theme = useTheme();
  let saving = $state(false);
  
  async function saveTheme() {
    saving = true;
    try {
      const themeJSON = theme.exportJSON();
      await api.post('/profile/theme', { theme_config: themeJSON });
      alert('Theme saved!');
    } catch (e) {
      alert('Failed to save theme');
    } finally {
      saving = false;
    }
  }
</script>

<div class="space-y-4">
  <h2 class="text-2xl font-bold">Theme Settings</h2>
  
  <!-- Preset selector -->
  <div>
    <label>Choose Preset</label>
    <select onchange={(e) => theme.setPreset(e.target.value)}>
      <option value="default">Default</option>
      <option value="dark">Dark</option>
      <option value="minimal">Minimal</option>
      <option value="vibrant">Vibrant</option>
    </select>
  </div>
  
  <!-- Color pickers -->
  <div>
    <label>Primary Color</label>
    <input 
      type="color" 
      value={theme.engine.getTokens().colors.primary.value}
      oninput={(e) => theme.update({
        colors: {
          ...theme.engine.getTokens().colors,
          primary: { value: e.target.value, opacity: 100 }
        }
      })}
    />
  </div>
  
  <!-- Save button -->
  <button onclick={saveTheme} disabled={saving}>
    {saving ? 'Saving...' : 'Save Theme'}
  </button>
</div>
```

## 6. Migrate existing inline styles

### Find all hardcoded styles:
```bash
# Search for hardcoded colors
grep -r "background-color: #" frontend/src/

# Search for hardcoded padding
grep -r "padding: [0-9]" frontend/src/

# Search for hardcoded border-radius
grep -r "border-radius: [0-9]" frontend/src/
```

### Replace với CSS variables:
```svelte
<!-- Before -->
<div style="background-color: #ffffff; padding: 16px; border-radius: 12px;">

<!-- After -->
<div style="background-color: var(--card-bg); padding: var(--card-padding); border-radius: var(--card-radius);">
```

## 7. Testing checklist

- [ ] Theme loads correctly from DB
- [ ] Theme changes apply instantly
- [ ] Export/Import works
- [ ] All presets work
- [ ] Public page uses user's theme
- [ ] Dashboard editor saves theme
- [ ] No hardcoded styles remain
- [ ] Mobile responsive
- [ ] Performance is good

## 8. Migration strategy

### Phase 1: Setup (1 day)
1. Add theme system files
2. Run migration (add theme_config column)
3. Create API endpoints
4. Test basic functionality

### Phase 2: Dashboard (2 days)
1. Create theme editor component
2. Add to dashboard settings
3. Test save/load
4. Add preset selector

### Phase 3: Public pages (2 days)
1. Wrap public pages with ThemeProvider
2. Load theme from profile
3. Replace hardcoded styles
4. Test all layouts (grid, carousel, card, list)

### Phase 4: Polish (1 day)
1. Add loading states
2. Error handling
3. Validation
4. Documentation

## 9. Performance tips

### ✅ DO:
- Use CSS variables (fast)
- Cache parsed styles
- Lazy load theme editor
- Debounce theme updates

### ❌ DON'T:
- Parse JSON on every render
- Create new objects in render
- Use inline styles everywhere
- Update theme too frequently

## 10. Rollback plan

Nếu có vấn đề, có thể rollback dễ dàng:

1. Theme system là optional - không break existing code
2. Database column có thể NULL - không ảnh hưởng existing users
3. Có thể disable ThemeProvider và dùng lại hardcoded styles
4. Theme JSON có thể export/backup trước khi deploy

## 11. Future enhancements

- [ ] Theme marketplace (users share themes)
- [ ] Visual theme builder (no code)
- [ ] AI theme generator
- [ ] Theme templates by industry
- [ ] A11y contrast checker
- [ ] Animation presets
- [ ] Font family support
- [ ] Background patterns/gradients
