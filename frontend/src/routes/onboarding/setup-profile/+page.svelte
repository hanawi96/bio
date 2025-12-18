<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { profileApi } from '$lib/api/profile';
	import { toast } from 'svelte-sonner';
	import { onboardingData } from '$lib/stores/onboarding';
	import { fade, fly } from 'svelte/transition';

	let bio = '';
	let avatarFile: File | null = null;
	let avatarPreview = '';
	let loading = false;
	let uploading = false;
	let mounted = false;

	onMount(() => {
		mounted = true;
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
				toast.error('Kích thước file phải nhỏ hơn 5MB');
				return;
			}

			if (!file.type.startsWith('image/')) {
				toast.error('Vui lòng chọn file ảnh');
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
			toast.error(err.message || 'Không thể cập nhật profile');
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
	<title>Hoàn thiện Profile - LinkBio</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100 px-4 py-12">
	{#if mounted}
		<div class="max-w-md w-full" in:fade={{ duration: 300 }}>
			<!-- App Icon -->
			<div class="flex justify-center mb-8" in:fly={{ y: -20, duration: 400, delay: 100 }}>
				<div class="relative">
					<div class="w-20 h-20 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-[22px] shadow-2xl flex items-center justify-center">
						<svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
						</svg>
					</div>
					<div class="absolute -inset-1 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-[24px] blur-xl opacity-30 -z-10"></div>
				</div>
			</div>

			<!-- Title -->
			<div class="text-center mb-10" in:fly={{ y: -20, duration: 400, delay: 200 }}>
				<h1 class="text-3xl font-bold text-gray-900 mb-3 tracking-tight">Hoàn thiện Profile</h1>
				<p class="text-gray-500 text-base">Thêm ảnh và bio để cá nhân hóa trang của bạn</p>
			</div>

			<!-- Card -->
			<div class="bg-white rounded-3xl shadow-xl overflow-hidden" in:fly={{ y: 20, duration: 400, delay: 300 }}>
				<!-- Avatar Section -->
				<div class="p-8 flex flex-col items-center border-b border-gray-100">
					<label class="cursor-pointer group">
						<input
							type="file"
							accept="image/*"
							on:change={handleAvatarChange}
							class="hidden"
						/>
						<div class="relative">
							{#if avatarPreview}
								<div class="w-32 h-32 rounded-full overflow-hidden border-4 border-gray-100 shadow-lg group-hover:border-indigo-200 transition-all">
									<img
										src={avatarPreview}
										alt="Avatar preview"
										class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
									/>
								</div>
							{:else}
								<div class="w-32 h-32 rounded-full bg-gradient-to-br from-indigo-100 to-purple-100 flex items-center justify-center border-4 border-gray-100 shadow-lg group-hover:border-indigo-200 group-hover:scale-105 transition-all">
									<svg class="w-12 h-12 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
									</svg>
								</div>
							{/if}
							
							<!-- Camera Badge -->
							<div class="absolute bottom-0 right-0 w-10 h-10 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-full flex items-center justify-center border-4 border-white shadow-lg group-hover:scale-110 transition-transform">
								<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
								</svg>
							</div>
						</div>
					</label>
					<p class="mt-4 text-sm font-medium text-gray-700">Nhấn để tải ảnh lên</p>
					<p class="text-xs text-gray-500 mt-1">JPG, PNG hoặc GIF (tối đa 5MB)</p>
				</div>

				<!-- Bio Section -->
				<div class="p-8">
					<label for="bio" class="block text-sm font-semibold text-gray-900 mb-3">
						Bio
					</label>
					<div class="relative">
						<textarea
							id="bio"
							bind:value={bio}
							placeholder="Kể về bản thân bạn..."
							rows="4"
							maxlength="150"
							class="w-full px-4 py-3 bg-gray-50 border-2 border-gray-200 rounded-2xl text-base text-gray-900 placeholder-gray-400 focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all resize-none outline-none"
						></textarea>
					</div>
					<div class="flex justify-between items-center mt-3 px-1">
						<p class="text-xs text-gray-500">Không bắt buộc</p>
						<p class="text-xs font-medium" class:text-gray-500={bio.length < 140} class:text-amber-600={bio.length >= 140 && bio.length < 150} class:text-red-600={bio.length >= 150}>
							{bio.length}/150
						</p>
					</div>
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="mt-8 space-y-3" in:fly={{ y: 20, duration: 400, delay: 400 }}>
				<button
					type="button"
					on:click={handleContinue}
					disabled={loading}
					class="w-full bg-gradient-to-r from-indigo-600 to-purple-600 text-white py-4 px-6 rounded-2xl font-semibold text-base shadow-lg shadow-indigo-500/30 hover:shadow-xl hover:shadow-indigo-500/40 disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none active:scale-[0.98] transition-all"
				>
					{#if uploading}
						<span class="flex items-center justify-center gap-2">
							<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Đang tải lên...
						</span>
					{:else if loading}
						<span class="flex items-center justify-center gap-2">
							<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Đang lưu...
						</span>
					{:else}
						Tiếp tục
					{/if}
				</button>

				<button
					type="button"
					on:click={handleSkip}
					disabled={loading}
					class="w-full bg-white text-gray-700 py-4 px-6 rounded-2xl font-semibold text-base border-2 border-gray-200 hover:border-gray-300 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed active:scale-[0.98] transition-all"
				>
					Bỏ qua bây giờ
				</button>
			</div>

			<!-- Progress Indicator -->
			<div class="mt-8 flex items-center justify-center gap-2" in:fly={{ y: 10, duration: 400, delay: 500 }}>
				<div class="w-2 h-2 rounded-full bg-indigo-600"></div>
				<div class="w-2 h-2 rounded-full bg-indigo-600"></div>
				<div class="w-2 h-2 rounded-full bg-gray-300"></div>
			</div>
		</div>
	{/if}
</div>
