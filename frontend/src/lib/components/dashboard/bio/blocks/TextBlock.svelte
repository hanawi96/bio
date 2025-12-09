<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let block: any;
	export let selected: boolean = false;
	
	const dispatch = createEventDispatcher();

	// Parse style JSON if exists
	let styleConfig: any = {};
	try {
		if (block.style) {
			styleConfig = JSON.parse(block.style);
		}
	} catch (e) {
		console.error('Failed to parse style:', e);
	}

	function handleDelete() {
		dispatch('delete', block.id);
	}

	function handleSelect() {
		dispatch('select', block.id);
	}

	function handleEdit() {
		dispatch('edit', block);
	}

	function handleToggle() {
		dispatch('toggle', block.id);
	}

	function handleDuplicate() {
		dispatch('duplicate', block.id);
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
					class="w-4 h-4 rounded border-2 border-gray-300 cursor-pointer transition-all checked:bg-indigo-600 checked:border-indigo-600 hover:border-indigo-400 focus:ring-2 focus:ring-indigo-500/20 focus:ring-offset-0 accent-indigo-600"
					style="color-scheme: light;"
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

			<div class="flex-shrink-0 w-12 h-12 rounded-lg bg-purple-100 flex items-center justify-center">
				<svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"/>
				</svg>
			</div>

			<div class="flex-1 min-w-0">
				<div class="flex items-center gap-2 mb-1">
					<h3 class="text-sm font-medium text-gray-900 truncate">
						{block.text_style === 'heading' ? 'Heading' : 'Paragraph'}
					</h3>
				</div>
				<div class="flex items-center gap-2 text-xs text-gray-500">
					<div class="flex items-center gap-1">
						<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
						</svg>
						<span class="truncate max-w-[200px]">{block.content || 'Empty'}</span>
					</div>
					{#if styleConfig.backgroundColor}
						<div class="flex items-center gap-1">
							<div class="w-3 h-3 rounded border border-gray-300" style="background-color: {styleConfig.backgroundColor}"></div>
						</div>
					{/if}
				</div>
			</div>

			<div class="flex items-center gap-0.5">
				<button 
					onclick={handleEdit}
					class="p-1 rounded text-gray-400 hover:text-gray-600 hover:bg-gray-50"
					title="Edit"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
					</svg>
				</button>

				<button 
					onclick={handleToggle}
					class="p-1 rounded {block.is_active ? 'text-green-600 hover:bg-green-50' : 'text-gray-400 hover:bg-gray-50'}"
					title={block.is_active ? 'Hide' : 'Show'}
				>
					{#if block.is_active}
						<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
						</svg>
					{:else}
						<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
						</svg>
					{/if}
				</button>

				<button 
					onclick={handleDuplicate}
					class="p-1 rounded text-gray-400 hover:text-gray-600 hover:bg-gray-50"
					title="Duplicate"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
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
