<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let block: any;
	export let selected: boolean = false;
	
	const dispatch = createEventDispatcher();
	let showEditDialog = false;
	let editedUrl = block.image_url || '';
	let editedAlt = block.alt_text || '';

	function handleDelete() {
		dispatch('delete', block.id);
	}

	function handleSelect() {
		dispatch('select', block.id);
	}

	function handleSave() {
		dispatch('update', { id: block.id, image_url: editedUrl, alt_text: editedAlt });
		showEditDialog = false;
	}
</script>

<div class="relative bg-white rounded-xl {selected ? 'ring-2 ring-indigo-500' : ''}">
	<div class="relative p-4">
		<div class="flex items-center gap-3">
			<label class="flex items-center cursor-pointer">
				<input 
					type="checkbox" 
					checked={selected}
					onchange={handleSelect}
					class="w-4 h-4 rounded border-2 border-gray-300 cursor-pointer transition-all checked:bg-indigo-600 checked:border-indigo-600 hover:border-indigo-400 focus:ring-2 focus:ring-indigo-500/20 accent-indigo-600"
				/>
			</label>
			
			<div class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400">
				<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
					<circle cx="9" cy="7" r="1.5"/>
					<circle cx="9" cy="12" r="1.5"/>
					<circle cx="9" cy="17" r="1.5"/>
					<circle cx="15" cy="7" r="1.5"/>
					<circle cx="15" cy="12" r="1.5"/>
					<circle cx="15" cy="17" r="1.5"/>
				</svg>
			</div>

			{#if block.image_url}
				<div class="flex-shrink-0 w-12 h-12 rounded-lg overflow-hidden bg-gray-100">
					<img src={block.image_url} alt={block.alt_text} class="w-full h-full object-cover"/>
				</div>
			{:else}
				<div class="flex-shrink-0 w-12 h-12 rounded-lg bg-gray-100 flex items-center justify-center">
					<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
					</svg>
				</div>
			{/if}

			<div class="flex-1 min-w-0">
				<h3 class="text-sm font-medium text-gray-900 mb-1">Image</h3>
				<p class="text-xs text-gray-500 truncate">{block.image_url || 'No image'}</p>
			</div>

			<div class="flex items-center gap-0.5">
				<button 
					onclick={() => showEditDialog = true}
					class="p-1 rounded text-gray-400 hover:text-gray-600 hover:bg-gray-50"
					title="Edit"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
					</svg>
				</button>
				
				<button 
					onclick={handleDelete} 
					class="p-1 rounded text-gray-400 hover:text-red-600 hover:bg-red-50"
					title="Delete"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
					</svg>
				</button>
			</div>
		</div>
	</div>
</div>

{#if showEditDialog}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onclick={() => showEditDialog = false}>
		<div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6" onclick={(e) => e.stopPropagation()}>
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-lg font-semibold text-gray-900">Edit Image</h3>
				<button 
					onclick={() => showEditDialog = false}
					class="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
					</svg>
				</button>
			</div>

			<div class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Image URL</label>
					<input
						type="text"
						bind:value={editedUrl}
						placeholder="https://example.com/image.jpg"
						class="w-full px-4 py-2.5 bg-gray-50 border-0 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					/>
				</div>

				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Alt Text</label>
					<input
						type="text"
						bind:value={editedAlt}
						placeholder="Description of image"
						class="w-full px-4 py-2.5 bg-gray-50 border-0 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					/>
				</div>
			</div>

			<div class="flex items-center gap-3 mt-6">
				<button
					onclick={() => showEditDialog = false}
					class="flex-1 px-4 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium rounded-xl"
				>
					Cancel
				</button>
				<button
					onclick={handleSave}
					class="flex-1 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl"
				>
					Save
				</button>
			</div>
		</div>
	</div>
{/if}
