# Link Layout Feature ✨

## Tổng quan
Tính năng cho phép người dùng chọn layout hiển thị cho từng link, tương tự như Linktree và các đối thủ khác. Với dialog chuyên nghiệp, preview trực quan và UX mượt mà.

## Các Layout Types

### 1. Classic Layout (Mặc định)
- Hiển thị compact, gọn gàng
- Thumbnail nhỏ (80x80px) bên trái
- Title và URL hiển thị bên phải
- Phù hợp cho danh sách links nhiều

### 2. Featured Layout
- Hiển thị nổi bật với thumbnail lớn (full width, 256px height)
- Title overlay trên ảnh với gradient overlay
- Tạo sự chú ý, phù hợp cho links quan trọng
- Hiệu ứng hover scale nhẹ

## Cách sử dụng

### Trong Dashboard
1. Mở trang Links
2. Click vào icon Layout (grid icon) trên mỗi link card
3. Dialog mở ra với 2 tabs:
   - **Link Settings**: Cấu hình link (coming soon)
   - **Layout**: Chọn layout hiển thị
4. Chọn layout với preview trực quan:
   - **Classic**: Compact, efficient - với demo teal card
   - **Featured**: Large, eye-catching - với demo orange gradient card
5. Click "Save Changes" để áp dụng
6. Layout được lưu tự động và hiển thị ngay trong preview

### Trong Public Profile
- Links được render theo layout đã chọn
- Featured links nổi bật hơn với ảnh lớn
- Classic links hiển thị gọn gàng

## Technical Implementation

### Database
```sql
ALTER TABLE links ADD COLUMN layout_type VARCHAR(20) DEFAULT 'classic';
```

### Backend (Go)
- Updated `Link` model với field `layout_type`
- Updated repository queries để include layout_type
- Updated API handlers để accept layout_type updates

### Frontend (Svelte)
- Updated `LinkCard.svelte` với conditional rendering
- Added layout selector dropdown với preview
- Updated `ProfilePreview.svelte` để render đúng layout
- Updated public profile page `[username]/+page.svelte`

## UI/UX Features

### Professional Dialog
- **Tabs Navigation**: Link Settings & Layout tabs
- **Large Preview Cards**: Full-size demo của mỗi layout
- **Radio Button Selection**: Clear visual feedback
- **Hover Effects**: Scale animation trên preview cards
- **Clean Footer**: Cancel & Save buttons

### Smooth Transitions
- Dialog fade-in animation
- Card hover scale effects (1.05x)
- Border color transitions
- Shadow transitions

### Visual Feedback
- Active layout: Dark border + gray background
- Radio button: Filled circle khi selected
- Hover state: Border darkens + shadow appears
- Preview cards: Realistic demo với gradients

### Responsive Design
- Dialog responsive với mobile
- Preview cards scale appropriately
- Touch-friendly buttons
- Proper spacing và padding

## Performance
- Layout change được update qua API call
- Optimistic UI update (instant feedback)
- No page reload required
- Smooth animations với CSS transforms

## Future Enhancements
- Thêm layout types khác (grid, carousel, etc.)
- Custom layout settings (colors, fonts)
- Layout templates
- Bulk layout change
