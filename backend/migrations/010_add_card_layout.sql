-- Add 'card' layout option to group_layout constraint
-- Migration 010: Add card layout support

-- Drop old constraints
ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_group_layout;
ALTER TABLE blocks DROP CONSTRAINT IF EXISTS chk_blocks_group_layout;

-- Add new constraints with 'card' option
ALTER TABLE links ADD CONSTRAINT chk_links_group_layout 
  CHECK (group_layout IN ('list', 'grid', 'carousel', 'card'));

ALTER TABLE blocks ADD CONSTRAINT chk_blocks_group_layout 
  CHECK (group_layout IN ('list', 'grid', 'carousel', 'card'));

-- Add comment
COMMENT ON CONSTRAINT chk_links_group_layout ON links IS 'Valid group layouts: list, grid, carousel, card';
COMMENT ON CONSTRAINT chk_blocks_group_layout ON blocks IS 'Valid group layouts: list, grid, carousel, card';
