-- Add theme_name column to profiles table
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS theme_name VARCHAR(50);

-- Set default value for existing records
UPDATE profiles SET theme_name = 'default' WHERE theme_name IS NULL;
