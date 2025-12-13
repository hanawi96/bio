-- Add header_config column to profiles table
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS header_config JSONB DEFAULT '{"layout":"centered","coverType":"gradient","coverColor":"#6366f1","coverGradientFrom":"#8b5cf6","coverGradientTo":"#ec4899","coverHeight":120,"avatarSize":96,"avatarBorder":4,"avatarBorderColor":"#ffffff","showCover":true,"bioAlign":"center","bioSize":"md"}'::jsonb;
