-- Create user_themes table for storing custom themes
-- This allows users to create, save, and manage multiple custom themes

CREATE TABLE IF NOT EXISTS user_themes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Theme identification
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100),
    description TEXT,
    
    -- Theme configuration (full theme preset)
    -- Includes: page styles, card styles, text styles, header styles
    config JSONB NOT NULL,
    
    -- Preview/thumbnail
    thumbnail_url TEXT,
    
    -- Sharing features (for future community marketplace)
    is_public BOOLEAN DEFAULT false,
    downloads_count INTEGER DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT user_themes_name_unique UNIQUE(user_id, name),
    CONSTRAINT user_themes_slug_unique UNIQUE(slug),
    CONSTRAINT user_themes_name_not_empty CHECK (LENGTH(TRIM(name)) > 0),
    CONSTRAINT user_themes_config_not_null CHECK (config IS NOT NULL)
);

-- Indexes for performance
CREATE INDEX idx_user_themes_user_id ON user_themes(user_id);
CREATE INDEX idx_user_themes_public ON user_themes(is_public) WHERE is_public = true;
CREATE INDEX idx_user_themes_slug ON user_themes(slug) WHERE slug IS NOT NULL;
CREATE INDEX idx_user_themes_created_at ON user_themes(created_at DESC);

-- GIN index for JSONB config queries (if needed for searching themes by properties)
CREATE INDEX idx_user_themes_config ON user_themes USING GIN (config);

-- Comments for documentation
COMMENT ON TABLE user_themes IS 'Stores user custom themes for reusability and sharing';
COMMENT ON COLUMN user_themes.config IS 'Full theme configuration including page, card, text, and header styles';
COMMENT ON COLUMN user_themes.is_public IS 'Whether theme is published to community marketplace';
COMMENT ON COLUMN user_themes.slug IS 'URL-friendly identifier for public themes (e.g., /themes/neon-glow)';
COMMENT ON COLUMN user_themes.downloads_count IS 'Number of times theme has been downloaded/imported by other users';

-- Function to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_user_themes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update updated_at
CREATE TRIGGER trigger_user_themes_updated_at
    BEFORE UPDATE ON user_themes
    FOR EACH ROW
    EXECUTE FUNCTION update_user_themes_updated_at();

-- Function to auto-generate slug from name
CREATE OR REPLACE FUNCTION generate_theme_slug()
RETURNS TRIGGER AS $$
DECLARE
    base_slug TEXT;
    final_slug TEXT;
    counter INTEGER := 0;
BEGIN
    -- Only generate slug if is_public is true and slug is null
    IF NEW.is_public = true AND NEW.slug IS NULL THEN
        -- Convert name to slug: lowercase, replace spaces with hyphens, remove special chars
        base_slug := LOWER(REGEXP_REPLACE(TRIM(NEW.name), '[^a-zA-Z0-9\s-]', '', 'g'));
        base_slug := REGEXP_REPLACE(base_slug, '\s+', '-', 'g');
        base_slug := REGEXP_REPLACE(base_slug, '-+', '-', 'g');
        base_slug := TRIM(BOTH '-' FROM base_slug);
        
        final_slug := base_slug;
        
        -- Check for uniqueness and append counter if needed
        WHILE EXISTS (SELECT 1 FROM user_themes WHERE slug = final_slug AND id != NEW.id) LOOP
            counter := counter + 1;
            final_slug := base_slug || '-' || counter;
        END LOOP;
        
        NEW.slug := final_slug;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-generate slug
CREATE TRIGGER trigger_generate_theme_slug
    BEFORE INSERT OR UPDATE ON user_themes
    FOR EACH ROW
    EXECUTE FUNCTION generate_theme_slug();
