<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { dndzone } from 'svelte-dnd-action';
	import type { Block } from '$lib/api/blocks';

	export let group: Block;
	export let expanded = false;
	
	const dispatch = createEventDispatcher();
	
	// Active tab
	let activeTab: 'texts' | 'layout' = 'texts';
	
	// Show input for adding text
	let showAddInput = false;
	
	// Edit mode for text items
	let editingTextId: string | null = null;
	let editingTextContent = '';
	
	// Group style from parent block
	$: groupStyle = (() => {
		try {
			if (!group.style) {
				return {
					textAlign: 'left',
					fontSize: 'text-medium',
					textColor: '#000000',
					isBold: false,
					isItalic: false
				};
			}
			if (typeof group.style === 'string') {
				return JSON.parse(group.style);
			}
			return group.style;
		} catch (e) {
			console.error('Failed to parse group.style:', group.style, e);
			return {
				textAlign: 'left',
				fontSize: 'text-medium',
				textColor: '#000000',
				isBold: false,
				isItalic: false
			};
		}
	})();
	
	let newTextContent = '';

	// Drag & drop items
	type DndItem = { id: string; data: Block };
	let dndItems: DndItem[] = [];
	
	// Sync dndItems from children
	$: {
		const children = group.children || [];
		if (dndItems.length === 0 || dndItems.length !== children.length) {
			dndItems = children
				.slice()
				.sort((a, b) => a.position - b.position)
				.map(block => ({ id: block.id, data: block }));
		} else {
			dndItems = dndItems.map(item => {
				const updated = children.find(c => c.id === item.id);
				return updated ? { ...item, data: updated } : item;
			});
		}
	}

	function handleDndConsider(e: CustomEvent) {
		dndItems = e.detail.items;
	}
	
	function handleDndFinalize(e: CustomEvent) {
		dndItems = e.detail.items;
		const blockIds = dndItems.map(item => item.id);
		dispatch('reordertexts', { groupId: group.id, blockIds });
	}

	function toggleExpand() {
		if (expanded) {
			dispatch('collapse', group.id);
		} else {
			dispatch('expand', group.id);
		}
	}

	function handleDelete() {
		if (confirm(`Delete "${group.group_title}" and all its text items?`)) {
			dispatch('delete', group.id);
		}
	}
	
	function handleToggleGroupVisibility() {
		dispatch('togglevisibility', group.id);
	}
	
	function handleToggleTextVisibility(textId: string) {
		dispatch('toggletextvisibility', { groupId: group.id, textId });
	}

	function handleAddText() {
		if (!newTextContent.trim()) return;
		dispatch('addtext', { groupId: group.id, content: newTextContent.trim() });
		newTextContent = '';
		showAddInput = false; // Hide input after adding
	}
	
	function toggleAddInput() {
		showAddInput = !showAddInput;
		if (showAddInput) {
			editingTextId = null; // Hide edit mode when showing add input
			editingTextContent = '';
		} else {
			newTextContent = ''; // Clear input when hiding
		}
	}

	function handleDeleteText(textId: string) {
		dispatch('deletetext', { groupId: group.id, textId });
	}
	
	function startEditText(textId: string, currentContent: string) {
		showAddInput = false; // Hide add input when editing
		editingTextId = textId;
		editingTextContent = currentContent;
	}
	
	function cancelEditText() {
		editingTextId = null;
		editingTextContent = '';
	}
	
	function saveEditText() {
		if (!editingTextContent.trim() || !editingTextId) return;
		dispatch('edittext', { 
			groupId: group.id, 
			textId: editingTextId, 
			content: editingTextContent.trim() 
		});
		editingTextId = null;
		editingTextContent = '';
	}
	
	function applyGroupStyle() {
		dispatch('updatestyle', { groupId: group.id, style: groupStyle });
	}
	
	// Presets
	const presets = {
		default: { textAlign: 'left', fontSize: 'text-medium', textColor: '#000000', isBold: false, isItalic: false },
		faq: { textAlign: 'left', fontSize: 'text-medium', textColor: '#1f2937', isBold: true, isItalic: false },
		quote: { textAlign: 'center', fontSize: 'text-large', textColor: '#4b5563', isBold: false, isItalic: true },
		feature: { textAlign: 'left', fontSize: 'text-medium', textColor: '#1f2937', isBold: true, isItalic: false }
	};
	
	function applyPreset(preset: keyof typeof presets) {
		groupStyle = { ...presets[preset] };
		applyGroupStyle();
	}
</script>

<div class="bg-white rounded-2xl shadow-sm hover:shadow-md transition-all border border-gray-100">
	<!-- Group Header -->
	<div class="flex items-center gap-3 p-4">
		<!-- Expand/Collapse Button -->
		<button
			onclick={toggleExpand}
			class="flex-shrink-0 w-8 h-8 rounded-lg hover:bg-gray-100 flex items-center justify-center transition-colors"
			aria-label={expanded ? 'Collapse group' : 'Expand group'}
		>
			<svg 
				class="w-5 h-5 text-gray-500 transition-transform {expanded ? 'rotate-90' : ''}" 
				fill="none" 
				stroke="currentColor" 
				viewBox="0 0 24 24"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
			</svg>
		</button>

		<!-- Icon -->
		<div class="flex-shrink-0 w-12 h-12 rounded-lg bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center shadow-sm">
			<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"/>
			</svg>
		</div>

		<!-- Title & Count -->
		<button onclick={toggleExpand} class="flex-1 text-left min-w-0">
			<h3 class="text-base font-semibold text-gray-900 truncate">
				{group.group_title || 'Text Group'}
			</h3>
			<p class="text-sm text-gray-500">
				{group.children?.length || 0} text {group.children?.length === 1 ? 'item' : 'items'}
			</p>
		</button>

		<!-- Actions -->
		<div class="flex items-center gap-1">
			<!-- Toggle Visibility -->
			<button
				onclick={handleToggleGroupVisibility}
				class="p-2 text-gray-400 hover:text-gray-600 transition-colors rounded-lg hover:bg-gray-100"
				title={group.is_active ? 'Hide group' : 'Show group'}
			>
				{#if group.is_active}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
					</svg>
				{:else}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
					</svg>
				{/if}
			</button>
			
			<!-- Drag Handle -->
			<button class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400 p-2 transition-colors" aria-label="Drag to reorder">
				<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
					<circle cx="9" cy="7" r="1.5"/>
					<circle cx="9" cy="12" r="1.5"/>
					<circle cx="9" cy="17" r="1.5"/>
					<circle cx="15" cy="7" r="1.5"/>
					<circle cx="15" cy="12" r="1.5"/>
					<circle cx="15" cy="17" r="1.5"/>
				</svg>
			</button>

			<!-- Delete -->
			<button
				onclick={handleDelete}
				class="p-2 text-gray-400 hover:text-red-500 transition-colors rounded-lg hover:bg-red-50"
				title="Delete group"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
				</svg>
			</button>
		</div>
	</div>

	<!-- Expanded Content -->
	{#if expanded}
		<div class="border-t border-gray-100">
			<!-- Tabs -->
			<div class="px-6 pt-6 pb-3 border-b border-gray-100">
				<div class="flex gap-6">
					<button
						onclick={() => activeTab = 'texts'}
						class="pb-3 text-sm font-semibold transition-colors relative"
						class:text-gray-900={activeTab === 'texts'}
						class:text-gray-500={activeTab !== 'texts'}
					>
						Texts
						{#if activeTab === 'texts'}
							<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-gray-900"></div>
						{/if}
					</button>
					<button
						onclick={() => activeTab = 'layout'}
						class="pb-3 text-sm font-semibold transition-colors relative"
						class:text-gray-900={activeTab === 'layout'}
						class:text-gray-500={activeTab !== 'layout'}
					>
						Layout
						{#if activeTab === 'layout'}
							<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-gray-900"></div>
						{/if}
					</button>
				</div>
			</div>

			<!-- Content -->
			<div class="p-6">
				{#if activeTab === 'texts'}
					<!-- Add Text Section -->
					<div class="mb-6">
						{#if !showAddInput}
							<!-- Add Text Button -->
							<button
								onclick={toggleAddInput}
								class="w-full py-3 px-4 bg-gradient-to-r from-purple-600 to-purple-700 hover:from-purple-700 hover:to-purple-800 text-white rounded-xl font-semibold transition-all shadow-md hover:shadow-lg flex items-center justify-center gap-2 group"
							>
								<svg class="w-5 h-5 transition-transform group-hover:rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
								</svg>
								Add Text
							</button>
						{:else}
							<!-- Add Text Input Form -->
							<div class="bg-gradient-to-br from-purple-50 via-white to-blue-50 rounded-2xl p-6 border-2 border-purple-200/60 shadow-sm">
								<div class="flex items-center gap-2 mb-3">
									<svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
									</svg>
									<span class="text-sm font-semibold text-gray-700">
										Add New Text Item
									</span>
								</div>
								<div class="flex gap-3">
									<div class="flex-1 relative group">
										<input 
											type="text" 
											bind:value={newTextContent}
											onkeydown={(e) => e.key === 'Enter' && handleAddText()}
											placeholder="Enter your text here..."
											class="w-full px-5 py-4 text-base bg-white border-2 border-gray-200 rounded-xl focus:outline-none focus:ring-4 focus:ring-purple-500/20 focus:border-purple-500 transition-all shadow-sm hover:shadow-md hover:border-gray-300 placeholder-gray-400"
											autofocus
										/>
										{#if newTextContent.trim()}
											<div class="absolute right-4 top-1/2 -translate-y-1/2 text-xs font-medium text-purple-600 bg-purple-100 px-2 py-1 rounded-full">
												Press Enter
											</div>
										{/if}
									</div>
									<button 
										onclick={handleAddText}
										disabled={!newTextContent.trim()}
										class="px-8 py-4 bg-gradient-to-r from-purple-600 to-purple-700 text-white rounded-xl text-base font-bold hover:from-purple-700 hover:to-purple-800 disabled:from-gray-300 disabled:to-gray-400 disabled:cursor-not-allowed transition-all shadow-lg shadow-purple-500/30 hover:shadow-xl hover:shadow-purple-500/40 disabled:shadow-none hover:scale-105 active:scale-95 flex items-center gap-2 group"
									>
										<svg class="w-5 h-5 transition-transform group-hover:rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
										</svg>
										Add
									</button>
								</div>
								<div class="flex items-center justify-between mt-3">
									<p class="text-xs text-gray-500 flex items-center gap-2">
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
										</svg>
										Press <kbd class="px-1.5 py-0.5 bg-gray-100 border border-gray-300 rounded text-xs font-mono">Enter</kbd> to add
									</p>
									<button
										onclick={toggleAddInput}
										class="text-sm text-gray-500 hover:text-gray-700 font-medium"
									>
										Cancel
									</button>
								</div>
							</div>
						{/if}
					</div>

					<!-- Text List -->
					{#if group.children && group.children.length > 0}
						<div 
							class="space-y-3"
							use:dndzone={{items: dndItems, flipDurationMs: 200, dropTargetStyle: {}}}
							onconsider={handleDndConsider}
							onfinalize={handleDndFinalize}
						>
							{#each dndItems as item (item.id)}
								{@const textBlock = item.data}
								{#if editingTextId === textBlock.id}
									<!-- Edit Mode -->
									<div class="bg-gradient-to-br from-purple-50 via-white to-blue-50 rounded-2xl p-6 border-2 border-purple-200/60 shadow-sm">
										<div class="flex items-center gap-2 mb-3">
											<svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
											</svg>
											<span class="text-sm font-semibold text-gray-700">
												Edit Text
											</span>
										</div>
										<div class="flex gap-3">
											<div class="flex-1 relative group">
												<input
													type="text"
													bind:value={editingTextContent}
													onkeydown={(e) => {
														if (e.key === 'Enter') saveEditText();
														if (e.key === 'Escape') cancelEditText();
													}}
													placeholder="Enter your text here..."
													class="w-full px-5 py-4 text-base bg-white border-2 border-gray-200 rounded-xl focus:outline-none focus:ring-4 focus:ring-purple-500/20 focus:border-purple-500 transition-all shadow-sm hover:shadow-md hover:border-gray-300 placeholder-gray-400"
													autofocus
												/>
												{#if editingTextContent.trim()}
													<div class="absolute right-4 top-1/2 -translate-y-1/2 text-xs font-medium text-purple-600 bg-purple-100 px-2 py-1 rounded-full">
														Press Enter
													</div>
												{/if}
											</div>
											<button
												onclick={saveEditText}
												disabled={!editingTextContent.trim()}
												class="px-8 py-4 bg-gradient-to-r from-purple-600 to-purple-700 text-white rounded-xl text-base font-bold hover:from-purple-700 hover:to-purple-800 disabled:from-gray-300 disabled:to-gray-400 disabled:cursor-not-allowed transition-all shadow-lg shadow-purple-500/30 hover:shadow-xl hover:shadow-purple-500/40 disabled:shadow-none hover:scale-105 active:scale-95 flex items-center gap-2 group"
											>
												<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
												</svg>
												Save
											</button>
										</div>
										<div class="flex items-center justify-between mt-3">
											<p class="text-xs text-gray-500 flex items-center gap-2">
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
												</svg>
												Press <kbd class="px-1.5 py-0.5 bg-gray-100 border border-gray-300 rounded text-xs font-mono">Enter</kbd> to save
											</p>
											<button
												onclick={cancelEditText}
												class="text-sm text-gray-500 hover:text-gray-700 font-medium"
											>
												Cancel
											</button>
										</div>
									</div>
								{:else}
									<!-- View Mode -->
									<div class="bg-gradient-to-r from-gray-50 to-gray-100/50 rounded-2xl shadow-sm hover:shadow-md transition-all border border-gray-200/60">
										<div class="flex items-center gap-4 p-3.5">
											<!-- Drag Handle -->
											<button
												class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400 p-1 transition-colors flex-shrink-0"
												onclick={(e) => e.stopPropagation()}
												aria-label="Drag to reorder text"
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

											<!-- Icon -->
											<div class="flex-shrink-0 w-10 h-10 rounded-lg bg-gradient-to-br from-purple-100 to-purple-200 flex items-center justify-center">
												<svg class="w-4 h-4 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"/>
												</svg>
											</div>

											<!-- Text Content (without style preview in card) -->
											<div class="flex-1 min-w-0">
												<p class="text-sm text-gray-900 break-words">
													{textBlock.content || 'Empty text'}
												</p>
											</div>

											<!-- Actions -->
											<div class="flex items-center gap-1">
												<!-- Toggle Visibility -->
												<button
													onclick={() => handleToggleTextVisibility(textBlock.id)}
													class="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors flex-shrink-0"
													title={textBlock.is_active ? 'Hide text' : 'Show text'}
												>
													{#if textBlock.is_active}
														<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
															<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
														</svg>
													{:else}
														<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
														</svg>
													{/if}
												</button>
												
												<!-- Edit Button -->
												<button
													onclick={() => startEditText(textBlock.id, textBlock.content || '')}
													class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors flex-shrink-0"
													title="Edit"
												>
													<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
													</svg>
												</button>
												
												<!-- Delete Button -->
												<button
													onclick={() => handleDeleteText(textBlock.id)}
													class="p-1.5 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors flex-shrink-0"
													title="Delete"
												>
													<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
													</svg>
												</button>
											</div>
										</div>
									</div>
								{/if}
							{/each}
						</div>
					{/if}

				{:else if activeTab === 'layout'}
					<!-- Layout Tab -->
					<div class="space-y-4">
						<!-- Presets -->
						<div>
							<label class="text-sm font-medium text-gray-700 mb-3 block">Quick Presets</label>
							<div class="flex gap-2 flex-wrap">
								<button onclick={() => applyPreset('default')} class="px-4 py-2 text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors">
									Default
								</button>
								<button onclick={() => applyPreset('faq')} class="px-4 py-2 text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors">
									FAQ
								</button>
								<button onclick={() => applyPreset('quote')} class="px-4 py-2 text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors">
									Quote
								</button>
								<button onclick={() => applyPreset('feature')} class="px-4 py-2 text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors">
									Feature
								</button>
							</div>
						</div>
						
						<!-- Text Alignment -->
						<div>
							<label class="text-sm font-medium text-gray-700 mb-3 block">Text Alignment</label>
							<div class="flex gap-2">
								<button 
									onclick={() => { groupStyle.textAlign = 'left'; applyGroupStyle(); }} 
									class="flex-1 p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textAlign === 'left' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									<svg class="w-5 h-5 mx-auto text-gray-700" fill="currentColor" viewBox="0 0 24 24">
										<path d="M3 4h18v2H3V4zm0 7h12v2H3v-2zm0 7h18v2H3v-2z"/>
									</svg>
								</button>
								<button 
									onclick={() => { groupStyle.textAlign = 'center'; applyGroupStyle(); }} 
									class="flex-1 p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textAlign === 'center' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									<svg class="w-5 h-5 mx-auto text-gray-700" fill="currentColor" viewBox="0 0 24 24">
										<path d="M3 4h18v2H3V4zm4 7h10v2H7v-2zm-4 7h18v2H3v-2z"/>
									</svg>
								</button>
								<button 
									onclick={() => { groupStyle.textAlign = 'right'; applyGroupStyle(); }} 
									class="flex-1 p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textAlign === 'right' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									<svg class="w-5 h-5 mx-auto text-gray-700" fill="currentColor" viewBox="0 0 24 24">
										<path d="M3 4h18v2H3V4zm6 7h12v2H9v-2zm-6 7h18v2H3v-2z"/>
									</svg>
								</button>
							</div>
						</div>
						
						<!-- Font Size -->
						<div>
							<label class="text-sm font-medium text-gray-700 mb-3 block">Font Size</label>
							<select 
								bind:value={groupStyle.fontSize} 
								onchange={applyGroupStyle} 
								class="w-full p-3 bg-white border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
							>
								<option value="text-small">Small</option>
								<option value="text-medium">Medium</option>
								<option value="text-large">Large</option>
								<option value="headline-small">Headline Small</option>
								<option value="headline-medium">Headline Medium</option>
								<option value="headline-large">Headline Large</option>
							</select>
						</div>
						
						<!-- Text Style -->
						<div>
							<label class="text-sm font-medium text-gray-700 mb-3 block">Text Style</label>
							<div class="flex gap-2">
								<button 
									onclick={() => { groupStyle.isBold = !groupStyle.isBold; applyGroupStyle(); }} 
									class="px-4 py-2 font-bold text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.isBold ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									B
								</button>
								<button 
									onclick={() => { groupStyle.isItalic = !groupStyle.isItalic; applyGroupStyle(); }} 
									class="px-4 py-2 italic text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.isItalic ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									I
								</button>
							</div>
						</div>
						
						<!-- Text Color -->
						<div>
							<label class="text-sm font-medium text-gray-700 mb-3 block">Text Color</label>
							<div class="flex gap-3 items-center">
								<input 
									type="color" 
									bind:value={groupStyle.textColor} 
									onchange={applyGroupStyle} 
									class="w-16 h-10 rounded-lg border cursor-pointer" 
								/>
								<span class="text-sm text-gray-600 font-mono">{groupStyle.textColor}</span>
							</div>
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
