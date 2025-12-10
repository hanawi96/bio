-- Add description column to links table
ALTER TABLE links ADD COLUMN IF NOT EXISTS description TEXT;
