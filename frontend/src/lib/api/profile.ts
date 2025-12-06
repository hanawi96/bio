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
		api.put<Profile>('/profile', data, token)
};
