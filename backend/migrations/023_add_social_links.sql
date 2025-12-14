-- Add social_links column to profiles table
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS social_links TEXT;

-- Add comment
COMMENT ON COLUMN profiles.social_links IS 'JSON array of social media links';
