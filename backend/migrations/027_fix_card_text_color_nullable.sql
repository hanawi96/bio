-- Fix card_text_color to be nullable for granular locking
-- This allows groups to inherit text color from theme when NULL

-- Change card_text_color to allow NULL
ALTER TABLE links ALTER COLUMN card_text_color DROP NOT NULL;

-- Set ALL existing values to NULL for groups (so they inherit from theme)
-- This resets all groups to use theme colors by default
UPDATE links 
SET card_text_color = NULL 
WHERE is_group = true;

-- Set default to NULL for new groups
ALTER TABLE links ALTER COLUMN card_text_color SET DEFAULT NULL;
