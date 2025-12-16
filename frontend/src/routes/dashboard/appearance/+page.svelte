<script lang="ts">
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { 
		globalTheme, 
		themePresets, 
		themeStyles,
		currentHeaderStyle,
		defaultHeaderStyles
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
	import { loadThemeFromProfile, syncPreviewStylesFromTheme } from '$lib/utils/themeLoader';
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
	let savedThemeName = $state<string>('default');
	let saving = $state(false);
	let loading = $state(true);
	let isInitialLoad = $state(true);
	let customHeaderPresets = $state<any[]>([]);
	let isUsingPreset = $state(true); // Track if user is using unmodified preset
	
	const hasUnsavedChanges = $derived($pendingChanges.hasChanges);
	
	// Sync preview styles whenever theme changes
	$effect(() => {
		if (!loading && !isInitialLoad) {
			syncPreviewStylesFromTheme();
		}
		const _ = $globalTheme;
	});
	
	// Auto-switch to custom theme when user modifies preset
	$effect(() => {
		if (loading || isInitialLoad || currentTheme === 'custom' || isUsingPreset) return;
		
		// User has modified theme → switch to custom
		currentTheme = 'custom';
		selectedCategory = 'custom';
		pendingChanges.updateTheme($globalTheme);
		pendingChanges.updateHeader($currentHeaderStyle);
		toast.info('Switched to custom theme');
	});
	

	
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



	async function handleDeletePreset(event: CustomEvent<{ presets: any[] }>) {
		const updatedPresets = event.detail.presets;
		customHeaderPresets = updatedPresets;
		
		// Save immediately to database
		try {
			const token = get(auth).token!;
			const updatedTheme = globalTheme.getCurrent();
			const updatedHeader = get(currentHeaderStyle);
			
			const customThemeWithPresets = {
				...updatedTheme,
				header: updatedHeader,
				customHeaderPresets: updatedPresets
			};
			
			if (currentTheme === 'custom') {
				await profileApi.updateProfile({ 
					theme_name: 'custom',
					custom_theme_config: JSON.stringify(customThemeWithPresets)
				}, token);
			} else {
				await profileApi.updateProfile({
					theme_name: currentTheme,
					theme_config: null,
					header_config: JSON.stringify(updatedHeader),
					custom_theme_config: JSON.stringify({ customHeaderPresets: updatedPresets })
				}, token);
			}
			
			toast.success('Custom header style deleted successfully!');
		} catch (e: any) {
			toast.error(e.message || 'Failed to delete header style');
			// Reload to restore state
			window.location.reload();
		}
	}

	async function saveAllChanges() {
		if (saving || !hasUnsavedChanges) return;
		
		saving = true;
		const minDelay = new Promise(resolve => setTimeout(resolve, 500));
		
		try {
			const token = get(auth).token!;
			const updatedTheme = globalTheme.getCurrent();
			const updatedHeader = get(currentHeaderStyle);
			const changes = pendingChanges.getAll();
			
			const savePromises = [];
			
			// IMPORTANT: Snapshot $state proxies before stringify to avoid losing data
			const customHeaderPresetsSnapshot = JSON.parse(JSON.stringify(customHeaderPresets));
			
			const customThemeWithPresets = {
				...updatedTheme,
				header: updatedHeader,
				customHeaderPresets: customHeaderPresetsSnapshot
			};
			
			if (currentTheme === 'custom') {
				savePromises.push(
					profileApi.updateProfile({ 
						theme_name: 'custom',
						custom_theme_config: JSON.stringify(customThemeWithPresets)
					}, token)
				);
				savedThemeName = 'custom';
			} else {
				savePromises.push(
					profileApi.updateProfile({
						theme_name: currentTheme,
						theme_config: null,
						header_config: JSON.stringify(updatedHeader),
						custom_theme_config: JSON.stringify({ customHeaderPresets: customHeaderPresetsSnapshot })
					}, token)
				);
				savedThemeName = currentTheme;
			}
			
			if (changes.linkStyles && Object.keys(changes.linkStyles).length > 0) {
				savePromises.push(linksApi.updateAllGroupStyles(changes.linkStyles, token));
			}
			
			await Promise.all([...savePromises, minDelay]);
			
			links = await linksApi.getLinks(token);
			previewStyles.reset();
			syncPreviewStylesFromTheme();
			pendingChanges.reset();
			
			// Update original values after successful save
			pendingChanges.setOriginal({
				theme: globalTheme.getCurrent(),
				header: get(currentHeaderStyle)
			});
			
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

	// Setup modification tracking
	let trackModifications = $state(false);
	function handleThemeModified() {
		if (trackModifications && isUsingPreset && currentTheme !== 'custom') {
			isUsingPreset = false;
		}
	}
	
	onMount(async () => {
		// Set callbacks to track user modifications
		globalTheme.setModifiedCallback(handleThemeModified);
		currentHeaderStyle.setModifiedCallback(handleThemeModified);
		
		try {
			const token = get(auth).token;
			
			const profileData = await profileApi.getMyProfile(token!).catch(() => null);
			
			const { themeName, category } = loadThemeFromProfile(profileData);
			savedThemeName = themeName;
			currentTheme = themeName;
			selectedCategory = category;
			isUsingPreset = themeName !== 'custom'; // Set flag based on loaded theme
			
			if (profileData?.custom_theme_config) {
				try {
					const customConfig = typeof profileData.custom_theme_config === 'string'
						? JSON.parse(profileData.custom_theme_config)
						: profileData.custom_theme_config;
					
					if (customConfig?.customHeaderPresets && Array.isArray(customConfig.customHeaderPresets)) {
						customHeaderPresets = customConfig.customHeaderPresets;
					}
				} catch (e) {
					console.warn('Failed to load custom header presets:', e);
				}
			}
			
			// Load other data after theme is set
			const [linksData, blocksData] = await Promise.all([
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
			
			syncPreviewStylesFromTheme();
			
			// Set original values for change detection
			pendingChanges.setOriginal({
				theme: globalTheme.getCurrent(),
				header: get(currentHeaderStyle)
			});
			
		} catch (e: any) {
			console.error('Failed to load data:', e);
		} finally {
			loading = false;
			// Allow auto-switch after initial load completes (short delay to avoid false triggers)
			setTimeout(() => {
				isInitialLoad = false;
				trackModifications = true; // Start tracking modifications
			}, 200);
		}
		
		return () => {
			// Cleanup
		};
	});

	async function selectTheme(themeId: string, presetName: string) {
		const isAlreadySaved = savedThemeName === themeId;
		
		if (themeId === 'custom') {
			// If already on custom theme, don't reload from database (keep current modifications)
			if (currentTheme === 'custom') {
				toast.info('Already using custom theme');
				return;
			}
			
			currentTheme = 'custom';
			selectedCategory = 'custom';
			isUsingPreset = false; // Custom theme is not a preset
			
			try {
				const token = get(auth).token;
				const profileData = await profileApi.getMyProfile(token!);
				
				// Temporarily disable tracking while loading custom theme
				trackModifications = false;
				
				if (profileData?.custom_theme_config) {
					const customConfig = typeof profileData.custom_theme_config === 'string' 
						? JSON.parse(profileData.custom_theme_config)
						: profileData.custom_theme_config;
					
					if (customConfig && Object.keys(customConfig).length > 0) {
						const { header, ...themeConfig } = customConfig;
						globalTheme.loadFromJSON(JSON.stringify(themeConfig));
						currentHeaderStyle.set(header && Object.keys(header).length > 0 ? header : defaultHeaderStyles);
					} else {
						globalTheme.setPreset('default');
						currentHeaderStyle.set(defaultHeaderStyles);
					}
				} else {
					globalTheme.setPreset('default');
					currentHeaderStyle.set(defaultHeaderStyles);
				}
				
				syncPreviewStylesFromTheme();
				
				// Re-enable tracking
				setTimeout(() => {
					trackModifications = true;
				}, 100);
				
				if (!isAlreadySaved) {
					pendingChanges.updateTheme(globalTheme.getCurrent());
					pendingChanges.updateHeader(get(currentHeaderStyle));
					toast.info('Custom theme loaded. Click "Save All" to apply.');
				} else {
					pendingChanges.reset();
					pendingChanges.setOriginal({
						theme: globalTheme.getCurrent(),
						header: get(currentHeaderStyle)
					});
					toast.info('Custom theme loaded');
				}
			} catch (e: any) {
				console.error('Failed to load custom theme:', e);
				toast.error('Failed to load custom theme');
			}
			return;
		}
		
		const preset = themePresets[presetName];
		if (!preset) return;

		currentTheme = themeId;
		isUsingPreset = true; // Mark as using unmodified preset
		
		// Temporarily disable tracking while setting preset
		trackModifications = false;
		globalTheme.setPreset(presetName);
		if (preset.header) {
			currentHeaderStyle.set(preset.header);
		}
		syncPreviewStylesFromTheme();
		
		// Re-enable tracking after a short delay
		setTimeout(() => {
			trackModifications = true;
		}, 100);

		if (isAlreadySaved) {
			pendingChanges.reset();
			pendingChanges.setOriginal({
				theme: globalTheme.getCurrent(),
				header: get(currentHeaderStyle)
			});
			toast.info('Theme loaded');
		} else {
			pendingChanges.updateTheme(globalTheme.getCurrent());
			pendingChanges.updateHeader(preset.header || get(currentHeaderStyle));
			toast.info('Theme preview applied. Click "Save All" to keep changes.');
		}
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
			// Filter out empty social links
			const validSocialLinks = socialLinks.filter(link => link.platform && link.url);
			const updated = await profileApi.updateProfile({ 
				bio,
				social_links: JSON.stringify(validSocialLinks)
			}, token!);
			profile = updated;
			socialLinks = validSocialLinks;
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
	{#if loading}
		<!-- Loading State -->
		<div class="h-full flex items-center justify-center">
			<div class="text-center">
				<svg class="w-12 h-12 animate-spin text-indigo-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				<p class="text-sm text-gray-600">Loading appearance settings...</p>
			</div>
		</div>
	{:else}
	<div in:fade={{ duration: 200 }}>
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
				<!-- Profile & Social Card -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6 hover:shadow-lg transition-shadow">
					<div class="flex items-start justify-between mb-4">
						<h3 class="text-base font-semibold text-gray-900">Profile & Social</h3>
						<button
							onclick={() => showProfileModal = true}
							class="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
							</svg>
							Edit
						</button>
					</div>
					
					<div class="flex items-start gap-4">
						<!-- Avatar -->
						<div class="relative flex-shrink-0">
							{#if avatarPreview}
								<img src={avatarPreview} alt="Avatar" class="w-20 h-20 rounded-2xl object-cover ring-2 ring-gray-100 shadow-sm" />
							{:else}
								<div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 ring-2 ring-gray-100 shadow-sm"></div>
							{/if}
						</div>
						
						<!-- Info -->
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium text-gray-500">@{profile?.username || 'username'}</p>
							{#if bio}
								<p class="text-sm text-gray-700 mt-1.5 line-clamp-2 leading-relaxed">{bio}</p>
							{:else}
								<p class="text-sm text-gray-400 mt-1.5 italic">No bio yet</p>
							{/if}
							
							<!-- Social Icons -->
							{#if socialLinks.length > 0}
								<div class="flex items-center gap-2 mt-3 flex-wrap">
									{#each socialLinks as link}
										{@const platform = [
											{ id: 'twitter', icon: 'M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z' },
											{ id: 'facebook', icon: 'M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z' },
											{ id: 'instagram', icon: 'M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z' },
											{ id: 'linkedin', icon: 'M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z' },
											{ id: 'github', icon: 'M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z' },
											{ id: 'youtube', icon: 'M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z' },
											{ id: 'tiktok', icon: 'M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z' },
											{ id: 'telegram', icon: 'M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z' },
											{ id: 'discord', icon: 'M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0 a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z' },
											{ id: 'website', icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z' }
										].find(p => p.id === link.platform)}
										<a 
											href={link.url} 
											target="_blank" 
											rel="noopener noreferrer"
											class="w-8 h-8 flex items-center justify-center rounded-lg bg-gray-100 hover:bg-gray-200 text-gray-700 hover:text-gray-900 transition-all hover:scale-110"
											title={link.platform}
										>
											<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
												<path d={platform?.icon || ''}/>
											</svg>
										</a>
									{/each}
								</div>
							{:else}
								<p class="text-xs text-gray-400 mt-3 italic">No social links added</p>
							{/if}
						</div>
					</div>
				</div>

				<!-- Theme Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-xl font-bold text-gray-900">Theme</h2>
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
					<HeaderStyleEditor 
						bind:customPresets={customHeaderPresets}
						currentThemeName={currentTheme}
						on:deletePreset={handleDeletePreset}
					/>
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

				<!-- Page Settings Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<PageSettingsEditor settings={pageSettings} />
				</div>
			</div>

			<!-- Preview Panel (1/3) -->
			<div class="lg:col-span-1">
				<div class="sticky top-24 space-y-4">
					<!-- Preview Header -->
					<div class="text-center">
						<h2 class="text-2xl font-bold text-gray-900 mb-2">Live Preview</h2>
						<p class="text-sm text-gray-500">See how your profile looks</p>
					</div>
					
					<!-- Preview Mockup -->
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
{/if}
</div>


<!-- Profile Edit Modal -->
{#if showProfileModal}
	<div class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4" onclick={(e) => e.target === e.currentTarget && (showProfileModal = false)}>
		<div class="bg-white rounded-3xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-hidden flex flex-col">
			<!-- Modal Header -->
			<div class="flex-shrink-0 bg-gradient-to-r from-indigo-600 to-purple-600 px-8 py-6 flex items-center justify-between">
				<div>
					<h2 class="text-2xl font-bold text-white">Edit Profile</h2>
					<p class="text-indigo-100 text-sm mt-1">Update your profile information and social links</p>
				</div>
				<button
					onclick={() => showProfileModal = false}
					class="w-10 h-10 flex items-center justify-center rounded-xl bg-white/10 hover:bg-white/20 text-white transition-colors"
				>
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>
			
			<!-- Modal Content -->
			<div class="flex-1 overflow-y-auto">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-8 p-8">
					<!-- Left Column: Avatar & Bio -->
					<div class="space-y-6">
						<!-- Avatar Section -->
						<div>
							<label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-4">Profile Picture</label>
							
							<div class="relative group">
								<div class="flex items-center justify-center mb-4">
									<div class="relative">
										{#if avatarPreview}
											<img 
												src={avatarPreview} 
												alt="Avatar" 
												class="w-40 h-40 rounded-3xl object-cover shadow-xl transition-all group-hover:shadow-2xl" 
											/>
										{:else}
											<div class="w-40 h-40 rounded-3xl bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 shadow-xl transition-all group-hover:shadow-2xl"></div>
										{/if}
										
										<div class="absolute inset-0 flex items-center justify-center bg-black/50 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
											<svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
											</svg>
										</div>
									</div>
								</div>
								
								<input
									type="file"
									accept="image/*"
									onchange={handleAvatarChange}
									class="hidden"
									id="modal-avatar-upload"
								/>
								
								<div class="space-y-2">
									<label
										for="modal-avatar-upload"
										class="w-full inline-flex items-center justify-center gap-2 px-4 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-xl text-sm font-medium cursor-pointer transition-all"
									>
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
										</svg>
										Choose Image
									</label>
									
									{#if avatarFile}
										<button
											onclick={uploadAvatar}
											disabled={uploadingAvatar}
											class="w-full inline-flex items-center justify-center gap-2 px-4 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-all"
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
												Upload
											{/if}
										</button>
									{/if}
									
									<p class="text-xs text-gray-500 text-center">JPG, PNG or GIF • Max 5MB</p>
								</div>
							</div>
						</div>

						<!-- Bio Section -->
						<div>
							<label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-3">Bio</label>
							<div class="relative">
								<textarea
									bind:value={bio}
									placeholder="Tell people about yourself..."
									rows="6"
									maxlength="200"
									class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-all resize-none text-sm"
								></textarea>
								
								<div class="absolute bottom-3 right-3">
									<span class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium {bio.length >= 180 ? 'bg-amber-100 text-amber-700' : bio.length >= 200 ? 'bg-red-100 text-red-700' : 'bg-gray-200 text-gray-600'}">
										{bio.length}/200
									</span>
								</div>
							</div>
						</div>
					</div>

					<!-- Right Column: Social Links -->

					<div class="space-y-6">
						<div>
							<label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-4">Social Links</label>
							
							<div class="flex items-center justify-between mb-4">
								<p class="text-sm text-gray-600">Connect your social profiles</p>
								<button
									onclick={() => {
										socialLinks = [...socialLinks, { platform: '', url: '' }];
									}}
									class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors"
								>
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
									</svg>
									Add
								</button>
							</div>

							{#if socialLinks.length > 0}
								<div class="space-y-3 max-h-[400px] overflow-y-auto pr-2">
									{#each socialLinks as link, index}
										<div class="flex gap-2 items-start p-3 bg-gray-50 rounded-xl border border-gray-200">
											<div class="flex-1 space-y-2">
												<select 
													bind:value={link.platform}
													class="w-full px-3 py-2.5 text-sm bg-white border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
												>
													<option value="">Select platform...</option>
													<option value="twitter">Twitter/X</option>
													<option value="facebook">Facebook</option>
													<option value="instagram">Instagram</option>
													<option value="linkedin">LinkedIn</option>
													<option value="github">GitHub</option>
													<option value="youtube">YouTube</option>
													<option value="tiktok">TikTok</option>
													<option value="telegram">Telegram</option>
													<option value="discord">Discord</option>
													<option value="website">Website</option>
												</select>
												<input
													type="url"
													bind:value={link.url}
													placeholder="https://..."
													class="w-full px-3 py-2.5 text-sm bg-white border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
												/>
											</div>
											<button
												onclick={() => {
													socialLinks = socialLinks.filter((_, i) => i !== index);
												}}
												class="p-2.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
											>
												<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
												</svg>
											</button>
										</div>
									{/each}
								</div>
							{:else}
								<div class="text-center py-12 px-4 bg-gray-50 rounded-xl border border-gray-200">
									<svg class="w-16 h-16 mx-auto text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
									</svg>
									<p class="text-sm text-gray-500 font-medium">No social links</p>
									<p class="text-xs text-gray-400 mt-1">Click "Add" to get started</p>
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
			
			<!-- Modal Footer -->
			<div class="flex-shrink-0 bg-gray-50 border-t border-gray-200 px-8 py-5 flex items-center justify-between">
				<p class="text-sm text-gray-500">Changes will be saved to your profile</p>
				<div class="flex items-center gap-3">
					<button
						onclick={() => showProfileModal = false}
						class="px-5 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-200 bg-gray-100 rounded-xl transition-colors"
					>
						Cancel
					</button>
					<button
						onclick={async () => {
							await saveProfile();
							showProfileModal = false;
						}}
						disabled={savingProfile}
						class="inline-flex items-center gap-2 px-6 py-2.5 bg-indigo-600 text-white rounded-xl text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors shadow-lg shadow-indigo-500/30"
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
	</div>
{/if}

