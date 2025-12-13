/**
 * Design Tokens - Core theme values
 * Tokens KHÔNG render CSS trực tiếp, mà là nguồn dữ liệu để sinh CSS
 */

export type ColorToken = {
	value: string; // hex color
	opacity?: number; // 0-100
};

export type SpacingToken = number | { top?: number; right?: number; bottom?: number; left?: number };

export type ShadowToken = {
	x: number;
	y: number;
	blur: number;
	color: string;
	opacity: number;
};

export type BorderToken = {
	width: number;
	style: 'solid' | 'dashed' | 'dotted' | 'none';
	color: string;
	radius: number;
};

export type TypographyToken = {
	fontSize: number; // px
	fontWeight: number | 'normal' | 'bold';
	lineHeight: number;
	letterSpacing?: number;
	textTransform?: 'none' | 'uppercase' | 'lowercase' | 'capitalize';
};

/**
 * Theme Tokens - Tất cả giá trị có thể customize
 */
export interface ThemeTokens {
	// Colors
	colors: {
		primary: ColorToken;
		secondary: ColorToken;
		accent: ColorToken;
		background: ColorToken;
		surface: ColorToken;
		text: ColorToken;
		textSecondary: ColorToken;
		border: ColorToken;
		error: ColorToken;
		success: ColorToken;
		warning: ColorToken;
	};

	// Typography Scale
	typography: {
		xs: TypographyToken;
		sm: TypographyToken;
		base: TypographyToken;
		lg: TypographyToken;
		xl: TypographyToken;
		'2xl': TypographyToken;
		'3xl': TypographyToken;
	};

	// Spacing Scale (px)
	spacing: {
		xs: number;
		sm: number;
		md: number;
		lg: number;
		xl: number;
		'2xl': number;
	};

	// Border Radius Scale (px)
	radius: {
		none: number;
		sm: number;
		md: number;
		lg: number;
		xl: number;
		full: number;
	};

	// Shadow Presets
	shadows: {
		none: ShadowToken | null;
		sm: ShadowToken;
		md: ShadowToken;
		lg: ShadowToken;
		xl: ShadowToken;
	};

	// Component-specific tokens
	components: {
		card: {
			background: ColorToken;
			padding: SpacingToken;
			radius: number;
			shadow: keyof ThemeTokens['shadows'];
			border: BorderToken;
		};
		button: {
			background: ColorToken;
			text: ColorToken;
			padding: SpacingToken;
			radius: number;
			shadow: keyof ThemeTokens['shadows'];
		};
		link: {
			background: ColorToken;
			text: ColorToken;
			padding: SpacingToken;
			radius: number;
			shadow: keyof ThemeTokens['shadows'];
			border: BorderToken;
		};
	};
}

/**
 * Default Theme - Base theme mặc định
 */
export const defaultTheme: ThemeTokens = {
	colors: {
		primary: { value: '#6366f1', opacity: 100 },
		secondary: { value: '#8b5cf6', opacity: 100 },
		accent: { value: '#ec4899', opacity: 100 },
		background: { value: '#ffffff', opacity: 100 },
		surface: { value: '#f9fafb', opacity: 100 },
		text: { value: '#111827', opacity: 100 },
		textSecondary: { value: '#6b7280', opacity: 100 },
		border: { value: '#e5e7eb', opacity: 100 },
		error: { value: '#ef4444', opacity: 100 },
		success: { value: '#10b981', opacity: 100 },
		warning: { value: '#f59e0b', opacity: 100 }
	},

	typography: {
		xs: { fontSize: 12, fontWeight: 'normal', lineHeight: 1.5 },
		sm: { fontSize: 14, fontWeight: 'normal', lineHeight: 1.5 },
		base: { fontSize: 16, fontWeight: 'normal', lineHeight: 1.5 },
		lg: { fontSize: 18, fontWeight: 'normal', lineHeight: 1.5 },
		xl: { fontSize: 20, fontWeight: 'normal', lineHeight: 1.4 },
		'2xl': { fontSize: 24, fontWeight: 'bold', lineHeight: 1.3 },
		'3xl': { fontSize: 32, fontWeight: 'bold', lineHeight: 1.2 }
	},

	spacing: {
		xs: 4,
		sm: 8,
		md: 16,
		lg: 24,
		xl: 32,
		'2xl': 48
	},

	radius: {
		none: 0,
		sm: 4,
		md: 8,
		lg: 12,
		xl: 16,
		full: 9999
	},

	shadows: {
		none: null,
		sm: { x: 0, y: 1, blur: 2, color: '#000000', opacity: 10 },
		md: { x: 0, y: 4, blur: 6, color: '#000000', opacity: 10 },
		lg: { x: 0, y: 10, blur: 15, color: '#000000', opacity: 10 },
		xl: { x: 0, y: 20, blur: 25, color: '#000000', opacity: 15 }
	},

	components: {
		card: {
			background: { value: '#ffffff', opacity: 100 },
			padding: 16,
			radius: 12,
			shadow: 'sm',
			border: { width: 1, style: 'solid', color: '#e5e7eb', radius: 12 }
		},
		button: {
			background: { value: '#6366f1', opacity: 100 },
			text: { value: '#ffffff', opacity: 100 },
			padding: { top: 12, right: 24, bottom: 12, left: 24 },
			radius: 8,
			shadow: 'sm'
		},
		link: {
			background: { value: '#ffffff', opacity: 100 },
			text: { value: '#111827', opacity: 100 },
			padding: 16,
			radius: 12,
			shadow: 'sm',
			border: { width: 1, style: 'solid', color: '#e5e7eb', radius: 12 }
		}
	}
};

/**
 * Theme Presets - Các theme có sẵn
 */
export const themePresets: Record<string, Partial<ThemeTokens>> = {
	default: defaultTheme,

	dark: {
		colors: {
			primary: { value: '#818cf8', opacity: 100 },
			secondary: { value: '#a78bfa', opacity: 100 },
			accent: { value: '#f472b6', opacity: 100 },
			background: { value: '#111827', opacity: 100 },
			surface: { value: '#1f2937', opacity: 100 },
			text: { value: '#f9fafb', opacity: 100 },
			textSecondary: { value: '#9ca3af', opacity: 100 },
			border: { value: '#374151', opacity: 100 },
			error: { value: '#f87171', opacity: 100 },
			success: { value: '#34d399', opacity: 100 },
			warning: { value: '#fbbf24', opacity: 100 }
		}
	},

	minimal: {
		colors: {
			primary: { value: '#000000', opacity: 100 },
			secondary: { value: '#666666', opacity: 100 },
			accent: { value: '#000000', opacity: 100 },
			background: { value: '#ffffff', opacity: 100 },
			surface: { value: '#fafafa', opacity: 100 },
			text: { value: '#000000', opacity: 100 },
			textSecondary: { value: '#666666', opacity: 100 },
			border: { value: '#e0e0e0', opacity: 100 },
			error: { value: '#000000', opacity: 100 },
			success: { value: '#000000', opacity: 100 },
			warning: { value: '#000000', opacity: 100 }
		},
		radius: {
			none: 0,
			sm: 0,
			md: 0,
			lg: 0,
			xl: 0,
			full: 0
		},
		shadows: {
			none: null,
			sm: { x: 0, y: 0, blur: 0, color: '#000000', opacity: 0 },
			md: { x: 0, y: 0, blur: 0, color: '#000000', opacity: 0 },
			lg: { x: 0, y: 0, blur: 0, color: '#000000', opacity: 0 },
			xl: { x: 0, y: 0, blur: 0, color: '#000000', opacity: 0 }
		}
	},

	vibrant: {
		colors: {
			primary: { value: '#f59e0b', opacity: 100 },
			secondary: { value: '#ec4899', opacity: 100 },
			accent: { value: '#8b5cf6', opacity: 100 },
			background: { value: '#fef3c7', opacity: 100 },
			surface: { value: '#ffffff', opacity: 100 },
			text: { value: '#78350f', opacity: 100 },
			textSecondary: { value: '#92400e', opacity: 100 },
			border: { value: '#fcd34d', opacity: 100 },
			error: { value: '#dc2626', opacity: 100 },
			success: { value: '#059669', opacity: 100 },
			warning: { value: '#d97706', opacity: 100 }
		},
		radius: {
			none: 0,
			sm: 8,
			md: 16,
			lg: 24,
			xl: 32,
			full: 9999
		}
	}
};
