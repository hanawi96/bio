<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { Link } from '$lib/api/links';
	
	export let open = false;
	export let link: Link | null = null;
	export let groupTitle = '';
	
	const dispatch = createEventDispatcher();
	
	let title = '';
	let url = '';
	let urlError = '';
	let fileInput: HTMLInputElement;
	
	// Mode: 'add' or 'edit'
	$: mode = link ? 'edit' : 'add';
	
	// Reactive thumbnail URL synced with link prop
	$: thumbnailUrl = link?.thumbnail_url || '';
	
	$: if (link && open) {
		title = link.title;
		url = link.url;
		urlError = '';
	}
	
	function isValidUrl(urlString: string): boolean {
		try {
			const urlObj = new URL(urlString);
			return urlObj.protocol === 'http:' || urlObj.protocol === 'https:';
		} catch {
			return false;
		}
	}
	
	function validateUrl() {
		if (!url.trim()) {
			urlError = '';
			return;
		}
		
		if (!isValidUrl(url.trim())) {
			urlError = 'Please enter a valid URL (e.g., https://example.com)';
		} else {
			urlError = '';
		}
	}
	
	function handleThumbnailClick() {
		fileInput?.click();
	}
	
	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (file && link) {
			dispatch('uploadThumbnail', { linkId: link.id, file });
		}
	}
	
	function handleRemoveThumbnail() {
		if (link) {
			dispatch('removeThumbnail', link.id);
		}
	}
	
	function handleSave() {
		if (!title.trim() || !url.trim()) return;
		
		validateUrl();
		if (urlError) return;
		
		if (mode === 'edit' && link) {
			dispatch('save', {
				id: link.id,
				title: title.trim(),
				url: url.trim()
			});
		} else {
			dispatch('add', {
				title: title.trim(),
				url: url.trim()
			});
		}
		
		open = false;
	}
	
	function resetForm() {
		title = '';
		url = '';
		urlError = '';
	}
	
	$: if (!open) {
		resetForm();
	}
	
	function handleClose() {
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-2xl p-0 overflow-hidden" showCloseButton={false}>
		<!-- Simple Header -->
		<div class="flex items-center justify-between px-6 py-4 border-b">
			<h2 class="text-xl font-bold text-gray-900">{mode === 'edit' ? 'Edit' : 'Add'}</h2>
			<button 
				onclick={handleClose}
				class="text-gray-400 hover:text-gray-600 transition-colors"
			>
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
		
		<!-- Main Content -->
		<div class="p-6">
			<!-- Two Column Layout -->
			<div class="grid grid-cols-[1fr_180px] gap-6 mb-6">
				<!-- Left: Input Fields -->
				<div class="space-y-4">
					<!-- Title Input -->
					<input
						id="edit-link-title"
						type="text"
						bind:value={title}
						placeholder="Title"
						class="w-full px-4 py-3 bg-gray-100 border-0 rounded-lg focus:bg-white focus:ring-2 focus:ring-gray-200 transition-all text-gray-900 placeholder-gray-500"
					/>

					<!-- URL Input -->
					<input
						id="edit-link-url"
						type="url"
						bind:value={url}
						placeholder="URL"
						class="w-full px-4 py-3 bg-gray-100 border-0 rounded-lg transition-all text-gray-900 placeholder-gray-500 {urlError ? 'ring-2 ring-red-400 bg-red-50' : 'focus:bg-white focus:ring-2 focus:ring-gray-200'}"
						onblur={validateUrl}
					/>
					{#if urlError}
						<p class="text-xs text-red-600">{urlError}</p>
					{/if}
				</div>

				<!-- Right: Picture Upload -->
				<div class="relative h-full min-h-[140px]">
					{#if thumbnailUrl}
						<!-- Display existing thumbnail -->
						<div class="relative h-full rounded-lg overflow-hidden border-2 border-gray-300 group">
							<img 
								src={thumbnailUrl} 
								alt="Thumbnail" 
								class="w-full h-full object-cover"
							/>
							<div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
								<button
									type="button"
									onclick={handleThumbnailClick}
									class="p-2 bg-white rounded-lg hover:bg-gray-100 transition-colors"
									title="Change image"
								>
									<svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
									</svg>
								</button>
								<button
									type="button"
									onclick={handleRemoveThumbnail}
									class="p-2 bg-white rounded-lg hover:bg-red-100 transition-colors"
									title="Remove image"
								>
									<svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
									</svg>
								</button>
							</div>
						</div>
					{:else}
						<!-- Upload placeholder -->
						<button 
							type="button"
							onclick={handleThumbnailClick}
							class="h-full min-h-[140px] border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center gap-2 hover:border-gray-400 transition-colors bg-gray-50 w-full"
						>
							<svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
							</svg>
							<span class="text-sm text-gray-500">Picture</span>
						</button>
					{/if}
					<input 
						bind:this={fileInput}
						type="file" 
						accept="image/*" 
						onchange={handleFileSelect}
						class="hidden"
					/>
				</div>
			</div>

			<!-- Toggle Switch -->
			<div class="flex items-center justify-between py-4">
				<span class="text-base text-gray-900">Make this a highlighted link</span>
				<button
					type="button"
					onclick={() => {}}
					class="relative inline-flex h-7 w-12 items-center rounded-full transition-colors bg-gray-300"
				>
					<span class="inline-block h-5 w-5 transform rounded-full bg-white shadow-sm transition-transform translate-x-1"></span>
				</button>
			</div>
		</div>

		<!-- Save Button -->
		<div class="px-6 pb-6">
			<button 
				type="button"
				onclick={handleSave}
				disabled={!title.trim() || !url.trim() || !!urlError}
				class="w-full py-4 bg-gradient-to-r from-red-500 via-pink-500 to-purple-500 hover:from-red-600 hover:via-pink-600 hover:to-purple-600 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all font-bold text-lg uppercase tracking-wide shadow-lg"
			>
				Save
			</button>
		</div>
	</Dialog.Content>
</Dialog.Root>
