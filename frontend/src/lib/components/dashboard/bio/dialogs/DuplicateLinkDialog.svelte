<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { Link } from '$lib/api/links';
	
	export let open = false;
	export let link: Link | null = null;
	export let groups: Link[] = []; // All available groups
	export let currentGroupId: string = '';
	
	const dispatch = createEventDispatcher();
	
	let duplicateMode: 'current' | 'other' = 'current';
	let selectedGroupId: string = '';
	
	$: if (open && link) {
		// Reset to default when opening
		duplicateMode = 'current';
		selectedGroupId = '';
	}
	
	$: availableGroups = groups.filter(g => g.id !== currentGroupId);
	
	function handleClose() {
		open = false;
	}
	
	function handleDuplicate() {
		if (!link) return;
		
		if (duplicateMode === 'current') {
			// Duplicate in current group
			dispatch('duplicate', { linkId: link.id, targetGroupId: currentGroupId });
		} else {
			// Duplicate to another group
			if (!selectedGroupId) return;
			dispatch('duplicate', { linkId: link.id, targetGroupId: selectedGroupId });
		}
		
		open = false;
	}
	
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		} else if (event.key === 'Enter' && (event.ctrlKey || event.metaKey)) {
			handleDuplicate();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-md p-0 overflow-hidden" showCloseButton={false}>
		<!-- Header -->
		<div class="flex items-center justify-between px-6 py-4 border-b border-gray-200 bg-gradient-to-r from-blue-50 to-indigo-50">
			<div class="flex items-center gap-3">
				<div class="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-sm">
					<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
					</svg>
				</div>
				<h2 class="text-lg font-bold text-gray-900">Duplicate Link</h2>
			</div>
			<button 
				onclick={handleClose}
				class="text-gray-500 hover:text-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-400 rounded-lg p-1"
				aria-label="Close dialog"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
		
		<!-- Content -->
		<div class="p-6 space-y-6">
			<!-- Link Preview -->
			{#if link}
				<div class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg border border-gray-200">
					{#if link.thumbnail_url}
						<img src={link.thumbnail_url} alt={link.title} class="w-12 h-12 rounded-lg object-cover flex-shrink-0" />
					{:else}
						<div class="w-12 h-12 rounded-lg bg-gradient-to-br from-gray-300 to-gray-400 flex items-center justify-center flex-shrink-0">
							<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
							</svg>
						</div>
					{/if}
					<div class="flex-1 min-w-0">
						<p class="font-semibold text-gray-900 truncate">{link.title}</p>
						<p class="text-sm text-gray-500 truncate">{link.url}</p>
					</div>
				</div>
			{/if}

			<!-- Options -->
			<div class="space-y-3">
				<label class="block">
					<input
						type="radio"
						name="duplicate-mode"
						value="current"
						bind:group={duplicateMode}
						class="sr-only peer"
					/>
					<div class="flex items-start gap-3 p-4 rounded-xl border-2 border-gray-200 cursor-pointer transition-all hover:border-indigo-300 hover:bg-indigo-50/50 peer-checked:border-indigo-500 peer-checked:bg-indigo-50 peer-checked:shadow-sm">
						<div class="flex-shrink-0 mt-0.5">
							<div class="w-5 h-5 rounded-full border-2 border-gray-300 peer-checked:border-indigo-500 peer-checked:bg-indigo-500 flex items-center justify-center transition-all">
								{#if duplicateMode === 'current'}
									<div class="w-2 h-2 rounded-full bg-white"></div>
								{/if}
							</div>
						</div>
						<div class="flex-1">
							<div class="flex items-center gap-2 mb-1">
								<svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
								</svg>
								<span class="font-semibold text-gray-900">Duplicate in current group</span>
							</div>
							<p class="text-sm text-gray-600">Create a copy in the same group</p>
						</div>
					</div>
				</label>

				<label class="block">
					<input
						type="radio"
						name="duplicate-mode"
						value="other"
						bind:group={duplicateMode}
						class="sr-only peer"
					/>
					<div class="flex items-start gap-3 p-4 rounded-xl border-2 border-gray-200 cursor-pointer transition-all hover:border-indigo-300 hover:bg-indigo-50/50 peer-checked:border-indigo-500 peer-checked:bg-indigo-50 peer-checked:shadow-sm">
						<div class="flex-shrink-0 mt-0.5">
							<div class="w-5 h-5 rounded-full border-2 border-gray-300 peer-checked:border-indigo-500 peer-checked:bg-indigo-500 flex items-center justify-center transition-all">
								{#if duplicateMode === 'other'}
									<div class="w-2 h-2 rounded-full bg-white"></div>
								{/if}
							</div>
						</div>
						<div class="flex-1">
							<div class="flex items-center gap-2 mb-1">
								<svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
								</svg>
								<span class="font-semibold text-gray-900">Duplicate to another group</span>
							</div>
							<p class="text-sm text-gray-600">Create a copy in a different group</p>
						</div>
					</div>
				</label>
			</div>

			<!-- Group Selector (shown when "other" is selected) -->
			{#if duplicateMode === 'other'}
				<div class="animate-fadeIn">
					<label for="target-group" class="block text-sm font-medium text-gray-700 mb-2">
						Select target group
					</label>
					{#if availableGroups.length > 0}
						<select
							id="target-group"
							bind:value={selectedGroupId}
							class="w-full px-4 py-3 bg-white border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 focus:outline-none transition-all text-gray-900"
						>
							<option value="" disabled>Choose a group...</option>
							{#each availableGroups as group}
								<option value={group.id}>
									{group.group_title || 'Untitled Group'} ({group.children?.length || 0} links)
								</option>
							{/each}
						</select>
					{:else}
						<div class="p-4 bg-amber-50 border border-amber-200 rounded-lg">
							<div class="flex items-start gap-3">
								<svg class="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
								</svg>
								<div>
									<p class="font-medium text-amber-900">No other groups available</p>
									<p class="text-sm text-amber-700 mt-1">Create more groups to duplicate links between them.</p>
								</div>
							</div>
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="flex items-center justify-end gap-3 px-6 py-4 bg-gray-50 border-t border-gray-200">
			<button
				onclick={handleClose}
				class="px-5 py-2.5 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-gray-300"
			>
				Cancel
			</button>
			<button
				onclick={handleDuplicate}
				disabled={duplicateMode === 'other' && !selectedGroupId}
				class="px-5 py-2.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white text-sm font-semibold rounded-lg shadow-sm hover:shadow-md transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
			>
				<span class="flex items-center gap-2">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
					</svg>
					Duplicate
				</span>
			</button>
		</div>
	</Dialog.Content>
</Dialog.Root>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
	
	.animate-fadeIn {
		animation: fadeIn 0.2s ease-out;
	}
	
	/* Custom radio button styling */
	input[type="radio"]:checked + div .w-5 {
		border-color: #6366f1;
		background-color: #6366f1;
	}
</style>
