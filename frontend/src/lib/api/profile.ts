import { api } from './client';

export interface Profile {
	id: string;
	user_id: string;
	username: string;
	avatar_url?: string;
	bio?: string;
	theme_config: Record<string, any>;
	custom_css?: string;
}

export const profileApi = {
	getMyProfile: (token: string) => api.get<Profile>('/profile', token),
	getPublicProfile: (username: string) => api.get<Profile>(`/p/${username}`),
	updateProfile: (data: Partial<Profile>, token: string) =>
		api.put<Profile>('/profile', data, token),
	uploadAvatar: async (file: File, token: string): Promise<Profile> => {
		const formData = new FormData();
		formData.append('avatar', file);
		const response = await fetch(`${api.baseURL}/profile/avatar`, {
			method: 'POST',
			headers: { Authorization: `Bearer ${token}` },
			body: formData
		});
		if (!response.ok) {
			const error = await response.json();
			throw new Error(error.message || 'Upload failed');
		}
		return response.json();
	}
};
