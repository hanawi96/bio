<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { themePresets } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { onboardingData } from '$lib/stores/onboarding';

	let loading = false;
	let selectedTheme = 'default';
	
	// Get onboarding data for preview
	$: profileData = $onboardingData;

	const themes = [
		{ 
			id: 'default', 
			name: 'Default', 
			description: 'Clean and modern',
			preview: { bg: 'linear-gradient(to top, #faf5ff, #eff6ff)', card: '#ffffff', text: '#111827' }
		},
		{ 
			id: 'dark', 
			name: 'Dark', 
			description: 'Sleek dark mode',
			preview: { bg: '#111827', card: '#1f2937', text: '#ffffff' }
		},
		{ 
			id: 'minimal', 
			name: 'Minimal', 
			description: 'Simple and clean',
			preview: { bg: '#ffffff', card: '#ffffff', text: '#000000' }
		},
		{ 
			id: 'vibrant', 
			name: 'Vibrant', 
			description: 'Bold and colorful',
			preview: { bg: '#fef3c7', card: '#ffffff', text: '#78350f' }
		},
		{ 
			id: 'mcalpine', 
			name: 'McAlpine', 
			description: 'Professional dark',
			preview: { bg: '#1a1a1a', card: '#333333', text: '#ffffff' }
		},
		{ 
			id: 'yoga', 
			name: 'Yoga', 
			description: 'Calm and peaceful',
			preview: { bg: '#b8c5d6', card: '#ffffff', text: '#4a6fa5' }
		},
		{ 
			id: 'jerry', 
			name: 'Jerry', 
			description: 'Minimalist dark',
			preview: { bg: '#000000', card: '#1a1a1a', text: '#ffffff' }
		}
	];

	onMount(() => {
		// TODO: Uncomment for production
		// if (!$auth.user) {
		// 	goto('/auth/login');
		// 	return;
		// }

		// if ($auth.user.username && $auth.user.username.startsWith('temp_')) {
		// 	goto('/onboarding/setup-url');
		// }
	});

	async function handleContinue() {
		loading = true;
		try {
			const preset = themePresets[selectedTheme];
			if (!preset) {
				throw new Error('Theme not found');
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
							{#each themes.filter(t => t.id === selectedTheme) as currentThemePreview}
								<div class="h-full overflow-y-auto relative" style="background: {currentThemePreview.preview.bg}">
									<!-- Subscribe & Share buttons -->
									<div class="sticky top-0 left-0 right-0 px-6 py-4 flex justify-between items-center z-10">
										<button class="flex items-center gap-1.5 px-3 py-1.5 bg-white/90 backdrop-blur-sm rounded-full text-xs font-semibold text-gray-800 shadow-sm">
											<svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20">
												<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
											</svg>
											Subscribe
										</button>
										<button class="w-8 h-8 bg-white/90 backdrop-blur-sm rounded-full flex items-center justify-center shadow-sm">
											<svg class="w-4 h-4 text-gray-800" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
											</svg>
										</button>
									</div>

									<!-- Profile Section -->
									<div class="px-6 pb-6">
										<!-- Avatar -->
										<div class="flex justify-center mb-4">
											{#if profileData.avatarPreview}
												<img
													src={profileData.avatarPreview}
													alt="Avatar"
													class="w-20 h-20 rounded-full object-cover shadow-lg"
												/>
											{:else}
												<div class="w-20 h-20 rounded-full bg-gradient-to-br from-violet-500 to-indigo-500 shadow-lg flex items-center justify-center">
													<svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
													</svg>
												</div>
											{/if}
										</div>

										<!-- Username & Bio -->
										<div class="text-center mb-6">
											<h3 class="text-lg font-bold mb-1" style="color: {currentThemePreview.preview.text}">
												@{$auth?.user?.username || 'yourname'}
											</h3>
											<p class="text-sm opacity-75" style="color: {currentThemePreview.preview.text}">
												{profileData.bio || 'Your bio goes here'}
											</p>
										</div>

										<!-- Links Section -->
										<div class="space-y-3 px-4">
											<div class="rounded-xl p-3.5 text-center font-medium text-sm shadow-sm transition-transform hover:scale-105" style="background: {currentThemePreview.preview.card}; color: {currentThemePreview.preview.text}">
												My Website
											</div>
											<div class="rounded-xl p-3.5 text-center font-medium text-sm shadow-sm transition-transform hover:scale-105" style="background: {currentThemePreview.preview.card}; color: {currentThemePreview.preview.text}">
												Instagram
											</div>
											<div class="rounded-xl p-3.5 text-center font-medium text-sm shadow-sm transition-transform hover:scale-105" style="background: {currentThemePreview.preview.card}; color: {currentThemePreview.preview.text}">
												YouTube
											</div>
											<div class="rounded-xl p-3.5 text-center font-medium text-sm shadow-sm transition-transform hover:scale-105" style="background: {currentThemePreview.preview.card}; color: {currentThemePreview.preview.text}">
												Twitter
											</div>
											<div class="rounded-xl p-3.5 text-center font-medium text-sm shadow-sm transition-transform hover:scale-105" style="background: {currentThemePreview.preview.card}; color: {currentThemePreview.preview.text}">
												TikTok
											</div>
										</div>

										<!-- Footer -->
										<div class="text-center mt-8 opacity-50">
											<p class="text-xs font-medium" style="color: {currentThemePreview.preview.text}">Made with LinkBio</p>
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
