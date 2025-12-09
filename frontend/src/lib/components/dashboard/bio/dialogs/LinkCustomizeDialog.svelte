<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';

	export let open = false;
	export let link: Link;

	const dispatch = createEventDispatcher();
	
	let layoutType: 'classic' | 'carousel' | 'grid' | 'card' = 'classic';
	let imagePlacement: 'left' | 'right' | 'top' | 'bottom' = 'left';
	let textAlignment: 'left' | 'center' | 'right' = 'left';
	let textSize: 'S' | 'M' | 'L' | 'XL' = 'M';
	let showOutline = false;
	let showShadow = false;

	$: if (open && link) {
		layoutType = (link.layout_type as any) || 'classic';
		imagePlacement = (link.image_placement as any) || 'left';
		textAlignment = (link.text_alignment as any) || 'left';
		textSize = (link.text_size as any) || 'M';
		showOutline = link.show_outline || false;
		showShadow = link.show_shadow || false;
	}

	function handleSave() {
		dispatch('updateLayout', {
			id: link.id,
			layout_type: layoutType,
			image_placement: imagePlacement,
			text_alignment: textAlignment,
			text_size: textSize,
			show_outline: showOutline,
			show_shadow: showShadow
		});
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-4xl max-h-[90vh] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title class="text-xl font-bold">Customize Link</Dialog.Title>
		</Dialog.Header>

		<div class="space-y-8 py-4">
			<!-- Layout -->
			<div>
				<h3 class="text-base font-bold text-gray-900 mb-4">Layout</h3>
				<div class="grid grid-cols-4 gap-4">
					<button
						onclick={() => layoutType = 'classic'}
						class="flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all hover:border-gray-400"
						class:border-gray-900={layoutType === 'classic'}
						class:bg-gray-50={layoutType === 'classic'}
						class:border-gray-200={layoutType !== 'classic'}
					>
						<div class="w-full aspect-square bg-white rounded-lg border border-gray-200 p-3 flex flex-col gap-2">
							<div class="h-2 bg-gray-300 rounded"></div>
							<div class="h-2 bg-gray-300 rounded"></div>
							<div class="h-2 bg-gray-300 rounded w-3/4"></div>
						</div>
						<span class="text-sm font-medium" class:text-gray-900={layoutType === 'classic'} class:text-gray-500={layoutType !== 'classic'}>Classic</span>
					</button>

					<button
						onclick={() => layoutType = 'carousel'}
						class="flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all hover:border-gray-400"
						class:border-gray-900={layoutType === 'carousel'}
						class:bg-gray-50={layoutType === 'carousel'}
						class:border-gray-200={layoutType !== 'carousel'}
					>
						<div class="w-full aspect-square bg-white rounded-lg border border-gray-200 p-3 flex gap-2">
							<div class="flex-1 bg-gray-300 rounded"></div>
							<div class="flex-1 bg-gray-300 rounded"></div>
							<div class="flex-1 bg-gray-300 rounded"></div>
						</div>
						<span class="text-sm font-medium" class:text-gray-900={layoutType === 'carousel'} class:text-gray-500={layoutType !== 'carousel'}>Carousel</span>
					</button>

					<button
						onclick={() => layoutType = 'grid'}
						class="flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all hover:border-gray-400"
						class:border-gray-900={layoutType === 'grid'}
						class:bg-gray-50={layoutType === 'grid'}
						class:border-gray-200={layoutType !== 'grid'}
					>
						<div class="w-full aspect-square bg-white rounded-lg border border-gray-200 p-3 grid grid-cols-3 gap-1.5">
							{#each Array(9) as _}
								<div class="bg-gray-300 rounded"></div>
							{/each}
						</div>
						<span class="text-sm font-medium" class:text-gray-900={layoutType === 'grid'} class:text-gray-500={layoutType !== 'grid'}>Image grid</span>
					</button>

					<button
						onclick={() => layoutType = 'card'}
						class="flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all hover:border-gray-400"
						class:border-gray-900={layoutType === 'card'}
						class:bg-gray-50={layoutType === 'card'}
						class:border-gray-200={layoutType !== 'card'}
					>
						<div class="w-full aspect-square bg-white rounded-lg border border-gray-200 p-3 grid grid-cols-2 gap-2">
							<div class="bg-gray-900 rounded"></div>
							<div class="bg-gray-300 rounded"></div>
							<div class="bg-gray-300 rounded"></div>
							<div class="bg-gray-300 rounded"></div>
						</div>
						<span class="text-sm font-medium" class:text-gray-900={layoutType === 'card'} class:text-gray-500={layoutType !== 'card'}>Card</span>
					</button>
				</div>
			</div>

			<!-- Image Placement -->
			<div>
				<h3 class="text-base font-bold text-gray-900 mb-4">Image placement</h3>
				<select
					bind:value={imagePlacement}
					class="w-full px-4 py-3 bg-white border-2 border-gray-200 rounded-xl text-sm focus:outline-none focus:border-gray-900 transition-colors"
				>
					<option value="left">Left</option>
					<option value="right">Right</option>
					<option value="top">Top</option>
					<option value="bottom">Bottom</option>
				</select>
			</div>

			<!-- Text Alignment -->
			<div>
				<h3 class="text-base font-bold text-gray-900 mb-4">Text alignment</h3>
				<div class="flex gap-3">
					<button
						onclick={() => textAlignment = 'left'}
						class="flex-1 p-4 rounded-xl border-2 transition-all hover:border-gray-400"
						class:border-gray-900={textAlignment === 'left'}
						class:bg-gray-50={textAlignment === 'left'}
						class:border-gray-200={textAlignment !== 'left'}
					>
						<svg class="w-6 h-6 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h10M4 18h14"/>
						</svg>
					</button>
					<button
						onclick={() => textAlignment = 'center'}
						class="flex-1 p-4 rounded-xl border-2 transition-all hover:border-gray-400"
						class:border-gray-900={textAlignment === 'center'}
						class:bg-gray-50={textAlignment === 'center'}
						class:border-gray-200={textAlignment !== 'center'}
					>
						<svg class="w-6 h-6 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M7 12h10M5 18h14"/>
						</svg>
					</button>
					<button
						onclick={() => textAlignment = 'right'}
						class="flex-1 p-4 rounded-xl border-2 transition-all hover:border-gray-400"
						class:border-gray-900={textAlignment === 'right'}
						class:bg-gray-50={textAlignment === 'right'}
						class:border-gray-200={textAlignment !== 'right'}
					>
						<svg class="w-6 h-6 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M10 12h10M6 18h14"/>
						</svg>
					</button>
				</div>
			</div>

			<!-- Text Size -->
			<div>
				<h3 class="text-base font-bold text-gray-900 mb-4">Text size</h3>
				<div class="flex gap-3">
					{#each ['S', 'M', 'L', 'XL'] as size}
						<button
							onclick={() => textSize = size}
							class="flex-1 p-4 rounded-xl border-2 transition-all hover:border-gray-400 font-bold"
							class:border-gray-900={textSize === size}
							class:bg-gray-50={textSize === size}
							class:border-gray-200={textSize !== size}
						>
							{size}
						</button>
					{/each}
				</div>
			</div>

			<!-- Toggles -->
			<div class="space-y-3">
				<label class="flex items-center justify-between p-4 rounded-xl border-2 border-gray-200 hover:border-gray-300 cursor-pointer transition-colors">
					<span class="text-sm font-medium text-gray-900">Link outline</span>
					<input
						type="checkbox"
						bind:checked={showOutline}
						class="w-12 h-6 rounded-full appearance-none cursor-pointer transition-colors relative"
						class:bg-gray-900={showOutline}
						class:bg-gray-300={!showOutline}
						style="box-shadow: inset 0 0 0 2px transparent;"
					/>
				</label>

				<label class="flex items-center justify-between p-4 rounded-xl border-2 border-gray-200 hover:border-gray-300 cursor-pointer transition-colors">
					<span class="text-sm font-medium text-gray-900">Link shadow</span>
					<input
						type="checkbox"
						bind:checked={showShadow}
						class="w-12 h-6 rounded-full appearance-none cursor-pointer transition-colors relative"
						class:bg-gray-900={showShadow}
						class:bg-gray-300={!showShadow}
					/>
				</label>
			</div>
		</div>

		<Dialog.Footer class="flex justify-end gap-3">
			<button
				onclick={() => open = false}
				class="px-6 py-2.5 text-sm font-medium text-gray-700 hover:text-gray-900 transition-colors"
			>
				Cancel
			</button>
			<button
				onclick={handleSave}
				class="px-6 py-2.5 bg-gray-900 hover:bg-gray-800 text-white text-sm font-medium rounded-xl transition-all"
			>
				Save
			</button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<style>
	input[type="checkbox"]:checked::after {
		content: '';
		position: absolute;
		right: 2px;
		top: 2px;
		width: 20px;
		height: 20px;
		background: white;
		border-radius: 50%;
		transition: all 0.2s;
	}

	input[type="checkbox"]:not(:checked)::after {
		content: '';
		position: absolute;
		left: 2px;
		top: 2px;
		width: 20px;
		height: 20px;
		background: white;
		border-radius: 50%;
		transition: all 0.2s;
	}
</style>
