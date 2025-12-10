<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	
	export let open = false;
	
	const dispatch = createEventDispatcher();
	
	let title = '';
	let layout: 'list' | 'grid' | 'carousel' | 'card' = 'list';
	
	function handleCreate() {
		if (!title.trim() || !layout) return;
		dispatch('create', { title: title.trim(), layout });
		open = false;
	}
	
	function resetForm() {
		title = '';
		layout = 'list';
	}
	
	$: if (!open) {
		resetForm();
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-xl p-8 overflow-hidden" showCloseButton={false}>
		<!-- Header with Icon -->
		<div class="flex items-start gap-3 mb-6">
			<div class="w-12 h-12 bg-gradient-to-br from-emerald-400 to-emerald-500 rounded-xl flex items-center justify-center flex-shrink-0 shadow-lg shadow-emerald-500/20">
				<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
					<path stroke-linecap="round" stroke-linejoin="round" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
				</svg>
			</div>
			<div>
				<h2 class="text-xl font-bold text-gray-900">Link block</h2>
				<p class="text-gray-500 text-sm mt-0.5">Select link layout</p>
			</div>
		</div>

		<!-- Layout Grid -->
		<div class="grid grid-cols-2 gap-3 mb-5">
			<!-- Classic Layout -->
			<button
				type="button"
				onclick={() => layout = 'list'}
				class="group relative bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'list'}
				class:ring-emerald-500={layout === 'list'}
				class:ring-offset-2={layout === 'list'}
			>
				<div class="aspect-[4/3] p-4 flex flex-col justify-center gap-2">
					<!-- Classic list items with icon -->
					<div class="flex items-center gap-2 bg-white rounded-lg p-2 shadow-sm">
						<div class="w-7 h-7 bg-gradient-to-br from-amber-200 to-amber-300 rounded-md flex-shrink-0"></div>
						<div class="flex-1 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
					<div class="flex items-center gap-2 bg-white rounded-lg p-2 shadow-sm">
						<div class="w-7 h-7 bg-gradient-to-br from-amber-200 to-amber-300 rounded-md flex-shrink-0"></div>
						<div class="flex-1 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
					<div class="flex items-center gap-2 bg-white rounded-lg p-2 shadow-sm">
						<div class="w-7 h-7 bg-gradient-to-br from-amber-200 to-amber-300 rounded-md flex-shrink-0"></div>
						<div class="flex-1 space-y-1">
							<div class="h-1.5 bg-gray-200 rounded-full w-3/4"></div>
							<div class="h-1.5 bg-gray-150 rounded-full w-1/2"></div>
						</div>
					</div>
				</div>
				<div class="p-3 text-center bg-white/50">
					<p class="text-sm font-semibold text-gray-700 group-hover:text-gray-900">Classic</p>
				</div>
			</button>

			<!-- Carousel Layout -->
			<button
				type="button"
				onclick={() => layout = 'carousel'}
				class="group relative bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'carousel'}
				class:ring-emerald-500={layout === 'carousel'}
				class:ring-offset-2={layout === 'carousel'}
			>
				<div class="aspect-[4/3] p-4 flex items-center justify-center gap-2">
					<!-- Carousel cards -->
					<div class="w-16 h-20 bg-gradient-to-br from-amber-100 to-orange-100 rounded-lg shadow-sm transform -rotate-6 opacity-60"></div>
					<div class="w-18 h-24 bg-gradient-to-br from-blue-100 to-cyan-100 rounded-lg shadow-md z-10"></div>
					<div class="w-16 h-20 bg-gradient-to-br from-green-100 to-emerald-100 rounded-lg shadow-sm transform rotate-6 opacity-60"></div>
				</div>
				<div class="p-3 text-center bg-white/50">
					<p class="text-sm font-semibold text-gray-700 group-hover:text-gray-900">Carousel</p>
				</div>
			</button>

			<!-- Image Grid Layout -->
			<button
				type="button"
				onclick={() => layout = 'grid'}
				class="group relative bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'grid'}
				class:ring-emerald-500={layout === 'grid'}
				class:ring-offset-2={layout === 'grid'}
			>
				<div class="aspect-[4/3] p-4">
					<!-- 3x3 Grid -->
					<div class="grid grid-cols-3 gap-1.5 h-full">
						<div class="bg-gradient-to-br from-orange-200 to-red-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-gray-200 to-gray-300 rounded-md"></div>
						<div class="bg-gradient-to-br from-amber-200 to-yellow-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-red-200 to-pink-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-yellow-200 to-amber-300 rounded-md"></div>
						<div class="bg-gradient-to-br from-amber-100 to-orange-200 rounded-md"></div>
						<div class="bg-gradient-to-br from-orange-300 to-amber-400 rounded-md"></div>
						<div class="bg-gradient-to-br from-gray-300 to-slate-300 rounded-md"></div>
						<div class="bg-gradient-to-br from-blue-200 to-cyan-200 rounded-md"></div>
					</div>
				</div>
				<div class="p-3 text-center bg-white/50">
					<p class="text-sm font-semibold text-gray-700 group-hover:text-gray-900">Image grid</p>
				</div>
			</button>

			<!-- Card Layout -->
			<button
				type="button"
				onclick={() => layout = 'card'}
				class="group relative bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all hover:shadow-md overflow-hidden"
				class:ring-2={layout === 'card'}
				class:ring-emerald-500={layout === 'card'}
				class:ring-offset-2={layout === 'card'}
			>
				<div class="aspect-[4/3] p-4 space-y-2">
					<!-- Card style with image on side -->
					<div class="bg-white rounded-lg overflow-hidden shadow-sm flex">
						<div class="w-14 h-14 bg-gradient-to-br from-amber-200 to-orange-200 flex-shrink-0"></div>
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
						<div class="w-14 h-14 bg-gradient-to-br from-blue-200 to-cyan-200 flex-shrink-0"></div>
					</div>
				</div>
				<div class="p-3 text-center bg-white/50">
					<p class="text-sm font-semibold text-gray-700 group-hover:text-gray-900">Card</p>
				</div>
			</button>
		</div>

		<!-- Group Title Input -->
		<div class="mb-5">
			<input
				type="text"
				bind:value={title}
				placeholder="Enter group name..."
				class="w-full px-4 py-3 bg-gray-50 rounded-xl focus:bg-white focus:ring-2 focus:ring-emerald-500/30 transition-all text-gray-900 placeholder-gray-400 outline-none"
				onkeydown={(e) => e.key === 'Enter' && handleCreate()}
			/>
		</div>

		<!-- Action Buttons -->
		<div class="flex gap-3">
			<button 
				type="button"
				onclick={() => open = false}
				class="flex-1 px-5 py-2.5 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-all font-medium"
			>
				Cancel
			</button>
			<button 
				type="button"
				onclick={handleCreate}
				disabled={!title.trim() || !layout}
				class="flex-1 px-5 py-2.5 bg-emerald-500 hover:bg-emerald-600 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all font-medium shadow-lg shadow-emerald-500/25 disabled:shadow-none"
			>
				Create
			</button>
		</div>
	</Dialog.Content>
</Dialog.Root>
