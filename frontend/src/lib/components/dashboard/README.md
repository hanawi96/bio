# Dashboard Components Structure

Cấu trúc components được tổ chức theo **feature-based architecture**, mirror với routes structure.

## 📁 Structure

```
dashboard/
├── bio/                    # Bio page components (/dashboard/bio)
│   ├── blocks/            # Content blocks (Link, Image, Text...)
│   │   ├── LinkBlock.svelte
│   │   └── index.ts
│   ├── dialogs/           # Dialogs for bio page
│   │   ├── LinkLayoutDialog.svelte
│   │   ├── LinkScheduleDialog.svelte
│   │   └── index.ts
│   ├── BioToolbar.svelte
│   ├── CalendarView.svelte
│   └── index.ts
│
├── overview/              # Dashboard home (/dashboard)
│   ├── StatsCard.svelte
│   ├── QuickActions.svelte
│   ├── RecentActivity.svelte
│   └── index.ts
│
├── preview/               # Preview components (shared)
│   ├── ProfilePreview.svelte
│   └── index.ts
│
├── shared/                # Shared across all dashboard
│   ├── EmptyState.svelte
│   ├── OnboardingTour.svelte
│   └── index.ts
│
└── index.ts              # Main barrel export
```

## 🎯 Import Examples

### Bio Page
```typescript
import { LinkBlock, BioToolbar, CalendarView } from '$lib/components/dashboard/bio';
import { LinkLayoutDialog, LinkScheduleDialog } from '$lib/components/dashboard/bio/dialogs';
```

### Overview Page
```typescript
import { StatsCard, QuickActions, RecentActivity } from '$lib/components/dashboard/overview';
```

### Shared Components
```typescript
import { ProfilePreview } from '$lib/components/dashboard/preview';
import { EmptyState, OnboardingTour } from '$lib/components/dashboard/shared';
```

### All at once
```typescript
import { 
  LinkBlock, 
  BioToolbar, 
  ProfilePreview, 
  StatsCard 
} from '$lib/components/dashboard';
```

## 📝 Naming Conventions

- **Blocks**: Suffix with `Block` (LinkBlock, ImageBlock, TextBlock)
- **Dialogs**: Suffix with `Dialog` (LinkLayoutDialog, ImageCropDialog)
- **Toolbars**: Suffix with `Toolbar` (BioToolbar, AppearanceToolbar)
- **Views**: Suffix with `View` (CalendarView, GridView)

## 🔄 Adding New Components

### Adding a new block type (e.g., ImageBlock)

1. Create component:
```bash
frontend/src/lib/components/dashboard/bio/blocks/ImageBlock.svelte
```

2. Export in index:
```typescript
// bio/blocks/index.ts
export { default as ImageBlock } from './ImageBlock.svelte';
```

3. Use in page:
```typescript
import { ImageBlock } from '$lib/components/dashboard/bio/blocks';
```

### Adding a new dashboard section (e.g., Analytics)

1. Create folder:
```bash
frontend/src/lib/components/dashboard/analytics/
```

2. Add components:
```bash
analytics/ClicksChart.svelte
analytics/TrafficChart.svelte
analytics/index.ts
```

3. Export in main index:
```typescript
// dashboard/index.ts
export * from './analytics';
```

## ✅ Benefits

- **Clear separation**: Each dashboard section has its own folder
- **Scalable**: Easy to add new sections and components
- **Maintainable**: Know exactly where to find components
- **Consistent**: Follows routes structure
- **Type-safe**: Barrel exports with TypeScript

## 🚀 Future Sections

Planned dashboard sections:
- `appearance/` - Theme & styling components
- `analytics/` - Charts & analytics components
- `settings/` - Settings page components
