-- Remove link_id column from blocks table
-- This column was not needed as blocks and link groups are separate concepts

ALTER TABLE blocks DROP COLUMN IF EXISTS link_id;
