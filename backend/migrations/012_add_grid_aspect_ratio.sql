-- Add grid_aspect_ratio field for grid layout
-- Migration 012: Add aspect ratio configuration

-- Add grid_aspect_ratio to links table
ALTER TABLE links ADD COLUMN IF NOT EXISTS grid_aspect_ratio VARCHAR(10) DEFAULT '3:2';
ALTER TABLE links ADD CONSTRAINT chk_links_grid_aspect_ratio CHECK (grid_aspect_ratio IN ('1:1', '3:2', '16:9', '3:1', '2:3'));

-- Add grid_aspect_ratio to blocks table
ALTER TABLE blocks ADD COLUMN IF NOT EXISTS grid_aspect_ratio VARCHAR(10) DEFAULT '3:2';
ALTER TABLE blocks ADD CONSTRAINT chk_blocks_grid_aspect_ratio CHECK (grid_aspect_ratio IN ('1:1', '3:2', '16:9', '3:1', '2:3'));

-- Add comments
COMMENT ON COLUMN links.grid_aspect_ratio IS 'Aspect ratio for grid layout thumbnails (1:1, 3:2, 16:9, 3:1, 2:3)';
COMMENT ON COLUMN blocks.grid_aspect_ratio IS 'Aspect ratio for grid layout thumbnails (1:1, 3:2, 16:9, 3:1, 2:3)';
