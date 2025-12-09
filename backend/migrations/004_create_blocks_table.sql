-- Create blocks table
CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    block_type VARCHAR(50) NOT NULL, -- 'text', 'image', 'video', 'social', 'divider', 'link'
    position INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- Common fields
    content TEXT,
    
    -- Text block
    text_style VARCHAR(20), -- 'heading', 'paragraph'
    
    -- Image block
    image_url TEXT,
    alt_text TEXT,
    
    -- Video block
    video_url TEXT,
    
    -- Social block (JSON array)
    social_links JSONB,
    
    -- Divider block
    divider_style VARCHAR(20), -- 'line', 'space', 'dots'
    
    -- Email collector block
    placeholder TEXT,
    
    -- Embed block
    embed_url TEXT,
    embed_type VARCHAR(50), -- 'spotify', 'soundcloud', 'maps', 'other'
    
    -- Link block (reference to links table)
    link_id UUID REFERENCES links(id) ON DELETE CASCADE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_blocks_profile_id ON blocks(profile_id);
CREATE INDEX idx_blocks_position ON blocks(profile_id, position);
CREATE INDEX idx_blocks_type ON blocks(block_type);
