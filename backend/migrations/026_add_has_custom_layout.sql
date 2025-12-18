-- Add has_custom_layout flag to track if link group has custom layout settings
-- This allows us to distinguish between theme defaults and user customizations

ALTER TABLE links ADD COLUMN IF NOT EXISTS has_custom_layout BOOLEAN DEFAULT FALSE;

-- Create index for faster queries
CREATE INDEX IF NOT EXISTS idx_links_has_custom_layout ON links(has_custom_layout) WHERE is_group = true;

-- Update existing groups that have non-null layout values to mark as custom
UPDATE links 
SET has_custom_layout = true 
WHERE is_group = true 
  AND (
    text_alignment IS NOT NULL 
    OR text_size IS NOT NULL 
    OR image_shape IS NOT NULL
  );
