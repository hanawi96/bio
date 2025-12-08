<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	
	export let currentImage: string | null = null;
	export let uploading = false;
	
	const dispatch = createEventDispatcher();
	let fileInput: HTMLInputElement;
	let previewUrl: string | null = currentImage;
	let dragOver = false;
	let errorMessage = '';
	
	$: previewUrl = currentImage;
	
	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (file) processFile(file);
	}
	
	function handleDrop(event: DragEvent) {
		event.preventDefault();
		dragOver = false;
		const file = event.dataTransfer?.files[0];
		if (file && file.type.startsWith('image/')) {
			processFile(file);
		}
	}
	
	function processFile(file: File) {
		errorMessage = '';
		
		// Validate file size (5MB)
		if (file.size > 5 * 1024 * 1024) {
			errorMessage = 'File size must be less than 5MB';
			return;
		}
		
		// Validate file type
		if (!['image/jpeg', 'image/png', 'image/webp', 'image/gif'].includes(file.type)) {
			errorMessage = 'Only JPEG, PNG, WebP, and GIF images are allowed';
			return;
		}
		
		// Create preview
		const reader = new FileReader();
		reader.onload = (e) => {
			previewUrl = e.target?.result as string;
		};
		reader.readAsDataURL(file);
		
		// Dispatch upload event
		dispatch('upload', file);
	}
	
	function handleRemove() {
		previewUrl = null;
		errorMessage = '';
		if (fileInput) fileInput.value = '';
		dispatch('remove');
	}
</script>

<div class="space-y-3">
	{#if errorMessage}
		<div class="px-4 py-3 bg-red-50 border border-red-200 rounded-lg flex items-center gap-2 text-red-700 text-sm">
			<svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
				<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
			</svg>
			<span>{errorMessage}</span>
		</div>
	{/if}

	{#if previewUrl}
		<div class="relative group">
			<div class="relative w-full aspect-square rounded-xl overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 border-2 border-gray-200">
				<img 
					src={previewUrl} 
					alt="Thumbnail preview" 
					class="w-full h-full object-cover"
				/>
				{#if uploading}
					<div class="absolute inset-0 bg-gradient-to-br from-indigo-600/90 to-blue-600/90 backdrop-blur-sm flex items-center justify-center">
						<div class="text-center">
							<div class="relative w-16 h-16 mx-auto mb-3">
								<div class="absolute inset-0 border-4 border-white/20 rounded-full"></div>
								<div class="absolute inset-0 border-4 border-transparent border-t-white rounded-full animate-spin"></div>
							</div>
							<p class="text-white text-sm font-semibold">Uploading...</p>
							<p class="text-white/80 text-xs mt-1">Please wait</p>
						</div>
					</div>
				{/if}
			</div>
			
			{#if !uploading}
				<div class="absolute top-3 right-3 flex gap-2 opacity-0 group-hover:opacity-100 transition-all duration-200">
					<button
						onclick={() => fileInput.click()}
						class="p-2.5 bg-white hover:bg-gray-50 rounded-lg shadow-lg transition-all hover:scale-105"
						title="Change image"
					>
						<svg class="w-4 h-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
						</svg>
					</button>
					<button
						onclick={handleRemove}
						class="p-2.5 bg-red-500 hover:bg-red-600 rounded-lg shadow-lg transition-all hover:scale-105"
						title="Remove image"
					>
						<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
						</svg>
					</button>
				</div>
			{/if}
		</div>
	{:else}
		<div
			class="relative border-2 border-dashed rounded-xl transition-all cursor-pointer group"
			class:border-indigo-500={dragOver}
			class:bg-indigo-50={dragOver}
			class:border-gray-300={!dragOver}
			class:hover:border-indigo-400={!dragOver && !uploading}
			class:hover:bg-gray-50={!dragOver && !uploading}
			ondragover={(e) => { e.preventDefault(); dragOver = true; }}
			ondragleave={() => dragOver = false}
			ondrop={handleDrop}
			onclick={() => !uploading && fileInput.click()}
		>
			<div class="p-10 text-center">
				<div class="w-20 h-20 mx-auto mb-4 bg-gradient-to-br from-indigo-500 to-blue-600 rounded-2xl flex items-center justify-center shadow-xl shadow-indigo-500/30 group-hover:scale-105 transition-transform">
					<svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
					</svg>
				</div>
				
				<p class="text-sm text-gray-600 mb-1">
					Drag & drop or click to browse
				</p>
				<p class="text-xs text-gray-500">
					PNG, JPG, WebP or GIF • Max 5MB
				</p>
			</div>
		</div>
	{/if}
	
	<input
		bind:this={fileInput}
		type="file"
		accept="image/jpeg,image/png,image/webp,image/gif"
		onchange={handleFileSelect}
		class="hidden"
	/>
</div>
