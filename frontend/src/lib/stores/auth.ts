import { writable } from 'svelte/store';
import { browser } from '$app/environment';

interface User {
	id: string;
	email: string;
	username: string;
}

interface AuthState {
	user: User | null;
	token: string | null;
	loading: boolean;
}

function getInitialState(): AuthState {
	if (!browser) {
		return { user: null, token: null, loading: false };
	}
	
	const token = localStorage.getItem('token');
	const userStr = localStorage.getItem('user');
	const user = userStr ? JSON.parse(userStr) : null;
	
	return { user, token, loading: false };
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(getInitialState());

	return {
		subscribe,
		login: (user: User, token: string) => {
			if (browser) {
				localStorage.setItem('token', token);
				localStorage.setItem('user', JSON.stringify(user));
			}
			set({ user, token, loading: false });
		},
		logout: () => {
			if (browser) {
				localStorage.removeItem('token');
				localStorage.removeItem('user');
			}
			set({ user: null, token: null, loading: false });
		},
		setLoading: (loading: boolean) => {
			update(state => ({ ...state, loading }));
		}
	};
}

export const auth = createAuthStore();
