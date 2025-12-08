# Quick Start: Link Layout Feature

## ✅ Đã hoàn thành

### Database
- ✅ Thêm column `layout_type` vào table `links`
- ✅ Default value: 'classic'

### Backend (Go)
- ✅ Updated `Link` model
- ✅ Updated repository queries
- ✅ Updated API handlers
- ✅ Added public profile endpoint với links

### Frontend (Svelte)
- ✅ Updated `LinkCard.svelte` với 2 layouts
- ✅ Added layout selector dropdown
- ✅ Updated `ProfilePreview.svelte`
- ✅ Updated public profile page
- ✅ Smooth animations & transitions

## 🚀 Cách test

### 1. Start Backend
```bash
cd backend
go run .
```

### 2. Start Frontend
```bash
cd frontend
npm run dev
```

### 3. Test trong Dashboard
1. Login vào dashboard
2. Vào trang Links
3. Click icon Layout (grid icon) trên link card
4. Chọn "Featured" hoặc "Classic"
5. Xem preview thay đổi ngay lập tức

### 4. Test Public Profile
1. Mở `http://localhost:5173/[username]`
2. Xem links hiển thị theo layout đã chọn

## 🎨 UI Features

### Classic Layout
- Compact design
- Small thumbnail (80x80px)
- Title & URL side by side
- Perfect for many links

### Featured Layout
- Large hero image (full width, 256px height)
- Title overlay on image
- Gradient overlay effect
- Eye-catching for important links

## 🔧 Troubleshooting

### Backend không compile?
```bash
cd backend
go mod tidy
go build .
```

### Frontend có lỗi TypeScript?
```bash
cd frontend
npm run check
```

### Database migration chưa chạy?
```bash
docker exec linkbio_postgres psql -U linkbio -d linkbio -c "ALTER TABLE links ADD COLUMN IF NOT EXISTS layout_type VARCHAR(20) DEFAULT 'classic';"
```

## 📝 Notes

- Layout được lưu tự động khi thay đổi
- Không cần reload page
- Smooth animations với CSS transforms
- Mobile responsive
- Accessibility compliant (có thể thêm aria-labels nếu cần)

## 🎯 Next Steps

Có thể mở rộng:
- Thêm layout types khác (grid, carousel)
- Custom layout settings
- Bulk layout change
- Layout templates
