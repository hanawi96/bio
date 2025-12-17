<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';
	import { themePresets, defaultCardStyles, globalTheme } from '$lib/stores/theme';
	
	export let group: Link;
	
	const dispatch = createEventDispatcher();
	
	// Reset group to default card styles
	function resetToTheme() {
		const card = defaultCardStyles;
		dispatch('update', {
			groupId: group.id,
			card_background_color: card.cardBackground,
			card_background_opacity: card.cardBackgroundOpacity,
			card_text_color: card.cardTextColor,
			card_border_radius: card.cardBorderRadius,
			has_card_border: card.cardBorder,
			card_border_color: card.cardBorderColor,
			card_border_width: card.cardBorderWidth,
			has_card_background: true
		});
	}
	
	// Debounce timer for slider updates
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;
	
	function dispatchDebounced(data: any) {
		if (debounceTimer) clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			dispatch('update', data);
		}, 300);
	}
	
	// Local state for real-time display
	let localMarginValue: number | null = null;
	let localPaddingValue: number | null = null;
	let lastGroupId: string = '';
	let showCustomPadding = false;
	let showCustomSpacing = false;
	let showCustomOpacity = false;
	let showCustomBorderRadius = false;
	
	// Reset local state when switching to different group
	$: if (group.id !== lastGroupId) {
		localPaddingValue = null;
		localMarginValue = null;
		showCustomPadding = false;
		showCustomSpacing = false;
		showCustomOpacity = false;
		showCustomBorderRadius = false;
		lastGroupId = group.id;
	}
	
	// Use group properties directly - no local state needed
	$: layout = group.group_layout || 'list';
	$: textAlignment = group.text_alignment || 'center';
	$: showText = group.show_text !== undefined ? group.show_text : true;
	
	// Reactive padding and margin values
	$: currentPadding = getPadding();
	$: isUniformPadding = currentPadding.top === currentPadding.right && currentPadding.right === currentPadding.bottom && currentPadding.bottom === currentPadding.left;
	$: currentMargin = getMargin();
	$: displayPadding = localPaddingValue !== null ? localPaddingValue : currentPadding.top;
	$: displayMargin = localMarginValue !== null ? localMarginValue : currentMargin.bottom;
	
	// When local padding is set, assume uniform (since we set all 4 sides equal)
	$: isDisplayUniform = localPaddingValue !== null ? true : isUniformPadding;
	
	function updateLayout(value: string) {
		dispatch('update', {
			groupId: group.id,
			group_layout: value
		});
	}
	
	function updateTextAlignment(value: string) {
		dispatch('update', {
			groupId: group.id,
			text_alignment: value
		});
	}
	
	function updateShowText(value: boolean) {
		dispatch('update', {
			groupId: group.id,
			show_text: value
		});
	}

	function getPadding(): { top: number; right: number; bottom: number; left: number } {
		if (!group.style) return { top: 16, right: 16, bottom: 16, left: 16 };
		try {
			const style = JSON.parse(group.style);
			if (typeof style.padding === 'number') {
				return { top: style.padding, right: style.padding, bottom: style.padding, left: style.padding };
			}
			return {
				top: style.padding?.top || 16,
				right: style.padding?.right || 16,
				bottom: style.padding?.bottom || 16,
				left: style.padding?.left || 16
			};
		} catch {
			return { top: 16, right: 16, bottom: 16, left: 16 };
		}
	}

	function updatePadding(value: number) {
		localPaddingValue = value;
		let style: any = {};
		if (group.style) {
			try {
				style = JSON.parse(group.style);
			} catch {}
		}
		style.padding = { top: value, right: value, bottom: value, left: value };
		dispatch('update', {
			groupId: group.id,
			style: JSON.stringify(style)
		});
	}

	function updatePaddingSide(side: 'top' | 'right' | 'bottom' | 'left', value: number) {
		localPaddingValue = null; // Reset uniform padding state when adjusting individual sides
		let style: any = {};
		if (group.style) {
			try {
				style = JSON.parse(group.style);
			} catch {}
		}
		if (typeof style.padding === 'number') {
			style.padding = { top: style.padding, right: style.padding, bottom: style.padding, left: style.padding };
		}
		if (!style.padding) style.padding = { top: 16, right: 16, bottom: 16, left: 16 };
		style.padding[side] = value;
		dispatchDebounced({
			groupId: group.id,
			style: JSON.stringify(style)
		});
	}

	function getMargin(): { top: number; bottom: number } {
		if (!group.style) return { top: 0, bottom: 0 };
		try {
			const style = JSON.parse(group.style);
			return {
				top: style.margin?.top ?? 0,
				bottom: style.margin?.bottom ?? 0
			};
		} catch {
			return { top: 0, bottom: 0 };
		}
	}

	function updateMargin(value: number) {
		localMarginValue = value;
		let style: any = {};
		if (group.style) {
			try {
				style = JSON.parse(group.style);
			} catch {}
		}
		style.margin = { top: 0, bottom: value };
		dispatch('update', {
			groupId: group.id,
			style: JSON.stringify(style)
		});
	}

	function updateOpacityPreset(value: number) {
		dispatch('update', {
			groupId: group.id,
			has_card_background: true,
			card_background_opacity: value
		});
	}

	function updateBorderRadiusPreset(value: number) {
		dispatch('update', {
			groupId: group.id,
			has_card_background: true,
			card_border_radius: value
		});
	}
	
	// Slider - debounced update with local state
	function updateMarginDebounced(value: number) {
		localMarginValue = value; // Update local state immediately for display
		let style: any = {};
		if (group.style) {
			try {
				style = JSON.parse(group.style);
			} catch {}
		}
		style.margin = { top: 0, bottom: value };
		dispatchDebounced({
			groupId: group.id,
			style: JSON.stringify(style)
		});
	}

</script>

<div class="p-6 space-y-6">
	<!-- Reset to Theme Button -->
	<div class="flex justify-end">
		<button
			onclick={resetToTheme}
			class="px-3 py-1.5 text-xs font-medium text-gray-600 hover:text-gray-900 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors flex items-center gap-1.5"
			title="Reset this group's styles to match the current theme"
		>
			<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
			</svg>
			Reset to Theme
		</button>
	</div>

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
		<h3 class="text-base font-semibold text-gray-900 mb-3">Text alignment</h3>
		<div class="flex gap-2">
			<button
				type="button"
				onclick={() => updateTextAlignment('left')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center"
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
			</button>
			
			<button
				type="button"
				onclick={() => updateTextAlignment('center')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center"
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
			</button>
			
			<button
				type="button"
				onclick={() => updateTextAlignment('right')}
				class="flex-1 px-4 py-2.5 border rounded-lg transition-all flex items-center justify-center"
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
			</button>
		</div>
	</div>

	<!-- Text Size Section -->
	<div>
		<h3 class="text-base font-semibold text-gray-900 mb-3">Text size</h3>
		<div class="flex gap-2">
			{#each ['S', 'M', 'L', 'XL'] as size}
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, text_size: size })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all font-semibold"
					class:border-gray-900={(group.text_size || 'M') === size}
					class:bg-gray-900={(group.text_size || 'M') === size}
					class:text-white={(group.text_size || 'M') === size}
					class:border-gray-200={(group.text_size || 'M') !== size}
					class:bg-white={(group.text_size || 'M') !== size}
					class:text-gray-700={(group.text_size || 'M') !== size}
					class:text-xs={size === 'S'}
					class:text-sm={size === 'M'}
					class:text-base={size === 'L'}
					class:text-lg={size === 'XL'}
				>
					{size}
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
			<h3 class="text-base font-semibold text-gray-900 mb-3">Image shape</h3>
			<div class="flex gap-2">
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'sharp' })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md"
					class:border-emerald-500={(group.image_shape || 'square') === 'sharp'}
					class:bg-emerald-50={(group.image_shape || 'square') === 'sharp'}
					class:border-gray-200={(group.image_shape || 'square') !== 'sharp'}
					class:bg-white={(group.image_shape || 'square') !== 'sharp'}
				>
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 flex-shrink-0"></div>
						<span class="text-sm font-medium text-gray-700">Sharp</span>
					</div>
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'square' })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md"
					class:border-emerald-500={(group.image_shape || 'square') === 'square'}
					class:bg-emerald-50={(group.image_shape || 'square') === 'square'}
					class:border-gray-200={(group.image_shape || 'square') !== 'square'}
					class:bg-white={(group.image_shape || 'square') !== 'square'}
				>
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 rounded-lg flex-shrink-0"></div>
						<span class="text-sm font-medium text-gray-700">Rounded</span>
					</div>
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'circle' })}
					class="flex-1 px-4 py-2.5 border rounded-lg transition-all hover:shadow-md"
					class:border-emerald-500={(group.image_shape || 'square') === 'circle'}
					class:bg-emerald-50={(group.image_shape || 'square') === 'circle'}
					class:border-gray-200={(group.image_shape || 'square') !== 'circle'}
					class:bg-white={(group.image_shape || 'square') !== 'circle'}
				>
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 rounded-full flex-shrink-0"></div>
						<span class="text-sm font-medium text-gray-700">Circle</span>
					</div>
				</button>
			</div>
		</div>
	{/if}

	<!-- Link Card Background Section -->
	<div class="space-y-4">
		<div>
			<div class="flex items-center gap-3 mb-3">
				<h3 class="text-base font-semibold text-gray-900">Link card background</h3>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, has_card_background: !(group.has_card_background || false) })}
					class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors"
					class:bg-gray-900={group.has_card_background}
					class:bg-gray-300={!group.has_card_background}
				>
					<span
						class="inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform"
						class:translate-x-5={group.has_card_background}
						class:translate-x-0.5={!group.has_card_background}
					></span>
				</button>
			</div>

		{#if group.has_card_background}
			<!-- Background Color -->
			<div>
				<div class="flex gap-3 items-center">
					<!-- Color Picker (Large) -->
					<div class="flex items-center gap-2">
						<input 
							type="color" 
							value={group.card_background_color || '#ffffff'}
							onchange={(e) => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: e.currentTarget.value })}
							class="w-12 h-12 rounded-lg border-2 border-gray-200 cursor-pointer overflow-hidden" 
						/>
						<span class="text-xs text-gray-600 font-mono">{group.card_background_color || '#ffffff'}</span>
					</div>
					
					<!-- Color Presets (Small) -->
					<div class="flex gap-1.5 flex-wrap flex-1">
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#ffffff' })} class="w-6 h-6 rounded-full border-2 bg-white hover:scale-110 transition-transform" title="White"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#f3f4f6' })} class="w-6 h-6 rounded-full border-2 bg-gray-100 hover:scale-110 transition-transform" title="Light Gray"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#000000' })} class="w-6 h-6 rounded-full border-2 bg-black hover:scale-110 transition-transform" title="Black"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#3b82f6' })} class="w-6 h-6 rounded-full border-2 bg-blue-500 hover:scale-110 transition-transform" title="Blue"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#8b5cf6' })} class="w-6 h-6 rounded-full border-2 bg-purple-500 hover:scale-110 transition-transform" title="Purple"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#10b981' })} class="w-6 h-6 rounded-full border-2 bg-green-500 hover:scale-110 transition-transform" title="Green"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#ef4444' })} class="w-6 h-6 rounded-full border-2 bg-red-500 hover:scale-110 transition-transform" title="Red"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#f59e0b' })} class="w-6 h-6 rounded-full border-2 bg-amber-500 hover:scale-110 transition-transform" title="Amber"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#ec4899' })} class="w-6 h-6 rounded-full border-2 bg-pink-500 hover:scale-110 transition-transform" title="Pink"></button>
					</div>
				</div>
			</div>

			<!-- Background Opacity -->
			<div>
				<h3 class="text-base font-semibold text-gray-900 mb-3 mt-2">Background opacity</h3>
				
				<!-- Preset Buttons -->
				<div class="flex gap-2 mb-3">
					<button 
						onclick={() => {
							showCustomOpacity = false;
							updateOpacityPreset(50);
						}} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={!showCustomOpacity && (group.card_background_opacity || 100) === 50}
						class:bg-emerald-50={!showCustomOpacity && (group.card_background_opacity || 100) === 50}
						class:border-gray-200={showCustomOpacity || (group.card_background_opacity || 100) !== 50}
						class:bg-white={showCustomOpacity || (group.card_background_opacity || 100) !== 50}
					>
						50%
					</button>
					<button 
						onclick={() => {
							showCustomOpacity = false;
							updateOpacityPreset(75);
						}} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={!showCustomOpacity && (group.card_background_opacity || 100) === 75}
						class:bg-emerald-50={!showCustomOpacity && (group.card_background_opacity || 100) === 75}
						class:border-gray-200={showCustomOpacity || (group.card_background_opacity || 100) !== 75}
						class:bg-white={showCustomOpacity || (group.card_background_opacity || 100) !== 75}
					>
						75%
					</button>
					<button 
						onclick={() => {
							showCustomOpacity = false;
							updateOpacityPreset(100);
						}} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={!showCustomOpacity && (group.card_background_opacity || 100) === 100}
						class:bg-emerald-50={!showCustomOpacity && (group.card_background_opacity || 100) === 100}
						class:border-gray-200={showCustomOpacity || (group.card_background_opacity || 100) !== 100}
						class:bg-white={showCustomOpacity || (group.card_background_opacity || 100) !== 100}
					>
						100%
					</button>
					<button 
						onclick={() => showCustomOpacity = !showCustomOpacity} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={showCustomOpacity}
						class:bg-emerald-50={showCustomOpacity}
						class:border-gray-200={!showCustomOpacity}
						class:bg-white={!showCustomOpacity}
					>
						Custom
					</button>
				</div>

				{#if showCustomOpacity}
					<div class="pl-4 border-l-2 border-gray-200">
						<div class="flex items-center justify-between mb-2">
							<span class="text-xs font-medium text-gray-700">Custom opacity</span>
							<span class="text-xs font-mono text-gray-500">{group.card_background_opacity || 100}%</span>
						</div>
						<input 
							type="range" 
							min="0" 
							max="100" 
							value={group.card_background_opacity || 100}
							oninput={(e) => dispatchDebounced({ groupId: group.id, has_card_background: true, card_background_opacity: parseInt(e.currentTarget.value) })}
							class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
						/>
					</div>
				{/if}
			</div>

			<!-- Border Radius -->
			<div>
				<h3 class="text-base font-semibold text-gray-900 mb-3">Border radius</h3>
				
				<!-- Preset Buttons -->
				<div class="flex gap-2 mb-3">
					<button 
						onclick={() => {
							showCustomBorderRadius = false;
							updateBorderRadiusPreset(0);
						}} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={!showCustomBorderRadius && (group.card_border_radius || 12) === 0}
						class:bg-emerald-50={!showCustomBorderRadius && (group.card_border_radius || 12) === 0}
						class:border-gray-200={showCustomBorderRadius || (group.card_border_radius || 12) !== 0}
						class:bg-white={showCustomBorderRadius || (group.card_border_radius || 12) !== 0}
					>
						None
					</button>
					<button 
						onclick={() => {
							showCustomBorderRadius = false;
							updateBorderRadiusPreset(8);
						}} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={!showCustomBorderRadius && (group.card_border_radius || 12) === 8}
						class:bg-emerald-50={!showCustomBorderRadius && (group.card_border_radius || 12) === 8}
						class:border-gray-200={showCustomBorderRadius || (group.card_border_radius || 12) !== 8}
						class:bg-white={showCustomBorderRadius || (group.card_border_radius || 12) !== 8}
					>
						Small
					</button>
					<button 
						onclick={() => {
							showCustomBorderRadius = false;
							updateBorderRadiusPreset(12);
						}} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={!showCustomBorderRadius && (group.card_border_radius || 12) === 12}
						class:bg-emerald-50={!showCustomBorderRadius && (group.card_border_radius || 12) === 12}
						class:border-gray-200={showCustomBorderRadius || (group.card_border_radius || 12) !== 12}
						class:bg-white={showCustomBorderRadius || (group.card_border_radius || 12) !== 12}
					>
						Medium
					</button>
					<button 
						onclick={() => {
							showCustomBorderRadius = false;
							updateBorderRadiusPreset(24);
						}} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={!showCustomBorderRadius && (group.card_border_radius || 12) === 24}
						class:bg-emerald-50={!showCustomBorderRadius && (group.card_border_radius || 12) === 24}
						class:border-gray-200={showCustomBorderRadius || (group.card_border_radius || 12) !== 24}
						class:bg-white={showCustomBorderRadius || (group.card_border_radius || 12) !== 24}
					>
						Large
					</button>
					<button 
						onclick={() => showCustomBorderRadius = !showCustomBorderRadius} 
						class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={showCustomBorderRadius}
						class:bg-emerald-50={showCustomBorderRadius}
						class:border-gray-200={!showCustomBorderRadius}
						class:bg-white={!showCustomBorderRadius}
					>
						Custom
					</button>
				</div>

				{#if showCustomBorderRadius}
					<div class="pl-4 border-l-2 border-gray-200">
						<div class="flex items-center justify-between mb-2">
							<span class="text-xs font-medium text-gray-700">Custom radius</span>
							<span class="text-xs font-mono text-gray-500">{group.card_border_radius || 12}px</span>
						</div>
						<input 
							type="range" 
							min="0" 
							max="32" 
							value={group.card_border_radius || 12}
							oninput={(e) => dispatchDebounced({ groupId: group.id, has_card_background: true, card_border_radius: parseInt(e.currentTarget.value) })}
							class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
						/>
					</div>
				{/if}
			</div>
		{/if}
		</div>
	</div>

	<!-- Link Card Border Section -->
	<div class="space-y-4">
		<div>
			<div class="flex items-center gap-3 mb-3">
				<h3 class="text-base font-semibold text-gray-900">Link card border</h3>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, has_card_border: !(group.has_card_border || false) })}
					class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors"
					class:bg-gray-900={group.has_card_border}
					class:bg-gray-300={!group.has_card_border}
				>
					<span
						class="inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform"
						class:translate-x-5={group.has_card_border}
						class:translate-x-0.5={!group.has_card_border}
					></span>
				</button>
			</div>

		{#if group.has_card_border}
			<!-- Border Color -->
			<div>
				<label class="text-sm font-medium text-gray-700 mb-2 block">Border Color</label>
				<div class="flex gap-3 items-center">
					<!-- Color Picker (Large) -->
					<div class="flex items-center gap-2">
						<input 
							type="color" 
							value={group.card_border_color || '#e5e7eb'}
							onchange={(e) => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: e.currentTarget.value })}
							class="w-12 h-12 rounded-lg border-2 border-gray-200 cursor-pointer overflow-hidden" 
						/>
						<span class="text-xs text-gray-600 font-mono">{group.card_border_color || '#e5e7eb'}</span>
					</div>
					
					<!-- Color Presets (Small) -->
					<div class="flex gap-1.5 flex-wrap flex-1">
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#e5e7eb' })} class="w-6 h-6 rounded-full border-2 bg-gray-200 hover:scale-110 transition-transform" title="Light Gray"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#9ca3af' })} class="w-6 h-6 rounded-full border-2 bg-gray-400 hover:scale-110 transition-transform" title="Gray"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#000000' })} class="w-6 h-6 rounded-full border-2 bg-black hover:scale-110 transition-transform" title="Black"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#3b82f6' })} class="w-6 h-6 rounded-full border-2 bg-blue-500 hover:scale-110 transition-transform" title="Blue"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#8b5cf6' })} class="w-6 h-6 rounded-full border-2 bg-purple-500 hover:scale-110 transition-transform" title="Purple"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#10b981' })} class="w-6 h-6 rounded-full border-2 bg-green-500 hover:scale-110 transition-transform" title="Green"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#ef4444' })} class="w-6 h-6 rounded-full border-2 bg-red-500 hover:scale-110 transition-transform" title="Red"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_color: '#f59e0b' })} class="w-6 h-6 rounded-full border-2 bg-amber-500 hover:scale-110 transition-transform" title="Amber"></button>
					</div>
				</div>
			</div>

			<!-- Border Style -->
			<div>
				<label class="text-sm font-medium text-gray-700 mt-2 mb-2 block">Border Style</label>
				<div class="flex gap-2">
					<button 
						onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_style: 'solid' })} 
						class="flex-1 px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={(group.card_border_style || 'solid') === 'solid'}
						class:bg-emerald-50={(group.card_border_style || 'solid') === 'solid'}
						class:border-gray-200={(group.card_border_style || 'solid') !== 'solid'}
					>
						<div class="w-full h-0.5 bg-gray-700"></div>
						<span class="block mt-1">Solid</span>
					</button>
					<button 
						onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_style: 'dashed' })} 
						class="flex-1 px-3 py-2 text-xs bg-white border-2 border-dashed rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={(group.card_border_style || 'solid') === 'dashed'}
						class:bg-emerald-50={(group.card_border_style || 'solid') === 'dashed'}
						class:border-gray-200={(group.card_border_style || 'solid') !== 'dashed'}
					>
						<div class="w-full h-0.5 border-t-2 border-dashed border-gray-700"></div>
						<span class="block mt-1">Dashed</span>
					</button>
					<button 
						onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_style: 'dotted' })} 
						class="flex-1 px-3 py-2 text-xs bg-white border-2 border-dotted rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={(group.card_border_style || 'solid') === 'dotted'}
						class:bg-emerald-50={(group.card_border_style || 'solid') === 'dotted'}
						class:border-gray-200={(group.card_border_style || 'solid') !== 'dotted'}
					>
						<div class="w-full h-0.5 border-t-2 border-dotted border-gray-700"></div>
						<span class="block mt-1">Dotted</span>
					</button>
				</div>
			</div>

			<!-- Border Width -->
			<div>
				<label class="text-sm font-medium text-gray-700 mb-2 mt-2 block">Border Width</label>
				<div class="flex gap-2">
					<button 
						onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_width: 1 })} 
						class="flex-1 px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={(group.card_border_width || 1) === 1}
						class:bg-emerald-50={(group.card_border_width || 1) === 1}
						class:border-gray-200={(group.card_border_width || 1) !== 1}
					>1px</button>
					<button 
						onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_width: 2 })} 
						class="flex-1 px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={(group.card_border_width || 1) === 2}
						class:bg-emerald-50={(group.card_border_width || 1) === 2}
						class:border-gray-200={(group.card_border_width || 1) !== 2}
					>2px</button>
					<button 
						onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_width: 3 })} 
						class="flex-1 px-3 py-2 text-xs bg-white border-4 rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={(group.card_border_width || 1) === 3}
						class:bg-emerald-50={(group.card_border_width || 1) === 3}
						class:border-gray-200={(group.card_border_width || 1) !== 3}
					>3px</button>
					<button 
						onclick={() => dispatch('update', { groupId: group.id, has_card_border: true, card_border_width: 4 })} 
						class="flex-1 px-3 py-2 text-xs bg-white border-4 rounded-lg hover:bg-gray-50 transition-colors"
						class:border-emerald-500={(group.card_border_width || 1) === 4}
						class:bg-emerald-50={(group.card_border_width || 1) === 4}
						class:border-gray-200={(group.card_border_width || 1) !== 4}
					>4px</button>
				</div>
			</div>
		{/if}
		</div>
	</div>

	<!-- Link Card Padding Section -->
	<div class="space-y-4">
		<h3 class="text-base font-semibold text-gray-900">Link Card Padding</h3>
		
		<!-- Preset Buttons -->
		<div class="flex gap-2">
				<button 
					onclick={() => {
						showCustomPadding = false;
						updatePadding(8);
					}} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={displayPadding === 8 && isDisplayUniform && !showCustomPadding}
					class:bg-emerald-50={displayPadding === 8 && isDisplayUniform && !showCustomPadding}
					class:border-gray-200={!(displayPadding === 8 && isDisplayUniform && !showCustomPadding)}
					class:bg-white={!(displayPadding === 8 && isDisplayUniform && !showCustomPadding)}
				>
					Small
				</button>
				<button 
					onclick={() => {
						showCustomPadding = false;
						updatePadding(16);
					}} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={displayPadding === 16 && isDisplayUniform && !showCustomPadding}
					class:bg-emerald-50={displayPadding === 16 && isDisplayUniform && !showCustomPadding}
					class:border-gray-200={!(displayPadding === 16 && isDisplayUniform && !showCustomPadding)}
					class:bg-white={!(displayPadding === 16 && isDisplayUniform && !showCustomPadding)}
				>
					Medium
				</button>
				<button 
					onclick={() => {
						showCustomPadding = false;
						updatePadding(24);
					}} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={displayPadding === 24 && isDisplayUniform && !showCustomPadding}
					class:bg-emerald-50={displayPadding === 24 && isDisplayUniform && !showCustomPadding}
					class:border-gray-200={!(displayPadding === 24 && isDisplayUniform && !showCustomPadding)}
					class:bg-white={!(displayPadding === 24 && isDisplayUniform && !showCustomPadding)}
				>
					Large
				</button>
				<button 
					onclick={() => {
						showCustomPadding = false;
						updatePadding(32);
					}} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={displayPadding === 32 && isDisplayUniform && !showCustomPadding}
					class:bg-emerald-50={displayPadding === 32 && isDisplayUniform && !showCustomPadding}
					class:border-gray-200={!(displayPadding === 32 && isDisplayUniform && !showCustomPadding)}
					class:bg-white={!(displayPadding === 32 && isDisplayUniform && !showCustomPadding)}
				>
					XL
				</button>
				<button 
					onclick={() => showCustomPadding = !showCustomPadding} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={showCustomPadding}
					class:bg-emerald-50={showCustomPadding}
					class:border-gray-200={!showCustomPadding}
					class:bg-white={!showCustomPadding}
				>
					Custom
				</button>
		</div>

		{#if showCustomPadding}
			<!-- Custom Padding Sliders - 2x2 Grid -->
			<div class="grid grid-cols-2 gap-4 pl-4 border-l-2 border-gray-200">
				<!-- Top -->
				<div>
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-sm font-medium text-gray-700">Top</span>
						<span class="text-xs font-mono text-gray-500">{currentPadding.top}px</span>
					</div>
					<input 
						type="range" 
						min="0" 
						max="64" 
						value={currentPadding.top}
						oninput={(e) => updatePaddingSide('top', parseInt(e.currentTarget.value))}
						class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
					/>
				</div>
				
				<!-- Right -->
				<div>
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-sm font-medium text-gray-700">Right</span>
						<span class="text-xs font-mono text-gray-500">{currentPadding.right}px</span>
					</div>
					<input 
						type="range" 
						min="0" 
						max="64" 
						value={currentPadding.right}
						oninput={(e) => updatePaddingSide('right', parseInt(e.currentTarget.value))}
						class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
					/>
				</div>
				
				<!-- Bottom -->
				<div>
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-sm font-medium text-gray-700">Bottom</span>
						<span class="text-xs font-mono text-gray-500">{currentPadding.bottom}px</span>
					</div>
					<input 
						type="range" 
						min="0" 
						max="64" 
						value={currentPadding.bottom}
						oninput={(e) => updatePaddingSide('bottom', parseInt(e.currentTarget.value))}
						class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
					/>
				</div>
				
				<!-- Left -->
				<div>
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-sm font-medium text-gray-700">Left</span>
						<span class="text-xs font-mono text-gray-500">{currentPadding.left}px</span>
					</div>
					<input 
						type="range" 
						min="0" 
						max="64" 
						value={currentPadding.left}
						oninput={(e) => updatePaddingSide('left', parseInt(e.currentTarget.value))}
						class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
					/>
				</div>
			</div>
		{/if}
	</div>

	<!-- Link Card Spacing Section -->
	<div class="space-y-4">
		<h3 class="text-base font-semibold text-gray-900">Link card spacing</h3>
		<p class="text-xs text-gray-500">Space between cards (uses margin-bottom)</p>
		
		<!-- Preset Buttons -->
		<div class="flex gap-2">
				<button 
					onclick={() => {
						showCustomSpacing = false;
						updateMargin(4);
					}} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={displayMargin === 4 && !showCustomSpacing}
					class:bg-emerald-50={displayMargin === 4 && !showCustomSpacing}
					class:border-gray-200={!(displayMargin === 4 && !showCustomSpacing)}
					class:bg-white={!(displayMargin === 4 && !showCustomSpacing)}
				>
					Small
				</button>
				<button 
					onclick={() => {
						showCustomSpacing = false;
						updateMargin(8);
					}} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={displayMargin === 8 && !showCustomSpacing}
					class:bg-emerald-50={displayMargin === 8 && !showCustomSpacing}
					class:border-gray-200={!(displayMargin === 8 && !showCustomSpacing)}
					class:bg-white={!(displayMargin === 8 && !showCustomSpacing)}
				>
					Medium
				</button>
				<button 
					onclick={() => {
						showCustomSpacing = false;
						updateMargin(16);
					}} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={displayMargin === 16 && !showCustomSpacing}
					class:bg-emerald-50={displayMargin === 16 && !showCustomSpacing}
					class:border-gray-200={!(displayMargin === 16 && !showCustomSpacing)}
					class:bg-white={!(displayMargin === 16 && !showCustomSpacing)}
				>
					Large
				</button>
				<button 
					onclick={() => showCustomSpacing = !showCustomSpacing} 
					class="flex-1 px-3 py-2 text-xs border rounded-lg hover:bg-gray-50 transition-colors"
					class:border-emerald-500={showCustomSpacing}
					class:bg-emerald-50={showCustomSpacing}
					class:border-gray-200={!showCustomSpacing}
					class:bg-white={!showCustomSpacing}
				>
					Custom
				</button>
			</div>

		{#if showCustomSpacing}
			<!-- Custom Spacing Slider -->
			<div class="pl-4 border-l-2 border-gray-200">
				<div class="flex items-center justify-between mb-2">
					<span class="text-sm font-medium text-gray-700">Custom spacing</span>
					<span class="text-xs font-mono text-gray-500">{localMarginValue !== null ? localMarginValue : currentMargin.bottom}px</span>
				</div>
				<input 
					type="range" 
					min="0" 
					max="32" 
					value={localMarginValue !== null ? localMarginValue : currentMargin.bottom}
					oninput={(e) => updateMarginDebounced(parseInt(e.currentTarget.value))}
					class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
				/>
			</div>
		{/if}
	</div>
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
