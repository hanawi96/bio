<script lang="ts">
	import Avatar from '$lib/components/ui/avatar.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';
	import type { Block } from '$lib/api/blocks';
	import { onMount } from 'svelte';

	export let profile: Partial<Profile> = {};
	export let links: Link[] = [];
	export let blocks: Block[] = [];
	export let showInactive: boolean = true;

	onMount(() => {
		// Carousel scroll functionality works natively with CSS
	});

	// Get link_ids from inactive blocks - these links should be hidden
	$: hiddenLinkIds = new Set(
		(blocks || [])
			.filter(block => !block.is_active && block.link_id)
			.map(block => block.link_id!)
	);

	// Filter out links that are:
	// 1. Hidden by inactive blocks (hiddenLinkIds)
	// 2. Inactive themselves (when showInactive is false)
	$: visibleLinks = (links || []).filter(link => 
		!hiddenLinkIds.has(link.id) && (showInactive || link.is_active)
	);

	// Combine and sort items: pinned links first, then by position
	type Item = { type: 'link'; data: Link; isPinned: boolean } | { type: 'block'; data: Block; isPinned: boolean };
	let items: Item[] = [];
	
	$: {
		console.log('[ProfilePreview] Rebuilding items, blocks:', blocks?.map(b => ({ id: b.id, style: b.style?.substring(0, 50) })));
		items = [
			...visibleLinks.map(link => ({ type: 'link' as const, data: link, position: link.position, isPinned: link.is_pinned || false })),
			...(blocks || []).map(block => ({ type: 'block' as const, data: block, position: block.position, isPinned: false }))
		].sort((a, b) => {
			// Pinned items first
			if (a.isPinned && !b.isPinned) return -1;
			if (!a.isPinned && b.isPinned) return 1;
			// Then by position
			return a.position - b.position;
		});
	}
</script>

<div class="w-full max-w-sm mx-auto">
	<!-- Phone Frame (iPhone Style) -->
	<div class="relative bg-gray-900 rounded-[3rem] p-2 shadow-2xl">
		<!-- Dynamic Island / Notch -->
		<div class="absolute top-0 left-1/2 -translate-x-1/2 w-32 h-7 bg-gray-900 rounded-b-3xl z-10"></div>
		
		<!-- Screen -->
		<div class="relative bg-gradient-to-br from-purple-50 to-blue-50 rounded-[2.5rem] overflow-hidden h-[650px]">
			<div class="h-full overflow-y-auto p-6 space-y-6">
				<!-- Profile Header -->
				<div class="text-center space-y-3">
					<div class="flex justify-center">
						<div class="relative">
							<Avatar 
								src={profile.avatar_url || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + (profile.username || 'user')}
								alt="Avatar"
								class="w-24 h-24 border-4 border-white shadow-lg"
							/>
							<div class="absolute -bottom-1 -right-1 w-6 h-6 bg-green-500 border-2 border-white rounded-full"></div>
						</div>
					</div>
					<div>
						<h2 class="text-xl font-bold text-gray-900">@{profile.username || 'username'}</h2>
						{#if profile.bio}
							<p class="text-sm text-gray-600 mt-1">{profile.bio}</p>
						{/if}
					</div>
				</div>

				<!-- Links & Blocks -->
				<div class="space-y-3">
					{#each items as item (item.type === 'block' ? item.data.id : item.data.id)}
						{#if item.type === 'block' && (showInactive || item.data.is_active)}
							{#if item.data.block_type === 'text' && !item.data.is_group}
								{@const styleConfig = item.data.style ? JSON.parse(item.data.style) : {}}
								<div 
									class="w-full"
									style="
										font-size: {
											styleConfig.fontSize === 'text-small' ? '14px' : 
											styleConfig.fontSize === 'text-medium' ? '16px' : 
											styleConfig.fontSize === 'text-large' ? '18px' :
											styleConfig.fontSize === 'headline-small' ? '20px' :
											styleConfig.fontSize === 'headline-medium' ? '24px' :
											styleConfig.fontSize === 'headline-large' ? '32px' :
											styleConfig.fontSize === 'small' ? '14px' :
											styleConfig.fontSize === 'large' ? '20px' : '16px'
										};
										text-align: {styleConfig.textAlign || 'left'};
										font-weight: {styleConfig.fontSize?.startsWith('headline') || styleConfig.isBold ? 'bold' : 'normal'};
										font-style: {styleConfig.isItalic ? 'italic' : 'normal'};
										text-decoration: {styleConfig.isUnderline && styleConfig.isStrikethrough ? 'underline line-through' : styleConfig.isUnderline ? 'underline' : styleConfig.isStrikethrough ? 'line-through' : 'none'};
										color: {styleConfig.textColor || '#000000'};
									"
								>
									{item.data.content || 'Empty text'}
								</div>
							{:else if item.data.block_type === 'text' && item.data.is_group}
								<!-- Text Group: Render all children with group style -->
								{#if item.data.children && item.data.children.length > 0}
									{@const sortedChildren = item.data.children
										.filter(c => c.is_active)
										.sort((a, b) => a.position - b.position)}
									{@const groupStyle = (() => {
										try {
											if (!item.data.style) {
												console.log('[ProfilePreview] No style for group:', item.data.id);
												return { textAlign: 'left', fontSize: 'text-medium', textColor: '#000000' };
											}
											if (typeof item.data.style === 'string') {
												const parsed = JSON.parse(item.data.style);
												console.log('[ProfilePreview] Parsed group style:', { groupId: item.data.id, style: parsed });
												return parsed;
											}
											console.log('[ProfilePreview] Using object style:', { groupId: item.data.id, style: item.data.style });
											return item.data.style;
										} catch (e) {
											console.error('[ProfilePreview] Failed to parse groupStyle:', item.data.style, e);
											return { textAlign: 'left', fontSize: 'text-medium', textColor: '#000000' };
										}
									})()}
									<div class="space-y-2">
										{#each sortedChildren as child}
											<div class="bg-gradient-to-r from-gray-50 to-gray-100/50 rounded-xl shadow-sm border border-gray-200/60 p-3.5">
												<p 
													class="text-gray-900 break-words"
													style="
														font-size: {groupStyle.fontSize === 'headline-large' ? '1.5rem' : groupStyle.fontSize === 'headline-medium' ? '1.25rem' : groupStyle.fontSize === 'headline-small' ? '1.125rem' : groupStyle.fontSize === 'text-large' ? '1.125rem' : groupStyle.fontSize === 'text-small' ? '0.875rem' : '1rem'};
														text-align: {groupStyle.textAlign || 'left'};
														font-weight: {groupStyle.isBold ? 'bold' : 'normal'};
														font-style: {groupStyle.isItalic ? 'italic' : 'normal'};
														color: {groupStyle.textColor || '#000000'};
													"
												>
													{child.content || 'Empty text'}
												</p>
											</div>
										{/each}
									</div>
								{/if}
							{:else if item.data.block_type === 'divider'}
								{#if item.data.divider_style === 'line'}
									<div class="border-t border-gray-200"></div>
								{:else if item.data.divider_style === 'dots'}
									<div class="flex justify-center gap-1">
										<div class="w-1.5 h-1.5 rounded-full bg-gray-300"></div>
										<div class="w-1.5 h-1.5 rounded-full bg-gray-300"></div>
										<div class="w-1.5 h-1.5 rounded-full bg-gray-300"></div>
									</div>
								{:else}
									<div class="h-4"></div>
								{/if}
							{/if}
						{:else if item.type === 'link' && (showInactive || item.data.is_active)}
							{@const link = item.data}
							{#if link.is_group}
								<!-- Group: Render all children according to layout -->
								{#if link.children && link.children.length > 0}
									{@const sortedChildren = link.children
										.filter(c => c.is_active)
										.sort((a, b) => {
											// Pinned first
											if (a.is_pinned && !b.is_pinned) return -1;
											if (!a.is_pinned && b.is_pinned) return 1;
											// Then by position
											return a.position - b.position;
										})}
									{#if link.group_layout === 'grid'}
										{@const gridCols = link.grid_columns || 2}
										{@const aspectRatio = link.grid_aspect_ratio || '3:2'}
										{@const textSize = link.text_size || 'M'}
										{@const getAspectStyle = (ratio) => {
											const map = { '1:1': '1/1', '3:2': '3/2', '16:9': '16/9', '3:1': '3/1', '2:3': '2/3' };
											return `aspect-ratio: ${map[ratio] || '3/2'}`;
										}}
										<div class="grid gap-3" class:grid-cols-1={gridCols === 1} class:grid-cols-2={gridCols === 2} class:grid-cols-3={gridCols === 3} class:grid-cols-4={gridCols === 4}>
											{#each sortedChildren as child}
												<a
													href={child.url}
													target="_blank"
													rel="noopener noreferrer"
													class="block bg-white hover:bg-gray-50 rounded-xl p-3 transition-all"
													class:shadow-sm={link.show_shadow}
													class:hover:shadow-md={link.show_shadow}
													class:border-2={link.show_outline}
													class:border-gray-200={link.show_outline}
												>
													{#if child.thumbnail_url}
														<img src={child.thumbnail_url} alt={child.title} class="w-full object-cover rounded-lg mb-2" style={getAspectStyle(aspectRatio)} />
													{/if}
													<p class="font-medium text-gray-900 truncate"
														class:text-xs={textSize === 'S'}
														class:text-sm={textSize === 'M'}
														class:text-base={textSize === 'L'}
														class:text-lg={textSize === 'XL'}
														style="text-align: {link.text_alignment || 'center'}"
													>{child.title}</p>
													{#if link.show_description !== false && child.description}
														<p class="text-xs text-gray-500 mt-1 line-clamp-2" style="text-align: {link.text_alignment || 'center'}">{child.description}</p>
													{/if}
												</a>
											{/each}
										</div>
									{:else if link.group_layout === 'carousel'}
										{@const textSize = link.text_size || 'M'}
										{@const aspectRatio = link.grid_aspect_ratio || '3:2'}
										{@const carouselId = `carousel-${link.id}`}
										{@const getAspectStyle = (ratio) => {
											const map = { '1:1': '1/1', '3:2': '3/2', '16:9': '16/9', '3:1': '3/1', '2:3': '2/3' };
											return `aspect-ratio: ${map[ratio] || '3/2'}`;
										}}

										<div class="relative group">
											<!-- Carousel Container -->
											<div id={carouselId} class="overflow-x-scroll snap-x snap-mandatory scrollbar-hide" style="scroll-behavior: smooth;">
												<div class="flex gap-3 px-4">
													{#each sortedChildren as child, idx}
														<a
															href={child.url}
															target="_blank"
															rel="noopener noreferrer"
															class="block bg-white hover:bg-gray-50 rounded-xl p-4 transition-all flex-shrink-0 snap-center w-[85%]"
															class:shadow-sm={link.show_shadow}
															class:hover:shadow-md={link.show_shadow}
															class:border-2={link.show_outline}
															class:border-gray-200={link.show_outline}
														>
															{#if child.thumbnail_url}
																<img src={child.thumbnail_url} alt={child.title} class="w-full object-cover rounded-lg mb-3" style={getAspectStyle(aspectRatio)}/>
															{/if}
															{#if link.show_text !== false}
																<p class="font-medium text-gray-900 mb-1"
																	class:text-xs={textSize === 'S'}
																	class:text-sm={textSize === 'M'}
																	class:text-base={textSize === 'L'}
																	class:text-lg={textSize === 'XL'}
																	style="text-align: {link.text_alignment || 'center'}"
																>{child.title}</p>
																{#if child.description}
																	<p class="text-xs text-gray-500 line-clamp-2" style="text-align: {link.text_alignment || 'center'}">{child.description}</p>
																{/if}
															{/if}
														</a>
													{/each}
												</div>
											</div>

											<!-- Navigation Arrows -->
											{#if sortedChildren.length > 1}
												<button
													onclick={() => {
														const el = document.getElementById(carouselId);
														if (el) el.scrollBy({ left: -el.clientWidth * 0.85, behavior: 'smooth' });
													}}
													class="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-white/90 hover:bg-white rounded-full shadow-lg flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10"
													aria-label="Previous"
												>
													<svg class="w-5 h-5 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
													</svg>
												</button>
												<button
													onclick={() => {
														const el = document.getElementById(carouselId);
														if (el) el.scrollBy({ left: el.clientWidth * 0.85, behavior: 'smooth' });
													}}
													class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-white/90 hover:bg-white rounded-full shadow-lg flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10"
													aria-label="Next"
												>
													<svg class="w-5 h-5 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
													</svg>
												</button>

												<!-- Dots Indicator -->
												<div class="flex justify-center gap-1.5 mt-3">
													{#each sortedChildren as _, idx}
														<button
															onclick={() => {
																const el = document.getElementById(carouselId);
																if (el) {
																	const cardWidth = el.clientWidth * 0.85;
																	el.scrollTo({ left: cardWidth * idx, behavior: 'smooth' });
																}
															}}
															class="w-1.5 h-1.5 rounded-full bg-gray-300 hover:bg-gray-500 transition-colors cursor-pointer"
															aria-label={`Go to slide ${idx + 1}`}
														></button>
													{/each}
												</div>
											{/if}
										</div>
									{:else if link.group_layout === 'card'}
										<!-- Card layout -->
										{@const textSize = link.text_size || 'M'}
										{@const imagePlacement = link.image_placement || 'alternating'}
										<div class="space-y-3">
											{#each sortedChildren as child, index}
												{@const shouldReverse = imagePlacement === 'right' || (imagePlacement === 'alternating' && index % 2 === 0)}
												<a
													href={child.url}
													target="_blank"
													rel="noopener noreferrer"
													class="block bg-white hover:bg-gray-50 rounded-xl overflow-hidden transition-all"
													class:shadow-sm={link.show_shadow}
													class:hover:shadow-md={link.show_shadow}
													class:border-2={link.show_outline}
													class:border-gray-200={link.show_outline}
												>
													<div class="flex items-stretch" class:flex-row-reverse={shouldReverse}>
														{#if child.thumbnail_url}
															<div class="w-1/2 bg-gray-100">
																<img src={child.thumbnail_url} alt={child.title} class="w-full h-full object-cover"/>
															</div>
														{/if}
														<div class="flex-1 p-4 flex flex-col justify-center">
															<p class="font-bold text-gray-900 mb-1"
																class:text-sm={textSize === 'S'}
																class:text-base={textSize === 'M'}
																class:text-lg={textSize === 'L'}
																class:text-xl={textSize === 'XL'}
																style="text-align: {link.text_alignment || 'left'}"
															>{child.title}</p>
															{#if link.show_description !== false && child.description}
																<p class="text-xs text-gray-500" style="text-align: {link.text_alignment || 'left'}">{child.description}</p>
															{/if}
														</div>
													</div>
												</a>
											{/each}
										</div>
									{:else}
										<!-- List layout (default) -->
										{@const textSize = link.text_size || 'M'}
										{@const imageShape = link.image_shape || 'square'}
										<div class="space-y-2">
											{#each sortedChildren as child}
												<a
													href={child.url}
													target="_blank"
													rel="noopener noreferrer"
													class="block bg-white hover:bg-gray-50 rounded-xl p-3 transition-all"
													class:shadow-sm={link.show_shadow}
													class:hover:shadow-md={link.show_shadow}
													class:border-2={link.show_outline}
													class:border-gray-200={link.show_outline}
												>
													<div class="flex items-center gap-3">
														{#if child.thumbnail_url}
															<img src={child.thumbnail_url} alt={child.title} class="w-10 h-10 object-cover flex-shrink-0" class:rounded-lg={imageShape === 'square'} class:rounded-full={imageShape === 'circle'}/>
														{/if}
														<div class="flex-1">
															<p class="font-medium text-gray-900"
																class:text-xs={textSize === 'S'}
																class:text-sm={textSize === 'M'}
																class:text-base={textSize === 'L'}
																class:text-lg={textSize === 'XL'}
																style="text-align: {link.text_alignment || 'left'}"
															>{child.title}</p>
															{#if link.show_description !== false && child.description}
																<p class="text-xs text-gray-500 mt-0.5" style="text-align: {link.text_alignment || 'left'}">{child.description}</p>
															{/if}
														</div>
													</div>
												</a>
											{/each}
										</div>
									{/if}
								{/if}
							{:else if !link.is_group && link.layout_type === 'featured'}
							<!-- Featured Layout (Non-group links only) -->
							<a
								href={link.url}
								target="_blank"
								rel="noopener noreferrer"
								class="block w-full bg-white hover:bg-gray-50 rounded-2xl overflow-hidden shadow-md hover:shadow-lg transition-all duration-200 {!link.is_active ? 'opacity-50' : ''}"
							>
								{#if link.thumbnail_url}
									<div class="relative h-32 bg-gradient-to-br from-indigo-100 to-blue-100">
										<img 
											src={link.thumbnail_url} 
											alt={link.title}
											class="w-full h-full object-cover"
										/>
										<div class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent"></div>
										<div class="absolute bottom-0 left-0 right-0 p-3">
											<span class="font-bold text-white text-sm drop-shadow-lg">{link.title}</span>
										</div>
									</div>
								{:else}
									<div class="relative h-32 bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center">
										<div class="absolute bottom-0 left-0 right-0 p-3">
											<span class="font-bold text-white text-sm drop-shadow-lg">{link.title}</span>
										</div>
									</div>
								{/if}
							</a>
							{:else if !link.is_group}
								<!-- Classic Layout (Non-group links only) -->
								<a
									href={link.url}
									target="_blank"
									rel="noopener noreferrer"
									class="block w-full bg-white hover:bg-gray-50 rounded-2xl p-4 shadow-sm hover:shadow-md transition-all duration-200 {!link.is_active ? 'opacity-50' : ''}"
								>
									<div class="flex items-center gap-3">
										{#if link.thumbnail_url}
											<div class="w-10 h-10 rounded-lg overflow-hidden flex-shrink-0 bg-gray-100">
												<img 
													src={link.thumbnail_url} 
													alt={link.title}
													class="w-full h-full object-cover"
												/>
											</div>
										{/if}
										<span class="font-medium text-gray-900 flex-1">{link.title}</span>
									</div>
								</a>
							{/if}
						{/if}
					{/each}
					
					{#if items.length === 0}
						<div class="text-center py-8 text-gray-400">
							<svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
							</svg>
							<p class="text-sm">No content yet</p>
						</div>
					{/if}
				</div>

				<!-- Footer -->
				<div class="text-center pt-4">
					<Badge variant="secondary" class="text-xs">
						Made with LinkBio
					</Badge>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.scrollbar-hide {
		-ms-overflow-style: none;
		scrollbar-width: none;
		-webkit-overflow-scrolling: touch;
	}
	.scrollbar-hide::-webkit-scrollbar {
		display: none;
	}
</style>
