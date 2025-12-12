-- Add link card background customization columns
-- These settings are stored on GROUP level and inherited by children when rendering

ALTER TABLE links ADD COLUMN IF NOT EXISTS has_card_background BOOLEAN DEFAULT true NOT NULL;
ALTER TABLE links ADD COLUMN IF NOT EXISTS card_background_color VARCHAR(7) DEFAULT '#ffffff';
ALTER TABLE links ADD COLUMN IF NOT EXISTS card_background_opacity INT DEFAULT 100;
ALTER TABLE links ADD COLUMN IF NOT EXISTS card_border_radius INT DEFAULT 12;

-- Add check constraints for valid values
ALTER TABLE links ADD CONSTRAINT chk_links_card_background_opacity 
  CHECK (card_background_opacity >= 0 AND card_background_opacity <= 100);

ALTER TABLE links ADD CONSTRAINT chk_links_card_border_radius 
  CHECK (card_border_radius >= 0 AND card_border_radius <= 32);

-- Add comments for documentation
COMMENT ON COLUMN links.has_card_background IS 'Enable custom background (stored on group, inherited by children)';
COMMENT ON COLUMN links.card_background_color IS 'Background color in hex format (e.g., #ffffff)';
COMMENT ON COLUMN links.card_background_opacity IS 'Background opacity percentage (0-100)';
COMMENT ON COLUMN links.card_border_radius IS 'Border radius in pixels (0-32)';
