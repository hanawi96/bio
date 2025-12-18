-- Quick fix: Reset all card_text_color to NULL for groups
-- Run this SQL directly in your database

-- Step 1: Make column nullable (if not already)
ALTER TABLE links ALTER COLUMN card_text_color DROP NOT NULL;

-- Step 2: Reset all groups to NULL (inherit from theme)
UPDATE links 
SET card_text_color = NULL 
WHERE is_group = true;

-- Step 3: Set default to NULL
ALTER TABLE links ALTER COLUMN card_text_color SET DEFAULT NULL;

-- Verify the change
SELECT id, group_title, card_text_color, is_group 
FROM links 
WHERE is_group = true 
LIMIT 10;
