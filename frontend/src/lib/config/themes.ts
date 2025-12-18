/**
 * Single source of truth for all theme configurations
 * Used by both onboarding and dashboard appearance pages
 */

import { themePresets } from '$lib/stores/theme';
import { generateSimplePreview } from '$lib/utils/themePreview';

export interface ThemeMetadata {
	id: string;
	name: string;
	description: string;
	category: 'classic' | 'vibrant' | 'cozy' | 'bold';
	preset: string;
	preview: {
		bg: string;
		card: string;
		text: string;
	};
}

export const themeCategories = [
	{ id: 'all', label: 'All' },
	{ id: 'classic', label: 'Classic' },
	{ id: 'vibrant', label: 'Vibrant' },
	{ id: 'cozy', label: 'Cozy' },
	{ id: 'bold', label: 'Bold' },
	{ id: 'custom', label: 'Custom' }
] as const;

/**
 * Base theme metadata (without preview colors)
 * Preview colors are auto-generated from actual presets
 */
const baseThemeMetadata: Omit<ThemeMetadata, 'preview'>[] = [
	{
		id: 'default',
		name: 'Default',
		description: 'Clean and modern',
		category: 'cozy',
		preset: 'default'
	},
	{
		id: 'dark',
		name: 'Dark',
		description: 'Sleek dark mode',
		category: 'bold',
		preset: 'dark'
	},
	{
		id: 'minimal',
		name: 'Minimal',
		description: 'Simple and clean',
		category: 'cozy',
		preset: 'minimal'
	},
	{
		id: 'vibrant',
		name: 'Vibrant',
		description: 'Bold and colorful',
		category: 'vibrant',
		preset: 'vibrant'
	},
	{
		id: 'mcalpine',
		name: 'McAlpine',
		description: 'Professional dark',
		category: 'classic',
		preset: 'mcalpine'
	},
	{
		id: 'yoga',
		name: 'Yoga',
		description: 'Calm and peaceful',
		category: 'classic',
		preset: 'yoga'
	},
	{
		id: 'jerry',
		name: 'Jerry',
		description: 'Minimalist dark',
		category: 'classic',
		preset: 'jerry'
	}
];

/**
 * Theme metadata with auto-generated preview colors
 * This ensures preview colors always match actual theme presets
 * Uses centralized generateSimplePreview() for consistency
 */
export const themeMetadata: ThemeMetadata[] = baseThemeMetadata.map(base => {
	const preset = themePresets[base.preset];
	if (!preset) {
		console.warn(`Theme preset '${base.preset}' not found for theme '${base.id}'`);
		return {
			...base,
			preview: { bg: '#ffffff', card: '#ffffff', text: '#000000' }
		};
	}
	
	return {
		...base,
		preview: generateSimplePreview(preset)
	};
});

/**
 * Get themes by category
 */
export function getThemesByCategory(category: string): ThemeMetadata[] {
	if (category === 'all') {
		return themeMetadata;
	}
	if (category === 'custom') {
		return [];
	}
	return themeMetadata.filter(theme => theme.category === category);
}

/**
 * Get theme metadata by ID
 */
export function getThemeById(id: string): ThemeMetadata | undefined {
	return themeMetadata.find(theme => theme.id === id);
}

/**
 * Get themes grouped by category (for appearance page)
 */
export function getThemesGroupedByCategory() {
	const grouped: Record<string, Array<{ id: string; name: string; preset: string; isCustom?: boolean }>> = {
		classic: [],
		vibrant: [],
		cozy: [],
		bold: [],
		custom: [
			{ id: 'custom', name: 'Your Custom Theme', preset: 'custom', isCustom: true }
		]
	};

	themeMetadata.forEach(theme => {
		grouped[theme.category].push({
			id: theme.id,
			name: theme.name,
			preset: theme.preset
		});
	});

	return grouped;
}
