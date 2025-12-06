<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { toast } from 'svelte-sonner';
	
	import { profileApi } from '$lib/api/profile';
	import type { Profile } from '$lib/api/profile';

	let profile: Profile | null = null;
	let initialLoading = true;

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
						<Label for="username">Username</Label>
						<Input id="username" value={displayProfile.username} disabled class="mt-2" />
						<p class="text-xs text-gray-500 mt-1">Your profile URL: linkbio.com/{displayProfile.username}</p>
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
