<script lang="ts">
	import { onMount } from 'svelte';
	import { dndzone } from 'svelte-dnd-action';
	import { auth } from '$lib/stores/auth';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Dialog from '$lib/components/ui/dialog';
	import { toast } from 'svelte-sonner';
	
	import LinkCard from '$lib/components/dashboard/LinkCard.svelte';
	import ProfilePreview from '$lib/components/dashboard/ProfilePreview.svelte';
	import { profileApi } from '$lib/api/profile';
	import { linksApi } from '$lib/api/links';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';

	let profile: Profile | null = null;
	let links: Link[] = [];
	let initialLoading = true;
	let isDragging = false;

	function handleCopyUrl() {
		if (!profile?.username) {
			toast.error('Profile not loaded yet');
			return;
		}
		const url = `${window.location.origin}/${profile.username}`;
		navigator.clipboard.writeText(url).then(() => {
			toast.success('Profile URL copied to clipboard!');
		}).catch(err => {
			console.error('Copy failed:', err);
			toast.error('Failed to copy URL');
		});
	}

	let showAddLinkDialog = false;
	let newLink = { title: '', url: '' };

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			initialLoading = true;
			const startTime = Date.now();
			
			let profileData = null;
			let linksData: Link[] = [];
			
			try {
				profileData = await profileApi.getMyProfile($auth.token!);
			} catch (profileError: any) {
				console.warn('Profile not found, will be created automatically:', profileError);
			}
			
			try {
				linksData = await linksApi.getLinks($auth.token!);
			} catch (linksError: any) {
				console.warn('No links found:', linksError);
				linksData = [];
			}
			
			profile = profileData;
			links = linksData.sort((a, b) => a.position - b.position);
			
			// Minimum loading time for smooth UX
			const elapsed = Date.now() - startTime;
			if (elapsed < 500) {
				await new Promise(resolve => setTimeout(resolve, 500 - elapsed));
			}
		} catch (error: any) {
			console.error('Failed to load dashboard data:', error);
			toast.error(error.message || 'Failed to load data');
		} finally {
			initialLoading = false;
		}
	}

	async function handleAddLink() {
		if (!newLink.title || !newLink.url) {
			toast.error('Please fill in all fields');
			return;
		}

		// Validate URL format
		try {
			new URL(newLink.url);
		} catch {
			toast.error('Please enter a valid URL (e.g., https://example.com)');
			return;
		}
		
		try {
			const link = await linksApi.createLink({
				title: newLink.title,
				url: newLink.url,
				position: links.length,
				is_active: true
			}, $auth.token!);
			
			links = [...links, link];
			newLink = { title: '', url: '' };
			showAddLinkDialog = false;
			toast.success('Link added successfully!');
		} catch (error: any) {
			console.error('Add link error:', error);
			toast.error(error.message || 'Failed to add link');
		}
	}

	async function handleUpdateLink(event: CustomEvent) {
		const { id, title, url } = event.detail;
		try {
			const updated = await linksApi.updateLink(id, { title, url }, $auth.token!);
			links = links.map(link => link.id === id ? updated : link);
			toast.success('Link updated!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to update link');
		}
	}

	async function handleDeleteLink(event: CustomEvent) {
		const id = event.detail;
		try {
			await linksApi.deleteLink(id, $auth.token!);
			links = links.filter(link => link.id !== id);
			toast.success('Link deleted!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete link');
		}
	}

	async function handleToggleLink(event: CustomEvent) {
		const id = event.detail;
		const link = links.find(l => l.id === id);
		if (!link) return;
		
		try {
			const updated = await linksApi.updateLink(id, { is_active: !link.is_active }, $auth.token!);
			links = links.map(l => l.id === id ? updated : l);
		} catch (error: any) {
			toast.error(error.message || 'Failed to toggle link');
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
	<title>Links - LinkBio</title>
</svelte:head>

<style>
	@keyframes fade-in {
		from { opacity: 0; transform: translateY(10px); }
		to { opacity: 1; transform: translateY(0); }
	}
	
	.animate-in {
		animation: fade-in 0.3s ease-out;
	}

	/* Remove all borders when dragging */
	:global(section > div) {
		outline: none !important;
	}
	
	:global(section > div[style*="cursor: grabbing"]),
	:global(section > div[style*="cursor: grabbing"] *) {
		outline: none !important;
		border: none !important;
		transition: none !important;
	}

	/* Disable all transitions when dragging */
	section.dragging :global(*) {
		transition: none !important;
		animation: none !important;
	}
</style>



<div class="h-full bg-gray-50">
	<!-- Page Header -->
	<div class="bg-white/80 backdrop-blur-xl border-b border-gray-200/50 px-8 h-16 sticky top-0 z-10">
		<div class="flex items-center justify-between h-full">
			<div>
				<h1 class="text-xl font-bold text-gray-900">
					Links
				</h1>
				<p class="text-xs text-gray-500">Manage and organize your links</p>
			</div>
			<div class="flex items-center gap-3">
				<div class="flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 rounded-xl text-sm min-w-[280px]">
					<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
					</svg>
					<span class="font-medium text-gray-900 flex-1">
						linkbio.com/{profile?.username || 'loading...'}
					</span>
				</div>

				<Button
					variant="outline"
					size="icon"
					onclick={handleCopyUrl}
					class="h-10 w-10 border-gray-200 hover:bg-gray-50"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
					</svg>
				</Button>

				<Button
					variant="outline"
					size="icon"
					onclick={() => {
						if (profile?.username) {
							window.open(`${window.location.origin}/${profile.username}`, '_blank');
						}
					}}
					class="h-10 w-10 border-gray-200 hover:bg-gray-50"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
					</svg>
				</Button>

				<Button
					onclick={() => {
						if (profile?.username) {
							const url = `${window.location.origin}/${profile.username}`;
							if (navigator.share) {
								navigator.share({ title: 'My LinkBio', url }).catch(() => {});
							} else {
								handleCopyUrl();
							}
						}
					}}
					class="h-10 px-6 bg-gray-900 hover:bg-gray-800 text-white font-semibold"
				>
					SHARE
				</Button>
			</div>
		</div>
	</div>

	<!-- Main Content -->
	<div class="p-8 relative min-h-[calc(100vh-4rem)]">
		{#if initialLoading}
			<div class="absolute inset-0 flex items-center justify-center">
				<div class="text-center">
					<div class="relative w-14 h-14 mx-auto mb-4">
						<div class="absolute inset-0 border-4 border-indigo-100 rounded-full"></div>
						<div class="absolute inset-0 border-4 border-transparent border-t-indigo-600 border-r-blue-600 rounded-full animate-spin"></div>
					</div>
					<p class="text-sm text-gray-500">Loading...</p>
				</div>
			</div>
		{:else}
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8 animate-in">
			<!-- Links Section (2/3) -->
			<div class="lg:col-span-2 space-y-6">
				<!-- Add Link Button - Full Width -->
				<button
					onclick={() => showAddLinkDialog = true}
					class="w-full flex items-center justify-center gap-3 px-6 py-8 bg-white hover:bg-gray-50 border border-dashed border-gray-300 hover:border-indigo-400 rounded-2xl text-gray-600 hover:text-indigo-600 font-semibold transition-all group"
				>
					<svg class="w-6 h-6 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					<span class="text-lg">Add</span>
				</button>

				<!-- Action Buttons -->
				<div class="flex items-center justify-between mb-6 px-4">
					<button class="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
						</svg>
						<span class="font-medium">Add collection</span>
					</button>
					<button class="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors">
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"/>
						</svg>
						<span class="font-medium">View archive</span>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
						</svg>
					</button>
				</div>

				<!-- Links List -->
				<section 
					class="space-y-3"
					class:dragging={isDragging}
					use:dndzone={{ 
						items: links, 
						flipDurationMs: 200, 
						dropTargetStyle: {}, 
						dragDisabled: false,
						type: 'links'
					}}
					onconsider={(e) => { 
						isDragging = true;
						links = e.detail.items; 
					}}
					onfinalize={async (e) => {
						isDragging = false;
						links = e.detail.items;
						const linkIds = links.map(l => l.id);
						
						try {
							await linksApi.reorderLinks(linkIds, $auth.token!);
						} catch (error) {
							console.error('Failed to reorder:', error);
							toast.error('Failed to reorder links');
						}
					}}
				>
					{#each links as link (link.id)}
						<div>
							<LinkCard 
								{link} 
								on:update={handleUpdateLink}
								on:delete={handleDeleteLink}
								on:toggle={handleToggleLink}
							/>
						</div>
					{/each}

					{#if links.length === 0}
						<div class="relative bg-gradient-to-br from-indigo-50 via-white to-blue-50 rounded-2xl p-16 text-center border-2 border-dashed border-indigo-200 overflow-hidden">
							<div class="absolute top-0 right-0 w-32 h-32 bg-indigo-100 rounded-full -mr-16 -mt-16 opacity-50"></div>
							<div class="absolute bottom-0 left-0 w-24 h-24 bg-blue-100 rounded-full -ml-12 -mb-12 opacity-50"></div>
							
							<div class="relative">
								<div class="w-20 h-20 mx-auto mb-6 bg-gradient-to-br from-indigo-500 to-blue-700 rounded-2xl flex items-center justify-center shadow-xl shadow-indigo-500/30">
									<svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
									</svg>
								</div>
								<h3 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent mb-2">
									No links yet
								</h3>
								<p class="text-gray-600 mb-6 max-w-sm mx-auto">Get started by adding your first link and share it with the world</p>
								<Button 
									onclick={() => showAddLinkDialog = true}
									class="bg-gradient-to-r from-indigo-700 to-blue-700 hover:from-indigo-800 hover:to-blue-800 shadow-lg shadow-indigo-500/30"
								>
									<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
									</svg>
									Add Your First Link
								</Button>
							</div>
						</div>
					{/if}
				</section>
			</div>

			<!-- Preview Section (1/3) -->
			<div class="lg:col-span-1">
				<div class="sticky top-32 space-y-6">
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

					<!-- Branding Toggle -->
					<div class="flex items-center justify-center gap-3">
						<label class="relative inline-flex items-center cursor-pointer">
							<input type="checkbox" class="sr-only peer" disabled>
							<div class="w-11 h-6 bg-gray-200 rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all"></div>
						</label>
						<span class="text-sm font-medium text-gray-700">Hide Beacons logo and footer</span>
						<button class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-pink-500 to-rose-500 hover:from-pink-600 hover:to-rose-600 text-white font-bold rounded-full text-sm transition-all shadow-lg shadow-pink-500/30">
							<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
								<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
							</svg>
							Pro
						</button>
					</div>
				</div>
			</div>
		</div>
		{/if}
	</div>
</div>

<!-- Add Link Dialog -->
<Dialog.Root bind:open={showAddLinkDialog}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Add New Link</Dialog.Title>
			<Dialog.Description>
				Add a new link to your profile. It will be visible to your visitors.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4 py-4">
			<div>
				<Label for="link-title">Title</Label>
				<Input id="link-title" bind:value={newLink.title} placeholder="My Awesome Link" class="mt-2" />
			</div>
			<div>
				<Label for="link-url">URL</Label>
				<Input id="link-url" bind:value={newLink.url} placeholder="https://example.com" class="mt-2" />
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => showAddLinkDialog = false}>Cancel</Button>
			<Button onclick={handleAddLink} class="bg-gradient-to-r from-indigo-700 to-blue-700 hover:from-indigo-800 hover:to-blue-800 shadow-lg shadow-indigo-500/30">Add Link</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
