-- Add image_shape column to links table
ALTER TABLE links 
ADD COLUMN IF NOT EXISTS image_shape VARCHAR(20) DEFAULT 'square';

-- Valid values: 'square', 'circle'
COMMENT ON COLUMN links.image_shape IS 'Shape of the thumbnail image: square or circle';
