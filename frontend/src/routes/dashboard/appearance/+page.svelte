<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { Button } from '$lib/components/ui/button';
	import Avatar from '$lib/components/ui/avatar.svelte';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import { toast } from 'svelte-sonner';
	
	import ProfilePreview from '$lib/components/dashboard/ProfilePreview.svelte';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';

	let profile: Profile | null = null;
	let links: Link[] = [];
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
			
			try {
				const linksData = await linksApi.getLinks($auth.token!);
				links = linksData.sort((a, b) => a.position - b.position);
			} catch (linksError: any) {
				console.warn('No links found:', linksError);
				links = [];
			}
		} catch (error: any) {
			console.error('Failed to load data:', error);
			toast.error(error.message || 'Failed to load data');
		} finally {
			initialLoading = false;
		}
	}

	async function handleUpdateProfile() {
		if (!profile) return;
		try {
			const updated = await profileApi.updateProfile({
				bio: profile.bio,
				avatar_url: profile.avatar_url
			}, $auth.token!);
			profile = updated;
			toast.success('Profile updated!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to update profile');
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
	<title>Appearance - LinkBio</title>
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
				<h1 class="text-xl font-bold text-gray-900">Appearance</h1>
				<p class="text-xs text-gray-500">Customize your profile look</p>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="p-8">
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<!-- Editor Section (2/3) -->
			<div class="lg:col-span-2">
				<div class="bg-white rounded-xl p-6 shadow-sm space-y-6">
					{#if profile}
					<div class="space-y-6">
						<div>
							<Label>Profile Picture</Label>
							<div class="flex items-center gap-4 mt-2">
								<Avatar 
									src={profile.avatar_url || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + profile.username}
									alt="Avatar"
									class="w-20 h-20"
								/>
								<Button variant="outline">Change Avatar</Button>
							</div>
						</div>

						<div>
							<Label for="bio">Bio</Label>
							<Textarea 
								id="bio" 
								bind:value={profile.bio} 
								placeholder="Tell people about yourself..." 
								class="mt-2" 
								rows={3} 
							/>
						</div>

						<div>
							<Label>Theme</Label>
							<div class="grid grid-cols-3 gap-3 mt-2">
								<button class="aspect-square rounded-xl bg-gradient-to-br from-purple-500 to-blue-500 hover:scale-105 transition-transform border-2 border-purple-600"></button>
								<button class="aspect-square rounded-xl bg-gradient-to-br from-pink-500 to-orange-500 hover:scale-105 transition-transform"></button>
								<button class="aspect-square rounded-xl bg-gradient-to-br from-green-500 to-indigo-500 hover:scale-105 transition-transform"></button>
							</div>
						</div>

						<Button class="w-full" on:click={handleUpdateProfile}>Save Changes</Button>
					</div>
					{/if}
				</div>
			</div>

			<!-- Preview Section (1/3) -->
			<div class="lg:col-span-1">
				<div class="sticky top-32">
					<div class="text-center mb-6">
						<h3 class="text-lg font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent">
							Live Preview
						</h3>
						<p class="text-sm text-gray-500 mt-1">See how your profile looks</p>
					</div>
					<div class="relative">
						<div class="absolute inset-0 bg-gradient-to-r from-indigo-500 to-blue-500 rounded-3xl blur-2xl opacity-20"></div>
						<div class="relative">
							<ProfilePreview profile={displayProfile} {links} />
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
{/if}
