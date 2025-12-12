<script lang="ts">
	import Avatar from '$lib/components/ui/avatar.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import type { Profile } from '$lib/api/profile';
	import type { Link } from '$lib/api/links';
	import type { Block } from '$lib/api/blocks';
	import { onMount, onDestroy } from 'svelte';

	export let profile: Partial<Profile> = {};
	export let links: Link[] = [];
	export let blocks: Block[] = [];
	export let showInactive: boolean = true;
	
	// Helper function to parse padding from style
	function getPaddingStyle(style: string | null | undefined): string {
		try {
			if (!style) return '16px';
			const parsed = JSON.parse(style);
			const p = parsed.padding;
			if (!p) return '16px';
			if (typeof p === 'number') return `${p}px`;
			return `${p.top || 16}px ${p.right || 16}px ${p.bottom || 16}px ${p.left || 16}px`;
		} catch {
			return '16px';
		}
	}
	
	// Default style object - reuse to avoid object creation
	const DEFAULT_STYLE = {
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

	// Memoize parsed styles to avoid re-parsing
	const styleCache = new Map<string, any>();
	
	// Helper to parse group style with caching
	function parseGroupStyle(styleData: any) {
		if (!styleData) return DEFAULT_STYLE;
		
		// Use cache for string styles
		if (typeof styleData === 'string') {
			if (styleCache.has(styleData)) {
				return styleCache.get(styleData);
			}
			
			try {
				const parsed = JSON.parse(styleData);
				const result = {
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
					backgroundOpacity: parsed.backgroundOpacity ?? 90,
					borderRadius: parsed.borderRadius ?? 12,
					padding: parsed.padding || 16,
					shadow: parsed.shadow || 'none',
					hasBorder: parsed.hasBorder || false,
					borderColor: parsed.borderColor || '#e5e7eb',
					borderWidth: parsed.borderWidth || 1,
					borderStyle: parsed.borderStyle || 'solid'
				};
				
				// Cache the result
				styleCache.set(styleData, result);
				return result;
			} catch (e) {
				return DEFAULT_STYLE;
			}
		}
		
		// For object styles, return directly with defaults
		return {
			textAlign: styleData.textAlign || 'left',
			fontSize: styleData.fontSize || 'text-medium',
			textColor: styleData.textColor || '#000000',
			isBold: styleData.isBold || false,
			isItalic: styleData.isItalic || false,
			isUnderline: styleData.isUnderline || false,
			isStrikethrough: styleData.isStrikethrough || false,
			textTransform: styleData.textTransform || 'none',
			hasBackground: styleData.hasBackground || false,
			backgroundColor: styleData.backgroundColor || '#ffffff',
			backgroundOpacity: styleData.backgroundOpacity ?? 90,
			borderRadius: styleData.borderRadius ?? 12,
			padding: styleData.padding || 16,
			shadow: styleData.shadow || 'none',
			hasBorder: styleData.hasBorder || false,
			borderColor: styleData.borderColor || '#e5e7eb',
			borderWidth: styleData.borderWidth || 1,
			borderStyle: styleData.borderStyle || 'solid'
		};
	}

	onMount(() => {
		// Carousel scroll functionality works natively with CSS
	});

	onDestroy(() => {
		// Clear cache to prevent memory leaks
		styleCache.clear();
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
		items = [
			...visibleLinks.map(link => ({ type: 'link' as const, data: link, position: link.position, isPinned: link.is_pinned || false })),
			...(blocks || []).map(block => ({ type: 'block' as const, data: block, position: block.position, isPinned: false }))
		].sort((a, b) => {
			if (a.isPinned && !b.isPinned) return -1;
			if (!a.isPinned && b.isPinned) return 1;
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
			<!-- Fixed Action Buttons -->
			<div class="absolute top-4 left-4 right-4 flex justify-between items-center z-20">
				<button 
					class="flex items-center gap-2 px-5 py-2.5 bg-white/90 backdrop-blur-sm hover:bg-white rounded-full shadow-lg transition-all hover:shadow-xl border border-gray-200"
					onclick={() => alert('Subscribe functionality')}
				>
					<svg class="w-5 h-5 text-gray-700" fill="currentColor" viewBox="0 0 20 20">
						<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
					</svg>
					<span class="text-sm font-semibold text-gray-800">Subscribe</span>
				</button>
				
				<button 
					class="flex items-center justify-center w-10 h-10 bg-white/90 backdrop-blur-sm hover:bg-white rounded-full shadow-lg transition-all hover:shadow-xl border border-gray-200 mr-1"
					onclick={() => alert('Share functionality')}
					title="Share"
				>
					<svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
					</svg>
				</button>
			</div>
			
			<div class="h-full overflow-y-auto p-6 space-y-6 pt-20">
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
					{#each items as item (item.data.id)}
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
									{@const groupStyle = parseGroupStyle(item.data.style)}
									
									<!-- Pre-compute all style values once -->
									{@const fontSize = groupStyle.fontSize === 'headline-large' ? '1.5rem' : 
										groupStyle.fontSize === 'headline-medium' ? '1.25rem' : 
										groupStyle.fontSize === 'headline-small' ? '1.125rem' : 
										groupStyle.fontSize === 'text-large' ? '1.125rem' : 
										groupStyle.fontSize === 'text-small' ? '0.875rem' : '1rem'}
					{@const textAlign = groupStyle.textAlign || 'left'}
					{@const fontWeight = groupStyle.isBold ? 'bold' : 'normal'}
					{@const fontStyle = groupStyle.isItalic ? 'italic' : 'normal'}
					{@const textDecoration = groupStyle.isUnderline && groupStyle.isStrikethrough ? 'underline line-through' : groupStyle.isUnderline ? 'underline' : groupStyle.isStrikethrough ? 'line-through' : 'none'}
					{@const textTransform = groupStyle.textTransform || 'none'}
					{@const textColor = groupStyle.textColor || '#000000'}								<!-- Pre-compute background styles if enabled -->
								{@const hasCustomStyle = groupStyle.hasBackground || groupStyle.hasBorder}
								{@const bgStyles = hasCustomStyle ? {
									backgroundColor: groupStyle.hasBackground ? `rgba(${parseInt(groupStyle.backgroundColor?.slice(1,3) || 'ff', 16)}, ${parseInt(groupStyle.backgroundColor?.slice(3,5) || 'ff', 16)}, ${parseInt(groupStyle.backgroundColor?.slice(5,7) || 'ff', 16)}, ${(groupStyle.backgroundOpacity ?? 90) / 100})` : 'transparent',
									borderRadius: `${groupStyle.borderRadius ?? 12}px`,
									padding: `${groupStyle.padding ?? 16}px`,
									border: groupStyle.hasBorder ? `${groupStyle.borderWidth ?? 1}px ${groupStyle.borderStyle || 'solid'} ${groupStyle.borderColor || '#e5e7eb'}` : ''
								} : null}									<div class="space-y-2">
										{#each sortedChildren as child (child.id)}
											{#if bgStyles}
												<div 
													class="rounded-xl"
													class:shadow-sm={groupStyle.shadow === 'sm'}
													class:shadow-md={groupStyle.shadow === 'md'}
													class:shadow-lg={groupStyle.shadow === 'lg'}
													class:shadow-xl={groupStyle.shadow === 'xl'}
													style="background-color: {bgStyles.backgroundColor}; border-radius: {bgStyles.borderRadius}; padding: {bgStyles.padding}; {bgStyles.border ? `border: ${bgStyles.border};` : ''}"
												>
													<p 
														class="text-gray-900 break-words"
														style="font-size: {fontSize}; text-align: {textAlign}; font-weight: {fontWeight}; font-style: {fontStyle}; text-decoration: {textDecoration}; text-transform: {textTransform}; color: {textColor};"
													>
														{child.content || 'Empty text'}
													</p>
												</div>
											{:else}
												<p 
													class="text-gray-900 break-words"
													style="font-size: {fontSize}; text-align: {textAlign}; font-weight: {fontWeight}; font-style: {fontStyle}; text-decoration: {textDecoration}; text-transform: {textTransform}; color: {textColor};"
												>
													{child.content || 'Empty text'}
												</p>
											{/if}
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
													class="block bg-white hover:bg-gray-50 rounded-xl transition-all"
													class:shadow-sm={link.show_shadow}
													class:hover:shadow-md={link.show_shadow}
													class:border-2={link.show_outline}
													class:border-gray-200={link.show_outline}
													style="padding: {getPaddingStyle(link.style)}"
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
														{@const hasCustomBg = link.has_card_background !== undefined ? link.has_card_background : true}
														{@const bgColor = link.card_background_color || '#ffffff'}
														{@const bgOpacity = link.card_background_opacity !== undefined ? link.card_background_opacity : 100}
														{@const borderRadius = link.card_border_radius !== undefined ? link.card_border_radius : 12}
														{@const textColor = link.card_text_color || '#000000'}
														{@const shadowX = link.shadow_x || 0}
														{@const shadowY = link.shadow_y || 4}
														{@const shadowBlur = link.shadow_blur || 10}
														{@const r = parseInt(bgColor.slice(1,3), 16)}
														{@const g = parseInt(bgColor.slice(3,5), 16)}
														{@const b = parseInt(bgColor.slice(5,7), 16)}
														{@const shadowStyle = link.show_shadow ? `box-shadow: ${shadowX}px ${shadowY}px ${shadowBlur}px rgba(0,0,0,0.2);` : ''}
														{@const borderStyle = link.has_card_border ? `border: ${link.card_border_width || 1}px ${link.card_border_style || 'solid'} ${link.card_border_color || '#e5e7eb'};` : ''}
														{@const bgStyle = hasCustomBg ? `background-color: rgba(${r}, ${g}, ${b}, ${bgOpacity / 100}); border-radius: ${borderRadius}px; padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}` : `padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}`}
														<a
															href={child.url}
															target="_blank"
															rel="noopener noreferrer"
															class="block hover:bg-gray-50 transition-all flex-shrink-0 snap-center w-[85%]"
															class:bg-white={!hasCustomBg}
															class:rounded-xl={!hasCustomBg}
															class:border-2={link.show_outline && !link.has_card_border}
															class:border-gray-200={link.show_outline && !link.has_card_border}
															style={bgStyle}
														>
															{#if child.thumbnail_url}
																<img src={child.thumbnail_url} alt={child.title} class="w-full object-cover rounded-lg mb-3" style={getAspectStyle(aspectRatio)}/>
															{/if}
															{#if link.show_text !== false}
																<p class="font-medium mb-1"
																	class:text-xs={textSize === 'S'}
																	class:text-sm={textSize === 'M'}
																	class:text-base={textSize === 'L'}
																	class:text-lg={textSize === 'XL'}
																	style="text-align: {link.text_alignment || 'center'}; color: {textColor};"
																>{child.title}</p>
																{#if child.description}
																	<p class="text-xs line-clamp-2" style="text-align: {link.text_alignment || 'center'}; color: {textColor}; opacity: 0.7;">{child.description}</p>
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
												{@const hasCustomBg = link.has_card_background !== undefined ? link.has_card_background : true}
												{@const bgColor = link.card_background_color || '#ffffff'}
												{@const bgOpacity = link.card_background_opacity !== undefined ? link.card_background_opacity : 100}
												{@const borderRadius = link.card_border_radius !== undefined ? link.card_border_radius : 12}
												{@const textColor = link.card_text_color || '#000000'}
												{@const shadowX = link.shadow_x || 0}
												{@const shadowY = link.shadow_y || 4}
												{@const shadowBlur = link.shadow_blur || 10}
												{@const r = parseInt(bgColor.slice(1,3), 16)}
												{@const g = parseInt(bgColor.slice(3,5), 16)}
												{@const b = parseInt(bgColor.slice(5,7), 16)}
												{@const shadowStyle = link.show_shadow ? `box-shadow: ${shadowX}px ${shadowY}px ${shadowBlur}px rgba(0,0,0,0.2);` : ''}
												{@const borderStyle = link.has_card_border ? `border: ${link.card_border_width || 1}px ${link.card_border_style || 'solid'} ${link.card_border_color || '#e5e7eb'};` : ''}
												{@const bgStyle = hasCustomBg ? `background-color: rgba(${r}, ${g}, ${b}, ${bgOpacity / 100}); border-radius: ${borderRadius}px; ${shadowStyle} ${borderStyle}` : `${shadowStyle} ${borderStyle}`}
												<a
													href={child.url}
													target="_blank"
													rel="noopener noreferrer"
													class="block hover:bg-gray-50 overflow-hidden transition-all"
													class:bg-white={!hasCustomBg}
													class:rounded-xl={!hasCustomBg}
													class:border-2={link.show_outline && !link.has_card_border}
													class:border-gray-200={link.show_outline && !link.has_card_border}
													style={bgStyle}
												>
													<div class="flex items-stretch" class:flex-row-reverse={shouldReverse}>
														{#if child.thumbnail_url}
															<div class="w-1/2 bg-gray-100">
																<img src={child.thumbnail_url} alt={child.title} class="w-full h-full object-cover"/>
															</div>
														{/if}
														<div class="flex-1 flex flex-col justify-center" style="padding: {getPaddingStyle(link.style)}">
															<p class="font-bold mb-1"
																class:text-sm={textSize === 'S'}
																class:text-base={textSize === 'M'}
																class:text-lg={textSize === 'L'}
																class:text-xl={textSize === 'XL'}
																style="text-align: {link.text_alignment || 'left'}; color: {textColor};"
															>{child.title}</p>
															{#if link.show_description !== false && child.description}
																<p class="text-xs" style="text-align: {link.text_alignment || 'left'}; color: {textColor}; opacity: 0.7;">{child.description}</p>
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
												{@const hasCustomBg = link.has_card_background !== undefined ? link.has_card_background : true}
												{@const bgColor = link.card_background_color || '#ffffff'}
												{@const bgOpacity = link.card_background_opacity !== undefined ? link.card_background_opacity : 100}
												{@const borderRadius = link.card_border_radius !== undefined ? link.card_border_radius : 12}
												{@const textColor = link.card_text_color || '#000000'}
												{@const shadowX = link.shadow_x || 0}
												{@const shadowY = link.shadow_y || 4}
												{@const shadowBlur = link.shadow_blur || 10}
												{@const r = parseInt(bgColor.slice(1,3), 16)}
												{@const g = parseInt(bgColor.slice(3,5), 16)}
												{@const b = parseInt(bgColor.slice(5,7), 16)}
												{@const shadowStyle = link.show_shadow ? `box-shadow: ${shadowX}px ${shadowY}px ${shadowBlur}px rgba(0,0,0,0.2);` : ''}
												{@const borderStyle = link.has_card_border ? `border: ${link.card_border_width || 1}px ${link.card_border_style || 'solid'} ${link.card_border_color || '#e5e7eb'};` : ''}
												{@const bgStyle = hasCustomBg ? `background-color: rgba(${r}, ${g}, ${b}, ${bgOpacity / 100}); border-radius: ${borderRadius}px; padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}` : `padding: ${getPaddingStyle(link.style)}; ${shadowStyle} ${borderStyle}`}

												<a
													href={child.url}
													target="_blank"
													rel="noopener noreferrer"
													class="block hover:bg-gray-50 transition-all"
													class:bg-white={!hasCustomBg}
													class:rounded-xl={!hasCustomBg}
													class:border-2={link.show_outline && !link.has_card_border}
													class:border-gray-200={link.show_outline && !link.has_card_border}
													style={bgStyle}
												>
													<div class="flex items-center gap-3">
														{#if child.thumbnail_url}
															<img src={child.thumbnail_url} alt={child.title} class="w-10 h-10 object-cover flex-shrink-0" class:rounded-lg={imageShape === 'square'} class:rounded-full={imageShape === 'circle'} class:rounded-none={imageShape === 'sharp'}/>
														{/if}
														<div class="flex-1">
															<p class="font-medium"
																class:text-xs={textSize === 'S'}
																class:text-sm={textSize === 'M'}
																class:text-base={textSize === 'L'}
																class:text-lg={textSize === 'XL'}
																style="text-align: {link.text_alignment || 'left'}; color: {textColor};"
															>{child.title}</p>
															{#if link.show_description !== false && child.description}
																<p class="text-xs mt-0.5" style="text-align: {link.text_alignment || 'left'}; color: {textColor}; opacity: 0.7;">{child.description}</p>
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
