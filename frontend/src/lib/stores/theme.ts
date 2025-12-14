import { writable, derived } from 'svelte/store';

export interface PageStyles {
	pageBackground: string;
	pageBackgroundType: 'solid' | 'gradient' | 'image' | 'video';
	pageGradientFrom: string;
	pageGradientTo: string;
	pageGradientDirection: 'up' | 'down' | 'radial';
	pageBackgroundImage?: string;
	pageBackgroundVideo?: string;
	textColor: string;
	textSecondaryColor: string;
	accentColor: string;
}

export interface CardStyles {
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
}

export interface TextStyles {
	hasBackground: boolean;
	backgroundColor: string;
	backgroundOpacity: number;
	textColor: string;
	borderRadius: number;
	padding: number;
	shadow: 'none' | 'sm' | 'md' | 'lg' | 'xl';
	hasBorder: boolean;
	borderColor: string;
	borderWidth: number;
	borderStyle: string;
	textAlign: string;
	fontSize: string;
	isBold: boolean;
	isItalic: boolean;
	isUnderline: boolean;
	isStrikethrough: boolean;
	textTransform: string;
}

export interface HeaderStyles {
	layout: 'centered' | 'overlap' | 'card' | 'split' | 'glass' | 'gradient' | 'minimal' | 'full' | 'side' | 'asymmetric';
	coverType: 'color' | 'gradient' | 'image';
	coverColor: string;
	coverGradientFrom: string;
	coverGradientTo: string;
	coverHeight: number;
	avatarSize: number;
	avatarBorder: number;
	avatarBorderColor: string;
	avatarShape: 'circle' | 'square' | 'rounded' | 'vertical' | 'horizontal';
	showCover: boolean;
	bioAlign: 'left' | 'center' | 'right';
	bioSize: 'sm' | 'md' | 'lg';
}

export interface ThemePreset { page: PageStyles; card: CardStyles; text: TextStyles; header: HeaderStyles; }
export interface ThemeConfig extends PageStyles, CardStyles {
	pageBackgroundImage?: string;
	pageBackgroundVideo?: string;
}

export const defaultPageStyles: PageStyles = {
	pageBackground: '#ffffff', pageBackgroundType: 'gradient', pageGradientFrom: '#faf5ff',
	pageGradientTo: '#eff6ff', pageGradientDirection: 'up', textColor: '#111827', textSecondaryColor: '#6b7280', accentColor: '#6366f1'
};

export const defaultCardStyles: CardStyles = {
	enableCardBackground: true, cardBackground: '#ffffff', cardBackgroundOpacity: 100, cardTextColor: '#000000', cardBorderRadius: 12,
	cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 10, cardBorder: false,
	cardBorderColor: '#e5e7eb', cardBorderWidth: 1
};

export const defaultTextStyles: TextStyles = {
	hasBackground: false, backgroundColor: '#ffffff', backgroundOpacity: 90, textColor: '#000000',
	borderRadius: 12, padding: 16, shadow: 'none', hasBorder: false, borderColor: '#e5e7eb',
	borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium',
	isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none'
};

export const defaultHeaderStyles: HeaderStyles = {
	layout: 'centered', coverType: 'gradient', coverColor: '#6366f1',
	coverGradientFrom: '#8b5cf6', coverGradientTo: '#ec4899',
	coverHeight: 140, avatarSize: 110, avatarBorder: 4, avatarBorderColor: '#ffffff',
	avatarShape: 'circle', showCover: true, bioAlign: 'center', bioSize: 'md'
};

export const themePresets: Record<string, ThemePreset> = {
	default: { page: { ...defaultPageStyles }, card: { ...defaultCardStyles }, text: { ...defaultTextStyles }, header: { ...defaultHeaderStyles } },
	mcalpine: {
		page: { pageBackground: '#1a1a1a', pageBackgroundType: 'solid', pageGradientFrom: '#1a1a1a', pageGradientTo: '#1a1a1a', pageGradientDirection: 'up', textColor: '#ffffff', textSecondaryColor: '#cccccc', accentColor: '#ffffff' },
		card: { enableCardBackground: true, cardBackground: '#333333', cardBackgroundOpacity: 100, cardTextColor: '#ffffff', cardBorderRadius: 12, cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 10, cardBorder: false, cardBorderColor: '#444444', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#333333', backgroundOpacity: 100, textColor: '#ffffff', borderRadius: 12, padding: 16, shadow: 'md', hasBorder: false, borderColor: '#444444', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' },
		header: { layout: 'overlap', coverType: 'color', coverColor: '#000000', coverGradientFrom: '#000000', coverGradientTo: '#1a1a1a', coverHeight: 160, avatarSize: 120, avatarBorder: 4, avatarBorderColor: '#ffffff', avatarShape: 'circle', showCover: true, bioAlign: 'center', bioSize: 'md' }
	},
	yoga: {
		page: { pageBackground: '#b8c5d6', pageBackgroundType: 'solid', pageGradientFrom: '#b8c5d6', pageGradientTo: '#b8c5d6', pageGradientDirection: 'up', textColor: '#2d3748', textSecondaryColor: '#4a5568', accentColor: '#4a6fa5' },
		card: { enableCardBackground: true, cardBackground: '#ffffff', cardBackgroundOpacity: 80, cardTextColor: '#4a6fa5', cardBorderRadius: 16, cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 12, cardBorder: false, cardBorderColor: '#e2e8f0', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#ffffff', backgroundOpacity: 80, textColor: '#4a6fa5', borderRadius: 16, padding: 16, shadow: 'md', hasBorder: false, borderColor: '#e2e8f0', borderWidth: 1, borderStyle: 'solid', textAlign: 'center', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' },
		header: { layout: 'card', coverType: 'gradient', coverColor: '#b8c5d6', coverGradientFrom: '#a8b5c6', coverGradientTo: '#c8d5e6', coverHeight: 100, avatarSize: 88, avatarBorder: 4, avatarBorderColor: '#ffffff', avatarShape: 'circle', showCover: true, bioAlign: 'center', bioSize: 'md' }
	},
	jerry: {
		page: { pageBackground: '#000000', pageBackgroundType: 'solid', pageGradientFrom: '#000000', pageGradientTo: '#000000', pageGradientDirection: 'up', textColor: '#ffffff', textSecondaryColor: '#a0a0a0', accentColor: '#ffffff' },
		card: { enableCardBackground: true, cardBackground: '#1a1a1a', cardBackgroundOpacity: 100, cardTextColor: '#ffffff', cardBorderRadius: 8, cardShadow: false, cardShadowX: 0, cardShadowY: 0, cardShadowBlur: 0, cardBorder: true, cardBorderColor: '#333333', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#1a1a1a', backgroundOpacity: 100, textColor: '#ffffff', borderRadius: 8, padding: 16, shadow: 'none', hasBorder: true, borderColor: '#333333', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' },
		header: { layout: 'minimal', coverType: 'color', coverColor: '#000000', coverGradientFrom: '#000000', coverGradientTo: '#000000', coverHeight: 0, avatarSize: 72, avatarBorder: 2, avatarBorderColor: '#333333', avatarShape: 'circle', showCover: false, bioAlign: 'left', bioSize: 'sm' }
	},
	dark: {
		page: { pageBackground: '#111827', pageBackgroundType: 'solid', pageGradientFrom: '#111827', pageGradientTo: '#111827', pageGradientDirection: 'up', textColor: '#f9fafb', textSecondaryColor: '#9ca3af', accentColor: '#818cf8' },
		card: { enableCardBackground: true, cardBackground: '#1f2937', cardBackgroundOpacity: 100, cardTextColor: '#ffffff', cardBorderRadius: 12, cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 10, cardBorder: false, cardBorderColor: '#374151', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#1f2937', backgroundOpacity: 100, textColor: '#ffffff', borderRadius: 12, padding: 16, shadow: 'md', hasBorder: false, borderColor: '#374151', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' },
		header: { layout: 'glass', coverType: 'gradient', coverColor: '#1f2937', coverGradientFrom: '#4f46e5', coverGradientTo: '#7c3aed', coverHeight: 130, avatarSize: 92, avatarBorder: 3, avatarBorderColor: 'rgba(255,255,255,0.2)', avatarShape: 'circle', showCover: true, bioAlign: 'center', bioSize: 'md' }
	},
	minimal: {
		page: { pageBackground: '#ffffff', pageBackgroundType: 'solid', pageGradientFrom: '#ffffff', pageGradientTo: '#ffffff', pageGradientDirection: 'up', textColor: '#000000', textSecondaryColor: '#666666', accentColor: '#000000' },
		card: { enableCardBackground: false, cardBackground: '#ffffff', cardBackgroundOpacity: 100, cardTextColor: '#000000', cardBorderRadius: 0, cardShadow: false, cardShadowX: 0, cardShadowY: 0, cardShadowBlur: 0, cardBorder: true, cardBorderColor: '#000000', cardBorderWidth: 1 },
		text: { hasBackground: false, backgroundColor: '#ffffff', backgroundOpacity: 100, textColor: '#000000', borderRadius: 0, padding: 16, shadow: 'none', hasBorder: true, borderColor: '#000000', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' },
		header: { layout: 'side', coverType: 'color', coverColor: '#f5f5f5', coverGradientFrom: '#ffffff', coverGradientTo: '#ffffff', coverHeight: 0, avatarSize: 64, avatarBorder: 1, avatarBorderColor: '#000000', avatarShape: 'circle', showCover: false, bioAlign: 'left', bioSize: 'sm' }
	},
	vibrant: {
		page: { pageBackground: '#fef3c7', pageBackgroundType: 'solid', pageGradientFrom: '#fef3c7', pageGradientTo: '#fef3c7', pageGradientDirection: 'up', textColor: '#78350f', textSecondaryColor: '#92400e', accentColor: '#f59e0b' },
		card: { enableCardBackground: true, cardBackground: '#ffffff', cardBackgroundOpacity: 100, cardTextColor: '#78350f', cardBorderRadius: 24, cardShadow: true, cardShadowX: 0, cardShadowY: 6, cardShadowBlur: 15, cardBorder: false, cardBorderColor: '#fcd34d', cardBorderWidth: 2 },
		text: { hasBackground: true, backgroundColor: '#ffffff', backgroundOpacity: 100, textColor: '#78350f', borderRadius: 24, padding: 20, shadow: 'lg', hasBorder: false, borderColor: '#fcd34d', borderWidth: 2, borderStyle: 'solid', textAlign: 'center', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' },
		header: { layout: 'gradient', coverType: 'gradient', coverColor: '#fbbf24', coverGradientFrom: '#f59e0b', coverGradientTo: '#ef4444', coverHeight: 150, avatarSize: 104, avatarBorder: 5, avatarBorderColor: '#ffffff', avatarShape: 'circle', showCover: true, bioAlign: 'center', bioSize: 'lg' }
	}
};


export function hexToRgba(hex: string, opacity: number): string {
	if (!hex || hex.length < 7) return `rgba(0, 0, 0, ${opacity / 100})`;
	const r = parseInt(hex.slice(1, 3), 16);
	const g = parseInt(hex.slice(3, 5), 16);
	const b = parseInt(hex.slice(5, 7), 16);
	return `rgba(${r}, ${g}, ${b}, ${opacity / 100})`;
}

export function cardStylesToLinkFields(card: CardStyles, text?: TextStyles): Record<string, any> {
	// Map text block alignment to link alignment
	const textAlign = text?.textAlign || 'center';
	const alignment = textAlign === 'left' ? 'left' : textAlign === 'right' ? 'right' : 'center';
	
	// Get padding from text block (default 16)
	const padding = text?.padding || 16;

	return {
		card_background_color: card.cardBackground, card_background_opacity: card.cardBackgroundOpacity,
		card_text_color: card.cardTextColor, card_border_radius: card.cardBorderRadius,
		show_shadow: card.cardShadow, shadow_x: card.cardShadowX, shadow_y: card.cardShadowY, shadow_blur: card.cardShadowBlur,
		has_card_border: card.cardBorder, card_border_color: card.cardBorderColor, card_border_width: card.cardBorderWidth, 
		has_card_background: card.enableCardBackground ?? true,
		// Apply typography from theme
		text_alignment: alignment,
		text_size: 'M',
		image_shape: 'square',
		// Apply spacing from theme
		style: JSON.stringify({
			padding: { top: padding, right: padding, bottom: padding, left: padding },
			margin: { top: 0, bottom: 8 }
		})
	};
}

export function textStylesToBlockStyle(text: TextStyles): string {
	return JSON.stringify({
		hasBackground: text.hasBackground, backgroundColor: text.backgroundColor, backgroundOpacity: text.backgroundOpacity,
		textColor: text.textColor, borderRadius: text.borderRadius, padding: text.padding, shadow: text.shadow,
		hasBorder: text.hasBorder, borderColor: text.borderColor, borderWidth: text.borderWidth, borderStyle: text.borderStyle,
		textAlign: text.textAlign, fontSize: text.fontSize, isBold: text.isBold, isItalic: text.isItalic,
		isUnderline: text.isUnderline, isStrikethrough: text.isStrikethrough, textTransform: text.textTransform
	});
}

export function applyCardStylesToGroup(group: any, card: CardStyles): any {
	return {
		...group, card_background_color: card.cardBackground, card_background_opacity: card.cardBackgroundOpacity,
		card_text_color: card.cardTextColor, card_border_radius: card.cardBorderRadius, show_shadow: card.cardShadow,
		shadow_x: card.cardShadowX, shadow_y: card.cardShadowY, shadow_blur: card.cardShadowBlur,
		has_card_border: card.cardBorder, card_border_color: card.cardBorderColor, card_border_width: card.cardBorderWidth, 
		has_card_background: card.enableCardBackground ?? true
	};
}

export function applyTextStylesToGroup(group: any, text: TextStyles): any {
	return { ...group, style: textStylesToBlockStyle(text) };
}

const defaultTheme: ThemeConfig = { ...defaultPageStyles, ...defaultCardStyles };

function createThemeStore() {
	const { subscribe, set, update } = writable<ThemeConfig>(defaultTheme);
	return {
		subscribe,
		setPreset(presetName: string) { const preset = themePresets[presetName]; if (preset) set({ ...preset.page, ...preset.card }); },
		getPreset(presetName: string): ThemePreset | null { return themePresets[presetName] || null; },
		update(partial: Partial<ThemeConfig>) { update(current => ({ ...current, ...partial })); },
		reset() { set(defaultTheme); },
		loadFromJSON(json: string | object | null | undefined) {
			if (!json) { set(defaultTheme); return; }
			try { const config = typeof json === 'string' ? JSON.parse(json) : json; set({ ...defaultTheme, ...config }); }
			catch (e) { console.warn('Failed to parse theme JSON:', e); set(defaultTheme); }
		},
		exportJSON(): string { let current: ThemeConfig = defaultTheme; subscribe(v => current = v)(); return JSON.stringify(current); },
		getCurrent(): ThemeConfig { let current: ThemeConfig = defaultTheme; subscribe(v => current = v)(); return current; }
	};
}

export const globalTheme = createThemeStore();

// Current header style store
export const currentHeaderStyle = writable<HeaderStyles>(defaultHeaderStyles);

export const themeStyles = derived(globalTheme, ($theme) => {
	// Always respect the current pageBackgroundType setting
	const bgType = $theme.pageBackgroundType || 'solid';
	let pageBackground = $theme.pageBackground || '#ffffff';

	// Generate background based on type
	if (bgType === 'gradient') {
		const from = $theme.pageGradientFrom || '#faf5ff';
		const to = $theme.pageGradientTo || '#eff6ff';
		const direction = $theme.pageGradientDirection || 'up';
		
		if (direction === 'radial') {
			pageBackground = `radial-gradient(circle at center, ${from}, ${to})`;
		} else if (direction === 'down') {
			pageBackground = `linear-gradient(to bottom, ${from}, ${to})`;
		} else {
			pageBackground = `linear-gradient(to top, ${from}, ${to})`;
		}
	} else if (bgType === 'image' && $theme.pageBackgroundImage) {
		pageBackground = `url(${$theme.pageBackgroundImage})`;
	} else if (bgType === 'video' && $theme.pageBackgroundVideo) {
		// For video, use solid color as fallback
		pageBackground = $theme.pageBackground || '#000000';
	}
	// For 'solid' type, just use pageBackground as-is

	const enableBg = $theme.enableCardBackground ?? true;
	const bgColor = enableBg ? hexToRgba($theme.cardBackground, $theme.cardBackgroundOpacity) : 'transparent';

	return {
		pageBackground,
		pageBackgroundType: bgType,
		pageBackgroundVideo: $theme.pageBackgroundVideo,
		pageBackgroundImage: $theme.pageBackgroundImage,
		cardStyle: `background-color: ${bgColor}; color: ${$theme.cardTextColor}; border-radius: ${$theme.cardBorderRadius}px; ${$theme.cardShadow ? `box-shadow: ${$theme.cardShadowX}px ${$theme.cardShadowY}px ${$theme.cardShadowBlur}px rgba(0,0,0,0.1);` : ''} ${$theme.cardBorder ? `border: ${$theme.cardBorderWidth}px solid ${$theme.cardBorderColor};` : ''}`.trim()
	};
});
