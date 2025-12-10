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

<div class="relative bg-white rounded-xl shadow-sm hover:shadow-md transition-all {selected ? 'ring-2 ring-indigo-500' : ''}">
	<div class="relative p-3.5">
		<div class="flex items-center gap-4">
			<!-- Thumbnail/Icon -->
			<div class="flex-shrink-0 w-14 h-14 rounded-lg bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center shadow-sm">
				<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"/>
				</svg>
			</div>

			<!-- Content -->
			<button
				onclick={handleEdit}
				class="flex-1 text-left min-w-0"
			>
				<h3 class="text-base font-semibold text-gray-900 truncate">
					{block.text_style === 'heading' ? 'Heading' : 'Paragraph'}
				</h3>
				<p class="text-sm text-gray-500 mt-0.5 truncate">
					{block.content || 'Empty text'}
				</p>
			</button>

			<!-- Actions (Linktree Style) -->
			<div class="flex items-center gap-2">
				<!-- Drag Handle -->
				<button
					class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400 p-1 transition-colors"
				>
					<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
						<circle cx="9" cy="7" r="1.5"/>
						<circle cx="9" cy="12" r="1.5"/>
						<circle cx="9" cy="17" r="1.5"/>
						<circle cx="15" cy="7" r="1.5"/>
						<circle cx="15" cy="12" r="1.5"/>
						<circle cx="15" cy="17" r="1.5"/>
					</svg>
				</button>
			</div>
		</div>
	</div>
</div>
