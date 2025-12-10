<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';
	
	export let group: Link;
	
	const dispatch = createEventDispatcher();
	
	// Use group properties directly - no local state needed
	$: layout = group.group_layout || 'list';
	$: textAlignment = group.text_alignment || 'center';
	$: showOutline = group.show_outline || false;
	$: showShadow = group.show_shadow || false;
	$: showDescription = group.show_description !== undefined ? group.show_description : true;
	
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
	
	function updateShowOutline() {
		dispatch('update', {
			groupId: group.id,
			show_outline: !showOutline
		});
	}
	
	function updateShowShadow() {
		dispatch('update', {
			groupId: group.id,
			show_shadow: !showShadow
		});
	}
	
	function updateShowDescription(value: boolean) {
		dispatch('update', {
			groupId: group.id,
			show_description: value
		});
	}
</script>

<div class="p-6 space-y-6">
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
						class="flex-1 px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md"
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

		<div class="mt-4">
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
						class="px-3 py-2.5 border-2 rounded-lg transition-all hover:shadow-md text-xs font-medium"
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

	<!-- Image Placement Section (only show for card layout) -->
	{#if layout === 'card'}
		<div>
			<h3 class="text-base font-semibold text-gray-900 mb-3">Image placement</h3>
			<div class="flex flex-col gap-2">
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_placement: 'left' })}
					class="px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md text-left"
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
					class="px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md text-left"
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
					class="px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md text-left"
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
				class="flex-1 px-4 py-2.5 border-2 rounded-lg transition-all flex items-center justify-center"
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
				class="flex-1 px-4 py-2.5 border-2 rounded-lg transition-all flex items-center justify-center"
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
				class="flex-1 px-4 py-2.5 border-2 rounded-lg transition-all flex items-center justify-center"
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
					class="flex-1 px-4 py-2.5 border-2 rounded-lg transition-all font-semibold"
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

	<!-- Image Shape Section (for Classic layout only) -->
	{#if layout === 'list'}
		<div>
			<h3 class="text-base font-semibold text-gray-900 mb-3">Image shape</h3>
			<div class="flex gap-2">
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'square' })}
					class="flex-1 px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md"
					class:border-emerald-500={(group.image_shape || 'square') === 'square'}
					class:bg-emerald-50={(group.image_shape || 'square') === 'square'}
					class:border-gray-200={(group.image_shape || 'square') !== 'square'}
					class:bg-white={(group.image_shape || 'square') !== 'square'}
				>
					<div class="flex flex-col items-center gap-2">
						<div class="w-10 h-10 bg-gradient-to-br from-amber-200 to-amber-300 rounded-lg"></div>
						<span class="text-sm font-medium text-gray-700">Square</span>
					</div>
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'circle' })}
					class="flex-1 px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md"
					class:border-emerald-500={(group.image_shape || 'square') === 'circle'}
					class:bg-emerald-50={(group.image_shape || 'square') === 'circle'}
					class:border-gray-200={(group.image_shape || 'square') !== 'circle'}
					class:bg-white={(group.image_shape || 'square') !== 'circle'}
				>
					<div class="flex flex-col items-center gap-2">
						<div class="w-10 h-10 bg-gradient-to-br from-amber-200 to-amber-300 rounded-full"></div>
						<span class="text-sm font-medium text-gray-700">Circle</span>
					</div>
				</button>
			</div>
		</div>
	{/if}

	<!-- Toggles Section -->
	<div class="space-y-3">
		<label class="flex items-center justify-between cursor-pointer">
			<span class="text-base font-medium text-gray-900">Link outline</span>
			<button
				type="button"
				onclick={updateShowOutline}
				class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors"
				class:bg-gray-900={showOutline}
				class:bg-gray-300={!showOutline}
			>
				<span
					class="inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform"
					class:translate-x-6={showOutline}
					class:translate-x-1={!showOutline}
				></span>
			</button>
		</label>
		
		<label class="flex items-center justify-between cursor-pointer">
			<span class="text-base font-medium text-gray-900">Link shadow</span>
			<button
				type="button"
				onclick={updateShowShadow}
				class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors"
				class:bg-gray-900={showShadow}
				class:bg-gray-300={!showShadow}
			>
				<span
					class="inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform"
					class:translate-x-6={showShadow}
					class:translate-x-1={!showShadow}
				></span>
			</button>
		</label>
	</div>

	<!-- Show Description Preset -->
	<div>
		<h3 class="text-base font-semibold text-gray-900 mb-3">Link description</h3>
		<div class="flex gap-2">
			<button
				type="button"
				onclick={() => updateShowDescription(true)}
				class="flex-1 px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md font-medium text-sm"
				class:border-emerald-500={showDescription}
				class:bg-emerald-50={showDescription}
				class:text-emerald-700={showDescription}
				class:border-gray-200={!showDescription}
				class:bg-white={!showDescription}
				class:text-gray-700={!showDescription}
			>
				Show
			</button>
			<button
				type="button"
				onclick={() => updateShowDescription(false)}
				class="flex-1 px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md font-medium text-sm"
				class:border-emerald-500={!showDescription}
				class:bg-emerald-50={!showDescription}
				class:text-emerald-700={!showDescription}
				class:border-gray-200={showDescription}
				class:bg-white={showDescription}
				class:text-gray-700={showDescription}
			>
				Hide
			</button>
		</div>
	</div>
</div>
