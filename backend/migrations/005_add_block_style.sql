-- Add style column to blocks table for text formatting
ALTER TABLE blocks ADD COLUMN IF NOT EXISTS style TEXT;
