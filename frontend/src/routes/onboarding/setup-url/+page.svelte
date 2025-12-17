<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { api } from '$lib/api/client';
	import { onMount } from 'svelte';

	let customUrl = '';
	let error = '';
	let loading = false;
	let checking = false;
	let isAvailable: boolean | null = null;

	onMount(() => {
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

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-violet-50 via-white to-indigo-50 px-4">
	<div class="max-w-lg w-full">
		<div class="text-center mb-8">
			<div class="w-16 h-16 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
				<svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
				</svg>
			</div>
			<h1 class="text-3xl font-bold text-gray-900 mb-2">Chọn URL cho trang của bạn</h1>
			<p class="text-gray-600">Đây sẽ là địa chỉ duy nhất để mọi người tìm thấy bạn</p>
		</div>

		<div class="bg-white rounded-2xl shadow-xl p-8">
			<form on:submit|preventDefault={handleSubmit} class="space-y-6">
				{#if error}
					<div class="bg-red-50 text-red-600 p-4 rounded-lg text-sm">
						{error}
					</div>
				{/if}

				<div>
					<label for="customUrl" class="block text-sm font-semibold text-gray-700 mb-2">
						URL của bạn
					</label>
					<div class="relative">
						<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-gray-500 text-sm">
							linkbio.com/
						</div>
						<input
							id="customUrl"
							type="text"
							bind:value={customUrl}
							placeholder="ten-cua-ban"
							required
							minlength="3"
							maxlength="30"
							class="w-full pl-32 pr-12 py-3 border-2 border-gray-200 rounded-xl focus:border-violet-500 focus:ring-2 focus:ring-violet-200 transition-all"
						/>
						<div class="absolute inset-y-0 right-0 pr-4 flex items-center">
							{#if checking}
								<div class="w-5 h-5 border-2 border-violet-600 border-t-transparent rounded-full animate-spin"></div>
							{:else if isAvailable === true}
								<svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
								</svg>
							{:else if isAvailable === false}
								<svg class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
								</svg>
							{/if}
						</div>
					</div>
					<p class="mt-2 text-xs text-gray-500">
						Chỉ sử dụng chữ cái, số, gạch ngang (-) và gạch dưới (_)
					</p>
					{#if isAvailable === true}
						<p class="mt-2 text-sm text-green-600 font-medium">
							✓ URL này có sẵn!
						</p>
					{:else if isAvailable === false}
						<p class="mt-2 text-sm text-red-600 font-medium">
							✗ URL này đã được sử dụng
						</p>
					{/if}
				</div>

				{#if customUrl}
					<div class="bg-violet-50 border border-violet-200 rounded-xl p-4">
						<p class="text-sm text-gray-600 mb-1">URL của bạn sẽ là:</p>
						<p class="text-lg font-semibold text-violet-600 break-all">
							linkbio.com/{customUrl}
						</p>
					</div>
				{/if}

				<button
					type="submit"
					disabled={loading || !customUrl || isAvailable !== true}
					class="w-full bg-gradient-to-r from-violet-600 to-indigo-600 text-white py-3 px-6 rounded-xl font-semibold hover:shadow-lg hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 transition-all"
				>
					{loading ? 'Đang xử lý...' : 'Tiếp tục'}
				</button>
			</form>
		</div>

		<p class="text-center text-sm text-gray-500 mt-6">
			Bạn có thể thay đổi URL này sau trong phần cài đặt
		</p>

		<!-- Progress Indicator -->
		<div class="mt-8 flex items-center justify-center gap-2">
			<div class="w-2 h-2 rounded-full bg-violet-600"></div>
			<div class="w-2 h-2 rounded-full bg-gray-300"></div>
			<div class="w-2 h-2 rounded-full bg-gray-300"></div>
		</div>
	</div>
</div>
