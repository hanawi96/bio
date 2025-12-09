# Text Block Modal Feature

## Tổng quan
Tính năng này cho phép người dùng nhập nội dung và tùy chỉnh text ngay khi thêm block text, thay vì tạo block trống rồi mới chỉnh sửa sau.

## Các thay đổi đã thực hiện

### 1. Frontend Components

#### TextBlockDialog.svelte (Mới)
- Modal với 3 tabs: Content, Settings, Section
- Toolbar định dạng text:
  - Font size: Small, Medium, Large
  - Text align: Left, Center, Right, Justify
  - Text style: Bold, Italic, Strikethrough, Underline
  - Color picker
- Preview real-time trong textarea
- Lưu tất cả formatting vào JSON

#### Cập nhật bio/+page.svelte
- Thêm `showTextBlockDialog` state
- Thêm `handleSaveTextBlock()` function
- Khi chọn "text" block → mở TextBlockDialog thay vì tạo block trống
- Lưu formatting vào field `style` dưới dạng JSON

#### Cập nhật TextBlock.svelte
- Parse JSON từ field `style`
- Hiển thị text với formatting đúng
- Preview trong danh sách blocks

### 2. Backend Changes

#### Models (models.go)
- Thêm field `Style *string` vào struct Block
- Lưu JSON formatting

#### Repository (block_repository.go)
- Cập nhật tất cả queries để bao gồm field `style`
- GetByUserID: SELECT thêm `style`
- Create: INSERT thêm `style`
- Update: UPDATE thêm `style`

#### Database Migration
- File: `005_add_block_style.sql`
- Thêm column `style TEXT` vào table `blocks`

### 3. API Types (blocks.ts)
- Thêm `style?: string` vào interface Block

## Cách sử dụng

### Tạo mới Text Block
1. Click nút "Add" trong bio page
2. Chọn "Text" block
3. Modal sẽ hiện ra với:
   - Tab Content: Nhập text và định dạng
   - Tab Settings: (Coming soon)
   - Tab Section: (Coming soon)
4. Tùy chỉnh text với toolbar
5. Click "Save changes"
6. Block text sẽ được tạo với formatting đã chọn

### Chỉnh sửa Text Block
1. Click nút Edit (icon bút) trên text block
2. Modal đầy đủ sẽ hiện ra với:
   - Nội dung hiện tại đã được load
   - Tất cả formatting hiện tại (font size, align, bold, italic, color, etc.)
   - Có thể chỉnh sửa và preview real-time
3. Thay đổi nội dung và formatting
4. Click "Save changes" để lưu

## Format JSON Structure

```json
{
  "fontSize": "small" | "medium" | "large",
  "textAlign": "left" | "center" | "right" | "justify",
  "isBold": boolean,
  "isItalic": boolean,
  "isUnderline": boolean,
  "isStrikethrough": boolean,
  "textColor": "#hex"
}
```

## Các file đã thay đổi

### Frontend
- `frontend/src/lib/components/dashboard/bio/dialogs/TextBlockDialog.svelte` (Mới)
- `frontend/src/routes/dashboard/bio/+page.svelte`
- `frontend/src/lib/components/dashboard/bio/blocks/TextBlock.svelte`
- `frontend/src/lib/api/blocks.ts`

### Backend
- `backend/repository/models.go`
- `backend/repository/block_repository.go`
- `backend/migrations/005_add_block_style.sql` (Mới)

## Testing

### Test tạo mới
1. Khởi động backend: `go run main.go` trong folder backend
2. Khởi động frontend: `npm run dev` trong folder frontend
3. Truy cập dashboard/bio
4. Click "Add" → chọn "Text"
5. Nhập nội dung và tùy chỉnh formatting
6. Lưu và kiểm tra preview

### Test chỉnh sửa
1. Click nút Edit trên text block đã tạo
2. Modal đầy đủ sẽ hiện ra với giá trị hiện tại
3. Thay đổi content hoặc formatting
4. Lưu và kiểm tra thay đổi được áp dụng

## Features Hoàn Thành

✅ Modal tạo mới text block với toolbar đầy đủ
✅ Modal chỉnh sửa text block với toolbar đầy đủ (thay vì modal đơn giản)
✅ Load giá trị hiện tại khi chỉnh sửa
✅ Preview real-time trong textarea
✅ Lưu formatting dưới dạng JSON
✅ Hiển thị text với formatting trong danh sách blocks
✅ Dark theme giống đối thủ
✅ 3 tabs: Content, Settings, Section
✅ Toolbar: Font size, Text align, Bold, Italic, Underline, Strikethrough, Color picker

## Notes

- Modal có dark theme giống như đối thủ
- Formatting được lưu dưới dạng JSON để dễ mở rộng
- Preview real-time trong textarea với tất cả formatting
- Tương thích với các block khác (image, video, social, etc.)
- Khi chỉnh sửa, tất cả giá trị hiện tại được load vào modal
- Modal edit và modal create có cùng giao diện và tính năng
