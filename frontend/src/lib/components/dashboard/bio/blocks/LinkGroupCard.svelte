<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { dndzone } from 'svelte-dnd-action';
	import type { Link } from '$lib/api/links';
	import LayoutSettingsSection from '../sections/LayoutSettingsSection.svelte';
	
	export let group: Link;
	export let selected: boolean = false;
	export let expanded: boolean = false;
	
	const dispatch = createEventDispatcher();
	
	let isExpanded = expanded;
	let showMenu = false;
	let activeTab: 'links' | 'layouts' = 'links';
	let linkMenuOpen: string | null = null; // Track which link's menu is open
	let groupMenuOpen = false; // Track if group menu is open
	let fileInput: HTMLInputElement;
	let uploadingLinkId: string | null = null;
	
	// Sync isExpanded with prop
	$: isExpanded = expanded;
	
	$: childrenCount = group.children?.length || 0;
	$: children = group.children || [];
	
	// Drag & drop items
	type DndItem = { id: string; data: Link };
	let dndItems: DndItem[] = [];
	
	// Sync dndItems from children - keep user's drag order, don't auto-sort by pin status
	$: {
		if (dndItems.length === 0 || dndItems.length !== children.length) {
			// Initial load or children count changed (add/remove)
			// Sort by position only (pinned status doesn't affect edit view order)
			dndItems = children
				.slice()
				.sort((a, b) => a.position - b.position)
				.map(link => ({ id: link.id, data: link }));
		} else {
			// Update data for existing items without changing order
			dndItems = dndItems.map(item => {
				const updated = children.find(c => c.id === item.id);
				return updated ? { ...item, data: updated } : item;
			});
		}
	}
	
	function handleThumbnailClick(linkId: string) {
		uploadingLinkId = linkId;
		fileInput?.click();
	}
	
	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (file && uploadingLinkId) {
			dispatch('uploadThumbnail', { linkId: uploadingLinkId, file });
			uploadingLinkId = null;
		}
	}
	
	function handleDndConsider(e: CustomEvent) {
		const newItems = e.detail.items;
		
		// Enforce constraint: unpinned items cannot be dragged above pinned items
		const pinnedItems = newItems.filter((item: DndItem) => item.data.is_pinned);
		const unpinnedItems = newItems.filter((item: DndItem) => !item.data.is_pinned);
		
		// Find the first unpinned item index in newItems
		const firstUnpinnedIndex = newItems.findIndex((item: DndItem) => !item.data.is_pinned);
		
		// Check if any pinned item appears after unpinned items (violation)
		let hasViolation = false;
		if (pinnedItems.length > 0 && firstUnpinnedIndex !== -1) {
			for (let i = firstUnpinnedIndex; i < newItems.length; i++) {
				if (newItems[i].data.is_pinned) {
					hasViolation = true;
					break;
				}
			}
		}
		
		// If constraint violated, enforce correct order
		if (hasViolation) {
			dndItems = [...pinnedItems, ...unpinnedItems];
		} else {
			dndItems = newItems;
		}
	}
	
	function handleDndFinalize(e: CustomEvent) {
		const newItems = e.detail.items;
		
		// Enforce final constraint: pinned first, then unpinned
		const pinnedItems = newItems.filter((item: DndItem) => item.data.is_pinned);
		const unpinnedItems = newItems.filter((item: DndItem) => !item.data.is_pinned);
		
		// Always maintain correct order: pinned first, then unpinned
		dndItems = [...pinnedItems, ...unpinnedItems];
		
		// Dispatch event to parent to save new order
		const linkIds = dndItems.map(item => item.id);
		dispatch('reorderlinks', { groupId: group.id, linkIds });
	}
	
	// Click outside handler for link menu
	function clickOutsideLinkMenu(node: HTMLElement) {
		const handleClick = (event: MouseEvent) => {
			const target = event.target as Node;
			// Check if click is outside the menu and not on a toggle switch
			if (node && !node.contains(target) && !event.defaultPrevented) {
				// Don't close if clicking on any button inside the menu
				const isInsideMenu = (target as Element).closest?.('.absolute.right-0.top-full');
				if (!isInsideMenu) {
					linkMenuOpen = null;
				}
			}
		};
		
		document.addEventListener('mousedown', handleClick, true);
		
		return {
			destroy() {
				document.removeEventListener('mousedown', handleClick, true);
			}
		};
	}
	
	// Click outside handler for group menu
	function clickOutsideGroupMenu(node: HTMLElement) {
		const handleClick = (event: MouseEvent) => {
			const target = event.target as Node;
			if (node && !node.contains(target) && !event.defaultPrevented) {
				groupMenuOpen = false;
			}
		};
		
		document.addEventListener('mousedown', handleClick, true);
		
		return {
			destroy() {
				document.removeEventListener('mousedown', handleClick, true);
			}
		};
	}
</script>

{#if !isExpanded}
	<!-- Collapsed View - Linktree Style -->
	<div 
		class="group/card bg-white rounded-xl transition-all hover:shadow-md"
		class:ring-2={selected}
		class:ring-indigo-500={selected}
		class:shadow-sm={!selected}
	>
		<div class="flex items-center gap-4 px-4 py-3.5">
			<!-- Thumbnail/Icon -->
			<button
				onclick={() => dispatch('expand', group.id)}
				class="flex-shrink-0 w-14 h-14 rounded-lg bg-gradient-to-br from-emerald-400 to-emerald-600 flex items-center justify-center shadow-sm hover:shadow-md transition-shadow"
			>
				<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
				</svg>
			</button>

			<!-- Content - Click to expand -->
			<button
				onclick={() => dispatch('expand', group.id)}
				class="flex-1 text-left min-w-0"
			>
				<h3 class="text-base font-semibold text-gray-900 truncate">
					{group.group_title || 'Links'}
				</h3>
				<p class="text-sm text-gray-500 mt-0.5 truncate">
					{childrenCount} {childrenCount === 1 ? 'link' : 'links'}
				</p>
			</button>

			<!-- Stats & Actions (Linktree Style) -->
			<div class="flex items-center gap-3">
				<!-- Toggle Visibility Button -->
				<button
					onclick={(e) => { 
						e.stopPropagation();
						dispatch('toggle', group.id);
					}}
					class="p-1.5 hover:bg-gray-100 rounded-lg transition-colors"
					class:text-gray-400={group.is_active}
					class:text-gray-300={!group.is_active}
					title={group.is_active ? 'Hide group' : 'Show group'}
				>
					{#if group.is_active}
						<!-- Eye Open - Visible -->
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
						</svg>
					{:else}
						<!-- Eye Closed - Hidden -->
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
						</svg>
					{/if}
				</button>

				<!-- 3-dot Menu Button -->
				<div class="relative" use:clickOutsideGroupMenu>
					<button
						onclick={(e) => { 
							e.stopPropagation();
							groupMenuOpen = !groupMenuOpen;
						}}
						class="p-1.5 hover:bg-gray-100 rounded-lg transition-colors text-gray-400 hover:text-gray-600"
						title="More options"
					>
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<circle cx="10" cy="4" r="1.5"/>
							<circle cx="10" cy="10" r="1.5"/>
							<circle cx="10" cy="16" r="1.5"/>
						</svg>
					</button>

					{#if groupMenuOpen}
						<div class="absolute right-0 top-full mt-1 w-56 bg-white rounded-xl shadow-xl border border-gray-200 py-2 z-20">
							<!-- Edit Group -->
							<button
								onclick={(e) => { 
									e.stopPropagation();
									dispatch('expand', group.id); 
									groupMenuOpen = false;
								}}
								class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-blue-50 flex items-center gap-3 transition-all group"
							>
								<svg class="w-4 h-4 text-blue-500 group-hover:text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
								</svg>
								<span class="font-medium">Edit group</span>
							</button>

							<!-- Duplicate Group -->
							<button
								onclick={(e) => { 
									e.stopPropagation();
									dispatch('duplicate', group.id); 
									groupMenuOpen = false;
								}}
								class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-green-50 flex items-center gap-3 transition-all group"
							>
								<svg class="w-4 h-4 text-green-500 group-hover:text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
								</svg>
								<span class="font-medium">Duplicate group</span>
							</button>

							<div class="border-t border-gray-100 my-1"></div>

							<!-- Delete Group -->
							<button
								onclick={(e) => { 
									e.stopPropagation();
									dispatch('delete', group.id); 
									groupMenuOpen = false;
								}}
								class="w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50 flex items-center gap-3 transition-all group"
							>
								<svg class="w-4 h-4 text-red-500 group-hover:text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
								</svg>
								<span class="font-medium">Delete group</span>
							</button>
						</div>
					{/if}
				</div>

				<!-- Drag Handle -->
				<button
					class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400 p-1 transition-colors"
					onclick={(e) => e.stopPropagation()}
				>
					<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
						<circle cx="9" cy="7" r="1.5"/>
						<circle cx="9" cy="12" r="1.5"/>
						<circle cx="9" cy="17" r="1.5"/>
						<circle cx="15" cy="7" r="1.5"/>
						<circle cx="15" cy="12" r="1.5"/>
						<circle cx="15" cy="17" r="1.5"/>
					</svg>
				</button>
			</div>
		</div>
	</div>
{:else}
	<!-- Expanded View (Inline Full UI) -->
	<div class="bg-white rounded-xl shadow-lg">
		<!-- Header -->
		<div class="border-b border-gray-100 px-6 py-4">
			<div class="flex items-center justify-between mb-4">
				<!-- Left side: Back button + Icon + Title -->
				<div class="flex items-center gap-3">
					<button
						onclick={() => dispatch('collapse', group.id)}
						class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
						title="Back"
					>
						<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
						</svg>
					</button>
					
					<div class="w-10 h-10 rounded-lg bg-emerald-100 flex items-center justify-center flex-shrink-0">
						<svg class="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
						</svg>
					</div>
					
					<h2 class="text-xl font-bold text-gray-900">{group.group_title || 'Links'}</h2>
				</div>

				<!-- Right side: Action buttons -->
				<div class="flex items-center gap-3">
					<button class="p-2 hover:bg-gray-100 rounded-lg transition-colors" title="Help">
						<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
						</svg>
					</button>
					<button 
						onclick={() => dispatch('duplicate', group.id)}
						class="p-2 hover:bg-gray-100 rounded-lg transition-colors" 
						title="Duplicate"
					>
						<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
						</svg>
					</button>
					<button 
						onclick={() => dispatch('delete', group.id)}
						class="p-2 hover:bg-gray-100 rounded-lg transition-colors" 
						title="Delete"
					>
						<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
						</svg>
					</button>
					<button 
						onclick={() => dispatch('toggle', group.id)}
						class="p-2 hover:bg-gray-100 rounded-lg transition-colors" 
						title={group.is_active ? 'Hide group' : 'Show group'}
					>
						{#if group.is_active}
							<!-- Eye Open - Visible -->
							<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
							</svg>
						{:else}
							<!-- Eye Closed - Hidden -->
							<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
							</svg>
						{/if}
					</button>
				</div>
			</div>

			<!-- Tabs -->
			<div class="flex gap-6">
				<button
					onclick={() => activeTab = 'links'}
					class="pb-3 text-sm font-semibold transition-colors relative"
					class:text-gray-900={activeTab === 'links'}
					class:text-gray-500={activeTab !== 'links'}
				>
					Links
					{#if activeTab === 'links'}
						<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-gray-900"></div>
					{/if}
				</button>
				<button
					onclick={() => activeTab = 'layouts'}
					class="pb-3 text-sm font-semibold transition-colors relative"
					class:text-gray-900={activeTab === 'layouts'}
					class:text-gray-500={activeTab !== 'layouts'}
				>
					Layouts
					{#if activeTab === 'layouts'}
						<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-gray-900"></div>
					{/if}
				</button>
			</div>
		</div>

		<!-- Content -->
		<div class="p-6">
			{#if activeTab === 'links'}
				<!-- Add Link Button -->
				<button
					onclick={() => dispatch('addlink', group.id)}
					class="w-full px-4 py-3 text-gray-700 hover:text-gray-900 hover:bg-gray-50 transition-colors flex items-center justify-center gap-2 text-sm font-medium mb-6"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
					</svg>
					Add link
				</button>

				<!-- Links List - Linktree Style -->
				<div 
					class="space-y-3"
					use:dndzone={{items: dndItems, flipDurationMs: 200, dropTargetStyle: {}}}
					onconsider={handleDndConsider}
					onfinalize={handleDndFinalize}
				>
					{#each dndItems as item (item.id)}
						{@const link = item.data}
						<div class="bg-gradient-to-r from-gray-50 to-gray-100/50 rounded-xl shadow-sm hover:shadow-md transition-all border border-gray-200/60">
							<div class="flex items-center gap-4 p-3.5">
								<!-- Drag Handle -->
								<button
									class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400 p-1 transition-colors flex-shrink-0"
									onclick={(e) => e.stopPropagation()}
								>
									<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
										<circle cx="9" cy="7" r="1.5"/>
										<circle cx="9" cy="12" r="1.5"/>
										<circle cx="9" cy="17" r="1.5"/>
										<circle cx="15" cy="7" r="1.5"/>
										<circle cx="15" cy="12" r="1.5"/>
										<circle cx="15" cy="17" r="1.5"/>
									</svg>
								</button>

								<!-- Thumbnail - Clickable to upload -->
								<button
									onclick={() => handleThumbnailClick(link.id)}
									class="flex-shrink-0 w-14 h-14 rounded-lg overflow-hidden hover:opacity-80 transition-opacity group/thumb relative"
									title="Click to upload thumbnail"
								>
									{#if link.thumbnail_url}
										<img 
											src={link.thumbnail_url} 
											alt={link.title}
											class="w-full h-full object-cover"
										/>
										<div class="absolute inset-0 bg-black/50 opacity-0 group-hover/thumb:opacity-100 transition-opacity flex items-center justify-center">
											<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
											</svg>
										</div>
									{:else}
										<div class="w-full h-full rounded-lg bg-gray-200 flex items-center justify-center">
											<svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
											</svg>
										</div>
									{/if}
								</button>

							<!-- Content - Click to edit -->
							<button
								onclick={() => dispatch('editlink', link)}
								class="flex-1 text-left min-w-0"
							>
								<h3 class="text-base font-semibold text-gray-900 truncate">{link.title}</h3>
								{#if link.description}
									<p class="text-sm text-gray-600 mt-0.5 line-clamp-2">{link.description}</p>
								{/if}
								<p class="text-xs text-gray-400 mt-0.5 truncate">{link.url}</p>
							</button>								<!-- Stats & Actions (Linktree Style) -->
								<div class="flex items-center gap-2">
									<!-- Toggle Visibility Button -->
									<button
										onclick={(e) => { 
											e.stopPropagation();
											dispatch('togglelinkvisibility', link.id);
										}}
										class="p-1.5 hover:bg-gray-100 rounded-lg transition-colors"
										class:text-gray-400={link.is_active}
										class:text-gray-300={!link.is_active}
										title={link.is_active ? 'Hide link' : 'Show link'}
									>
										{#if link.is_active}
											<!-- Eye Open - Visible -->
											<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
											</svg>
										{:else}
											<!-- Eye Closed - Hidden -->
											<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
											</svg>
										{/if}
									</button>

									<!-- Pin Link Button -->
									<button
										onclick={(e) => { 
											e.stopPropagation();
											dispatch('pinlink', link.id);
										}}
										class="p-1.5 hover:bg-gray-100 rounded-lg transition-colors"
										class:text-amber-500={link.is_pinned}
										class:text-gray-400={!link.is_pinned}
										title={link.is_pinned ? 'Unpin link' : 'Pin link'}
									>
										<svg class="w-5 h-5" fill={link.is_pinned ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
											<path stroke-linecap="round" stroke-linejoin="round" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
										</svg>
									</button>

									<!-- More Menu (3 dots) -->
									<div class="relative" use:clickOutsideLinkMenu>
										<button
											onclick={(e) => { 
												e.stopPropagation(); 
												linkMenuOpen = linkMenuOpen === link.id ? null : link.id;
											}}
											class="p-1.5 hover:bg-gray-100 rounded-lg transition-colors text-gray-400 hover:text-gray-600"
											title="More options"
										>
											<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
												<circle cx="10" cy="4" r="1.5"/>
												<circle cx="10" cy="10" r="1.5"/>
												<circle cx="10" cy="16" r="1.5"/>
											</svg>
										</button>

										{#if linkMenuOpen === link.id}
											<div class="absolute right-0 top-full mt-1 w-56 bg-white rounded-xl shadow-xl border border-gray-200 py-2 z-20">
												<!-- Edit Link -->
												<button
													onclick={(e) => { 
														e.stopPropagation();
														dispatch('editlink', link); 
														linkMenuOpen = null; 
													}}
													class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-blue-50 flex items-center gap-3 transition-all group"
												>
													<svg class="w-4 h-4 text-blue-500 group-hover:text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
													</svg>
													<span class="font-medium">Edit link</span>
												</button>

												<!-- Duplicate Link -->
												<button
													onclick={(e) => { 
														e.stopPropagation();
														dispatch('duplicatelink', link.id); 
														linkMenuOpen = null; 
													}}
													class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-green-50 flex items-center gap-3 transition-all group"
												>
													<svg class="w-4 h-4 text-green-500 group-hover:text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
													</svg>
													<span class="font-medium">Duplicate link</span>
												</button>

												<div class="border-t border-gray-100 my-1"></div>

												<!-- Open in New Tab Toggle -->
												<div 
													class="px-4 py-2.5"
													onclick={(e) => e.stopPropagation()}
												>
													<label class="flex items-center justify-between cursor-pointer group">
														<div class="flex items-center gap-3">
															<svg class="w-4 h-4 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
																<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
															</svg>
															<span class="text-sm font-medium text-gray-700">Open in new tab</span>
														</div>
														<!-- Toggle Switch -->
														<button
															type="button"
															onclick={(e) => {
																e.stopPropagation();
																e.preventDefault();
																console.log('🔄 Toggle clicked for link:', link.id, 'Current state:', link.open_in_new_tab);
																dispatch('toggleNewTab', link.id);
															}}
															class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none"
															class:bg-indigo-600={link.open_in_new_tab}
															class:bg-gray-200={!link.open_in_new_tab}
															role="switch"
															aria-checked={link.open_in_new_tab}
														>
															<span
																class="inline-block h-4 w-4 transform rounded-full bg-white shadow-sm transition-transform"
																class:translate-x-5={link.open_in_new_tab}
																class:translate-x-0.5={!link.open_in_new_tab}
															></span>
														</button>
													</label>
												</div>

												<div class="border-t border-gray-100 my-1"></div>

												<!-- Delete Link -->
												<button
													onclick={(e) => { 
														e.stopPropagation();
														dispatch('removefromgroup', link.id); 
														linkMenuOpen = null; 
													}}
													class="w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50 flex items-center gap-3 transition-all group"
												>
													<svg class="w-4 h-4 text-red-500 group-hover:text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
													</svg>
													<span class="font-medium">Delete link</span>
												</button>
											</div>
										{/if}
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{:else if activeTab === 'layouts'}
				<!-- Layouts Tab -->
				<LayoutSettingsSection 
					{group}
					on:update={(e) => dispatch('updatelayout', e.detail)}
				/>
			{/if}
		</div>
	</div>
{/if}

<!-- Hidden file input for thumbnail upload -->
<input
	bind:this={fileInput}
	type="file"
	accept="image/*"
	onchange={handleFileSelect}
	class="hidden"
/>

<!-- Click outside to close menus -->
{#if showMenu || linkMenuOpen || groupMenuOpen}
	<button
		class="fixed inset-0 z-0"
		onclick={() => { showMenu = false; linkMenuOpen = null; groupMenuOpen = false; }}
		aria-label="Close menu"
	></button>
{/if}

<style>
	/* Disable transitions during drag */
	:global(.space-y-2),
	:global(.space-y-2 *) {
		transition: none !important;
	}
	
	/* Keep opacity for dragging item */
	:global([aria-grabbed="true"]) {
		opacity: 0.5;
		cursor: grabbing !important;
	}
</style>
