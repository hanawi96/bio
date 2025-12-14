-- Add custom_theme_config column to profiles table
-- This field stores the user's custom theme configuration separately from preset themes
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS custom_theme_config JSONB;

-- Add comment for documentation
COMMENT ON COLUMN profiles.custom_theme_config IS 'Stores custom theme configuration when user modifies a preset theme';
