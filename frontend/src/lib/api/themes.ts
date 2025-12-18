import { api } from './client';

export interface UserTheme {
	id: string;
	user_id: string;
	name: string;
	slug?: string;
	description?: string;
	config: ThemeConfig;
	thumbnail_url?: string;
	is_public: boolean;
	downloads_count: number;
	created_at: string;
	updated_at: string;
}

export interface ThemeConfig {
	// Page styles
	pageBackground: string;
	pageBackgroundType: 'solid' | 'gradient' | 'image' | 'video';
	pageGradientFrom: string;
	pageGradientTo: string;
	pageGradientDirection: 'up' | 'down' | 'radial';
	textColor: string;
	textSecondaryColor: string;
	accentColor: string;
	
	// Card styles
	enableCardBackground: boolean;
	cardBackground: string;
	cardBackgroundOpacity: number;
	cardTextColor: string;
	cardBorderRadius: number;
	cardShadow: boolean;
	cardShadowX: number;
	cardShadowY: number;
	cardShadowBlur: number;
	cardBorder: boolean;
	cardBorderColor: string;
	cardBorderWidth: number;
	cardPadding: number;
	cardSpacing: number;
	
	// Typography
	textAlignment?: 'left' | 'center' | 'right';
	textSize?: 'S' | 'M' | 'L' | 'XL';
	imageShape?: 'sharp' | 'square' | 'circle';
	
	// Header styles
	header?: {
		layout: string;
		coverType: string;
		coverColor: string;
		coverGradientFrom: string;
		coverGradientTo: string;
		coverHeight: number;
		avatarSize: number;
		avatarBorder: number;
		avatarBorderColor: string;
		avatarShape: string;
		avatarAlign: string;
		showCover: boolean;
		bioAlign: string;
		bioSize: string;
		bioTextColor: string;
	};
	
	// Custom header presets
	customHeaderPresets?: any[];
}

export interface CreateThemeRequest {
	name: string;
	description?: string;
	config: ThemeConfig;
}

export interface ImportThemeRequest {
	name: string;
	description?: string;
	config: ThemeConfig;
	source_theme_id?: string;
}

export interface UpdateThemeRequest {
	name?: string;
	description?: string;
	config?: ThemeConfig;
	thumbnail_url?: string;
	is_public?: boolean;
}

export const themesApi = {
	// Get user's custom themes
	getMyThemes: (token: string) => 
		api.get<UserTheme[]>('/themes/my', token),
	
	// Get a specific theme
	getThemeById: (themeId: string, token: string) => 
		api.get<UserTheme>(`/themes/${themeId}`, token),
	
	// Create a new theme
	createTheme: (data: CreateThemeRequest, token: string) => 
		api.post<UserTheme>('/themes', data, token),
	
	// Update a theme
	updateTheme: (themeId: string, data: UpdateThemeRequest, token: string) => 
		api.put<UserTheme>(`/themes/${themeId}`, data, token),
	
	// Delete a theme
	deleteTheme: (themeId: string, token: string) => 
		api.delete(`/themes/${themeId}`, token),
	
	// Export theme to JSON
	exportTheme: (themeId: string, token: string) => 
		api.get<{ name: string; description?: string; config: ThemeConfig; exported_at: string }>(`/themes/${themeId}/export`, token),
	
	// Import theme from JSON
	importTheme: (data: ImportThemeRequest, token: string) => 
		api.post<UserTheme>('/themes/import', data, token),
	
	// Publish theme (make public)
	publishTheme: (themeId: string, token: string) => 
		api.post<UserTheme>(`/themes/${themeId}/publish`, {}, token),
	
	// Unpublish theme (make private)
	unpublishTheme: (themeId: string, token: string) => 
		api.post<UserTheme>(`/themes/${themeId}/unpublish`, {}, token),
	
	// Get public themes (marketplace)
	getPublicThemes: (limit = 20, offset = 0) => 
		api.get<UserTheme[]>(`/themes/public?limit=${limit}&offset=${offset}`),
	
	// Get public theme by slug
	getThemeBySlug: (slug: string) => 
		api.get<UserTheme>(`/themes/slug/${slug}`)
};
