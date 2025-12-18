<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { Link } from '$lib/api/links';
	
	export let open = false;
	export let link: Link | null = null;
	export let isUploading = false;
	
	const dispatch = createEventDispatcher();
	
	let title = '';
	let url = '';
	let description = '';
	let urlError = '';
	let fileInput: HTMLInputElement;
	let selectedFile: File | null = null;
	let previewUrl: string = '';
	
	// Mode: 'add' or 'edit'
	$: mode = link ? 'edit' : 'add';
	
	// Thumbnail URL: use preview if available, otherwise use link's thumbnail
	$: thumbnailUrl = previewUrl || link?.thumbnail_url || '';
	
	$: if (link && open) {
		title = link.title;
		url = link.url;
		description = link.description || '';
		urlError = '';
		selectedFile = null;
		previewUrl = '';
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
		if (!file) return;
		
		selectedFile = file;
		previewUrl = URL.createObjectURL(file);
	}
	
	function handleRemoveThumbnail() {
		// Clear local preview and file
		if (previewUrl) {
			URL.revokeObjectURL(previewUrl);
			previewUrl = '';
		}
		selectedFile = null;
		
		// If editing existing link with saved thumbnail, remove from server
		if (link && link.thumbnail_url) {
			dispatch('removeThumbnail', link.id);
		}
	}
	
	function handleSave() {
		if (!title.trim() || !url.trim()) return;
		
		validateUrl();
		if (urlError) return;
		
		console.log('💾 handleSave called', { hasFile: !!selectedFile, mode });
		
		if (mode === 'edit' && link) {
			dispatch('save', {
				id: link.id,
				title: title.trim(),
				url: url.trim(),
				description: description.trim() || null,
				file: selectedFile
			});
		} else {
			dispatch('add', {
				title: title.trim(),
				url: url.trim(),
				description: description.trim() || null,
				file: selectedFile
			});
		}
		
		// Only close modal immediately if no file upload needed
		if (!selectedFile) {
			console.log('ℹ️ No file, closing modal immediately');
			open = false;
		} else {
			console.log('⏳ File selected, keeping modal open for upload');
		}
	}
	
	function resetForm() {
		title = '';
		url = '';
		description = '';
		urlError = '';
		selectedFile = null;
		if (previewUrl) {
			URL.revokeObjectURL(previewUrl);
			previewUrl = '';
		}
	}
	
	$: if (!open) {
		resetForm();
	}
	
	function handleClose() {
		open = false;
	}
	
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		} else if (event.key === 'Enter' && (event.ctrlKey || event.metaKey)) {
			// Ctrl+Enter or Cmd+Enter to save
			handleSave();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-2xl p-0 overflow-hidden" showCloseButton={false}>
		<!-- Simple Header -->
		<div class="flex items-center justify-between px-6 py-4 border-b">
			<h2 class="text-xl font-bold text-gray-900">{mode === 'edit' ? 'Edit' : 'Add'}</h2>
			<button 
				onclick={handleClose}
				class="text-gray-500 hover:text-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-300 rounded-lg p-1"
				aria-label="Close dialog"
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
					<label for="edit-link-title" class="sr-only">Link Title</label>
					<input
						id="edit-link-title"
						type="text"
						bind:value={title}
						placeholder="Title"
						required
						class="w-full px-4 py-3 bg-gray-100 border-0 rounded-lg focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:outline-none transition-all text-gray-900 placeholder-gray-500"
					/>

					<!-- URL Input -->
					<label for="edit-link-url" class="sr-only">Link URL</label>
					<input
						id="edit-link-url"
						type="url"
						bind:value={url}
						placeholder="URL (e.g., https://example.com)"
						required
						aria-invalid={!!urlError}
						aria-describedby={urlError ? 'url-error' : undefined}
						class="w-full px-4 py-3 bg-gray-100 border-0 rounded-lg transition-all text-gray-900 placeholder-gray-500 focus:outline-none {urlError ? 'ring-2 ring-red-500 bg-red-50' : 'focus:bg-white focus:ring-2 focus:ring-indigo-500'}"
						onblur={validateUrl}
					/>
					{#if urlError}
						<p id="url-error" class="text-xs text-red-700 font-medium" role="alert">{urlError}</p>
					{/if}

					<!-- Description Input -->
					<label for="edit-link-description" class="sr-only">Link Description (optional)</label>
					<textarea
						id="edit-link-description"
						bind:value={description}
						placeholder="Description (optional)"
						rows="2"
						class="w-full px-4 py-3 bg-gray-100 border-0 rounded-lg focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:outline-none transition-all text-gray-900 placeholder-gray-500 resize-none"
					></textarea>
				</div>

				<!-- Right: Picture Upload -->
				<div class="relative w-full aspect-square">
					{#if thumbnailUrl}
						<!-- Display existing thumbnail -->
						<div class="relative w-full h-full rounded-lg overflow-hidden border-2 border-gray-300 group">
							<img 
								src={thumbnailUrl} 
								alt="Thumbnail" 
								class="w-full h-full object-cover"
							/>
							
							<!-- Loading Overlay -->
							{#if isUploading}
								<div class="absolute inset-0 bg-black/60 flex flex-col items-center justify-center gap-2">
									<svg class="w-10 h-10 text-white animate-spin" fill="none" viewBox="0 0 24 24">
										<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
										<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
									</svg>
									<span class="text-white text-sm font-medium">Uploading...</span>
								</div>
							{:else}
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
							{/if}
						</div>
					{:else}
						<!-- Upload placeholder -->
						<button 
							type="button"
							onclick={handleThumbnailClick}
							class="w-full h-full border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center gap-2 hover:border-gray-400 transition-colors bg-gray-50"
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
					role="switch"
					aria-checked="false"
					aria-label="Toggle highlighted link"
				>
					<span class="inline-block h-5 w-5 transform rounded-full bg-white shadow-sm transition-transform translate-x-1"></span>
				</button>
			</div>
		</div>

		<!-- Save Button -->
		<div class="px-6 pb-6">
			<button 
				type="submit"
				onclick={handleSave}
				disabled={!title.trim() || !url.trim() || !!urlError || isUploading}
				class="w-full py-4 bg-gradient-to-r from-red-500 via-pink-500 to-purple-500 hover:from-red-600 hover:via-pink-600 hover:to-purple-600 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all font-bold text-lg uppercase tracking-wide shadow-lg focus:outline-none focus:ring-4 focus:ring-purple-300"
				aria-label="{isUploading ? 'Uploading...' : 'Save link'}"
			>
				{isUploading ? 'Uploading...' : 'Save'}
			</button>
		</div>
	</Dialog.Content>
</Dialog.Root>
