<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	
	export let open = false;
	export let groupTitle = '';
	
	const dispatch = createEventDispatcher();
	
	let title = '';
	let url = '';
	
	function handleAdd() {
		if (!title.trim() || !url.trim()) return;
		
		dispatch('add', { 
			title: title.trim(), 
			url: url.trim() 
		});
	}
	
	function resetForm() {
		title = '';
		url = '';
	}
	
	// Reset form when dialog closes
	$: if (!open) {
		resetForm();
	}
	
	function handleClose() {
		open = false;
	}
</script>

<Dialog.Root {open} onOpenChange={(newOpen) => open = newOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Add Link to Group</Dialog.Title>
			<Dialog.Description>
				Adding to: <span class="font-semibold text-gray-900">{groupTitle}</span>
			</Dialog.Description>
		</Dialog.Header>
		
		<div class="space-y-4 py-4">
			<!-- Link Title -->
			<div class="space-y-2">
				<Label for="link-title">Link Title</Label>
				<Input
					id="link-title"
					bind:value={title}
					placeholder="e.g., Facebook, Twitter"
					class="w-full"
				/>
			</div>

			<!-- Link URL -->
			<div class="space-y-2">
				<Label for="link-url">URL</Label>
				<Input
					id="link-url"
					type="url"
					bind:value={url}
					placeholder="https://example.com"
					class="w-full"
					onkeydown={(e) => e.key === 'Enter' && handleAdd()}
				/>
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
				onclick={handleAdd}
				disabled={!title.trim() || !url.trim()}
				class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
				</svg>
				Add Link
			</button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
