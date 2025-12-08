-- Add thumbnail_url column to links table
ALTER TABLE links ADD COLUMN IF NOT EXISTS thumbnail_url TEXT;

-- Add index for faster queries
CREATE INDEX IF NOT EXISTS idx_links_thumbnail_url ON links(thumbnail_url) WHERE thumbnail_url IS NOT NULL;
