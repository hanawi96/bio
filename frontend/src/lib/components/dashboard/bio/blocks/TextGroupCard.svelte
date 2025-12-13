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
	
	// Menu state
	let groupMenuOpen = false;
	let groupMenuButtonRef: HTMLElement | null = null;
	let shouldGroupMenuOpenUpward = false;
	
	// Title editing
	let isEditingTitle = false;
	let editedTitle = '';
	let titleInputRef: HTMLInputElement | null = null;
	let isSavingTitle = false;
	
	// Group style from parent block - using let instead of $: for mutability
	let groupStyle = {
		textAlign: 'left',
		fontSize: 'text-medium',
		textColor: '#000000',
		isBold: false,
		isItalic: false,
		isUnderline: false,
		isStrikethrough: false,
		textTransform: 'none',
		hasBackground: false,
		backgroundColor: '#ffffff',
		backgroundOpacity: 90,
		borderRadius: 12,
		padding: 16,
		shadow: 'none',
		hasBorder: false,
		borderColor: '#e5e7eb',
		borderWidth: 1,
		borderStyle: 'solid'
	};
	
	// Parse and update groupStyle whenever group.style changes
	$: {
		try {
			if (!group.style) {
				groupStyle = {
					textAlign: 'left',
					fontSize: 'text-medium',
					textColor: '#000000',
					isBold: false,
					isItalic: false,
					isUnderline: false,
					isStrikethrough: false,
					textTransform: 'none',
					hasBackground: false,
					backgroundColor: '#ffffff',
					backgroundOpacity: 90,
					borderRadius: 12,
					padding: 16,
					shadow: 'none',
					hasBorder: false,
					borderColor: '#e5e7eb',
					borderWidth: 1,
					borderStyle: 'solid'
				};
			} else if (typeof group.style === 'string') {
				const parsed = JSON.parse(group.style);
				groupStyle = {
					textAlign: parsed.textAlign || 'left',
					fontSize: parsed.fontSize || 'text-medium',
					textColor: parsed.textColor || '#000000',
					isBold: parsed.isBold || false,
					isItalic: parsed.isItalic || false,
					isUnderline: parsed.isUnderline || false,
					isStrikethrough: parsed.isStrikethrough || false,
					textTransform: parsed.textTransform || 'none',
					hasBackground: parsed.hasBackground || false,
					backgroundColor: parsed.backgroundColor || '#ffffff',
					backgroundOpacity: parsed.backgroundOpacity !== undefined ? parsed.backgroundOpacity : 90,
					borderRadius: parsed.borderRadius !== undefined ? parsed.borderRadius : 12,
					padding: parsed.padding || 16,
					shadow: parsed.shadow || 'none',
					hasBorder: parsed.hasBorder || false,
					borderColor: parsed.borderColor || '#e5e7eb',
					borderWidth: parsed.borderWidth || 1,
					borderStyle: parsed.borderStyle || 'solid'
				};
			} else {
				const style = group.style as any;
				groupStyle = {
					textAlign: style.textAlign || 'left',
					fontSize: style.fontSize || 'text-medium',
					textColor: style.textColor || '#000000',
					isBold: style.isBold || false,
					isItalic: style.isItalic || false,
					isUnderline: style.isUnderline || false,
					isStrikethrough: style.isStrikethrough || false,
					textTransform: style.textTransform || 'none',
					hasBackground: style.hasBackground || false,
					backgroundColor: style.backgroundColor || '#ffffff',
					backgroundOpacity: style.backgroundOpacity !== undefined ? style.backgroundOpacity : 90,
					borderRadius: style.borderRadius !== undefined ? style.borderRadius : 12,
					padding: style.padding || 16,
					shadow: style.shadow || 'none',
					hasBorder: style.hasBorder || false,
					borderColor: style.borderColor || '#e5e7eb',
					borderWidth: style.borderWidth || 1,
					borderStyle: style.borderStyle || 'solid'
				};
			}
		} catch (e) {
			console.error('Failed to parse group.style:', group.style, e);
			groupStyle = {
				textAlign: 'left',
				fontSize: 'text-medium',
				textColor: '#000000',
				isBold: false,
				isItalic: false,
				isUnderline: false,
				isStrikethrough: false,
				textTransform: 'none',
				hasBackground: false,
				backgroundColor: '#ffffff',
				backgroundOpacity: 90,
				borderRadius: 12,
				padding: 16,
				shadow: 'none',
				hasBorder: false,
				borderColor: '#e5e7eb',
				borderWidth: 1,
				borderStyle: 'solid'
			};
		}
	}
	
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
		dispatch('delete', group.id);
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
	
	// Debounce timer for slider API calls
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;
	
	// For sliders - debounced to avoid too many API calls
	function applyGroupStyleDebounced() {
		if (debounceTimer) clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			dispatch('updatestyle', { groupId: group.id, style: groupStyle });
		}, 300);
	}
	
	// For buttons/toggles - immediate, no debounce
	function applyGroupStyle() {
		if (debounceTimer) clearTimeout(debounceTimer);
		dispatch('updatestyle', { groupId: group.id, style: groupStyle });
	}
	
	// Modern Creative Presets - Professional Color Combinations
	const presets = {
		default: { 
			textAlign: 'left', fontSize: 'text-medium', textColor: '#000000', isBold: false, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: false, backgroundColor: '#ffffff', backgroundOpacity: 90, borderRadius: 12, padding: 16, shadow: 'none', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Gradient-inspired with depth
		'ocean-depth': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#06B6D4', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Warm sunset with complementary
		'sunset-blaze': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#F97316', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Purple luxury with gold accent
		'royal-purple': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#FDE047', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#7C3AED', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Mint with deep teal
		'mint-teal': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#34D399', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Hot pink with purple
		'pink-passion': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#EC4899', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Navy with cyan accent
		'navy-cyan': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#67E8F9', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#1E3A8A', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Coral with deep orange
		'coral-dream': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#FB923C', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Electric blue with pink
		'electric-pop': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#3B82F6', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Lime with emerald
		'lime-fresh': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#14532D', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#BEF264', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Deep red with dark border
		'crimson-bold': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#DC2626', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Indigo with violet
		'indigo-wave': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#E0E7FF', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#6366F1', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Teal with sky blue
		'teal-sky': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#ffffff', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#14B8A6', backgroundOpacity: 100, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Frosted glass premium
		'glass-elite': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#1E293B', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#F8FAFC', backgroundOpacity: 60, borderRadius: 24, padding: 24, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Neon cyberpunk
		'neon-cyber': { 
			textAlign: 'center', fontSize: 'text-medium', textColor: '#00FFF0', isBold: true, isItalic: false, isUnderline: false, isStrikethrough: false, textTransform: 'uppercase',
			hasBackground: true, backgroundColor: '#0F172A', backgroundOpacity: 100, borderRadius: 20, padding: 22, shadow: 'xl', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		},
		// Elegant quote
		'quote-elegant': { 
			textAlign: 'center', fontSize: 'text-large', textColor: '#475569', isBold: false, isItalic: true, isUnderline: false, isStrikethrough: false, textTransform: 'none',
			hasBackground: true, backgroundColor: '#F1F5F9', backgroundOpacity: 100, borderRadius: 20, padding: 24, shadow: 'md', hasBorder: false, borderColor: '#e5e7eb', borderWidth: 1, borderStyle: 'solid'
		}
	};
	
	const presetCategories = [
		{ name: 'Gradient', presets: ['ocean-depth', 'sunset-blaze', 'mint-teal', 'coral-dream', 'indigo-wave'] },
		{ name: 'Luxury', presets: ['royal-purple', 'navy-cyan', 'glass-elite', 'teal-sky', 'quote-elegant'] },
		{ name: 'Bold', presets: ['pink-passion', 'electric-pop', 'lime-fresh', 'crimson-bold', 'neon-cyber'] }
	];
	
	const presetLabels: Record<string, string> = {
		'default': 'Default',
		'ocean-depth': 'Ocean',
		'sunset-blaze': 'Sunset',
		'royal-purple': 'Royal',
		'mint-teal': 'Mint',
		'pink-passion': 'Pink',
		'navy-cyan': 'Navy',
		'coral-dream': 'Coral',
		'electric-pop': 'Electric',
		'lime-fresh': 'Lime',
		'crimson-bold': 'Crimson',
		'indigo-wave': 'Indigo',
		'teal-sky': 'Teal',
		'glass-elite': 'Glass',
		'neon-cyber': 'Neon',
		'quote-elegant': 'Quote'
	};
	
	function applyPreset(preset: keyof typeof presets) {
		groupStyle = { ...presets[preset] };
		applyGroupStyle();
	}
	
	// Helper function to convert hex to rgba
	function hexToRgba(hex: string, opacity: number) {
		const r = parseInt(hex.slice(1, 3), 16);
		const g = parseInt(hex.slice(3, 5), 16);
		const b = parseInt(hex.slice(5, 7), 16);
		return `rgba(${r}, ${g}, ${b}, ${opacity / 100})`;
	}
	
	// Click outside handler for group menu
	function clickOutsideGroupMenu(node: HTMLElement) {
		const handleClick = (event: MouseEvent) => {
			const target = event.target as Node;
			if (node && !node.contains(target) && !event.defaultPrevented) {
				groupMenuOpen = false;
			}
		};
		
		document.addEventListener('mousedown', handleClick, true);
		
		return {
			destroy() {
				document.removeEventListener('mousedown', handleClick, true);
			}
		};
	}
	
	function handleDuplicate() {
		dispatch('duplicate', group.id);
		groupMenuOpen = false;
	}
	
	function startEditTitle() {
		isEditingTitle = true;
		editedTitle = group.group_title || '';
		setTimeout(() => titleInputRef?.focus(), 0);
	}
	
	function saveTitle() {
		if (isSavingTitle) return; // Prevent double save
		isSavingTitle = true;
		
		if (!editedTitle.trim()) {
			editedTitle = group.group_title || 'Text Group';
		}
		if (editedTitle.trim() !== group.group_title) {
			dispatch('updatetitle', { groupId: group.id, title: editedTitle.trim() });
		}
		isEditingTitle = false;
		
		// Reset saving flag after a short delay
		setTimeout(() => {
			isSavingTitle = false;
		}, 100);
	}
	
	function cancelEditTitle() {
		if (isSavingTitle) return; // Don't cancel if saving
		isEditingTitle = false;
		editedTitle = '';
	}
	
	function handleTitleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			e.preventDefault();
			titleInputRef?.blur(); // Blur first to prevent onblur trigger
			saveTitle();
		} else if (e.key === 'Escape') {
			e.preventDefault();
			titleInputRef?.blur();
			cancelEditTitle();
		}
	}
	
	function handleTitleBlur() {
		// Only save on blur if not already saving (from Enter key)
		if (!isSavingTitle) {
			saveTitle();
		}
	}
</script>

<div class="bg-white rounded-2xl shadow-sm hover:shadow-md transition-all border border-gray-100">
	<!-- Group Header -->
	<div class="flex items-center gap-4 px-4 py-3.5">
		<!-- Drag Handle -->
		<button class="drag-handle cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-400 p-1 transition-colors flex-shrink-0" aria-label="Drag to reorder">
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
		<div class="flex-shrink-0 w-14 h-14 rounded-lg bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center shadow-sm">
			<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"/>
			</svg>
		</div>

		<!-- Title & Count -->
		<div onclick={expanded ? undefined : toggleExpand} class="flex-1 min-w-0" class:cursor-pointer={!expanded}>
			{#if expanded && isEditingTitle}
				<input
					bind:this={titleInputRef}
					bind:value={editedTitle}
					onkeydown={handleTitleKeydown}
					onblur={handleTitleBlur}
					class="text-base font-semibold text-gray-900 border-b-2 border-purple-500 focus:outline-none bg-transparent w-full"
					placeholder="Group title"
				/>
			{:else}
				<button 
					onclick={expanded ? startEditTitle : toggleExpand} 
					class="text-left w-full"
					class:hover:text-purple-600={expanded}
					class:transition-colors={expanded}
				>
					<h3 class="text-base font-semibold text-gray-900 truncate">
						{group.group_title || 'Text Group'}
					</h3>
				</button>
			{/if}
			<p class="text-sm text-gray-500">
				{group.children?.length || 0} text {group.children?.length === 1 ? 'item' : 'items'}
			</p>
		</div>

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

			<!-- 3-dot Menu Button -->
			<div class="relative" use:clickOutsideGroupMenu>
				<button
					bind:this={groupMenuButtonRef}
					onclick={(e) => { 
						e.stopPropagation();
						
						// Detect if menu should open upward
						if (groupMenuButtonRef) {
							const rect = groupMenuButtonRef.getBoundingClientRect();
							const windowHeight = window.innerHeight;
							const spaceBelow = windowHeight - rect.bottom;
							const menuHeight = 200; // Approximate menu height
							
							// If not enough space below, open upward
							shouldGroupMenuOpenUpward = spaceBelow < menuHeight + 60;
						}
						
						groupMenuOpen = !groupMenuOpen;
					}}
					class="p-2 text-gray-400 hover:text-gray-600 transition-colors rounded-lg hover:bg-gray-100"
					title="More options"
				>
					<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
						<circle cx="10" cy="4" r="1.5"/>
						<circle cx="10" cy="10" r="1.5"/>
						<circle cx="10" cy="16" r="1.5"/>
					</svg>
				</button>

				{#if groupMenuOpen}
					<div class="absolute right-0 {shouldGroupMenuOpenUpward ? 'bottom-full mb-1' : 'top-full mt-1'} w-56 bg-white rounded-xl shadow-xl border border-gray-200 py-2 z-20">
						<!-- Edit Group -->
						<button
							onclick={(e) => { 
								e.stopPropagation();
								toggleExpand();
								groupMenuOpen = false;
							}}
							class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-blue-50 flex items-center gap-3 transition-all group"
						>
							<svg class="w-4 h-4 text-blue-500 group-hover:text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
							</svg>
							<span class="font-medium">Edit group</span>
						</button>

						<!-- Duplicate Group -->
						<button
							onclick={(e) => { 
								e.stopPropagation();
								handleDuplicate();
							}}
							class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-green-50 flex items-center gap-3 transition-all group"
						>
							<svg class="w-4 h-4 text-green-500 group-hover:text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
							</svg>
							<span class="font-medium">Duplicate group</span>
						</button>

						<div class="border-t border-gray-100 my-1"></div>

						<!-- Delete Group -->
						<button
							onclick={(e) => { 
								e.stopPropagation();
								handleDelete();
								groupMenuOpen = false;
							}}
							class="w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50 flex items-center gap-3 transition-all group"
						>
							<svg class="w-4 h-4 text-red-500 group-hover:text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
							</svg>
							<span class="font-medium">Delete group</span>
						</button>
					</div>
				{/if}
			</div>
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
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"/>
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
					<div class="space-y-5">
						<!-- Presets Section -->
						<div>
							<div class="text-sm font-medium text-gray-700 mb-4">Quick Style Presets</div>
							<div class="space-y-4">
								{#each presetCategories as category}
									<div>
										<h4 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">{category.name}</h4>
										<div class="grid grid-cols-5 gap-2">
											{#each category.presets as presetKey}
												{@const preset = presets[presetKey as keyof typeof presets]}
												<button
													onclick={() => applyPreset(presetKey as keyof typeof presets)}
													class="p-3 text-center rounded-xl transition-colors"
													style="
														background-color: {preset.hasBackground ? hexToRgba(preset.backgroundColor, preset.backgroundOpacity) : 'transparent'};
													"
													class:shadow-sm={preset.shadow === 'sm'}
													class:shadow-md={preset.shadow === 'md'}
													class:shadow-lg={preset.shadow === 'lg'}
													class:shadow-xl={preset.shadow === 'xl'}
													title="Apply {presetLabels[presetKey]} style"
												>
													<div 
														class="text-xs font-medium truncate"
														class:font-bold={preset.isBold}
														class:italic={preset.isItalic}
														class:underline={preset.isUnderline}
														class:uppercase={preset.textTransform === 'uppercase'}
														style="color: {preset.textColor}; text-align: {preset.textAlign};"
													>
														{presetLabels[presetKey]}
													</div>
												</button>
											{/each}
										</div>
									</div>
								{/each}
							</div>
						</div>
						
						<!-- Divider -->
						<div class="border-t border-gray-200"></div>
						<!-- Text Alignment -->
						<div>
							<label class="text-sm font-medium text-gray-700 mb-3 block">Text Alignment</label>
							<div class="grid grid-cols-4 gap-2">
								<button 
									onclick={() => { groupStyle.textAlign = 'left'; applyGroupStyle(); }} 
									class="p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textAlign === 'left' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									<svg class="w-5 h-5 mx-auto text-gray-700" fill="currentColor" viewBox="0 0 24 24">
										<path d="M3 4h18v2H3V4zm0 7h12v2H3v-2zm0 7h18v2H3v-2z"/>
									</svg>
								</button>
								<button 
									onclick={() => { groupStyle.textAlign = 'center'; applyGroupStyle(); }} 
									class="p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textAlign === 'center' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									<svg class="w-5 h-5 mx-auto text-gray-700" fill="currentColor" viewBox="0 0 24 24">
										<path d="M3 4h18v2H3V4zm4 7h10v2H7v-2zm-4 7h18v2H3v-2z"/>
									</svg>
								</button>
								<button 
									onclick={() => { groupStyle.textAlign = 'right'; applyGroupStyle(); }} 
									class="p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textAlign === 'right' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
								>
									<svg class="w-5 h-5 mx-auto text-gray-700" fill="currentColor" viewBox="0 0 24 24">
										<path d="M3 4h18v2H3V4zm6 7h12v2H9v-2zm-6 7h18v2H3v-2z"/>
									</svg>
								</button>
								<button 
									onclick={() => { groupStyle.textAlign = 'justify'; applyGroupStyle(); }} 
									class="p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textAlign === 'justify' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
									title="Justify"
								>
									<svg class="w-5 h-5 mx-auto text-gray-700" fill="currentColor" viewBox="0 0 24 24">
										<path d="M3 4h18v2H3V4zm0 7h18v2H3v-2zm0 7h18v2H3v-2z"/>
									</svg>
								</button>
							</div>
						</div>
						
					<!-- Font Size -->
					<div>
						<label class="text-sm font-medium text-gray-700 mb-3 block">Font Size</label>
						<div class="grid grid-cols-3 gap-2">
							<button
								onclick={() => { groupStyle.fontSize = 'text-small'; applyGroupStyle(); }}
								class="px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.fontSize === 'text-small' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
							>
								Small
							</button>
							<button
								onclick={() => { groupStyle.fontSize = 'text-medium'; applyGroupStyle(); }}
								class="px-3 py-2 text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.fontSize === 'text-medium' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
							>
								Medium
							</button>
							<button
								onclick={() => { groupStyle.fontSize = 'text-large'; applyGroupStyle(); }}
								class="px-3 py-2 text-base bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.fontSize === 'text-large' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
							>
								Large
							</button>
							<button
								onclick={() => { groupStyle.fontSize = 'headline-small'; applyGroupStyle(); }}
								class="px-3 py-2 text-sm font-bold bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.fontSize === 'headline-small' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
							>
								H3
							</button>
							<button
								onclick={() => { groupStyle.fontSize = 'headline-medium'; applyGroupStyle(); }}
								class="px-3 py-2 text-base font-bold bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.fontSize === 'headline-medium' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
							>
								H2
							</button>
							<button
								onclick={() => { groupStyle.fontSize = 'headline-large'; applyGroupStyle(); }}
								class="px-3 py-2 text-lg font-bold bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.fontSize === 'headline-large' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
							>
								H1
							</button>
						</div>
					</div>						<!-- Text Style -->
						<div>
							<div class="flex items-center justify-between mb-3">
								<label class="text-sm font-medium text-gray-700">Text Style</label>
								<button
									onclick={() => { 
										groupStyle.isBold = false;
										groupStyle.isItalic = false;
										groupStyle.isUnderline = false;
										groupStyle.isStrikethrough = false;
										groupStyle.textTransform = 'none';
										applyGroupStyle();
									}}
									class="px-2 py-1 text-xs text-gray-600 hover:text-purple-600 hover:bg-purple-50 rounded transition-colors"
									title="Reset to default"
								>
									Reset
								</button>
							</div>
							<div class="flex gap-2 flex-wrap">
								<button
									onclick={() => { 
										groupStyle.isBold = false;
										groupStyle.isItalic = false;
										groupStyle.isUnderline = false;
										groupStyle.isStrikethrough = false;
										groupStyle.textTransform = 'none';
										applyGroupStyle();
									}}
									class="px-4 py-2 text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {!groupStyle.isBold && !groupStyle.isItalic && !groupStyle.isUnderline && !groupStyle.isStrikethrough && groupStyle.textTransform === 'none' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
									title="Normal text"
								>
									Normal
								</button>
								<button 
									onclick={() => { groupStyle.isBold = !groupStyle.isBold; applyGroupStyle(); }} 
									class="px-4 py-2 font-bold text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.isBold ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
									title="Bold"
								>
									B
								</button>
								<button 
									onclick={() => { groupStyle.isItalic = !groupStyle.isItalic; applyGroupStyle(); }} 
									class="px-4 py-2 italic text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.isItalic ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
									title="Italic"
								>
									I
								</button>
								<button 
									onclick={() => { groupStyle.isUnderline = !groupStyle.isUnderline; applyGroupStyle(); }} 
									class="px-4 py-2 underline text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.isUnderline ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
									title="Underline"
								>
									U
								</button>
								<button 
									onclick={() => { groupStyle.isStrikethrough = !groupStyle.isStrikethrough; applyGroupStyle(); }} 
									class="px-4 py-2 line-through text-sm bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.isStrikethrough ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
									title="Strikethrough"
								>
									S
								</button>
								<button 
									onclick={() => { 
										groupStyle.textTransform = groupStyle.textTransform === 'uppercase' ? 'none' : 'uppercase'; 
										applyGroupStyle(); 
									}} 
									class="px-3 py-2 uppercase text-xs font-bold bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.textTransform === 'uppercase' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
									title="Uppercase"
								>
									ABC
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
								class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer overflow-hidden" 
							/>
							<span class="text-sm text-gray-600 font-mono">{groupStyle.textColor}</span>
						</div>
					</div>						<!-- Divider -->
						<div class="border-t border-gray-200 my-6"></div>
						
						<!-- Text Card Appearance -->
						<div class="space-y-4">
							
							
							<!-- Enable Background -->
							<div class="flex items-center justify-between">
								<label class="text-sm font-medium text-gray-700">Enable Background</label>
								<button
									onclick={() => { groupStyle.hasBackground = !groupStyle.hasBackground; applyGroupStyle(); }}
									class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {groupStyle.hasBackground ? 'bg-purple-600' : 'bg-gray-300'}"
								>
									<span class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform {groupStyle.hasBackground ? 'translate-x-6' : 'translate-x-1'}"></span>
								</button>
							</div>
							
							{#if groupStyle.hasBackground}
							<!-- Background Color -->
							<div>
								<label class="text-sm font-medium text-gray-700 mb-2 block">Background Color</label>
								<div class="space-y-3">
									<div class="flex gap-3 items-center">
										<input 
											type="color" 
											bind:value={groupStyle.backgroundColor} 
											onchange={applyGroupStyle} 
											class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer overflow-hidden" 
										/>
										<span class="text-sm text-gray-600 font-mono">{groupStyle.backgroundColor}</span>
									</div>
									
									<!-- Color Presets -->
									<div class="flex gap-2 flex-wrap">
										<button onclick={() => { groupStyle.backgroundColor = '#ffffff'; applyGroupStyle(); }} class="w-7 h-7 rounded-full border-2 bg-white hover:scale-110 transition-transform" title="White"></button>
										<button onclick={() => { groupStyle.backgroundColor = '#f3f4f6'; applyGroupStyle(); }} class="w-7 h-7 rounded-full border-2 bg-gray-100 hover:scale-110 transition-transform" title="Light Gray"></button>
										<button onclick={() => { groupStyle.backgroundColor = '#000000'; applyGroupStyle(); }} class="w-7 h-7 rounded-full border-2 bg-black hover:scale-110 transition-transform" title="Black"></button>
										<button onclick={() => { groupStyle.backgroundColor = '#3b82f6'; applyGroupStyle(); }} class="w-7 h-7 rounded-full border-2 bg-blue-500 hover:scale-110 transition-transform" title="Blue"></button>
										<button onclick={() => { groupStyle.backgroundColor = '#8b5cf6'; applyGroupStyle(); }} class="w-7 h-7 rounded-full border-2 bg-purple-500 hover:scale-110 transition-transform" title="Purple"></button>
										<button onclick={() => { groupStyle.backgroundColor = '#10b981'; applyGroupStyle(); }} class="w-7 h-7 rounded-full border-2 bg-green-500 hover:scale-110 transition-transform" title="Green"></button>
										<button onclick={() => { groupStyle.backgroundColor = '#ef4444'; applyGroupStyle(); }} class="w-7 h-7 rounded-full border-2 bg-red-500 hover:scale-110 transition-transform" title="Red"></button>
									</div>
								</div>
							</div>								<!-- Opacity -->
								<div>
									<label class="text-sm font-medium text-gray-700 mb-2 block">Opacity: {groupStyle.backgroundOpacity}%</label>
									<input 
										type="range" 
										min="0" 
										max="100" 
										bind:value={groupStyle.backgroundOpacity} 
										oninput={applyGroupStyleDebounced} 
										class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
									/>
								</div>
								
								<!-- Border Radius -->
								<div>
									<label class="text-sm font-medium text-gray-700 mb-2 block">Border Radius: {groupStyle.borderRadius}px</label>
									<input 
										type="range" 
										min="0" 
										max="32" 
										bind:value={groupStyle.borderRadius} 
										oninput={applyGroupStyleDebounced} 
										class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
									/>
									<div class="flex gap-2 mt-2">
										<button onclick={() => { groupStyle.borderRadius = 0; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Sharp</button>
										<button onclick={() => { groupStyle.borderRadius = 8; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Small</button>
										<button onclick={() => { groupStyle.borderRadius = 16; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Medium</button>
										<button onclick={() => { groupStyle.borderRadius = 24; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Large</button>
										<button onclick={() => { groupStyle.borderRadius = 32; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">XL</button>
									</div>
								</div>
								
								<!-- Padding -->
								<div>
									<label class="text-sm font-medium text-gray-700 mb-2 block">Padding: {groupStyle.padding}px</label>
									<input 
										type="range" 
										min="8" 
										max="24" 
										bind:value={groupStyle.padding} 
										oninput={applyGroupStyleDebounced} 
										class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
									/>
									<div class="flex gap-2 mt-2">
										<button onclick={() => { groupStyle.padding = 8; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Compact</button>
										<button onclick={() => { groupStyle.padding = 16; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Normal</button>
										<button onclick={() => { groupStyle.padding = 24; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors">Relaxed</button>
									</div>
								</div>
								
								<!-- Shadow -->
								<div>
									<label class="text-sm font-medium text-gray-700 mb-2 block">Shadow</label>
									<div class="grid grid-cols-5 gap-2">
										<button 
											onclick={() => { groupStyle.shadow = 'none'; applyGroupStyle(); }} 
											class="px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.shadow === 'none' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
										>
											None
										</button>
										<button 
											onclick={() => { groupStyle.shadow = 'sm'; applyGroupStyle(); }} 
											class="px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors shadow-sm {groupStyle.shadow === 'sm' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
										>
											Soft
										</button>
										<button 
											onclick={() => { groupStyle.shadow = 'md'; applyGroupStyle(); }} 
											class="px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors shadow-md {groupStyle.shadow === 'md' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
										>
											Medium
										</button>
										<button 
											onclick={() => { groupStyle.shadow = 'lg'; applyGroupStyle(); }} 
											class="px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors shadow-lg {groupStyle.shadow === 'lg' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
										>
											Strong
										</button>
										<button 
											onclick={() => { groupStyle.shadow = 'xl'; applyGroupStyle(); }} 
											class="px-3 py-2 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors shadow-xl {groupStyle.shadow === 'xl' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
										>
											XL
										</button>
									</div>
								</div>
								
							{/if}
							
							<!-- Border -->
							<div class="space-y-3">
								<div class="flex items-center justify-between">
									<label class="text-sm font-medium text-gray-700">Enable Border</label>
									<button
										onclick={() => { groupStyle.hasBorder = !groupStyle.hasBorder; applyGroupStyle(); }}
										class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {groupStyle.hasBorder ? 'bg-purple-600' : 'bg-gray-300'}"
									>
										<span class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform {groupStyle.hasBorder ? 'translate-x-6' : 'translate-x-1'}"></span>
									</button>
								</div>
								
								{#if groupStyle.hasBorder}
									<div>
										<label class="text-sm font-medium text-gray-700 mb-2 block">Border Color</label>
										<div class="flex gap-3 items-center">
											<input 
												type="color" 
												bind:value={groupStyle.borderColor} 
												onchange={applyGroupStyle} 
												class="w-10 h-10 rounded-full border-2 border-gray-200 cursor-pointer overflow-hidden" 
											/>
											<span class="text-sm text-gray-600 font-mono">{groupStyle.borderColor}</span>
										</div>
									</div>
									
									<!-- Border Style -->
									<div>
										<label class="text-sm font-medium text-gray-700 mb-2 block">Border Style</label>
										<div class="flex gap-2">
											<button 
												onclick={() => { 
													groupStyle = { ...groupStyle, borderStyle: 'solid' }; 
													applyGroupStyle(); 
												}} 
												class="flex-1 px-3 py-2 text-xs bg-white border-2 rounded-lg hover:bg-gray-50 transition-colors {groupStyle.borderStyle === 'solid' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
											>
												<div class="w-full h-0.5 bg-gray-700"></div>
												<span class="block mt-1">Solid</span>
											</button>
											<button 
												onclick={() => { 
													groupStyle = { ...groupStyle, borderStyle: 'dashed' }; 
													applyGroupStyle(); 
												}} 
												class="flex-1 px-3 py-2 text-xs bg-white border-2 border-dashed rounded-lg hover:bg-gray-50 transition-colors {groupStyle.borderStyle === 'dashed' ? 'ring-2 ring-purple-500 bg-purple-50' : ''}"
											>
												<div class="w-full h-0.5 border-t-2 border-dashed border-gray-700"></div>
												<span class="block mt-1">Dashed</span>
											</button>
										</div>
									</div>
									
									<!-- Border Width -->
									<div>
										<label class="text-sm font-medium text-gray-700 mb-2 block">Border Width: {groupStyle.borderWidth}px</label>
										<div class="flex gap-2">
											<button onclick={() => { groupStyle.borderWidth = 1; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border rounded-lg hover:bg-gray-50 transition-colors {groupStyle.borderWidth === 1 ? 'ring-2 ring-purple-500 bg-purple-50' : ''}">1px</button>
											<button onclick={() => { groupStyle.borderWidth = 2; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border-2 rounded-lg hover:bg-gray-50 transition-colors {groupStyle.borderWidth === 2 ? 'ring-2 ring-purple-500 bg-purple-50' : ''}">2px</button>
											<button onclick={() => { groupStyle.borderWidth = 3; applyGroupStyle(); }} class="px-3 py-1 text-xs bg-white border-4 rounded-lg hover:bg-gray-50 transition-colors {groupStyle.borderWidth === 3 ? 'ring-2 ring-purple-500 bg-purple-50' : ''}">3px</button>
										</div>
									</div>
								{/if}
							</div>
						</div>
				</div>
			{/if}
		</div>
	</div>
{/if}
</div>

<style>
	/* Remove default browser styling from color input to make color fill the entire circle */
	input[type="color"] {
		padding: 0;
		background: none;
	}
	
	input[type="color"]::-webkit-color-swatch-wrapper {
		padding: 0;
		border: none;
	}
	
	input[type="color"]::-webkit-color-swatch {
		border: none;
		border-radius: 50%;
	}
	
	input[type="color"]::-moz-color-swatch {
		border: none;
		border-radius: 50%;
	}
</style>
