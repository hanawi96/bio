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
	$: showText = group.show_text !== undefined ? group.show_text : true;
	
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
	
	function updateShowText(value: boolean) {
		dispatch('update', {
			groupId: group.id,
			show_text: value
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

	<!-- Text Visibility Toggle (only for carousel layout) -->
	{#if layout === 'carousel'}
		<div>
			<h3 class="text-base font-semibold text-gray-900 mb-3">Text visibility</h3>
			<div class="flex gap-2">
				<button
					type="button"
					onclick={() => updateShowText(true)}
					class="flex-1 px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md font-medium text-sm"
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
					class="flex-1 px-4 py-3 border-2 rounded-lg transition-all hover:shadow-md font-medium text-sm"
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
					class="flex-1 px-4 py-2.5 border-2 rounded-lg transition-all hover:shadow-md"
					class:border-emerald-500={(group.image_shape || 'square') === 'square'}
					class:bg-emerald-50={(group.image_shape || 'square') === 'square'}
					class:border-gray-200={(group.image_shape || 'square') !== 'square'}
					class:bg-white={(group.image_shape || 'square') !== 'square'}
				>
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-gradient-to-br from-amber-200 to-amber-300 rounded-lg flex-shrink-0"></div>
						<span class="text-sm font-medium text-gray-700">Square</span>
					</div>
				</button>
				<button
					type="button"
					onclick={() => dispatch('update', { groupId: group.id, image_shape: 'circle' })}
					class="flex-1 px-4 py-2.5 border-2 rounded-lg transition-all hover:shadow-md"
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

	<!-- Toggles Section -->
	<div class="space-y-4">
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
		
		<div>
			<label class="flex items-center justify-between cursor-pointer mb-3">
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
			
			{#if showShadow}
				<div class="space-y-3 pl-4 border-l-2 border-gray-200">
					<!-- Shadow X -->
					<div>
						<div class="flex items-center justify-between mb-1.5">
							<div class="flex items-center gap-2">
								<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18"/>
								</svg>
								<span class="text-sm font-medium text-gray-700">Horizontal</span>
							</div>
							<span class="text-xs font-mono text-gray-500">{group.shadow_x || 0}px</span>
						</div>
						<input 
							type="range" 
							min="-20" 
							max="20" 
							value={group.shadow_x || 0}
							oninput={(e) => dispatch('update', { groupId: group.id, show_shadow: true, shadow_x: parseInt(e.currentTarget.value) })}
							class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-gray-900"
						/>
					</div>
					
					<!-- Shadow Y -->
					<div>
						<div class="flex items-center justify-between mb-1.5">
							<div class="flex items-center gap-2">
								<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 17l4-4m0 0l-4-4m4 4H4"/>
								</svg>
								<span class="text-sm font-medium text-gray-700">Vertical</span>
							</div>
							<span class="text-xs font-mono text-gray-500">{group.shadow_y || 4}px</span>
						</div>
						<input 
							type="range" 
							min="0" 
							max="20" 
							value={group.shadow_y || 4}
							oninput={(e) => dispatch('update', { groupId: group.id, show_shadow: true, shadow_y: parseInt(e.currentTarget.value) })}
							class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-gray-900"
						/>
					</div>
					
					<!-- Shadow Blur -->
					<div>
						<div class="flex items-center justify-between mb-1.5">
							<div class="flex items-center gap-2">
								<svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
								</svg>
								<span class="text-sm font-medium text-gray-700">Blur</span>
							</div>
							<span class="text-xs font-mono text-gray-500">{group.shadow_blur || 10}px</span>
						</div>
						<input 
							type="range" 
							min="0" 
							max="40" 
							value={group.shadow_blur || 10}
							oninput={(e) => dispatch('update', { groupId: group.id, show_shadow: true, shadow_blur: parseInt(e.currentTarget.value) })}
							class="w-full h-1.5 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-gray-900"
						/>
					</div>
				</div>
			{/if}
		</div>
	</div>

	<!-- Divider -->
	<div class="border-t border-gray-200 my-6"></div>

	<!-- Link Card Background Section -->
	<div class="space-y-4">
		<h3 class="text-base font-semibold text-gray-900">Link Card Background</h3>
		
		<!-- Enable Background -->
		<div class="flex items-center justify-between">
			<label class="text-sm font-medium text-gray-700">Enable Custom Background</label>
			<button
				type="button"
				onclick={() => dispatch('update', { groupId: group.id, has_card_background: !(group.has_card_background || false) })}
				class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors"
				class:bg-emerald-600={group.has_card_background}
				class:bg-gray-300={!group.has_card_background}
			>
				<span
					class="inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform"
					class:translate-x-6={group.has_card_background}
					class:translate-x-1={!group.has_card_background}
				></span>
			</button>
		</div>

		{#if group.has_card_background}
			<!-- Background Color -->
			<div>
				<label class="text-sm font-medium text-gray-700 mb-2 block">Background Color</label>
				<div class="space-y-3">
					<div class="flex gap-3 items-center">
						<input 
							type="color" 
							value={group.card_background_color || '#ffffff'}
							onchange={(e) => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: e.currentTarget.value })}
							class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer overflow-hidden" 
						/>
						<span class="text-sm text-gray-600 font-mono">{group.card_background_color || '#ffffff'}</span>
					</div>
					
					<!-- Color Presets -->
					<div class="flex gap-2 flex-wrap">
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#ffffff' })} class="w-7 h-7 rounded-full border-2 bg-white hover:scale-110 transition-transform" title="White"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#f3f4f6' })} class="w-7 h-7 rounded-full border-2 bg-gray-100 hover:scale-110 transition-transform" title="Light Gray"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#000000' })} class="w-7 h-7 rounded-full border-2 bg-black hover:scale-110 transition-transform" title="Black"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#3b82f6' })} class="w-7 h-7 rounded-full border-2 bg-blue-500 hover:scale-110 transition-transform" title="Blue"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#8b5cf6' })} class="w-7 h-7 rounded-full border-2 bg-purple-500 hover:scale-110 transition-transform" title="Purple"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#10b981' })} class="w-7 h-7 rounded-full border-2 bg-green-500 hover:scale-110 transition-transform" title="Green"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#ef4444' })} class="w-7 h-7 rounded-full border-2 bg-red-500 hover:scale-110 transition-transform" title="Red"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#f59e0b' })} class="w-7 h-7 rounded-full border-2 bg-amber-500 hover:scale-110 transition-transform" title="Amber"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_background_color: '#ec4899' })} class="w-7 h-7 rounded-full border-2 bg-pink-500 hover:scale-110 transition-transform" title="Pink"></button>
					</div>
				</div>
			</div>

			<!-- Opacity -->
			<div>
				<label class="text-sm font-medium text-gray-700 mb-2 block">Opacity: {group.card_background_opacity || 100}%</label>
				<input 
					type="range" 
					min="0" 
					max="100" 
					value={group.card_background_opacity || 100}
					oninput={(e) => dispatch('update', { groupId: group.id, has_card_background: true, card_background_opacity: parseInt(e.currentTarget.value) })}
					class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
				/>
			</div>

			<!-- Border Radius -->
			<div>
				<label class="text-sm font-medium text-gray-700 mb-2 block">Border Radius: {group.card_border_radius || 12}px</label>
				<input 
					type="range" 
					min="0" 
					max="32" 
					value={group.card_border_radius || 12}
					oninput={(e) => dispatch('update', { groupId: group.id, has_card_background: true, card_border_radius: parseInt(e.currentTarget.value) })}
					class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-emerald-600"
				/>
				<div class="flex gap-2 mt-2">
					<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_border_radius: 0 })} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Sharp</button>
					<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_border_radius: 8 })} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Small</button>
					<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_border_radius: 12 })} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Medium</button>
					<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_border_radius: 16 })} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Large</button>
					<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_border_radius: 24 })} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">XL</button>
				</div>
			</div>

			<!-- Text Color -->
			<div>
				<label class="text-sm font-medium text-gray-700 mb-2 block">Text Color</label>
				<div class="space-y-3">
					<div class="flex gap-3 items-center">
						<input 
							type="color" 
							value={group.card_text_color || '#000000'}
							onchange={(e) => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: e.currentTarget.value })}
							class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer overflow-hidden" 
						/>
						<span class="text-sm text-gray-600 font-mono">{group.card_text_color || '#000000'}</span>
					</div>
					
					<!-- Text Color Presets -->
					<div class="flex gap-2 flex-wrap">
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#000000' })} class="w-7 h-7 rounded-full border-2 bg-black hover:scale-110 transition-transform" title="Black"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#ffffff' })} class="w-7 h-7 rounded-full border-2 bg-white hover:scale-110 transition-transform" title="White"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#6b7280' })} class="w-7 h-7 rounded-full border-2 bg-gray-500 hover:scale-110 transition-transform" title="Gray"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#3b82f6' })} class="w-7 h-7 rounded-full border-2 bg-blue-500 hover:scale-110 transition-transform" title="Blue"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#10b981' })} class="w-7 h-7 rounded-full border-2 bg-green-500 hover:scale-110 transition-transform" title="Green"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#ef4444' })} class="w-7 h-7 rounded-full border-2 bg-red-500 hover:scale-110 transition-transform" title="Red"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#f59e0b' })} class="w-7 h-7 rounded-full border-2 bg-amber-500 hover:scale-110 transition-transform" title="Amber"></button>
						<button onclick={() => dispatch('update', { groupId: group.id, has_card_background: true, card_text_color: '#8b5cf6' })} class="w-7 h-7 rounded-full border-2 bg-purple-500 hover:scale-110 transition-transform" title="Purple"></button>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Remove default browser styling from color input to make color fill the entire circle */
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
		border-radius: 50%;
	}
	
	input[type="color"]::-moz-color-swatch {
		border: none;
		border-radius: 50%;
	}
</style>
