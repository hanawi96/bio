-- Remove default values for text_alignment, text_size, image_shape
-- This allows new groups to inherit from theme instead of having hardcoded defaults

-- Change text_alignment default from 'left' to NULL
ALTER TABLE links ALTER COLUMN text_alignment DROP DEFAULT;
ALTER TABLE links ALTER COLUMN text_alignment SET DEFAULT NULL;

-- Change text_size default from 'M' to NULL
ALTER TABLE links ALTER COLUMN text_size DROP DEFAULT;
ALTER TABLE links ALTER COLUMN text_size SET DEFAULT NULL;

-- Change image_placement default from 'left' to NULL
ALTER TABLE links ALTER COLUMN image_placement DROP DEFAULT;
ALTER TABLE links ALTER COLUMN image_placement SET DEFAULT NULL;

-- Update existing groups that have default values to NULL (so they inherit from theme)
-- Only update groups that haven't been explicitly customized
-- We can identify these by checking if they have the exact default values
UPDATE links 
SET text_alignment = NULL 
WHERE is_group = true 
  AND text_alignment = 'left'
  AND updated_at = created_at;  -- Only update if never modified

UPDATE links 
SET text_size = NULL 
WHERE is_group = true 
  AND text_size = 'M'
  AND updated_at = created_at;

UPDATE links 
SET image_placement = NULL 
WHERE is_group = true 
  AND image_placement = 'left'
  AND updated_at = created_at;
