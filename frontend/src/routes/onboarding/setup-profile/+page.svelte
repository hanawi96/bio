<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { profileApi } from '$lib/api/profile';
	import { toast } from 'svelte-sonner';
	import { onboardingData } from '$lib/stores/onboarding';

	let bio = '';
	let avatarFile: File | null = null;
	let avatarPreview = '';
	let loading = false;
	let uploading = false;

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

	function handleAvatarChange(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		
		if (file) {
			if (file.size > 5 * 1024 * 1024) {
				toast.error('File size must be less than 5MB');
				return;
			}

			if (!file.type.startsWith('image/')) {
				toast.error('Please select an image file');
				return;
			}

			avatarFile = file;
			const reader = new FileReader();
			reader.onload = (e) => {
				avatarPreview = e.target?.result as string;
			};
			reader.readAsDataURL(file);
		}
	}

	async function handleContinue() {
		loading = true;
		try {
			// Save to store for preview in next step
			onboardingData.setData({
				avatarPreview,
				bio: bio.trim()
			});

			// Upload avatar if selected
			if (avatarFile) {
				uploading = true;
				await profileApi.uploadAvatar(avatarFile, $auth.token!);
			}

			// Update bio if provided
			if (bio.trim()) {
				await profileApi.updateProfile({ bio: bio.trim() }, $auth.token!);
			}

			goto('/onboarding/choose-theme');
		} catch (err: any) {
			toast.error(err.message || 'Failed to update profile');
		} finally {
			loading = false;
			uploading = false;
		}
	}

	function handleSkip() {
		// Clear store when skipping
		onboardingData.clear();
		goto('/onboarding/choose-theme');
	}
</script>

<svelte:head>
	<title>Setup Profile - LinkBio</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-violet-50 via-white to-indigo-50 px-4 py-12">
	<div class="max-w-lg mx-auto">
		<!-- Header -->
		<div class="text-center mb-8">
			<div class="w-16 h-16 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
				<svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
				</svg>
			</div>
			<h1 class="text-3xl font-bold text-gray-900 mb-2">Complete Your Profile</h1>
			<p class="text-gray-600">Add a photo and bio to personalize your page</p>
		</div>

		<!-- Form Card - iOS Style -->
		<div class="bg-white rounded-3xl shadow-xl overflow-hidden">
			<!-- Avatar Section -->
			<div class="p-8 border-b border-gray-100">
				<label class="block text-center">
					<input
						type="file"
						accept="image/*"
						on:change={handleAvatarChange}
						class="hidden"
					/>
					<div class="relative inline-block cursor-pointer group">
						{#if avatarPreview}
							<img
								src={avatarPreview}
								alt="Avatar preview"
								class="w-32 h-32 rounded-full object-cover border-4 border-white shadow-lg group-hover:opacity-90 transition-opacity"
							/>
						{:else}
							<div class="w-32 h-32 rounded-full bg-gradient-to-br from-violet-500 to-indigo-500 flex items-center justify-center border-4 border-white shadow-lg group-hover:scale-105 transition-transform">
								<svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
								</svg>
							</div>
						{/if}
						
						<!-- Camera Icon Overlay -->
						<div class="absolute bottom-0 right-0 w-10 h-10 bg-violet-600 rounded-full flex items-center justify-center border-4 border-white shadow-lg group-hover:bg-violet-700 transition-colors">
							<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
							</svg>
						</div>
					</div>
					<p class="mt-3 text-sm text-gray-600">Tap to upload photo</p>
					<p class="text-xs text-gray-500 mt-1">JPG, PNG or GIF (max 5MB)</p>
				</label>
			</div>

			<!-- Bio Section -->
			<div class="p-8">
				<label for="bio" class="block text-sm font-semibold text-gray-700 mb-3">
					Bio
				</label>
				<textarea
					id="bio"
					bind:value={bio}
					placeholder="Tell people about yourself..."
					rows="4"
					maxlength="150"
					class="w-full px-4 py-3 border-2 border-gray-200 rounded-2xl focus:border-violet-500 focus:ring-2 focus:ring-violet-200 transition-all resize-none"
				></textarea>
				<div class="flex justify-between items-center mt-2">
					<p class="text-xs text-gray-500">Optional</p>
					<p class="text-xs text-gray-500">{bio.length}/150</p>
				</div>
			</div>
		</div>

		<!-- Action Buttons -->
		<div class="mt-8 space-y-3">
			<button
				type="button"
				on:click={handleContinue}
				disabled={loading}
				class="w-full bg-gradient-to-r from-violet-600 to-indigo-600 text-white py-4 px-6 rounded-2xl font-semibold hover:shadow-lg hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 transition-all"
			>
				{#if uploading}
					<span class="flex items-center justify-center gap-2">
						<svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Uploading...
					</span>
				{:else if loading}
					Saving...
				{:else}
					Continue
				{/if}
			</button>

			<button
				type="button"
				on:click={handleSkip}
				disabled={loading}
				class="w-full bg-white text-gray-700 py-4 px-6 rounded-2xl font-semibold border-2 border-gray-200 hover:border-gray-300 hover:bg-gray-50 transition-all disabled:opacity-50"
			>
				Skip for now
			</button>
		</div>

		<!-- Progress Indicator -->
		<div class="mt-8 flex items-center justify-center gap-2">
			<div class="w-2 h-2 rounded-full bg-violet-600"></div>
			<div class="w-2 h-2 rounded-full bg-violet-600"></div>
			<div class="w-2 h-2 rounded-full bg-gray-300"></div>
		</div>
	</div>
</div>
