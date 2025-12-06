import { api } from './client';

export interface Link {
	id: string;
	profile_id: string;
	title: string;
	url: string;
	position: number;
	clicks: number;
	is_active: boolean;
	scheduled_at?: string;
	expires_at?: string;
}

export const linksApi = {
	getLinks: (token: string) => api.get<Link[]>('/links', token),
	createLink: (data: Partial<Link>, token: string) => api.post<Link>('/links', data, token),
	updateLink: (id: string, data: Partial<Link>, token: string) =>
		api.put<Link>(`/links/${id}`, data, token),
	deleteLink: (id: string, token: string) => api.delete(`/links/${id}`, token),
	reorderLinks: (linkIds: string[], token: string) =>
		api.put('/links/reorder', { link_ids: linkIds }, token)
};
