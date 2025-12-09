<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { Link } from '$lib/api/links';
	import ImageUploader from '$lib/components/ui/ImageUploader.svelte';
	
	export let open = false;
	export let link: Link | null = null;
	export let groupTitle = '';
	
	const dispatch = createEventDispatcher();
	
	let title = '';
	let url = '';
	let thumbnailUrl: string | null = null;
	let uploadingThumbnail = false;
	
	// Mode: 'add' or 'edit'
	$: mode = link ? 'edit' : 'add';
	
	$: if (link && open) {
		title = link.title;
		url = link.url;
		thumbnailUrl = link.thumbnail_url || null;
	}
	
	function handleSave() {
		if (!title.trim() || !url.trim()) return;
		
		if (mode === 'edit' && link) {
			dispatch('save', {
				id: link.id,
				title: title.trim(),
				url: url.trim(),
				thumbnail_url: thumbnailUrl
			});
		} else {
			dispatch('add', {
				title: title.trim(),
				url: url.trim()
			});
		}
		
		open = false;
	}
	
	async function handleUploadThumbnail(event: CustomEvent) {
		const file = event.detail as File;
		uploadingThumbnail = true;
		
		// Dispatch to parent to handle actual upload
		dispatch('uploadThumbnail', { linkId: link?.id, file });
	}
	
	// Reset uploading state when thumbnail changes (upload complete)
	$: if (link?.thumbnail_url && uploadingThumbnail && thumbnailUrl !== link.thumbnail_url) {
		thumbnailUrl = link.thumbnail_url;
		uploadingThumbnail = false;
	}
	
	function handleRemoveThumbnail() {
		thumbnailUrl = null;
		dispatch('removeThumbnail', link?.id);
	}
	
	function resetForm() {
		title = '';
		url = '';
		thumbnailUrl = null;
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
					class="w-full"
				/>
			</div>

			<!-- Thumbnail (only in edit mode) -->
			{#if mode === 'edit'}
				<div class="space-y-2">
					<Label>Thumbnail Image</Label>
					<ImageUploader 
						currentImage={thumbnailUrl}
						uploading={uploadingThumbnail}
						on:upload={handleUploadThumbnail}
						on:remove={handleRemoveThumbnail}
					/>
					<p class="text-xs text-gray-500">Recommended size: 400x400px</p>
				</div>
			{/if}
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
				disabled={!title.trim() || !url.trim()}
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
