-- Add layout_type column to links table
-- This allows links to be displayed in different layouts (classic or featured)

ALTER TABLE links ADD COLUMN IF NOT EXISTS layout_type VARCHAR(20) DEFAULT 'classic';

-- Add comment for documentation
COMMENT ON COLUMN links.layout_type IS 'Layout type for link display: classic (compact) or featured (large with prominent image)';
