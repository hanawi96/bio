<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';
	import { themePresets, defaultCardStyles, globalTheme } from '$lib/stores/theme';
	
	export let group: Link;
	
	const dispatch = createEventDispatcher();
	
	// Reset group to theme defaults - only reset properties that are currently custom
	function resetToTheme() {
		const updates: any = {
			groupId: group.id,
			has_custom_layout: false
		};
		
		// Only reset properties that have custom values (not null)
		if (group.text_alignment !== null && group.text_alignment !== undefined) {
			updates.text_alignment = null;
		}
		if (group.text_size !== null && group.text_size !== undefined) {
			updates.text_size = null;
		}
		if (group.image_shape !== null && group.image_shape !== undefined) {
			updates.image_shape = null;
		}
		
		dispatch('update', updates);
	}
	
	// Check if any property is custom
	$: hasCustomProperties = 
		(group.text_alignment !== null && group.text_alignment !== undefined) ||
		(group.text_size !== null && group.text_size !== undefined) ||
		(group.image_shape !== null && group.image_shape !== undefined);
	
	// Debounce timer for slider updates
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;
	
	function dispatchDebounced(data: any) {
		if (debounceTimer) clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			dispatch('update', data);
		}, 300);
	}
	
	// Get current theme for fallback values
	$: currentTheme = globalTheme.getCurrent();
	
	// Use group properties with theme fallback
	// Always use theme as fallback if field is null, regardless of has_custom_layout
	$: layout = group.group_layout || 'list';
	$: textAlignment = group.text_alignment || currentTheme?.textAlignment || 'center';
	$: imageShape = group.image_shape || currentTheme?.imageShape || 'square';
	$: textSize = group.text_size || currentTheme?.textSize || 'M';
	$: showText = group.show_text !== undefined ? group.show_text : true;
	
	console.log('[LAYOUT_SETTINGS]', group.group_title, {
		has_custom: group.has_custom_layout,
		text_alignment: group.text_alignment,
		image_shape: group.image_shape,
		text_size: group.text_size,
		theme_align: currentTheme?.textAlignment,
		final_textAlignment: textAlignment,
		final_imageShape: imageShape
	});
	

	function updateLayout(value: string) {
		dispatch('update', {
			groupId: group.id,
			group_layout: value
		});
	}
	
	function updateTextAlignment(value: string) {
		dispatch('update', {
			groupId: group.id,
			text_alignment: value,
			has_custom_layout: true
		});
	}
	
	function updateShowText(value: boolean) {
		dispatch('update', {
			groupId: group.id,
			show_text: value
		});
	}



</script>

<div class="p-6 space-y-6">
	<!-- Reset to Theme Button - Only show if there are custom properties -->
	{#if hasCustomProperties}
		<div class="flex justify-end">
			<button
				onclick={resetToTheme}
				class="px-3 py-1.5 text-xs font-medium text-indigo-600 hover:text-indigo-700 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors flex items-center gap-1.5"
				title="Reset all custom styles to match the current theme"
			>
				<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
				</svg>
				Reset All to Theme
			</button>
		</div>
	{/if}

	<!-- Layout Section -->
	<div>
		<h3 class="text-base font-semibold text-gray-900 mb-3">Layout</h3>
		<div class="flex gap-3">
			<!-- Classic -->
			<button
				type="button"
				onclick={() => updateLayout('list')}
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'list'}
				class:ring-emerald-500={layout === 'list'}
				class:ring-offset-2={layout === 'list'}
			>
				<div class="aspect-[2/1] p-3 flex flex-col justify-center gap-1.5">
					<div class="flex items-center gap-1.5 bg-white rounded-md p-1.5 shadow-sm">
						<div class="w-5 h-5 bg-gradient-to-br from-amber-200 to-amber-300 rounded flex-shrink-0"></div>
						<div class="flex-1 space-y-0.5">
							<div class="h-1 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
					<div class="flex items-center gap-1.5 bg-white rounded-md p-1.5 shadow-sm">
						<div class="w-5 h-5 bg-gradient-to-br from-amber-200 to-amber-300 rounded flex-shrink-0"></div>
						<div class="flex-1 space-y-0.5">
							<div class="h-1 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
					<div class="flex items-center gap-1.5 bg-white rounded-md p-1.5 shadow-sm">
						<div class="w-5 h-5 bg-gradient-to-br from-amber-200 to-amber-300 rounded flex-shrink-0"></div>
						<div class="flex-1 space-y-0.5">
							<div class="h-1 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
				</div>
				<div class="p-1.5 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-700">Classic</p>
				</div>
			</button>

			<!-- Carousel -->
			<button
				type="button"
				onclick={() => updateLayout('carousel')}
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'carousel'}
				class:ring-emerald-500={layout === 'carousel'}
				class:ring-offset-2={layout === 'carousel'}
			>
				<div class="aspect-[2/1] p-3 flex items-center justify-center gap-1.5">
					<div class="w-12 h-14 bg-gradient-to-br from-yellow-100 to-yellow-200 rounded-md shadow-sm opacity-60"></div>
					<div class="w-14 h-16 bg-gradient-to-br from-emerald-100 to-emerald-200 rounded-md shadow-md"></div>
				</div>
				<div class="p-1.5 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-700">Carousel</p>
				</div>
			</button>

			<!-- Image Grid -->
			<button
				type="button"
				onclick={() => updateLayout('grid')}
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'grid'}
				class:ring-emerald-500={layout === 'grid'}
				class:ring-offset-2={layout === 'grid'}
			>
				<div class="aspect-[2/1] p-3">
					<div class="grid grid-cols-3 gap-1 h-full">
						<div class="bg-gradient-to-br from-orange-200 to-red-200 rounded-sm"></div>
						<div class="bg-gradient-to-br from-gray-200 to-gray-300 rounded-sm"></div>
						<div class="bg-gradient-to-br from-yellow-200 to-amber-200 rounded-sm"></div>
						<div class="bg-gradient-to-br from-pink-200 to-rose-200 rounded-sm"></div>
						<div class="bg-gradient-to-br from-yellow-300 to-amber-300 rounded-sm"></div>
						<div class="bg-gradient-to-br from-amber-100 to-orange-200 rounded-sm"></div>
						<div class="bg-gradient-to-br from-orange-300 to-amber-400 rounded-sm"></div>
						<div class="bg-gradient-to-br from-gray-300 to-slate-300 rounded-sm"></div>
						<div class="bg-gradient-to-br from-blue-200 to-cyan-200 rounded-sm"></div>
					</div>
				</div>
				<div class="p-1.5 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-700">Image grid</p>
				</div>
			</button>

			<!-- Card -->
			<button
				type="button"
				onclick={() => updateLayout('card')}
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'card'}
				class:ring-emerald-500={layout === 'card'}
				class:ring-offset-2={layout === 'card'}
			>
				<div class="aspect-[2/1] p-3 space-y-1.5 flex flex-col justify-center">
					<div class="bg-white rounded-md overflow-hidden shadow-sm flex">
						<div class="w-10 h-10 bg-gradient-to-br from-amber-200 to-orange-200 flex-shrink-0"></div>
						<div class="flex-1 p-1.5 space-y-0.5">
							<div class="h-1 bg-gray-200 rounded-full w-4/5"></div>
							<div class="h-1 bg-gray-150 rounded-full w-3/5"></div>
						</div>
					</div>
					<div class="bg-white rounded-md overflow-hidden shadow-sm flex">
						<div class="flex-1 p-1.5 space-y-0.5">
							<div class="h-1 bg-gray-200 rounded-full w-4/5"></div>
							<div class="h-1 bg-gray-150 rounded-full w-3/5"></div>
						</div>
						<div class="w-10 h-10 bg-gradient-to-br from-blue-200 to-cyan-200 flex-shrink-0"></div>
					</div>
				</div>
				<div class="p-1.5 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-700">Card</p>
				</div>
			</button>
		</div>
	</div>

	<!-- Grid Columns Preset (only show for grid layout) -->
	{#if layout === 'grid'}
		<div>
			<h3 class="text-base font-semibold text-gray-900 mb-3">Grid columns</h3>
			<div class="flex gap-2">
				{#each [1, 2, 3, 4] as cols}
					<button
						type="button"
						onclick={() => dispatch('update', { groupId: group.id, grid_columns: cols })}
						class="flex-1 px-4 py-3 border rounded-lg transition-all hover:shadow-md"
						class:border-emerald-500={(group.grid_columns || 2) === cols}
						class:bg-emerald-50={(group.grid_columns || 2) === cols}
						class:border-gray-200={(group.grid_columns || 2) !== cols}
						class:bg-white={(group.grid_columns || 2) !== cols}
					>
						<div class="text-center">
							<div class="text-lg font-bold text-gray-900">{cols}</div>
							<div class="text-xs text-gray-500 mt-1">
								{cols === 1 ? 'col' : 'cols'}
							</div>
						</div>
					</button>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Aspect Ratio Section (for Grid and Carousel layouts) -->
	{#if layout === 'grid' || layout === 'carousel'}
		<div class:mt-4={layout === 'grid'}>
			<h3 class="text-base font-semibold text-gray-900 mb-3">Aspect ratio</h3>
			<div class="grid grid-cols-3 gap-2">
				{#each [
					{ value: '1:1', label: '1:1 Square' },
					{ value: '3:2', label: '3:2 Horizontal' },
					{ value: '16:9', label: '16:9 Horizontal' },
					{ value: '3:1', label: '3:1 Horizontal' },
					{ value: '2:3', label: '2:3 Vertical' }
				] as ratio}
					<button
						type="button"
						onclick={() => dispatch('update', { groupId: group.id, grid_aspect_ratio: ratio.value })}
						class="px-3 py-2.5 border rounded-lg transition-all hover:shadow-md text-xs font-medium"
						class:border-emerald-500={(group.grid_aspect_ratio || '3:2') === ratio.value}
						class:bg-emerald-50={(group.grid_aspect_ratio || '3:2') === ratio.value}
						class:text-emerald-700={(group.grid_aspect_ratio || '3:2') === ratio.value}
						class:border-gray-200={(group.grid_aspect_ratio || '3:2') !== ratio.value}
						class:bg-white={(group.grid_aspect_ratio || '3:2') !== ratio.value}
						class:text-gray-700={(group.grid_aspect_ratio || '3:2') !== ratio.value}
					>
						{ratio.label}
					</button>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Text Visibility Toggle (only for carousel layout) -->
	{#if layout === 'carousel'}
		<div>
			<h3 class="text-base font-semibold text-gray-900 mb-3">Text visibility</h3>
			<div class="flex gap-2">
				<button
					type="button"
					onclick={() => updateShowText(true)}
					class="flex-1 px-4 py-3 border rounded-lg transition-all hover:shadow-md font-medium text-sm"
					class:border-emerald-500={showText}
					class:bg-emerald-50={showText}
					class:text-emerald-700={showText}
					class:border-gray-200={!showText}
					class:bg-white={!showText}
					class:text-gray-700={!showText}
				>
					Show
				</button>
				<button
					type="button"
					onclick={() => updateShowText(false)}
					class="flex-1 px-4 py-3 border rounded-lg transition-all hover:shadow-md font-medium text-sm"
					class:border-emerald-500={!showText}
					class:bg-emerald-50={!showText}
					class:text-emerald-700={!showText}
					class:border-gray-200={showText}
					class:bg-white={showText}
					class:text-gray-700={showText}
				>
					Hide
				</button>
			</div>
			<p class="text-xs text-gray-500 mt-2">Toggle title and description visibility in carousel cards</p>
		</div>
	{/if}

	<!-- Image Placement Section (only show for card layout) -->
	{#if layout === 'card'}
		<div>
			<h3 class="text-base font-semibold text-gray-900 mb-3">Image placement</h3>
			<div class="flex flex-col gap-2">
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_placement: 'left' })}
					class="px-4 py-3 border rounded-lg transition-all hover:shadow-md text-left"
					class:border-emerald-500={(group.image_placement || 'alternating') === 'left'}
					class:bg-emerald-50={(group.image_placement || 'alternating') === 'left'}
					class:border-gray-200={(group.image_placement || 'alternating') !== 'left'}
					class:bg-white={(group.image_placement || 'alternating') !== 'left'}
				>
					<div class="font-semibold text-gray-900">Left</div>
					<div class="text-xs text-gray-500 mt-0.5">Image on the left, text on the right</div>
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_placement: 'right' })}
					class="px-4 py-3 border rounded-lg transition-all hover:shadow-md text-left"
					class:border-emerald-500={(group.image_placement || 'alternating') === 'right'}
					class:bg-emerald-50={(group.image_placement || 'alternating') === 'right'}
					class:border-gray-200={(group.image_placement || 'alternating') !== 'right'}
					class:bg-white={(group.image_placement || 'alternating') !== 'right'}
				>
					<div class="font-semibold text-gray-900">Right</div>
					<div class="text-xs text-gray-500 mt-0.5">Image on the right, text on the left</div>
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_placement: 'alternating' })}
					class="px-4 py-3 border rounded-lg transition-all hover:shadow-md text-left"
					class:border-emerald-500={(group.image_placement || 'alternating') === 'alternating'}
					class:bg-emerald-50={(group.image_placement || 'alternating') === 'alternating'}
					class:border-gray-200={(group.image_placement || 'alternating') !== 'alternating'}
					class:bg-white={(group.image_placement || 'alternating') !== 'alternating'}
				>
					<div class="font-semibold text-gray-900">Alternating</div>
					<div class="text-xs text-gray-500 mt-0.5">Image alternates between left and right</div>
				</button>
			</div>
		</div>
	{/if}

	<!-- Text Alignment Section -->
	<div>
		<div class="flex items-center justify-between mb-3">
			<h3 class="text-base font-semibold text-gray-900">Text alignment</h3>
			{#if group.text_alignment !== null && group.text_alignment !== undefined}
				<button
					onclick={() => dispatch('update', { groupId: group.id, text_alignment: null })}
					class="text-xs text-indigo-600 hover:text-indigo-700 font-medium flex items-center gap-1"
					title="Reset to theme default"
				>
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
					</svg>
					Reset
				</button>
			{/if}
		</div>
		<div class="flex gap-2">
			<button
				type="button"
				onclick={() => updateTextAlignment('left')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center relative"
				class:border-gray-900={textAlignment === 'left'}
				class:bg-gray-900={textAlignment === 'left'}
				class:text-white={textAlignment === 'left'}
				class:border-gray-200={textAlignment !== 'left'}
				class:bg-white={textAlignment !== 'left'}
				class:text-gray-700={textAlignment !== 'left'}
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h10M4 18h14"/>
				</svg>
				{#if group.text_alignment === 'left'}
					<span class="absolute -top-1 -right-1 w-2 h-2 bg-indigo-500 rounded-full" title="Custom value"></span>
				{/if}
			</button>
			
			<button
				type="button"
				onclick={() => updateTextAlignment('center')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center relative"
				class:border-gray-900={textAlignment === 'center'}
				class:bg-gray-900={textAlignment === 'center'}
				class:text-white={textAlignment === 'center'}
				class:border-gray-200={textAlignment !== 'center'}
				class:bg-white={textAlignment !== 'center'}
				class:text-gray-700={textAlignment !== 'center'}
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M7 12h10M6 18h12"/>
				</svg>
				{#if group.text_alignment === 'center'}
					<span class="absolute -top-1 -right-1 w-2 h-2 bg-indigo-500 rounded-full" title="Custom value"></span>
				{/if}
			</button>
			
			<button
				type="button"
				onclick={() => updateTextAlignment('right')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center relative"
				class:border-gray-900={textAlignment === 'right'}
				class:bg-gray-900={textAlignment === 'right'}
				class:text-white={textAlignment === 'right'}
				class:border-gray-200={textAlignment !== 'right'}
				class:bg-white={textAlignment !== 'right'}
				class:text-gray-700={textAlignment !== 'right'}
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M10 12h10M8 18h12"/>
				</svg>
				{#if group.text_alignment === 'right'}
					<span class="absolute -top-1 -right-1 w-2 h-2 bg-indigo-500 rounded-full" title="Custom value"></span>
				{/if}
			</button>
		</div>
	</div>

	<!-- Text Size Section -->
	<div>
		<div class="flex items-center justify-between mb-3">
			<h3 class="text-base font-semibold text-gray-900">Text size</h3>
			{#if group.text_size !== null && group.text_size !== undefined}
				<button
					onclick={() => dispatch('update', { groupId: group.id, text_size: null })}
					class="text-xs text-indigo-600 hover:text-indigo-700 font-medium flex items-center gap-1"
					title="Reset to theme default"
				>
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
					</svg>
					Reset
				</button>
			{/if}
		</div>
		<div class="flex gap-2">
			{#each ['S', 'M', 'L', 'XL'] as size}
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, text_size: size, has_custom_layout: true })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all font-semibold relative"
					class:border-gray-900={textSize === size}
					class:bg-gray-900={textSize === size}
					class:text-white={textSize === size}
					class:border-gray-200={textSize !== size}
					class:bg-white={textSize !== size}
					class:text-gray-700={textSize !== size}
					class:text-xs={size === 'S'}
					class:text-sm={size === 'M'}
					class:text-base={size === 'L'}
					class:text-lg={size === 'XL'}
				>
					{size}
					{#if group.text_size === size}
						<span class="absolute -top-1 -right-1 w-2 h-2 bg-indigo-500 rounded-full" title="Custom value"></span>
					{/if}
				</button>
			{/each}
		</div>
	</div>

	<!-- Text Color Section -->
	<div>
		<h3 class="text-base font-semibold text-gray-900 mb-3">Text color</h3>
		<div class="flex gap-3 items-center">
			<!-- Color Picker (Large) -->
			<div class="flex items-center gap-2">
				<input 
					type="color" 
					value={group.card_text_color || '#000000'}
					onchange={(e) => dispatch('update', { groupId: group.id, card_text_color: e.currentTarget.value })}
					class="w-12 h-12 rounded-lg border-2 border-gray-200 cursor-pointer overflow-hidden" 
				/>
				<span class="text-xs text-gray-600 font-mono">{group.card_text_color || '#000000'}</span>
			</div>
			
			<!-- Color Presets (Small) -->
			<div class="flex gap-1.5 flex-wrap flex-1">
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#000000' })} class="w-6 h-6 rounded-full border-2 bg-black hover:scale-110 transition-transform" title="Black"></button>
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#ffffff' })} class="w-6 h-6 rounded-full border-2 bg-white hover:scale-110 transition-transform" title="White"></button>
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#6b7280' })} class="w-6 h-6 rounded-full border-2 bg-gray-500 hover:scale-110 transition-transform" title="Gray"></button>
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#3b82f6' })} class="w-6 h-6 rounded-full border-2 bg-blue-500 hover:scale-110 transition-transform" title="Blue"></button>
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#10b981' })} class="w-6 h-6 rounded-full border-2 bg-green-500 hover:scale-110 transition-transform" title="Green"></button>
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#ef4444' })} class="w-6 h-6 rounded-full border-2 bg-red-500 hover:scale-110 transition-transform" title="Red"></button>
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#f59e0b' })} class="w-6 h-6 rounded-full border-2 bg-amber-500 hover:scale-110 transition-transform" title="Amber"></button>
				<button onclick={() => dispatch('update', { groupId: group.id, card_text_color: '#8b5cf6' })} class="w-6 h-6 rounded-full border-2 bg-purple-500 hover:scale-110 transition-transform" title="Purple"></button>
			</div>
		</div>
	</div>

	<!-- Image Shape Section (for Classic layout only) -->
	{#if layout === 'list'}
		<div>
			<div class="flex items-center justify-between mb-3">
				<h3 class="text-base font-semibold text-gray-900">Image shape</h3>
				{#if group.image_shape !== null && group.image_shape !== undefined}
					<button
						onclick={() => dispatch('update', { groupId: group.id, image_shape: null })}
						class="text-xs text-indigo-600 hover:text-indigo-700 font-medium flex items-center gap-1"
						title="Reset to theme default"
					>
						<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
						</svg>
						Reset
					</button>
				{/if}
			</div>
			<div class="flex gap-2">
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'sharp', has_custom_layout: true })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md relative"
					class:border-emerald-500={imageShape === 'sharp'}
					class:bg-emerald-50={imageShape === 'sharp'}
					class:border-gray-200={imageShape !== 'sharp'}
					class:bg-white={imageShape !== 'sharp'}
				>
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 flex-shrink-0"></div>
						<span class="text-sm font-medium text-gray-700">Sharp</span>
					</div>
					{#if group.image_shape === 'sharp'}
						<span class="absolute -top-1 -right-1 w-2 h-2 bg-indigo-500 rounded-full" title="Custom value"></span>
					{/if}
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'square', has_custom_layout: true })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md relative"
					class:border-emerald-500={imageShape === 'square'}
					class:bg-emerald-50={imageShape === 'square'}
					class:border-gray-200={imageShape !== 'square'}
					class:bg-white={imageShape !== 'square'}
				>
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 rounded-lg flex-shrink-0"></div>
						<span class="text-sm font-medium text-gray-700">Rounded</span>
					</div>
					{#if group.image_shape === 'square'}
						<span class="absolute -top-1 -right-1 w-2 h-2 bg-indigo-500 rounded-full" title="Custom value"></span>
					{/if}
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'circle', has_custom_layout: true })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md relative"
					class:border-emerald-500={imageShape === 'circle'}
					class:bg-emerald-50={imageShape === 'circle'}
					class:border-gray-200={imageShape !== 'circle'}
					class:bg-white={imageShape !== 'circle'}
				>
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 rounded-full flex-shrink-0"></div>
						<span class="text-sm font-medium text-gray-700">Circle</span>
					</div>
					{#if group.image_shape === 'circle'}
						<span class="absolute -top-1 -right-1 w-2 h-2 bg-indigo-500 rounded-full" title="Custom value"></span>
					{/if}
				</button>
			</div>
		</div>
	{/if}
</div>

<style>
	/* Remove default browser styling from color input to make color fill the entire container */
	input[type="color"] {
		padding: 0;
		background: none;
	}
	
	input[type="color"]::-webkit-color-swatch-wrapper {
		padding: 0;
		border: none;
	}
	
	input[type="color"]::-webkit-color-swatch {
		border: none;
		border-radius: 8px;
	}
	
	input[type="color"]::-moz-color-swatch {
		border: none;
		border-radius: 8px;
	}
</style>
