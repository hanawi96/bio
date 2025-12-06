import { browser } from '$app/environment';
import { goto } from '$app/navigation';

// Disable SSR and prerendering for dashboard
export const ssr = false;
export const prerender = false;
export const csr = true;

export async function load() {
	if (browser) {
		const token = localStorage.getItem('token');
		if (!token) {
			goto('/auth/login');
			return {};
		}
	}
	return {};
}
