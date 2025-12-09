<script lang="ts">
	import LinkLayoutDialog from '../dialogs/LinkLayoutDialog.svelte';
	import LinkScheduleDialog from '../dialogs/LinkScheduleDialog.svelte';
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';

	export let link: Link;
	export let selected: boolean = false;
	
	const dispatch = createEventDispatcher();
	let editedTitle = link.title;
	let editedUrl = link.url;
	let showLayoutDialog = false;
	let showEditDialog = false;
	let showScheduleDialog = false;
	
	// Local display values - update immediately for instant UI feedback
	let displayTitle = link.title;
	let displayUrl = link.url;
	
	// Sync display values when prop changes
	$: {
		displayTitle = link.title;
		displayUrl = link.url;
		editedTitle = link.title;
		editedUrl = link.url;
	}

	function handleDelete() {
		dispatch('delete', link.id);
	}

	function handleToggle() {
		dispatch('toggle', link.id);
	}

	function handleEditThumbnail() {
		dispatch('editThumbnail', link.id);
	}

	function handleUpdateLayout(event: CustomEvent) {
		dispatch('updateLayout', event.detail);
	}

	function handleDuplicate() {
		dispatch('duplicate', link.id);
	}

	function handleSelect() {
		dispatch('select', link.id);
	}

	function handleSchedule() {
		showScheduleDialog = true;
	}

	function handleSaveSchedule(event: CustomEvent) {
		dispatch('updateSchedule', event.detail);
	}

	function handlePin() {
		dispatch('pin', link.id);
	}
</script>

<!-- Link Block -->
<div class="relative bg-white rounded-xl {selected ? 'ring-2 ring-indigo-500' : ''} {link.is_pinned ? 'bg-yellow-50/50' : ''}">
	<div class="relative p-4">
		<div class="flex items-center gap-3">
			<!-- Checkbox - Minimal -->
			<label class="flex items-center cursor-pointer">
				<input 
					type="checkbox" 
					checked={selected}
					onchange={handleSelect}
					class="w-4 h-4 rounded border-2 border-gray-300 cursor-pointer transition-all
						checked:bg-indigo-600 checked:border-indigo-600 
						hover:border-indigo-400
						focus:ring-2 focus:ring-indigo-500/20 focus:ring-offset-0
						accent-indigo-600"
					style="color-scheme: light;"
				/>
			</label>
			
			<!-- Drag Handle - Subtle -->
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

			<!-- Thumbnail - Smaller & Cleaner -->
			{#if link.thumbnail_url}
				<div class="flex-shrink-0 w-12 h-12 rounded-lg overflow-hidden bg-gray-100">
					<img 
						src={link.thumbnail_url} 
						alt={link.title}
						class="w-full h-full object-cover object-center"
					/>
				</div>
			{:else}
				<!-- Placeholder - Minimal -->
				<div class="flex-shrink-0 w-12 h-12 rounded-lg bg-gray-100 flex items-center justify-center">
					<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
					</svg>
				</div>
			{/if}

			<!-- Content - Compact -->
			<div class="flex-1 min-w-0 space-y-1">
				<!-- Title - Simple Display -->
				<div class="flex items-center gap-2">
					<h3 class="text-sm font-medium text-gray-900 truncate flex-1">
						{displayTitle}
					</h3>
					
					<!-- Pinned Badge -->
					{#if link.is_pinned}
						<span class="flex-shrink-0 px-2 py-0.5 bg-yellow-100 text-yellow-700 text-xs font-medium rounded">
							Pinned
						</span>
					{/if}
				</div>



				<!-- Stats & Actions Bar - Compact -->
				<div class="flex items-center gap-2 pt-1">
					<!-- Clicks counter - Minimal -->
					<div class="flex items-center gap-1 text-xs text-gray-500">
						<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
						</svg>
						<span>{link.clicks}</span>
					</div>

					<!-- Schedule Badge -->
					{#if link.scheduled_at || link.expires_at}
						<div class="flex items-center gap-1 px-2 py-0.5 bg-purple-100 text-purple-700 text-xs font-medium rounded">
							<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
							</svg>
							<span>Scheduled</span>
						</div>
					{/if}

					<!-- Divider -->
					<div class="w-px h-3 bg-gray-200"></div>

					<!-- All Action buttons - Smaller -->
					<div class="flex items-center gap-0.5">
						<!-- Edit Button -->
						<button 
							onclick={() => showEditDialog = true}
							class="p-1 rounded text-gray-400 hover:text-gray-600 hover:bg-gray-50"
							title="Edit"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
							</svg>
						</button>
						
						<!-- Pin Toggle - Smaller -->
						<button 
							onclick={handlePin}
							class="p-1 rounded {link.is_pinned ? 'text-yellow-600 hover:bg-yellow-50' : 'text-gray-400 hover:bg-gray-50'}"
							title={link.is_pinned ? 'Unpin' : 'Pin'}
						>
							<svg class="w-3.5 h-3.5" fill={link.is_pinned ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
							</svg>
						</button>

						<!-- Visibility Toggle -->
						<button 
							onclick={handleToggle}
							class="p-1 rounded {link.is_active ? 'text-green-600 hover:bg-green-50' : 'text-gray-400 hover:bg-gray-50'}"
							title={link.is_active ? 'Hide' : 'Show'}
						>
							{#if link.is_active}
								<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
								</svg>
							{:else}
								<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
								</svg>
							{/if}
						</button>

						<button 
							onclick={handleEditThumbnail}
							class="p-1 rounded text-gray-400 hover:text-gray-600 hover:bg-gray-50"
							title="Edit image"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
							</svg>
						</button>

						<!-- Schedule Button -->
						<button 
							onclick={handleSchedule}
							class="p-1 rounded {link.scheduled_at || link.expires_at ? 'text-purple-600 hover:bg-purple-50' : 'text-gray-400 hover:bg-gray-50'}"
							title="Schedule"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
							</svg>
						</button>
						
						<button 
							onclick={() => showLayoutDialog = true}
							class="relative p-1 rounded text-gray-400 hover:text-gray-600 hover:bg-gray-50"
							title="Change layout"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h4a1 1 0 011 1v7a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v7a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 16a1 1 0 011-1h4a1 1 0 011 1v3a1 1 0 01-1 1H5a1 1 0 01-1-1v-3zM14 16a1 1 0 011-1h4a1 1 0 011 1v3a1 1 0 01-1 1h-4a1 1 0 01-1-1v-3z"/>
							</svg>
							{#if link.layout_type === 'featured'}
								<span class="absolute -top-0.5 -right-0.5 w-1.5 h-1.5 bg-indigo-600 rounded-full"></span>
							{/if}
						</button>

						<button 
							onclick={handleDuplicate}
							class="p-1 rounded text-gray-400 hover:text-gray-600 hover:bg-gray-50"
							title="Duplicate"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
							</svg>
						</button>
						
						<button 
							onclick={handleDelete} 
							class="p-1 rounded text-gray-400 hover:text-red-600 hover:bg-red-50"
							title="Delete"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
							</svg>
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Layout Dialog -->
{#key link.layout_type}
	<LinkLayoutDialog 
		bind:open={showLayoutDialog} 
		{link}
		on:updateLayout={handleUpdateLayout}
	/>
{/key}

<!-- Schedule Dialog -->
<LinkScheduleDialog 
	bind:open={showScheduleDialog} 
	{link}
	on:save={handleSaveSchedule}
/>

<!-- Edit Link Dialog -->
{#if showEditDialog}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onclick={() => showEditDialog = false}>
		<div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6" onclick={(e) => e.stopPropagation()}>
			<!-- Header -->
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-lg font-semibold text-gray-900">Edit Link</h3>
				<button 
					onclick={() => showEditDialog = false}
					class="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
					</svg>
				</button>
			</div>

			<!-- Form -->
			<div class="space-y-4">
				<!-- Title Input -->
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">
						Title
					</label>
					<input
						type="text"
						bind:value={editedTitle}
						placeholder="Enter link title"
						class="w-full px-4 py-2.5 bg-gray-50 border-0 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:bg-white transition-all"
					/>
				</div>

				<!-- URL Input -->
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">
						URL
					</label>
					<input
						type="text"
						bind:value={editedUrl}
						placeholder="https://example.com"
						class="w-full px-4 py-2.5 bg-gray-50 border-0 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:bg-white transition-all"
					/>
				</div>
			</div>

			<!-- Actions -->
			<div class="flex items-center gap-3 mt-6">
				<button
					onclick={() => showEditDialog = false}
					class="flex-1 px-4 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium rounded-xl transition-colors"
				>
					Cancel
				</button>
				<button
					onclick={() => {
						if (editedTitle !== link.title) {
							displayTitle = editedTitle;
							dispatch('update', { id: link.id, title: editedTitle, url: link.url });
						}
						if (editedUrl !== link.url) {
							displayUrl = editedUrl;
							dispatch('update', { id: link.id, title: editedTitle, url: editedUrl });
						}
						showEditDialog = false;
					}}
					class="flex-1 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl transition-colors"
				>
					Save Changes
				</button>
			</div>
		</div>
	</div>
{/if}
