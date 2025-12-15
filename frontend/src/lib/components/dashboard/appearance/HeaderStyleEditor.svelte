<script lang="ts">
	import { currentHeaderStyle, themePresets, type HeaderStyles } from '$lib/stores/theme';
	import { pendingChanges } from '$lib/stores/pendingChanges';
	import { onDestroy } from 'svelte';
	
	type CustomPreset = { id: string; name: string; layout: HeaderStyles['layout']; preview: { cover: string; avatar: string }; settings: Partial<HeaderStyles> };
	
	const headerStyle = $derived($currentHeaderStyle);
	let selectedCategory = $state<string>('all');
	let applying = $state(false);
	let manuallySelected = $state<string | null>(null);
	let expandedPreset = $state<string | null>(null);
	let customPresets = $state<CustomPreset[]>([]);
	let saveTimer: ReturnType<typeof setTimeout> | null = null;
	
	// Load custom presets from localStorage
	if (typeof window !== 'undefined') {
		const saved = localStorage.getItem('customHeaderPresets');
		if (saved) customPresets = JSON.parse(saved);
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
		if (applying) return;
		
		const settings = customSettings || presetMap[layout];
		const newStyle = { ...headerStyle, ...settings, avatarBorderColor: headerStyle.avatarBorderColor || '#ffffff' };
		
		manuallySelected = customId || layout;
		expandedPreset = null;
		currentHeaderStyle.set(newStyle);
		pendingChanges.updateHeader(newStyle);
	}

	function toggleCustomize(id: string) {
		if (expandedPreset === id) {
			expandedPreset = null;
		} else {
			expandedPreset = id;
		}
	}

	function updateHeaderSetting(updates: Partial<HeaderStyles>) {
		const newStyle = { ...headerStyle, ...updates };
		console.log('[HeaderStyleEditor] updateHeaderSetting:', {
			updates,
			newStyle: {
				bioAlign: newStyle.bioAlign,
				avatarAlign: newStyle.avatarAlign,
				layout: newStyle.layout
			}
		});
		currentHeaderStyle.set(newStyle);
		pendingChanges.updateHeader(newStyle);
	}

	function saveAsCustom() {
		const id = `custom-${Date.now()}`;
		const basePreset = allPresets.find(p => p.layout === headerStyle.layout);
		
		const customPreset: CustomPreset = {
			id,
			name: `Custom ${customPresets.length + 1}`,
			layout: headerStyle.layout,
			preview: basePreset?.preview || { cover: headerStyle.coverGradientFrom || '#667eea', avatar: '#ffffff' },
			settings: { ...headerStyle }
		};
		
		customPresets = [...customPresets, customPreset];
		manuallySelected = id as any;
		expandedPreset = null;
		
		if (typeof window !== 'undefined') {
			localStorage.setItem('customHeaderPresets', JSON.stringify(customPresets));
		}
	}

	function deleteCustom(id: string) {
		customPresets = customPresets.filter(p => p.id !== id);
		if (typeof window !== 'undefined') {
			localStorage.setItem('customHeaderPresets', JSON.stringify(customPresets));
		}
	}

	function resetToPreset(layout: HeaderStyles['layout']) {
		const presetSettings = presetMap[layout];
		if (!presetSettings) return;
		
		const newStyle = { ...headerStyle, ...presetSettings, avatarBorderColor: '#ffffff' };
		currentHeaderStyle.set(newStyle);
		pendingChanges.updateHeader(newStyle);
	}
	onDestroy(() => {
		if (saveTimer) clearTimeout(saveTimer);
	});
	
	$effect(() => {
		if (headerStyle.layout && !manuallySelected) {
			manuallySelected = headerStyle.layout;
		}
	});
</script>

<div class="space-y-6">
	<div>
		<div class="flex items-center justify-between mb-4">
			<h3 class="text-lg font-semibold text-gray-900">Header Style</h3>
			{#if applying}
				<span class="text-sm text-indigo-600 flex items-center gap-2">
					<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Saving...
				</span>
			{/if}
		</div>
		
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
				<div class="relative">
					<button
						onclick={() => isCustom ? selectLayout(preset.layout, preset.settings, preset.id) : selectLayout(preset.layout)}
						class="group relative w-full"
					>
						<div class="aspect-[9/16] rounded-lg overflow-hidden border-2 transition-all {manuallySelected === preset.layout ? 'border-indigo-600 shadow-md ring-2 ring-indigo-200' : 'border-gray-200 hover:border-gray-300'}"
							style="background: {preset.preview.cover};"
						>
							<div class="h-full flex flex-col items-center justify-center p-1.5 relative">
								<div class="w-5 h-5 rounded-full" style="border: 1px solid {preset.preview.avatar}; background: rgba(255,255,255,0.3);"></div>
								<div class="mt-1 w-6 h-0.5 rounded" style="background: {preset.preview.avatar}; opacity: 0.7;"></div>
								<div class="mt-0.5 w-4 h-0.5 rounded" style="background: {preset.preview.avatar}; opacity: 0.5;"></div>
								
								<!-- Customize Button - only show when selected -->
								{#if manuallySelected === (isCustom ? preset.id : preset.layout)}
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
						{#if manuallySelected === (isCustom ? preset.id : preset.layout)}
							<div class="absolute top-1 right-1 w-4 h-4 bg-indigo-600 rounded-full flex items-center justify-center z-10">
								<svg class="w-2.5 h-2.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
								</svg>
							</div>
						{/if}
						{#if isCustom}
							<button
								onclick={(e) => { e.stopPropagation(); deleteCustom(preset.id); }}
								class="absolute top-1 left-1 w-4 h-4 bg-red-600 rounded-full flex items-center justify-center z-10 hover:bg-red-700"
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
		
		<!-- Close button for expanded section -->
		{#if expandedPreset}
			<div class="mt-3 flex justify-end">
				<button
					onclick={() => expandedPreset = null}
					class="px-3 py-1.5 text-xs font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all"
				>
					Close Customization
				</button>
			</div>
		{/if}
		
		<!-- Expand Section (FULL WIDTH BELOW GRID) -->
		{#if expandedPreset}
			{@const preset = allPresets.find(p => ('settings' in p ? p.id : p.layout) === expandedPreset)}
			{#if preset}
				<div class="mt-4 p-6 bg-gradient-to-br from-indigo-50 to-purple-50 rounded-xl border-2 border-indigo-200 space-y-4">
							{#if supportsCover}
								<!-- Show Cover Toggle -->
								<div class="flex items-center justify-between p-2 bg-white rounded-lg">
									<span class="text-xs font-medium text-gray-900">Show Cover</span>
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
									<!-- Cover Height -->
									<div>
										<div class="flex items-center justify-between mb-2">
											<label class="text-xs font-semibold text-gray-900">Cover Height</label>
											<span class="text-xs font-mono text-indigo-600">{headerStyle.coverHeight}px</span>
										</div>
										<input
											type="range"
											min="0"
											max="300"
											value={headerStyle.coverHeight}
											oninput={(e) => updateHeaderSetting({ coverHeight: parseInt(e.currentTarget.value) })}
											class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
										/>
									</div>
								{/if}
							{/if}

							<!-- Avatar Size -->
							<div>
								<div class="flex items-center justify-between mb-2">
									<label class="text-xs font-semibold text-gray-900">Avatar Size</label>
									<span class="text-xs font-mono text-indigo-600">{headerStyle.avatarSize}px</span>
								</div>
								<input
									type="range"
									min="60"
									max="180"
									value={headerStyle.avatarSize}
									oninput={(e) => updateHeaderSetting({ avatarSize: parseInt(e.currentTarget.value) })}
									class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
								/>
							</div>

							<!-- Avatar Border -->
							<div>
								<div class="flex items-center justify-between mb-2">
									<label class="text-xs font-semibold text-gray-900">Avatar Border</label>
									<span class="text-xs font-mono text-indigo-600">{headerStyle.avatarBorder}px</span>
								</div>
								<input
									type="range"
									min="0"
									max="10"
									value={headerStyle.avatarBorder}
									oninput={(e) => updateHeaderSetting({ avatarBorder: parseInt(e.currentTarget.value) })}
									class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
								/>
							</div>

							<!-- Avatar Shape -->
							<div>
								<label class="text-xs font-semibold text-gray-900 mb-2 block">Avatar Shape</label>
								<div class="grid grid-cols-3 gap-1.5">
									{#each [
										{ value: 'circle', label: 'Circle', icon: '●' },
										{ value: 'square', label: 'Square', icon: '■' },
										{ value: 'rounded', label: 'Rounded', icon: '▢' },
										{ value: 'vertical', label: 'Vertical', icon: '▯' },
										{ value: 'horizontal', label: 'Horizontal', icon: '▭' }
									] as shape}
										<button
											onclick={() => updateHeaderSetting({ avatarShape: shape.value as any })}
											class="px-2 py-1.5 text-xs border rounded transition-all flex flex-col items-center gap-0.5 {headerStyle.avatarShape === shape.value ? 'border-indigo-600 bg-indigo-50 text-indigo-600' : 'border-gray-200 text-gray-700 hover:bg-gray-100'}"
										>
											<span class="text-base">{shape.icon}</span>
											<span class="text-[10px]">{shape.label}</span>
										</button>
									{/each}
								</div>
							</div>

							<!-- Avatar Alignment -->
							<div>
								<label class="text-xs font-semibold text-gray-900 mb-2 block">Avatar Alignment</label>
								<div class="flex gap-1.5">
									{#each ['left', 'center', 'right'] as align}
										<button
											onclick={() => updateHeaderSetting({ avatarAlign: align as any })}
											class="flex-1 px-2 py-1.5 text-xs border rounded transition-all {headerStyle.avatarAlign === align ? 'border-indigo-600 bg-indigo-50 text-indigo-600' : 'border-gray-200 text-gray-700 hover:bg-gray-100'}"
										>
											{align === 'left' ? 'Left' : align === 'center' ? 'Center' : 'Right'}
										</button>
									{/each}
								</div>
							</div>

							<!-- Bio Alignment -->
							<div>
								<label class="text-xs font-semibold text-gray-900 mb-2 block">Bio Alignment</label>
								<div class="flex gap-1.5">
									{#each ['left', 'center', 'right'] as align}
										<button
											onclick={() => updateHeaderSetting({ bioAlign: align as any })}
											class="flex-1 px-2 py-1.5 text-xs border rounded transition-all {headerStyle.bioAlign === align ? 'border-indigo-600 bg-indigo-50 text-indigo-600' : 'border-gray-200 text-gray-700 hover:bg-gray-100'}"
										>
											{align === 'left' ? 'Left' : align === 'center' ? 'Center' : 'Right'}
										</button>
									{/each}
								</div>
							</div>

							<!-- Bio Size -->
							<div>
								<label class="text-xs font-semibold text-gray-900 mb-2 block">Bio Size</label>
								<div class="flex gap-1.5">
									{#each ['sm', 'md', 'lg'] as size}
										<button
											onclick={() => updateHeaderSetting({ bioSize: size as any })}
											class="flex-1 px-2 py-1.5 text-xs border rounded transition-all {headerStyle.bioSize === size ? 'border-indigo-600 bg-indigo-50 text-indigo-600' : 'border-gray-200 text-gray-700 hover:bg-gray-100'}"
										>
											{size.toUpperCase()}
										</button>
									{/each}
								</div>
							</div>

							<!-- Bio Text Color -->
							<div>
								<label class="text-xs font-semibold text-gray-900 mb-2 block">Bio Text Color</label>
								<div class="flex gap-2 items-center">
									<input 
										type="color" 
										value={headerStyle.bioTextColor || '#6b7280'}
										onchange={(e) => updateHeaderSetting({ bioTextColor: e.currentTarget.value })}
										class="w-10 h-10 rounded-lg border-2 border-gray-200 cursor-pointer" 
									/>
									<div class="flex-1 flex gap-1 flex-wrap">
										{#each ['#000000', '#ffffff', '#6b7280', '#3b82f6', '#8b5cf6', '#10b981', '#f59e0b', '#ef4444'] as color}
											<button 
												onclick={() => updateHeaderSetting({ bioTextColor: color })}
												class="w-6 h-6 rounded-md border-2 hover:scale-110 transition-transform {headerStyle.bioTextColor === color ? 'border-indigo-600 ring-2 ring-indigo-200' : 'border-gray-200'}" 
												style="background-color: {color};"
												title={color}
											></button>
										{/each}
									</div>
								</div>
							</div>

					<div class="flex gap-2">
						<button
							onclick={() => resetToPreset(expandedPreset)}
							class="flex-1 px-4 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 transition-all"
						>
							Reset
						</button>
						<button
							onclick={saveAsCustom}
							class="flex-1 px-4 py-2.5 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 transition-all"
						>
							Save as Custom
						</button>
					</div>
				</div>
			{/if}
		{/if}
	</div>
</div>
