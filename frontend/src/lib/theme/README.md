# 🎨 Theme System

Hệ thống theme dựa trên **Design Tokens + Preset Rules + CSS Variables**.

## 📐 Kiến trúc

```
Theme = Design Tokens + Theme Engine + CSS Variables
         ↓                ↓              ↓
    JSON/Object      Conversion      Runtime CSS
```

### Luồng hoạt động:

1. **Design Tokens** (JSON) → Định nghĩa giá trị core (colors, spacing, typography)
2. **Theme Engine** (TS) → Convert tokens thành CSS variables
3. **CSS Variables** (CSS) → Apply vào components tại runtime
4. **Components** → Sử dụng CSS variables, không hardcode

## 🚀 Cách sử dụng

### 1. Wrap app với ThemeProvider

```svelte
<!-- +layout.svelte -->
<script>
  import { ThemeProvider } from '$lib/theme';
</script>

<ThemeProvider>
  <slot />
</ThemeProvider>
```

### 2. Sử dụng themed components

```svelte
<script>
  import { ThemedCard, ThemedLink } from '$lib/theme';
</script>

<ThemedCard variant="elevated">
  <h2>My Card</h2>
  <p>Content here</p>
</ThemedCard>

<ThemedLink 
  href="https://example.com"
  title="My Link"
  description="Click here"
/>
```

### 3. Sử dụng CSS variables trực tiếp

```svelte
<div style="
  background-color: var(--color-primary);
  padding: var(--spacing-md);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
">
  Custom styled element
</div>
```

### 4. Thay đổi theme động

```svelte
<script>
  import { useTheme } from '$lib/theme';
  
  const theme = useTheme();
  
  // Đổi preset
  theme.setPreset('dark');
  
  // Update một phần tokens
  theme.update({
    colors: {
      primary: { value: '#ff0000', opacity: 100 }
    }
  });
  
  // Export theme
  const json = theme.exportJSON();
  
  // Import theme
  theme.loadFromJSON(json);
</script>
```

## 🎯 Ưu điểm

### ✅ Không hardcode CSS
- Tất cả styles đều từ tokens
- Dễ maintain và scale

### ✅ Runtime theming
- Đổi theme không cần rebuild
- CSS variables update tức thì

### ✅ Type-safe
- Full TypeScript support
- Autocomplete cho tokens

### ✅ Export/Import
- Lưu theme dưới dạng JSON
- Có thể lưu vào DB hoặc file

### ✅ Preset system
- Có sẵn nhiều theme (default, dark, minimal, vibrant)
- Dễ tạo preset mới

### ✅ Component-level overrides
- Có thể override styles cho từng component
- Vẫn giữ consistency

## 📦 Structure

```
theme/
├── tokens.ts              # Design tokens definitions
├── engine.ts              # Theme engine (tokens → CSS)
├── context.svelte.ts      # Svelte 5 runes store
├── index.ts               # Main exports
├── components/
│   ├── ThemeProvider.svelte
│   ├── ThemedCard.svelte
│   └── ThemedLink.svelte
└── README.md
```

## 🎨 Available Tokens

### Colors
- `primary`, `secondary`, `accent`
- `background`, `surface`
- `text`, `textSecondary`
- `border`, `error`, `success`, `warning`

### Typography
- `xs`, `sm`, `base`, `lg`, `xl`, `2xl`, `3xl`

### Spacing
- `xs`, `sm`, `md`, `lg`, `xl`, `2xl`

### Radius
- `none`, `sm`, `md`, `lg`, `xl`, `full`

### Shadows
- `none`, `sm`, `md`, `lg`, `xl`

## 🔧 Tích hợp với Backend

### Lưu theme vào DB

```typescript
// Frontend
const themeJSON = theme.exportJSON();

// POST to backend
await fetch('/api/profile/theme', {
  method: 'POST',
  body: JSON.stringify({ theme: themeJSON })
});
```

### Load theme từ DB

```typescript
// GET from backend
const response = await fetch('/api/profile/theme');
const { theme: themeJSON } = await response.json();

// Apply theme
theme.loadFromJSON(themeJSON);
```

### Database Schema

```sql
ALTER TABLE profiles ADD COLUMN theme_config JSONB;
```

## 🎭 Tạo preset mới

```typescript
// tokens.ts
export const themePresets = {
  // ... existing presets
  
  myCustomTheme: {
    colors: {
      primary: { value: '#your-color', opacity: 100 },
      // ... other colors
    },
    radius: {
      none: 0,
      sm: 4,
      // ... other radius
    }
  }
};
```

## 🧪 Testing

Truy cập `/theme-demo` để test theme system:
- Đổi preset
- Customize colors, radius
- Export/Import JSON
- Xem preview components

## 📝 Best Practices

1. **Luôn dùng CSS variables** thay vì hardcode colors
2. **Sử dụng themed components** khi có thể
3. **Override cẩn thận** - chỉ override khi thực sự cần
4. **Test với nhiều themes** - đảm bảo UI work với tất cả presets
5. **Export theme** trước khi deploy - backup theme config

## 🔮 Roadmap

- [ ] Theme marketplace (users có thể share themes)
- [ ] Visual theme editor (drag & drop)
- [ ] Animation tokens
- [ ] Gradient support
- [ ] Dark mode auto-detect
- [ ] A11y contrast checker
