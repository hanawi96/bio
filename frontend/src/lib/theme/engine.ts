/**
 * Theme Engine - Chuyển đổi tokens thành CSS
 * Engine KHÔNG render trực tiếp, mà sinh ra CSS từ tokens
 */

import type { ThemeTokens, ColorToken, ShadowToken, SpacingToken } from './tokens';
import { defaultTheme } from './tokens';

/**
 * Deep merge hai objects
 */
function deepMerge<T>(target: T, source: Partial<T>): T {
	const result = { ...target };
	for (const key in source) {
		if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
			result[key] = deepMerge(result[key] || ({} as any), source[key] as any);
		} else if (source[key] !== undefined) {
			result[key] = source[key] as any;
		}
	}
	return result;
}

/**
 * Theme Engine Class
 */
export class ThemeEngine {
	private tokens: ThemeTokens;

	constructor(customTokens?: Partial<ThemeTokens>) {
		// Merge custom tokens với default theme
		this.tokens = customTokens ? deepMerge(defaultTheme, customTokens) : defaultTheme;
	}

	/**
	 * Get raw tokens (để export hoặc lưu DB)
	 */
	getTokens(): ThemeTokens {
		return this.tokens;
	}

	/**
	 * Update tokens (partial update)
	 */
	updateTokens(updates: Partial<ThemeTokens>): void {
		this.tokens = deepMerge(this.tokens, updates);
	}

	/**
	 * Convert ColorToken to CSS color string
	 */
	private colorToCSS(color: ColorToken): string {
		const hex = color.value;
		const opacity = color.opacity ?? 100;

		if (opacity === 100) {
			return hex;
		}

		// Convert hex to rgba
		const r = parseInt(hex.slice(1, 3), 16);
		const g = parseInt(hex.slice(3, 5), 16);
		const b = parseInt(hex.slice(5, 7), 16);
		return `rgba(${r}, ${g}, ${b}, ${opacity / 100})`;
	}

	/**
	 * Convert ShadowToken to CSS box-shadow string
	 */
	private shadowToCSS(shadow: ShadowToken | null): string {
		if (!shadow) return 'none';
		const color = this.hexToRgba(shadow.color, shadow.opacity);
		return `${shadow.x}px ${shadow.y}px ${shadow.blur}px ${color}`;
	}

	/**
	 * Convert hex color to rgba
	 */
	private hexToRgba(hex: string, opacity: number): string {
		const r = parseInt(hex.slice(1, 3), 16);
		const g = parseInt(hex.slice(3, 5), 16);
		const b = parseInt(hex.slice(5, 7), 16);
		return `rgba(${r}, ${g}, ${b}, ${opacity / 100})`;
	}

	/**
	 * Convert SpacingToken to CSS padding/margin string
	 */
	private spacingToCSS(spacing: SpacingToken): string {
		if (typeof spacing === 'number') {
			return `${spacing}px`;
		}
		const { top = 0, right = 0, bottom = 0, left = 0 } = spacing;
		return `${top}px ${right}px ${bottom}px ${left}px`;
	}

	/**
	 * Generate CSS Variables từ tokens
	 * Đây là cách tối ưu nhất: tokens → CSS vars → components
	 */
	generateCSSVariables(): string {
		const vars: string[] = [];

		// Colors
		Object.entries(this.tokens.colors).forEach(([key, color]) => {
			vars.push(`--color-${key}: ${this.colorToCSS(color)};`);
		});

		// Typography
		Object.entries(this.tokens.typography).forEach(([key, typo]) => {
			vars.push(`--font-size-${key}: ${typo.fontSize}px;`);
			vars.push(`--font-weight-${key}: ${typo.fontWeight};`);
			vars.push(`--line-height-${key}: ${typo.lineHeight};`);
		});

		// Spacing
		Object.entries(this.tokens.spacing).forEach(([key, value]) => {
			vars.push(`--spacing-${key}: ${value}px;`);
		});

		// Radius
		Object.entries(this.tokens.radius).forEach(([key, value]) => {
			vars.push(`--radius-${key}: ${value}px;`);
		});

		// Shadows
		Object.entries(this.tokens.shadows).forEach(([key, shadow]) => {
			vars.push(`--shadow-${key}: ${this.shadowToCSS(shadow)};`);
		});

		// Component tokens
		vars.push(`--card-bg: ${this.colorToCSS(this.tokens.components.card.background)};`);
		vars.push(`--card-padding: ${this.spacingToCSS(this.tokens.components.card.padding)};`);
		vars.push(`--card-radius: ${this.tokens.components.card.radius}px;`);
		vars.push(`--card-shadow: ${this.shadowToCSS(this.tokens.shadows[this.tokens.components.card.shadow])};`);

		vars.push(`--link-bg: ${this.colorToCSS(this.tokens.components.link.background)};`);
		vars.push(`--link-text: ${this.colorToCSS(this.tokens.components.link.text)};`);
		vars.push(`--link-padding: ${this.spacingToCSS(this.tokens.components.link.padding)};`);
		vars.push(`--link-radius: ${this.tokens.components.link.radius}px;`);
		vars.push(`--link-shadow: ${this.shadowToCSS(this.tokens.shadows[this.tokens.components.link.shadow])};`);

		return `:root {\n  ${vars.join('\n  ')}\n}`;
	}

	/**
	 * Generate inline styles cho một component cụ thể
	 * Dùng khi cần override styles động
	 */
	getComponentStyles(component: 'card' | 'link' | 'button', overrides?: any): string {
		const comp = this.tokens.components[component];
		const styles: string[] = [];

		if (overrides?.background || comp.background) {
			const bg = overrides?.background || comp.background;
			styles.push(`background-color: ${this.colorToCSS(bg)}`);
		}

		if (overrides?.text || (comp as any).text) {
			const text = overrides?.text || (comp as any).text;
			styles.push(`color: ${this.colorToCSS(text)}`);
		}

		if (overrides?.padding !== undefined || comp.padding) {
			const padding = overrides?.padding ?? comp.padding;
			styles.push(`padding: ${this.spacingToCSS(padding)}`);
		}

		if (overrides?.radius !== undefined || comp.radius) {
			const radius = overrides?.radius ?? comp.radius;
			styles.push(`border-radius: ${radius}px`);
		}

		if (overrides?.shadow || comp.shadow) {
			const shadowKey = (overrides?.shadow || comp.shadow) as keyof ThemeTokens['shadows'];
			const shadow = this.tokens.shadows[shadowKey];
			styles.push(`box-shadow: ${this.shadowToCSS(shadow)}`);
		}

		if (overrides?.border || (comp as any).border) {
			const border = overrides?.border || (comp as any).border;
			if (border.width > 0) {
				styles.push(
					`border: ${border.width}px ${border.style} ${border.color}`
				);
			}
		}

		return styles.join('; ');
	}

	/**
	 * Get color value
	 */
	getColor(key: keyof ThemeTokens['colors']): string {
		return this.colorToCSS(this.tokens.colors[key]);
	}

	/**
	 * Get spacing value
	 */
	getSpacing(key: keyof ThemeTokens['spacing']): number {
		return this.tokens.spacing[key];
	}

	/**
	 * Get radius value
	 */
	getRadius(key: keyof ThemeTokens['radius']): number {
		return this.tokens.radius[key];
	}

	/**
	 * Get typography
	 */
	getTypography(key: keyof ThemeTokens['typography']) {
		return this.tokens.typography[key];
	}

	/**
	 * Export theme as JSON (để lưu DB hoặc file)
	 */
	exportJSON(): string {
		return JSON.stringify(this.tokens, null, 2);
	}

	/**
	 * Import theme from JSON
	 */
	static fromJSON(json: string): ThemeEngine {
		const tokens = JSON.parse(json) as Partial<ThemeTokens>;
		return new ThemeEngine(tokens);
	}
}

/**
 * Create theme engine instance
 */
export function createTheme(customTokens?: Partial<ThemeTokens>): ThemeEngine {
	return new ThemeEngine(customTokens);
}
