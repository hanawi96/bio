<script lang="ts">
	import { currentHeaderStyle, themePresets, type HeaderStyles } from '$lib/stores/theme';
	import { pendingChanges } from '$lib/stores/pendingChanges';
	import { onDestroy, createEventDispatcher } from 'svelte';
	import { toast } from 'svelte-sonner';
	
	const dispatch = createEventDispatcher();
	
	type CustomPreset = { id: string; name: string; layout: HeaderStyles['layout']; preview: { cover: string; avatar: string }; settings: Partial<HeaderStyles> };
	
	let { customPresets = $bindable([]), currentThemeName = 'default' } = $props();
	
	const headerStyle = $derived($currentHeaderStyle);
	let selectedCategory = $state<string>('all');
	let manuallySelected = $state<string | null>(null);
	let expandedPreset = $state<string | null>(null);
	let isModifiedFromPreset = $state(false);
	let activeTab = $state<string>('avatar');
	let showDeleteModal = $state(false);
	let presetToDelete = $state<string | null>(null);
	let isUserEditing = $state(false); // Track if user is actively editing (not auto-updates)
	let createPresetTimeout: ReturnType<typeof setTimeout> | null = null;
	
	function findMatchingCustomPreset(): string | null {
		if (!customPresets || customPresets.length === 0) return null;
		
		for (const preset of customPresets) {
			if (!preset.settings) continue;
			
			// Compare all important properties for exact match
			const matches = 
				preset.settings.layout === headerStyle.layout &&
				preset.settings.avatarSize === headerStyle.avatarSize &&
				preset.settings.avatarShape === headerStyle.avatarShape &&
				preset.settings.avatarAlign === headerStyle.avatarAlign &&
				preset.settings.avatarBorder === headerStyle.avatarBorder &&
				preset.settings.coverHeight === headerStyle.coverHeight &&
				preset.settings.showCover === headerStyle.showCover &&
				preset.settings.bioAlign === headerStyle.bioAlign &&
				preset.settings.bioSize === headerStyle.bioSize &&
				preset.settings.bioTextColor === headerStyle.bioTextColor;
			
			if (matches) return preset.id;
		}
		
		return null;
	}

	const categories = [
		{ id: 'all', label: 'All' },
		{ id: 'classic', label: 'Classic' },
		{ id: 'modern', label: 'Modern' },
		{ id: 'creative', label: 'Creative' },
		{ id: 'minimal', label: 'Minimal' },
		{ id: 'custom', label: 'Custom' }
	];

	const headerPresets: Record<string, Array<{ id: string; name: string; layout: HeaderStyles['layout']; preview: { cover: string; avatar: string } }>> = {
		classic: [
			{ id: 'centered', name: 'Centered', layout: 'centered', preview: { cover: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', avatar: '#ffffff' } },
			{ id: 'overlap', name: 'Overlap', layout: 'overlap', preview: { cover: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', avatar: '#ffffff' } },
			{ id: 'side', name: 'Side Profile', layout: 'side', preview: { cover: '#ffffff', avatar: '#6366f1' } }
		],
		modern: [
			{ id: 'card', name: 'Card Float', layout: 'card', preview: { cover: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)', avatar: '#ffffff' } },
			{ id: 'glass', name: 'Glass', layout: 'glass', preview: { cover: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)', avatar: 'rgba(255,255,255,0.2)' } },
			{ id: 'gradient', name: 'Gradient', layout: 'gradient', preview: { cover: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)', avatar: '#ffffff' } }
		],
		creative: [
			{ id: 'split', name: 'Split Screen', layout: 'split', preview: { cover: 'linear-gradient(90deg, #ff9a9e 0%, #fecfef 100%)', avatar: '#ffffff' } },
			{ id: 'asymmetric', name: 'Asymmetric', layout: 'asymmetric', preview: { cover: 'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)', avatar: '#ffffff' } },
			{ id: 'full', name: 'Full Bleed', layout: 'full', preview: { cover: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)', avatar: '#ffffff' } }
		],
		minimal: [
			{ id: 'minimal', name: 'Minimal', layout: 'minimal', preview: { cover: '#f5f5f5', avatar: '#000000' } }
		]
	};

	const layoutsWithCover = ['centered', 'overlap', 'card', 'glass', 'gradient', 'full', 'split', 'asymmetric'];
	const supportsCover = $derived(layoutsWithCover.includes(headerStyle.layout));

	const allPresets = $derived([...Object.values(headerPresets).flat(), ...customPresets]);
	
	const displayedPresets = $derived.by(() => {
		if (selectedCategory === 'all') return allPresets;
		if (selectedCategory === 'custom') return customPresets;
		return headerPresets[selectedCategory] || [];
	});

	const presetMap: Record<HeaderStyles['layout'], Partial<HeaderStyles>> = {
		centered: { layout: 'centered', coverType: 'gradient', coverGradientFrom: '#667eea', coverGradientTo: '#764ba2', coverHeight: 160, avatarSize: 120, avatarBorder: 4, avatarShape: 'circle', avatarAlign: 'center', showCover: true, bioAlign: 'center', bioSize: 'md', bioTextColor: '#6b7280' },
		overlap: { layout: 'overlap', coverType: 'gradient', coverGradientFrom: '#f093fb', coverGradientTo: '#f5576c', coverHeight: 180, avatarSize: 130, avatarBorder: 4, avatarShape: 'circle', avatarAlign: 'center', showCover: true, bioAlign: 'center', bioSize: 'md', bioTextColor: '#6b7280' },
		card: { layout: 'card', coverType: 'gradient', coverGradientFrom: '#4facfe', coverGradientTo: '#00f2fe', coverHeight: 140, avatarSize: 110, avatarBorder: 4, avatarShape: 'circle', avatarAlign: 'center', showCover: true, bioAlign: 'center', bioSize: 'md', bioTextColor: '#6b7280' },
		glass: { layout: 'glass', coverType: 'gradient', coverGradientFrom: '#fa709a', coverGradientTo: '#fee140', coverHeight: 170, avatarSize: 120, avatarBorder: 4, avatarBorderColor: 'rgba(255,255,255,0.2)', avatarShape: 'circle', avatarAlign: 'center', showCover: true, bioAlign: 'center', bioSize: 'md', bioTextColor: '#ffffff' },
		gradient: { layout: 'gradient', coverType: 'gradient', coverGradientFrom: '#30cfd0', coverGradientTo: '#330867', coverHeight: 190, avatarSize: 130, avatarBorder: 5, avatarShape: 'circle', avatarAlign: 'center', showCover: true, bioAlign: 'center', bioSize: 'lg', bioTextColor: '#ffffff' },
		minimal: { layout: 'minimal', coverType: 'color', coverColor: '#f5f5f5', coverHeight: 0, avatarSize: 90, avatarBorder: 2, avatarBorderColor: '#000000', avatarShape: 'circle', avatarAlign: 'left', showCover: false, bioAlign: 'left', bioSize: 'sm', bioTextColor: '#000000' },
		full: { layout: 'full', coverType: 'gradient', coverGradientFrom: '#a8edea', coverGradientTo: '#fed6e3', coverHeight: 240, avatarSize: 150, avatarBorder: 6, avatarShape: 'circle', avatarAlign: 'center', showCover: true, bioAlign: 'center', bioSize: 'lg', bioTextColor: '#6b7280' },
		side: { layout: 'side', coverType: 'color', coverColor: '#ffffff', coverHeight: 0, avatarSize: 90, avatarBorder: 2, avatarBorderColor: '#6366f1', avatarShape: 'circle', avatarAlign: 'left', showCover: false, bioAlign: 'left', bioSize: 'sm', bioTextColor: '#000000' },
		split: { layout: 'split', coverType: 'gradient', coverGradientFrom: '#ff9a9e', coverGradientTo: '#fecfef', coverHeight: 200, avatarSize: 120, avatarBorder: 4, avatarShape: 'circle', avatarAlign: 'center', showCover: true, bioAlign: 'center', bioSize: 'md', bioTextColor: '#6b7280' },
		asymmetric: { layout: 'asymmetric', coverType: 'gradient', coverGradientFrom: '#ffecd2', coverGradientTo: '#fcb69f', coverHeight: 220, avatarSize: 140, avatarBorder: 5, avatarShape: 'circle', avatarAlign: 'left', showCover: true, bioAlign: 'left', bioSize: 'lg', bioTextColor: '#6b7280' }
	};

	function selectLayout(layout: HeaderStyles['layout'], customSettings?: Partial<HeaderStyles>, customId?: string) {
		const settings = customSettings || presetMap[layout];
		const newStyle = { ...headerStyle, ...settings, avatarBorderColor: headerStyle.avatarBorderColor || '#ffffff' };
		
		manuallySelected = customId || layout;
		expandedPreset = null;
		isModifiedFromPreset = false;
		currentHeaderStyle.set(newStyle);
	}

	function toggleCustomize(id: string) {
		if (expandedPreset === id) {
			expandedPreset = null;
		} else {
			expandedPreset = id;
			activeTab = supportsCover ? 'cover' : 'avatar';
		}
	}

	function updateHeaderSetting(updates: Partial<HeaderStyles>) {
		const newStyle = { ...headerStyle, ...updates };
		isUserEditing = true; // Mark as user-initiated edit
		currentHeaderStyle.set(newStyle);
		
		// Auto-create custom preset if editing from a preset (only once)
		// Set flag FIRST to prevent race conditions
		if (!isModifiedFromPreset && manuallySelected && !manuallySelected.startsWith('custom-')) {
			isModifiedFromPreset = true; // Set immediately to prevent duplicate calls
			
			// Clear any pending timeout
			if (createPresetTimeout) {
				clearTimeout(createPresetTimeout);
			}
			
			// Use setTimeout to batch multiple rapid updates
			createPresetTimeout = setTimeout(() => {
				// Check if a custom preset with same config already exists
				const matchingCustomId = findMatchingCustomPreset();
				if (matchingCustomId) {
					// Switch to existing custom preset instead of creating new one
					manuallySelected = matchingCustomId;
					expandedPreset = matchingCustomId;
					toast.info('Switched to existing custom style');
				} else {
					// Create new custom preset
					saveAsCustom();
				}
				createPresetTimeout = null;
			}, 300); // Debounce to batch rapid updates
		}
		
		// Reset flag after a short delay
		setTimeout(() => {
			isUserEditing = false;
		}, 100);
	}

	function saveAsCustom() {
		// Check if exact same custom preset already exists
		const existingMatch = findMatchingCustomPreset();
		if (existingMatch) {
			// Don't create duplicate, just switch to existing
			manuallySelected = existingMatch;
			expandedPreset = existingMatch;
			toast.info('Using existing custom style');
			return;
		}
		
		const id = `custom-${Date.now()}`;
		const basePreset = allPresets.find(p => p.layout === headerStyle.layout);
		
		// Deep clone to avoid Svelte proxy issues
		const settingsSnapshot = JSON.parse(JSON.stringify(headerStyle));
		
		const customPreset: CustomPreset = {
			id,
			name: `Custom ${customPresets.length + 1}`,
			layout: headerStyle.layout,
			preview: basePreset?.preview || { cover: headerStyle.coverGradientFrom || '#667eea', avatar: '#ffffff' },
			settings: settingsSnapshot
		};
		
		customPresets = [...customPresets, customPreset];
		manuallySelected = id as any;
		expandedPreset = id;
		isModifiedFromPreset = false;
		
		toast.success('Custom header style created');
	}

	function confirmDelete(id: string) {
		presetToDelete = id;
		showDeleteModal = true;
	}

	async function deleteCustom() {
		if (!presetToDelete) return;
		
		const deletingId = presetToDelete;
		const isCurrentlySelected = manuallySelected === deletingId;
		
		console.log('🗑️ [deleteCustom] Deleting preset:', deletingId);
		console.log('🗑️ manuallySelected:', manuallySelected);
		console.log('🗑️ isCurrentlySelected:', isCurrentlySelected);
		
		showDeleteModal = false;
		presetToDelete = null;
		
		// Remove from local state
		customPresets = customPresets.filter(p => p.id !== deletingId);
		
		// If deleting currently selected preset, switch to theme's default header
		if (isCurrentlySelected) {
			console.log('🗑️ Switching to theme default header...');
			console.log('🗑️ currentThemeName:', currentThemeName);
			
			// Get theme's default header from themePresets
			const themePreset = themePresets[currentThemeName];
			console.log('🗑️ themePreset:', themePreset);
			console.log('🗑️ themePreset.header:', themePreset?.header);
			
			if (themePreset?.header) {
				const newStyle = { ...themePreset.header };
				console.log('🗑️ newStyle from theme:', newStyle);
				console.log('🗑️ newStyle.coverHeight:', newStyle.coverHeight);
				console.log('🗑️ newStyle.layout:', newStyle.layout);
				
				currentHeaderStyle.set(newStyle);
				manuallySelected = null;
				
				console.log('✅ Switched to theme default header');
				toast.success('Custom header style deleted. Switched to theme default.');
			} else {
				// Fallback to layout preset if theme has no header
				const currentLayout = headerStyle.layout;
				const defaultPreset = presetMap[currentLayout];
				console.log('⚠️ No theme header, using layout preset:', currentLayout);
				
				if (defaultPreset) {
					const newStyle = { ...headerStyle, ...defaultPreset };
					currentHeaderStyle.set(newStyle);
					manuallySelected = currentLayout;
					toast.success('Custom header style deleted.');
				}
			}
		} else {
			console.log('ℹ️ Not currently selected, just deleting');
		}
		
		// Dispatch event to parent to save immediately
		dispatch('deletePreset', { presets: customPresets });
	}

	function resetToPreset(layout: HeaderStyles['layout']) {
		const presetSettings = presetMap[layout];
		if (!presetSettings) return;
		
		const newStyle = { ...headerStyle, ...presetSettings, avatarBorderColor: '#ffffff' };
		currentHeaderStyle.set(newStyle);
	}
	onDestroy(() => {
		// Cleanup
	});
	
	// Auto-detect which preset matches current header style
	$effect(() => {
		if (headerStyle.layout) {
			// First check if it matches a custom preset (exact match on all properties)
			const matchingCustomId = findMatchingCustomPreset();
			if (matchingCustomId) {
				manuallySelected = matchingCustomId;
			} else {
				// No exact custom match → use layout as identifier for built-in presets
				// This ensures highlight works immediately when theme is selected
				manuallySelected = headerStyle.layout;
			}
		}
	});
</script>

<div class="space-y-6">
	<div>
		<div class="flex items-center justify-between mb-4">
			<h3 class="text-lg font-semibold text-gray-900">Header Style</h3>
			{#if expandedPreset}
				<button
					onclick={() => expandedPreset = null}
					class="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
					</svg>
					Back
				</button>
			{/if}
		</div>
		
		{#if !expandedPreset}
			<!-- Category Tabs -->
			<div class="flex gap-2 mb-3">
				{#each categories as cat}
					<button
						onclick={() => selectedCategory = cat.id}
						class="px-3 py-1.5 rounded-full text-xs font-medium transition-all {selectedCategory === cat.id
							? 'bg-gray-900 text-white'
							: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
					>
						{cat.label}
					</button>
				{/each}
			</div>

			<!-- Header Presets Grid -->
			<div class="grid grid-cols-3 md:grid-cols-5 lg:grid-cols-7 gap-2">
			{#each displayedPresets as preset}
				{@const isCustom = 'settings' in preset}
				{@const isSelected = isCustom 
					? manuallySelected === preset.id 
					: manuallySelected === preset.layout}
				<div class="relative">
					<button
						onclick={() => isCustom ? selectLayout(preset.layout, preset.settings, preset.id) : selectLayout(preset.layout)}
						class="group relative w-full"
					>
						<div class="aspect-[9/16] rounded-lg overflow-hidden border-2 transition-all {isSelected ? 'border-indigo-600 shadow-md ring-2 ring-indigo-200' : 'border-gray-200 hover:border-gray-300'}"
							style="background: {preset.preview.cover};"
						>
							<div class="h-full flex flex-col items-center justify-center p-1.5 relative">
								<div class="w-5 h-5 rounded-full" style="border: 1px solid {preset.preview.avatar}; background: rgba(255,255,255,0.3);"></div>
								<div class="mt-1 w-6 h-0.5 rounded" style="background: {preset.preview.avatar}; opacity: 0.7;"></div>
								<div class="mt-0.5 w-4 h-0.5 rounded" style="background: {preset.preview.avatar}; opacity: 0.5;"></div>
								
								<!-- Customize Button - only show when selected -->
								{#if isSelected}
									<button
										onclick={(e) => { e.stopPropagation(); toggleCustomize(isCustom ? preset.id : preset.layout); }}
										class="absolute bottom-1 left-1 right-1 px-1.5 py-1 text-[9px] font-medium rounded transition-all flex items-center justify-center gap-1 backdrop-blur-sm {expandedPreset === (isCustom ? preset.id : preset.layout) ? 'bg-white/95 text-indigo-600 border border-indigo-600' : 'bg-white/90 text-gray-700 border border-white/50 hover:bg-white'}"
									>
										<svg class="w-2.5 h-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
										</svg>
										{expandedPreset === (isCustom ? preset.id : preset.layout) ? 'Hide' : 'Edit'}
									</button>
								{/if}
							</div>
						</div>
						{#if isSelected}
							<div class="absolute top-1 right-1 w-4 h-4 bg-indigo-600 rounded-full flex items-center justify-center z-10">
								<svg class="w-2.5 h-2.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
								</svg>
							</div>
						{/if}
						{#if isCustom}
							<button
								onclick={(e) => { e.stopPropagation(); confirmDelete(preset.id); }}
								class="absolute top-1 left-1 w-4 h-4 bg-red-600 rounded-full flex items-center justify-center z-10 hover:bg-red-700 transition-colors"
								title="Delete custom preset"
							>
								<svg class="w-2.5 h-2.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							</button>
						{/if}
						<p class="mt-1 text-[10px] font-medium text-gray-600 text-center truncate">{preset.name}</p>
					</button>
				</div>
			{/each}
			</div>
		{:else}
			<!-- Edit Section with Vertical Tabs -->
			{@const preset = allPresets.find(p => ('settings' in p ? p.id : p.layout) === expandedPreset)}
			{#if preset}
				<div class="p-4 bg-gradient-to-br from-indigo-50 to-purple-50 rounded-xl border-2 border-indigo-200">
					<div class="flex gap-4">
						<!-- Vertical Tabs -->
						<div class="flex flex-col gap-2 min-w-[120px]">
							{#if supportsCover}
								<button
									onclick={() => activeTab = 'cover'}
									class="px-4 py-2.5 text-sm font-medium rounded-lg transition-all text-left {activeTab === 'cover' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-700 hover:bg-white/50'}"
								>
									<div class="flex items-center gap-2">
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
										</svg>
										Cover
									</div>
								</button>
							{/if}
							<button
								onclick={() => activeTab = 'avatar'}
								class="px-4 py-2.5 text-sm font-medium rounded-lg transition-all text-left {activeTab === 'avatar' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-700 hover:bg-white/50'}"
							>
								<div class="flex items-center gap-2">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
									</svg>
									Avatar
								</div>
							</button>
							<button
								onclick={() => activeTab = 'bio'}
								class="px-4 py-2.5 text-sm font-medium rounded-lg transition-all text-left {activeTab === 'bio' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-700 hover:bg-white/50'}"
							>
								<div class="flex items-center gap-2">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
									</svg>
									Bio
								</div>
							</button>
						</div>

						<!-- Tab Content -->
						<div class="flex-1 bg-white rounded-lg p-4 space-y-4">
							{#if activeTab === 'cover' && supportsCover}
								<div class="flex items-center justify-between p-2 bg-gray-50 rounded-lg">
									<span class="text-sm font-medium text-gray-900">Show Cover</span>
									<label class="relative inline-flex items-center cursor-pointer">
										<input
											type="checkbox"
											checked={headerStyle.showCover}
											onchange={(e) => updateHeaderSetting({ showCover: e.currentTarget.checked })}
											class="sr-only peer"
										/>
										<div class="w-9 h-5 bg-gray-300 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-indigo-600"></div>
									</label>
								</div>
								{#if headerStyle.showCover}
									<div>
										<div class="flex items-center justify-between mb-2">
											<label class="text-sm font-medium text-gray-900">Cover Height</label>
											<span class="text-sm font-mono text-indigo-600">{headerStyle.coverHeight}px</span>
										</div>
										<input
											type="range"
											min="0"
											max="300"
											value={headerStyle.coverHeight}
											oninput={(e) => updateHeaderSetting({ coverHeight: parseInt(e.currentTarget.value) })}
											class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
										/>
									</div>
								{/if}
							{:else if activeTab === 'avatar'}
								<div>
									<div class="flex items-center justify-between mb-2">
										<label class="text-sm font-medium text-gray-900">Size</label>
										<span class="text-sm font-mono text-indigo-600">{headerStyle.avatarSize}px</span>
									</div>
									<input
										type="range"
										min="60"
										max="180"
										value={headerStyle.avatarSize}
										oninput={(e) => updateHeaderSetting({ avatarSize: parseInt(e.currentTarget.value) })}
										class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
									/>
								</div>
								<div>
									<div class="flex items-center justify-between mb-2">
										<label class="text-sm font-medium text-gray-900">Border Width</label>
										<span class="text-sm font-mono text-indigo-600">{headerStyle.avatarBorder}px</span>
									</div>
									<input
										type="range"
										min="0"
										max="10"
										value={headerStyle.avatarBorder}
										oninput={(e) => updateHeaderSetting({ avatarBorder: parseInt(e.currentTarget.value) })}
										class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
									/>
								</div>
								<div>
									<label class="text-sm font-medium text-gray-900 mb-2 block">Shape</label>
									<div class="grid grid-cols-3 gap-2">
										{#each [
											{ value: 'circle', label: 'Circle', icon: '●' },
											{ value: 'square', label: 'Square', icon: '■' },
											{ value: 'rounded', label: 'Rounded', icon: '▢' },
											{ value: 'vertical', label: 'Vertical', icon: '▯' },
											{ value: 'horizontal', label: 'Horizontal', icon: '▭' }
										] as shape}
											<button
												onclick={() => updateHeaderSetting({ avatarShape: shape.value as any })}
												class="px-3 py-2 text-sm border rounded-lg transition-all flex flex-col items-center gap-1 {headerStyle.avatarShape === shape.value ? 'border-indigo-600 bg-indigo-50 text-indigo-600' : 'border-gray-200 text-gray-700 hover:bg-gray-50'}"
											>
												<span class="text-lg">{shape.icon}</span>
												<span class="text-xs">{shape.label}</span>
											</button>
										{/each}
									</div>
								</div>
								<div>
									<label class="text-sm font-medium text-gray-900 mb-2 block">Alignment</label>
									<div class="flex gap-2">
										{#each ['left', 'center', 'right'] as align}
											<button
												onclick={() => updateHeaderSetting({ avatarAlign: align as any })}
												class="flex-1 px-3 py-2 text-sm border rounded-lg transition-all {headerStyle.avatarAlign === align ? 'border-indigo-600 bg-indigo-50 text-indigo-600 font-medium' : 'border-gray-200 text-gray-700 hover:bg-gray-50'}"
											>
												{align === 'left' ? 'Left' : align === 'center' ? 'Center' : 'Right'}
											</button>
										{/each}
									</div>
								</div>
							{:else if activeTab === 'bio'}
								<div>
									<label class="text-sm font-medium text-gray-900 mb-2 block">Alignment</label>
									<div class="flex gap-2">
										{#each ['left', 'center', 'right'] as align}
											<button
												onclick={() => updateHeaderSetting({ bioAlign: align as any })}
												class="flex-1 px-3 py-2 text-sm border rounded-lg transition-all {headerStyle.bioAlign === align ? 'border-indigo-600 bg-indigo-50 text-indigo-600 font-medium' : 'border-gray-200 text-gray-700 hover:bg-gray-50'}"
											>
												{align === 'left' ? 'Left' : align === 'center' ? 'Center' : 'Right'}
											</button>
										{/each}
									</div>
								</div>
								<div>
									<label class="text-sm font-medium text-gray-900 mb-2 block">Size</label>
									<div class="flex gap-2">
										{#each ['sm', 'md', 'lg'] as size}
											<button
												onclick={() => updateHeaderSetting({ bioSize: size as any })}
												class="flex-1 px-3 py-2 text-sm border rounded-lg transition-all {headerStyle.bioSize === size ? 'border-indigo-600 bg-indigo-50 text-indigo-600 font-medium' : 'border-gray-200 text-gray-700 hover:bg-gray-50'}"
											>
												{size.toUpperCase()}
											</button>
										{/each}
									</div>
								</div>
								<div>
									<label class="text-sm font-medium text-gray-900 mb-2 block">Text Color</label>
									<div class="flex gap-2 items-center">
										<input 
											type="color" 
											value={headerStyle.bioTextColor || '#6b7280'}
											onchange={(e) => updateHeaderSetting({ bioTextColor: e.currentTarget.value })}
											class="w-12 h-12 rounded-lg border-2 border-gray-200 cursor-pointer" 
										/>
										<div class="flex-1 flex gap-2 flex-wrap">
											{#each ['#000000', '#ffffff', '#6b7280', '#3b82f6', '#8b5cf6', '#10b981', '#f59e0b', '#ef4444'] as color}
												<button 
													onclick={() => updateHeaderSetting({ bioTextColor: color })}
													class="w-8 h-8 rounded-lg border-2 hover:scale-110 transition-transform {headerStyle.bioTextColor === color ? 'border-indigo-600 ring-2 ring-indigo-200' : 'border-gray-200'}" 
													style="background-color: {color};"
													title={color}
												></button>
											{/each}
										</div>
									</div>
								</div>
							{/if}

							<!-- Action Buttons -->
							<div class="flex gap-2 pt-2 border-t">
								<button
									onclick={() => resetToPreset(expandedPreset)}
									class="flex-1 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-all"
								>
									Reset
								</button>
								<button
									onclick={saveAsCustom}
									class="flex-1 px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 transition-all"
								>
									Save as Custom
								</button>
							</div>
						</div>
					</div>
				</div>
			{/if}
		{/if}
	</div>
</div>

<!-- Delete Confirmation Modal -->
{#if showDeleteModal}
	<div class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4" onclick={() => showDeleteModal = false}>
		<div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6" onclick={(e) => e.stopPropagation()}>
			<div class="flex items-start gap-4">
				<div class="flex-shrink-0 w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
					<svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				</div>
				<div class="flex-1">
					<h3 class="text-lg font-semibold text-gray-900 mb-2">Delete Custom Header Style?</h3>
					<p class="text-sm text-gray-600 mb-6">
						This action cannot be undone. The custom header style will be permanently deleted immediately.
					</p>
					<div class="flex gap-3">
						<button
							onclick={() => { showDeleteModal = false; presetToDelete = null; }}
							class="flex-1 px-4 py-2.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
						>
							Cancel
						</button>
						<button
							onclick={deleteCustom}
							class="flex-1 px-4 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
						>
							Delete
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
