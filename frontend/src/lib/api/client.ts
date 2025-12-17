const API_BASE = '/api';

interface RequestOptions extends RequestInit {
	token?: string;
}

async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
	const { token, ...fetchOptions } = options;

	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(fetchOptions.headers as Record<string, string>)
	};

	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(`${API_BASE}${endpoint}`, {
		...fetchOptions,
		headers
	});

	if (!response.ok) {
		const errorText = await response.text();
		let error;
		try {
			error = JSON.parse(errorText);
		} catch {
			error = { error: errorText || 'Request failed' };
		}
		throw new Error(error.error || error.message || `HTTP ${response.status}`);
	}

	// Handle 204 No Content or empty response
	if (response.status === 204 || response.headers.get('content-length') === '0') {
		return {} as T;
	}

	// Check if response is JSON
	const contentType = response.headers.get('content-type');
	if (contentType && contentType.includes('application/json')) {
		return response.json();
	}

	// For non-JSON responses (like plain text "OK"), return empty object
	return {} as T;
}

export const api = {
	baseURL: API_BASE,
	
	get: <T>(endpoint: string, token?: string) => 
		request<T>(endpoint, { method: 'GET', token }),
	
	post: <T>(endpoint: string, data: unknown, token?: string) =>
		request<T>(endpoint, { method: 'POST', body: JSON.stringify(data), token }),
	
	put: <T>(endpoint: string, data: unknown, token?: string) =>
		request<T>(endpoint, { method: 'PUT', body: JSON.stringify(data), token }),
	
	patch: <T>(endpoint: string, data: unknown, token?: string) =>
		request<T>(endpoint, { method: 'PATCH', body: JSON.stringify(data), token }),
	
	delete: <T>(endpoint: string, token?: string) =>
		request<T>(endpoint, { method: 'DELETE', token })
};
