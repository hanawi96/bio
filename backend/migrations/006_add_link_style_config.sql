-- Add link customization columns
ALTER TABLE links ADD COLUMN IF NOT EXISTS image_placement VARCHAR(20) DEFAULT 'left';
ALTER TABLE links ADD COLUMN IF NOT EXISTS text_alignment VARCHAR(20) DEFAULT 'left';
ALTER TABLE links ADD COLUMN IF NOT EXISTS text_size VARCHAR(10) DEFAULT 'M';
ALTER TABLE links ADD COLUMN IF NOT EXISTS show_outline BOOLEAN DEFAULT false;
ALTER TABLE links ADD COLUMN IF NOT EXISTS show_shadow BOOLEAN DEFAULT false;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_links_image_placement ON links(image_placement);
CREATE INDEX IF NOT EXISTS idx_links_text_alignment ON links(text_alignment);

-- Add check constraints for valid values
ALTER TABLE links ADD CONSTRAINT chk_image_placement 
  CHECK (image_placement IN ('left', 'right', 'top', 'bottom'));

ALTER TABLE links ADD CONSTRAINT chk_text_alignment 
  CHECK (text_alignment IN ('left', 'center', 'right'));

ALTER TABLE links ADD CONSTRAINT chk_text_size 
  CHECK (text_size IN ('S', 'M', 'L', 'XL'));
