<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { api } from '$lib/api/client';
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';

	let customUrl = '';
	let error = '';
	let loading = false;
	let checking = false;
	let isAvailable: boolean | null = null;
	let mounted = false;

	onMount(() => {
		mounted = true;
		
		// Nếu chưa login thì redirect về trang chủ
		if (!$auth.user) {
			goto('/auth/login');
			return;
		}
		
		// Nếu username không phải temp (đã setup rồi) thì redirect về dashboard
		if ($auth.user.username && !$auth.user.username.startsWith('temp_')) {
			goto('/dashboard');
		}
	});

	let checkTimeout: any;
	async function checkAvailability() {
		if (!customUrl || customUrl.length < 3) {
			isAvailable = null;
			return;
		}

		clearTimeout(checkTimeout);
		checkTimeout = setTimeout(async () => {
			checking = true;
			try {
				const response = await api.get(`/auth/check-username/${customUrl}`);
				isAvailable = response.available;
			} catch (err) {
				isAvailable = false;
			} finally {
				checking = false;
			}
		}, 500);
	}

	$: if (customUrl) checkAvailability();

	async function handleSubmit() {
		if (!customUrl || customUrl.length < 3) {
			error = 'URL phải có ít nhất 3 ký tự';
			return;
		}

		if (!isAvailable) {
			error = 'URL này đã được sử dụng';
			return;
		}

		loading = true;
		error = '';

		try {
			await api.patch('/auth/setup-username', { username: customUrl }, $auth.token!);
			
			// Cập nhật user trong store
			auth.updateUser({ ...$auth.user, username: customUrl });
			
			goto('/onboarding/setup-profile');
		} catch (err: any) {
			error = err.message || 'Có lỗi xảy ra';
		} finally {
			loading = false;
		}
	}

	function sanitizeInput(value: string) {
		return value.toLowerCase().replace(/[^a-z0-9_-]/g, '');
	}

	$: customUrl = sanitizeInput(customUrl);
</script>

<svelte:head>
	<title>Thiết lập URL - LinkBio</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100 px-4 py-12">
	{#if mounted}
		<div class="max-w-md w-full" in:fade={{ duration: 300 }}>
			<!-- App Icon -->
			<div class="flex justify-center mb-8" in:fly={{ y: -20, duration: 400, delay: 100 }}>
				<div class="relative">
					<div class="w-20 h-20 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-[22px] shadow-2xl flex items-center justify-center">
						<svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
						</svg>
					</div>
					<div class="absolute -inset-1 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-[24px] blur-xl opacity-30 -z-10"></div>
				</div>
			</div>

			<!-- Title -->
			<div class="text-center mb-10" in:fly={{ y: -20, duration: 400, delay: 200 }}>
				<h1 class="text-3xl font-bold text-gray-900 mb-3 tracking-tight">Chọn URL cho trang của bạn</h1>
				<p class="text-gray-500 text-base">Đây sẽ là địa chỉ duy nhất để mọi người tìm thấy bạn</p>
			</div>

			<!-- Card -->
			<div class="bg-white rounded-3xl shadow-xl overflow-hidden" in:fly={{ y: 20, duration: 400, delay: 300 }}>
				<form on:submit|preventDefault={handleSubmit} class="p-8">
					{#if error}
						<div class="mb-6 bg-red-50 border border-red-100 text-red-600 px-4 py-3 rounded-2xl text-sm flex items-center gap-3" in:fly={{ y: -10, duration: 200 }}>
							<svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
							</svg>
							<span>{error}</span>
						</div>
					{/if}

					<div class="space-y-6">
						<!-- URL Input -->
						<div>
							<label for="customUrl" class="block text-sm font-semibold text-gray-900 mb-3">
								URL của bạn
							</label>
							<div class="relative group">
								<div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
									<span class="text-gray-400 text-base font-medium">linkbio.com/</span>
								</div>
								<input
									id="customUrl"
									type="text"
									bind:value={customUrl}
									placeholder="ten-cua-ban"
									required
									minlength="3"
									maxlength="30"
									class="w-full pl-36 pr-14 py-4 bg-gray-50 border-2 border-gray-200 rounded-2xl text-base font-medium text-gray-900 placeholder-gray-400 focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all outline-none"
								/>
								<div class="absolute inset-y-0 right-0 pr-5 flex items-center">
									{#if checking}
										<div class="w-5 h-5 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
									{:else if isAvailable === true}
										<div class="w-6 h-6 bg-green-500 rounded-full flex items-center justify-center" in:fly={{ scale: 0.5, duration: 200 }}>
											<svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
											</svg>
										</div>
									{:else if isAvailable === false}
										<div class="w-6 h-6 bg-red-500 rounded-full flex items-center justify-center" in:fly={{ scale: 0.5, duration: 200 }}>
											<svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
											</svg>
										</div>
									{/if}
								</div>
							</div>
							<p class="mt-3 text-xs text-gray-500 px-1">
								Chỉ sử dụng chữ cái, số, gạch ngang (-) và gạch dưới (_)
							</p>
						</div>

						<!-- Preview -->
						{#if customUrl && customUrl.length >= 3}
							<div class="bg-gradient-to-br from-indigo-50 to-purple-50 border border-indigo-100 rounded-2xl p-5" in:fly={{ y: 10, duration: 300 }}>
								<div class="flex items-start gap-3">
									<div class="w-10 h-10 bg-white rounded-xl flex items-center justify-center flex-shrink-0 shadow-sm">
										<svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
										</svg>
									</div>
									<div class="flex-1 min-w-0">
										<p class="text-xs font-medium text-gray-600 mb-1">URL của bạn sẽ là:</p>
										<p class="text-base font-bold text-indigo-600 break-all">
											linkbio.com/{customUrl}
										</p>
									</div>
								</div>
							</div>
						{/if}

						<!-- Submit Button -->
						<button
							type="submit"
							disabled={loading || !customUrl || isAvailable !== true}
							class="w-full bg-gradient-to-r from-indigo-600 to-purple-600 text-white py-4 px-6 rounded-2xl font-semibold text-base shadow-lg shadow-indigo-500/30 hover:shadow-xl hover:shadow-indigo-500/40 disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none active:scale-[0.98] transition-all"
						>
							{#if loading}
								<span class="flex items-center justify-center gap-2">
									<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
										<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
										<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
									</svg>
									Đang xử lý...
								</span>
							{:else}
								Tiếp tục
							{/if}
						</button>
					</div>
				</form>
			</div>

			<!-- Footer Note -->
			<p class="text-center text-sm text-gray-500 mt-6 px-4" in:fly={{ y: 10, duration: 400, delay: 400 }}>
				Bạn có thể thay đổi URL này sau trong phần cài đặt
			</p>

			<!-- Progress Indicator -->
			<div class="mt-8 flex items-center justify-center gap-2" in:fly={{ y: 10, duration: 400, delay: 500 }}>
				<div class="w-2 h-2 rounded-full bg-indigo-600"></div>
				<div class="w-2 h-2 rounded-full bg-gray-300"></div>
				<div class="w-2 h-2 rounded-full bg-gray-300"></div>
			</div>
		</div>
	{/if}
</div>
