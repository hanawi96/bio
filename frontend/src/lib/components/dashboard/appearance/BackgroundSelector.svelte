<script lang="ts">
	import { globalTheme } from '$lib/stores/theme';
	import { profileApi } from '$lib/api/profile';
	import { auth } from '$lib/stores/auth';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	let selectedType = $state<'solid' | 'gradient' | 'image' | 'video'>('solid');
	let solidColor = $state('#ffffff');
	let gradientFrom = $state('#faf5ff');
	let gradientTo = $state('#eff6ff');
	let imageUrl = $state('');
	let videoUrl = $state('');
	let uploading = $state(false);
	
	// Track if user manually selected a type
	let manuallySelectedType = $state<string | null>(null);

	// Ensure color is always in HEX format
	function normalizeHex(color: string): string {
		if (!color) return '#ffffff';
		if (color.startsWith('#')) return color.toUpperCase();
		return `#${color}`.toUpperCase();
	}

	// Reactive subscription to globalTheme store
	let unsubscribe: (() => void) | null = null;
	
	$effect(() => {
		// Subscribe to theme store changes
		unsubscribe = globalTheme.subscribe((theme) => {
			selectedType = theme.pageBackgroundType || 'solid';
			solidColor = normalizeHex(theme.pageBackground || '#ffffff');
			gradientFrom = normalizeHex(theme.pageGradientFrom || '#faf5ff');
			gradientTo = normalizeHex(theme.pageGradientTo || '#eff6ff');
			imageUrl = theme.pageBackgroundImage || '';
			videoUrl = theme.pageBackgroundVideo || '';
			
			// Reset manual selection when theme changes externally
			if (manuallySelectedType && selectedType !== manuallySelectedType) {
				manuallySelectedType = null;
			}
		});
		
		// Cleanup on unmount
		return () => {
			if (unsubscribe) unsubscribe();
		};
	});
	
	// Handle manual type selection
	function selectType(type: 'solid' | 'gradient' | 'image' | 'video') {
		manuallySelectedType = type;
		selectedType = type;
	}

	let saving = $state(false);

	async function applyBackground() {
		if (saving) return;
		
		// Validate
		if (selectedType === 'image' && !imageUrl) {
			toast.error('Please enter an image URL or upload an image');
			return;
		}
		if (selectedType === 'video' && !videoUrl) {
			toast.error('Please enter a video URL or upload a video');
			return;
		}
		
		saving = true;
		const theme = globalTheme.getCurrent();
		const updates: any = {
			...theme,
			pageBackgroundType: selectedType,
			pageBackground: solidColor,
			pageGradientFrom: gradientFrom,
			pageGradientTo: gradientTo,
			pageBackgroundImage: imageUrl,
			pageBackgroundVideo: videoUrl
		};

		globalTheme.update(updates);

		try {
			await profileApi.updateProfile({
				theme_config: updates
			}, get(auth).token!);
			toast.success('Background updated!');
		} catch (e: any) {
			toast.error(e.message || 'Failed to update background');
		} finally {
			saving = false;
		}
	}

	async function uploadFile(e: Event, type: 'image' | 'video') {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		uploading = true;
		try {
			const formData = new FormData();
			formData.append('file', file);

			const response = await fetch('/api/upload', {
				method: 'POST',
				headers: { Authorization: `Bearer ${get(auth).token}` },
				body: formData
			});

			if (!response.ok) {
				const error = await response.json();
				throw new Error(error.message || 'Upload failed');
			}
			
			const data = await response.json();
			
			if (type === 'image') {
				imageUrl = data.url;
				selectedType = 'image';
			} else {
				videoUrl = data.url;
				selectedType = 'video';
			}
			
			await applyBackground();
		} catch (e: any) {
			toast.error(e.message || 'Upload failed');
		} finally {
			uploading = false;
		}
	}
</script>

<div class="space-y-6">
	<div>
		<h3 class="text-lg font-bold text-gray-900 mb-4">Background</h3>
		
		<!-- Type Selector -->
		<div class="grid grid-cols-5 gap-3 mb-6">
			<button
				onclick={() => selectType('solid')}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {manuallySelectedType === 'solid' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
			>
				<div class="w-12 h-12 rounded-lg bg-gradient-to-br from-indigo-100 to-indigo-200 flex items-center justify-center">
					<div class="w-6 h-6 rounded bg-indigo-600"></div>
				</div>
				<span class="text-xs font-medium text-gray-700">Solid</span>
			</button>

			<button
				onclick={() => selectType('gradient')}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {manuallySelectedType === 'gradient' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
			>
				<div class="w-12 h-12 rounded-lg bg-gradient-to-br from-purple-400 to-pink-400 flex items-center justify-center">
					<svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
						<path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z"/>
						<path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd"/>
					</svg>
				</div>
				<span class="text-xs font-medium text-gray-700">Gradient</span>
			</button>

			<button
				onclick={() => selectType('image')}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {manuallySelectedType === 'image' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
			>
				<div class="w-12 h-12 rounded-lg bg-gray-200 flex items-center justify-center">
					<svg class="w-6 h-6 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd"/>
					</svg>
				</div>
				<span class="text-xs font-medium text-gray-700">Image</span>
			</button>

			<button
				onclick={() => selectType('video')}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {manuallySelectedType === 'video' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
			>
				<div class="w-12 h-12 rounded-lg bg-gray-300 flex items-center justify-center">
					<svg class="w-6 h-6 text-gray-700" fill="currentColor" viewBox="0 0 20 20">
						<path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z"/>
					</svg>
				</div>
				<span class="text-xs font-medium text-gray-700">Video</span>
			</button>

			<button
				disabled
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 border-gray-200 opacity-50 cursor-not-allowed relative group"
				title="AI-generated backgrounds coming soon!"
			>
				<div class="w-12 h-12 rounded-lg bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center">
					<svg class="w-6 h-6 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
						<path d="M13 7H7v6h6V7z"/>
						<path fill-rule="evenodd" d="M7 2a1 1 0 012 0v1h2V2a1 1 0 112 0v1h2a2 2 0 012 2v2h1a1 1 0 110 2h-1v2h1a1 1 0 110 2h-1v2a2 2 0 01-2 2h-2v1a1 1 0 11-2 0v-1H9v1a1 1 0 11-2 0v-1H5a2 2 0 01-2-2v-2H2a1 1 0 110-2h1V9H2a1 1 0 010-2h1V5a2 2 0 012-2h2V2zM5 5h10v10H5V5z" clip-rule="evenodd"/>
					</svg>
				</div>
				<span class="text-xs font-medium text-gray-400">AI (Soon)</span>
			</button>
		</div>

		<!-- Settings (only show when manually selected) -->
		{#if manuallySelectedType === 'solid'}
			<div class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Color</label>
					<div class="flex gap-2 items-center">
						<input
							type="color"
							value={solidColor}
							oninput={(e) => solidColor = normalizeHex(e.currentTarget.value)}
							class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer"
							style="padding: 0; -webkit-appearance: none; appearance: none;"
						/>
						<input
							type="text"
							value={solidColor}
							oninput={(e) => {
								const val = e.currentTarget.value.toUpperCase();
								if (/^#[0-9A-F]{6}$/.test(val)) solidColor = val;
								else if (/^#[0-9A-F]{0,6}$/.test(val)) e.currentTarget.value = val;
							}}
							class="flex-1 px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-colors"
							placeholder="#FFFFFF"
							maxlength="7"
						/>
					</div>
				</div>
				
				<!-- Color Presets -->
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Trendy Colors</label>
					<div class="flex flex-wrap gap-2">
						{#each [
							{ color: '#FFFFFF', name: 'White' },
							{ color: '#F8FAFC', name: 'Slate' },
							{ color: '#FEF3C7', name: 'Cream' },
							{ color: '#DBEAFE', name: 'Sky' },
							{ color: '#E0E7FF', name: 'Indigo' },
							{ color: '#FCE7F3', name: 'Pink' },
							{ color: '#DCFCE7', name: 'Green' },
							{ color: '#FEF3C7', name: 'Amber' },
							{ color: '#1E293B', name: 'Dark' },
							{ color: '#0F172A', name: 'Night' },
							{ color: '#4F46E5', name: 'Purple' },
							{ color: '#EC4899', name: 'Rose' },
							{ color: '#10B981', name: 'Emerald' },
							{ color: '#F59E0B', name: 'Orange' },
							{ color: '#3B82F6', name: 'Blue' },
							{ color: '#EF4444', name: 'Red' }
						] as preset}
							<button
								onclick={() => solidColor = preset.color}
								class="relative w-7 h-7 rounded-full border-2 transition-all hover:scale-125 {solidColor.toUpperCase() === preset.color ? 'border-indigo-600 ring-2 ring-indigo-200' : 'border-gray-300 hover:border-gray-400'}"
								style="background: {preset.color};"
								title={preset.name}
							>
								{#if solidColor.toUpperCase() === preset.color}
									<div class="absolute inset-0 flex items-center justify-center">
										<svg class="w-3 h-3 {['#FFFFFF', '#F8FAFC', '#FEF3C7', '#DBEAFE', '#E0E7FF', '#FCE7F3', '#DCFCE7'].includes(preset.color) ? 'text-gray-700' : 'text-white'}" fill="currentColor" viewBox="0 0 20 20">
											<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
										</svg>
									</div>
								{/if}
							</button>
						{/each}
					</div>
				</div>
				
				<!-- Preview -->
				<div class="w-full h-24 rounded-lg border border-gray-200 shadow-sm" style="background: {solidColor};"></div>
			</div>
		{:else if manuallySelectedType === 'gradient'}
			<div class="space-y-4">
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-2">From Color</label>
						<div class="flex gap-2 items-center">
							<input
								type="color"
								value={gradientFrom}
								oninput={(e) => gradientFrom = normalizeHex(e.currentTarget.value)}
								class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer"
								style="padding: 0; -webkit-appearance: none; appearance: none;"
							/>
							<input
								type="text"
								value={gradientFrom}
								oninput={(e) => {
									const val = e.currentTarget.value.toUpperCase();
									if (/^#[0-9A-F]{6}$/.test(val)) gradientFrom = val;
									else if (/^#[0-9A-F]{0,6}$/.test(val)) e.currentTarget.value = val;
								}}
								class="flex-1 px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-colors"
								placeholder="#FAF5FF"
								maxlength="7"
							/>
						</div>
					</div>
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-2">To Color</label>
						<div class="flex gap-2 items-center">
							<input
								type="color"
								value={gradientTo}
								oninput={(e) => gradientTo = normalizeHex(e.currentTarget.value)}
								class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer"
								style="padding: 0; -webkit-appearance: none; appearance: none;"
							/>
							<input
								type="text"
								value={gradientTo}
								oninput={(e) => {
									const val = e.currentTarget.value.toUpperCase();
									if (/^#[0-9A-F]{6}$/.test(val)) gradientTo = val;
									else if (/^#[0-9A-F]{0,6}$/.test(val)) e.currentTarget.value = val;
								}}
								class="flex-1 px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-colors"
								placeholder="#EFF6FF"
								maxlength="7"
							/>
						</div>
					</div>
				</div>
				<div class="w-full h-24 rounded-lg border border-gray-200 shadow-sm" style="background: linear-gradient(to bottom right, {gradientFrom}, {gradientTo});"></div>
			</div>
		{:else if manuallySelectedType === 'image'}
			<div class="space-y-3">
				<div class="flex gap-2">
					<input
						type="text"
						bind:value={imageUrl}
						class="flex-1 px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-colors"
						placeholder="Paste image URL or upload..."
					/>
					<input
						type="file"
						accept="image/*"
						onchange={(e) => uploadFile(e, 'image')}
						class="hidden"
						id="bg-image-upload"
					/>
					<label
						for="bg-image-upload"
						class="flex items-center gap-2 px-4 py-2 bg-white border border-gray-200 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-50 hover:border-indigo-500 hover:text-indigo-600 cursor-pointer transition-all whitespace-nowrap"
					>
						{#if uploading}
							<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
						{:else}
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
							</svg>
						{/if}
						{uploading ? 'Uploading...' : 'Upload'}
					</label>
				</div>
				{#if imageUrl}
					<div class="relative rounded-lg overflow-hidden border border-gray-200 shadow-sm group">
						<img src={imageUrl} alt="Background preview" class="w-full h-32 object-cover" onerror={(e) => e.target.style.display = 'none'} />
						<button
							onclick={() => imageUrl = ''}
							class="absolute top-2 right-2 w-6 h-6 bg-black/50 hover:bg-black/70 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
							title="Remove image"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
							</svg>
						</button>
					</div>
				{/if}
			</div>
		{:else if manuallySelectedType === 'video'}
			<div class="space-y-3">
				<div class="flex gap-2">
					<input
						type="text"
						bind:value={videoUrl}
						class="flex-1 px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-colors"
						placeholder="Paste video URL or upload..."
					/>
					<input
						type="file"
						accept="video/*"
						onchange={(e) => uploadFile(e, 'video')}
						class="hidden"
						id="bg-video-upload"
					/>
					<label
						for="bg-video-upload"
						class="flex items-center gap-2 px-4 py-2 bg-white border border-gray-200 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-50 hover:border-indigo-500 hover:text-indigo-600 cursor-pointer transition-all whitespace-nowrap"
					>
						{#if uploading}
							<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
						{:else}
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4a3 3 0 016 0v4a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z"/>
							</svg>
						{/if}
						{uploading ? 'Uploading...' : 'Upload'}
					</label>
				</div>
				{#if videoUrl}
					<div class="relative rounded-lg overflow-hidden border border-gray-200 shadow-sm group">
						<video src={videoUrl} class="w-full h-32 object-cover" muted loop autoplay></video>
						<button
							onclick={() => videoUrl = ''}
							class="absolute top-2 right-2 w-6 h-6 bg-black/50 hover:bg-black/70 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
							title="Remove video"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
							</svg>
						</button>
					</div>
				{/if}
			</div>
		{/if}

		<!-- Apply Button (only show when manually selected) -->
		{#if manuallySelectedType}
			<button
				onclick={applyBackground}
				disabled={uploading || saving}
				class="w-full mt-6 px-5 py-3 bg-indigo-600 text-white rounded-xl text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-sm hover:shadow-md flex items-center justify-center gap-2"
			>
				{#if saving}
					<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Saving...
				{:else}
					Apply Background
				{/if}
			</button>
		{/if}
	</div>
</div>


<style>
	/* Custom color picker styling */
	input[type="color"] {
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background-color: transparent;
		border: none;
		cursor: pointer;
	}
	
	input[type="color"]::-webkit-color-swatch-wrapper {
		padding: 0;
		border-radius: 50%;
		overflow: hidden;
	}
	
	input[type="color"]::-webkit-color-swatch {
		border: none;
		border-radius: 50%;
	}
	
	input[type="color"]::-moz-color-swatch {
		border: none;
		border-radius: 50%;
	}
</style>
