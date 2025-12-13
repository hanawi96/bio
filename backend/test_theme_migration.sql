-- Test theme_config migration
-- Run this to verify the migration worked

-- Check if theme_config column exists
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'profiles' 
  AND column_name = 'theme_config';

-- Check if index exists
SELECT 
    indexname,
    indexdef
FROM pg_indexes 
WHERE tablename = 'profiles' 
  AND indexname = 'idx_profiles_theme_config';

-- Test insert a sample theme
-- UPDATE profiles 
-- SET theme_config = '{
--   "colors": {
--     "primary": {"value": "#6366f1", "opacity": 100}
--   }
-- }'::jsonb
-- WHERE id = 1;

-- Query theme config
-- SELECT id, username, theme_config 
-- FROM profiles 
-- WHERE theme_config IS NOT NULL;
