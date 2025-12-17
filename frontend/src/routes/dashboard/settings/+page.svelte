<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { toast } from 'svelte-sonner';
	
	import { profileApi } from '$lib/api/profile';
	import type { Profile } from '$lib/api/profile';
	import { api } from '$lib/api/client';

	let profile: Profile | null = null;
	let initialLoading = true;
	let editingUsername = false;
	let newUsername = '';
	let checkingUsername = false;
	let usernameAvailable: boolean | null = null;
	let savingUsername = false;

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			initialLoading = true;
			
			try {
				profile = await profileApi.getMyProfile($auth.token!);
			} catch (profileError: any) {
				console.warn('Profile not found:', profileError);
			}
		} catch (error: any) {
			console.error('Failed to load data:', error);
			toast.error(error.message || 'Failed to load data');
		} finally {
			initialLoading = false;
		}
	}

	$: displayProfile = profile || { 
		username: 'loading...', 
		bio: '', 
		avatar_url: null,
		id: '',
		user_id: '',
		theme_config: {},
		custom_css: null,
		created_at: new Date(),
		updated_at: new Date()
	};

	function startEditUsername() {
		newUsername = displayProfile.username;
		editingUsername = true;
		usernameAvailable = null;
	}

	function cancelEditUsername() {
		editingUsername = false;
		newUsername = '';
		usernameAvailable = null;
	}

	let checkTimeout: any;
	async function checkUsernameAvailability() {
		if (!newUsername || newUsername.length < 3 || newUsername === displayProfile.username) {
			usernameAvailable = null;
			return;
		}

		clearTimeout(checkTimeout);
		checkTimeout = setTimeout(async () => {
			checkingUsername = true;
			try {
				const response = await api.get(`/auth/check-username/${newUsername}`);
				usernameAvailable = response.available;
			} catch (err) {
				usernameAvailable = false;
			} finally {
				checkingUsername = false;
			}
		}, 500);
	}

	$: if (newUsername && editingUsername) checkUsernameAvailability();

	async function saveUsername() {
		if (!newUsername || usernameAvailable !== true) return;

		savingUsername = true;
		try {
			await api.patch('/auth/setup-username', { username: newUsername }, $auth.token!);
			
			// Update user in store
			auth.updateUser({ ...$auth.user!, username: newUsername });
			
			// Update profile
			if (profile) {
				profile = { ...profile, username: newUsername };
			}
			
			toast.success('Username updated successfully');
			editingUsername = false;
		} catch (err: any) {
			toast.error(err.message || 'Failed to update username');
		} finally {
			savingUsername = false;
		}
	}

	function sanitizeUsername(value: string) {
		return value.toLowerCase().replace(/[^a-z0-9_-]/g, '');
	}

	$: newUsername = sanitizeUsername(newUsername);
</script>

<svelte:head>
	<title>Settings - LinkBio</title>
</svelte:head>

{#if initialLoading}
	<div class="h-full flex items-center justify-center bg-gray-50">
		<div class="text-center">
			<div class="relative w-20 h-20 mx-auto mb-6">
				<div class="absolute inset-0 border-4 border-purple-100 rounded-full"></div>
				<div class="absolute inset-0 border-4 border-transparent border-t-purple-600 border-r-blue-600 rounded-full animate-spin"></div>
			</div>
			<h3 class="text-lg font-semibold text-gray-900 mb-2">Loading</h3>
		</div>
	</div>
{:else}
<div class="h-full bg-gray-50">
	<!-- Page Header -->
	<div class="bg-white/80 backdrop-blur-xl border-b border-gray-200/50 px-8 h-16 sticky top-0 z-10">
		<div class="flex items-center justify-between h-full">
			<div>
				<h1 class="text-xl font-bold text-gray-900">Settings</h1>
				<p class="text-xs text-gray-500">Manage your account settings</p>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="p-8">
		<div class="max-w-2xl">
			<div class="bg-white rounded-xl p-6 shadow-sm space-y-6">
				<div>
					<h2 class="text-2xl font-bold text-gray-900 mb-1">Account Settings</h2>
					<p class="text-sm text-gray-600">Manage your account information</p>
				</div>

				<div class="space-y-4">
					<div>
						<div class="flex items-center justify-between mb-2">
							<Label for="username">Username</Label>
							{#if !editingUsername}
								<button
									type="button"
									on:click={startEditUsername}
									class="text-sm text-violet-600 hover:text-violet-700 font-medium"
								>
									Edit
								</button>
							{/if}
						</div>
						
						{#if editingUsername}
							<div class="space-y-3">
								<div class="relative">
									<Input 
										id="username" 
										bind:value={newUsername}
										placeholder="your-username"
										minlength="3"
										maxlength="30"
										class="pr-10"
									/>
									<div class="absolute inset-y-0 right-0 pr-3 flex items-center">
										{#if checkingUsername}
											<div class="w-4 h-4 border-2 border-violet-600 border-t-transparent rounded-full animate-spin"></div>
										{:else if usernameAvailable === true}
											<svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
											</svg>
										{:else if usernameAvailable === false}
											<svg class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
											</svg>
										{/if}
									</div>
								</div>
								
								{#if usernameAvailable === true}
									<p class="text-sm text-green-600">✓ Username available</p>
								{:else if usernameAvailable === false}
									<p class="text-sm text-red-600">✗ Username already taken</p>
								{/if}
								
								<div class="flex gap-2">
									<button
										type="button"
										on:click={saveUsername}
										disabled={savingUsername || usernameAvailable !== true}
										class="px-4 py-2 bg-violet-600 text-white rounded-lg hover:bg-violet-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium"
									>
										{savingUsername ? 'Saving...' : 'Save'}
									</button>
									<button
										type="button"
										on:click={cancelEditUsername}
										class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 text-sm font-medium"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<Input id="username" value={displayProfile.username} disabled class="mt-2" />
							<p class="text-xs text-gray-500 mt-1">Your profile URL: linkbio.com/{displayProfile.username}</p>
						{/if}
					</div>

					<div>
						<Label for="email">Email</Label>
						<Input id="email" type="email" value={$auth.user?.email} disabled class="mt-2" />
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
{/if}
