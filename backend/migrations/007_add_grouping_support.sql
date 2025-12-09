-- Add grouping support to links and blocks tables
-- This enables parent-child relationships for organizing items into groups

-- ============================================
-- LINKS TABLE - Add grouping columns
-- ============================================

-- Add parent_id for self-referencing (child links belong to parent group)
ALTER TABLE links ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES links(id) ON DELETE CASCADE;

-- Add is_group flag to identify group containers
ALTER TABLE links ADD COLUMN IF NOT EXISTS is_group BOOLEAN DEFAULT false NOT NULL;

-- Add group metadata
ALTER TABLE links ADD COLUMN IF NOT EXISTS group_title VARCHAR(255);
ALTER TABLE links ADD COLUMN IF NOT EXISTS group_layout VARCHAR(20) DEFAULT 'list';

-- ============================================
-- BLOCKS TABLE - Add grouping columns
-- ============================================

-- Add parent_id for self-referencing (child blocks belong to parent group)
ALTER TABLE blocks ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES blocks(id) ON DELETE CASCADE;

-- Add is_group flag to identify group containers
ALTER TABLE blocks ADD COLUMN IF NOT EXISTS is_group BOOLEAN DEFAULT false NOT NULL;

-- Add group metadata
ALTER TABLE blocks ADD COLUMN IF NOT EXISTS group_title VARCHAR(255);
ALTER TABLE blocks ADD COLUMN IF NOT EXISTS group_layout VARCHAR(20) DEFAULT 'list';

-- ============================================
-- INDEXES for performance
-- ============================================

-- Index for querying children of a group
CREATE INDEX IF NOT EXISTS idx_links_parent_id ON links(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_blocks_parent_id ON blocks(parent_id) WHERE parent_id IS NOT NULL;

-- Index for querying top-level items (groups and standalone)
CREATE INDEX IF NOT EXISTS idx_links_top_level ON links(profile_id, position) WHERE parent_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_blocks_top_level ON blocks(profile_id, position) WHERE parent_id IS NULL;

-- Index for group identification
CREATE INDEX IF NOT EXISTS idx_links_is_group ON links(is_group) WHERE is_group = true;
CREATE INDEX IF NOT EXISTS idx_blocks_is_group ON blocks(is_group) WHERE is_group = true;

-- ============================================
-- CONSTRAINTS for data integrity
-- ============================================

-- Drop existing constraints if they exist (for re-running migration)
DO $$ 
BEGIN
    ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_group_layout;
    ALTER TABLE blocks DROP CONSTRAINT IF EXISTS chk_blocks_group_layout;
    ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_group_no_parent;
    ALTER TABLE blocks DROP CONSTRAINT IF EXISTS chk_blocks_group_no_parent;
    ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_group_title;
    ALTER TABLE blocks DROP CONSTRAINT IF EXISTS chk_blocks_group_title;
EXCEPTION
    WHEN undefined_object THEN NULL;
END $$;

-- Ensure group_layout has valid values
ALTER TABLE links ADD CONSTRAINT chk_links_group_layout 
  CHECK (group_layout IN ('list', 'grid', 'carousel'));

ALTER TABLE blocks ADD CONSTRAINT chk_blocks_group_layout 
  CHECK (group_layout IN ('list', 'grid', 'carousel'));

-- Prevent nested groups (groups cannot have parent_id)
ALTER TABLE links ADD CONSTRAINT chk_links_group_no_parent 
  CHECK (NOT (is_group = true AND parent_id IS NOT NULL));

ALTER TABLE blocks ADD CONSTRAINT chk_blocks_group_no_parent 
  CHECK (NOT (is_group = true AND parent_id IS NOT NULL));

-- Ensure groups have a title
ALTER TABLE links ADD CONSTRAINT chk_links_group_title 
  CHECK (NOT (is_group = true AND (group_title IS NULL OR group_title = '')));

ALTER TABLE blocks ADD CONSTRAINT chk_blocks_group_title 
  CHECK (NOT (is_group = true AND (group_title IS NULL OR group_title = '')));

-- ============================================
-- COMMENTS for documentation
-- ============================================

COMMENT ON COLUMN links.parent_id IS 'Reference to parent group (NULL for top-level items)';
COMMENT ON COLUMN links.is_group IS 'True if this is a group container, false for regular links';
COMMENT ON COLUMN links.group_title IS 'Display title for the group (required if is_group=true)';
COMMENT ON COLUMN links.group_layout IS 'Layout style for displaying group items: list, grid, or carousel';

COMMENT ON COLUMN blocks.parent_id IS 'Reference to parent group (NULL for top-level items)';
COMMENT ON COLUMN blocks.is_group IS 'True if this is a group container, false for regular blocks';
COMMENT ON COLUMN blocks.group_title IS 'Display title for the group (required if is_group=true)';
COMMENT ON COLUMN blocks.group_layout IS 'Layout style for displaying group items: list, grid, or carousel';
