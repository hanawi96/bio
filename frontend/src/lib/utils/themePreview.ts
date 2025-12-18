/**
 * Utility functions for generating theme preview colors
 * Ensures 100% synchronization between onboarding and appearance pages
 */

import type { ThemePreset } from '$lib/stores/theme';

export interface ThemePreviewColors {
	bg: string;
	accent: string;
	cardBg: string;
	cardText: string;
	textColor: string;
}

/**
 * Generate preview colors from theme preset
 * This is the single source of truth for preview color generation
 * 
 * @param preset - Theme preset from themePresets
 * @returns Preview colors object
 */
export function generatePreviewColors(preset: ThemePreset): ThemePreviewColors {
	if (!preset?.page || !preset?.card) {
		return {
			bg: '#ffffff',
			accent: '#000000',
			cardBg: '#ffffff',
			cardText: '#000000',
			textColor: '#000000'
		};
	}
	
	// Generate background preview
	let bgPreview = preset.page.pageBackground;
	
	if (preset.page.pageBackgroundType === 'gradient') {
		const direction = preset.page.pageGradientDirection === 'radial' 
			? 'radial-gradient(circle at center'
			: preset.page.pageGradientDirection === 'down'
			? 'linear-gradient(to bottom'
			: 'linear-gradient(to top';
		bgPreview = `${direction}, ${preset.page.pageGradientFrom}, ${preset.page.pageGradientTo})`;
	}
	
	return {
		bg: bgPreview,
		accent: preset.page.accentColor,
		cardBg: preset.card.cardBackground,
		cardText: preset.card.cardTextColor,
		textColor: preset.page.textColor
	};
}

/**
 * Simplified preview for onboarding page (only bg, card, text)
 */
export function generateSimplePreview(preset: ThemePreset): { bg: string; card: string; text: string } {
	const full = generatePreviewColors(preset);
	return {
		bg: full.bg,
		card: full.cardBg,
		text: full.cardText
	};
}
