-- Update image_placement constraint to include 'alternating'
ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_image_placement;

ALTER TABLE links ADD CONSTRAINT chk_image_placement 
CHECK (image_placement IN ('left', 'right', 'top', 'bottom', 'alternating'));
