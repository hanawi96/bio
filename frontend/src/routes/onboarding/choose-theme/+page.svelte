<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { themePresets } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { onboardingData } from '$lib/stores/onboarding';
	import { themeMetadata } from '$lib/config/themes';

	let loading = false;
	let selectedTheme = 'default';
	
	// Get onboarding data for preview
	$: profileData = $onboardingData;

	// Use centralized theme metadata
	const themes = themeMetadata;



	async function handleContinue() {
		loading = true;
		try {
			// Find the selected theme metadata
			const themeInfo = themes.find(t => t.id === selectedTheme);
			if (!themeInfo) {
				throw new Error('Theme not found');
			}

			// Get the preset from theme store
			const preset = themePresets[themeInfo.preset];
			if (!preset) {
				throw new Error('Theme preset not found');
			}

			// Prepare theme data
			const themeData = {
				theme_name: selectedTheme,
				theme_config: preset.page,
				card_styles: preset.card,
				text_styles: JSON.stringify(preset.text),
				header_config: preset.header
			};

			// Apply theme
			await profileApi.applyTheme(themeData, $auth.token!);
			
			// Clear onboarding data
			onboardingData.clear();
			
			toast.success('Theme applied successfully');
			goto('/dashboard');
		} catch (err: any) {
			toast.error(err.message || 'Failed to apply theme');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Choose Theme - LinkBio</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-violet-50 via-white to-indigo-50 px-4 py-12">
	<div class="max-w-5xl mx-auto">
		<div class="text-center mb-12">
			<div class="w-16 h-16 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
				<svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"/>
				</svg>
			</div>
			<h1 class="text-3xl font-bold text-gray-900 mb-2">Choose Your Theme</h1>
			<p class="text-gray-600">Pick a style that represents you. You can change it anytime.</p>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
			<!-- Theme Selection -->
			<div class="space-y-4">
				<h2 class="text-xl font-bold text-gray-900 mb-4">Choose a theme</h2>
				<div class="grid grid-cols-1 gap-4">
					{#each themes as theme}
						<button
							type="button"
							on:click={() => selectedTheme = theme.id}
							class="relative bg-white rounded-xl p-4 shadow hover:shadow-lg transition-all border-2 group text-left"
							class:border-violet-600={selectedTheme === theme.id}
							class:border-gray-200={selectedTheme !== theme.id}
						>
							{#if selectedTheme === theme.id}
								<div class="absolute top-3 right-3 w-6 h-6 bg-violet-600 rounded-full flex items-center justify-center">
									<svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
									</svg>
								</div>
							{/if}

							<div class="flex items-center gap-4">
								<div class="w-16 h-16 rounded-lg overflow-hidden flex-shrink-0" style="background: {theme.preview.bg}">
									<div class="w-full h-full flex flex-col items-center justify-center p-2 gap-1">
										<div class="w-2 h-2 rounded-full" style="background: {theme.preview.card}; opacity: 0.8;"></div>
										<div class="w-full h-2 rounded" style="background: {theme.preview.card}"></div>
										<div class="w-full h-2 rounded" style="background: {theme.preview.card}"></div>
									</div>
								</div>
								<div class="flex-1">
									<h3 class="text-base font-bold text-gray-900">{theme.name}</h3>
									<p class="text-sm text-gray-600">{theme.description}</p>
								</div>
							</div>
						</button>
					{/each}
				</div>
			</div>

			<!-- Live Preview -->
			<div class="sticky top-8">
				<div class="text-center mb-6">
					<h2 class="text-2xl font-bold text-gray-900 mb-2">Live Preview</h2>
					<p class="text-sm text-gray-600">See how your profile looks</p>
				</div>
				
				<!-- iPhone Mockup -->
				<div class="relative mx-auto" style="width: 320px;">
					<!-- Phone Frame -->
					<div class="relative bg-gray-900 rounded-[3rem] p-3 shadow-2xl">
						<!-- Screen -->
						<div class="bg-white rounded-[2.5rem] overflow-hidden" style="height: 600px;">
							{#each themes.filter(t => t.id === selectedTheme) as currentTheme}
								{@const preset = themePresets[currentTheme.preset]}
								{@const header = preset?.header || {}}
								<div class="h-full overflow-y-auto relative" style="background: {currentTheme.preview.bg}">
									<!-- Header/Cover (if theme has it) -->
									{#if header.showCover && header.coverHeight > 0}
										<div 
											class="w-full relative"
											style="height: {header.coverHeight}px; {
												header.coverType === 'gradient' 
													? `background: linear-gradient(to bottom, ${header.coverGradientFrom}, ${header.coverGradientTo});`
													: `background: ${header.coverColor};`
											}"
										></div>
									{/if}
									
									<!-- Profile Section -->
									<div class="px-6 pb-6" style="margin-top: {header.showCover && header.coverHeight > 0 ? `-${(header.avatarSize || 96) / 2}px` : '0'}">
										<!-- Avatar -->
										<div class="flex {header.avatarAlign === 'left' ? 'justify-start' : header.avatarAlign === 'right' ? 'justify-end' : 'justify-center'} mb-4" style="position: relative; z-index: 10;">
											{#if profileData.avatarPreview}
												<img
													src={profileData.avatarPreview}
													alt="Avatar"
													class="object-cover shadow-lg"
													style="width: {header.avatarSize || 96}px; height: {header.avatarSize || 96}px; border-radius: {header.avatarShape === 'circle' ? '50%' : header.avatarShape === 'rounded' ? '20%' : '0'}; border: {header.avatarBorder || 4}px solid {header.avatarBorderColor || '#ffffff'};"
												/>
											{:else}
												<div 
													class="bg-gradient-to-br from-violet-500 to-indigo-500 shadow-lg flex items-center justify-center"
													style="width: {header.avatarSize || 96}px; height: {header.avatarSize || 96}px; border-radius: {header.avatarShape === 'circle' ? '50%' : header.avatarShape === 'rounded' ? '20%' : '0'}; border: {header.avatarBorder || 4}px solid {header.avatarBorderColor || '#ffffff'};"
												>
													<svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
													</svg>
												</div>
											{/if}
										</div>

										<!-- Username & Bio -->
										<div class="mb-6" style="text-align: {header.bioAlign || 'center'}">
											<h3 class="text-lg font-bold mb-1" style="color: {currentTheme.preview.text}">
												@{$auth?.user?.username || 'yourname'}
											</h3>
											<p class="text-sm opacity-75" style="color: {currentTheme.preview.text}">
												{profileData.bio || 'Your bio goes here'}
											</p>
										</div>

										<!-- Links Section -->
										<div class="space-y-3">
											{#each ['My Website', 'Instagram', 'YouTube', 'Twitter', 'TikTok'] as linkTitle}
												<div 
													class="p-3.5 text-center font-medium text-sm shadow-sm transition-transform hover:scale-105" 
													style="background: {currentTheme.preview.card}; color: {currentTheme.preview.text}; border-radius: {preset?.card?.cardBorderRadius || 12}px;"
												>
													{linkTitle}
												</div>
											{/each}
										</div>

										<!-- Footer -->
										<div class="text-center mt-8 opacity-50">
											<p class="text-xs font-medium" style="color: {currentTheme.preview.text}">Made with LinkBio</p>
										</div>
									</div>
								</div>
							{/each}
						</div>
						
						<!-- Notch -->
						<div class="absolute top-0 left-1/2 -translate-x-1/2 w-32 h-6 bg-gray-900 rounded-b-2xl"></div>
					</div>
				</div>
			</div>
		</div>

		<div class="flex items-center justify-center gap-4">
			<button
				type="button"
				on:click={() => goto('/onboarding/setup-profile')}
				class="px-8 py-4 bg-white text-gray-700 rounded-2xl font-semibold border-2 border-gray-200 hover:border-gray-300 hover:bg-gray-50 transition-all flex items-center gap-2"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
				</svg>
				Back
			</button>
			<button
				type="button"
				on:click={handleContinue}
				disabled={loading}
				class="px-12 py-4 bg-gradient-to-r from-violet-600 to-indigo-600 text-white rounded-2xl font-semibold hover:shadow-lg hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 transition-all"
			>
				{loading ? 'Applying theme...' : 'Continue to Dashboard'}
			</button>
		</div>

		<!-- Progress Indicator -->
		<div class="mt-8 flex items-center justify-center gap-2">
			<div class="w-2 h-2 rounded-full bg-violet-600"></div>
			<div class="w-2 h-2 rounded-full bg-violet-600"></div>
			<div class="w-2 h-2 rounded-full bg-violet-600"></div>
		</div>
	</div>
</div>
