-- Add show_text column to links table for carousel text visibility
ALTER TABLE links ADD COLUMN IF NOT EXISTS show_text BOOLEAN NOT NULL DEFAULT true;
