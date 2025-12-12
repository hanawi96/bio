-- Add style column to links table for background and visual customization
-- This enables the same styling capabilities as blocks (background color, opacity, border radius, etc.)

ALTER TABLE links ADD COLUMN IF NOT EXISTS style TEXT;

COMMENT ON COLUMN links.style IS 'JSON string containing style configuration: hasBackground, backgroundColor, backgroundOpacity, borderRadius, padding, etc.';
