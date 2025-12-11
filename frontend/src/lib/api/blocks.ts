import { api } from './client';

export interface Block {
	id: string;
	profile_id: string;
	parent_id?: string | null;
	is_group: boolean;
	group_title?: string | null;
	group_layout: 'list' | 'grid' | 'carousel';
	block_type: 'text' | 'image' | 'video' | 'social' | 'divider' | 'email' | 'embed' | 'link';
	position: number;
	is_active: boolean;
	content?: string;
	text_style?: 'heading' | 'paragraph';
	style?: string;
	image_url?: string;
	alt_text?: string;
	video_url?: string;
	social_links?: Array<{ platform: string; url: string }>;
	divider_style?: 'line' | 'space' | 'dots';
	placeholder?: string;
	embed_url?: string;
	embed_type?: 'spotify' | 'soundcloud' | 'maps' | 'other';
	link_id?: string;
	created_at: Date;
	updated_at: Date;
	children?: Block[];
}

export const blocksApi = {
	getBlocks: (token: string) => api.get<Block[]>('/blocks', token),
	
	createBlock: (data: Partial<Block>, token: string) => 
		api.post<Block>('/blocks', data, token),
	
	updateBlock: (id: string, data: Partial<Block>, token: string) =>
		api.put<Block>(`/blocks/${id}`, data, token),
	
	deleteBlock: (id: string, token: string) => 
		api.delete(`/blocks/${id}`, token),
	
	reorderBlocks: (blockIds: string[], token: string) =>
		api.put('/blocks/reorder', { block_ids: blockIds }, token),
	
	reorderGroupBlocks: (groupId: string, blockIds: string[], token: string) =>
		api.put(`/blocks/groups/${groupId}/reorder`, { block_ids: blockIds }, token),
	
	bulkDelete: (blockIds: string[], token: string) =>
		api.post('/blocks/bulk-delete', { block_ids: blockIds }, token)
};
