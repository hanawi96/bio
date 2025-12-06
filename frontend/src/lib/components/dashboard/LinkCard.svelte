<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import { createEventDispatcher } from 'svelte';
	import type { Link } from '$lib/api/links';

	export let link: Link;
	
	const dispatch = createEventDispatcher();
	let isEditingTitle = false;
	let isEditingUrl = false;
	let editedTitle = link.title;
	let editedUrl = link.url;

	function handleSaveTitle() {
		if (editedTitle !== link.title) {
			dispatch('update', { id: link.id, title: editedTitle, url: link.url });
		}
		isEditingTitle = false;
	}

	function handleSaveUrl() {
		if (editedUrl !== link.url) {
			dispatch('update', { id: link.id, title: link.title, url: editedUrl });
		}
		isEditingUrl = false;
	}

	function handleDelete() {
		dispatch('delete', link.id);
	}

	function handleToggle() {
		dispatch('toggle', link.id);
	}
</script>

<div class="group relative bg-white rounded-2xl p-6">
	<div class="flex items-start gap-4">
		<!-- Drag Handle -->
		<div class="drag-handle mt-2 cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400">
			<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
				<circle cx="8" cy="6" r="1.5"/>
				<circle cx="8" cy="12" r="1.5"/>
				<circle cx="8" cy="18" r="1.5"/>
				<circle cx="14" cy="6" r="1.5"/>
				<circle cx="14" cy="12" r="1.5"/>
				<circle cx="14" cy="18" r="1.5"/>
			</svg>
		</div>

		<!-- Content -->
		<div class="flex-1 space-y-3">
			<!-- Title -->
			<div class="flex items-center gap-2">
				{#if isEditingTitle}
					<Input 
						bind:value={editedTitle} 
						onblur={handleSaveTitle}
						onkeydown={(e) => e.key === 'Enter' && handleSaveTitle()}
						class="text-xl font-bold"
						autofocus
					/>
				{:else}
					<h3 class="text-xl font-bold text-gray-900">{link.title}</h3>
					<button onclick={() => isEditingTitle = true} class="text-gray-400 hover:text-gray-600">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
						</svg>
					</button>
				{/if}
			</div>

			<!-- URL -->
			<div class="flex items-center gap-2">
				{#if isEditingUrl}
					<Input 
						bind:value={editedUrl} 
						onblur={handleSaveUrl}
						onkeydown={(e) => e.key === 'Enter' && handleSaveUrl()}
						class="text-sm"
						autofocus
					/>
				{:else}
					<p class="text-sm text-gray-600 truncate">{link.url}</p>
					<button onclick={() => isEditingUrl = true} class="text-gray-400 hover:text-gray-600">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
						</svg>
					</button>
				{/if}
			</div>

			<!-- Action Icons -->
			<div class="flex items-center gap-4 pt-2">
				<button class="text-gray-400 hover:text-gray-600 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
					</svg>
				</button>
				<button class="text-gray-400 hover:text-gray-600 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
					</svg>
				</button>
				<button class="text-gray-400 hover:text-gray-600 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
					</svg>
				</button>
				<button class="text-gray-400 hover:text-gray-600 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>
					</svg>
				</button>
				<button class="text-gray-400 hover:text-gray-600 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
					</svg>
				</button>
				<button class="text-gray-400 hover:text-gray-600 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
					</svg>
				</button>
				<button class="text-gray-400 hover:text-gray-600 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
					</svg>
				</button>
				<span class="text-sm text-gray-500 font-medium">{link.clicks} clicks</span>
			</div>
		</div>

		<!-- Right Actions - Top Right -->
		<div class="flex items-center gap-2">
			<!-- Share Icon -->
			<button class="text-gray-400 hover:text-gray-600 transition-colors">
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
				</svg>
			</button>
			<!-- Toggle Switch -->
			<Switch checked={link.is_active} on:checkedChange={handleToggle} />
		</div>
	</div>
	
	<!-- Delete Icon - Absolute Position Bottom Right -->
	<button 
		onclick={handleDelete} 
		class="absolute bottom-4 right-4 text-gray-400 hover:text-red-600 transition-colors"
	>
		<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
		</svg>
	</button>
</div>
