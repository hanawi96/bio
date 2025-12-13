<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { 
		globalTheme, 
		themePresets, 
		themeStyles,
		cardStylesToLinkFields,
		textStylesToBlockStyle,
		applyCardStylesToGroup,
		applyTextStylesToGroup
	} from '$lib/stores/theme';
	import type { CardStyles, TextStyles, ThemePreset } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import { blocksApi } from '$lib/api/blocks';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';
	import type { Block } from '$lib/api/blocks';
	import { toast } from 'svelte-sonner';
	import ProfilePreview from '$lib/components/dashboard/preview/ProfilePreview.svelte';

	let profile = $state<Profile | null>(null);
	let links = $state<Link[]>([]);
	let blocks = $state<Block[]>([]);
	let selectedCategory = $state<string>('classic');
	let currentTheme = $state<string>('default');
	let saving = $state(false);
	let applying = $state(false);
	
	// Profile edit
	let bio = $state('');
	let avatarFile = $state<File | null>(null);
	let avatarPreview = $state<string>('');
	let uploadingAvatar = $state(false);
	let savingProfile = $state(false);

	const categories = [
		{ id: 'classic', label: 'Classic' },
		{ id: 'vibrant', label: 'Vibrant' },
		{ id: 'cozy', label: 'Cozy' },
		{ id: 'bold', label: 'Bold' }
	];

	const themes: Record<string, Array<{ id: string; name: string; preset: string }>> = {
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
		]
	};

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
			
			// Load theme from profile
			if (profileData?.theme_config) {
				try {
					const themeStr = typeof profileData.theme_config === 'string' 
						? profileData.theme_config 
						: JSON.stringify(profileData.theme_config);
					
					if (themeStr && themeStr !== '{}' && themeStr !== 'null') {
						globalTheme.loadFromJSON(themeStr);
					}
				} catch (themeError) {
					console.warn('Failed to load theme:', themeError);
				}
			}
		} catch (e: any) {
			console.error('Failed to load data:', e);
		}
	});

	/**
	 * Apply theme preset - Optimistic update + API call
	 */
	async function selectTheme(themeId: string, presetName: string) {
		const preset = themePresets[presetName];
		if (!preset) return;

		// Save old state for rollback
		const oldTheme = currentTheme;
		const oldLinks = links ? links.slice() : [];
		const oldBlocks = blocks ? blocks.slice() : [];
		const oldThemeConfig = globalTheme.getCurrent();

		// Optimistic update
		currentTheme = themeId;
		
		// 1. Update global theme store (page + card styles)
		globalTheme.setPreset(presetName);
		
		// 2. Update all link groups locally
		if (links) {
			links = links.map(link => {
				if (link.is_group) {
					return applyCardStylesToGroup(link, preset.card);
				}
				return link;
			});
		}
		
		// 3. Update all text groups locally
		if (blocks) {
			blocks = blocks.map(block => {
				if (block.is_group) {
					return applyTextStylesToGroup(block, preset.text);
				}
				return block;
			});
		}

		// API call in background
		applying = true;
		try {
			const token = get(auth).token;
			const result = await profileApi.applyTheme({
				theme_config: { ...preset.page, ...preset.card },
				card_styles: cardStylesToLinkFields(preset.card),
				text_styles: textStylesToBlockStyle(preset.text)
			}, token!);

			// Update with server response
			profile = result.profile;
			links = result.links || [];
			blocks = result.blocks || [];
			
			toast.success('Theme applied!');
		} catch (e: any) {
			// Rollback on error
			currentTheme = oldTheme;
			links = oldLinks;
			blocks = oldBlocks;
			globalTheme.loadFromJSON(JSON.stringify(oldThemeConfig));
			
			toast.error(e.message || 'Failed to apply theme');
		} finally {
			applying = false;
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
			const updated = await profileApi.updateProfile({ bio }, token!);
			profile = updated;
			toast.success('Profile updated!');
		} catch (e: any) {
			toast.error(e.message || 'Failed to update profile');
		} finally {
			savingProfile = false;
		}
	}

	// Get preview colors from preset
	function getPresetColors(presetName: string) {
		const preset = themePresets[presetName];
		if (!preset?.page || !preset?.card) {
			return { bg: '#ffffff', accent: '#000000', cardBg: '#ffffff', cardText: '#000000', textColor: '#000000' };
		}
		return {
			bg: preset.page.pageBackground,
			accent: preset.page.accentColor,
			cardBg: preset.card.cardBackground,
			cardText: preset.card.cardTextColor,
			textColor: preset.page.textColor
		};
	}
</script>

<div class="h-full bg-gray-50">
	<!-- Page Header -->
	<div class="bg-white/80 backdrop-blur-xl border-b border-gray-200/50 px-8 h-16 sticky top-0 z-10">
		<div class="flex items-center justify-between h-full">
			<div>
				<h1 class="text-xl font-bold text-gray-900">Appearance</h1>
				<p class="text-xs text-gray-500">Customize your profile appearance</p>
			</div>
			<button
				onclick={() => {
					if (profile?.username) {
						window.open(`/${profile.username}`, '_blank');
					}
				}}
				class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors"
			>
				View Live Preview
			</button>
		</div>
	</div>

	<!-- Main Content -->
	<div class="p-8">
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<!-- Settings Section (2/3) -->
			<div class="lg:col-span-2 space-y-6">
				<!-- Profile Settings -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<h2 class="text-xl font-bold text-gray-900 mb-6">Profile</h2>
					
					<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
						<!-- Avatar -->
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-3">Avatar</label>
							<div class="flex items-center gap-4">
								<div class="relative">
									{#if avatarPreview}
										<img src={avatarPreview} alt="Avatar" class="w-20 h-20 rounded-full object-cover" />
									{:else}
										<div class="w-20 h-20 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500"></div>
									{/if}
								</div>
								<div class="flex-1">
									<input
										type="file"
										accept="image/*"
										onchange={handleAvatarChange}
										class="hidden"
										id="avatar-upload"
									/>
									<label
										for="avatar-upload"
										class="inline-block px-4 py-2 bg-gray-100 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-200 cursor-pointer transition-colors"
									>
										Choose Image
									</label>
									{#if avatarFile}
										<button
											onclick={uploadAvatar}
											disabled={uploadingAvatar}
											class="ml-2 px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 transition-colors"
										>
											{uploadingAvatar ? 'Uploading...' : 'Upload'}
										</button>
									{/if}
								</div>
							</div>
						</div>

						<!-- Bio -->
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-3">Bio</label>
							<textarea
								bind:value={bio}
								placeholder="Tell people about yourself..."
								rows="3"
								class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
							></textarea>
							<button
								onclick={saveProfile}
								disabled={savingProfile}
								class="mt-2 px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 transition-colors"
							>
								{savingProfile ? 'Saving...' : 'Save Bio'}
							</button>
						</div>
					</div>
				</div>

				<!-- Theme Section -->
				<div class="bg-white rounded-2xl border border-gray-200 p-6">
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-xl font-bold text-gray-900">Theme</h2>
						{#if applying}
							<span class="text-sm text-indigo-600 flex items-center gap-2">
								<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Applying...
							</span>
						{/if}
					</div>
					
					<div class="flex gap-3 mb-6">
						{#each categories as cat}
							<button
								onclick={() => selectedCategory = cat.id}
								class="px-6 py-2.5 rounded-full text-sm font-medium transition-all {selectedCategory === cat.id
									? 'bg-gray-900 text-white'
									: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
							>
								{cat.label}
							</button>
						{/each}
					</div>

					<!-- Theme Grid -->
					<div class="grid grid-cols-2 md:grid-cols-3 gap-4">
						{#each themes[selectedCategory] || [] as themeItem}
							{@const colors = getPresetColors(themeItem.preset)}
							<button
								onclick={() => selectTheme(themeItem.id, themeItem.preset)}
								disabled={applying}
								class="relative group cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
							>
								<!-- Theme Preview Card -->
								<div
									class="aspect-[3/4] rounded-2xl overflow-hidden border-2 transition-all {currentTheme === themeItem.id
										? 'border-indigo-600 shadow-xl'
										: 'border-gray-200 hover:border-gray-300 hover:shadow-lg'}"
									style="background: {colors.bg};"
								>
									<!-- Mock Profile Preview -->
									<div class="h-full flex flex-col items-center justify-start p-4 pt-8">
										<!-- Avatar -->
										<div class="w-12 h-12 rounded-full bg-gray-300 mb-2"></div>
										
										<!-- Name -->
										<div
											class="text-sm font-bold mb-1"
											style="color: {colors.textColor};"
										>
											{themeItem.name}
										</div>
										
										<!-- Social Icons -->
										<div class="flex gap-1 mb-4">
											{#each [1, 2, 3] as _}
												<div class="w-3 h-3 rounded-full" style="background: {colors.accent}; opacity: 0.7;"></div>
											{/each}
										</div>

										<!-- Links Preview -->
										<div class="w-full space-y-2">
											{#each [1, 2] as _}
												<div
													class="w-full h-8 rounded-lg"
													style="background: {colors.cardBg};"
												>
													<div class="h-full flex items-center justify-center">
														<div class="w-2/3 h-2 rounded" style="background: {colors.cardText}; opacity: 0.3;"></div>
													</div>
												</div>
											{/each}
										</div>
									</div>
								</div>

								<!-- Selected Indicator -->
								{#if currentTheme === themeItem.id}
									<div class="absolute top-2 right-2 w-6 h-6 bg-indigo-600 rounded-full flex items-center justify-center">
										<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
										</svg>
									</div>
								{/if}

								<!-- Theme Name -->
								<p class="mt-2 text-sm font-medium text-gray-700 text-center">{themeItem.name}</p>
							</button>
						{/each}
					</div>

					<p class="mt-4 text-sm text-gray-500">
						Selecting a theme will apply styles to your page background and all link/text groups. You can customize individual groups in the Bio page.
					</p>
				</div>
			</div>

			<!-- Preview Panel (1/3) -->
			<div class="lg:col-span-1">
				<div class="sticky top-24">
					<ProfilePreview {profile} {links} {blocks} showInactive={true} />
				</div>
			</div>
		</div>
	</div>
</div>
