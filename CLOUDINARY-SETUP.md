# Cloudinary Setup Guide

## 🚀 Quick Setup (5 minutes)

### Step 1: Create Cloudinary Account
1. Go to https://cloudinary.com/users/register/free
2. Sign up for FREE account (no credit card required)
3. Verify your email

### Step 2: Get API Credentials
1. Login to Cloudinary Dashboard
2. Go to **Dashboard** → You'll see:
   - Cloud Name
   - API Key
   - API Secret
3. Copy these values

### Step 3: Create Upload Preset
1. Go to **Settings** → **Upload**
2. Scroll to **Upload presets**
3. Click **Add upload preset**
4. Configure:
   - **Preset name**: `linkbio_thumbnails`
   - **Signing Mode**: `Unsigned` (for easier setup)
   - **Folder**: `linkbio/thumbnails`
   - **Allowed formats**: `jpg, png, webp, gif`
   - **Transformation**: 
     - Width: 400
     - Height: 400
     - Crop: Fill
5. Click **Save**

### Step 4: Configure Backend
1. Copy `.env.example` to `.env`:
   ```bash
   cp backend/.env.example backend/.env
   ```

2. Add your Cloudinary credentials to `backend/.env`:
   ```env
   CLOUDINARY_CLOUD_NAME=your_cloud_name_here
   CLOUDINARY_API_KEY=your_api_key_here
   CLOUDINARY_API_SECRET=your_api_secret_here
   CLOUDINARY_UPLOAD_PRESET=linkbio_thumbnails
   ```

### Step 5: Run Database Migration
```bash
# Connect to PostgreSQL
psql -U linkbio -d linkbio

# Run migration
\i backend/db/migrations/003_add_thumbnail_url.sql

# Verify
\d links
# You should see 'thumbnail_url' column
```

Or use your database tool to run:
```sql
ALTER TABLE links ADD COLUMN IF NOT EXISTS thumbnail_url TEXT;
CREATE INDEX IF NOT EXISTS idx_links_thumbnail_url ON links(thumbnail_url) WHERE thumbnail_url IS NOT NULL;
```

### Step 6: Restart Backend
```bash
cd backend
go run main.go
```

## ✅ Test Upload

1. Go to http://localhost:5173/dashboard/links
2. Click on any link's **image icon** (camera icon)
3. Upload an image
4. Should see preview and upload progress
5. Image should appear in link card

## 📊 Cloudinary Free Tier Limits

- ✅ **25 GB** storage
- ✅ **25 GB** bandwidth/month
- ✅ **Unlimited** transformations
- ✅ **Unlimited** requests

**Enough for:**
- ~25,000 thumbnail images (1MB each)
- ~250,000 page views/month (100KB per view)

## 🔧 Troubleshooting

### Error: "cloudinary credentials not configured"
- Check `.env` file has all 3 values
- Restart backend server

### Error: "Upload failed: Invalid signature"
- Check API Secret is correct
- Make sure upload preset is set to "Unsigned"

### Error: "File too large"
- Max file size is 5MB
- Compress image before upload

### Images not showing
- Check browser console for CORS errors
- Verify Cloudinary URL is accessible
- Check `thumbnail_url` in database

## 🎨 Image Optimization Tips

Cloudinary automatically:
- ✅ Converts to WebP for modern browsers
- ✅ Resizes to optimal dimensions
- ✅ Compresses without quality loss
- ✅ Serves via global CDN

You can customize URLs:
```
Original: https://res.cloudinary.com/{cloud}/image/upload/v123/linkbio/thumbnails/abc.jpg

Resize to 200x200:
https://res.cloudinary.com/{cloud}/image/upload/w_200,h_200,c_fill/v123/linkbio/thumbnails/abc.jpg

WebP format:
https://res.cloudinary.com/{cloud}/image/upload/f_webp/v123/linkbio/thumbnails/abc.jpg

Quality 80%:
https://res.cloudinary.com/{cloud}/image/upload/q_80/v123/linkbio/thumbnails/abc.jpg
```

## 📈 Monitor Usage

1. Go to Cloudinary Dashboard
2. Check **Reports** → **Usage**
3. Monitor:
   - Storage used
   - Bandwidth used
   - Transformations

## 🔐 Security Best Practices

1. **Never commit `.env` to git**
   - Already in `.gitignore`

2. **Use signed uploads in production**
   - Change upload preset to "Signed"
   - Implement signature generation in backend

3. **Set upload restrictions**
   - Max file size: 5MB
   - Allowed formats: jpg, png, webp, gif
   - Rate limiting on upload endpoint

## 🚀 Next Steps

After setup works:
1. Test with different image formats
2. Test with large files (should reject >5MB)
3. Test thumbnail display in ProfilePreview
4. Monitor Cloudinary usage dashboard

## 💡 Alternative: Local Storage

If you don't want to use Cloudinary, you can store images locally:

1. Create `backend/uploads/thumbnails/` folder
2. Modify `backend/pkg/utils/cloudinary.go` to save locally
3. Serve via `/static/uploads/` endpoint
4. Add backup strategy

But Cloudinary is recommended for:
- ✅ Better performance (CDN)
- ✅ Auto optimization
- ✅ No server storage needed
- ✅ Free tier is generous
