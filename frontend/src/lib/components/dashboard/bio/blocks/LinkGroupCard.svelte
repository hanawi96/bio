<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { dndzone } from 'svelte-dnd-action';
	import type { Link } from '$lib/api/links';
	
	export let group: Link;
	export let selected: boolean = false;
	
	const dispatch = createEventDispatcher();
	
	let isExpanded = false;
	let showMenu = false;
	let activeTab: 'links' | 'layouts' | 'settings' = 'links';
	let linkMenuOpen: string | null = null; // Track which link's menu is open
	let fileInput: HTMLInputElement;
	let uploadingLinkId: string | null = null;
	
	$: childrenCount = group.children?.length || 0;
	$: children = group.children || [];
	
	// Drag & drop items
	type DndItem = { id: string; data: Link };
	let dndItems: DndItem[] = [];
	
	// Sync dndItems from children, but update data without changing order during drag
	$: {
		if (dndItems.length === 0 || dndItems.length !== children.length) {
			// Initial load or children count changed (add/remove)
			dndItems = children
				.slice()
				.sort((a, b) => {
					// Pinned items first, then by position
					if (a.is_pinned && !b.is_pinned) return -1;
					if (!a.is_pinned && b.is_pinned) return 1;
					return a.position - b.position;
				})
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
		dndItems = e.detail.items;
	}
	
	function handleDndFinalize(e: CustomEvent) {
		const newItems = e.detail.items;
		
		// Separate pinned and unpinned items
		const pinnedItems = newItems.filter((item: DndItem) => item.data.is_pinned);
		const unpinnedItems = newItems.filter((item: DndItem) => !item.data.is_pinned);
		
		// Reconstruct order: pinned first, then unpinned
		dndItems = [...pinnedItems, ...unpinnedItems];
		
		// Dispatch event to parent to save new order
		const linkIds = dndItems.map(item => item.id);
		dispatch('reorderlinks', { groupId: group.id, linkIds });
	}
</script>

{#if !isExpanded}
	<!-- Collapsed View -->
	<div 
		class="group/card bg-white rounded-xl transition-all shadow-sm"
		class:ring-2={selected}
		class:ring-blue-500={selected}
	>
		<div class="flex items-center gap-3 px-4 py-4">
			<!-- Checkbox -->
			<label class="flex items-center cursor-pointer">
				<input 
					type="checkbox" 
					checked={selected}
					onchange={() => dispatch('select', group.id)}
					class="w-4 h-4 rounded border-2 border-gray-300 cursor-pointer transition-all checked:bg-indigo-600 checked:border-indigo-600 hover:border-indigo-400 focus:ring-2 focus:ring-indigo-500/20 focus:ring-offset-0 accent-indigo-600"
					style="color-scheme: light;"
				/>
			</label>
			
			<!-- Drag Handle -->
			<div class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400">
				<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
					<circle cx="9" cy="7" r="1.5"/>
					<circle cx="9" cy="12" r="1.5"/>
					<circle cx="9" cy="17" r="1.5"/>
					<circle cx="15" cy="7" r="1.5"/>
					<circle cx="15" cy="12" r="1.5"/>
					<circle cx="15" cy="17" r="1.5"/>
				</svg>
			</div>

			<!-- Icon -->
			<div class="flex-shrink-0 w-10 h-10 rounded-lg bg-emerald-100 flex items-center justify-center">
				<svg class="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
				</svg>
			</div>

			<!-- Title - Click to expand -->
			<button
				onclick={() => isExpanded = true}
				class="flex-1 text-left hover:opacity-70 transition-opacity"
			>
				<h3 class="text-[17px] font-semibold text-gray-900">
					{group.group_title || 'Links'}
				</h3>
				<p class="text-sm text-gray-500 mt-0.5">
					{childrenCount} {childrenCount === 1 ? 'link' : 'links'}
				</p>
			</button>

			<!-- Actions -->
			<div class="flex items-center gap-2">
				<!-- More Menu -->
				<div class="relative">
					<button
						onclick={() => showMenu = !showMenu}
						class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
					>
						<svg class="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
							<circle cx="10" cy="4" r="1.5"/>
							<circle cx="10" cy="10" r="1.5"/>
							<circle cx="10" cy="16" r="1.5"/>
						</svg>
					</button>

					{#if showMenu}
						<div class="absolute right-0 top-full mt-1 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-10">
							<button
								onclick={() => { dispatch('addlink', group.id); showMenu = false; }}
								class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
								</svg>
								Add link
							</button>
							<button
								onclick={() => { dispatch('duplicate', group.id); showMenu = false; }}
								class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
								</svg>
								Duplicate
							</button>
							<div class="border-t border-gray-100 my-1"></div>
							<button
								onclick={() => { dispatch('delete', group.id); showMenu = false; }}
								class="w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-red-50 flex items-center gap-2"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
								</svg>
								Delete
							</button>
						</div>
					{/if}
				</div>

				<!-- Visibility Toggle -->
				<button
					onclick={() => dispatch('toggle', group.id)}
					class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
					title={group.is_active ? 'Hide' : 'Show'}
				>
					<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						{#if group.is_active}
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
						{:else}
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
						{/if}
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
				<button
					onclick={() => isExpanded = false}
					class="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
					</svg>
					<span class="text-sm font-medium">Back</span>
				</button>

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
						title="Visibility"
					>
						<svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
						</svg>
					</button>
				</div>
			</div>

			<div class="flex items-center gap-3 mb-4">
				<div class="w-12 h-12 rounded-xl bg-emerald-100 flex items-center justify-center">
					<svg class="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
					</svg>
				</div>
				<h2 class="text-2xl font-bold text-gray-900">{group.group_title || 'Links'}</h2>
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
				<button
					onclick={() => activeTab = 'settings'}
					class="pb-3 text-sm font-semibold transition-colors relative"
					class:text-gray-900={activeTab === 'settings'}
					class:text-gray-500={activeTab !== 'settings'}
				>
					Settings
					{#if activeTab === 'settings'}
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
					class="w-full px-6 py-4 bg-gray-900 hover:bg-gray-800 text-white rounded-xl transition-all flex items-center justify-center gap-2 text-base font-semibold shadow-sm mb-6"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
					</svg>
					Add link
				</button>

				<!-- Links List - Compact -->
				<div 
					class="space-y-2"
					use:dndzone={{items: dndItems, flipDurationMs: 200, dropTargetStyle: {}}}
					onconsider={handleDndConsider}
					onfinalize={handleDndFinalize}
				>
					{#each dndItems as item (item.id)}
						{@const link = item.data}
						<div class="bg-white rounded-xl shadow-sm p-3">
							<div class="flex gap-3 items-start">
								<!-- Drag Handle -->
								<button class="drag-handle pt-0.5 cursor-grab active:cursor-grabbing">
									<svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
										<circle cx="6" cy="5" r="1.5"/>
										<circle cx="6" cy="10" r="1.5"/>
										<circle cx="6" cy="15" r="1.5"/>
										<circle cx="10" cy="5" r="1.5"/>
										<circle cx="10" cy="10" r="1.5"/>
										<circle cx="10" cy="15" r="1.5"/>
									</svg>
								</button>

								<!-- Content -->
								<div class="flex-1 min-w-0">
									<h3 class="text-sm font-semibold text-gray-900 mb-0.5">{link.title}</h3>
									<p class="text-xs text-gray-500 mb-1">Đây là mô tả cho link</p>
									<p class="text-xs text-gray-400 truncate">{link.url}</p>
								</div>

								<!-- Thumbnail - Clickable to upload -->
								<button
									onclick={() => handleThumbnailClick(link.id)}
									class="w-16 h-16 rounded-lg flex-shrink-0 hover:opacity-80 transition-opacity group/thumb relative"
									title="Click to upload thumbnail"
								>
									{#if link.thumbnail_url}
										<img 
											src={link.thumbnail_url} 
											alt={link.title}
											class="w-full h-full rounded-lg object-cover"
										/>
										<div class="absolute inset-0 bg-black/50 opacity-0 group-hover/thumb:opacity-100 transition-opacity rounded-lg flex items-center justify-center">
											<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
											</svg>
										</div>
									{:else}
										<div class="w-full h-full rounded-lg bg-gray-100 flex items-center justify-center">
											<svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
											</svg>
										</div>
									{/if}
								</button>
							</div>

							<!-- Actions - Compact -->
							<div class="flex items-center justify-between mt-2 pt-2 border-t border-gray-100">
								<div class="flex items-center gap-1">
									<button 
										onclick={() => dispatch('pinlink', link.id)}
										class="p-1.5 hover:bg-gray-100 rounded-md transition-colors" 
										title={link.is_pinned ? 'Unpin' : 'Pin'}
									>
										<svg class="w-4 h-4 transition-colors" class:text-yellow-500={link.is_pinned} class:text-gray-500={!link.is_pinned} fill={link.is_pinned ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
										</svg>
									</button>
									<button class="p-1.5 hover:bg-gray-100 rounded-md transition-colors" title="Schedule">
										<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
										</svg>
									</button>
									<button class="p-1.5 hover:bg-gray-100 rounded-md transition-colors" title="Analytics">
										<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
										</svg>
									</button>
									<button class="p-1.5 hover:bg-gray-100 rounded-md transition-colors" title="Share">
										<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
										</svg>
									</button>
								</div>

								<div class="flex items-center gap-1">
									<button class="p-1.5 hover:bg-gray-100 rounded-md transition-colors" title="Analytics">
										<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 8v8m-4-5v5m-4-2v2m-2 4h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
										</svg>
									</button>
									
									<!-- More Menu (3 dots) -->
									<div class="relative">
										<button 
											onclick={() => linkMenuOpen = linkMenuOpen === link.id ? null : link.id}
											class="p-1.5 hover:bg-gray-100 rounded-md transition-colors" 
											title="More"
										>
											<svg class="w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
												<circle cx="10" cy="4" r="1.5"/>
												<circle cx="10" cy="10" r="1.5"/>
												<circle cx="10" cy="16" r="1.5"/>
											</svg>
										</button>
										
										{#if linkMenuOpen === link.id}
											<div class="absolute right-0 top-full mt-1 w-44 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-20">
												<!-- Edit Link -->
												<button
													onclick={() => { dispatch('editlink', link); linkMenuOpen = null; }}
													class="w-full px-3 py-2 text-left text-xs text-gray-700 hover:bg-gray-50 flex items-center gap-2"
												>
													<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
													</svg>
													Edit link
												</button>
												
												<!-- Toggle Open in New Tab -->
												<button
													onclick={() => { dispatch('toggleNewTab', link.id); }}
													class="w-full px-3 py-2 text-left text-xs text-gray-700 hover:bg-gray-50 flex items-center justify-between gap-2"
												>
													<div class="flex items-center gap-2">
														<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
														</svg>
														<span>Open in new tab</span>
													</div>
													{#if link.open_in_new_tab}
														<svg class="w-3.5 h-3.5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
														</svg>
													{/if}
												</button>
												
												<!-- Duplicate Link -->
												<button
													onclick={() => { dispatch('duplicatelink', link.id); linkMenuOpen = null; }}
													class="w-full px-3 py-2 text-left text-xs text-gray-700 hover:bg-gray-50 flex items-center gap-2"
												>
													<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
													</svg>
													Duplicate link
												</button>
												
												<div class="border-t border-gray-100 my-1"></div>
												
												<!-- Delete Link -->
												<button
													onclick={() => { dispatch('removefromgroup', link.id); linkMenuOpen = null; }}
													class="w-full px-3 py-2 text-left text-xs text-red-600 hover:bg-red-50 flex items-center gap-2"
												>
													<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
													</svg>
													Delete link
												</button>
											</div>
										{/if}
									</div>
									
									<div class="w-px h-4 bg-gray-200 mx-0.5"></div>
									<button 
										onclick={() => dispatch('togglelinkvisibility', link.id)}
										class="p-1.5 hover:bg-gray-100 rounded-md transition-colors" 
										title={link.is_active ? 'Hide' : 'Show'}
									>
										<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											{#if link.is_active}
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
											{:else}
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
											{/if}
										</svg>
									</button>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{:else if activeTab === 'layouts'}
				<div class="text-center py-12 text-gray-500">
					Layout settings coming soon...
				</div>
			{:else}
				<div class="text-center py-12 text-gray-500">
					Group settings coming soon...
				</div>
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
{#if showMenu || linkMenuOpen}
	<button
		class="fixed inset-0 z-0"
		onclick={() => { showMenu = false; linkMenuOpen = null; }}
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
