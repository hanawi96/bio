<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { 
		globalTheme, 
		themePresets, 
		themeStyles,
		currentHeaderStyle,
		cardStylesToLinkFields,
		textStylesToBlockStyle,
		applyCardStylesToGroup,
		applyTextStylesToGroup
	} from '$lib/stores/theme';
	import { previewStyles } from '$lib/stores/previewStyles';
	import { pendingChanges } from '$lib/stores/pendingChanges';
	import type { CardStyles, TextStyles, ThemePreset } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import { blocksApi } from '$lib/api/blocks';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';
	import type { Block } from '$lib/api/blocks';
	import { toast } from 'svelte-sonner';
	import ProfilePreview from '$lib/components/dashboard/preview/ProfilePreview.svelte';
	import HeaderStyleEditor from '$lib/components/dashboard/appearance/HeaderStyleEditor.svelte';
	import BackgroundSelector from '$lib/components/dashboard/appearance/BackgroundSelector.svelte';
	import SocialLinksEditor from '$lib/components/dashboard/appearance/SocialLinksEditor.svelte';
	import PageSettingsEditor from '$lib/components/dashboard/appearance/PageSettingsEditor.svelte';
	import CardStyleEditor from '$lib/components/dashboard/appearance/CardStyleEditor.svelte';
	import TypographyEditor from '$lib/components/dashboard/appearance/TypographyEditor.svelte';
	import ImageStyleEditor from '$lib/components/dashboard/appearance/ImageStyleEditor.svelte';
	import SpacingEditor from '$lib/components/dashboard/appearance/SpacingEditor.svelte';
	import ShadowEditor from '$lib/components/dashboard/appearance/ShadowEditor.svelte';
	import BorderEditor from '$lib/components/dashboard/appearance/BorderEditor.svelte';

	let profile = $state<Profile | null>(null);
	let links = $state<Link[]>([]);
	let blocks = $state<Block[]>([]);
	let selectedCategory = $state<string>('all');
	let currentTheme = $state<string>('default');
	let saving = $state(false);
	let applying = $state(false);
	
	const hasUnsavedChanges = $derived($pendingChanges.hasChanges);
	
	// Profile edit
	let bio = $state('');
	let avatarFile = $state<File | null>(null);
	let avatarPreview = $state<string>('');
	let uploadingAvatar = $state(false);
	let savingProfile = $state(false);
	let showProfileModal = $state(false);
	let socialLinks = $state<Array<{ platform: string; url: string }>>([]);
	
	// Page settings - use object for better reactivity
	let pageSettings = $state({
		showShareButton: true,
		showSubscribeButton: true,
		hideBranding: false
	});
	


	const presetColorsCache = new Map<string, ReturnType<typeof getPresetColors>>();

	// Check if current theme is customized
	function isThemeCustomized(): boolean {
		if (!currentTheme || currentTheme.startsWith('custom-')) return false;
		
		const currentConfig = globalTheme.getCurrent();
		const preset = themePresets[currentTheme];
		
		if (!preset) return false;
		
		// Compare key properties
		const isModified = 
			currentConfig.pageBackground !== preset.page.pageBackground ||
			currentConfig.pageBackgroundType !== preset.page.pageBackgroundType ||
			currentConfig.cardBackground !== preset.card.cardBackground ||
			currentConfig.cardTextColor !== preset.card.cardTextColor ||
			currentConfig.cardBorderRadius !== preset.card.cardBorderRadius;
		
		return isModified;
	}

	// Auto-switch to custom when user modifies theme
	$effect(() => {
		if (hasUnsavedChanges && isThemeCustomized() && currentTheme !== 'custom') {
			selectedCategory = 'custom';
			currentTheme = 'custom';
		}
	});

	async function saveAllChanges() {
		if (saving || !hasUnsavedChanges) return;
		
		saving = true;
		try {
			const changes = pendingChanges.getAll();
			const token = get(auth).token!;
			
			console.log('[SAVE] saveAllChanges called:', {
				currentTheme,
				hasThemeChanges: !!changes.theme,
				themeConfig: globalTheme.getCurrent()
			});
			
			// Save theme changes
			if (changes.theme) {
				// Only save theme config if it's custom theme
				// For preset themes, just save the theme name
				if (currentTheme === 'custom') {
					const updatedTheme = globalTheme.getCurrent();
					console.log('[SAVE] Saving CUSTOM theme config:', updatedTheme);
					await profileApi.updateProfile({ 
						theme_config: JSON.stringify(updatedTheme),
						theme_name: 'custom'
					}, token);
				} else {
					console.log('[SAVE] Saving PRESET theme name only:', currentTheme);
					await profileApi.updateProfile({ 
						theme_name: currentTheme
					}, token);
				}
			}
			
			// Save header changes
			if (changes.header) {
				await profileApi.updateProfile({ 
					header_config: JSON.stringify(get(currentHeaderStyle)) 
				}, token);
			}
			
			// Save link styles
			if (changes.linkStyles) {
				await linksApi.updateAllGroupStyles(changes.linkStyles, token);
			}
			
			// Reload data
			links = await linksApi.getLinks(token);
			syncPreviewStyles(links);
			
			// If currently on custom theme, keep it selected after save
			if (currentTheme === 'custom') {
				selectedCategory = 'custom';
			}
			
			pendingChanges.reset();
			toast.success('All changes saved!');
		} catch (e: any) {
			toast.error(e.message || 'Failed to save changes');
		} finally {
			saving = false;
		}
	}

	const categories = [
		{ id: 'all', label: 'All' },
		{ id: 'classic', label: 'Classic' },
		{ id: 'vibrant', label: 'Vibrant' },
		{ id: 'cozy', label: 'Cozy' },
		{ id: 'bold', label: 'Bold' },
		{ id: 'custom', label: 'Custom' }
	];

	// Helper to sync previewStyles from links data
	function syncPreviewStyles(linksData: Link[]) {
		const firstGroup = linksData?.find(l => l.is_group);
		if (firstGroup) {
			previewStyles.update({
				text_alignment: firstGroup.text_alignment,
				text_size: firstGroup.text_size,
				card_text_color: firstGroup.card_text_color,
				image_shape: firstGroup.image_shape,
				show_shadow: firstGroup.show_shadow,
				shadow_x: firstGroup.shadow_x,
				shadow_y: firstGroup.shadow_y,
				shadow_blur: firstGroup.shadow_blur,
				card_border_radius: firstGroup.card_border_radius,
				has_card_border: firstGroup.has_card_border,
				card_border_color: firstGroup.card_border_color,
				card_border_width: firstGroup.card_border_width,
				card_background_color: firstGroup.card_background_color,
				card_background_opacity: firstGroup.card_background_opacity,
				style: firstGroup.style
			});
		}
	}

	const themes: Record<string, Array<{ id: string; name: string; preset: string; isCustom?: boolean }>> = {
		classic: [
			{ id: 'mcalpine', name: 'McAlpine', preset: 'mcalpine' },
			{ id: 'yoga', name: 'Yoga', preset: 'yoga' },
			{ id: 'jerry', name: 'Jerry', preset: 'jerry' }
		],
		vibrant: [
			{ id: 'vibrant', name: 'Vibrant', preset: 'vibrant' }
		],
		cozy: [
			{ id: 'minimal', name: 'Minimal', preset: 'minimal' },
			{ id: 'default', name: 'Default', preset: 'default' }
		],
		bold: [
			{ id: 'dark', name: 'Dark Mode', preset: 'dark' }
		],
		custom: [
			{ id: 'custom', name: 'Your Custom Theme', preset: 'custom', isCustom: true }
		]
	};

	async function loadProfile() {
		const token = get(auth).token;
		const profileData = await profileApi.getMyProfile(token!);
		
		profile = profileData;
		bio = profileData?.bio || '';
		avatarPreview = profileData?.avatar_url || '';
		
		// Load social links
		if (profileData?.social_links) {
			try {
				socialLinks = typeof profileData.social_links === 'string' 
					? JSON.parse(profileData.social_links) 
					: profileData.social_links;
			} catch {
				socialLinks = [];
			}
		} else {
			socialLinks = [];
		}
		
		// Load page settings
		pageSettings.showShareButton = profileData?.show_share_button ?? true;
		pageSettings.showSubscribeButton = profileData?.show_subscribe_button ?? true;
		pageSettings.hideBranding = profileData?.hide_branding ?? false;
	}

	onMount(async () => {
		try {
			const token = get(auth).token;
			const [profileData, linksData, blocksData] = await Promise.all([
				profileApi.getMyProfile(token!).catch(() => null),
				linksApi.getLinks(token!).catch(() => []),
				blocksApi.getBlocks(token!).catch(() => [])
			]);
			
			profile = profileData;
			links = linksData || [];
			blocks = blocksData || [];
			bio = profileData?.bio || '';
			avatarPreview = profileData?.avatar_url || '';
			
			// Load social links
			if (profileData?.social_links) {
				try {
					socialLinks = typeof profileData.social_links === 'string' 
						? JSON.parse(profileData.social_links) 
						: profileData.social_links;
				} catch {
					socialLinks = [];
				}
			}
			
			// Load page settings
			pageSettings.showShareButton = profileData?.show_share_button ?? true;
			pageSettings.showSubscribeButton = profileData?.show_subscribe_button ?? true;
			pageSettings.hideBranding = profileData?.hide_branding ?? false;
			
			// Sync previewStyles from links data
			syncPreviewStyles(linksData || []);
			
			// Load theme from profile
			const savedThemeName = profileData?.theme_name || 'default';
			console.log('[LOAD] Loading theme:', {
				theme_name: savedThemeName,
				has_theme_config: !!profileData?.theme_config,
				profileData_keys: Object.keys(profileData || {})
			});
			
			if (savedThemeName && savedThemeName !== 'custom') {
				// Load preset theme
				console.log('[LOAD] Loading PRESET theme:', savedThemeName);
				const preset = themePresets[savedThemeName];
				console.log('[LOAD] Preset found:', !!preset, preset?.page);
				if (preset) {
					globalTheme.setPreset(savedThemeName);
					currentTheme = savedThemeName;
					
					// Find category
					for (const [categoryId, categoryThemes] of Object.entries(themes)) {
						const found = categoryThemes.find(t => t.id === savedThemeName);
						if (found) {
							selectedCategory = categoryId;
							break;
						}
					}
					console.log('[LOAD] Preset theme loaded:', { 
						currentTheme, 
						selectedCategory,
						themeConfig: globalTheme.getCurrent()
					});
				} else {
					console.error('[LOAD] Preset not found for:', savedThemeName);
				}
			} else if (savedThemeName === 'custom' && profileData?.theme_config) {
				// Load custom theme
				console.log('[LOAD] Loading CUSTOM theme');
				try {
					const themeStr = typeof profileData.theme_config === 'string' 
						? profileData.theme_config 
						: JSON.stringify(profileData.theme_config);
					
					if (themeStr && themeStr !== '{}' && themeStr !== 'null') {
						globalTheme.loadFromJSON(themeStr);
						currentTheme = 'custom';
						selectedCategory = 'custom';
						console.log('[LOAD] Custom theme loaded:', globalTheme.getCurrent());
					}
				} catch (themeError) {
					console.warn('[LOAD] Failed to load custom theme:', themeError);
				}
			}
			
			// Load header config from profile
			if (profileData?.header_config) {
				try {
					const headerConfig = typeof profileData.header_config === 'string' 
						? JSON.parse(profileData.header_config) 
						: profileData.header_config;
					
					if (headerConfig && Object.keys(headerConfig).length > 0) {
						currentHeaderStyle.set(headerConfig);
					}
				} catch (headerError) {
					console.warn('Failed to load header config:', headerError);
				}
			}
		} catch (e: any) {
			console.error('Failed to load data:', e);
		}
	});

	/**
	 * Apply theme preset or load custom theme
	 */
	async function selectTheme(themeId: string, presetName: string) {
		if (applying) return;
		
		console.log('🎨 selectTheme called:', { themeId, presetName });
		
		// If selecting custom, reload saved config from profile
		if (themeId === 'custom') {
			console.log('📦 Loading custom theme from DB...');
			currentTheme = 'custom';
			selectedCategory = 'custom';
			
			// Reload theme config from profile (DB)
			try {
				const token = get(auth).token;
				const profileData = await profileApi.getMyProfile(token!);
				
				if (profileData?.theme_config) {
					const themeStr = typeof profileData.theme_config === 'string' 
						? profileData.theme_config 
						: JSON.stringify(profileData.theme_config);
					
					console.log('✅ Custom theme loaded from DB:', JSON.parse(themeStr));
					globalTheme.loadFromJSON(themeStr);
					
					// Reload header config
					if (profileData?.header_config) {
						const headerConfig = typeof profileData.header_config === 'string' 
							? JSON.parse(profileData.header_config) 
							: profileData.header_config;
						currentHeaderStyle.set(headerConfig);
					}
					
					// Reload links to sync preview
					const freshLinks = await linksApi.getLinks(token!);
					links = freshLinks || [];
					syncPreviewStyles(freshLinks || []);
				}
			} catch (e: any) {
				console.error('❌ Failed to load custom theme:', e);
				toast.error('Failed to load custom theme');
			}
			
			pendingChanges.reset();
			return;
		}
		
		const preset = themePresets[presetName];
		if (!preset) {
			console.warn('⚠️ Preset not found:', presetName);
			return;
		}

		console.log('🎭 Applying preset:', presetName);

		// Reset previewStyles to clear any custom overrides
		previewStyles.reset();

		// Optimistic update - reset everything to theme preset
		const oldTheme = currentTheme;
		currentTheme = themeId;
		globalTheme.setPreset(presetName);
		console.log('✅ Theme set to:', themeId);
		
		// Update header style
		if (preset.header) {
			currentHeaderStyle.set(preset.header);
		}
		
		// Update groups only if they exist
		const hasLinkGroups = links?.some(l => l.is_group);
		const hasBlockGroups = blocks?.some(b => b.is_group);
		
		if (hasLinkGroups) {
			links = links.map(link => link.is_group ? applyCardStylesToGroup(link, preset.card) : link);
		}
		if (hasBlockGroups) {
			blocks = blocks.map(block => block.is_group ? applyTextStylesToGroup(block, preset.text) : block);
		}

		// Mark as having changes (will need to save)
		pendingChanges.updateTheme({ ...preset.page, ...preset.card });
		pendingChanges.updateHeader(preset.header);
		
		const cardStyles = cardStylesToLinkFields(preset.card, preset.text);
		pendingChanges.updateLinkStyles(cardStyles);
		
		console.log('⚠️ Theme changed - need to Save All to persist');
		toast.info('Theme preview applied. Click "Save All" to keep changes.');
	}

	function handleAvatarChange(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (file) {
			avatarFile = file;
			avatarPreview = URL.createObjectURL(file);
		}
	}

	async function uploadAvatar() {
		if (!avatarFile) return;
		
		uploadingAvatar = true;
		try {
			const token = get(auth).token;
			const updated = await profileApi.uploadAvatar(avatarFile, token!);
			profile = updated;
			avatarPreview = updated.avatar_url || '';
			avatarFile = null;
			toast.success('Avatar updated!');
		} catch (e: any) {
			toast.error(e.message || 'Failed to upload avatar');
		} finally {
			uploadingAvatar = false;
		}
	}

	async function saveProfile() {
		savingProfile = true;
		try {
			const token = get(auth).token;
			const updated = await profileApi.updateProfile({ bio }, token!);
			profile = updated;
			toast.success('Profile updated!');
		} catch (e: any) {
			toast.error(e.message || 'Failed to update profile');
		} finally {
			savingProfile = false;
		}
	}

	// Get preview colors from preset (memoized)
	function getPresetColors(presetName: string) {
		if (presetColorsCache.has(presetName)) {
			return presetColorsCache.get(presetName)!;
		}
		
		const preset = themePresets[presetName];
		const colors = !preset?.page || !preset?.card
			? { bg: '#ffffff', accent: '#000000', cardBg: '#ffffff', cardText: '#000000', textColor: '#000000' }
			: {
				bg: preset.page.pageBackground,
				accent: preset.page.accentColor,
				cardBg: preset.card.cardBackground,
				cardText: preset.card.cardTextColor,
				textColor: preset.page.textColor
			};
		
		presetColorsCache.set(presetName, colors);
		return colors;
	}
</script>

<div class="h-full bg-gray-50">
	<!-- Page Header -->
	<div class="bg-white/80 backdrop-blur-xl border-b border-gray-200/50 px-8 h-16 sticky top-0 z-50">
		<div class="flex items-center justify-between h-full">
			<div>
				<h1 class="text-xl font-bold text-gray-900">Appearance</h1>
				<p class="text-xs text-gray-500">Customize your profile appearance</p>
			</div>
			<div class="flex items-center gap-3">
				{#if hasUnsavedChanges}
					<div class="flex items-center gap-2 px-3 py-1.5 bg-amber-50 border border-amber-200 rounded-lg">
						<div class="w-2 h-2 bg-amber-500 rounded-full animate-pulse"></div>
						<span class="text-xs font-medium text-amber-700">Unsaved changes</span>
					</div>
				{/if}
				<button
					onclick={saveAllChanges}
					disabled={saving || !hasUnsavedChanges}
					class="px-5 py-2 bg-indigo-600 text-white rounded-lg text-sm font-semibold hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center gap-2 shadow-sm hover:shadow-md"
				>
					{#if saving}
						<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Saving...
					{:else}
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
						</svg>
						Save All
					{/if}
				</button>
				<button
					onclick={() => {
						if (profile?.username) {
							window.open(`/${profile.username}`, '_blank');
						}
					}}
					class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-200 transition-colors"
				>
					View Live
				</button>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="p-8">
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<!-- Settings Section (2/3) -->
			<div class="lg:col-span-2 space-y-6">
				<!-- Profile Settings Card -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6 hover:shadow-lg transition-shadow">
					<div class="flex items-start justify-between">
						<div class="flex items-start gap-4 flex-1">
							<!-- Avatar -->
							<div class="relative">
								{#if avatarPreview}
									<img src={avatarPreview} alt="Avatar" class="w-16 h-16 rounded-xl object-cover ring-2 ring-gray-100" />
								{:else}
									<div class="w-16 h-16 rounded-xl bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 ring-2 ring-gray-100"></div>
								{/if}
							</div>
							
							<!-- Info -->
							<div class="flex-1 min-w-0">
								<h3 class="text-base font-semibold text-gray-900">Profile Information</h3>
								<p class="text-sm text-gray-500 mt-0.5">@{profile?.username || 'username'}</p>
								{#if bio}
									<p class="text-sm text-gray-600 mt-2 line-clamp-2">{bio}</p>
								{:else}
									<p class="text-sm text-gray-400 mt-2 italic">No bio yet</p>
								{/if}
							</div>
						</div>
						
						<!-- Edit Button -->
						<button
							onclick={() => showProfileModal = true}
							class="flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
							</svg>
							Edit
						</button>
					</div>
				</div>

				<!-- Theme Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-xl font-bold text-gray-900">Theme</h2>
						<span class="text-sm text-indigo-600 flex items-center gap-2 transition-opacity duration-200" style="opacity: {applying ? 1 : 0}; pointer-events: {applying ? 'auto' : 'none'};">
							<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Applying...
						</span>
					</div>
					
					<div class="flex gap-2 mb-4">
						{#each categories as cat}
							<button
								onclick={() => selectedCategory = cat.id}
								class="px-4 py-2 rounded-full text-xs font-medium transition-all {selectedCategory === cat.id
									? cat.id === 'custom' 
										? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-md'
										: 'bg-gray-900 text-white'
									: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
							>
								{#if cat.id === 'custom'}
									<span class="flex items-center gap-1.5">
										<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"/>
										</svg>
										{cat.label}
									</span>
								{:else}
									{cat.label}
								{/if}
							</button>
						{/each}
					</div>

					<!-- Theme Grid -->
					<div class="grid grid-cols-3 md:grid-cols-5 lg:grid-cols-7 gap-2">
						{#each (selectedCategory === 'all' ? Object.values(themes).flat() : themes[selectedCategory] || []) as themeItem (themeItem.id)}
							{@const colors = getPresetColors(themeItem.preset)}
							{@const isSelected = currentTheme === themeItem.id}
							<button
								onclick={() => selectTheme(themeItem.id, themeItem.preset)}
								class="relative group cursor-pointer"
							>
								<!-- Theme Preview Card - Phone Shape -->
								<div
									class="aspect-[9/16] rounded-lg overflow-hidden border-2 transition-all relative {isSelected
										? 'border-indigo-600 shadow-md ring-2 ring-indigo-200'
										: themeItem.isCustom
										? 'border-dashed border-gray-300'
										: 'border-gray-200 hover:border-gray-300 hover:shadow-sm'}"
									style="background: {themeItem.isCustom ? 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)' : colors.bg};"
								>
									{#if themeItem.isCustom}
										<!-- Custom Icon -->
										<div class="h-full flex items-center justify-center">
											<svg class="w-6 h-6 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"/>
											</svg>
										</div>
									{:else}
										<!-- Minimal Preview -->
										<div class="h-full flex flex-col p-1.5 gap-1">
											<div class="w-full h-2 rounded-sm" style="background: {colors.cardBg}; opacity: 0.8;"></div>
											<div class="w-full h-2 rounded-sm" style="background: {colors.cardBg}; opacity: 0.6;"></div>
											<div class="flex-1"></div>
											<div class="w-3 h-3 rounded-full mx-auto" style="background: {colors.accent}; opacity: 0.5;"></div>
										</div>
									{/if}

									<!-- Selected Indicator -->
									{#if isSelected}
										<div class="absolute inset-0 flex items-center justify-center bg-indigo-600/10">
											<div class="w-5 h-5 bg-indigo-600 rounded-full flex items-center justify-center shadow-lg">
												<svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
												</svg>
											</div>
										</div>
									{/if}
								</div>

								<!-- Theme Name -->
								<p class="mt-1 text-[10px] font-medium text-center truncate {themeItem.isCustom ? 'text-indigo-600' : 'text-gray-600'}">{themeItem.name}</p>
							</button>
						{/each}
					</div>

					<p class="mt-3 text-xs text-gray-500">
						Selecting a theme will apply styles to your page background and all link/text groups.
					</p>
				</div>

				<!-- Background Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<BackgroundSelector />
				</div>

				<!-- Header Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<HeaderStyleEditor />
				</div>

				<!-- Card Colors Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<CardStyleEditor />
				</div>

				<!-- Typography Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<TypographyEditor {links} />
				</div>

				<!-- Image Shape Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<ImageStyleEditor {links} />
				</div>

				<!-- Shadow Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<ShadowEditor {links} />
				</div>

				<!-- Border & Radius Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<BorderEditor {links} />
				</div>

				<!-- Spacing Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<SpacingEditor {links} />
				</div>

				<!-- Social Links Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<SocialLinksEditor {socialLinks} onUpdate={loadProfile} />
				</div>

				<!-- Page Settings Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<PageSettingsEditor settings={pageSettings} />
				</div>
			</div>

			<!-- Preview Panel (1/3) -->
			<div class="lg:col-span-1">
				<div class="sticky top-24">
					<ProfilePreview 
						{profile} 
						{links} 
						{blocks} 
						{socialLinks} 
						showShareButton={pageSettings.showShareButton}
						showSubscribeButton={pageSettings.showSubscribeButton}
						hideBranding={pageSettings.hideBranding}
						showInactive={true} 
					/>
				</div>
			</div>
		</div>
	</div>
</div>


<!-- Profile Edit Modal -->
{#if showProfileModal}
	<div class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4" onclick={(e) => e.target === e.currentTarget && (showProfileModal = false)}>
		<div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
			<!-- Modal Header -->
			<div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
				<h2 class="text-lg font-bold text-gray-900">Edit Profile</h2>
				<button
					onclick={() => showProfileModal = false}
					class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 transition-colors"
				>
					<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>
			
			<!-- Modal Content -->
			<div class="p-6 space-y-8">
				<!-- Avatar Section -->
				<div>
					<label class="block text-sm font-semibold text-gray-900 mb-4">Profile Picture</label>
					
					<!-- Modern Avatar Upload Area -->
					<div class="relative group">
						<!-- Avatar Display with Hover Effect -->
						<div class="flex items-center justify-center mb-4">
							<div class="relative">
								{#if avatarPreview}
									<img 
										src={avatarPreview} 
										alt="Avatar" 
										class="w-32 h-32 rounded-2xl object-cover ring-4 ring-indigo-100 shadow-lg transition-all group-hover:ring-indigo-200 group-hover:shadow-xl" 
									/>
								{:else}
									<div class="w-32 h-32 rounded-2xl bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 ring-4 ring-indigo-100 shadow-lg transition-all group-hover:ring-indigo-200 group-hover:shadow-xl"></div>
								{/if}
								
								<!-- Camera Icon Overlay -->
								<div class="absolute inset-0 flex items-center justify-center bg-black/40 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity">
									<svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
									</svg>
								</div>
							</div>
						</div>
						
						<!-- Upload Controls -->
						<input
							type="file"
							accept="image/*"
							onchange={handleAvatarChange}
							class="hidden"
							id="modal-avatar-upload"
						/>
						
						<div class="space-y-3">
							<!-- Choose File Button -->
							<label
								for="modal-avatar-upload"
								class="w-full inline-flex items-center justify-center gap-2 px-5 py-3 bg-gradient-to-r from-indigo-50 to-purple-50 border-2 border-indigo-200 text-indigo-700 rounded-xl text-sm font-semibold hover:from-indigo-100 hover:to-purple-100 hover:border-indigo-300 cursor-pointer transition-all hover:shadow-md"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
								</svg>
								Choose Image
							</label>
							
							<!-- Upload Button (shown when file selected) -->
							{#if avatarFile}
								<button
									onclick={uploadAvatar}
									disabled={uploadingAvatar}
									class="w-full inline-flex items-center justify-center gap-2 px-5 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-xl text-sm font-semibold hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-xl"
								>
									{#if uploadingAvatar}
										<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
										Uploading...
									{:else}
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
										</svg>
										Upload Image
									{/if}
								</button>
							{/if}
							
							<!-- File Info -->
							<div class="flex items-center justify-center gap-2 text-xs text-gray-500">
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
								</svg>
								<span>JPG, PNG or GIF • Max 5MB</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Divider -->
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-gray-200"></div>
					</div>
					<div class="relative flex justify-center">
						<span class="px-3 bg-white text-xs text-gray-500 font-medium">ABOUT YOU</span>
					</div>
				</div>

				<!-- Bio Section -->
				<div>
					<label class="block text-sm font-semibold text-gray-900 mb-3">Bio</label>
					<div class="relative">
						<textarea
							bind:value={bio}
							placeholder="Tell people about yourself..."
							rows="5"
							maxlength="200"
							class="w-full px-4 py-3 bg-gradient-to-br from-gray-50 to-gray-100/50 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-all resize-none text-sm placeholder:text-gray-400"
						></textarea>
						
						<!-- Character Count Badge -->
						<div class="absolute bottom-3 right-3">
							<span class="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-medium transition-colors {bio.length >= 180 ? 'bg-amber-100 text-amber-700' : bio.length >= 200 ? 'bg-red-100 text-red-700' : 'bg-gray-100 text-gray-600'}">
								{bio.length}/200
							</span>
						</div>
					</div>
					<p class="text-xs text-gray-500 mt-2 flex items-center gap-1.5">
						<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
						</svg>
						Share your story, interests, or what makes you unique
					</p>
				</div>
			</div>
			
			<!-- Modal Footer -->
			<div class="sticky bottom-0 bg-gray-50 border-t border-gray-200 px-6 py-4 flex items-center justify-end gap-3">
				<button
					onclick={() => showProfileModal = false}
					class="px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
				>
					Cancel
				</button>
				<button
					onclick={async () => {
						await saveProfile();
						showProfileModal = false;
					}}
					disabled={savingProfile}
					class="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
				>
					{#if savingProfile}
						<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Saving...
					{:else}
						Save Changes
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}
