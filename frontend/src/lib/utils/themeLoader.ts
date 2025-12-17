import { globalTheme, themePresets, currentHeaderStyle, defaultHeaderStyles } from '$lib/stores/theme';
import { previewStyles } from '$lib/stores/previewStyles';
import type { Profile } from '$lib/api/profile';

/**
 * Load theme from profile data into stores
 * This ensures consistent theme loading across all pages (appearance + bio)
 */
export function loadThemeFromProfile(profileData: Profile | null) {
	if (!profileData) {
		// No profile - use defaults
		globalTheme.setPreset('default');
		currentHeaderStyle.set(defaultHeaderStyles);
		syncPreviewStylesFromTheme(); // Sync preview styles
		return { themeName: 'default', category: 'cozy' };
	}

	const themeName = profileData.theme_name || 'default';

	if (themeName === 'custom') {
		// Load custom theme from custom_theme_config (includes header and custom presets)
		if (profileData.custom_theme_config) {
			try {
				const customConfig = typeof profileData.custom_theme_config === 'string'
					? JSON.parse(profileData.custom_theme_config)
					: profileData.custom_theme_config;

				if (customConfig && Object.keys(customConfig).length > 0) {
					// Extract header and custom presets from custom config
					const { header, customHeaderPresets, ...themeConfig } = customConfig;

					// Load theme config
					globalTheme.loadFromJSON(JSON.stringify(themeConfig));

					// Load header from custom config, merge with defaults to fill missing fields
					if (header && Object.keys(header).length > 0) {
						currentHeaderStyle.set({ ...defaultHeaderStyles, ...header });
					} else {
						currentHeaderStyle.set(defaultHeaderStyles);
					}
				} else {
					// Empty custom config - fallback to default
					globalTheme.setPreset('default');
					currentHeaderStyle.set(defaultHeaderStyles);
				}
			} catch (error) {
				console.warn('Failed to parse custom theme:', error);
				globalTheme.setPreset('default');
				currentHeaderStyle.set(defaultHeaderStyles);
			}
		} else {
			// No custom config - fallback to default
			globalTheme.setPreset('default');
			currentHeaderStyle.set(defaultHeaderStyles);
		}

		syncPreviewStylesFromTheme(); // Sync preview styles
		return { themeName: 'custom', category: 'custom' };
	} else {
		// Load preset theme
		globalTheme.setPreset(themeName);

		// Load header config from header_config for preset themes
		if (profileData.header_config) {
			try {
				const headerConfig = typeof profileData.header_config === 'string'
					? JSON.parse(profileData.header_config)
					: profileData.header_config;

				if (headerConfig && Object.keys(headerConfig).length > 0) {
					// Merge with defaults to fill missing fields (like avatarAlign)
					const preset = themePresets[themeName];
					const presetHeader = preset?.header || defaultHeaderStyles;
					currentHeaderStyle.set({ ...presetHeader, ...headerConfig });
				} else {
					// Use preset's default header
					const preset = themePresets[themeName];
					if (preset?.header) {
						currentHeaderStyle.set(preset.header);
					} else {
						currentHeaderStyle.set(defaultHeaderStyles);
					}
				}
			} catch (error) {
				console.warn('Failed to load header config:', error);
				// Use preset's default header
				const preset = themePresets[themeName];
				if (preset?.header) {
					currentHeaderStyle.set(preset.header);
				} else {
					currentHeaderStyle.set(defaultHeaderStyles);
				}
			}
		} else {
			// No header config - use preset's default header
			const preset = themePresets[themeName];
			if (preset?.header) {
				currentHeaderStyle.set(preset.header);
			} else {
				currentHeaderStyle.set(defaultHeaderStyles);
			}
		}

		// Find category for preset
		const themes = {
			classic: ['mcalpine', 'yoga', 'jerry'],
			vibrant: ['vibrant'],
			cozy: ['minimal', 'default'],
			bold: ['dark']
		};

		let category = 'cozy';
		for (const [cat, themeList] of Object.entries(themes)) {
			if (themeList.includes(themeName)) {
				category = cat;
				break;
			}
		}

		syncPreviewStylesFromTheme(); // Sync preview styles
		return { themeName, category };
	}
}

/**
 * Sync previewStyles from globalTheme
 * This ensures preview components show the correct theme
 */
export function syncPreviewStylesFromTheme() {
	const theme = globalTheme.getCurrent();
	// Sync all theme properties to preview store for real-time preview updates
	previewStyles.update({
		text_alignment: theme.textAlignment,
		text_size: theme.textSize,
		image_shape: theme.imageShape,
		enableCardBackground: theme.enableCardBackground,
		card_background_color: theme.cardBackground,
		card_background_opacity: theme.cardBackgroundOpacity,
		card_text_color: theme.cardTextColor,
		card_border_radius: theme.cardBorderRadius,
		show_shadow: theme.cardShadow,
		shadow_x: theme.cardShadowX,
		shadow_y: theme.cardShadowY,
		shadow_blur: theme.cardShadowBlur,
		has_card_border: theme.cardBorder,
		card_border_color: theme.cardBorderColor,
		card_border_width: theme.cardBorderWidth
	});
}
