-- Add theme configuration to profiles
-- Theme config stored as JSONB for flexibility

-- Add theme_config column to profiles table
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS theme_config JSONB DEFAULT NULL;

-- Add index for faster queries
CREATE INDEX IF NOT EXISTS idx_profiles_theme_config ON profiles USING GIN (theme_config);

-- Add comment
COMMENT ON COLUMN profiles.theme_config IS 'User theme configuration stored as JSON (Design Tokens)';

-- Example theme_config structure:
-- {
--   "colors": {
--     "primary": { "value": "#6366f1", "opacity": 100 },
--     "secondary": { "value": "#8b5cf6", "opacity": 100 },
--     ...
--   },
--   "typography": {
--     "base": { "fontSize": 16, "fontWeight": "normal", "lineHeight": 1.5 },
--     ...
--   },
--   "spacing": { "md": 16, "lg": 24, ... },
--   "radius": { "md": 8, "lg": 12, ... },
--   "shadows": { ... },
--   "components": { ... }
-- }
