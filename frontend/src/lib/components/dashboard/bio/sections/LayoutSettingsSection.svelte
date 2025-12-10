<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';
	
	export let group: Link;
	
	const dispatch = createEventDispatcher();
	
	let layout = group.group_layout || 'list';
	let textAlignment = group.text_alignment || 'center';
	let showOutline = group.show_outline || false;
	let showShadow = group.show_shadow || false;
	
	let updateTimeout: ReturnType<typeof setTimeout>;
	
	function handleUpdate(field: string, value: any) {
		clearTimeout(updateTimeout);
		updateTimeout = setTimeout(() => {
			dispatch('update', {
				groupId: group.id,
				[field]: value
			});
		}, 500);
	}
	
	$: handleUpdate('group_layout', layout);
	$: handleUpdate('text_alignment', textAlignment);
	$: handleUpdate('show_outline', showOutline);
	$: handleUpdate('show_shadow', showShadow);
</script>

<div class="p-6 space-y-6">
	<!-- Layout Section -->
	<div>
		<h3 class="text-base font-semibold text-gray-900 mb-3">Layout</h3>
		<div class="flex gap-3">
			<!-- Classic -->
			<button
				type="button"
				onclick={() => layout = 'list'}
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'list'}
				class:ring-emerald-500={layout === 'list'}
				class:ring-offset-2={layout === 'list'}
			>
				<div class="aspect-[4/3] p-4 flex flex-col justify-center gap-2">
					<div class="flex items-center gap-2 bg-white rounded-lg p-2 shadow-sm">
						<div class="w-6 h-6 bg-gradient-to-br from-amber-200 to-amber-300 rounded flex-shrink-0"></div>
						<div class="flex-1 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
					<div class="flex items-center gap-2 bg-white rounded-lg p-2 shadow-sm">
						<div class="w-6 h-6 bg-gradient-to-br from-amber-200 to-amber-300 rounded flex-shrink-0"></div>
						<div class="flex-1 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
					<div class="flex items-center gap-2 bg-white rounded-lg p-2 shadow-sm">
						<div class="w-6 h-6 bg-gradient-to-br from-amber-200 to-amber-300 rounded flex-shrink-0"></div>
						<div class="flex-1 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
				</div>
				<div class="p-2 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-700">Classic</p>
				</div>
			</button>

			<!-- Carousel -->
			<button
				type="button"
				onclick={() => layout = 'carousel'}
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'carousel'}
				class:ring-emerald-500={layout === 'carousel'}
				class:ring-offset-2={layout === 'carousel'}
			>
				<div class="aspect-[4/3] p-4 flex items-center justify-center gap-2">
					<div class="w-16 h-20 bg-gradient-to-br from-yellow-100 to-yellow-200 rounded-lg shadow-sm opacity-60"></div>
					<div class="w-20 h-24 bg-gradient-to-br from-emerald-100 to-emerald-200 rounded-lg shadow-md"></div>
				</div>
				<div class="p-2 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-700">Carousel</p>
				</div>
			</button>

			<!-- Image Grid -->
			<button
				type="button"
				onclick={() => layout = 'grid'}
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'grid'}
				class:ring-emerald-500={layout === 'grid'}
				class:ring-offset-2={layout === 'grid'}
			>
				<div class="aspect-[4/3] p-4">
					<div class="grid grid-cols-3 gap-1.5 h-full">
						<div class="bg-gradient-to-br from-orange-200 to-red-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-gray-200 to-gray-300 rounded-md"></div>
						<div class="bg-gradient-to-br from-yellow-200 to-amber-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-pink-200 to-rose-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-yellow-300 to-amber-300 rounded-md"></div>
						<div class="bg-gradient-to-br from-amber-100 to-orange-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-orange-300 to-amber-400 rounded-md"></div>
						<div class="bg-gradient-to-br from-gray-300 to-slate-300 rounded-md"></div>
						<div class="bg-gradient-to-br from-blue-200 to-cyan-200 rounded-md"></div>
					</div>
				</div>
				<div class="p-2 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-700">Image grid</p>
				</div>
			</button>

			<!-- Card (disabled) -->
			<button
				type="button"
				disabled
				class="flex-1 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl overflow-hidden opacity-50 cursor-not-allowed"
			>
				<div class="aspect-[4/3] p-4 space-y-2 flex flex-col justify-center">
					<div class="bg-white rounded-lg overflow-hidden shadow-sm flex">
						<div class="w-12 h-12 bg-gradient-to-br from-amber-200 to-orange-200 flex-shrink-0"></div>
						<div class="flex-1 p-2 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-4/5"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-3/5"></div>
						</div>
					</div>
					<div class="bg-white rounded-lg overflow-hidden shadow-sm flex">
						<div class="flex-1 p-2 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-4/5"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-3/5"></div>
						</div>
						<div class="w-12 h-12 bg-gradient-to-br from-blue-200 to-cyan-200 flex-shrink-0"></div>
					</div>
				</div>
				<div class="p-2 text-center bg-white/60">
					<p class="text-xs font-semibold text-gray-400">Card</p>
				</div>
			</button>
		</div>
	</div>

	<!-- Text Alignment Section -->
	<div>
		<h3 class="text-base font-semibold text-gray-900 mb-3">Text alignment</h3>
		<div class="flex gap-2">
			<button
				type="button"
				onclick={() => textAlignment = 'left'}
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
				onclick={() => textAlignment = 'center'}
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
				onclick={() => textAlignment = 'right'}
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

	<!-- Toggles Section -->
	<div class="space-y-3">
		<label class="flex items-center justify-between cursor-pointer">
			<span class="text-base font-medium text-gray-900">Link outline</span>
			<button
				type="button"
				onclick={() => showOutline = !showOutline}
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
				onclick={() => showShadow = !showShadow}
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
</div>
