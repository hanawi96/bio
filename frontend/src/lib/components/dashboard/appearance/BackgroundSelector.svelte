<script lang="ts">
	import { globalTheme } from '$lib/stores/theme';
	import { pendingChanges } from '$lib/stores/pendingChanges';

	let selectedType = $state<'solid' | 'gradient' | 'image' | 'video'>('solid');
	let solidColor = $state('#ffffff');
	let gradientFrom = $state('#faf5ff');
	let gradientTo = $state('#eff6ff');
	let gradientDirection = $state<'up' | 'down' | 'radial'>('up');
	let imageUrl = $state('');
	let videoUrl = $state('');
	let uploading = $state(false);
	let showGradientCustom = $state(false);

	// Ensure color is always in HEX format
	function normalizeHex(color: string): string {
		if (!color) return '#ffffff';
		if (color.startsWith('#')) return color.toUpperCase();
		return `#${color}`.toUpperCase();
	}

	let isSyncing = false;

	// Sync from theme store (when theme preset changes)
	$effect(() => {
		const theme = $globalTheme;
		isSyncing = true;
		selectedType = theme.pageBackgroundType || 'solid';
		solidColor = normalizeHex(theme.pageBackground || '#ffffff');
		gradientFrom = normalizeHex(theme.pageGradientFrom || '#faf5ff');
		gradientTo = normalizeHex(theme.pageGradientTo || '#eff6ff');
		gradientDirection = theme.pageGradientDirection || 'up';
		imageUrl = theme.pageBackgroundImage || '';
		videoUrl = theme.pageBackgroundVideo || '';
		setTimeout(() => isSyncing = false, 0);
	});

	// Update theme on user changes only
	function updateTheme() {
		if (!isSyncing) {
			const updates: any = {
				pageBackgroundType: selectedType,
				pageBackground: solidColor,
				pageGradientFrom: gradientFrom,
				pageGradientTo: gradientTo,
				pageGradientDirection: gradientDirection,
				pageBackgroundImage: imageUrl,
				pageBackgroundVideo: videoUrl
			};
			globalTheme.update(updates);
			pendingChanges.updateTheme(updates);
		}
	}

	$effect(() => {
		selectedType;
		solidColor;
		gradientFrom;
		gradientTo;
		gradientDirection;
		imageUrl;
		videoUrl;
		updateTheme();
	});

	async function uploadFile(e: Event, type: 'image' | 'video') {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		uploading = true;
		try {
			const formData = new FormData();
			formData.append('file', file);
			
			// Get token from localStorage directly
			const token = localStorage.getItem('token');
			const response = await fetch('/api/upload', {
				method: 'POST',
				headers: { Authorization: `Bearer ${token}` },
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
		} catch (e: any) {
			console.error('Upload failed:', e);
		} finally {
			uploading = false;
		}
	}
</script>

<div class="space-y-6">
	<div>
		<h3 class="text-lg font-bold text-gray-900 mb-4">Wallpaper</h3>
		
		<!-- Type Selector -->
		<div class="grid grid-cols-5 gap-3 mb-6">
			<button
				onclick={() => selectedType = 'solid'}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {selectedType === 'solid' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
			>
				<div class="w-12 h-12 rounded-lg bg-gradient-to-br from-indigo-100 to-indigo-200 flex items-center justify-center">
					<div class="w-6 h-6 rounded bg-indigo-600"></div>
				</div>
				<span class="text-xs font-medium text-gray-700">Solid</span>
			</button>

			<button
				onclick={() => selectedType = 'gradient'}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {selectedType === 'gradient' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
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
				onclick={() => selectedType = 'image'}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {selectedType === 'image' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
			>
				<div class="w-12 h-12 rounded-lg bg-gray-200 flex items-center justify-center">
					<svg class="w-6 h-6 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd"/>
					</svg>
				</div>
				<span class="text-xs font-medium text-gray-700">Image</span>
			</button>

			<button
				onclick={() => selectedType = 'video'}
				class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all {selectedType === 'video' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
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

		<!-- Settings -->
		{#if selectedType === 'solid'}
			<div class="space-y-4">
				<div class="flex gap-3 items-center">
					<!-- Color Picker & Hex Input -->
					<div class="flex gap-2 items-center flex-shrink-0">
						<div class="w-10 h-10 rounded-lg border-2 border-gray-300 overflow-hidden bg-white">
							<input
								type="color"
								value={solidColor}
								oninput={(e) => solidColor = normalizeHex(e.currentTarget.value)}
								class="w-full h-full cursor-pointer block"
								style="-webkit-appearance: none; appearance: none; border: none; padding: 0; margin: 0;"
							/>
						</div>
						<input
							type="text"
							value={solidColor}
							oninput={(e) => {
								const val = e.currentTarget.value.toUpperCase();
								if (/^#[0-9A-F]{6}$/.test(val)) solidColor = val;
								else if (/^#[0-9A-F]{0,6}$/.test(val)) e.currentTarget.value = val;
							}}
							class="w-24 px-2 py-2 bg-gray-50 border border-gray-200 rounded-lg text-xs font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:bg-white transition-colors"
							placeholder="#FFFFFF"
							maxlength="7"
						/>
					</div>
					
					<!-- Color Presets -->
					<div class="flex-1 flex flex-wrap gap-1.5 items-center">
						{#each [
							{ color: '#FFFFFF', name: 'White' },
							{ color: '#F8FAFC', name: 'Slate' },
							{ color: '#FEF3C7', name: 'Cream' },
							{ color: '#DBEAFE', name: 'Sky' },
							{ color: '#E0E7FF', name: 'Indigo' },
							{ color: '#FCE7F3', name: 'Pink' },
							{ color: '#DCFCE7', name: 'Green' },
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
								class="relative w-7 h-7 rounded-full border-2 transition-all hover:scale-110 flex-shrink-0 {solidColor.toUpperCase() === preset.color ? 'border-indigo-600 ring-2 ring-indigo-200' : 'border-gray-300 hover:border-gray-400'}"
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
			</div>
		{:else if selectedType === 'gradient'}
			<div class="space-y-4">
				<!-- Gradient Presets -->
				<div class="flex flex-wrap gap-1.5 items-center">
					{#each [
						{ from: '#667EEA', to: '#764BA2', name: 'Purple Dream' },
						{ from: '#F093FB', to: '#F5576C', name: 'Pink Sunset' },
						{ from: '#4FACFE', to: '#00F2FE', name: 'Ocean Blue' },
						{ from: '#43E97B', to: '#38F9D7', name: 'Mint Fresh' },
						{ from: '#FA709A', to: '#FEE140', name: 'Peach Glow' },
						{ from: '#30CFD0', to: '#330867', name: 'Deep Sea' },
						{ from: '#A8EDEA', to: '#FED6E3', name: 'Cotton Candy' },
						{ from: '#FF9A9E', to: '#FAD0C4', name: 'Soft Coral' },
						{ from: '#FBC2EB', to: '#A6C1EE', name: 'Lavender Sky' },
						{ from: '#FDCBF1', to: '#E6DEE9', name: 'Pastel Pink' },
						{ from: '#A1C4FD', to: '#C2E9FB', name: 'Sky Blue' },
						{ from: '#D299C2', to: '#FEF9D7', name: 'Spring Bloom' },
						{ from: '#FEE140', to: '#FA709A', name: 'Sunset Vibes' },
						{ from: '#CE9FFC', to: '#7367F0', name: 'Royal Purple' },
						{ from: '#90F7EC', to: '#32CCBC', name: 'Aqua Marine' },
						{ from: '#FFF1EB', to: '#ACE0F9', name: 'Soft Sky' }
					] as preset}
						<button
							onclick={() => {
								gradientFrom = preset.from;
								gradientTo = preset.to;
								showGradientCustom = false;
							}}
							class="relative w-6 h-10 rounded border-2 transition-all hover:scale-110 flex-shrink-0 {gradientFrom.toUpperCase() === preset.from && gradientTo.toUpperCase() === preset.to && !showGradientCustom ? 'border-indigo-600 ring-2 ring-indigo-200' : 'border-gray-300 hover:border-gray-400'}"
							style="background: linear-gradient(to bottom, {preset.from}, {preset.to});"
							title={preset.name}
						>
							{#if gradientFrom.toUpperCase() === preset.from && gradientTo.toUpperCase() === preset.to && !showGradientCustom}
								<div class="absolute inset-0 flex items-center justify-center">
									<svg class="w-3 h-3 text-white drop-shadow-lg" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
									</svg>
								</div>
							{/if}
						</button>
					{/each}
					
					<!-- Custom Button -->
					<button
						onclick={() => showGradientCustom = !showGradientCustom}
						class="w-6 h-10 rounded border-2 transition-all hover:scale-110 flex-shrink-0 flex items-center justify-center {showGradientCustom ? 'border-indigo-600 bg-indigo-50' : 'border-gray-300 hover:border-gray-400 bg-gray-50'}"
						title="Custom Gradient"
					>
						<svg class="w-4 h-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
						</svg>
					</button>
				</div>
				
				<!-- Direction Selector - Always visible -->
				<div>
					<div class="grid grid-cols-3 gap-2">
						<button
							onclick={() => gradientDirection = 'up'}
							class="flex flex-col items-center gap-1 p-2 rounded-lg border-2 transition-all {gradientDirection === 'up' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
						>
							<svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
							</svg>
							<span class="text-xs font-medium text-gray-700">Up</span>
						</button>
						<button
							onclick={() => gradientDirection = 'down'}
							class="flex flex-col items-center gap-1 p-2 rounded-lg border-2 transition-all {gradientDirection === 'down' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
						>
							<svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"/>
							</svg>
							<span class="text-xs font-medium text-gray-700">Down</span>
						</button>
						<button
							onclick={() => gradientDirection = 'radial'}
							class="flex flex-col items-center gap-1 p-2 rounded-lg border-2 transition-all {gradientDirection === 'radial' ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'}"
						>
							<svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<circle cx="12" cy="12" r="10" stroke-width="2"/>
								<circle cx="12" cy="12" r="6" stroke-width="2"/>
								<circle cx="12" cy="12" r="2" stroke-width="2"/>
							</svg>
							<span class="text-xs font-medium text-gray-700">Radial</span>
						</button>
					</div>
				</div>
				
				{#if showGradientCustom}
					<div class="space-y-3 p-4 bg-gray-50 rounded-lg">
						<div class="grid grid-cols-2 gap-4">
							<div>
								<label class="block text-xs font-medium text-gray-700 mb-2">From Color</label>
								<div class="flex gap-2 items-center">
									<div class="w-8 h-8 rounded border-2 border-gray-300 overflow-hidden bg-white flex-shrink-0">
										<input
											type="color"
											value={gradientFrom}
											oninput={(e) => gradientFrom = normalizeHex(e.currentTarget.value)}
											class="w-full h-full cursor-pointer block"
											style="-webkit-appearance: none; appearance: none; border: none; padding: 0; margin: 0;"
										/>
									</div>
									<input
										type="text"
										value={gradientFrom}
										oninput={(e) => {
											const val = e.currentTarget.value.toUpperCase();
											if (/^#[0-9A-F]{6}$/.test(val)) gradientFrom = val;
											else if (/^#[0-9A-F]{0,6}$/.test(val)) e.currentTarget.value = val;
										}}
										class="flex-1 px-2 py-1.5 bg-white border border-gray-200 rounded text-xs font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
										placeholder="#FAF5FF"
										maxlength="7"
									/>
								</div>
							</div>
							<div>
								<label class="block text-xs font-medium text-gray-700 mb-2">To Color</label>
								<div class="flex gap-2 items-center">
									<div class="w-8 h-8 rounded border-2 border-gray-300 overflow-hidden bg-white flex-shrink-0">
										<input
											type="color"
											value={gradientTo}
											oninput={(e) => gradientTo = normalizeHex(e.currentTarget.value)}
											class="w-full h-full cursor-pointer block"
											style="-webkit-appearance: none; appearance: none; border: none; padding: 0; margin: 0;"
										/>
									</div>
									<input
										type="text"
										value={gradientTo}
										oninput={(e) => {
											const val = e.currentTarget.value.toUpperCase();
											if (/^#[0-9A-F]{6}$/.test(val)) gradientTo = val;
											else if (/^#[0-9A-F]{0,6}$/.test(val)) e.currentTarget.value = val;
										}}
										class="flex-1 px-2 py-1.5 bg-white border border-gray-200 rounded text-xs font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
										placeholder="#EFF6FF"
										maxlength="7"
									/>
								</div>
							</div>
						</div>
					</div>
				{/if}
			</div>
		{:else if selectedType === 'image'}
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
		{:else if selectedType === 'video'}
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
