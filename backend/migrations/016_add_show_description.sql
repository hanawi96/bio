-- Add show_description column to links table
-- This allows controlling visibility of link descriptions in all layouts

ALTER TABLE links ADD COLUMN IF NOT EXISTS show_description BOOLEAN DEFAULT true;

-- Add comment for documentation
COMMENT ON COLUMN links.show_description IS 'Controls visibility of link description in all layout types';
