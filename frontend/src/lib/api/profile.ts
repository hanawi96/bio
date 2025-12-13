import { api } from './client';
import type { Link } from './links';
import type { Block } from './blocks';

export interface Profile {
	id: string;
	user_id: string;
	username: string;
	avatar_url?: string;
	bio?: string;
	theme_config?: string | object;
	custom_css?: string;
}

export interface ApplyThemeRequest {
	theme_config: object;
	card_styles: object;
	text_styles: string;
}

export interface ApplyThemeResponse {
	profile: Profile;
	links: Link[];
	blocks: Block[];
}

export const profileApi = {
	getMyProfile: (token: string) => api.get<Profile>('/profile', token),
	getPublicProfile: (username: string) => api.get<Profile>(`/p/${username}`),
	updateProfile: (data: Partial<Profile>, token: string) =>
		api.put<Profile>('/profile', data, token),
	
	/**
	 * Apply theme preset to profile and all groups
	 * This will update:
	 * - profile.theme_config (page styles)
	 * - All link groups (card styles)
	 * - All text groups (text styles)
	 */
	applyTheme: (data: ApplyThemeRequest, token: string) =>
		api.post<ApplyThemeResponse>('/profile/apply-theme', data, token),
	
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
