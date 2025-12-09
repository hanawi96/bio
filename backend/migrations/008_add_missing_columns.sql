-- Add missing columns to links table
ALTER TABLE links ADD COLUMN IF NOT EXISTS thumbnail_url TEXT;
ALTER TABLE links ADD COLUMN IF NOT EXISTS is_pinned BOOLEAN DEFAULT false NOT NULL;

-- Add index for pinned links
CREATE INDEX IF NOT EXISTS idx_links_is_pinned ON links(is_pinned) WHERE is_pinned = true;

-- Add comment for documentation
COMMENT ON COLUMN links.thumbnail_url IS 'URL of the thumbnail image for the link';
COMMENT ON COLUMN links.is_pinned IS 'True if this link is pinned to the top';
