import { api } from './client';

export interface Link {
	id: string;
	profile_id: string;
	parent_id?: string | null;
	is_group: boolean;
	group_title?: string | null;
	group_layout: 'list' | 'grid' | 'carousel' | 'card';
	grid_columns?: number;
	grid_aspect_ratio?: string;
	title: string;
	url: string;
	description?: string | null;
	thumbnail_url?: string | null;
	image_shape?: 'sharp' | 'square' | 'circle';
	layout_type: 'classic' | 'featured' | 'carousel' | 'grid' | 'card';
	image_placement: 'left' | 'right' | 'top' | 'bottom' | 'alternating';
	text_alignment: 'left' | 'center' | 'right';
	text_size: 'S' | 'M' | 'L' | 'XL';
	show_outline: boolean;
	show_shadow: boolean;
	shadow_x?: number;
	shadow_y?: number;
	shadow_blur?: number;
	show_description: boolean;
	show_text: boolean;
	has_card_background?: boolean;
	card_background_color?: string;
	card_background_opacity?: number;
	card_border_radius?: number;
	card_text_color?: string;
	style?: string | null;
	position: number;
	clicks: number;
	is_active: boolean;
	is_pinned: boolean;
	open_in_new_tab?: boolean;
	scheduled_at?: string;
	expires_at?: string;
	children?: Link[];
}

export interface LinkFilters {
	search?: string;
	status?: 'active' | 'inactive' | '';
	layout_type?: 'classic' | 'featured' | '';
	sort_by?: 'position' | 'clicks' | 'created' | 'updated' | 'title' | '';
}

export const linksApi = {
	getLinks: (token: string, filters?: LinkFilters) => {
		const params = new URLSearchParams();
		if (filters?.search) params.append('search', filters.search);
		if (filters?.status) params.append('status', filters.status);
		if (filters?.layout_type) params.append('layout_type', filters.layout_type);
		if (filters?.sort_by) params.append('sort_by', filters.sort_by);
		const query = params.toString();
		const url = `/links${query ? '?' + query : ''}`;
		console.log('🌐 API Request URL:', url);
		console.log('📦 Filters:', filters);
		return api.get<Link[]>(url, token);
	},
	createLink: (data: Partial<Link>, token: string) => api.post<Link>('/links', data, token),
	updateLink: (id: string, data: Partial<Link>, token: string) =>
		api.put<Link>(`/links/${id}`, data, token),
	deleteLink: (id: string, token: string) => api.delete(`/links/${id}`, token),
	uploadThumbnail: async (id: string, file: File, token: string): Promise<Link> => {
		const formData = new FormData();
		formData.append('thumbnail', file);
		
		const response = await fetch(`${api.baseURL}/links/${id}/thumbnail`, {
			method: 'POST',
			headers: {
				'Authorization': `Bearer ${token}`
			},
			body: formData
		});
		
		if (!response.ok) {
			const error = await response.json();
			throw new Error(error.message || 'Upload failed');
		}
		
		return response.json();
	},
	deleteThumbnail: (id: string, token: string) => api.delete(`/links/${id}/thumbnail`, token),
	duplicateLink: (id: string, token: string) => api.post<Link>(`/links/${id}/duplicate`, {}, token),
	bulkAction: (linkIds: string[], action: 'delete' | 'activate' | 'deactivate', token: string) =>
		api.post('/links/bulk', { link_ids: linkIds, action }, token),
	togglePin: (id: string, token: string) => api.post<Link>(`/links/${id}/pin`, {}, token),
	
	// Group management
	createGroup: (title: string, layout: 'list' | 'grid' | 'carousel' | 'card', token: string) =>
		api.post<Link>('/links/groups', { title, layout }, token),
	addToGroup: (groupId: string, data: { title: string; url: string; description?: string | null }, token: string) =>
		api.post<Link>(`/links/groups/${groupId}/items`, data, token),
	moveToGroup: (linkId: string, groupId: string, token: string) =>
		api.put<Link>(`/links/${linkId}/move-to-group`, { group_id: groupId }, token),
	removeFromGroup: (linkId: string, token: string) =>
		api.put<Link>(`/links/${linkId}/remove-from-group`, {}, token),
	duplicateGroup: (groupId: string, token: string) =>
		api.post<Link>(`/links/groups/${groupId}/duplicate`, {}, token),
	reorderGroupLinks: (groupId: string, linkIds: string[], token: string) =>
		api.put(`/links/groups/${groupId}/reorder`, { link_ids: linkIds }, token)
};
