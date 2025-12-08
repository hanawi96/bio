<script lang="ts">
	import LinkLayoutDialog from './LinkLayoutDialog.svelte';
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';

	export let link: Link;
	export let selected: boolean = false;
	
	const dispatch = createEventDispatcher();
	let isEditingTitle = false;
	let isEditingUrl = false;
	let editedTitle = link.title;
	let editedUrl = link.url;
	let showLayoutDialog = false;
	let showMoreMenu = false;

	function handleSaveTitle() {
		if (editedTitle !== link.title) {
			dispatch('update', { id: link.id, title: editedTitle, url: link.url });
		}
		isEditingTitle = false;
	}

	function handleSaveUrl() {
		if (editedUrl !== link.url) {
			dispatch('update', { id: link.id, title: link.title, url: editedUrl });
		}
		isEditingUrl = false;
	}

	function handleDelete() {
		dispatch('delete', link.id);
		showMoreMenu = false;
	}

	function handleToggle() {
		dispatch('toggle', link.id);
	}

	function handleEditThumbnail() {
		dispatch('editThumbnail', link.id);
		showMoreMenu = false;
	}

	function handleUpdateLayout(event: CustomEvent) {
		dispatch('updateLayout', event.detail);
	}

	function handleDuplicate() {
		dispatch('duplicate', link.id);
		showMoreMenu = false;
	}

	function handleSelect() {
		dispatch('select', link.id);
	}

	function handlePin() {
		dispatch('pin', link.id);
	}

	function handleCopyUrl() {
		navigator.clipboard.writeText(link.url);
		showMoreMenu = false;
	}

	function handleOpenUrl() {
		window.open(link.url, '_blank');
		showMoreMenu = false;
	}
</script>

<!-- Modern LinkCard với thiết kế tinh tế và chuyên nghiệp -->
<div class="relative bg-white rounded-2xl shadow-sm overflow-hidden transition-all {selected ? 'ring-2 ring-indigo-500' : ''} {link.is_pinned ? 'ring-1 ring-yellow-400' : ''}">
	{#if link.is_pinned}
		<div class="absolute top-0 right-0 bg-gradient-to-l from-yellow-400 to-transparent px-4 py-1 rounded-bl-xl">
			<span class="text-xs font-bold text-yellow-900">PINNED</span>
		</div>
	{/if}
	
	<div class="relative p-5">
		<div class="flex items-center gap-4">
			<!-- Checkbox -->
			<label class="flex items-center cursor-pointer">
				<input 
					type="checkbox" 
					checked={selected}
					onchange={handleSelect}
					class="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 cursor-pointer"
				/>
			</label>
			
			<!-- Drag Handle -->
			<div class="drag-handle cursor-grab active:cursor-grabbing text-gray-300">
				<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
					<circle cx="8" cy="6" r="1.5"/>
					<circle cx="8" cy="12" r="1.5"/>
					<circle cx="8" cy="18" r="1.5"/>
					<circle cx="14" cy="6" r="1.5"/>
					<circle cx="14" cy="12" r="1.5"/>
					<circle cx="14" cy="18" r="1.5"/>
				</svg>
			</div>

			<!-- Thumbnail -->
			{#if link.thumbnail_url}
				<div class="flex-shrink-0 w-16 h-16 rounded-xl overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 ring-2 ring-gray-200 shadow-sm">
					<img 
						src={link.thumbnail_url} 
						alt={link.title}
						class="w-full h-full object-cover object-center"
					/>
				</div>
			{:else}
				<!-- Placeholder khi không có thumbnail -->
				<div class="flex-shrink-0 w-16 h-16 rounded-xl bg-gradient-to-br from-indigo-100 to-purple-100 ring-2 ring-indigo-200 flex items-center justify-center shadow-sm">
					<svg class="w-7 h-7 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
					</svg>
				</div>
			{/if}

			<!-- Content -->
			<div class="flex-1 min-w-0 space-y-2.5">
				<!-- Title -->
				<div class="flex items-center gap-2">
					{#if isEditingTitle}
						<input
							type="text"
							bind:value={editedTitle}
							onblur={handleSaveTitle}
							onkeydown={(e) => e.key === 'Enter' && handleSaveTitle()}
							class="flex-1 px-3 py-2 text-lg font-semibold border border-indigo-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
							autofocus
						/>
					{:else}
						<h3 class="text-lg font-semibold text-gray-900 truncate">
							{link.title}
						</h3>
						<button 
							onclick={() => isEditingTitle = true} 
							class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-indigo-600 hover:bg-indigo-50"
							title="Chỉnh sửa tiêu đề"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
							</svg>
						</button>
					{/if}
				</div>

				<!-- URL -->
				<div class="flex items-center gap-2">
					{#if isEditingUrl}
						<input
							type="text"
							bind:value={editedUrl}
							onblur={handleSaveUrl}
							onkeydown={(e) => e.key === 'Enter' && handleSaveUrl()}
							class="flex-1 px-3 py-2 text-sm border border-indigo-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
							autofocus
						/>
					{:else}
						<div class="flex items-center gap-1.5 flex-1 min-w-0">
							<svg class="w-3.5 h-3.5 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
							</svg>
							<p class="text-sm text-gray-500 truncate">
								{link.url}
							</p>
						</div>
						<button 
							onclick={() => isEditingUrl = true} 
							class="flex-shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-indigo-600 hover:bg-indigo-50"
							title="Chỉnh sửa URL"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
							</svg>
						</button>
					{/if}
				</div>

				<!-- Stats & Actions Bar -->
				<div class="flex items-center gap-3 pt-1">
					<!-- Clicks counter với icon -->
					<div class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-gray-50">
						<svg class="w-3.5 h-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
						</svg>
						<span class="text-xs font-medium text-gray-600">
							{link.clicks}
						</span>
					</div>

					<!-- Divider -->
					<div class="w-px h-4 bg-gray-200"></div>

					<!-- All Action buttons in one row -->
					<div class="flex items-center gap-1">
						<!-- Pin Toggle -->
						<button 
							onclick={handlePin}
							class="p-1.5 rounded-lg {link.is_pinned ? 'text-yellow-600 hover:bg-yellow-50' : 'text-gray-400 hover:bg-gray-50'}"
							title={link.is_pinned ? 'Bỏ ghim' : 'Ghim lên đầu'}
						>
							<svg class="w-4 h-4" fill={link.is_pinned ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
							</svg>
						</button>

						<!-- Visibility Toggle - Eye Icon -->
						<button 
							onclick={handleToggle}
							class="p-1.5 rounded-lg {link.is_active ? 'text-green-600 hover:bg-green-50' : 'text-gray-400 hover:bg-gray-50'}"
							title={link.is_active ? 'Ẩn link' : 'Hiện link'}
						>
							{#if link.is_active}
								<!-- Eye Open -->
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
								</svg>
							{:else}
								<!-- Eye Closed (with slash) -->
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
								</svg>
							{/if}
						</button>

						<button 
							onclick={handleEditThumbnail}
							class="p-1.5 rounded-lg text-gray-400 hover:text-indigo-600 hover:bg-indigo-50"
							title="Chỉnh sửa hình ảnh"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
							</svg>
						</button>
						
						<button 
							onclick={() => showLayoutDialog = true}
							class="relative p-1.5 rounded-lg text-gray-400 hover:text-indigo-600 hover:bg-indigo-50"
							title="Thay đổi bố cục"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h4a1 1 0 011 1v7a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v7a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 16a1 1 0 011-1h4a1 1 0 011 1v3a1 1 0 01-1 1H5a1 1 0 01-1-1v-3zM14 16a1 1 0 011-1h4a1 1 0 011 1v3a1 1 0 01-1 1h-4a1 1 0 01-1-1v-3z"/>
							</svg>
							{#if link.layout_type === 'featured'}
								<span class="absolute -top-0.5 -right-0.5 w-2 h-2 bg-gradient-to-br from-indigo-500 to-purple-500 rounded-full ring-2 ring-white"></span>
							{/if}
						</button>

						<button 
							onclick={handleDuplicate}
							class="p-1.5 rounded-lg text-gray-400 hover:text-blue-600 hover:bg-blue-50"
							title="Nhân bản link"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
							</svg>
						</button>
						
						<button 
							onclick={handleDelete} 
							class="p-1.5 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50"
							title="Xóa"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
