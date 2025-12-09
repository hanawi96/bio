<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let block: any;
	export let selected: boolean = false;
	
	const dispatch = createEventDispatcher();
	let showEditDialog = false;
	let editedUrl = block.embed_url || '';
	let editedType = block.embed_type || 'spotify';

	function handleDelete() {
		dispatch('delete', block.id);
	}

	function handleSelect() {
		dispatch('select', block.id);
	}

	function handleSave() {
		dispatch('update', { id: block.id, embed_url: editedUrl, embed_type: editedType });
		showEditDialog = false;
	}

	function getEmbedIcon(type: string) {
		switch(type) {
			case 'spotify': return 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm3.61 14.34c-.23.39-.72.51-1.11.28-3.06-1.87-6.91-2.29-11.44-1.25-.42.1-.85-.13-.95-.56-.1-.42.13-.85.56-.95 4.98-1.15 9.27-.65 12.64 1.45.39.23.51.72.28 1.11zm1.58-3.51c-.29.47-.91.61-1.38.32-3.5-2.15-8.83-2.77-12.97-1.52-.51.16-1.05-.13-1.21-.64-.16-.51.13-1.05.64-1.21 4.72-1.43 10.61-.74 14.61 1.73.47.29.61.91.32 1.38zm.14-3.65c-4.2-2.49-11.12-2.72-15.12-1.5-.61.19-1.26-.15-1.45-.76-.19-.61.15-1.26.76-1.45 4.58-1.4 12.23-1.13 17.09 1.74.56.33.75 1.06.42 1.62-.33.56-1.06.75-1.62.42z';
			case 'soundcloud': return 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z';
			case 'maps': return 'M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z';
			default: return 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z';
		}
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
					class="w-4 h-4 rounded border-2 border-gray-300 cursor-pointer transition-all checked:bg-indigo-600 checked:border-indigo-600"
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

			<div class="flex-shrink-0 w-12 h-12 rounded-lg bg-gradient-to-br from-purple-100 to-pink-100 flex items-center justify-center">
				<svg class="w-5 h-5 text-purple-600" fill="currentColor" viewBox="0 0 24 24">
					<path d={getEmbedIcon(block.embed_type || 'spotify')}/>
				</svg>
			</div>

			<div class="flex-1 min-w-0">
				<h3 class="text-sm font-medium text-gray-900 mb-1">Embed - {block.embed_type || 'Spotify'}</h3>
				<p class="text-xs text-gray-500 truncate">{block.embed_url || 'No URL'}</p>
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
				<h3 class="text-lg font-semibold text-gray-900">Edit Embed</h3>
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
					<label class="block text-sm font-medium text-gray-700 mb-2">Type</label>
					<select
						bind:value={editedType}
						class="w-full px-4 py-2.5 bg-gray-50 border-0 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					>
						<option value="spotify">Spotify</option>
						<option value="soundcloud">SoundCloud</option>
						<option value="maps">Google Maps</option>
						<option value="other">Other</option>
					</select>
				</div>

				<div>
					<label class="block text-sm font-medium text-gray-700 mb-2">Embed URL</label>
					<input
						type="text"
						bind:value={editedUrl}
						placeholder="https://open.spotify.com/embed/..."
						class="w-full px-4 py-2.5 bg-gray-50 border-0 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					/>
					<p class="text-xs text-gray-500 mt-2">Paste the embed URL from the service</p>
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
