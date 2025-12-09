<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { Link } from '$lib/api/links';
	
	export let open = false;
	export let link: Link | null = null;
	export let groupTitle = '';
	
	const dispatch = createEventDispatcher();
	
	let title = '';
	let url = '';
	let urlError = '';
	
	// Mode: 'add' or 'edit'
	$: mode = link ? 'edit' : 'add';
	
	$: if (link && open) {
		title = link.title;
		url = link.url;
		urlError = '';
	}
	
	function isValidUrl(urlString: string): boolean {
		try {
			const urlObj = new URL(urlString);
			return urlObj.protocol === 'http:' || urlObj.protocol === 'https:';
		} catch {
			return false;
		}
	}
	
	function validateUrl() {
		if (!url.trim()) {
			urlError = '';
			return;
		}
		
		if (!isValidUrl(url.trim())) {
			urlError = 'Please enter a valid URL (e.g., https://example.com)';
		} else {
			urlError = '';
		}
	}
	
	function handleSave() {
		if (!title.trim() || !url.trim()) return;
		
		validateUrl();
		if (urlError) return;
		
		if (mode === 'edit' && link) {
			dispatch('save', {
				id: link.id,
				title: title.trim(),
				url: url.trim()
			});
		} else {
			dispatch('add', {
				title: title.trim(),
				url: url.trim()
			});
		}
		
		open = false;
	}
	
	function resetForm() {
		title = '';
		url = '';
		urlError = '';
	}
	
	$: if (!open) {
		resetForm();
	}
	
	function handleClose() {
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>{mode === 'edit' ? 'Edit Link' : 'Add Link to Group'}</Dialog.Title>
			<Dialog.Description>
				{#if mode === 'edit'}
					Customize this link's appearance and details
				{:else}
					Adding to: <span class="font-semibold text-gray-900">{groupTitle}</span>
				{/if}
			</Dialog.Description>
		</Dialog.Header>
		
		<div class="space-y-4 py-4">
			<!-- Link Title -->
			<div class="space-y-2">
				<Label for="edit-link-title">Link Title</Label>
				<Input
					id="edit-link-title"
					bind:value={title}
					placeholder="e.g., Facebook, Twitter"
					class="w-full"
				/>
			</div>

			<!-- Link URL -->
			<div class="space-y-2">
				<Label for="edit-link-url">URL</Label>
				<Input
					id="edit-link-url"
					type="url"
					bind:value={url}
					placeholder="https://example.com"
					class={urlError ? 'w-full border-red-500' : 'w-full'}
					onblur={validateUrl}
				/>
				{#if urlError}
					<p class="text-xs text-red-600">{urlError}</p>
				{/if}
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
				onclick={handleSave}
				disabled={!title.trim() || !url.trim() || !!urlError}
				class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
			>
				{#if mode === 'add'}
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					Add Link
				{:else}
					Save Changes
				{/if}
			</button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
