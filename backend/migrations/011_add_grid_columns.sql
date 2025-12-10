-- Add grid_columns field for grid layout preset
-- Migration 011: Add grid columns configuration

-- Add grid_columns to links table
ALTER TABLE links ADD COLUMN IF NOT EXISTS grid_columns INT DEFAULT 2;
ALTER TABLE links ADD CONSTRAINT chk_links_grid_columns CHECK (grid_columns >= 1 AND grid_columns <= 4);

-- Add grid_columns to blocks table
ALTER TABLE blocks ADD COLUMN IF NOT EXISTS grid_columns INT DEFAULT 2;
ALTER TABLE blocks ADD CONSTRAINT chk_blocks_grid_columns CHECK (grid_columns >= 1 AND grid_columns <= 4);

-- Add comments
COMMENT ON COLUMN links.grid_columns IS 'Number of columns for grid layout (1-4)';
COMMENT ON COLUMN blocks.grid_columns IS 'Number of columns for grid layout (1-4)';
