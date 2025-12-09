const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api';

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
	async getBlocks(token: string): Promise<Block[]> {
		const res = await fetch(`${API_URL}/blocks`, {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (!res.ok) throw new Error('Failed to fetch blocks');
		return res.json();
	},

	async createBlock(data: Partial<Block>, token: string): Promise<Block> {
		const res = await fetch(`${API_URL}/blocks`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(data)
		});
		if (!res.ok) throw new Error('Failed to create block');
		return res.json();
	},

	async updateBlock(id: string, data: Partial<Block>, token: string): Promise<Block> {
		const res = await fetch(`${API_URL}/blocks/${id}`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(data)
		});
		if (!res.ok) throw new Error('Failed to update block');
		return res.json();
	},

	async deleteBlock(id: string, token: string): Promise<void> {
		const res = await fetch(`${API_URL}/blocks/${id}`, {
			method: 'DELETE',
			headers: { Authorization: `Bearer ${token}` }
		});
		if (!res.ok) throw new Error('Failed to delete block');
	},

	async reorderBlocks(blockIds: string[], token: string): Promise<void> {
		const res = await fetch(`${API_URL}/blocks/reorder`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify({ block_ids: blockIds })
		});
		if (!res.ok) throw new Error('Failed to reorder blocks');
	},

	async bulkDelete(blockIds: string[], token: string): Promise<void> {
		const res = await fetch(`${API_URL}/blocks/bulk-delete`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify({ block_ids: blockIds })
		});
		if (!res.ok) throw new Error('Failed to bulk delete blocks');
	}
};
