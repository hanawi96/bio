<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	
	export let open = false;
	
	const dispatch = createEventDispatcher();
	
	let title = '';
	let layout: 'list' | 'grid' | 'carousel' = 'list';
	
	function handleCreate() {
		if (!title.trim()) return;
		dispatch('create', { title: title.trim(), layout });
	}
	
	function resetForm() {
		title = '';
		layout = 'list';
	}
	
	// Reset form when dialog closes
	$: if (!open) {
		resetForm();
	}
	
	function handleClose() {
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Create Link Group</Dialog.Title>
			<Dialog.Description>
				Organize multiple links together in a group with custom layout
			</Dialog.Description>
		</Dialog.Header>
		
		<div class="space-y-4 py-4">
			<!-- Group Title -->
			<div class="space-y-2">
				<Label for="group-title">Group Title</Label>
				<Input
					id="group-title"
					bind:value={title}
					placeholder="e.g., Social Media, Contact Links"
					class="w-full"
					onkeydown={(e) => e.key === 'Enter' && handleCreate()}
				/>
			</div>

			<!-- Layout Selection -->
			<div class="space-y-2">
				<Label>Display Layout</Label>
				<div class="grid grid-cols-3 gap-3">
					<!-- List Layout -->
					<button
						type="button"
						onclick={() => layout = 'list'}
						class="relative p-4 border-2 rounded-lg transition-all hover:border-indigo-300"
						class:border-indigo-600={layout === 'list'}
						class:bg-indigo-50={layout === 'list'}
						class:border-gray-200={layout !== 'list'}
					>
						<div class="space-y-1.5 mb-2">
							<div class="h-2 bg-gray-300 rounded"></div>
							<div class="h-2 bg-gray-300 rounded"></div>
							<div class="h-2 bg-gray-300 rounded"></div>
						</div>
						<p class="text-xs font-medium text-gray-700">List</p>
						{#if layout === 'list'}
							<div class="absolute top-2 right-2 w-5 h-5 bg-indigo-600 rounded-full flex items-center justify-center">
								<svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
								</svg>
							</div>
						{/if}
					</button>

					<!-- Grid Layout -->
					<button
						type="button"
						onclick={() => layout = 'grid'}
						class="relative p-4 border-2 rounded-lg transition-all hover:border-indigo-300"
						class:border-indigo-600={layout === 'grid'}
						class:bg-indigo-50={layout === 'grid'}
						class:border-gray-200={layout !== 'grid'}
					>
						<div class="grid grid-cols-2 gap-1.5 mb-2">
							<div class="h-4 bg-gray-300 rounded"></div>
							<div class="h-4 bg-gray-300 rounded"></div>
							<div class="h-4 bg-gray-300 rounded"></div>
							<div class="h-4 bg-gray-300 rounded"></div>
						</div>
						<p class="text-xs font-medium text-gray-700">Grid</p>
						{#if layout === 'grid'}
							<div class="absolute top-2 right-2 w-5 h-5 bg-indigo-600 rounded-full flex items-center justify-center">
								<svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
								</svg>
							</div>
						{/if}
					</button>

					<!-- Carousel Layout -->
					<button
						type="button"
						onclick={() => layout = 'carousel'}
						class="relative p-4 border-2 rounded-lg transition-all hover:border-indigo-300"
						class:border-indigo-600={layout === 'carousel'}
						class:bg-indigo-50={layout === 'carousel'}
						class:border-gray-200={layout !== 'carousel'}
					>
						<div class="flex gap-1 mb-2 overflow-hidden">
							<div class="w-8 h-6 bg-gray-300 rounded flex-shrink-0"></div>
							<div class="w-8 h-6 bg-gray-400 rounded flex-shrink-0"></div>
							<div class="w-8 h-6 bg-gray-300 rounded flex-shrink-0"></div>
						</div>
						<p class="text-xs font-medium text-gray-700">Carousel</p>
						{#if layout === 'carousel'}
							<div class="absolute top-2 right-2 w-5 h-5 bg-indigo-600 rounded-full flex items-center justify-center">
								<svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
								</svg>
							</div>
						{/if}
					</button>
				</div>
			</div>

			<!-- Preview Info -->
			<div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
				<div class="flex gap-2">
					<svg class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
					</svg>
					<div class="text-sm text-blue-800">
						<p class="font-medium mb-1">After creating the group:</p>
						<ul class="text-xs space-y-0.5 text-blue-700">
							<li>• Add links directly to the group</li>
							<li>• Or move existing links into it</li>
							<li>• Change layout anytime</li>
						</ul>
					</div>
				</div>
			</div>
		</div>

		<Dialog.Footer>
			<button 
				type="button"
				onclick={handleClose}
				class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50"
			>
				Cancel
			</button>
			<button 
				type="button"
				onclick={handleCreate}
				disabled={!title.trim()}
				class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
				</svg>
				Create Group
			</button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
