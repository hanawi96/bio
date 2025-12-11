-- Fix old blocks without grid columns/aspect ratio
-- Ensure all existing blocks have proper default values

UPDATE blocks 
SET grid_columns = 2 
WHERE grid_columns IS NULL;

UPDATE blocks 
SET grid_aspect_ratio = '3:2' 
WHERE grid_aspect_ratio IS NULL OR grid_aspect_ratio = '';

UPDATE blocks
SET group_layout = 'list'
WHERE is_group = true AND (group_layout IS NULL OR group_layout = '');
