<script lang="ts">
	import { globalTheme } from '$lib/stores/theme';

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
			<div class="space-y-3">
				<!-- iOS-style Segmented Control for Direction -->
				<div class="inline-flex p-0.5 bg-gray-100 rounded-lg">
					<button
						onclick={() => gradientDirection = 'up'}
						class="px-4 py-1.5 rounded-md text-xs font-semibold transition-all {gradientDirection === 'up' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
					>
						↑ Up
					</button>
					<button
						onclick={() => gradientDirection = 'down'}
						class="px-4 py-1.5 rounded-md text-xs font-semibold transition-all {gradientDirection === 'down' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
					>
						↓ Down
					</button>
					<button
						onclick={() => gradientDirection = 'radial'}
						class="px-4 py-1.5 rounded-md text-xs font-semibold transition-all {gradientDirection === 'radial' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
					>
						◉ Radial
					</button>
				</div>

				<!-- Toggle between Presets and Custom -->
				{#if !showGradientCustom}
					<!-- Presets Grid - iOS Style -->
					<div>
						<div class="flex items-center justify-between mb-2">
							<span class="text-xs font-medium text-gray-600">Presets</span>
							<button
								onclick={() => showGradientCustom = true}
								class="flex items-center gap-1 px-2 py-1 rounded-md text-xs font-semibold text-indigo-600 hover:bg-indigo-50 transition-all"
							>
								<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
								</svg>
								Custom
							</button>
						</div>
						<div class="grid grid-cols-8 gap-2">
							{#each [
								{ from: '#667EEA', to: '#764BA2', name: 'Purple' },
								{ from: '#F093FB', to: '#F5576C', name: 'Pink' },
								{ from: '#4FACFE', to: '#00F2FE', name: 'Ocean' },
								{ from: '#43E97B', to: '#38F9D7', name: 'Mint' },
								{ from: '#FA709A', to: '#FEE140', name: 'Peach' },
								{ from: '#30CFD0', to: '#330867', name: 'Deep' },
								{ from: '#A8EDEA', to: '#FED6E3', name: 'Cotton' },
								{ from: '#FF9A9E', to: '#FAD0C4', name: 'Coral' },
								{ from: '#FBC2EB', to: '#A6C1EE', name: 'Lavender' },
								{ from: '#FDCBF1', to: '#E6DEE9', name: 'Pastel' },
								{ from: '#A1C4FD', to: '#C2E9FB', name: 'Sky' },
								{ from: '#D299C2', to: '#FEF9D7', name: 'Spring' },
								{ from: '#FEE140', to: '#FA709A', name: 'Sunset' },
								{ from: '#CE9FFC', to: '#7367F0', name: 'Royal' },
								{ from: '#90F7EC', to: '#32CCBC', name: 'Aqua' },
								{ from: '#FFF1EB', to: '#ACE0F9', name: 'Soft' }
							] as preset}
								<button
									onclick={() => {
										gradientFrom = preset.from;
										gradientTo = preset.to;
									}}
									class="relative aspect-square rounded-lg border-2 transition-all hover:scale-105 {gradientFrom.toUpperCase() === preset.from && gradientTo.toUpperCase() === preset.to ? 'border-indigo-500 ring-2 ring-indigo-200' : 'border-gray-200 hover:border-gray-300'}"
									style="background: linear-gradient(135deg, {preset.from}, {preset.to});"
									title={preset.name}
								>
									{#if gradientFrom.toUpperCase() === preset.from && gradientTo.toUpperCase() === preset.to}
										<div class="absolute inset-0 flex items-center justify-center">
											<div class="w-5 h-5 bg-white rounded-full shadow-lg flex items-center justify-center">
												<svg class="w-3 h-3 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
													<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
												</svg>
											</div>
										</div>
									{/if}
								</button>
							{/each}
						</div>
					</div>
				{:else}
					<!-- iOS-style Custom Color Picker -->
					<div class="space-y-3">
						<!-- Back to Presets Button -->
						<button
							onclick={() => showGradientCustom = false}
							class="flex items-center gap-1.5 text-xs font-semibold text-indigo-600 hover:text-indigo-700 transition-colors"
						>
							<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
							</svg>
							Back to Presets
						</button>

						<!-- Gradient Preview Strip -->
						<div 
							class="w-full h-20 rounded-2xl border-2 border-gray-200 shadow-sm overflow-hidden"
							style="background: {gradientDirection === 'radial' 
								? `radial-gradient(circle at center, ${gradientFrom}, ${gradientTo})` 
								: gradientDirection === 'down'
								? `linear-gradient(to bottom, ${gradientFrom}, ${gradientTo})`
								: `linear-gradient(to top, ${gradientFrom}, ${gradientTo})`};"
						></div>

						<!-- Color Pickers - iOS Style -->
						<div class="grid grid-cols-2 gap-4">
							<!-- Start Color -->
							<div class="text-center">
								<div class="relative inline-block">
									<div class="w-20 h-20 rounded-full border-4 border-white shadow-lg overflow-hidden cursor-pointer hover:scale-105 transition-transform mx-auto">
										<input
											type="color"
											value={gradientFrom}
											oninput={(e) => gradientFrom = normalizeHex(e.currentTarget.value)}
											class="w-full h-full cursor-pointer"
											style="-webkit-appearance: none; appearance: none; border: none;"
										/>
									</div>
								</div>
								<div class="mt-2 text-xs font-medium text-gray-500">Start Color</div>
								<div class="mt-1 text-sm font-mono font-bold text-gray-900">{gradientFrom}</div>
							</div>

							<!-- End Color -->
							<div class="text-center">
								<div class="relative inline-block">
									<div class="w-20 h-20 rounded-full border-4 border-white shadow-lg overflow-hidden cursor-pointer hover:scale-105 transition-transform mx-auto">
										<input
											type="color"
											value={gradientTo}
											oninput={(e) => gradientTo = normalizeHex(e.currentTarget.value)}
											class="w-full h-full cursor-pointer"
											style="-webkit-appearance: none; appearance: none; border: none;"
										/>
									</div>
								</div>
								<div class="mt-2 text-xs font-medium text-gray-500">End Color</div>
								<div class="mt-1 text-sm font-mono font-bold text-gray-900">{gradientTo}</div>
							</div>
						</div>
					</div>
				{/if}
			</div>
		{:else if selectedType === 'image'}
			<div class="space-y-4">
				{#if !imageUrl}
					<!-- iOS-style Upload Area -->
					<div class="relative">
						<input
							type="file"
							accept="image/*"
							onchange={(e) => uploadFile(e, 'image')}
							class="hidden"
							id="bg-image-upload"
						/>
						<label
							for="bg-image-upload"
							class="block relative overflow-hidden rounded-2xl border-2 border-dashed border-gray-300 hover:border-indigo-400 transition-all cursor-pointer group bg-gradient-to-br from-gray-50 to-gray-100/50 hover:from-indigo-50 hover:to-purple-50"
						>
							<div class="p-12 flex flex-col items-center justify-center text-center">
								{#if uploading}
									<!-- Uploading State -->
									<div class="w-20 h-20 rounded-full bg-indigo-100 flex items-center justify-center mb-4 animate-pulse">
										<svg class="w-10 h-10 text-indigo-600 animate-spin" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
									</div>
									<p class="text-base font-semibold text-indigo-600 mb-1">Uploading...</p>
									<p class="text-sm text-gray-500">Please wait</p>
								{:else}
									<!-- Upload Icon -->
									<div class="w-20 h-20 rounded-full bg-gradient-to-br from-indigo-100 to-purple-100 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
										<svg class="w-10 h-10 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
										</svg>
									</div>
									
									<!-- Text -->
									<p class="text-base font-semibold text-gray-900 mb-1">Choose an image</p>
									<p class="text-sm text-gray-500 mb-4">or drag and drop here</p>
									
									<!-- File Info -->
									<div class="inline-flex items-center gap-2 px-4 py-2 bg-white/80 backdrop-blur-sm rounded-full border border-gray-200 shadow-sm">
										<svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
										</svg>
										<span class="text-xs text-gray-600 font-medium">JPG, PNG, GIF • Max 10MB</span>
									</div>
								{/if}
							</div>
						</label>
					</div>
					
					<!-- Divider with "OR" -->
					<div class="relative flex items-center">
						<div class="flex-1 border-t border-gray-200"></div>
						<span class="px-4 text-xs font-medium text-gray-500 bg-white">OR</span>
						<div class="flex-1 border-t border-gray-200"></div>
					</div>
					
					<!-- URL Input -->
					<div class="space-y-2">
						<label class="block text-xs font-medium text-gray-700 px-1">Paste image URL</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
								</svg>
							</div>
							<input
								type="text"
								bind:value={imageUrl}
								class="w-full pl-10 pr-4 py-3 bg-white border border-gray-200 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all placeholder:text-gray-400"
								placeholder="https://example.com/image.jpg"
							/>
						</div>
					</div>
				{:else}
					<!-- Preview with iOS-style Card -->
					<div class="space-y-3">
						<div class="relative rounded-2xl overflow-hidden border border-gray-200 shadow-lg group bg-white">
							<!-- Image Preview -->
							<div class="relative aspect-video bg-gray-100">
								<img 
									src={imageUrl} 
									alt="Background preview" 
									class="w-full h-full object-cover" 
									onerror={(e) => {
										e.target.style.display = 'none';
										e.target.parentElement.innerHTML = '<div class="w-full h-full flex items-center justify-center"><svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg></div>';
									}} 
								/>
								
								<!-- Overlay Actions -->
								<div class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/0 to-black/20 opacity-0 group-hover:opacity-100 transition-opacity">
									<div class="absolute top-3 right-3 flex gap-2">
										<!-- Change Button -->
										<input
											type="file"
											accept="image/*"
											onchange={(e) => uploadFile(e, 'image')}
											class="hidden"
											id="bg-image-change"
										/>
										<label
											for="bg-image-change"
											class="flex items-center gap-1.5 px-3 py-2 bg-white/95 backdrop-blur-sm hover:bg-white text-gray-700 rounded-lg text-xs font-medium cursor-pointer transition-all shadow-lg"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
											</svg>
											Change
										</label>
										
										<!-- Remove Button -->
										<button
											onclick={() => imageUrl = ''}
											class="flex items-center gap-1.5 px-3 py-2 bg-red-500/95 backdrop-blur-sm hover:bg-red-600 text-white rounded-lg text-xs font-medium transition-all shadow-lg"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
											</svg>
											Remove
										</button>
									</div>
								</div>
							</div>
							
							<!-- Info Bar -->
							<div class="px-4 py-3 bg-gray-50 border-t border-gray-100">
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2">
										<div class="w-8 h-8 rounded-lg bg-indigo-100 flex items-center justify-center">
											<svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
											</svg>
										</div>
										<div>
											<p class="text-xs font-semibold text-gray-900">Image uploaded</p>
											<p class="text-xs text-gray-500">Background is ready</p>
										</div>
									</div>
									<div class="flex items-center gap-1 px-2 py-1 bg-green-100 rounded-full">
										<div class="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
										<span class="text-xs font-medium text-green-700">Active</span>
									</div>
								</div>
							</div>
						</div>
						
						<!-- URL Display (collapsible) -->
						<details class="group/details">
							<summary class="flex items-center justify-between px-4 py-2.5 bg-gray-50 hover:bg-gray-100 rounded-lg cursor-pointer transition-colors">
								<span class="text-xs font-medium text-gray-700">Image URL</span>
								<svg class="w-4 h-4 text-gray-500 transition-transform group-open/details:rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
								</svg>
							</summary>
							<div class="mt-2 px-4 py-3 bg-gray-50 rounded-lg border border-gray-200">
								<div class="flex items-center gap-2">
									<input
										type="text"
										bind:value={imageUrl}
										class="flex-1 px-3 py-2 bg-white border border-gray-200 rounded-lg text-xs font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
										placeholder="https://..."
									/>
									<button
										onclick={() => {
											navigator.clipboard.writeText(imageUrl);
										}}
										class="p-2 hover:bg-gray-200 rounded-lg transition-colors"
										title="Copy URL"
									>
										<svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
										</svg>
									</button>
								</div>
							</div>
						</details>
					</div>
				{/if}
			</div>
		{:else if selectedType === 'video'}
			<div class="space-y-4">
				{#if !videoUrl}
					<!-- iOS-style Upload Area -->
					<div class="relative">
						<input
							type="file"
							accept="video/*"
							onchange={(e) => uploadFile(e, 'video')}
							class="hidden"
							id="bg-video-upload"
						/>
						<label
							for="bg-video-upload"
							class="block relative overflow-hidden rounded-2xl border-2 border-dashed border-gray-300 hover:border-indigo-400 transition-all cursor-pointer group bg-gradient-to-br from-gray-50 to-gray-100/50 hover:from-indigo-50 hover:to-purple-50"
						>
							<div class="p-12 flex flex-col items-center justify-center text-center">
								{#if uploading}
									<!-- Uploading State -->
									<div class="w-20 h-20 rounded-full bg-indigo-100 flex items-center justify-center mb-4 animate-pulse">
										<svg class="w-10 h-10 text-indigo-600 animate-spin" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
									</div>
									<p class="text-base font-semibold text-indigo-600 mb-1">Uploading...</p>
									<p class="text-sm text-gray-500">Please wait</p>
								{:else}
									<!-- Upload Icon -->
									<div class="w-20 h-20 rounded-full bg-gradient-to-br from-purple-100 to-pink-100 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
										<svg class="w-10 h-10 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
										</svg>
									</div>
									
									<!-- Text -->
									<p class="text-base font-semibold text-gray-900 mb-1">Choose a video</p>
									<p class="text-sm text-gray-500 mb-4">or drag and drop here</p>
									
									<!-- File Info -->
									<div class="inline-flex items-center gap-2 px-4 py-2 bg-white/80 backdrop-blur-sm rounded-full border border-gray-200 shadow-sm">
										<svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
										</svg>
										<span class="text-xs text-gray-600 font-medium">MP4, WebM • Max 50MB</span>
									</div>
								{/if}
							</div>
						</label>
					</div>
					
					<!-- Divider with "OR" -->
					<div class="relative flex items-center">
						<div class="flex-1 border-t border-gray-200"></div>
						<span class="px-4 text-xs font-medium text-gray-500 bg-white">OR</span>
						<div class="flex-1 border-t border-gray-200"></div>
					</div>
					
					<!-- URL Input -->
					<div class="space-y-2">
						<label class="block text-xs font-medium text-gray-700 px-1">Paste video URL</label>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
								<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
								</svg>
							</div>
							<input
								type="text"
								bind:value={videoUrl}
								class="w-full pl-10 pr-4 py-3 bg-white border border-gray-200 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all placeholder:text-gray-400"
								placeholder="https://example.com/video.mp4"
							/>
						</div>
					</div>
				{:else}
					<!-- Preview with iOS-style Card -->
					<div class="space-y-3">
						<div class="relative rounded-2xl overflow-hidden border border-gray-200 shadow-lg group bg-white">
							<!-- Video Preview -->
							<div class="relative aspect-video bg-gray-100">
								<video 
									src={videoUrl} 
									class="w-full h-full object-cover" 
									muted 
									loop 
									autoplay
									onerror={(e) => {
										e.target.style.display = 'none';
										e.target.parentElement.innerHTML = '<div class="w-full h-full flex items-center justify-center"><svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg></div>';
									}}
								></video>
								
								<!-- Play Indicator -->
								<div class="absolute bottom-3 left-3 flex items-center gap-2 px-3 py-1.5 bg-black/60 backdrop-blur-sm rounded-full">
									<div class="w-2 h-2 bg-red-500 rounded-full animate-pulse"></div>
									<span class="text-xs font-medium text-white">Playing</span>
								</div>
								
								<!-- Overlay Actions -->
								<div class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/0 to-black/20 opacity-0 group-hover:opacity-100 transition-opacity">
									<div class="absolute top-3 right-3 flex gap-2">
										<!-- Change Button -->
										<input
											type="file"
											accept="video/*"
											onchange={(e) => uploadFile(e, 'video')}
											class="hidden"
											id="bg-video-change"
										/>
										<label
											for="bg-video-change"
											class="flex items-center gap-1.5 px-3 py-2 bg-white/95 backdrop-blur-sm hover:bg-white text-gray-700 rounded-lg text-xs font-medium cursor-pointer transition-all shadow-lg"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
											</svg>
											Change
										</label>
										
										<!-- Remove Button -->
										<button
											onclick={() => videoUrl = ''}
											class="flex items-center gap-1.5 px-3 py-2 bg-red-500/95 backdrop-blur-sm hover:bg-red-600 text-white rounded-lg text-xs font-medium transition-all shadow-lg"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
											</svg>
											Remove
										</button>
									</div>
								</div>
							</div>
							
							<!-- Info Bar -->
							<div class="px-4 py-3 bg-gray-50 border-t border-gray-100">
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2">
										<div class="w-8 h-8 rounded-lg bg-purple-100 flex items-center justify-center">
											<svg class="w-4 h-4 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
											</svg>
										</div>
										<div>
											<p class="text-xs font-semibold text-gray-900">Video uploaded</p>
											<p class="text-xs text-gray-500">Background is ready</p>
										</div>
									</div>
									<div class="flex items-center gap-1 px-2 py-1 bg-green-100 rounded-full">
										<div class="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
										<span class="text-xs font-medium text-green-700">Active</span>
									</div>
								</div>
							</div>
						</div>
						
						<!-- URL Display (collapsible) -->
						<details class="group/details">
							<summary class="flex items-center justify-between px-4 py-2.5 bg-gray-50 hover:bg-gray-100 rounded-lg cursor-pointer transition-colors">
								<span class="text-xs font-medium text-gray-700">Video URL</span>
								<svg class="w-4 h-4 text-gray-500 transition-transform group-open/details:rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
								</svg>
							</summary>
							<div class="mt-2 px-4 py-3 bg-gray-50 rounded-lg border border-gray-200">
								<div class="flex items-center gap-2">
									<input
										type="text"
										bind:value={videoUrl}
										class="flex-1 px-3 py-2 bg-white border border-gray-200 rounded-lg text-xs font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
										placeholder="https://..."
									/>
									<button
										onclick={() => {
											navigator.clipboard.writeText(videoUrl);
										}}
										class="p-2 hover:bg-gray-200 rounded-lg transition-colors"
										title="Copy URL"
									>
										<svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
										</svg>
									</button>
								</div>
							</div>
						</details>
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
