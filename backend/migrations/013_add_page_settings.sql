-- Add page settings columns to profiles table
ALTER TABLE profiles 
ADD COLUMN IF NOT EXISTS show_share_button BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS show_subscribe_button BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS hide_branding BOOLEAN DEFAULT false;
