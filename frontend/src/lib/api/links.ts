import { api } from './client';

export interface Link {
	id: string;
	profile_id: string;
	title: string;
	url: string;
	thumbnail_url?: string | null;
	layout_type: 'classic' | 'featured';
	position: number;
	clicks: number;
	is_active: boolean;
	is_pinned: boolean;
	scheduled_at?: string;
	expires_at?: string;
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
	reorderLinks: (linkIds: string[], token: string) =>
		api.put('/links/reorder', { link_ids: linkIds }, token),
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
	togglePin: (id: string, token: string) => api.post<Link>(`/links/${id}/pin`, {}, token)
};
