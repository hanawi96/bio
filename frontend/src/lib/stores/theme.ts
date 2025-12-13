import { writable, derived } from 'svelte/store';

export interface PageStyles {
pageBackground: string;
pageBackgroundType: 'solid' | 'gradient';
pageGradientFrom: string;
pageGradientTo: string;
textColor: string;
textSecondaryColor: string;
accentColor: string;
}

export interface CardStyles {
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

export interface ThemePreset { page: PageStyles; card: CardStyles; text: TextStyles; }
export interface ThemeConfig extends PageStyles, CardStyles {}

export const defaultPageStyles: PageStyles = {
pageBackground: '#ffffff', pageBackgroundType: 'gradient', pageGradientFrom: '#faf5ff',
pageGradientTo: '#eff6ff', textColor: '#111827', textSecondaryColor: '#6b7280', accentColor: '#6366f1'
};

export const defaultCardStyles: CardStyles = {
cardBackground: '#ffffff', cardBackgroundOpacity: 100, cardTextColor: '#000000', cardBorderRadius: 12,
cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 10, cardBorder: false,
cardBorderColor: '#e5e7eb', cardBorderWidth: 1
};

export const defaultTextStyles: TextStyles = {
hasBackground: false, backgroundColor: '#ffffff', backgroundOpacity: 90, textColor: '#000000',
borderRadius: 12, padding: 16, shadow: 'none', hasBorder: false, borderColor: '#e5e7eb',
borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium',
isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none'
};

export const themePresets: Record<string, ThemePreset> = {
	default: { page: { ...defaultPageStyles }, card: { ...defaultCardStyles }, text: { ...defaultTextStyles } },
	mcalpine: {
		page: { pageBackground: '#1a1a1a', pageBackgroundType: 'solid', pageGradientFrom: '#1a1a1a', pageGradientTo: '#1a1a1a', textColor: '#ffffff', textSecondaryColor: '#cccccc', accentColor: '#ffffff' },
		card: { cardBackground: '#333333', cardBackgroundOpacity: 100, cardTextColor: '#ffffff', cardBorderRadius: 12, cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 10, cardBorder: false, cardBorderColor: '#444444', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#333333', backgroundOpacity: 100, textColor: '#ffffff', borderRadius: 12, padding: 16, shadow: 'md', hasBorder: false, borderColor: '#444444', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' }
	},
	yoga: {
		page: { pageBackground: '#b8c5d6', pageBackgroundType: 'solid', pageGradientFrom: '#b8c5d6', pageGradientTo: '#b8c5d6', textColor: '#2d3748', textSecondaryColor: '#4a5568', accentColor: '#4a6fa5' },
		card: { cardBackground: '#ffffff', cardBackgroundOpacity: 80, cardTextColor: '#4a6fa5', cardBorderRadius: 16, cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 12, cardBorder: false, cardBorderColor: '#e2e8f0', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#ffffff', backgroundOpacity: 80, textColor: '#4a6fa5', borderRadius: 16, padding: 16, shadow: 'md', hasBorder: false, borderColor: '#e2e8f0', borderWidth: 1, borderStyle: 'solid', textAlign: 'center', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' }
	},
	jerry: {
		page: { pageBackground: '#000000', pageBackgroundType: 'solid', pageGradientFrom: '#000000', pageGradientTo: '#000000', textColor: '#ffffff', textSecondaryColor: '#a0a0a0', accentColor: '#ffffff' },
		card: { cardBackground: '#1a1a1a', cardBackgroundOpacity: 100, cardTextColor: '#ffffff', cardBorderRadius: 8, cardShadow: false, cardShadowX: 0, cardShadowY: 0, cardShadowBlur: 0, cardBorder: true, cardBorderColor: '#333333', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#1a1a1a', backgroundOpacity: 100, textColor: '#ffffff', borderRadius: 8, padding: 16, shadow: 'none', hasBorder: true, borderColor: '#333333', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' }
	},
	dark: {
		page: { pageBackground: '#111827', pageBackgroundType: 'solid', pageGradientFrom: '#111827', pageGradientTo: '#111827', textColor: '#f9fafb', textSecondaryColor: '#9ca3af', accentColor: '#818cf8' },
		card: { cardBackground: '#1f2937', cardBackgroundOpacity: 100, cardTextColor: '#ffffff', cardBorderRadius: 12, cardShadow: true, cardShadowX: 0, cardShadowY: 4, cardShadowBlur: 10, cardBorder: false, cardBorderColor: '#374151', cardBorderWidth: 1 },
		text: { hasBackground: true, backgroundColor: '#1f2937', backgroundOpacity: 100, textColor: '#ffffff', borderRadius: 12, padding: 16, shadow: 'md', hasBorder: false, borderColor: '#374151', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' }
	},
	minimal: {
		page: { pageBackground: '#ffffff', pageBackgroundType: 'solid', pageGradientFrom: '#ffffff', pageGradientTo: '#ffffff', textColor: '#000000', textSecondaryColor: '#666666', accentColor: '#000000' },
		card: { cardBackground: '#ffffff', cardBackgroundOpacity: 100, cardTextColor: '#000000', cardBorderRadius: 0, cardShadow: false, cardShadowX: 0, cardShadowY: 0, cardShadowBlur: 0, cardBorder: true, cardBorderColor: '#000000', cardBorderWidth: 1 },
		text: { hasBackground: false, backgroundColor: '#ffffff', backgroundOpacity: 100, textColor: '#000000', borderRadius: 0, padding: 16, shadow: 'none', hasBorder: true, borderColor: '#000000', borderWidth: 1, borderStyle: 'solid', textAlign: 'left', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' }
	},
	vibrant: {
		page: { pageBackground: '#fef3c7', pageBackgroundType: 'solid', pageGradientFrom: '#fef3c7', pageGradientTo: '#fef3c7', textColor: '#78350f', textSecondaryColor: '#92400e', accentColor: '#f59e0b' },
		card: { cardBackground: '#ffffff', cardBackgroundOpacity: 100, cardTextColor: '#78350f', cardBorderRadius: 24, cardShadow: true, cardShadowX: 0, cardShadowY: 6, cardShadowBlur: 15, cardBorder: false, cardBorderColor: '#fcd34d', cardBorderWidth: 2 },
		text: { hasBackground: true, backgroundColor: '#ffffff', backgroundOpacity: 100, textColor: '#78350f', borderRadius: 24, padding: 20, shadow: 'lg', hasBorder: false, borderColor: '#fcd34d', borderWidth: 2, borderStyle: 'solid', textAlign: 'center', fontSize: 'text-medium', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none' }
	}
};


export function hexToRgba(hex: string, opacity: number): string {
	if (!hex || hex.length < 7) return `rgba(0, 0, 0, ${opacity / 100})`;
	const r = parseInt(hex.slice(1, 3), 16);
	const g = parseInt(hex.slice(3, 5), 16);
	const b = parseInt(hex.slice(5, 7), 16);
	return `rgba(${r}, ${g}, ${b}, ${opacity / 100})`;
}

export function cardStylesToLinkFields(card: CardStyles): Record<string, any> {
	return {
		card_background_color: card.cardBackground, card_background_opacity: card.cardBackgroundOpacity,
		card_text_color: card.cardTextColor, card_border_radius: card.cardBorderRadius,
		show_shadow: card.cardShadow, shadow_x: card.cardShadowX, shadow_y: card.cardShadowY, shadow_blur: card.cardShadowBlur,
		has_card_border: card.cardBorder, card_border_color: card.cardBorderColor, card_border_width: card.cardBorderWidth, has_card_background: true
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
	return { ...group, card_background_color: card.cardBackground, card_background_opacity: card.cardBackgroundOpacity,
		card_text_color: card.cardTextColor, card_border_radius: card.cardBorderRadius, show_shadow: card.cardShadow,
		shadow_x: card.cardShadowX, shadow_y: card.cardShadowY, shadow_blur: card.cardShadowBlur,
		has_card_border: card.cardBorder, card_border_color: card.cardBorderColor, card_border_width: card.cardBorderWidth, has_card_background: true };
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

export const themeStyles = derived(globalTheme, ($theme) => ({
	pageBackground: $theme.pageBackgroundType === 'gradient'
		? `linear-gradient(to bottom right, ${$theme.pageGradientFrom}, ${$theme.pageGradientTo})`
		: $theme.pageBackground,
	cardStyle: `background-color: ${hexToRgba($theme.cardBackground, $theme.cardBackgroundOpacity)}; color: ${$theme.cardTextColor}; border-radius: ${$theme.cardBorderRadius}px; ${$theme.cardShadow ? `box-shadow: ${$theme.cardShadowX}px ${$theme.cardShadowY}px ${$theme.cardShadowBlur}px rgba(0,0,0,0.1);` : ''} ${$theme.cardBorder ? `border: ${$theme.cardBorderWidth}px solid ${$theme.cardBorderColor};` : ''}`.trim()
}));
