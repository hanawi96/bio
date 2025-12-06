import { api } from './client';

export interface RegisterData {
	email: string;
	username: string;
	password: string;
}

export interface LoginData {
	email: string;
	password: string;
}

export interface AuthResponse {
	user: {
		id: string;
		email: string;
		username: string;
	};
	token: string;
}

export const authApi = {
	register: (data: RegisterData) => api.post<AuthResponse>('/auth/register', data),
	login: (data: LoginData) => api.post<AuthResponse>('/auth/login', data)
};
