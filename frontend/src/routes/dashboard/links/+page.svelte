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
	import ImageUploader from '$lib/components/ui/ImageUploader.svelte';
	import LinksToolbar from '$lib/components/dashboard/LinksToolbar.svelte';
	import { profileApi } from '$lib/api/profile';
	import { linksApi, type LinkFilters } from '$lib/api/links';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';

	let profile: Profile | null = null;
	let links: Link[] = [];
	let allLinks: Link[] = []; // Store unfiltered links
	let initialLoading = true;
	let selectedIds = new Set<string>();

	// Search & Filter state
	let searchQuery = '';
	let statusFilter: 'active' | 'inactive' | '' = '';
	let layoutFilter: 'classic' | 'featured' | '' = '';
	let sortBy: 'position' | 'clicks' | 'created' | 'updated' | 'title' | '' = '';
	let showFilters = false;

	$: selectedCount = selectedIds.size;
	$: hasActiveFilters = searchQuery || statusFilter || layoutFilter || sortBy;

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
	
	let showThumbnailDialog = false;
	let editingLinkId: string | null = null;
	let uploadingThumbnail = false;

	onMount(async () => {
		console.log('🎬 Component mounted');
		await loadData();
		// Set initialized after data loaded
		setTimeout(() => {
			console.log('✅ Setting isInitialized to true');
			isInitialized = true;
		}, 100);
	});

	async function loadData() {
		try {
			initialLoading = true;
			
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
			links = linksData;
			allLinks = linksData; // Store original unfiltered list
		} catch (error: any) {
			console.error('Failed to load dashboard data:', error);
			toast.error(error.message || 'Failed to load data');
		} finally {
			initialLoading = false;
		}
	}

	async function applyFilters() {
		try {
			const filters: LinkFilters = {
				search: searchQuery || undefined,
				status: statusFilter || undefined,
				layout_type: layoutFilter || undefined,
				sort_by: sortBy || undefined
			};
			
			console.log('🔍 Applying filters:', filters);
			console.log('📊 Current links count before filter:', links.length);
			
			links = await linksApi.getLinks($auth.token!, filters);
			
			console.log('✅ Links after filter:', links.length, links);
		} catch (error: any) {
			console.error('❌ Failed to apply filters:', error);
			toast.error('Failed to apply filters');
		}
	}

	function clearFilters() {
		searchQuery = '';
		statusFilter = '';
		layoutFilter = '';
		sortBy = '';
		links = allLinks;
	}

	// Debounce search
	let searchTimeout: any;
	let isInitialized = false;
	
	$: {
		console.log('🔄 Search query changed:', searchQuery, 'isInitialized:', isInitialized);
		if (isInitialized && searchQuery !== undefined) {
			clearTimeout(searchTimeout);
			console.log('⏱️ Setting timeout for search...');
			searchTimeout = setTimeout(() => {
				console.log('🚀 Timeout fired, calling applyFilters');
				applyFilters();
			}, 300);
		}
	}

	$: {
		console.log('🔄 Filters changed - status:', statusFilter, 'layout:', layoutFilter, 'sort:', sortBy);
		if (isInitialized && (statusFilter !== undefined || layoutFilter !== undefined || sortBy !== undefined)) {
			console.log('🚀 Calling applyFilters from filter change');
			applyFilters();
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
			await linksApi.createLink({
				title: newLink.title,
				url: newLink.url,
				position: allLinks.length,
				is_active: true
			}, $auth.token!);
			
			// Reload data to get fresh list
			await loadData();
			if (hasActiveFilters) {
				await applyFilters();
			}
			
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
			allLinks = allLinks.map(link => link.id === id ? updated : link);
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
			allLinks = allLinks.filter(link => link.id !== id);
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
			allLinks = allLinks.map(l => l.id === id ? updated : l);
		} catch (error: any) {
			toast.error(error.message || 'Failed to toggle link');
		}
	}

	function handleEditThumbnail(event: CustomEvent) {
		editingLinkId = event.detail;
		showThumbnailDialog = true;
	}

	async function handleUploadThumbnail(event: CustomEvent) {
		if (!editingLinkId) return;
		
		const file = event.detail as File;
		uploadingThumbnail = true;
		
		try {
			const updated = await linksApi.uploadThumbnail(editingLinkId, file, $auth.token!);
			links = links.map(l => l.id === editingLinkId ? updated : l);
			allLinks = allLinks.map(l => l.id === editingLinkId ? updated : l);
			toast.success('Thumbnail uploaded!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to upload thumbnail');
		} finally {
			uploadingThumbnail = false;
		}
	}

	async function handleRemoveThumbnail() {
		if (!editingLinkId) return;
		
		try {
			const updated: Link = await linksApi.deleteThumbnail(editingLinkId, $auth.token!);
			links = links.map(l => l.id === editingLinkId ? updated : l);
			allLinks = allLinks.map(l => l.id === editingLinkId ? updated : l);
			toast.success('Thumbnail removed!');
			showThumbnailDialog = false;
		} catch (error: any) {
			toast.error(error.message || 'Failed to remove thumbnail');
		}
	}

	async function handleUpdateLayout(event: CustomEvent) {
		const { id, layout_type } = event.detail;
		try {
			const updated = await linksApi.updateLink(id, { layout_type }, $auth.token!);
			links = links.map(l => l.id === id ? updated : l);
			allLinks = allLinks.map(l => l.id === id ? updated : l);
			toast.success(`Layout changed to ${layout_type}!`);
		} catch (error: any) {
			toast.error(error.message || 'Failed to update layout');
		}
	}

	async function handleDuplicate(event: CustomEvent) {
		const id = event.detail;
		try {
			await linksApi.duplicateLink(id, $auth.token!);
			await loadData();
			if (hasActiveFilters) {
				await applyFilters();
			}
			toast.success('Link duplicated successfully!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to duplicate link');
		}
	}

	function handleSelect(event: CustomEvent) {
		const id = event.detail;
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = selectedIds;
	}

	function selectAll() {
		selectedIds = new Set(links.map(l => l.id));
	}

	function clearSelection() {
		selectedIds.clear();
		selectedIds = selectedIds;
	}

	async function bulkDelete() {
		if (selectedIds.size === 0) return;
		
		try {
			await linksApi.bulkAction(Array.from(selectedIds), 'delete', $auth.token!);
			links = links.filter(l => !selectedIds.has(l.id));
			allLinks = allLinks.filter(l => !selectedIds.has(l.id));
			toast.success(`Deleted ${selectedIds.size} links`);
			clearSelection();
		} catch (error: any) {
			toast.error(error.message || 'Failed to delete links');
		}
	}

	async function bulkActivate() {
		if (selectedIds.size === 0) return;
		
		try {
			await linksApi.bulkAction(Array.from(selectedIds), 'activate', $auth.token!);
			links = links.map(l => selectedIds.has(l.id) ? {...l, is_active: true} : l);
			allLinks = allLinks.map(l => selectedIds.has(l.id) ? {...l, is_active: true} : l);
			toast.success(`Activated ${selectedIds.size} links`);
			clearSelection();
		} catch (error: any) {
			toast.error(error.message || 'Failed to activate links');
		}
	}

	async function bulkDeactivate() {
		if (selectedIds.size === 0) return;
		
		try {
			await linksApi.bulkAction(Array.from(selectedIds), 'deactivate', $auth.token!);
			links = links.map(l => selectedIds.has(l.id) ? {...l, is_active: false} : l);
			allLinks = allLinks.map(l => selectedIds.has(l.id) ? {...l, is_active: false} : l);
			toast.success(`Deactivated ${selectedIds.size} links`);
			clearSelection();
		} catch (error: any) {
			toast.error(error.message || 'Failed to deactivate links');
		}
	}

	async function handlePin(event: CustomEvent) {
		const id = event.detail;
		try {
			const updated = await linksApi.togglePin(id, $auth.token!);
			// Reload to get correct order
			await loadData();
			toast.success(updated.is_pinned ? 'Link pinned!' : 'Link unpinned!');
		} catch (error: any) {
			toast.error(error.message || 'Failed to pin link');
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

	/* Remove all borders and outlines when dragging */
	:global([aria-grabbed="true"]),
	:global([aria-grabbed="true"] *),
	:global(section div:focus),
	:global(section div:focus-visible),
	:global(section div[tabindex]:focus),
	:global(section div[tabindex]:focus-visible) {
		outline: none !important;
		border-color: transparent !important;
		box-shadow: none !important;
	}

	/* Ensure no focus styles on wrapper divs */
	:global(section > div) {
		outline: none !important;
	}

	/* Tắt hoàn toàn mọi transition/animation trong drag zone */
	:global(section.space-y-3),
	:global(section.space-y-3 *) {
		transition: none !important;
		animation: none !important;
	}

	/* Chỉ giữ opacity cho item đang drag */
	:global([aria-grabbed="true"]) {
		opacity: 0.5;
		cursor: grabbing !important;
	}
</style>

<div class="h-full bg-gray-50">
	<!-- Page Header -->
	<div class="bg-white/80 backdrop-blur-xl border-b border-gray-200/50 px-8 h-16 sticky top-0 z-10">
		<div class="flex items-center justify-between h-full">
			<div>
				<h1 class="text-xl font-bold text-gray-900">Links</h1>
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
			<div class="absolute inset-0 flex items-center justify-center bg-gray-50">
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
				<!-- Compact Toolbar -->
				<LinksToolbar
					bind:searchQuery
					bind:showFilters
					hasActiveFilters={[searchQuery, statusFilter, layoutFilter, sortBy].filter(Boolean).length}
					linksCount={links.length}
					on:toggleFilters={() => showFilters = !showFilters}
					on:selectAll={selectAll}
					on:addCollection={() => toast.info('Add collection feature coming soon!')}
					on:viewArchive={() => toast.info('View archive feature coming soon!')}
				/>

				<!-- Filter Options (when expanded) -->
				{#if showFilters}
					<div class="bg-white rounded-xl border border-gray-200 p-4 shadow-sm animate-in">
						<div class="flex items-center gap-3 flex-wrap">
							<!-- Status Filter -->
							<select
								bind:value={statusFilter}
								class="flex-1 sm:flex-none px-3 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all cursor-pointer"
							>
								<option value="">All Status</option>
								<option value="active">✓ Active</option>
								<option value="inactive">✕ Inactive</option>
							</select>

							<!-- Layout Filter -->
							<select
								bind:value={layoutFilter}
								class="flex-1 sm:flex-none px-3 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all cursor-pointer"
							>
								<option value="">All Layouts</option>
								<option value="classic">Classic</option>
								<option value="featured">Featured</option>
							</select>

							<!-- Sort By -->
							<select
								bind:value={sortBy}
								class="flex-1 sm:flex-none px-3 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all cursor-pointer"
							>
								<option value="">Default Order</option>
								<option value="clicks">🔥 Most Clicks</option>
								<option value="created">🆕 Newest First</option>
								<option value="updated">⏱️ Recently Updated</option>
								<option value="title">🔤 Alphabetical</option>
							</select>

							<!-- Clear Filters -->
							{#if hasActiveFilters}
								<button
									onclick={clearFilters}
									class="ml-auto flex items-center gap-1.5 px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors"
								>
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
									</svg>
									Clear All
								</button>
							{/if}
						</div>
					</div>
				{/if}

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

				<!-- Drag Disabled Notice -->
				{#if hasActiveFilters && links.length > 0}
					<div class="bg-amber-50 border border-amber-200 rounded-xl p-3 flex items-center gap-3 text-sm">
						<svg class="w-5 h-5 text-amber-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
						</svg>
						<span class="text-amber-800">
							<strong>Drag & drop disabled</strong> while filters are active. Clear filters to reorder links.
						</span>
					</div>
				{/if}

				<!-- Links List -->
				<section 
					class="space-y-3"
					use:dndzone={{ 
						items: links,
						dragDisabled: hasActiveFilters,
						dropTargetStyle: {},
						type: 'links',
						flipDurationMs: 0
					}}
					onconsider={(e) => {
						links = e.detail.items;
					}}
					onfinalize={async (e) => {
						const newOrder = e.detail.items;
						const oldIds = links.map(l => l.id).join(',');
						const newIds = newOrder.map(l => l.id).join(',');
						const hasChanged = oldIds !== newIds;
						
						links = newOrder;
						allLinks = newOrder;
						
						if (hasChanged) {
							try {
								await linksApi.reorderLinks(newOrder.map(l => l.id), $auth.token!);
							} catch (error) {
								console.error('Failed to reorder:', error);
								toast.error('Failed to reorder links');
							}
						}
					}}
				>
					{#each links as link (link.id)}
						<div style="outline: none;">
							<LinkCard 
								{link}
								selected={selectedIds.has(link.id)}
								on:update={handleUpdateLink}
								on:delete={handleDeleteLink}
								on:toggle={handleToggleLink}
								on:editThumbnail={handleEditThumbnail}
								on:updateLayout={handleUpdateLayout}
								on:duplicate={handleDuplicate}
								on:select={handleSelect}
								on:pin={handlePin}
							/>
						</div>
					{/each}
				</section>

				{#if links.length === 0}
					{#if hasActiveFilters}
						<!-- No Results from Filter -->
						<div class="relative bg-white rounded-2xl p-12 text-center border border-gray-200">
							<div class="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
								<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
								</svg>
							</div>
							<h3 class="text-lg font-bold text-gray-900 mb-2">No links found</h3>
							<p class="text-gray-600 mb-4">Try adjusting your filters or search query</p>
							<Button 
								onclick={clearFilters}
								variant="outline"
								class="border-gray-300"
							>
								Clear Filters
							</Button>
						</div>
					{:else}
						<!-- No Links at All -->
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
				{/if}
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

<!-- Floating Bulk Action Bar -->
{#if selectedCount > 0}
	<div class="fixed bottom-8 left-1/2 -translate-x-1/2 z-50 animate-in">
		<div class="bg-gray-900 text-white rounded-2xl shadow-2xl px-6 py-4 flex items-center gap-6">
			<div class="flex items-center gap-2">
				<span class="font-semibold">{selectedCount}</span>
				<span class="text-gray-300">selected</span>
			</div>
			
			<div class="w-px h-6 bg-gray-700"></div>
			
			<div class="flex items-center gap-2">
				<button
					onclick={bulkActivate}
					class="px-4 py-2 bg-green-600 hover:bg-green-700 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
					</svg>
					Show
				</button>
				
				<button
					onclick={bulkDeactivate}
					class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
					</svg>
					Hide
				</button>
				
				<button
					onclick={bulkDelete}
					class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
					</svg>
					Delete
				</button>
			</div>
			
			<div class="w-px h-6 bg-gray-700"></div>
			
			<button
				onclick={clearSelection}
				class="p-2 hover:bg-gray-800 rounded-lg transition-colors"
				title="Clear selection"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
	</div>
{/if}

<!-- Thumbnail Upload Dialog -->
<Dialog.Root bind:open={showThumbnailDialog}>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Link Thumbnail</Dialog.Title>
			<Dialog.Description>
				Upload an image to make your link stand out. Recommended size: 400x400px
			</Dialog.Description>
		</Dialog.Header>
		<div class="py-4">
			<ImageUploader 
				currentImage={editingLinkId ? links.find(l => l.id === editingLinkId)?.thumbnail_url : null}
				uploading={uploadingThumbnail}
				on:upload={handleUploadThumbnail}
				on:remove={handleRemoveThumbnail}
			/>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => showThumbnailDialog = false}>Close</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
